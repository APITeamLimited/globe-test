// Copyright 2014 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package internal // import "github.com/garyburd/redigo/internal"

import (
	"strings"
)

const (
	WatchState = 1 << iota
	MultiState
	SubscribeState
	MonitorState
)

type CommandInfo struct ***REMOVED***
	Set, Clear int
***REMOVED***

var commandInfos = map[string]CommandInfo***REMOVED***
	"WATCH":      ***REMOVED***Set: WatchState***REMOVED***,
	"UNWATCH":    ***REMOVED***Clear: WatchState***REMOVED***,
	"MULTI":      ***REMOVED***Set: MultiState***REMOVED***,
	"EXEC":       ***REMOVED***Clear: WatchState | MultiState***REMOVED***,
	"DISCARD":    ***REMOVED***Clear: WatchState | MultiState***REMOVED***,
	"PSUBSCRIBE": ***REMOVED***Set: SubscribeState***REMOVED***,
	"SUBSCRIBE":  ***REMOVED***Set: SubscribeState***REMOVED***,
	"MONITOR":    ***REMOVED***Set: MonitorState***REMOVED***,
***REMOVED***

func init() ***REMOVED***
	for n, ci := range commandInfos ***REMOVED***
		commandInfos[strings.ToLower(n)] = ci
	***REMOVED***
***REMOVED***

func LookupCommandInfo(commandName string) CommandInfo ***REMOVED***
	if ci, ok := commandInfos[commandName]; ok ***REMOVED***
		return ci
	***REMOVED***
	return commandInfos[strings.ToUpper(commandName)]
***REMOVED***
