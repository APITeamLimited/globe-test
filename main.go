package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/js"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/sampler"
	"github.com/loadimpact/speedboat/sampler/influxdb"
	"github.com/loadimpact/speedboat/sampler/stream"
	"github.com/loadimpact/speedboat/simple"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"time"
)

const (
	typeURL = "url"
	typeJS  = "js"
)

func pollVURamping(ctx context.Context, t lib.Test) <-chan int ***REMOVED***
	ch := make(chan int)
	startTime := time.Now()

	go func() ***REMOVED***
		defer close(ch)

		ticker := time.NewTicker(1 * time.Second)
		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				ch <- t.VUsAt(time.Since(startTime))
			case <-ctx.Done():
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***

func parseOutput(out, format string) (sampler.Output, error) ***REMOVED***
	switch ***REMOVED***
	case out == "-":
		return stream.New(format, os.Stdout)
	case strings.HasPrefix(out, "influxdb+"):
		url := strings.TrimPrefix(out, "influxdb+")
		return influxdb.NewFromURL(url)
	default:
		f, err := os.Create(out)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return stream.New(format, f)
	***REMOVED***
***REMOVED***

func guessType(arg string) string ***REMOVED***
	switch ***REMOVED***
	case strings.Contains(arg, "://"):
		return typeURL
	case strings.HasSuffix(arg, ".js"):
		return typeJS
	***REMOVED***
	return ""
***REMOVED***

func makeRunner(t lib.Test, filename, typ string) (lib.Runner, error) ***REMOVED***
	if typ == typeURL ***REMOVED***
		return simple.New(t), nil
	***REMOVED***

	bytes, err := ioutil.ReadFile(filename)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	switch typ ***REMOVED***
	case typeJS:
		return js.New(t, filename, string(bytes)), nil
	default:
		return nil, errors.New("Unknown type specified")
	***REMOVED***
***REMOVED***

func action(cc *cli.Context) error ***REMOVED***
	if cc.IsSet("verbose") ***REMOVED***
		log.SetLevel(log.DebugLevel)
	***REMOVED***

	sampler.DefaultSampler.OnError = func(err error) ***REMOVED***
		log.WithError(err).Error("[Sampler error]")
	***REMOVED***

	outFormat := cc.String("format")
	for _, out := range cc.StringSlice("metrics") ***REMOVED***
		output, err := parseOutput(out, outFormat)
		if err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
		sampler.DefaultSampler.Outputs = append(sampler.DefaultSampler.Outputs, output)
	***REMOVED***

	var t lib.Test
	var r lib.Runner

	// TODO: Majorly simplify this, along with the Test structure; the URL field is going
	// away in favor of environment variables (or something of the sort), which means 90%
	// of this code goes out the window - once things elsewhere stop depending on it >_>
	switch len(cc.Args()) ***REMOVED***
	case 0:
		cli.ShowAppHelp(cc)
		return nil
	case 1, 2:
		filename := cc.Args()[0]
		typ := cc.String("type")
		if typ == "" ***REMOVED***
			typ = guessType(filename)
		***REMOVED***

		switch typ ***REMOVED***
		case typeJS:
			t.Script = filename
		case typeURL:
			t.URL = filename
		case "":
			return cli.NewExitError("Ambiguous argument, please specify -t/--type", 1)
		default:
			return cli.NewExitError("Unknown type specified", 1)
		***REMOVED***

		if typ != typeURL && len(cc.Args()) > 1 ***REMOVED***
			t.URL = cc.Args()[1]
		***REMOVED***

		r_, err := makeRunner(t, filename, typ)
		if err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
		r = r_

	default:
		return cli.NewExitError("Too many arguments!", 1)
	***REMOVED***

	t.Stages = []lib.TestStage***REMOVED***
		lib.TestStage***REMOVED***
			Duration: cc.Duration("duration"),
			StartVUs: cc.Int("vus"),
			EndVUs:   cc.Int("vus"),
		***REMOVED***,
	***REMOVED***

	vus := lib.VUGroup***REMOVED***
		Pool: lib.VUPool***REMOVED***
			New: r.NewVU,
		***REMOVED***,
		RunOnce: func(ctx context.Context, vu lib.VU) ***REMOVED***
			if err := vu.RunOnce(ctx); err != nil ***REMOVED***
				log.WithError(err).Error("Uncaught Error")
			***REMOVED***
		***REMOVED***,
	***REMOVED***

	for i := 0; i < t.MaxVUs(); i++ ***REMOVED***
		vu, err := vus.Pool.New()
		if err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
		vus.Pool.Put(vu)
	***REMOVED***

	ctx, cancel := context.WithTimeout(context.Background(), t.TotalDuration())
	vus.Start(ctx)

	quit := make(chan os.Signal)
	signal.Notify(quit)

	scaleTo := pollVURamping(ctx, t)
mainLoop:
	for ***REMOVED***
		select ***REMOVED***
		case num := <-scaleTo:
			vus.Scale(num)
		case <-quit:
			cancel()
		case <-ctx.Done():
			break mainLoop
		***REMOVED***
	***REMOVED***

	vus.Stop()

	return nil
***REMOVED***

func main() ***REMOVED***
	// Free up -v and -h for our own flags
	cli.VersionFlag.Name = "version"
	cli.HelpFlag.Name = "help, ?"

	// Bootstrap the app from commandline flags
	app := cli.NewApp()
	app.Name = "speedboat"
	app.Usage = "A next-generation load generator"
	app.Version = "1.0.0-mvp1"
	app.Flags = []cli.Flag***REMOVED***
		cli.StringFlag***REMOVED***
			Name:  "type, t",
			Usage: "Input file type, if not evident (url or js)",
		***REMOVED***,
		cli.IntFlag***REMOVED***
			Name:  "vus, u",
			Usage: "Number of VUs to simulate",
			Value: 10,
		***REMOVED***,
		cli.DurationFlag***REMOVED***
			Name:  "duration, d",
			Usage: "Test duration",
			Value: time.Duration(10) * time.Second,
		***REMOVED***,
		cli.BoolFlag***REMOVED***
			Name:  "verbose, v",
			Usage: "More verbose output",
		***REMOVED***,
		cli.BoolFlag***REMOVED***
			Name:  "quiet, q",
			Usage: "Suppress the summary at the end of a test",
		***REMOVED***,
		cli.StringSliceFlag***REMOVED***
			Name:  "metrics, m",
			Usage: "Write metrics to a file or database",
		***REMOVED***,
		cli.StringFlag***REMOVED***
			Name:  "format, f",
			Usage: "Metric output format (json or csv)",
			Value: "json",
		***REMOVED***,
	***REMOVED***
	app.Action = action
	app.Run(os.Args)
***REMOVED***
