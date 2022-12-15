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

func GetCreditsClient(standalone bool) *redis.Client {
	if !standalone {
		return nil
	}

	enableCreditsSystem := GetEnvVariable("ENABLE_CREDITS_SYSTEM", "false") == "true"

	if !enableCreditsSystem {
		return nil
	}

	clientHost := GetEnvVariable("CREDITS_REDIS_HOST", "localhost")
	clientPort := GetEnvVariable("CREDITS_REDIS_PORT", "6379")

	isSecure := GetEnvVariable("CREDITS_REDIS_IS_SECURE", "false") == "true"

	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", clientHost, clientPort),
		Username: "default",
		Password: GetEnvVariable("CREDITS_REDIS_PASSWORD", ""),
	}

	if isSecure {
		clientCert := GetEnvVariable("CREDITS_REDIS_CERT", "")
		clientKey := GetEnvVariable("CREDITS_REDIS_KEY", "")

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			panic(fmt.Errorf("error loading credits cert: %s", err))
		}

		options.TLSConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: GetEnvVariable("CREDITS_REDIS_INSECURE_SKIP_VERIFY", "false") == "true",
			Certificates:       []tls.Certificate{cert},
		}
	}

	return redis.NewClient(options)
}

type CreditsManager struct {
	ctx           context.Context
	creditsClient *redis.Client

	workspaceName string

	// Interval timer
	ticker *time.Ticker

	oldCredits   int64
	usedCredits  int64
	creditsMutex sync.Mutex

	lastBillingTime time.Time
}

func CreateCreditsManager(ctx context.Context, variant string, variantTargetId string,
	creditsClient *redis.Client) *CreditsManager {
	creditsManager := &CreditsManager{
		ctx:           ctx,
		creditsClient: creditsClient,
		workspaceName: fmt.Sprintf("%s:%s", variant, variantTargetId),
		ticker:        time.NewTicker(1 * time.Second),
		creditsMutex:  sync.Mutex{},
	}

	go func() {
		for range creditsManager.ticker.C {
			creditsManager.captureCredits()
		}
	}()

	// Perform initial credits capture
	creditsManager.captureCredits()

	return creditsManager
}

func (creditsManager *CreditsManager) GetCredits() int64 {
	if creditsManager == nil {
		return math.MaxInt64
	}

	creditsManager.creditsMutex.Lock()
	defer creditsManager.creditsMutex.Unlock()

	return creditsManager.oldCredits - creditsManager.usedCredits
}

func (creditsManager *CreditsManager) StopCreditsCapturing() {
	if creditsManager == nil {
		return
	}

	creditsManager.ticker.Stop()

	// Capture credits one last time
	creditsManager.captureCredits()
}

func (creditsManager *CreditsManager) UseCredits(credits int64) bool {
	if creditsManager == nil {
		return true
	}

	creditsManager.creditsMutex.Lock()
	defer creditsManager.creditsMutex.Unlock()

	// Check delta is not greater than credits
	if credits+creditsManager.usedCredits > creditsManager.oldCredits {
		return false
	}

	creditsManager.usedCredits += credits

	return true
}

func (creditsManager *CreditsManager) ForceDeductCredits(credits int64, setLastBillingTime bool) {
	if creditsManager == nil {
		return
	}

	creditsManager.creditsMutex.Lock()
	defer creditsManager.creditsMutex.Unlock()

	creditsManager.usedCredits += credits

	// Check delta is not greater than credits
	if credits >= creditsManager.oldCredits {
		// Just set credits to 0
		creditsManager.usedCredits = creditsManager.oldCredits
	}

	if setLastBillingTime {
		creditsManager.lastBillingTime = time.Now()
	}
}

// captureCredits captures credits from the credits pool and stores them in the credits store
func (creditsManager *CreditsManager) captureCredits() {
	if creditsManager == nil {
		return
	}

	creditsManager.creditsMutex.Lock()
	defer creditsManager.creditsMutex.Unlock()

	// DECR usedCredits from credits pool
	newCredits, err := creditsManager.creditsClient.DecrBy(creditsManager.ctx, creditsManager.workspaceName, creditsManager.usedCredits).Result()
	if err != nil {
		fmt.Println("Error capturing credits: ", err)
		return
	}

	// If newCredits is negative, set it to 0
	if newCredits < 0 {
		newCredits = 0
		creditsManager.creditsClient.Set(creditsManager.ctx, creditsManager.workspaceName, strconv.FormatInt(newCredits, 10), 0)
	}

	// Add the credits to the credits store
	creditsManager.oldCredits = newCredits
	creditsManager.usedCredits = 0
}

func (creditsManager *CreditsManager) LastBillingTime() time.Time {
	if creditsManager == nil {
		return time.Now()
	}

	return creditsManager.lastBillingTime
}
