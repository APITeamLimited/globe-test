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

package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	stdlog "log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"go.k6.io/k6/lib/consts"
	"go.k6.io/k6/log"
)

var BannerColor = color.New(color.FgCyan)

//TODO: remove these global variables
//nolint:gochecknoglobals
var (
	outMutex   = &sync.Mutex***REMOVED******REMOVED***
	isDumbTerm = os.Getenv("TERM") == "dumb"
	stdoutTTY  = !isDumbTerm && (isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()))
	stderrTTY  = !isDumbTerm && (isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd()))
	stdout     = &consoleWriter***REMOVED***colorable.NewColorableStdout(), stdoutTTY, outMutex, nil***REMOVED***
	stderr     = &consoleWriter***REMOVED***colorable.NewColorableStderr(), stderrTTY, outMutex, nil***REMOVED***
)

const (
	defaultConfigFileName   = "config.json"
	waitRemoteLoggerTimeout = time.Second * 5
)

//TODO: remove these global variables
//nolint:gochecknoglobals
var defaultConfigFilePath = defaultConfigFileName // Updated with the user's config folder in the init() function below
//nolint:gochecknoglobals
var configFilePath = os.Getenv("K6_CONFIG") // Overridden by `-c`/`--config` flag!

//nolint:gochecknoglobals
var (
	// TODO: have environment variables for configuring these? hopefully after we move away from global vars though...
	quiet   bool
	noColor bool
	address string
)

// This is to keep all fields needed for the main/root k6 command
type rootCommand struct ***REMOVED***
	ctx            context.Context
	logger         *logrus.Logger
	fallbackLogger logrus.FieldLogger
	cmd            *cobra.Command
	loggerStopped  <-chan struct***REMOVED******REMOVED***
	logOutput      string
	logFmt         string
	loggerIsRemote bool
	verbose        bool
***REMOVED***

func newRootCommand(ctx context.Context, logger *logrus.Logger, fallbackLogger logrus.FieldLogger) *rootCommand ***REMOVED***
	c := &rootCommand***REMOVED***
		ctx:            ctx,
		logger:         logger,
		fallbackLogger: fallbackLogger,
	***REMOVED***
	// the base command when called without any subcommands.
	c.cmd = &cobra.Command***REMOVED***
		Use:               "k6",
		Short:             "a next-generation load generator",
		Long:              BannerColor.Sprintf("\n%s", consts.Banner()),
		SilenceUsage:      true,
		SilenceErrors:     true,
		PersistentPreRunE: c.persistentPreRunE,
	***REMOVED***

	confDir, err := os.UserConfigDir()
	if err != nil ***REMOVED***
		logrus.WithError(err).Warn("could not get config directory")
		confDir = ".config"
	***REMOVED***
	defaultConfigFilePath = filepath.Join(
		confDir,
		"loadimpact",
		"k6",
		defaultConfigFileName,
	)

	c.cmd.PersistentFlags().AddFlagSet(c.rootCmdPersistentFlagSet())
	return c
***REMOVED***

func (c *rootCommand) persistentPreRunE(cmd *cobra.Command, args []string) error ***REMOVED***
	var err error
	if !cmd.Flags().Changed("log-output") ***REMOVED***
		if envLogOutput, ok := os.LookupEnv("K6_LOG_OUTPUT"); ok ***REMOVED***
			c.logOutput = envLogOutput
		***REMOVED***
	***REMOVED***
	c.loggerStopped, err = c.setupLoggers()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	select ***REMOVED***
	case <-c.loggerStopped:
	default:
		c.loggerIsRemote = true
	***REMOVED***

	if noColor ***REMOVED***
		// TODO: figure out something else... currently, with the wrappers
		// below, we're stripping any colors from the output after we've
		// added them. The problem is that, besides being very inefficient,
		// this actually also strips other special characters from the
		// intended output, like the progressbar formatting ones, which
		// would otherwise be fine (in a TTY).
		//
		// It would be much better if we avoid messing with the output and
		// instead have a parametrized instance of the color library. It
		// will return colored output if colors are enabled and simply
		// return the passed input as-is (i.e. be a noop) if colors are
		// disabled...
		stdout.Writer = colorable.NewNonColorable(os.Stdout)
		stderr.Writer = colorable.NewNonColorable(os.Stderr)
	***REMOVED***
	stdlog.SetOutput(c.logger.Writer())
	c.logger.Debugf("k6 version: v%s", consts.FullVersion())
	return nil
***REMOVED***

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := &logrus.Logger***REMOVED***
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	***REMOVED***

	var fallbackLogger logrus.FieldLogger = &logrus.Logger***REMOVED***
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	***REMOVED***

	c := newRootCommand(ctx, logger, fallbackLogger)

	loginCmd := getLoginCmd()
	loginCmd.AddCommand(getLoginCloudCommand(logger), getLoginInfluxDBCommand(logger))
	c.cmd.AddCommand(
		getArchiveCmd(logger),
		getCloudCmd(ctx, logger),
		getConvertCmd(),
		getInspectCmd(logger),
		loginCmd,
		getPauseCmd(ctx),
		getResumeCmd(ctx),
		getScaleCmd(ctx),
		getRunCmd(ctx, logger),
		getStatsCmd(ctx),
		getStatusCmd(ctx),
		getVersionCmd(),
	)

	if err := c.cmd.Execute(); err != nil ***REMOVED***
		fields := logrus.Fields***REMOVED******REMOVED***
		code := -1
		if e, ok := err.(ExitCode); ok ***REMOVED***
			code = e.Code
			if e.Hint != "" ***REMOVED***
				fields["hint"] = e.Hint
			***REMOVED***
		***REMOVED***

		logger.WithFields(fields).Error(err)
		if c.loggerIsRemote ***REMOVED***
			fallbackLogger.WithFields(fields).Error(err)
			cancel()
			c.waitRemoteLogger()
		***REMOVED***

		os.Exit(code)
	***REMOVED***

	cancel()
	c.waitRemoteLogger()
***REMOVED***

func (c *rootCommand) waitRemoteLogger() ***REMOVED***
	if c.loggerIsRemote ***REMOVED***
		select ***REMOVED***
		case <-c.loggerStopped:
		case <-time.After(waitRemoteLoggerTimeout):
			c.fallbackLogger.Error("Remote logger didn't stop in %s", waitRemoteLoggerTimeout)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *rootCommand) rootCmdPersistentFlagSet() *pflag.FlagSet ***REMOVED***
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	// TODO: figure out a better way to handle the CLI flags - global variables are not very testable... :/
	flags.BoolVarP(&c.verbose, "verbose", "v", false, "enable verbose logging")
	flags.BoolVarP(&quiet, "quiet", "q", false, "disable progress updates")
	flags.BoolVar(&noColor, "no-color", false, "disable colored output")
	flags.StringVar(&c.logOutput, "log-output", "stderr",
		"change the output for k6 logs, possible values are stderr,stdout,none,loki[=host:port]")
	flags.StringVar(&c.logFmt, "logformat", "", "log output format") // TODO rename to log-format and warn on old usage
	flags.StringVarP(&address, "address", "a", "localhost:6565", "address for the api server")

	// TODO: Fix... This default value needed, so both CLI flags and environment variables work
	flags.StringVarP(&configFilePath, "config", "c", configFilePath, "JSON config file")
	// And we also need to explicitly set the default value for the usage message here, so things
	// like `K6_CONFIG="blah" k6 run -h` don't produce a weird usage message
	flags.Lookup("config").DefValue = defaultConfigFilePath
	must(cobra.MarkFlagFilename(flags, "config"))
	return flags
***REMOVED***

// fprintf panics when where's an error writing to the supplied io.Writer
func fprintf(w io.Writer, format string, a ...interface***REMOVED******REMOVED***) (n int) ***REMOVED***
	n, err := fmt.Fprintf(w, format, a...)
	if err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
	return n
***REMOVED***

// RawFormatter it does nothing with the message just prints it
type RawFormatter struct***REMOVED******REMOVED***

// Format renders a single log entry
func (f RawFormatter) Format(entry *logrus.Entry) ([]byte, error) ***REMOVED***
	return append([]byte(entry.Message), '\n'), nil
***REMOVED***

// The returned channel will be closed when the logger has finished flushing and pushing logs after
// the provided context is closed. It is closed if the logger isn't buffering and sending messages
// Asynchronously
func (c *rootCommand) setupLoggers() (<-chan struct***REMOVED******REMOVED***, error) ***REMOVED***
	ch := make(chan struct***REMOVED******REMOVED***)
	close(ch)

	if c.verbose ***REMOVED***
		c.logger.SetLevel(logrus.DebugLevel)
	***REMOVED***
	switch c.logOutput ***REMOVED***
	case "stderr":
		c.logger.SetOutput(stderr)
	case "stdout":
		c.logger.SetOutput(stdout)
	case "none":
		c.logger.SetOutput(ioutil.Discard)
	default:
		if !strings.HasPrefix(c.logOutput, "loki") ***REMOVED***
			return nil, fmt.Errorf("unsupported log output `%s`", c.logOutput)
		***REMOVED***
		ch = make(chan struct***REMOVED******REMOVED***)
		hook, err := log.LokiFromConfigLine(c.ctx, c.fallbackLogger, c.logOutput, ch)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		c.logger.AddHook(hook)
		c.logger.SetOutput(ioutil.Discard) // don't output to anywhere else
		c.logFmt = "raw"
		noColor = true // disable color
	***REMOVED***

	switch c.logFmt ***REMOVED***
	case "raw":
		c.logger.SetFormatter(&RawFormatter***REMOVED******REMOVED***)
		c.logger.Debug("Logger format: RAW")
	case "json":
		c.logger.SetFormatter(&logrus.JSONFormatter***REMOVED******REMOVED***)
		c.logger.Debug("Logger format: JSON")
	default:
		c.logger.SetFormatter(&logrus.TextFormatter***REMOVED***ForceColors: stderrTTY, DisableColors: noColor***REMOVED***)
		c.logger.Debug("Logger format: TEXT")
	***REMOVED***
	return ch, nil
***REMOVED***
