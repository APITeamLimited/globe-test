package aggregator

import (
	"fmt"
	"regexp"
)

// Aggregates intervals and corrects rates
func (aggregator *aggregator) aggregateIntervals(iwsfs []*intervalWithSubfraction, includeGlobalLocation bool) (*Interval, error) {
	// Get all possible locations
	locations := make([]string, 0)
	for _, iwsf := range iwsfs {
		containsLocation := false
		for _, location := range locations {
			if location == iwsf.location {
				containsLocation = true
				break
			}
		}

		if !containsLocation {
			locations = append(locations, iwsf.location)
		}
	}

	if includeGlobalLocation {
		locations = append(locations, "global")
	}

	iwsfByLocation := make([]*intervalWithSubfraction, len(locations))

	for i, location := range locations {
		intevalsToCombine := make([]*intervalWithSubfraction, 0)

		if location == "global" {
			intevalsToCombine = iwsfs
		} else {
			for _, iwsf := range iwsfs {
				if iwsf.location == location {
					intevalsToCombine = append(intevalsToCombine, iwsf)
				}
			}
		}

		combinedInterval, err := combinePeriodIntervals(intevalsToCombine, location)
		if err != nil {
			return nil, err
		}

		iwsfByLocation[i] = combinedInterval
	}

	sinkCount := 0
	for _, iwsf := range iwsfByLocation {
		sinkCount += len(iwsf.interval.Sinks)
	}

	unifiedInterval := Interval{
		Period: iwsfByLocation[0].interval.Period,
		Sinks:  make(map[string]*Sink, sinkCount),
	}

	for _, iwsf := range iwsfByLocation {
		for sinkName, sink := range iwsf.interval.Sinks {
			unifiedInterval.Sinks[sinkName] = sink
		}
	}

	return &unifiedInterval, nil
}

// Combines intervals frrom the same period into a single interval
func combinePeriodIntervals(iwsfs []*intervalWithSubfraction, sinkPrefix string) (*intervalWithSubfraction, error) {
	totalFraction := 0.0
	for _, iwsf := range iwsfs {
		totalFraction += iwsf.subFraction
	}

	newIwsf := intervalWithSubfraction{
		subFraction: totalFraction,
		interval: &Interval{
			Period: iwsfs[0].interval.Period,
			// Init sinks with the size of the first interval to minimise allocations
			// Future intervals will still require a reallocation if they have more sinks
			Sinks: make(map[string]*Sink, len(iwsfs[0].interval.Sinks)),
		},
	}

	var aggregatedSubfraction = 0.0

	for _, iwf := range iwsfs {
		for sinkName, sink := range iwf.interval.Sinks {
			prefixedName := fmt.Sprintf("%s::%s", sinkPrefix, sinkName)

			// Add sink if it doesn't exist
			aggregatedSink, ok := newIwsf.interval.Sinks[prefixedName]
			if !ok {
				aggregatedSink = &Sink{
					Type:   sink.Type,
					Labels: make(map[string]float64, len(sink.Labels)),
				}

				newIwsf.interval.Sinks[prefixedName] = aggregatedSink
			}

			for key, value := range sink.Labels {
				aggregatedValue, ok := aggregatedSink.Labels[key]
				if !ok {
					aggregatedValue = 0
				}

				newValue, err := combineSinkValues(sink.Type, key, aggregatedValue, aggregatedSubfraction, value, iwf.subFraction)
				if err != nil {
					// Unknown sink key, just ignore it as not critical
					continue
				}

				aggregatedSink.Labels[key] = newValue
			}
		}

		aggregatedSubfraction += iwf.subFraction
	}

	return &newIwsf, nil
}

var percentileRegex = regexp.MustCompile(`p\([1-9][0-9]?|100\)`)

// Combines sink values from multiple intervals
func combineSinkValues(sinkType SinkType, key string, value1, fraction1, value2, fraction2 float64) (float64, error) {
	switch sinkType {
	case SinkType_Counter:
		switch key {
		case "count":
			return value1 + value2, nil
		default:
			return 0, fmt.Errorf("unknown key in counter sink: '%s'", key)
		}
	case SinkType_Rate:
		switch key {
		case "pass", "total":
			return value1 + value2, nil
		case "rate":
			// Return mean of the rates
			return ((value1 * fraction1) + (value2 * fraction2)) / (fraction1 + fraction2), nil
		default:
			return 0, fmt.Errorf("unknown key in rate sink: '%s'", key)
		}
	case SinkType_Trend:
		switch key {
		case "avg", "med":
			return ((value1 * fraction1) + (value2 * fraction2)) / (fraction1 + fraction2), nil
		case "max":
			if value1 > value2 {
				return value1, nil
			}
			return value2, nil
		case "min":
			if value1 < value2 {
				return value1, nil
			}
			return value2, nil
		default:
			if percentileRegex.MatchString(key) {
				return ((value1 * fraction1) + (value2 * fraction2)) / (fraction1 + fraction2), nil
			}
			return 0, fmt.Errorf("unknown key in trend sink: '%s'", key)
		}
	case SinkType_Gauge:
		switch key {
		case "value":
			// Gauge just returns the latest value
			return value2, nil
		default:
			return 0, fmt.Errorf("unknown key in gauge sink: '%s'", key)
		}
	}

	return 0, fmt.Errorf("unknown sink type: '%s'", sinkType.String())
}

func calculateSinkRates(interval *Interval) (*Interval, error) {
	for _, sink := range interval.Sinks {
		if sink.Type != SinkType_Rate {
			continue
		}

		pass, ok := sink.Labels["pass"]
		if !ok {
			return nil, fmt.Errorf("missing 'pass' label in rate sink")
		}

		total, ok := sink.Labels["total"]
		if !ok {
			return nil, fmt.Errorf("missing 'total' label in rate sink")
		}

		sink.Labels["rate"] = pass / total
	}

	return interval, nil
}
