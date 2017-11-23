// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"io"
	"os"

	"github.com/gin-gonic/gin/binding"
)

const ENV_GIN_MODE = "GIN_MODE"

const (
	DebugMode   string = "debug"
	ReleaseMode string = "release"
	TestMode    string = "test"
)
const (
	debugCode = iota
	releaseCode
	testCode
)

// DefaultWriter is the default io.Writer used the Gin for debug output and
// middleware output like Logger() or Recovery().
// Note that both Logger and Recovery provides custom ways to configure their
// output io.Writer.
// To support coloring in Windows use:
// 		import "github.com/mattn/go-colorable"
// 		gin.DefaultWriter = colorable.NewColorableStdout()
var DefaultWriter io.Writer = os.Stdout
var DefaultErrorWriter io.Writer = os.Stderr

var ginMode = debugCode
var modeName = DebugMode

func init() ***REMOVED***
	mode := os.Getenv(ENV_GIN_MODE)
	if len(mode) == 0 ***REMOVED***
		SetMode(DebugMode)
	***REMOVED*** else ***REMOVED***
		SetMode(mode)
	***REMOVED***
***REMOVED***

func SetMode(value string) ***REMOVED***
	switch value ***REMOVED***
	case DebugMode:
		ginMode = debugCode
	case ReleaseMode:
		ginMode = releaseCode
	case TestMode:
		ginMode = testCode
	default:
		panic("gin mode unknown: " + value)
	***REMOVED***
	modeName = value
***REMOVED***

func DisableBindValidation() ***REMOVED***
	binding.Validator = nil
***REMOVED***

func Mode() string ***REMOVED***
	return modeName
***REMOVED***
