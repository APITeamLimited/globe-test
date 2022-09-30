package http

import (
	"context"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/APITeamLimited/globe-test/worker/js/common"
	"github.com/APITeamLimited/globe-test/worker/js/modulestest"
	"github.com/APITeamLimited/globe-test/worker/libWorker/netext/httpext"
)

//nolint:golint, revive
func getTestModuleInstance(
	t testing.TB,
) (*goja.Runtime, *ModuleInstance, *modulestest.VU) ***REMOVED***
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	root := New()
	mockVU := &modulestest.VU***REMOVED***
		RuntimeField: rt,
		InitEnvField: &common.InitEnvironment***REMOVED***
			Registry: workerMetrics.NewRegistry(),
		***REMOVED***,
		CtxField: ctx,
	***REMOVED***
	mi, ok := root.NewModuleInstance(mockVU).(*ModuleInstance)
	require.True(t, ok)

	require.NoError(t, rt.Set("http", mi.Exports().Default))

	return rt, mi, mockVU
***REMOVED***

func TestTagURL(t *testing.T) ***REMOVED***
	t.Parallel()

	testdata := map[string]struct***REMOVED*** u, n string ***REMOVED******REMOVED***
		`http://localhost/anything/`:               ***REMOVED***"http://localhost/anything/", "http://localhost/anything/"***REMOVED***,
		`http://localhost/anything/$***REMOVED***1+1***REMOVED***`:         ***REMOVED***"http://localhost/anything/2", "http://localhost/anything/$***REMOVED******REMOVED***"***REMOVED***,
		`http://localhost/anything/$***REMOVED***1+1***REMOVED***/`:        ***REMOVED***"http://localhost/anything/2/", "http://localhost/anything/$***REMOVED******REMOVED***/"***REMOVED***,
		`http://localhost/anything/$***REMOVED***1+1***REMOVED***/$***REMOVED***1+2***REMOVED***`:  ***REMOVED***"http://localhost/anything/2/3", "http://localhost/anything/$***REMOVED******REMOVED***/$***REMOVED******REMOVED***"***REMOVED***,
		`http://localhost/anything/$***REMOVED***1+1***REMOVED***/$***REMOVED***1+2***REMOVED***/`: ***REMOVED***"http://localhost/anything/2/3/", "http://localhost/anything/$***REMOVED******REMOVED***/$***REMOVED******REMOVED***/"***REMOVED***,
	***REMOVED***
	for expr, data := range testdata ***REMOVED***
		expr, data := expr, data
		t.Run("expr="+expr, func(t *testing.T) ***REMOVED***
			t.Parallel()
			rt, _, _ := getTestModuleInstance(t)
			tag, err := httpext.NewURL(data.u, data.n)
			require.NoError(t, err)
			v, err := rt.RunString("http.url`" + expr + "`")
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, tag, v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
