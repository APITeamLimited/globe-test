package lib

import (
	"context"
	"errors"
	"sync"
)

type VUPool struct ***REMOVED***
	New func() (VU, error)

	vus   []VU
	mutex sync.Mutex
***REMOVED***

func (p *VUPool) Get() (VU, error) ***REMOVED***
	p.mutex.Lock()
	defer p.mutex.Unlock()

	l := len(p.vus)
	if l == 0 ***REMOVED***
		return p.New()
	***REMOVED***

	vu := p.vus[l-1]
	p.vus = p.vus[:l-1]
	return vu, nil
***REMOVED***

func (p *VUPool) Put(vu VU) ***REMOVED***
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.vus = append(p.vus, vu)
***REMOVED***

func (p *VUPool) Count() int ***REMOVED***
	return len(p.vus)
***REMOVED***

type VUGroup struct ***REMOVED***
	Pool    VUPool
	RunOnce func(ctx context.Context, vu VU)

	ctx       context.Context
	cancelAll context.CancelFunc
	cancelers []context.CancelFunc
***REMOVED***

func (g *VUGroup) Start(ctx context.Context) ***REMOVED***
	g.ctx, g.cancelAll = context.WithCancel(ctx)
***REMOVED***

func (g *VUGroup) Stop() ***REMOVED***
	g.cancelAll()
	g.ctx = nil
***REMOVED***

func (g *VUGroup) Scale(count int) error ***REMOVED***
	if g.ctx == nil ***REMOVED***
		panic(errors.New("Group not running"))
	***REMOVED***

	for len(g.cancelers) < count ***REMOVED***
		vu, err := g.Pool.Get()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		id := int64(len(g.cancelers) + 1)
		if err := vu.Reconfigure(id); err != nil ***REMOVED***
			return err
		***REMOVED***

		ctx, cancel := context.WithCancel(g.ctx)
		g.cancelers = append(g.cancelers, cancel)

		go g.runVU(ctx, vu)
	***REMOVED***

	for len(g.cancelers) > count ***REMOVED***
		g.cancelers[len(g.cancelers)-1]()
		g.cancelers = g.cancelers[:len(g.cancelers)-1]
	***REMOVED***

	return nil
***REMOVED***

func (g *VUGroup) runVU(ctx context.Context, vu VU) ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return
		default:
			g.RunOnce(ctx, vu)
		***REMOVED***
	***REMOVED***
***REMOVED***
