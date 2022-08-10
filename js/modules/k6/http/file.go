package http

import (
	"fmt"
	"strings"
	"time"

	"go.k6.io/k6/js/common"
)

// FileData represents a binary file requiring multipart request encoding
type FileData struct ***REMOVED***
	Data        []byte
	Filename    string
	ContentType string
***REMOVED***

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string ***REMOVED***
	return quoteEscaper.Replace(s)
***REMOVED***

// File returns a FileData object.
func (mi *ModuleInstance) file(data interface***REMOVED******REMOVED***, args ...string) FileData ***REMOVED***
	// supply valid default if filename and content-type are not specified
	fname, ct := fmt.Sprintf("%d", time.Now().UnixNano()), "application/octet-stream"

	if len(args) > 0 ***REMOVED***
		fname = args[0]

		if len(args) > 1 ***REMOVED***
			ct = args[1]
		***REMOVED***
	***REMOVED***

	dt, err := common.ToBytes(data)
	if err != nil ***REMOVED***
		common.Throw(mi.vu.Runtime(), err)
	***REMOVED***

	return FileData***REMOVED***
		Data:        dt,
		Filename:    fname,
		ContentType: ct,
	***REMOVED***
***REMOVED***
