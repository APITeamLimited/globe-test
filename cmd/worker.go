package cmd

import (
	"github.com/spf13/cobra"
	"go.k6.io/k6/worker"
)

func getCmdWorker(gs *globalState) *cobra.Command ***REMOVED***
	workerCmd := &cobra.Command***REMOVED***
		Use:   "worker",
		Short: "Runs a k6 worker to accept jobs",
		Long:  "Runs a k6 worker to accept jobs and execute them",
		Example: `
			k6 worker
		`,
		Run: func(cmd *cobra.Command, args []string) ***REMOVED***
			worker.Run()
		***REMOVED***,
	***REMOVED***

	return workerCmd
***REMOVED***
