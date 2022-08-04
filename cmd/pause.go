package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/guregu/null.v3"

	v1 "go.k6.io/k6/api/v1"
	"go.k6.io/k6/api/v1/client"
)

func getCmdPause(globalState *globalState) *cobra.Command ***REMOVED***
	// pauseCmd represents the pause command
	pauseCmd := &cobra.Command***REMOVED***
		Use:   "pause",
		Short: "Pause a running test",
		Long: `Pause a running test.

  Use the global --address flag to specify the URL to the API server.`,
		RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
			c, err := client.New(globalState.flags.address)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			status, err := c.SetStatus(globalState.ctx, v1.Status***REMOVED***
				Paused: null.BoolFrom(true),
			***REMOVED***)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			return yamlPrint(globalState.stdOut, status)
		***REMOVED***,
	***REMOVED***
	return pauseCmd
***REMOVED***
