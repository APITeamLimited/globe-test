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

const (
	OUT_OF_CREDITS_MESSAGE = "out of credits"
	creditsCapturePeriod   = 5 * time.Second
)

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

		DialTimeout:  20 * time.Second,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,

		MaxRetries:   20,
		MinIdleConns: 5,
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

	funcModeInfo FuncModeInfo

	freeCreditsName string
	paidCreditsName string

	// Interval timer
	captureTicker *time.Ticker

	// Prevent multiple captures from running at the same time
	isCapturingMutex sync.Mutex

	oldCredits   int64
	usedCredits  int64
	creditsMutex sync.Mutex

	lastBillingTime time.Time

	billingTicker *time.Ticker
}

func CreateCreditsManager(ctx context.Context, variant string, variantTargetId string,
	creditsClient *redis.Client, funcModeInfo FuncModeInfo) *CreditsManager {
	creditsManager := &CreditsManager{
		ctx:              ctx,
		creditsClient:    creditsClient,
		freeCreditsName:  fmt.Sprintf("%s:%s:freeCredits", variant, variantTargetId),
		paidCreditsName:  fmt.Sprintf("%s:%s:paidCredits", variant, variantTargetId),
		captureTicker:    time.NewTicker(creditsCapturePeriod),
		creditsMutex:     sync.Mutex{},
		isCapturingMutex: sync.Mutex{},
		funcModeInfo:     funcModeInfo,
		lastBillingTime:  time.Now(),
	}

	go func() {
		for range creditsManager.captureTicker.C {
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

	creditsManager.captureTicker.Stop()

	// Capture credits one last time
	creditsManager.captureCredits()
}

func (creditsManager *CreditsManager) ForceDeductCredits(credits int64, setLastBillingTime bool) {
	if creditsManager == nil {
		return
	}

	creditsManager.creditsMutex.Lock()
	defer creditsManager.creditsMutex.Unlock()

	creditsManager.usedCredits += credits

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

		// Add newFreeCredits toback to redis, like setting to 0 but this way we don't have to lock
		newFreeCredits = creditsManager.creditsClient.IncrBy(creditsManager.ctx, creditsManager.freeCreditsName, -newFreeCredits).Val()
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
			// Add newPaidCredits toback to redis, like setting to 0 but this way we don't have to lock
			newPaidCredits = creditsManager.creditsClient.IncrBy(creditsManager.ctx, creditsManager.paidCreditsName, -newPaidCredits).Val()
		}
	}

	// DECR usedCredits from credits pool
	newCredits := newFreeCredits + newPaidCredits

	// Add the credits to the credits store
	creditsManager.oldCredits = newCredits
	creditsManager.usedCredits = 0
}

func (creditsManager *CreditsManager) BillFinalCredits() {
	if creditsManager == nil {
		return
	}

	timeSinceLastBilling := time.Since(creditsManager.lastBillingTime)
	billingCycleCount := int64(math.Ceil(float64(timeSinceLastBilling.Milliseconds()) / 100))

	if billingCycleCount <= 0 {
		billingCycleCount = 1
	}

	fractionCost := billingCycleCount * creditsManager.funcModeInfo.Instance100MSUnitRate

	creditsManager.ForceDeductCredits(fractionCost, false)

	creditsManager.billingTicker.Stop()

	creditsManager.StopCreditsCapturing()
}

func (creditsManager *CreditsManager) StartMonitoringCredits(outOfCreditsCallback func()) {
	if creditsManager == nil {
		return
	}

	// Every second check if we have enough credits to continue
	creditsManager.billingTicker = time.NewTicker(100 * time.Millisecond)

	go func() {
		for range creditsManager.billingTicker.C {
			if creditsManager.GetCredits() < creditsManager.funcModeInfo.Instance100MSUnitRate {
				outOfCreditsCallback()
				return
			}

			creditsManager.ForceDeductCredits(creditsManager.funcModeInfo.Instance100MSUnitRate, true)
		}
	}()
}
