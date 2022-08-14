package har

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.k6.io/k6/js"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/loader"
	"go.k6.io/k6/metrics"
)

func TestBuildK6Headers(t *testing.T) ***REMOVED***
	headers := []struct ***REMOVED***
		values   []Header
		expected []string
	***REMOVED******REMOVED***
		***REMOVED***[]Header***REMOVED******REMOVED***"name", "1"***REMOVED***, ***REMOVED***"name", "2"***REMOVED******REMOVED***, []string***REMOVED***`"name": "1"`***REMOVED******REMOVED***,
		***REMOVED***[]Header***REMOVED******REMOVED***"name", "1"***REMOVED***, ***REMOVED***"name2", "2"***REMOVED******REMOVED***, []string***REMOVED***`"name": "1"`, `"name2": "2"`***REMOVED******REMOVED***,
		***REMOVED***[]Header***REMOVED******REMOVED***":host", "localhost"***REMOVED******REMOVED***, []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***

	for _, pair := range headers ***REMOVED***
		v := buildK6Headers(pair.values)
		assert.Equal(t, len(v), len(pair.expected), fmt.Sprintf("params: %v", pair.values))
	***REMOVED***
***REMOVED***

func TestBuildK6RequestObject(t *testing.T) ***REMOVED***
	req := &Request***REMOVED***
		Method:  "get",
		URL:     "http://www.google.es",
		Headers: []Header***REMOVED******REMOVED***"accept-language", "es-ES,es;q=0.8"***REMOVED******REMOVED***,
		Cookies: []Cookie***REMOVED******REMOVED***Name: "a", Value: "b"***REMOVED******REMOVED***,
	***REMOVED***
	v, err := buildK6RequestObject(req)
	assert.NoError(t, err)
	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	_, err = js.New(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, &loader.SourceData***REMOVED***
			URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
			Data: []byte(fmt.Sprintf("export default function() ***REMOVED*** res = http.batch([%v]); ***REMOVED***", v)),
		***REMOVED***, nil, lib.GetTestWorkerInfo())
	assert.NoError(t, err)
***REMOVED***

func TestBuildK6Body(t *testing.T) ***REMOVED***
	bodyText := "ccustemail=ppcano%40gmail.com&size=medium&topping=cheese&delivery=12%3A00&comments="

	req := &Request***REMOVED***
		Method: "post",
		URL:    "http://www.google.es",
		PostData: &PostData***REMOVED***
			MimeType: "application/x-www-form-urlencoded",
			Text:     bodyText,
		***REMOVED***,
	***REMOVED***
	postParams, plainText, err := buildK6Body(req)
	assert.NoError(t, err)
	assert.Equal(t, len(postParams), 0, "postParams should be empty")
	assert.Equal(t, bodyText, plainText)

	email := "user@mail.es"
	expectedEmailParam := fmt.Sprintf(`"email": %q`, email)

	req = &Request***REMOVED***
		Method: "post",
		URL:    "http://www.google.es",
		PostData: &PostData***REMOVED***
			MimeType: "application/x-www-form-urlencoded",
			Params: []Param***REMOVED***
				***REMOVED***Name: "email", Value: url.QueryEscape(email)***REMOVED***,
				***REMOVED***Name: "pw", Value: "hola"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	postParams, plainText, err = buildK6Body(req)
	assert.NoError(t, err)
	assert.Equal(t, plainText, "", "expected empty plainText")
	assert.Equal(t, len(postParams), 2, "postParams should have two items")
	assert.Equal(t, postParams[0], expectedEmailParam, "expected unescaped value")
***REMOVED***
