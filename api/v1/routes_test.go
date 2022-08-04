package v1

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.k6.io/k6/api/common"
	"go.k6.io/k6/core"
)

func newRequestWithEngine(engine *core.Engine, method, target string, body io.Reader) *http.Request ***REMOVED***
	r := httptest.NewRequest(method, target, body)
	return r.WithContext(common.WithEngine(r.Context(), engine))
***REMOVED***

func TestNewHandler(t *testing.T) ***REMOVED***
	assert.NotNil(t, NewHandler())
***REMOVED***
