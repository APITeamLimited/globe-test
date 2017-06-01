/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package dummy

import (
	"context"
	"sync"

	"github.com/loadimpact/k6/stats"
)

type Collector struct ***REMOVED***
	Samples []stats.Sample
	running bool

	lock sync.Mutex
***REMOVED***

func (c *Collector) Init() error ***REMOVED*** return nil ***REMOVED***

func (c *Collector) MakeConfig() interface***REMOVED******REMOVED*** ***REMOVED*** return nil ***REMOVED***

func (c *Collector) Run(ctx context.Context) ***REMOVED***
	c.lock.Lock()
	c.running = true
	c.lock.Unlock()

	<-ctx.Done()

	c.lock.Lock()
	c.running = false
	c.lock.Unlock()
***REMOVED***

func (c *Collector) Collect(samples []stats.Sample) ***REMOVED***
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.running ***REMOVED***
		panic("attempted to collect while not running")
	***REMOVED***
	c.Samples = append(c.Samples, samples...)
***REMOVED***

func (c *Collector) IsRunning() bool ***REMOVED***
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.running
***REMOVED***
