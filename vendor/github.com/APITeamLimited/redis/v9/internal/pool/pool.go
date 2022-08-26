package pool

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/APITeamLimited/redis/v9/internal"
)

var (
	// ErrClosed performs any operation on the closed client will return this error.
	ErrClosed = errors.New("redis: client is closed")

	// ErrPoolTimeout timed out waiting to get a connection from the connection pool.
	ErrPoolTimeout = errors.New("redis: connection pool timeout")
)

var timers = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		t := time.NewTimer(time.Hour)
		t.Stop()
		return t
	***REMOVED***,
***REMOVED***

// Stats contains pool state information and accumulated stats.
type Stats struct ***REMOVED***
	Hits     uint32 // number of times free connection was found in the pool
	Misses   uint32 // number of times free connection was NOT found in the pool
	Timeouts uint32 // number of times a wait timeout occurred

	TotalConns uint32 // number of total connections in the pool
	IdleConns  uint32 // number of idle connections in the pool
	StaleConns uint32 // number of stale connections removed from the pool
***REMOVED***

type Pooler interface ***REMOVED***
	NewConn(context.Context) (*Conn, error)
	CloseConn(*Conn) error

	Get(context.Context) (*Conn, error)
	Put(context.Context, *Conn)
	Remove(context.Context, *Conn, error)

	Len() int
	IdleLen() int
	Stats() *Stats

	Close() error
***REMOVED***

type Options struct ***REMOVED***
	Dialer  func(context.Context) (net.Conn, error)
	OnClose func(*Conn) error

	PoolFIFO        bool
	PoolSize        int
	PoolTimeout     time.Duration
	MinIdleConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
***REMOVED***

type lastDialErrorWrap struct ***REMOVED***
	err error
***REMOVED***

type ConnPool struct ***REMOVED***
	cfg *Options

	dialErrorsNum uint32 // atomic
	lastDialError atomic.Value

	queue chan struct***REMOVED******REMOVED***

	connsMu   sync.Mutex
	conns     []*Conn
	idleConns []*Conn

	poolSize     int
	idleConnsLen int

	stats Stats

	_closed  uint32 // atomic
	closedCh chan struct***REMOVED******REMOVED***
***REMOVED***

var _ Pooler = (*ConnPool)(nil)

func NewConnPool(opt *Options) *ConnPool ***REMOVED***
	p := &ConnPool***REMOVED***
		cfg: opt,

		queue:     make(chan struct***REMOVED******REMOVED***, opt.PoolSize),
		conns:     make([]*Conn, 0, opt.PoolSize),
		idleConns: make([]*Conn, 0, opt.PoolSize),
		closedCh:  make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	p.connsMu.Lock()
	p.checkMinIdleConns()
	p.connsMu.Unlock()

	return p
***REMOVED***

func (p *ConnPool) checkMinIdleConns() ***REMOVED***
	if p.cfg.MinIdleConns == 0 ***REMOVED***
		return
	***REMOVED***
	for p.poolSize < p.cfg.PoolSize && p.idleConnsLen < p.cfg.MinIdleConns ***REMOVED***
		p.poolSize++
		p.idleConnsLen++

		go func() ***REMOVED***
			err := p.addIdleConn()
			if err != nil && err != ErrClosed ***REMOVED***
				p.connsMu.Lock()
				p.poolSize--
				p.idleConnsLen--
				p.connsMu.Unlock()
			***REMOVED***
		***REMOVED***()
	***REMOVED***
***REMOVED***

func (p *ConnPool) addIdleConn() error ***REMOVED***
	cn, err := p.dialConn(context.TODO(), true)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	p.connsMu.Lock()
	defer p.connsMu.Unlock()

	// It is not allowed to add new connections to the closed connection pool.
	if p.closed() ***REMOVED***
		_ = cn.Close()
		return ErrClosed
	***REMOVED***

	p.conns = append(p.conns, cn)
	p.idleConns = append(p.idleConns, cn)
	return nil
***REMOVED***

func (p *ConnPool) NewConn(ctx context.Context) (*Conn, error) ***REMOVED***
	return p.newConn(ctx, false)
***REMOVED***

func (p *ConnPool) newConn(ctx context.Context, pooled bool) (*Conn, error) ***REMOVED***
	cn, err := p.dialConn(ctx, pooled)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	p.connsMu.Lock()
	defer p.connsMu.Unlock()

	// It is not allowed to add new connections to the closed connection pool.
	if p.closed() ***REMOVED***
		_ = cn.Close()
		return nil, ErrClosed
	***REMOVED***

	p.conns = append(p.conns, cn)
	if pooled ***REMOVED***
		// If pool is full remove the cn on next Put.
		if p.poolSize >= p.cfg.PoolSize ***REMOVED***
			cn.pooled = false
		***REMOVED*** else ***REMOVED***
			p.poolSize++
		***REMOVED***
	***REMOVED***

	return cn, nil
***REMOVED***

func (p *ConnPool) dialConn(ctx context.Context, pooled bool) (*Conn, error) ***REMOVED***
	if p.closed() ***REMOVED***
		return nil, ErrClosed
	***REMOVED***

	if atomic.LoadUint32(&p.dialErrorsNum) >= uint32(p.cfg.PoolSize) ***REMOVED***
		return nil, p.getLastDialError()
	***REMOVED***

	netConn, err := p.cfg.Dialer(ctx)
	if err != nil ***REMOVED***
		p.setLastDialError(err)
		if atomic.AddUint32(&p.dialErrorsNum, 1) == uint32(p.cfg.PoolSize) ***REMOVED***
			go p.tryDial()
		***REMOVED***
		return nil, err
	***REMOVED***

	cn := NewConn(netConn)
	cn.pooled = pooled
	return cn, nil
***REMOVED***

func (p *ConnPool) tryDial() ***REMOVED***
	for ***REMOVED***
		if p.closed() ***REMOVED***
			return
		***REMOVED***

		conn, err := p.cfg.Dialer(context.Background())
		if err != nil ***REMOVED***
			p.setLastDialError(err)
			time.Sleep(time.Second)
			continue
		***REMOVED***

		atomic.StoreUint32(&p.dialErrorsNum, 0)
		_ = conn.Close()
		return
	***REMOVED***
***REMOVED***

func (p *ConnPool) setLastDialError(err error) ***REMOVED***
	p.lastDialError.Store(&lastDialErrorWrap***REMOVED***err: err***REMOVED***)
***REMOVED***

func (p *ConnPool) getLastDialError() error ***REMOVED***
	err, _ := p.lastDialError.Load().(*lastDialErrorWrap)
	if err != nil ***REMOVED***
		return err.err
	***REMOVED***
	return nil
***REMOVED***

// Get returns existed connection from the pool or creates a new one.
func (p *ConnPool) Get(ctx context.Context) (*Conn, error) ***REMOVED***
	if p.closed() ***REMOVED***
		return nil, ErrClosed
	***REMOVED***

	if err := p.waitTurn(ctx); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for ***REMOVED***
		p.connsMu.Lock()
		cn, err := p.popIdle()
		p.connsMu.Unlock()

		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if cn == nil ***REMOVED***
			break
		***REMOVED***

		if !p.isHealthyConn(cn) ***REMOVED***
			_ = p.CloseConn(cn)
			continue
		***REMOVED***

		atomic.AddUint32(&p.stats.Hits, 1)
		return cn, nil
	***REMOVED***

	atomic.AddUint32(&p.stats.Misses, 1)

	newcn, err := p.newConn(ctx, true)
	if err != nil ***REMOVED***
		p.freeTurn()
		return nil, err
	***REMOVED***

	return newcn, nil
***REMOVED***

func (p *ConnPool) waitTurn(ctx context.Context) error ***REMOVED***
	select ***REMOVED***
	case <-ctx.Done():
		return ctx.Err()
	default:
	***REMOVED***

	select ***REMOVED***
	case p.queue <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
		return nil
	default:
	***REMOVED***

	timer := timers.Get().(*time.Timer)
	timer.Reset(p.cfg.PoolTimeout)

	select ***REMOVED***
	case <-ctx.Done():
		if !timer.Stop() ***REMOVED***
			<-timer.C
		***REMOVED***
		timers.Put(timer)
		return ctx.Err()
	case p.queue <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
		if !timer.Stop() ***REMOVED***
			<-timer.C
		***REMOVED***
		timers.Put(timer)
		return nil
	case <-timer.C:
		timers.Put(timer)
		atomic.AddUint32(&p.stats.Timeouts, 1)
		return ErrPoolTimeout
	***REMOVED***
***REMOVED***

func (p *ConnPool) freeTurn() ***REMOVED***
	<-p.queue
***REMOVED***

func (p *ConnPool) popIdle() (*Conn, error) ***REMOVED***
	if p.closed() ***REMOVED***
		return nil, ErrClosed
	***REMOVED***
	n := len(p.idleConns)
	if n == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	var cn *Conn
	if p.cfg.PoolFIFO ***REMOVED***
		cn = p.idleConns[0]
		copy(p.idleConns, p.idleConns[1:])
		p.idleConns = p.idleConns[:n-1]
	***REMOVED*** else ***REMOVED***
		idx := n - 1
		cn = p.idleConns[idx]
		p.idleConns = p.idleConns[:idx]
	***REMOVED***
	p.idleConnsLen--
	p.checkMinIdleConns()
	return cn, nil
***REMOVED***

func (p *ConnPool) Put(ctx context.Context, cn *Conn) ***REMOVED***
	if cn.rd.Buffered() > 0 ***REMOVED***
		internal.Logger.Printf(ctx, "Conn has unread data")
		p.Remove(ctx, cn, BadConnError***REMOVED******REMOVED***)
		return
	***REMOVED***

	if !cn.pooled ***REMOVED***
		p.Remove(ctx, cn, nil)
		return
	***REMOVED***

	var shouldCloseConn bool

	p.connsMu.Lock()

	if p.cfg.MaxIdleConns == 0 || p.idleConnsLen < p.cfg.MaxIdleConns ***REMOVED***
		p.idleConns = append(p.idleConns, cn)
		p.idleConnsLen++
	***REMOVED*** else ***REMOVED***
		p.removeConn(cn)
		shouldCloseConn = true
	***REMOVED***

	p.connsMu.Unlock()

	p.freeTurn()

	if shouldCloseConn ***REMOVED***
		_ = p.closeConn(cn)
	***REMOVED***
***REMOVED***

func (p *ConnPool) Remove(ctx context.Context, cn *Conn, reason error) ***REMOVED***
	p.removeConnWithLock(cn)
	p.freeTurn()
	_ = p.closeConn(cn)
***REMOVED***

func (p *ConnPool) CloseConn(cn *Conn) error ***REMOVED***
	p.removeConnWithLock(cn)
	return p.closeConn(cn)
***REMOVED***

func (p *ConnPool) removeConnWithLock(cn *Conn) ***REMOVED***
	p.connsMu.Lock()
	defer p.connsMu.Unlock()
	p.removeConn(cn)
***REMOVED***

func (p *ConnPool) removeConn(cn *Conn) ***REMOVED***
	for i, c := range p.conns ***REMOVED***
		if c == cn ***REMOVED***
			p.conns = append(p.conns[:i], p.conns[i+1:]...)
			if cn.pooled ***REMOVED***
				p.poolSize--
				p.checkMinIdleConns()
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *ConnPool) closeConn(cn *Conn) error ***REMOVED***
	if p.cfg.OnClose != nil ***REMOVED***
		_ = p.cfg.OnClose(cn)
	***REMOVED***
	return cn.Close()
***REMOVED***

// Len returns total number of connections.
func (p *ConnPool) Len() int ***REMOVED***
	p.connsMu.Lock()
	n := len(p.conns)
	p.connsMu.Unlock()
	return n
***REMOVED***

// IdleLen returns number of idle connections.
func (p *ConnPool) IdleLen() int ***REMOVED***
	p.connsMu.Lock()
	n := p.idleConnsLen
	p.connsMu.Unlock()
	return n
***REMOVED***

func (p *ConnPool) Stats() *Stats ***REMOVED***
	idleLen := p.IdleLen()
	return &Stats***REMOVED***
		Hits:     atomic.LoadUint32(&p.stats.Hits),
		Misses:   atomic.LoadUint32(&p.stats.Misses),
		Timeouts: atomic.LoadUint32(&p.stats.Timeouts),

		TotalConns: uint32(p.Len()),
		IdleConns:  uint32(idleLen),
		StaleConns: atomic.LoadUint32(&p.stats.StaleConns),
	***REMOVED***
***REMOVED***

func (p *ConnPool) closed() bool ***REMOVED***
	return atomic.LoadUint32(&p._closed) == 1
***REMOVED***

func (p *ConnPool) Filter(fn func(*Conn) bool) error ***REMOVED***
	p.connsMu.Lock()
	defer p.connsMu.Unlock()

	var firstErr error
	for _, cn := range p.conns ***REMOVED***
		if fn(cn) ***REMOVED***
			if err := p.closeConn(cn); err != nil && firstErr == nil ***REMOVED***
				firstErr = err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return firstErr
***REMOVED***

func (p *ConnPool) Close() error ***REMOVED***
	if !atomic.CompareAndSwapUint32(&p._closed, 0, 1) ***REMOVED***
		return ErrClosed
	***REMOVED***
	close(p.closedCh)

	var firstErr error
	p.connsMu.Lock()
	for _, cn := range p.conns ***REMOVED***
		if err := p.closeConn(cn); err != nil && firstErr == nil ***REMOVED***
			firstErr = err
		***REMOVED***
	***REMOVED***
	p.conns = nil
	p.poolSize = 0
	p.idleConns = nil
	p.idleConnsLen = 0
	p.connsMu.Unlock()

	return firstErr
***REMOVED***

func (p *ConnPool) isHealthyConn(cn *Conn) bool ***REMOVED***
	now := time.Now()

	if p.cfg.ConnMaxLifetime > 0 && now.Sub(cn.createdAt) >= p.cfg.ConnMaxLifetime ***REMOVED***
		return false
	***REMOVED***
	if p.cfg.ConnMaxIdleTime > 0 && now.Sub(cn.UsedAt()) >= p.cfg.ConnMaxIdleTime ***REMOVED***
		atomic.AddUint32(&p.stats.IdleConns, 1)
		return false
	***REMOVED***

	if connCheck(cn.netConn) != nil ***REMOVED***
		return false
	***REMOVED***

	cn.SetUsedAt(now)
	return true
***REMOVED***
