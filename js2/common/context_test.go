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

package common

import (
	"context"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
)

func TestContextState(t *testing.T) ***REMOVED***
	st := &State***REMOVED******REMOVED***
	assert.Equal(t, st, GetState(WithState(context.Background(), st)))
***REMOVED***

func TestContextStateNil(t *testing.T) ***REMOVED***
	assert.Nil(t, GetState(context.Background()))
***REMOVED***

func TestContextRuntime(t *testing.T) ***REMOVED***
	rt := goja.New()
	assert.Equal(t, rt, GetRuntime(WithRuntime(context.Background(), rt)))
***REMOVED***

func TestContextRuntimeNil(t *testing.T) ***REMOVED***
	assert.Nil(t, GetRuntime(context.Background()))
***REMOVED***
