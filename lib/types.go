package lib

import (
	"encoding/json"
	"github.com/robertkrimen/otto"
	"gopkg.in/guregu/null.v3"
)

type Options struct ***REMOVED***
	Run      null.Bool   `json:"run"`
	VUs      null.Int    `json:"vus"`
	VUsMax   null.Int    `json:"vus-max"`
	Duration null.String `json:"duration"`

	Quit        null.Bool  `json:"quit"`
	QuitOnTaint null.Bool  `json:"quit-on-taint"`
	Acceptance  null.Float `json:"acceptance"`

	Thresholds map[string][]*Threshold `json:"thresholds"`
***REMOVED***

func (o Options) Apply(opts Options) Options ***REMOVED***
	if opts.Run.Valid ***REMOVED***
		o.Run = opts.Run
	***REMOVED***
	if opts.VUs.Valid ***REMOVED***
		o.VUs = opts.VUs
	***REMOVED***
	if opts.VUsMax.Valid ***REMOVED***
		o.VUsMax = opts.VUsMax
	***REMOVED***
	if opts.Duration.Valid ***REMOVED***
		o.Duration = opts.Duration
	***REMOVED***
	if opts.Quit.Valid ***REMOVED***
		o.Quit = opts.Quit
	***REMOVED***
	if opts.QuitOnTaint.Valid ***REMOVED***
		o.QuitOnTaint = opts.QuitOnTaint
	***REMOVED***
	if opts.Acceptance.Valid ***REMOVED***
		o.Acceptance = opts.Acceptance
	***REMOVED***
	if opts.Thresholds != nil ***REMOVED***
		o.Thresholds = opts.Thresholds
	***REMOVED***
	return o
***REMOVED***

func (o Options) SetAllValid(valid bool) Options ***REMOVED***
	o.Run.Valid = valid
	o.VUs.Valid = valid
	o.VUsMax.Valid = valid
	o.Duration.Valid = valid
	o.Quit.Valid = valid
	o.QuitOnTaint.Valid = valid
	return o
***REMOVED***

type Threshold struct ***REMOVED***
	Source string
	Script *otto.Script
	Failed bool
***REMOVED***

func (t Threshold) String() string ***REMOVED***
	return t.Source
***REMOVED***

func (t Threshold) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal(t.Source)
***REMOVED***

func (t *Threshold) UnmarshalJSON(data []byte) error ***REMOVED***
	var src string
	if err := json.Unmarshal(data, &src); err != nil ***REMOVED***
		return err
	***REMOVED***
	t.Source = src
	t.Script = nil
	t.Failed = false
	return nil
***REMOVED***
