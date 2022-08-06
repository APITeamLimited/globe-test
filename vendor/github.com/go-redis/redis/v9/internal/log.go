package internal

import (
	"context"
	"fmt"
	"log"
	"os"
)

type Logging interface ***REMOVED***
	Printf(ctx context.Context, format string, v ...interface***REMOVED******REMOVED***)
***REMOVED***

type logger struct ***REMOVED***
	log *log.Logger
***REMOVED***

func (l *logger) Printf(ctx context.Context, format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	_ = l.log.Output(2, fmt.Sprintf(format, v...))
***REMOVED***

// Logger calls Output to print to the stderr.
// Arguments are handled in the manner of fmt.Print.
var Logger Logging = &logger***REMOVED***
	log: log.New(os.Stderr, "redis: ", log.LstdFlags|log.Lshortfile),
***REMOVED***
