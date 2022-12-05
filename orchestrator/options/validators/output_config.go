package validators

import (
	"fmt"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

func OutputConfig(options *libWorker.Options) error ***REMOVED***
	if !options.OutputConfig.Valid && options.ExecutionMode.Value == types.HTTPMultipleExecutionMode ***REMOVED***
		options.OutputConfig = types.DefaultOutputConfig()
	***REMOVED***

	err := validateMetricGraphs(options.OutputConfig.Value.Graphs)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	localhost output config

	return nil
***REMOVED***

func validateMetricGraphs(metricGraphs []types.MetricGraph) error ***REMOVED***
	for _, metricGraph := range metricGraphs ***REMOVED***
		if metricGraph.Name == "" ***REMOVED***
			return fmt.Errorf("metric graph name cannot be empty")
		***REMOVED***

		if metricGraph.DesiredWidth < 1 || metricGraph.DesiredWidth > 3 ***REMOVED***
			return fmt.Errorf("metric graph desiredWidth must be between 1 and 3")
		***REMOVED***

		for _, series := range metricGraph.Series ***REMOVED***

			if series.LoadZone == "" ***REMOVED***
				return fmt.Errorf("loadZone cannot be empty")
			***REMOVED***

			if !types.IsValidSeriesKind(series.Kind) ***REMOVED***
				return fmt.Errorf("metric graph kind must be one of %s, %s, %s",
					types.AreaGraphSeriesType, types.LineGraphSeriesType, types.ColumnGraphSeriesType)
			***REMOVED***

			if !types.IsBuiltinMetric(series.Metric) ***REMOVED***
				return fmt.Errorf("metric must be one of the builtin metric types, got '%s'", series.Metric)
			***REMOVED***

			if !types.ValidSeriesColor(series.Color) ***REMOVED***
				return fmt.Errorf("graph series colour invalid, got '%s'", series.Color)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
