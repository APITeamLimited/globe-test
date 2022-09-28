package stats

import (
	"strconv"
	"time"
)

// LoadRawData parses and converts a slice of mixed data types to floats
func LoadRawData(raw interface***REMOVED******REMOVED***) (f Float64Data) ***REMOVED***
	var r []interface***REMOVED******REMOVED***
	var s Float64Data

	switch t := raw.(type) ***REMOVED***
	case []interface***REMOVED******REMOVED***:
		r = t
	case []uint:
		for _, v := range t ***REMOVED***
			s = append(s, float64(v))
		***REMOVED***
		return s
	case []uint8:
		for _, v := range t ***REMOVED***
			s = append(s, float64(v))
		***REMOVED***
		return s
	case []uint16:
		for _, v := range t ***REMOVED***
			s = append(s, float64(v))
		***REMOVED***
		return s
	case []uint32:
		for _, v := range t ***REMOVED***
			s = append(s, float64(v))
		***REMOVED***
		return s
	case []uint64:
		for _, v := range t ***REMOVED***
			s = append(s, float64(v))
		***REMOVED***
		return s
	case []bool:
		for _, v := range t ***REMOVED***
			if v == true ***REMOVED***
				s = append(s, 1.0)
			***REMOVED*** else ***REMOVED***
				s = append(s, 0.0)
			***REMOVED***
		***REMOVED***
		return s
	case []float64:
		return Float64Data(t)
	case []int:
		for _, v := range t ***REMOVED***
			s = append(s, float64(v))
		***REMOVED***
		return s
	case []int8:
		for _, v := range t ***REMOVED***
			s = append(s, float64(v))
		***REMOVED***
		return s
	case []int16:
		for _, v := range t ***REMOVED***
			s = append(s, float64(v))
		***REMOVED***
		return s
	case []int32:
		for _, v := range t ***REMOVED***
			s = append(s, float64(v))
		***REMOVED***
		return s
	case []int64:
		for _, v := range t ***REMOVED***
			s = append(s, float64(v))
		***REMOVED***
		return s
	case []string:
		for _, v := range t ***REMOVED***
			r = append(r, v)
		***REMOVED***
	case []time.Duration:
		for _, v := range t ***REMOVED***
			r = append(r, v)
		***REMOVED***
	case map[int]int:
		for i := 0; i < len(t); i++ ***REMOVED***
			s = append(s, float64(t[i]))
		***REMOVED***
		return s
	case map[int]int8:
		for i := 0; i < len(t); i++ ***REMOVED***
			s = append(s, float64(t[i]))
		***REMOVED***
		return s
	case map[int]int16:
		for i := 0; i < len(t); i++ ***REMOVED***
			s = append(s, float64(t[i]))
		***REMOVED***
		return s
	case map[int]int32:
		for i := 0; i < len(t); i++ ***REMOVED***
			s = append(s, float64(t[i]))
		***REMOVED***
		return s
	case map[int]int64:
		for i := 0; i < len(t); i++ ***REMOVED***
			s = append(s, float64(t[i]))
		***REMOVED***
		return s
	case map[int]string:
		for i := 0; i < len(t); i++ ***REMOVED***
			r = append(r, t[i])
		***REMOVED***
	case map[int]uint:
		for i := 0; i < len(t); i++ ***REMOVED***
			s = append(s, float64(t[i]))
		***REMOVED***
		return s
	case map[int]uint8:
		for i := 0; i < len(t); i++ ***REMOVED***
			s = append(s, float64(t[i]))
		***REMOVED***
		return s
	case map[int]uint16:
		for i := 0; i < len(t); i++ ***REMOVED***
			s = append(s, float64(t[i]))
		***REMOVED***
		return s
	case map[int]uint32:
		for i := 0; i < len(t); i++ ***REMOVED***
			s = append(s, float64(t[i]))
		***REMOVED***
		return s
	case map[int]uint64:
		for i := 0; i < len(t); i++ ***REMOVED***
			s = append(s, float64(t[i]))
		***REMOVED***
		return s
	case map[int]bool:
		for i := 0; i < len(t); i++ ***REMOVED***
			if t[i] == true ***REMOVED***
				s = append(s, 1.0)
			***REMOVED*** else ***REMOVED***
				s = append(s, 0.0)
			***REMOVED***
		***REMOVED***
		return s
	case map[int]float64:
		for i := 0; i < len(t); i++ ***REMOVED***
			s = append(s, t[i])
		***REMOVED***
		return s
	case map[int]time.Duration:
		for i := 0; i < len(t); i++ ***REMOVED***
			r = append(r, t[i])
		***REMOVED***
	***REMOVED***

	for _, v := range r ***REMOVED***
		switch t := v.(type) ***REMOVED***
		case int:
			a := float64(t)
			f = append(f, a)
		case uint:
			f = append(f, float64(t))
		case float64:
			f = append(f, t)
		case string:
			fl, err := strconv.ParseFloat(t, 64)
			if err == nil ***REMOVED***
				f = append(f, fl)
			***REMOVED***
		case bool:
			if t == true ***REMOVED***
				f = append(f, 1.0)
			***REMOVED*** else ***REMOVED***
				f = append(f, 0.0)
			***REMOVED***
		case time.Duration:
			f = append(f, float64(t))
		***REMOVED***
	***REMOVED***
	return f
***REMOVED***
