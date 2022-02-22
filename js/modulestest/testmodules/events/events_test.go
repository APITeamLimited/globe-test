package events

import (
	"context"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/eventloop"
	"go.k6.io/k6/js/modulestest"
)

func TestSetTimeout(t *testing.T) ***REMOVED***
	t.Parallel()
	rt := goja.New()
	vu := &modulestest.VU***REMOVED***
		RuntimeField: rt,
		InitEnvField: &common.InitEnvironment***REMOVED******REMOVED***,
		CtxField:     context.Background(),
		StateField:   nil,
	***REMOVED***

	m, ok := New().NewModuleInstance(vu).(*Events)
	require.True(t, ok)
	var log []string
	require.NoError(t, rt.Set("events", m.Exports().Named))
	require.NoError(t, rt.Set("print", func(s string) ***REMOVED*** log = append(log, s) ***REMOVED***))
	loop := eventloop.New(vu)
	vu.RegisterCallbackField = loop.RegisterCallback

	err := loop.Start(func() error ***REMOVED***
		_, err := vu.Runtime().RunString(`
      events.setTimeout(()=> ***REMOVED***
        print("in setTimeout")
      ***REMOVED***)
      print("outside setTimeout")
      `)
		return err
	***REMOVED***)
	require.NoError(t, err)
	require.Equal(t, []string***REMOVED***"outside setTimeout", "in setTimeout"***REMOVED***, log)
***REMOVED***
