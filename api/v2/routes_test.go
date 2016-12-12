package v2

import (
	"github.com/loadimpact/k6/api/common"
	"github.com/loadimpact/k6/lib"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newRequestWithEngine(engine *lib.Engine, method, target string, body io.Reader) *http.Request ***REMOVED***
	r := httptest.NewRequest(method, target, body)
	return r.WithContext(common.WithEngine(r.Context(), engine))
***REMOVED***

func TestNewHandler(t *testing.T) ***REMOVED***
	assert.NotNil(t, NewHandler())
***REMOVED***
