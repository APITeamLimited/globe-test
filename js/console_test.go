package js

import (
	"context"
	"fmt"
	"testing"

	log "github.com/Sirupsen/logrus"
	logtest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestConsoleLog(t *testing.T) ***REMOVED***
	levels := map[string]log.Level***REMOVED***
		"log":   log.InfoLevel,
		"debug": log.DebugLevel,
		"info":  log.InfoLevel,
		"warn":  log.WarnLevel,
		"error": log.ErrorLevel,
	***REMOVED***
	argsets := map[string]log.Fields***REMOVED***
		`"a"`:       ***REMOVED***"arg0": "a"***REMOVED***,
		`"a","b"`:   ***REMOVED***"arg0": "a", "arg1": "b"***REMOVED***,
		`***REMOVED***a:1***REMOVED***`:     ***REMOVED***"a": "1"***REMOVED***,
		`***REMOVED***a:1,b:2***REMOVED***`: ***REMOVED***"a": "1", "b": "2"***REMOVED***,
		`"a",***REMOVED***a:1***REMOVED***`: ***REMOVED***"arg0": "a", "a": "1"***REMOVED***,
		`***REMOVED***a:1***REMOVED***,"a"`: ***REMOVED***"a": "1", "arg1": "a"***REMOVED***,
	***REMOVED***
	for name, level := range levels ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			rt, err := New()
			assert.NoError(t, err)

			logger, hook := logtest.NewNullLogger()
			logger.Level = log.DebugLevel
			_ = rt.VM.Set("__console__", &Console***REMOVED***logger***REMOVED***)

			for args, fields := range argsets ***REMOVED***
				t.Run(args, func(t *testing.T) ***REMOVED***
					_ = rt.VM.Set("__initapi__", &InitAPI***REMOVED***r: rt***REMOVED***)
					exp, err := rt.load("__snippet__", []byte(fmt.Sprintf(`
					console.%s("init",%s);
					export default function() ***REMOVED***
						console.%s("default",%s);
					***REMOVED***
					`, name, args, name, args)))
					if !assert.NoError(t, err) ***REMOVED***
						return
					***REMOVED***
					_ = rt.VM.Set("__initapi__", nil)

					initEntry := hook.LastEntry()
					if assert.NotNil(t, initEntry, "nothing logged from init") ***REMOVED***
						assert.Equal(t, "init", initEntry.Message)
						assert.Equal(t, level, initEntry.Level)
						assert.EqualValues(t, fields, initEntry.Data)
					***REMOVED***

					r, err := NewRunner(rt, exp)
					if !assert.NoError(t, err) ***REMOVED***
						return
					***REMOVED***

					vu, err := r.NewVU()
					assert.NoError(t, err)

					_, err = vu.RunOnce(context.Background())
					assert.NoError(t, err)

					entry := hook.LastEntry()
					if assert.NotNil(t, entry, "nothing logged from default") ***REMOVED***
						assert.Equal(t, "default", entry.Message)
						assert.Equal(t, level, entry.Level)
						assert.EqualValues(t, fields, entry.Data)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
