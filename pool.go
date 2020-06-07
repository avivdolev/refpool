// Package refpool implements a resource pool with reference counting.
// Based on sync.Pool and sync/atomic
package refpool

import (
	"sync"
	"sync/atomic"
)

// Element is a resource that holds a reference counter
type Element interface {
	// Counter should return pointer to an integer for atomic reference counting.
	// Pool methods will use sync.Atomic functions to modify the counter.
	Counter() *int64
}

// Refpool is the main package type
type Refpool struct {
	p sync.Pool
}

// IncElement adds n to e.Counter() and return the updated value.
func (rp *Refpool) IncElement(e Element, n int64) int64 {
	return atomic.AddInt64(e.Counter(), n)
}

// SetElement value of counter
func (rp *Refpool) SetElement(e Element, n int64) {
	atomic.StoreInt64(e.Counter(), n)
}

// New returns a new Refpool.
// The argument new should be a function which allocates a new element to be returned from the pool if it's empty.
func New(new func() Element) *Refpool {
	return &Refpool{
		p: sync.Pool{
			New: func() interface{} {
				return new()
			},
		},
	}
}

// Get an element from the pool, allocated if needed.
// The element's counter is set to 0.
func (rp *Refpool) Get() Element {
	e := rp.p.Get().(Element)
	rp.SetElement(e, 0)
	return e
}

// Put an element back into the pool.
// First lower the counter by 1, and only put it if the value is 0 or lower.
func (rp *Refpool) Put(e Element) {
	if rp.IncElement(e, -1) > 0 {
		return
	}
	rp.p.Put(e)
}

// Drop put the element back regardless its counter.
func (rp *Refpool) Drop(e Element) {
	rp.p.Put(e)
}
