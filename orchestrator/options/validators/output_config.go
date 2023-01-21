package validators

import (
	"fmt"

	"github.com/APITeamLimited/globe-test/lib/agent"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

func OutputConfig(options *libWorker.Options, funcMode, standalone bool) error {
	localhost := !funcMode && !standalone

	if !options.OutputConfig.Valid && options.ExecutionMode.Value == types.HTTPMultipleExecutionMode {
		options.OutputConfig = types.DefaultOutputConfig(localhost)
	}

	// Override global output config with localhost output config
	if localhost {
		for i := range options.OutputConfig.Value.Graphs {
			for j := range options.OutputConfig.Value.Graphs[i].Series {
				if options.OutputConfig.Value.Graphs[i].Series[j].LoadZone == libOrch.GlobalName {
					options.OutputConfig.Value.Graphs[i].Series[j].LoadZone = agent.AgentWorkerName
				}
			}
		}
	}

	options.OutputConfig.Value.Graphs = applyDefaultDesiredWidth(options.OutputConfig.Value.Graphs)

	err := validateMetricGraphs(options.OutputConfig.Value.Graphs)
	if err != nil {
		return err
	}

	return nil
}

func applyDefaultDesiredWidth(metricGraphs []types.MetricGraph) []types.MetricGraph {
	for i := range metricGraphs {
		if metricGraphs[i].DesiredWidth == 0 {
			metricGraphs[i].DesiredWidth = 1
		}
	}

	return metricGraphs
}

func validateMetricGraphs(metricGraphs []types.MetricGraph) error {
	for _, metricGraph := range metricGraphs {
		if metricGraph.Name == "" {
			return fmt.Errorf("metric graph name cannot be empty")
		}

		if metricGraph.DesiredWidth == 0 {
			metricGraph.DesiredWidth = 1
		}

		if metricGraph.DesiredWidth < 1 || metricGraph.DesiredWidth > 3 {
			return fmt.Errorf("metric graph desiredWidth must be between 1 and 3")
		}

		for _, series := range metricGraph.Series {

			if series.LoadZone == "" {
				return fmt.Errorf("loadZone cannot be empty")
			}

			if !types.IsValidSeriesKind(series.Kind) {
				return fmt.Errorf("metric graph kind must be one of %s, %s, %s",
					types.AreaGraphSeriesType, types.LineGraphSeriesType, types.ColumnGraphSeriesType)
			}

			if !types.IsBuiltinMetric(series.Metric) {
				return fmt.Errorf("metric must be one of the builtin metric types, got '%s'", series.Metric)
			}

			if !types.ValidSeriesColor(series.Color) {
				return fmt.Errorf("graph series colour invalid, got '%s'", series.Color)
			}
		}
	}

	return nil
}
