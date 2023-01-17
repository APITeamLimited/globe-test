package lib

import (
	"context"
	"crypto/tls"
	"crypto/x509"
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

	enableCreditsSystem := GetEnvVariableBool("ENABLE_CREDITS_SYSTEM", false)

	if !enableCreditsSystem {
		return nil
	}

	clientHost := GetEnvVariable("CREDITS_REDIS_HOST", "localhost")
	clientPort := GetEnvVariable("CREDITS_REDIS_PORT", "6379")

	isSecure := GetEnvVariableBool("CREDITS_REDIS_IS_SECURE", false)

	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", clientHost, clientPort),
		Username: "default",
		Password: GetEnvVariable("CREDITS_REDIS_PASSWORD", ""),
	}

	if isSecure {
		clientCert := GetHexEnvVariable("CREDITS_REDIS_CERT_HEX", "")
		clientKey := GetHexEnvVariable("CREDITS_REDIS_KEY_HEX", "")

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			panic(fmt.Errorf("error loading credits cert: %s", err))
		}

		// Load CA cert
		caCertPool := x509.NewCertPool()
		caCert := GetHexEnvVariable("CREDITS_REDIS_CA_CERT_HEX", "")
		ok := caCertPool.AppendCertsFromPEM([]byte(caCert))
		if !ok {
			panic(fmt.Errorf("failed to parse root certificate"))
		}

		options.TLSConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: GetEnvVariableBool("CREDITS_REDIS_INSECURE_SKIP_VERIFY", false),
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
		}
	}

	return redis.NewClient(options)
}

type CreditsManager struct {
	ctx           context.Context
	creditsClient *redis.Client

	freeCreditsName string
	paidCreditsName string

	// Interval timer
	ticker *time.Ticker

	oldCredits   int64
	usedCredits  int64
	creditsMutex sync.Mutex

	lastBillingTime time.Time

	// Prevent multiple captures from running at the same time
	isCapturingMutex sync.Mutex
}

func CreateCreditsManager(ctx context.Context, variant string, variantTargetId string,
	creditsClient *redis.Client) *CreditsManager {
	creditsManager := &CreditsManager{
		ctx:              ctx,
		creditsClient:    creditsClient,
		freeCreditsName:  fmt.Sprintf("%s:%s:freeCredits", variant, variantTargetId),
		paidCreditsName:  fmt.Sprintf("%s:%s:paidCredits", variant, variantTargetId),
		ticker:           time.NewTicker(1 * time.Second),
		creditsMutex:     sync.Mutex{},
		isCapturingMutex: sync.Mutex{},
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

	creditsManager.isCapturingMutex.Lock()
	defer creditsManager.isCapturingMutex.Unlock()

	creditsManager.creditsMutex.Lock()
	defer creditsManager.creditsMutex.Unlock()

	newFreeCredits, err := creditsManager.creditsClient.DecrBy(creditsManager.ctx, creditsManager.freeCreditsName, creditsManager.usedCredits).Result()
	if err != nil {
		fmt.Println("Error capturing credits: ", err)
		return
	}

	newPaidCredits := int64(0)

	// If newFreeCredits is negative, deduct from paidCredits
	if newFreeCredits < 0 {
		newPaidCredits, err = creditsManager.creditsClient.DecrBy(creditsManager.ctx, creditsManager.paidCreditsName, -newFreeCredits).Result()
		if err != nil {
			fmt.Println("Error capturing credits: ", err)
			return
		}

		// Set newFreeCredits to 0
		creditsManager.creditsClient.Set(creditsManager.ctx, creditsManager.freeCreditsName, 0, 0)
		newFreeCredits = 0
	} else {
		newPaidCreditsStr, err := creditsManager.creditsClient.Get(creditsManager.ctx, creditsManager.paidCreditsName).Result()
		// Nil error can occur here if user hasn't purchased any paid credits
		if err != nil && err != redis.Nil {
			fmt.Println("Error capturing credits: ", err)
			return
		} else if err == redis.Nil {
			newPaidCreditsStr = "0"
		}

		newPaidCredits, err = strconv.ParseInt(newPaidCreditsStr, 10, 64)
		if err != nil {
			fmt.Println("Error capturing credits: ", err)
			return
		}

		if newPaidCredits < 0 {
			creditsManager.creditsClient.Set(creditsManager.ctx, creditsManager.paidCreditsName, "0", 0)
			newPaidCredits = 0
		}
	}

	// DECR usedCredits from credits pool
	newCredits := newFreeCredits + newPaidCredits

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
