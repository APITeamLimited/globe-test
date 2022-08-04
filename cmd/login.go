package cmd

import (
	"github.com/spf13/cobra"
)

// getCmdLogin returns the `k6 login` sub-command, together with its children.
func getCmdLogin(gs *globalState) *cobra.Command ***REMOVED***
	loginCmd := &cobra.Command***REMOVED***
		Use:   "login",
		Short: "Authenticate with a service",
		Long: `Authenticate with a service.

Logging into a service changes the default when just "-o [type]" is passed with
no parameters, you can always override the stored credentials by passing some
on the commandline.`,
		RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
			return cmd.Usage()
		***REMOVED***,
	***REMOVED***
	loginCmd.AddCommand(
		getCmdLoginCloud(gs),
		getCmdLoginInfluxDB(gs),
	)

	return loginCmd
***REMOVED***
