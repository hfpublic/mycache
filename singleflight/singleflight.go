package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mx sync.Mutex
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mx.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, has := g.m[key]; has {
		g.mx.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := &call{}
	c.wg.Add(1)
	g.m[key] = c
	g.mx.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mx.Lock()
	delete(g.m, key)
	g.mx.Unlock()

	return c.val, c.err
}
