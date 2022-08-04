package cmd

import (
	"github.com/spf13/cobra"

	"go.k6.io/k6/api/v1/client"
)

func getCmdStats(globalState *globalState) *cobra.Command ***REMOVED***
	// statsCmd represents the stats command
	statsCmd := &cobra.Command***REMOVED***
		Use:   "stats",
		Short: "Show test metrics",
		Long: `Show test metrics.

  Use the global --address flag to specify the URL to the API server.`,
		RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
			c, err := client.New(globalState.flags.address)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			metrics, err := c.Metrics(globalState.ctx)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			return yamlPrint(globalState.stdOut, metrics)
		***REMOVED***,
	***REMOVED***
	return statsCmd
***REMOVED***
