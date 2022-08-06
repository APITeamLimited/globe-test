package pool

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
)

const (
	stateDefault = 0
	stateInited  = 1
	stateClosed  = 2
)

type BadConnError struct ***REMOVED***
	wrapped error
***REMOVED***

var _ error = (*BadConnError)(nil)

func (e BadConnError) Error() string ***REMOVED***
	s := "redis: Conn is in a bad state"
	if e.wrapped != nil ***REMOVED***
		s += ": " + e.wrapped.Error()
	***REMOVED***
	return s
***REMOVED***

func (e BadConnError) Unwrap() error ***REMOVED***
	return e.wrapped
***REMOVED***

//------------------------------------------------------------------------------

type StickyConnPool struct ***REMOVED***
	pool   Pooler
	shared int32 // atomic

	state uint32 // atomic
	ch    chan *Conn

	_badConnError atomic.Value
***REMOVED***

var _ Pooler = (*StickyConnPool)(nil)

func NewStickyConnPool(pool Pooler) *StickyConnPool ***REMOVED***
	p, ok := pool.(*StickyConnPool)
	if !ok ***REMOVED***
		p = &StickyConnPool***REMOVED***
			pool: pool,
			ch:   make(chan *Conn, 1),
		***REMOVED***
	***REMOVED***
	atomic.AddInt32(&p.shared, 1)
	return p
***REMOVED***

func (p *StickyConnPool) NewConn(ctx context.Context) (*Conn, error) ***REMOVED***
	return p.pool.NewConn(ctx)
***REMOVED***

func (p *StickyConnPool) CloseConn(cn *Conn) error ***REMOVED***
	return p.pool.CloseConn(cn)
***REMOVED***

func (p *StickyConnPool) Get(ctx context.Context) (*Conn, error) ***REMOVED***
	// In worst case this races with Close which is not a very common operation.
	for i := 0; i < 1000; i++ ***REMOVED***
		switch atomic.LoadUint32(&p.state) ***REMOVED***
		case stateDefault:
			cn, err := p.pool.Get(ctx)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if atomic.CompareAndSwapUint32(&p.state, stateDefault, stateInited) ***REMOVED***
				return cn, nil
			***REMOVED***
			p.pool.Remove(ctx, cn, ErrClosed)
		case stateInited:
			if err := p.badConnError(); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			cn, ok := <-p.ch
			if !ok ***REMOVED***
				return nil, ErrClosed
			***REMOVED***
			return cn, nil
		case stateClosed:
			return nil, ErrClosed
		default:
			panic("not reached")
		***REMOVED***
	***REMOVED***
	return nil, fmt.Errorf("redis: StickyConnPool.Get: infinite loop")
***REMOVED***

func (p *StickyConnPool) Put(ctx context.Context, cn *Conn) ***REMOVED***
	defer func() ***REMOVED***
		if recover() != nil ***REMOVED***
			p.freeConn(ctx, cn)
		***REMOVED***
	***REMOVED***()
	p.ch <- cn
***REMOVED***

func (p *StickyConnPool) freeConn(ctx context.Context, cn *Conn) ***REMOVED***
	if err := p.badConnError(); err != nil ***REMOVED***
		p.pool.Remove(ctx, cn, err)
	***REMOVED*** else ***REMOVED***
		p.pool.Put(ctx, cn)
	***REMOVED***
***REMOVED***

func (p *StickyConnPool) Remove(ctx context.Context, cn *Conn, reason error) ***REMOVED***
	defer func() ***REMOVED***
		if recover() != nil ***REMOVED***
			p.pool.Remove(ctx, cn, ErrClosed)
		***REMOVED***
	***REMOVED***()
	p._badConnError.Store(BadConnError***REMOVED***wrapped: reason***REMOVED***)
	p.ch <- cn
***REMOVED***

func (p *StickyConnPool) Close() error ***REMOVED***
	if shared := atomic.AddInt32(&p.shared, -1); shared > 0 ***REMOVED***
		return nil
	***REMOVED***

	for i := 0; i < 1000; i++ ***REMOVED***
		state := atomic.LoadUint32(&p.state)
		if state == stateClosed ***REMOVED***
			return ErrClosed
		***REMOVED***
		if atomic.CompareAndSwapUint32(&p.state, state, stateClosed) ***REMOVED***
			close(p.ch)
			cn, ok := <-p.ch
			if ok ***REMOVED***
				p.freeConn(context.TODO(), cn)
			***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	return errors.New("redis: StickyConnPool.Close: infinite loop")
***REMOVED***

func (p *StickyConnPool) Reset(ctx context.Context) error ***REMOVED***
	if p.badConnError() == nil ***REMOVED***
		return nil
	***REMOVED***

	select ***REMOVED***
	case cn, ok := <-p.ch:
		if !ok ***REMOVED***
			return ErrClosed
		***REMOVED***
		p.pool.Remove(ctx, cn, ErrClosed)
		p._badConnError.Store(BadConnError***REMOVED***wrapped: nil***REMOVED***)
	default:
		return errors.New("redis: StickyConnPool does not have a Conn")
	***REMOVED***

	if !atomic.CompareAndSwapUint32(&p.state, stateInited, stateDefault) ***REMOVED***
		state := atomic.LoadUint32(&p.state)
		return fmt.Errorf("redis: invalid StickyConnPool state: %d", state)
	***REMOVED***

	return nil
***REMOVED***

func (p *StickyConnPool) badConnError() error ***REMOVED***
	if v := p._badConnError.Load(); v != nil ***REMOVED***
		if err := v.(BadConnError); err.wrapped != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (p *StickyConnPool) Len() int ***REMOVED***
	switch atomic.LoadUint32(&p.state) ***REMOVED***
	case stateDefault:
		return 0
	case stateInited:
		return 1
	case stateClosed:
		return 0
	default:
		panic("not reached")
	***REMOVED***
***REMOVED***

func (p *StickyConnPool) IdleLen() int ***REMOVED***
	return len(p.ch)
***REMOVED***

func (p *StickyConnPool) Stats() *Stats ***REMOVED***
	return &Stats***REMOVED******REMOVED***
***REMOVED***
