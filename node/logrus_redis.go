package node

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

// HookConfig stores configuration needed to setup the hook
type HookConfig struct ***REMOVED***
	Key      string
	Format   string
	App      string
	Host     string
	Password string
	Hostname string
	Port     int
	DB       int
	TTL      int
***REMOVED***

// RedisHook to sends logs to Redis server
type RedisHook struct ***REMOVED***
	RedisPool      *redis.Pool
	RedisHost      string
	RedisKey       string
	LogstashFormat string
	AppName        string
	Hostname       string
	RedisPort      int
	TTL            int
	DialOptions    []redis.DialOption
***REMOVED***

// NewHook creates a hook to be added to an instance of logger
func NewHook(config HookConfig, options ...redis.DialOption) (*RedisHook, error) ***REMOVED***
	pool := newRedisConnectionPool(config.Host, config.Password, config.Port, config.DB, options...)

	if config.Format != "v0" && config.Format != "v1" && config.Format != "access" ***REMOVED***
		return nil, fmt.Errorf("unknown message format")
	***REMOVED***

	// test if connection with REDIS can be established
	conn := pool.Get()
	defer conn.Close()

	// check connection
	_, err := conn.Do("PING")
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("unable to connect to REDIS: %s", err)
	***REMOVED***

	return &RedisHook***REMOVED***
		RedisHost:      config.Host,
		RedisPool:      pool,
		RedisKey:       config.Key,
		LogstashFormat: config.Format,
		AppName:        config.App,
		Hostname:       config.Hostname,
		TTL:            config.TTL,
		DialOptions:    options,
	***REMOVED***, nil

***REMOVED***

// Fire is called when a log event is fired.
func (hook *RedisHook) Fire(entry *logrus.Entry) error ***REMOVED***
	var msg interface***REMOVED******REMOVED***

	switch hook.LogstashFormat ***REMOVED***
	case "v0":
		msg = createV0Message(entry, hook.AppName, hook.Hostname)
	case "v1":
		msg = createV1Message(entry, hook.AppName, hook.Hostname)
	case "access":
		msg = createAccessLogMessage(entry, hook.AppName, hook.Hostname)
	default:
		fmt.Println("Invalid LogstashFormat")
	***REMOVED***

	// Marshal into json message
	js, err := json.Marshal(msg)
	if err != nil ***REMOVED***
		return fmt.Errorf("error creating message for REDIS: %s", err)
	***REMOVED***

	// get connection from pool
	conn := hook.RedisPool.Get()
	defer conn.Close()

	// send message
	_, err = conn.Do("RPUSH", hook.RedisKey, js)
	if err != nil ***REMOVED***
		return fmt.Errorf("error sending message to REDIS: %s", err)
	***REMOVED***

	if hook.TTL != 0 ***REMOVED***
		_, err = conn.Do("EXPIRE", hook.RedisKey, hook.TTL)
		if err != nil ***REMOVED***
			return fmt.Errorf("error setting TTL to key: %s, %s", hook.RedisKey, err)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Levels returns the available logging levels.
func (hook *RedisHook) Levels() []logrus.Level ***REMOVED***
	return []logrus.Level***REMOVED***
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	***REMOVED***
***REMOVED***

func createV0Message(entry *logrus.Entry, appName, hostname string) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	m := make(map[string]interface***REMOVED******REMOVED***)
	m["@timestamp"] = entry.Time.UTC().Format(time.RFC3339Nano)
	m["@source_host"] = hostname
	m["@message"] = entry.Message

	fields := make(map[string]interface***REMOVED******REMOVED***)
	fields["level"] = entry.Level.String()
	fields["application"] = appName

	for k, v := range entry.Data ***REMOVED***
		fields[k] = v
	***REMOVED***
	m["@fields"] = fields

	return m
***REMOVED***

func createV1Message(entry *logrus.Entry, appName, hostname string) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	m := make(map[string]interface***REMOVED******REMOVED***)
	m["@timestamp"] = entry.Time.UTC().Format(time.RFC3339Nano)
	m["host"] = hostname
	m["message"] = entry.Message
	m["level"] = entry.Level.String()
	m["application"] = appName
	for k, v := range entry.Data ***REMOVED***
		m[k] = v
	***REMOVED***

	return m
***REMOVED***

func createAccessLogMessage(entry *logrus.Entry, appName, hostname string) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	m := make(map[string]interface***REMOVED******REMOVED***)
	m["message"] = entry.Message
	m["@source_host"] = hostname

	fields := make(map[string]interface***REMOVED******REMOVED***)
	fields["application"] = appName

	for k, v := range entry.Data ***REMOVED***
		fields[k] = v
	***REMOVED***
	m["@fields"] = fields

	return m
***REMOVED***

func newRedisConnectionPool(server, password string, port int, db int, options ...redis.DialOption) *redis.Pool ***REMOVED***
	hostPort := fmt.Sprintf("%s:%d", server, port)
	return &redis.Pool***REMOVED***
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) ***REMOVED***
			dialOptions := append([]redis.DialOption***REMOVED***
				redis.DialDatabase(db),
				redis.DialPassword(password),
			***REMOVED***, options...)
			c, err := redis.Dial("tcp", hostPort, dialOptions...)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return c, err
		***REMOVED***,
		TestOnBorrow: func(c redis.Conn, t time.Time) error ***REMOVED***
			_, err := c.Do("PING")
			return err
		***REMOVED***,
	***REMOVED***
***REMOVED***
