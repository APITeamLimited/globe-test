package lib

import (
	"context"
	"crypto/tls"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/APITeamLimited/redis/v9"
)

const OUT_OF_CREDITS_MESSAGE = "out of credits"

func GetCreditsClient(standalone bool) *redis.Client ***REMOVED***
	if !standalone ***REMOVED***
		return nil
	***REMOVED***

	enableCreditsSystem := GetEnvVariable("ENABLE_CREDITS_SYSTEM", "false") == "true"

	if !enableCreditsSystem ***REMOVED***
		return nil
	***REMOVED***

	clientHost := GetEnvVariable("CREDITS_REDIS_HOST", "localhost")
	clientPort := GetEnvVariable("CREDITS_REDIS_PORT", "6379")

	isSecure := GetEnvVariable("CREDITS_REDIS_IS_SECURE", "false") == "true"

	options := &redis.Options***REMOVED***
		Addr:     fmt.Sprintf("%s:%s", clientHost, clientPort),
		Username: "default",
		Password: GetEnvVariable("CREDITS_REDIS_PASSWORD", ""),
	***REMOVED***

	if isSecure ***REMOVED***
		clientCert := GetEnvVariable("CREDITS_REDIS_CERT", "")
		clientKey := GetEnvVariable("CREDITS_REDIS_KEY", "")

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil ***REMOVED***
			panic(fmt.Errorf("error loading credits cert: %s", err))
		***REMOVED***

		options.TLSConfig = &tls.Config***REMOVED***
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: GetEnvVariable("CREDITS_REDIS_INSECURE_SKIP_VERIFY", "false") == "true",
			Certificates:       []tls.Certificate***REMOVED***cert***REMOVED***,
		***REMOVED***
	***REMOVED***

	return redis.NewClient(options)
***REMOVED***

type CreditsManager struct ***REMOVED***
	ctx           context.Context
	creditsClient *redis.Client

	workspaceName string

	// Interval timer
	ticker *time.Ticker

	oldCredits   int64
	usedCredits  int64
	creditsMutex sync.Mutex

	lastBillingTime time.Time
***REMOVED***

func CreateCreditsManager(ctx context.Context, variant string, variantTargetId string,
	creditsClient *redis.Client) *CreditsManager ***REMOVED***
	creditsManager := &CreditsManager***REMOVED***
		ctx:           ctx,
		creditsClient: creditsClient,
		workspaceName: fmt.Sprintf("%s:%s", variant, variantTargetId),
		ticker:        time.NewTicker(1 * time.Second),
		creditsMutex:  sync.Mutex***REMOVED******REMOVED***,
	***REMOVED***

	go func() ***REMOVED***
		for range creditsManager.ticker.C ***REMOVED***
			creditsManager.captureCredits()
		***REMOVED***
	***REMOVED***()

	// Perform initial credits capture
	creditsManager.captureCredits()

	return creditsManager
***REMOVED***

func (creditsManager *CreditsManager) GetCredits() int64 ***REMOVED***
	if creditsManager == nil ***REMOVED***
		return math.MaxInt64
	***REMOVED***

	creditsManager.creditsMutex.Lock()
	defer creditsManager.creditsMutex.Unlock()

	return creditsManager.oldCredits - creditsManager.usedCredits
***REMOVED***

func (creditsManager *CreditsManager) StopCreditsCapturing() ***REMOVED***
	if creditsManager == nil ***REMOVED***
		return
	***REMOVED***

	creditsManager.ticker.Stop()

	// Capture credits one last time
	creditsManager.captureCredits()
***REMOVED***

func (creditsManager *CreditsManager) UseCredits(credits int64) bool ***REMOVED***
	if creditsManager == nil ***REMOVED***
		return true
	***REMOVED***

	creditsManager.creditsMutex.Lock()
	defer creditsManager.creditsMutex.Unlock()

	// Check delta is not greater than credits
	if credits+creditsManager.usedCredits > creditsManager.oldCredits ***REMOVED***
		return false
	***REMOVED***

	creditsManager.usedCredits += credits

	return true
***REMOVED***

func (creditsManager *CreditsManager) ForceDeductCredits(credits int64, setLastBillingTime bool) ***REMOVED***
	if creditsManager == nil ***REMOVED***
		return
	***REMOVED***

	creditsManager.creditsMutex.Lock()
	defer creditsManager.creditsMutex.Unlock()

	creditsManager.usedCredits += credits

	// Check delta is not greater than credits
	if credits >= creditsManager.oldCredits ***REMOVED***
		// Just set credits to 0
		creditsManager.usedCredits = creditsManager.oldCredits
	***REMOVED***

	if setLastBillingTime ***REMOVED***
		creditsManager.lastBillingTime = time.Now()
	***REMOVED***
***REMOVED***

// captureCredits captures credits from the credits pool and stores them in the credits store
func (creditsManager *CreditsManager) captureCredits() ***REMOVED***
	if creditsManager == nil ***REMOVED***
		return
	***REMOVED***

	creditsManager.creditsMutex.Lock()
	defer creditsManager.creditsMutex.Unlock()

	// DECR usedCredits from credits pool
	newCredits, err := creditsManager.creditsClient.DecrBy(creditsManager.ctx, creditsManager.workspaceName, creditsManager.usedCredits).Result()
	if err != nil ***REMOVED***
		fmt.Println("Error capturing credits: ", err)
		return
	***REMOVED***

	// If newCredits is negative, set it to 0
	if newCredits < 0 ***REMOVED***
		newCredits = 0
		creditsManager.creditsClient.Set(creditsManager.ctx, creditsManager.workspaceName, strconv.FormatInt(newCredits, 10), 0)
	***REMOVED***

	// Add the credits to the credits store
	creditsManager.oldCredits = newCredits
	creditsManager.usedCredits = 0
***REMOVED***

func (creditsManager *CreditsManager) LastBillingTime() time.Time ***REMOVED***
	if creditsManager == nil ***REMOVED***
		return time.Now()
	***REMOVED***

	return creditsManager.lastBillingTime
***REMOVED***
