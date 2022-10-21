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
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
)

//nolint:golint, revive
func getTestModuleInstance(
	t testing.TB,
) (*goja.Runtime, *ModuleInstance, *modulestest.VU) {
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper{})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	root := New()
	mockVU := &modulestest.VU{
		RuntimeField: rt,
		InitEnvField: &common.InitEnvironment{
			Registry: workerMetrics.NewRegistry(),
		},
		CtxField: ctx,
	}
	mi, ok := root.NewModuleInstance(mockVU).(*ModuleInstance)
	require.True(t, ok)

	require.NoError(t, rt.Set("http", mi.Exports().Default))

	return rt, mi, mockVU
}

func TestTagURL(t *testing.T) {
	t.Parallel()

	testdata := map[string]struct{ u, n string }{
		`http://localhost/anything/`:               {"http://localhost/anything/", "http://localhost/anything/"},
		`http://localhost/anything/${1+1}`:         {"http://localhost/anything/2", "http://localhost/anything/${}"},
		`http://localhost/anything/${1+1}/`:        {"http://localhost/anything/2/", "http://localhost/anything/${}/"},
		`http://localhost/anything/${1+1}/${1+2}`:  {"http://localhost/anything/2/3", "http://localhost/anything/${}/${}"},
		`http://localhost/anything/${1+1}/${1+2}/`: {"http://localhost/anything/2/3/", "http://localhost/anything/${}/${}/"},
	}
	for expr, data := range testdata {
		expr, data := expr, data
		t.Run("expr="+expr, func(t *testing.T) {
			t.Parallel()
			rt, _, _ := getTestModuleInstance(t)
			tag, err := httpext.NewURL(data.u, data.n)
			require.NoError(t, err)
			v, err := rt.RunString("http.url`" + expr + "`")
			if assert.NoError(t, err) {
				assert.Equal(t, tag, v.Export())
			}
		})
	}
}
