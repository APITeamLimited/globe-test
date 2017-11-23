// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"html/template"
	"log"
)

func init() ***REMOVED***
	log.SetFlags(0)
***REMOVED***

// IsDebugging returns true if the framework is running in debug mode.
// Use SetMode(gin.Release) to switch to disable the debug mode.
func IsDebugging() bool ***REMOVED***
	return ginMode == debugCode
***REMOVED***

func debugPrintRoute(httpMethod, absolutePath string, handlers HandlersChain) ***REMOVED***
	if IsDebugging() ***REMOVED***
		nuHandlers := len(handlers)
		handlerName := nameOfFunction(handlers.Last())
		debugPrint("%-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
	***REMOVED***
***REMOVED***

func debugPrintLoadTemplate(tmpl *template.Template) ***REMOVED***
	if IsDebugging() ***REMOVED***
		var buf bytes.Buffer
		for _, tmpl := range tmpl.Templates() ***REMOVED***
			buf.WriteString("\t- ")
			buf.WriteString(tmpl.Name())
			buf.WriteString("\n")
		***REMOVED***
		debugPrint("Loaded HTML Templates (%d): \n%s\n", len(tmpl.Templates()), buf.String())
	***REMOVED***
***REMOVED***

func debugPrint(format string, values ...interface***REMOVED******REMOVED***) ***REMOVED***
	if IsDebugging() ***REMOVED***
		log.Printf("[GIN-debug] "+format, values...)
	***REMOVED***
***REMOVED***

func debugPrintWARNINGNew() ***REMOVED***
	debugPrint(`[WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

`)
***REMOVED***

func debugPrintWARNINGSetHTMLTemplate() ***REMOVED***
	debugPrint(`[WARNING] Since SetHTMLTemplate() is NOT thread-safe. It should only be called
at initialization. ie. before any route is registered or the router is listening in a socket:

	router := gin.Default()
	router.SetHTMLTemplate(template) // << good place

`)
***REMOVED***

func debugPrintError(err error) ***REMOVED***
	if err != nil ***REMOVED***
		debugPrint("[ERROR] %v\n", err)
	***REMOVED***
***REMOVED***
