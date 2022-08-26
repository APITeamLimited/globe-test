package pool

import (
	"bufio"
	"context"
	"net"
	"sync/atomic"
	"time"

	"github.com/APITeamLimited/redis/v9/internal/proto"
)

var noDeadline = time.Time***REMOVED******REMOVED***

type Conn struct ***REMOVED***
	usedAt  int64 // atomic
	netConn net.Conn

	rd *proto.Reader
	bw *bufio.Writer
	wr *proto.Writer

	Inited    bool
	pooled    bool
	createdAt time.Time
***REMOVED***

func NewConn(netConn net.Conn) *Conn ***REMOVED***
	cn := &Conn***REMOVED***
		netConn:   netConn,
		createdAt: time.Now(),
	***REMOVED***
	cn.rd = proto.NewReader(netConn)
	cn.bw = bufio.NewWriter(netConn)
	cn.wr = proto.NewWriter(cn.bw)
	cn.SetUsedAt(time.Now())
	return cn
***REMOVED***

func (cn *Conn) UsedAt() time.Time ***REMOVED***
	unix := atomic.LoadInt64(&cn.usedAt)
	return time.Unix(unix, 0)
***REMOVED***

func (cn *Conn) SetUsedAt(tm time.Time) ***REMOVED***
	atomic.StoreInt64(&cn.usedAt, tm.Unix())
***REMOVED***

func (cn *Conn) SetNetConn(netConn net.Conn) ***REMOVED***
	cn.netConn = netConn
	cn.rd.Reset(netConn)
	cn.bw.Reset(netConn)
***REMOVED***

func (cn *Conn) Write(b []byte) (int, error) ***REMOVED***
	return cn.netConn.Write(b)
***REMOVED***

func (cn *Conn) RemoteAddr() net.Addr ***REMOVED***
	if cn.netConn != nil ***REMOVED***
		return cn.netConn.RemoteAddr()
	***REMOVED***
	return nil
***REMOVED***

func (cn *Conn) WithReader(ctx context.Context, timeout time.Duration, fn func(rd *proto.Reader) error) error ***REMOVED***
	if timeout != 0 ***REMOVED***
		if err := cn.netConn.SetReadDeadline(cn.deadline(ctx, timeout)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return fn(cn.rd)
***REMOVED***

func (cn *Conn) WithWriter(
	ctx context.Context, timeout time.Duration, fn func(wr *proto.Writer) error,
) error ***REMOVED***
	if timeout != 0 ***REMOVED***
		if err := cn.netConn.SetWriteDeadline(cn.deadline(ctx, timeout)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if cn.bw.Buffered() > 0 ***REMOVED***
		cn.bw.Reset(cn.netConn)
	***REMOVED***

	if err := fn(cn.wr); err != nil ***REMOVED***
		return err
	***REMOVED***

	return cn.bw.Flush()
***REMOVED***

func (cn *Conn) Close() error ***REMOVED***
	return cn.netConn.Close()
***REMOVED***

func (cn *Conn) deadline(ctx context.Context, timeout time.Duration) time.Time ***REMOVED***
	tm := time.Now()
	cn.SetUsedAt(tm)

	if timeout > 0 ***REMOVED***
		tm = tm.Add(timeout)
	***REMOVED***

	if ctx != nil ***REMOVED***
		deadline, ok := ctx.Deadline()
		if ok ***REMOVED***
			if timeout == 0 ***REMOVED***
				return deadline
			***REMOVED***
			if deadline.Before(tm) ***REMOVED***
				return deadline
			***REMOVED***
			return tm
		***REMOVED***
	***REMOVED***

	if timeout > 0 ***REMOVED***
		return tm
	***REMOVED***

	return noDeadline
***REMOVED***
