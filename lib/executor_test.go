/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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

package lib

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExecutionStateVUIDs(t *testing.T) ***REMOVED***
	es := NewExecutionState(Options***REMOVED******REMOVED***, 0, 0) //TODO: fix
	assert.Equal(t, uint64(1), es.GetUniqueVUIdentifier())
	assert.Equal(t, uint64(2), es.GetUniqueVUIdentifier())
	assert.Equal(t, uint64(3), es.GetUniqueVUIdentifier())
	wg := sync.WaitGroup***REMOVED******REMOVED***
	rand.Seed(time.Now().UnixNano())
	count := rand.Intn(50)
	wg.Add(count)
	for i := 0; i < count; i++ ***REMOVED***
		go func() ***REMOVED***
			es.GetUniqueVUIdentifier()
			wg.Done()
		***REMOVED***()
	***REMOVED***
	wg.Wait()
	assert.Equal(t, uint64(4+count), es.GetUniqueVUIdentifier())
***REMOVED***

//TODO: way more tests...
