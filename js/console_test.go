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

package js

import (
	"context"
	"fmt"
	"testing"

	log "github.com/Sirupsen/logrus"
	logtest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/loadimpact/k6/lib"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestConsole(t *testing.T) ***REMOVED***
	levels := map[string]log.Level***REMOVED***
		"log":   log.InfoLevel,
		"debug": log.DebugLevel,
		"info":  log.InfoLevel,
		"warn":  log.WarnLevel,
		"error": log.ErrorLevel,
	***REMOVED***
	argsets := map[string]struct ***REMOVED***
		Message string
		Data    log.Fields
	***REMOVED******REMOVED***
		`"string"`:         ***REMOVED***Message: "string"***REMOVED***,
		`"string","a","b"`: ***REMOVED***Message: "string", Data: log.Fields***REMOVED***"0": "a", "1": "b"***REMOVED******REMOVED***,
		`"string",1,2`:     ***REMOVED***Message: "string", Data: log.Fields***REMOVED***"0": "1", "1": "2"***REMOVED******REMOVED***,
		`***REMOVED******REMOVED***`:               ***REMOVED***Message: "[object Object]"***REMOVED***,
	***REMOVED***
	for name, level := range levels ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			for args, result := range argsets ***REMOVED***
				t.Run(args, func(t *testing.T) ***REMOVED***
					r, err := New(&lib.SourceData***REMOVED***
						Filename: "/script",
						Data: []byte(fmt.Sprintf(
							`export default function() ***REMOVED*** console.%s(%s); ***REMOVED***`,
							name, args,
						)),
					***REMOVED***, afero.NewMemMapFs())
					assert.NoError(t, err)

					vu, err := r.newVU()
					assert.NoError(t, err)

					logger, hook := logtest.NewNullLogger()
					logger.Level = log.DebugLevel
					vu.VUContext.Console.Logger = logger

					_, err = vu.RunOnce(context.Background())
					assert.NoError(t, err)

					entry := hook.LastEntry()
					if assert.NotNil(t, entry, "nothing logged") ***REMOVED***
						assert.Equal(t, level, entry.Level)
						assert.Equal(t, result.Message, entry.Message)

						data := result.Data
						if data == nil ***REMOVED***
							data = make(log.Fields)
						***REMOVED***
						assert.Equal(t, data, entry.Data)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
