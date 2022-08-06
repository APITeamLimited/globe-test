package pool

import "context"

type SingleConnPool struct ***REMOVED***
	pool      Pooler
	cn        *Conn
	stickyErr error
***REMOVED***

var _ Pooler = (*SingleConnPool)(nil)

func NewSingleConnPool(pool Pooler, cn *Conn) *SingleConnPool ***REMOVED***
	return &SingleConnPool***REMOVED***
		pool: pool,
		cn:   cn,
	***REMOVED***
***REMOVED***

func (p *SingleConnPool) NewConn(ctx context.Context) (*Conn, error) ***REMOVED***
	return p.pool.NewConn(ctx)
***REMOVED***

func (p *SingleConnPool) CloseConn(cn *Conn) error ***REMOVED***
	return p.pool.CloseConn(cn)
***REMOVED***

func (p *SingleConnPool) Get(ctx context.Context) (*Conn, error) ***REMOVED***
	if p.stickyErr != nil ***REMOVED***
		return nil, p.stickyErr
	***REMOVED***
	return p.cn, nil
***REMOVED***

func (p *SingleConnPool) Put(ctx context.Context, cn *Conn) ***REMOVED******REMOVED***

func (p *SingleConnPool) Remove(ctx context.Context, cn *Conn, reason error) ***REMOVED***
	p.cn = nil
	p.stickyErr = reason
***REMOVED***

func (p *SingleConnPool) Close() error ***REMOVED***
	p.cn = nil
	p.stickyErr = ErrClosed
	return nil
***REMOVED***

func (p *SingleConnPool) Len() int ***REMOVED***
	return 0
***REMOVED***

func (p *SingleConnPool) IdleLen() int ***REMOVED***
	return 0
***REMOVED***

func (p *SingleConnPool) Stats() *Stats ***REMOVED***
	return &Stats***REMOVED******REMOVED***
***REMOVED***
