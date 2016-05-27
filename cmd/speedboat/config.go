package main

import (
	"errors"
	"github.com/loadimpact/speedboat"
	"time"
)

// Configuration type for a state.
type ConfigStage struct ***REMOVED***
	Duration interface***REMOVED******REMOVED*** `yaml:"duration"`
	VUs      interface***REMOVED******REMOVED*** `yaml:"vus"`
***REMOVED***

type Config struct ***REMOVED***
	Duration string        `yaml:"duration"`
	Script   string        `yaml:"script"`
	URL      string        `yaml:"url"`
	VUs      interface***REMOVED******REMOVED***   `yaml:"vus"`
	Stages   []ConfigStage `yaml:"stages"`
***REMOVED***

func parseVUs(vus interface***REMOVED******REMOVED***) (int, int, error) ***REMOVED***
	switch v := vus.(type) ***REMOVED***
	case int:
		return v, v, nil
	case []interface***REMOVED******REMOVED***:
		switch len(v) ***REMOVED***
		case 1:
			n, ok := v[0].(int)
			if !ok ***REMOVED***
				return 0, 0, errors.New("VU counts must be integers")
			***REMOVED***
			return n, n, nil
		case 2:
			n1, ok1 := v[0].(int)
			n2, ok2 := v[1].(int)
			if !ok1 || !ok2 ***REMOVED***
				return 0, 0, errors.New("VU counts must be integers")
			***REMOVED***
			return n1, n2, nil
		default:
			return 0, 0, errors.New("Only one or two VU steps allowed per stage")
		***REMOVED***
	case nil:
		return 0, 0, nil
	default:
		return 0, 0, errors.New("VUs must be either an integer or [integer, integer]")
	***REMOVED***
***REMOVED***

func (c *Config) MakeTest() (t speedboat.Test, err error) ***REMOVED***
	t.Script = c.Script
	t.URL = c.URL
	if t.Script == "" && t.URL == "" ***REMOVED***
		return t, errors.New("Neither script nor URL specified")
	***REMOVED***

	fullDuration := 10 * time.Second
	if c.Duration != "" ***REMOVED***
		fullDuration, err = time.ParseDuration(c.Duration)
		if err != nil ***REMOVED***
			return t, err
		***REMOVED***
	***REMOVED***

	if len(c.Stages) > 0 ***REMOVED***
		var totalFluid int
		var totalFixed time.Duration

		for _, stage := range c.Stages ***REMOVED***
			tStage := speedboat.TestStage***REMOVED******REMOVED***

			switch v := stage.Duration.(type) ***REMOVED***
			case int:
				totalFluid += v
			case string:
				dur, err := time.ParseDuration(v)
				if err != nil ***REMOVED***
					return t, err
				***REMOVED***
				tStage.Duration = dur
				totalFixed += dur
			default:
				return t, errors.New("Stage durations must be integers or strings")
			***REMOVED***

			start, end, err := parseVUs(stage.VUs)
			if err != nil ***REMOVED***
				return t, err
			***REMOVED***
			tStage.StartVUs = start
			tStage.EndVUs = end

			t.Stages = append(t.Stages, tStage)
		***REMOVED***

		if totalFixed > fullDuration ***REMOVED***
			if totalFluid == 0 ***REMOVED***
				fullDuration = totalFixed
			***REMOVED*** else ***REMOVED***
				return t, errors.New("Stages exceed test duration")
			***REMOVED***
		***REMOVED***

		remainder := fullDuration - totalFixed
		if remainder > 0 ***REMOVED***
			for i, stage := range c.Stages ***REMOVED***
				chunk, ok := stage.Duration.(int)
				if !ok ***REMOVED***
					continue
				***REMOVED***
				t.Stages[i].Duration = time.Duration(chunk) / remainder
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		start, end, err := parseVUs(c.VUs)
		if err != nil ***REMOVED***
			return t, err
		***REMOVED***

		t.Stages = []speedboat.TestStage***REMOVED***
			speedboat.TestStage***REMOVED***Duration: fullDuration, StartVUs: start, EndVUs: end***REMOVED***,
		***REMOVED***
	***REMOVED***

	return t, nil
***REMOVED***
