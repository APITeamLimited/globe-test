package orchestrator

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/aggregator"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/proto"
)

var TEST_INFO_KEYS = []string{"INTERVAL", "CONSOLE"}

// Cleans up the worker and orchestrator clients, storing all results in storeMongo
func cleanup(gs libOrch.BaseGlobalState, job libOrch.Job, childJobs map[string]libOrch.ChildJobDistribution, storeMongoDB *mongo.Database,
	scope libOrch.Scope, testInfoStoreReceipt *primitive.ObjectID) error {
	// Store results in MongoDB
	bucketName := fmt.Sprintf("%s:%s", scope.Variant, scope.VariantTargetId)
	var err error
	var jobBucket *gridfs.Bucket
	if gs.Standalone() {
		jobBucket, err = gridfs.NewBucket(storeMongoDB, options.GridFSBucket().SetName(bucketName))
		if err != nil {
			return err
		}
	}

	updatesKey := fmt.Sprintf("%s:updates", job.Id)

	unparsedMessages, err := gs.OrchestratorClient().SMembers(gs.Ctx(), updatesKey).Result()
	if err != nil {
		return err
	}

	var testData []libOrch.OrchestratorOrWorkerMessage

	for _, value := range unparsedMessages {
		// Declare here else fields will be inherited from previous iteration
		var message libOrch.OrchestratorOrWorkerMessage

		err := json.Unmarshal([]byte(value), &message)
		if err != nil {
			fmt.Println("error unmarshalling message: when cleaning up job", job.Id, err.Error())
			continue
		}

		if lib.StringInSlice(TEST_INFO_KEYS, message.MessageType) {
			testData = append(testData, message)
		}
	}

	postTestInfo, err := aggregator.DeterminePostTestInfo(gs, &testData)
	if err != nil {
		return fmt.Errorf("error determining post test info: %s", err.Error())
	}

	encodedBytes, err := proto.Marshal(postTestInfo)
	if err != nil {
		return fmt.Errorf("error marshalling post test info: %s", err.Error())
	}

	filename := fmt.Sprintf("GlobeTest:%s:postTestInfo", job.Id)

	if gs.Standalone() {
		err = libOrch.SetInBucket(jobBucket, filename, encodedBytes, "application/protobuf", testInfoStoreReceipt)
		if err != nil {
			// Can't alert client here, as the client has already been cleaned up
			return fmt.Errorf("error setting logs in bucket: %s", err.Error())
		}
	} else {
		// TODO
		localhostFile := libOrch.LocalhostFile{
			FileName: filename,
			Contents: base64.StdEncoding.EncodeToString(encodedBytes),
			Kind:     "TEST_INFO",
		}

		marshalledLocalhostFile, err := json.Marshal(localhostFile)
		if err != nil {
			// Can't alert client here, as the client has already been cleaned up
			return fmt.Errorf("error setting logs in bucket: %s", err.Error())

		}

		libOrch.DispatchMessage(gs, string(marshalledLocalhostFile), "LOCALHOST_FILE")
	}

	// Clean up orchestrator
	// Set types to expire so lagging users can access environment variables
	gs.OrchestratorClient().Expire(gs.Ctx(), updatesKey, time.Second*10)
	gs.OrchestratorClient().Expire(gs.Ctx(), job.Id, time.Second*10)
	gs.OrchestratorClient().SRem(gs.Ctx(), "orchestrator:executionHistory", job.Id)

	return nil
}
