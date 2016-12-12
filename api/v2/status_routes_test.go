package v2

import (
	"encoding/json"
	"github.com/loadimpact/k6/lib"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStatus(t *testing.T) ***REMOVED***
	engine, err := lib.NewEngine(nil)
	assert.NoError(t, err)

	rw := httptest.NewRecorder()
	NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "GET", "/v2/status", nil))
	res := rw.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	t.Run("document", func(t *testing.T) ***REMOVED***
		var doc jsonapi.Document
		assert.NoError(t, json.Unmarshal(rw.Body.Bytes(), &doc))
		if !assert.NotNil(t, doc.Data.DataObject) ***REMOVED***
			return
		***REMOVED***
		assert.Equal(t, "status", doc.Data.DataObject.Type)
	***REMOVED***)

	t.Run("status", func(t *testing.T) ***REMOVED***
		var status Status
		assert.NoError(t, jsonapi.Unmarshal(rw.Body.Bytes(), &status))
		assert.True(t, status.Running.Valid)
		assert.True(t, status.Tainted.Valid)
		assert.True(t, status.VUs.Valid)
		assert.True(t, status.VUsMax.Valid)
	***REMOVED***)
***REMOVED***
