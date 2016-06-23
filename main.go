package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/js"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/simple"
	"github.com/loadimpact/speedboat/stats"
	"github.com/loadimpact/speedboat/stats/influxdb"
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

func parseBackend(out string) (stats.Backend, error) ***REMOVED***
	switch ***REMOVED***
	case out == "-":
		return stats.NewJSONBackend(os.Stdout), nil
	case strings.HasPrefix(out, "influxdb+"):
		url := strings.TrimPrefix(out, "influxdb+")
		return influxdb.NewFromURL(url)
	default:
		f, err := os.Create(out)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return stats.NewJSONBackend(f), nil
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

func readAll(filename string) ([]byte, error) ***REMOVED***
	if filename == "-" ***REMOVED***
		return ioutil.ReadAll(os.Stdin)
	***REMOVED***

	return ioutil.ReadFile(filename)
***REMOVED***

func makeRunner(t lib.Test, filename, typ string) (lib.Runner, error) ***REMOVED***
	if typ == typeURL ***REMOVED***
		return simple.New(t), nil
	***REMOVED***

	bytes, err := readAll(filename)
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

	for _, out := range cc.StringSlice("metrics") ***REMOVED***
		backend, err := parseBackend(out)
		if err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
		stats.DefaultRegistry.Backends = append(stats.DefaultRegistry.Backends, backend)
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

		if filename == "-" && typ == "" ***REMOVED***
			return cli.NewExitError("Reading from stdin requires a -t/--type flag", 1)
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

	go func() ***REMOVED***
		ticker := time.NewTicker(1 * time.Second)
		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				if err := stats.Submit(); err != nil ***REMOVED***
					log.WithError(err).Error("[Couldn't submit stats]")
				***REMOVED***
			case <-ctx.Done():
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	quit := make(chan os.Signal)
	signal.Notify(quit)

	vus.Start(ctx)
	scaleTo := pollVURamping(ctx, t)
	mVUs := stats.Stat***REMOVED***Name: "vus", Type: stats.GaugeType***REMOVED***
mainLoop:
	for ***REMOVED***
		select ***REMOVED***
		case num := <-scaleTo:
			vus.Scale(num)
			stats.Add(stats.Point***REMOVED***
				Stat:   &mVUs,
				Values: stats.Value(float64(num)),
			***REMOVED***)
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
	***REMOVED***
	app.Action = action
	app.Run(os.Args)
***REMOVED***
