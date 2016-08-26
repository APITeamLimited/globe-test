package lib

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

type PoolTestVU struct ***REMOVED***
	ID int64
***REMOVED***

func (u *PoolTestVU) Reconfigure(id int64) error        ***REMOVED*** u.ID = id; return nil ***REMOVED***
func (u *PoolTestVU) RunOnce(ctx context.Context) error ***REMOVED*** return nil ***REMOVED***

func TestGetEmptyPool(t *testing.T) ***REMOVED***
	pool := VUPool***REMOVED***New: func() (VU, error) ***REMOVED*** return &PoolTestVU***REMOVED******REMOVED***, nil ***REMOVED******REMOVED***
	assert.Equal(t, 0, pool.Count())

	vu, _ := pool.Get()
	assert.IsType(t, &PoolTestVU***REMOVED******REMOVED***, vu)
	assert.Equal(t, 0, pool.Count())
***REMOVED***

func TestPutThenGet(t *testing.T) ***REMOVED***
	pool := VUPool***REMOVED***New: func() (VU, error) ***REMOVED*** return &PoolTestVU***REMOVED******REMOVED***, nil ***REMOVED******REMOVED***
	assert.Equal(t, 0, pool.Count())

	pool.Put(&PoolTestVU***REMOVED***ID: 1***REMOVED***)
	assert.Equal(t, 1, pool.Count())

	vu, _ := pool.Get()
	assert.IsType(t, &PoolTestVU***REMOVED******REMOVED***, vu)
	assert.Equal(t, int64(1), vu.(*PoolTestVU).ID)
	assert.Equal(t, 0, pool.Count())
***REMOVED***
