package v2

import (
	"github.com/julienschmidt/httprouter"
	"github.com/loadimpact/k6/api/common"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/guregu/null.v3"
	"net/http"
)

func HandleGetStatus(rw http.ResponseWriter, r *http.Request, p httprouter.Params) ***REMOVED***
	engine := common.GetEngine(r.Context())

	status := Status***REMOVED***
		Running: null.BoolFrom(engine.Status.Running.Bool),
		Tainted: null.BoolFrom(engine.Status.Tainted.Bool),
		VUs:     null.IntFrom(engine.Status.VUs.Int64),
		VUsMax:  null.IntFrom(engine.Status.VUsMax.Int64),
	***REMOVED***
	data, err := jsonapi.Marshal(status)
	if err != nil ***REMOVED***
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***
	_, _ = rw.Write(data)
***REMOVED***
