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
	"testing"
	"time"

	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
)

func TestCollectorRun(t *testing.T) ***REMOVED***
	var wg sync.WaitGroup
	c := &Collector***REMOVED******REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() ***REMOVED***
		defer wg.Done()
		c.Run(ctx)
	***REMOVED***()
	cancel()
	wg.Wait()
***REMOVED***

func TestCollectorCollect(t *testing.T) ***REMOVED***
	c := &Collector***REMOVED******REMOVED***
	t.Run("no context", func(t *testing.T) ***REMOVED***
		assert.Panics(t, func() ***REMOVED*** c.Collect([]stats.Sample***REMOVED******REMOVED******REMOVED******REMOVED***) ***REMOVED***)
	***REMOVED***)
	t.Run("context", func(t *testing.T) ***REMOVED***
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() ***REMOVED*** c.Run(ctx) ***REMOVED***()
		for !c.IsRunning() ***REMOVED***
			time.Sleep(1 * time.Millisecond)
		***REMOVED***
		c.Collect([]stats.Sample***REMOVED******REMOVED******REMOVED******REMOVED***)
		assert.Len(t, c.Samples, 1)
	***REMOVED***)
***REMOVED***
