/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package worker

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const (
	defaultConfigFileName   = "config.json"
	waitRemoteLoggerTimeout = time.Second * 5
)

// globalFlags contains global config values that apply for all k6 sub-commands.
type globalFlags struct ***REMOVED***
	configFilePath string
	noColor        bool
	address        string
	logOutput      string
	logFormat      string
***REMOVED***

// globalState contains the globalFlags and accessors for most of the global
// process-external state like CLI arguments, env vars, standard input, output
// and error, etc. In practice, most of it is normally accessed through the `os`
// package from the Go stdlib.
//
// We group them here so we can prevent direct access to them from the rest of
// the k6 codebase. This gives us the ability to mock them and have robust and
// easy-to-write integration-like tests to check the k6 end-to-end behavior in
// any simulated conditions.
//
// `newGlobalState()` returns a globalState object with the real `os`
// parameters, while `newGlobalTestState()` can be used in tests to create
// simulated environments.
type globalState struct ***REMOVED***
	ctx context.Context

	fs      afero.Fs
	getwd   func() (string, error)
	args    []string
	envVars map[string]string

	defaultFlags, flags globalFlags

	stdOut, stdErr *consoleWriter
	stdIn          io.Reader

	osExit       func(int)
	signalNotify func(chan<- os.Signal, ...os.Signal)
	signalStop   func(chan<- os.Signal)

	logger         *logrus.Logger
	fallbackLogger logrus.FieldLogger
***REMOVED***

// Ideally, this should be the only function in the whole codebase where we use
// global variables and functions from the os package. Anywhere else, things
// like os.Stdout, os.Stderr, os.Stdin, os.Getenv(), etc. should be removed and
// the respective properties of globalState used instead.
func newGlobalState(ctx context.Context, client *redis.Client, jobId string, workerId string) *globalState ***REMOVED***
	stdOut := &consoleWriter***REMOVED***ctx, client, jobId, workerId***REMOVED***
	stdErr := &consoleWriter***REMOVED***ctx, client, jobId, workerId***REMOVED***

	envVars := buildEnvMap(os.Environ())

	hook, err := NewRedisHook(client, ctx, jobId, workerId)
	if err != nil ***REMOVED***
		panic(err)

	***REMOVED***

	logger := &logrus.Logger***REMOVED***
		Out:       stdErr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	***REMOVED***

	confDir, err := os.UserConfigDir()
	if err != nil ***REMOVED***
		logger.WithError(err).Warn("could not get config directory")
		confDir = ".config"
	***REMOVED***

	defaultFlags := getDefaultFlags(confDir)

	logrus.AddHook(hook)

	logrus.SetOutput(ioutil.Discard)

	return &globalState***REMOVED***
		ctx:            ctx,
		fs:             afero.NewMemMapFs(),
		getwd:          os.Getwd,
		args:           append(make([]string, 0, len(os.Args)), os.Args...), // copy
		envVars:        envVars,
		defaultFlags:   defaultFlags,
		flags:          getFlags(defaultFlags, envVars),
		stdOut:         stdOut,
		stdErr:         stdErr,
		stdIn:          os.Stdin,
		osExit:         os.Exit,
		signalNotify:   signal.Notify,
		signalStop:     signal.Stop,
		logger:         logger,
		fallbackLogger: logger,
	***REMOVED***
***REMOVED***

func (w *consoleWriter) Write(p []byte) (n int, err error) ***REMOVED***
	origLen := len(p)

	// Intercept the write message so can assess log errors

	stringP := string(p)

	fmt.Println("stringP:", stringP)

	// See if stringP contains 'source=stacktrace' or 'source=console'
	if strings.Contains(stringP, "source=stacktrace") || strings.Contains(stringP, "source=console") ***REMOVED***
		// Determine level of message
		if strings.Contains(stringP, "level=error") ***REMOVED***
			// Get bit after msg=" and before next "
			msg1 := strings.Split(stringP, "msg=\"")[1]
			msg2 := strings.Split(msg1, "\"")[0]

			go handleStringError(w.ctx, w.client, w.jobId, w.workerId, msg2)
		***REMOVED*** else ***REMOVED***
			go DispatchMessage(w.ctx, w.client, w.jobId, w.workerId, string(p), "MESSAGE")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		go DispatchMessage(w.ctx, w.client, w.jobId, w.workerId, string(p), "MESSAGE")
	***REMOVED***

	return origLen, err
***REMOVED***

func getDefaultFlags(homeFolder string) globalFlags ***REMOVED***
	return globalFlags***REMOVED***
		address:        "localhost:6565",
		configFilePath: filepath.Join(homeFolder, "loadimpact", "k6", defaultConfigFileName),
		logOutput:      "stderr",
	***REMOVED***
***REMOVED***

func getFlags(defaultFlags globalFlags, env map[string]string) globalFlags ***REMOVED***
	result := defaultFlags

	// TODO: add env vars for the rest of the values (after adjusting
	// rootCmdPersistentFlagSet(), of course)

	if val, ok := env["K6_CONFIG"]; ok ***REMOVED***
		result.configFilePath = val
	***REMOVED***
	if val, ok := env["K6_LOG_OUTPUT"]; ok ***REMOVED***
		result.logOutput = val
	***REMOVED***
	if val, ok := env["K6_LOG_FORMAT"]; ok ***REMOVED***
		result.logFormat = val
	***REMOVED***
	if env["K6_NO_COLOR"] != "" ***REMOVED***
		result.noColor = true
	***REMOVED***
	// Support https://no-color.org/, even an empty value should disable the
	// color output from k6.
	if _, ok := env["NO_COLOR"]; ok ***REMOVED***
		result.noColor = true
	***REMOVED***
	return result
***REMOVED***

func parseEnvKeyValue(kv string) (string, string) ***REMOVED***
	if idx := strings.IndexRune(kv, '='); idx != -1 ***REMOVED***
		return kv[:idx], kv[idx+1:]
	***REMOVED***
	return kv, ""
***REMOVED***

func buildEnvMap(environ []string) map[string]string ***REMOVED***
	env := make(map[string]string, len(environ))
	for _, kv := range environ ***REMOVED***
		k, v := parseEnvKeyValue(kv)
		env[k] = v
	***REMOVED***
	return env
***REMOVED***

// RawFormatter it does nothing with the message just prints it
type RawFormatter struct***REMOVED******REMOVED***

// Format renders a single log entry
func (f RawFormatter) Format(entry *logrus.Entry) ([]byte, error) ***REMOVED***
	return append([]byte(entry.Message), '\n'), nil
***REMOVED***
