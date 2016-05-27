package main

import (
	"errors"
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

/*func parseVUs(vus interface***REMOVED******REMOVED***) (VUSpec, error) ***REMOVED***
	switch v := vus.(type) ***REMOVED***
	case nil:
		return VUSpec***REMOVED******REMOVED***, nil
	case int:
		return VUSpec***REMOVED***Start: v, End: v***REMOVED***, nil
	case []interface***REMOVED******REMOVED***:
		switch len(v) ***REMOVED***
		case 1:
			v0, ok := v[0].(int)
			if !ok ***REMOVED***
				return VUSpec***REMOVED******REMOVED***, errors.New("Item in VU declaration is not an int")
			***REMOVED***
			return VUSpec***REMOVED***Start: v0, End: v0***REMOVED***, nil
		case 2:
			v0, ok0 := v[0].(int)
			v1, ok1 := v[1].(int)
			if !ok0 || !ok1 ***REMOVED***
				return VUSpec***REMOVED******REMOVED***, errors.New("Item in VU declaration is not an int")
			***REMOVED***
			return VUSpec***REMOVED***Start: v0, End: v1***REMOVED***, nil
		default:
			return VUSpec***REMOVED******REMOVED***, errors.New("Wrong number of values in [start, end] VU ramp")
		***REMOVED***
	default:
		return VUSpec***REMOVED******REMOVED***, errors.New("VUs must be either a single int or [start, end]")
	***REMOVED***
***REMOVED***

func (c *Config) Compile() (t LoadTest, err error) ***REMOVED***
	// Script/URL
	t.Script = c.Script
	t.URL = c.URL
	if t.Script == "" && t.URL == "" ***REMOVED***
		return t, errors.New("Script or URL must be specified")
	***REMOVED***

	// Root VU definitions
	rootVUs, err := parseVUs(c.VUs)
	if err != nil ***REMOVED***
		return t, err
	***REMOVED***

	// Duration
	rootDurationS := c.Duration
	if rootDurationS == "" ***REMOVED***
		rootDurationS = "10s"
	***REMOVED***
	rootDuration, err := time.ParseDuration(rootDurationS)
	if err != nil ***REMOVED***
		return t, err
	***REMOVED***

	// Stages
	if len(c.Stages) > 0 ***REMOVED***
		// Figure out the scale for flexible durations
		totalFluidDuration := 0
		totalFixedDuration := time.Duration(0)
		for i := 0; i < len(c.Stages); i++ ***REMOVED***
			switch v := c.Stages[i].Duration.(type) ***REMOVED***
			case int:
				totalFluidDuration += v
			case string:
				duration, err := time.ParseDuration(v)
				if err != nil ***REMOVED***
					return t, err
				***REMOVED***
				totalFixedDuration += duration
			***REMOVED***
		***REMOVED***

		// Make sure the fixed segments don't exceed the test length
		available := time.Duration(rootDuration.Nanoseconds() - totalFixedDuration.Nanoseconds())
		if available.Nanoseconds() < 0 ***REMOVED***
			return t, errors.New("Fixed stages are exceeding the test duration")
		***REMOVED***

		// Compile stage definitions
		for i := 0; i < len(c.Stages); i++ ***REMOVED***
			cStage := &c.Stages[i]
			stage := Stage***REMOVED******REMOVED***

			// Stage duration
			switch v := cStage.Duration.(type) ***REMOVED***
			case int:
				claim := float64(v) / float64(totalFluidDuration)
				stage.Duration = time.Duration(available.Seconds()*claim) * time.Second
			case string:
				stage.Duration, err = time.ParseDuration(v)
			***REMOVED***
			if err != nil ***REMOVED***
				return t, err
			***REMOVED***

			// VU curve
			stage.VUs, err = parseVUs(cStage.VUs)
			if err != nil ***REMOVED***
				return t, err
			***REMOVED***
			if stage.VUs.Start == 0 && stage.VUs.End == 0 ***REMOVED***
				if i > 0 ***REMOVED***
					stage.VUs = VUSpec***REMOVED***
						Start: t.Stages[i-1].VUs.End,
						End:   t.Stages[i-1].VUs.End,
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					stage.VUs = rootVUs
				***REMOVED***
			***REMOVED***

			t.Stages = append(t.Stages, stage)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Create an implicit, full-duration stage
		t.Stages = []Stage***REMOVED***
			Stage***REMOVED***
				Duration: rootDuration,
				VUs:      rootVUs,
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	return t, nil
***REMOVED****/
