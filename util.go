package main

import (
	"gopkg.in/guregu/null.v3"
	"gopkg.in/urfave/cli.v1"
	"time"
)

// cliBool returns a CLI argument as a bool, which is invalid if not given.
func cliBool(cc *cli.Context, name string) null.Bool ***REMOVED***
	return null.NewBool(cc.Bool(name), cc.IsSet(name))
***REMOVED***

// cliInt64 returns a CLI argument as an int64, which is invalid if not given.
func cliInt64(cc *cli.Context, name string) null.Int ***REMOVED***
	return null.NewInt(cc.Int64(name), cc.IsSet(name))
***REMOVED***

// cliFloat64 returns a CLI argument as a float64, which is invalid if not given.
func cliFloat64(cc *cli.Context, name string) null.Float ***REMOVED***
	return null.NewFloat(cc.Float64(name), cc.IsSet(name))
***REMOVED***

// cliDuration returns a CLI argument as a duration string, which is invalid if not given.
func cliDuration(cc *cli.Context, name string) null.String ***REMOVED***
	return null.NewString(cc.Duration(name).String(), cc.IsSet(name))
***REMOVED***

func roundDuration(d, to time.Duration) time.Duration ***REMOVED***
	return d - (d % to)
***REMOVED***
