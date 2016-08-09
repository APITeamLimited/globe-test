package postman

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/lib"
	"golang.org/x/net/context"
	"strings"
)

type ErrorWithLineNumber struct ***REMOVED***
	Wrapped error
	Line    int
***REMOVED***

func (e ErrorWithLineNumber) Error() string ***REMOVED***
	return fmt.Sprintf("%s (line %d)", e.Wrapped.Error(), e.Line)
***REMOVED***

type Runner struct ***REMOVED***
	Collection Collection
***REMOVED***

type VU struct ***REMOVED***
	Runner *Runner
***REMOVED***

func New(source []byte) (*Runner, error) ***REMOVED***
	var collection Collection
	if err := json.Unmarshal(source, &collection); err != nil ***REMOVED***
		switch e := err.(type) ***REMOVED***
		case *json.SyntaxError:
			src := string(source)
			line := strings.Count(src[:e.Offset], "\n") + 1
			return nil, ErrorWithLineNumber***REMOVED***Wrapped: e, Line: line***REMOVED***
		case *json.UnmarshalTypeError:
			src := string(source)
			line := strings.Count(src[:e.Offset], "\n") + 1
			return nil, ErrorWithLineNumber***REMOVED***Wrapped: e, Line: line***REMOVED***
		***REMOVED***
		return nil, err
	***REMOVED***

	return &Runner***REMOVED***
		Collection: collection,
	***REMOVED***, nil
***REMOVED***

func (r *Runner) NewVU() (lib.VU, error) ***REMOVED***
	return &VU***REMOVED***Runner: r***REMOVED***, nil
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	return nil
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) error ***REMOVED***
	for _, item := range u.Runner.Collection.Item ***REMOVED***
		if err := u.runItem(item, u.Runner.Collection.Auth); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (u *VU) runItem(i Item, a Auth) error ***REMOVED***
	if i.Auth.Type != "" ***REMOVED***
		a = i.Auth
	***REMOVED***

	if i.Request.URL != "" ***REMOVED***
		log.WithField("url", i.Request.URL).Info("Request!")
	***REMOVED***

	for _, item := range i.Item ***REMOVED***
		if err := u.runItem(item, a); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
