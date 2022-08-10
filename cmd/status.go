package cmd

import (
	"github.com/spf13/cobra"

	"go.k6.io/k6/api/v1/client"
)

func getCmdStatus(globalState *globalState) *cobra.Command ***REMOVED***
	// statusCmd represents the status command
	statusCmd := &cobra.Command***REMOVED***
		Use:   "status",
		Short: "Show test status",
		Long: `Show test status.

  Use the global --address flag to specify the URL to the API server.`,
		RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
			c, err := client.New(globalState.flags.address)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			status, err := c.Status(globalState.ctx)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			return yamlPrint(globalState.stdOut, status)
		***REMOVED***,
	***REMOVED***
	return statusCmd
***REMOVED***
