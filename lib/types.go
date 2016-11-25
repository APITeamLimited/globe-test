package lib

import (
	"gopkg.in/guregu/null.v3"
)

type Options struct ***REMOVED***
	Run      null.Bool   `json:"run"`
	VUs      null.Int    `json:"vus"`
	VUsMax   null.Int    `json:"vus-max"`
	Duration null.String `json:"duration"`

	Quit        null.Bool `json:"quit"`
	QuitOnTaint null.Bool `json:"quit-on-taint"`

	// Thresholds are JS snippets keyed by metrics.
	Thresholds map[string][]string `json:"thresholds"`
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
	if len(opts.Thresholds) > 0 ***REMOVED***
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
