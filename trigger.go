// Copyright (c) 2016 Bob Ziuchkovski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package context

import (
	"sync"
)

type triggerSubscriber interface {
	Fire(value error)
}

type empty struct{}

var emptyValue = empty{}

type errorTrigger struct {
	mu sync.Mutex

	value       error
	waiter      chan struct{}
	subscribers map[triggerSubscriber]empty
}

func newTrigger() *errorTrigger {
	return &errorTrigger{
		waiter:      make(chan struct{}),
		subscribers: make(map[triggerSubscriber]empty),
	}
}

func (t *errorTrigger) Subscribe(s triggerSubscriber) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.value != nil {
		s.Fire(t.value)
		return
	}
	t.subscribers[s] = emptyValue
	return
}

func (t *errorTrigger) Unsubscribe(s triggerSubscriber) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.subscribers, s)
}

func (t *errorTrigger) Fire(value error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.value != nil {
		// Only fire once
		return
	}

	t.value = value
	for subscriber := range t.subscribers {
		subscriber.Fire(value)
	}
	t.subscribers = nil
	close(t.waiter)
}

func (t *errorTrigger) Fired() <-chan struct{} {
	return t.waiter
}

func (t *errorTrigger) Get() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.value
}
