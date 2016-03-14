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
	"errors"
	"runtime"
	"time"
)

var nilPairs = (*pairs)(nil)

var (
	// ErrTimeout is returned from Context.Error() when the context is canceled
	// due to timeout.
	ErrTimeout = errors.New("context timed-out")

	// ErrCanceled is returned from Context.Error() when the context is
	// canceled via Context.Cancel().
	ErrCanceled = errors.New("context cancelled")
)

// Fields represents a set of key/value pairs within a Context.
type Fields map[interface{}]interface{}

// Context carries key/value pairs and provides signaling capabilities to
// cancel unneeded work.
type Context interface {
	// WithFields returns a new Context with fields added to the existing
	// key/value pairs.
	WithFields(fields Fields) Context

	// WithValue returns a new Context with key and value added to the existing
	// key/value pairs.
	WithValue(key interface{}, value interface{}) Context

	// GetValue returns the value associated with key, setting present to
	// true if the value is found.
	GetValue(key interface{}) (value interface{}, present bool)

	// WithTimeout returns a new Context with an expiry.  If timeout elapses
	// and the context has not already been canceled, it is canceled with
	// ErrTimeout.  If an ancestor context has an existing expiration prior to
	// the timeout parameter, the ancestor's timeout is inherited instead.
	WithTimeout(timeout time.Duration) Context

	// TimeRemaining returns the time remaining before the Context expires.
	// If a timeout has not been set for the context or any of its ancestors,
	// timeoutPresent is false.  If the timeout has already elapsed,
	// timeoutPresent is true and remaining is set to zero.
	TimeRemaining() (remaining time.Duration, timeoutPresent bool)

	// Cancel cancels the current context and all child contexts with
	// ErrCanceled.  The call is a no-op for contexts that have already been
	// canceled.
	Cancel()

	// Terminated returns a channel that is closed when the current context
	// is canceled.
	Terminated() <-chan struct{}

	// Error returns the cancelation reason for the context.  It blocks until
	// the context is canceled.
	Error() error
}

type context struct {
	// Data
	parent *context
	fields *pairs

	// Signaling
	trigger *errorTrigger

	// Expiry
	timer    *time.Timer
	deadline time.Time
}

func (c *context) newChild() *context {
	parent := c
	child := &context{
		parent:  parent,
		fields:  parent.fields,
		trigger: newTrigger(),
		// timer: should not be inherited
		deadline: parent.deadline,
	}
	parent.trigger.Subscribe(child.trigger)

	// When child is garbage collected, we unsubscribe the child trigger.
	// Otherwise the parent might keep a reference to the trigger and leak mem.
	runtime.SetFinalizer(child, func(garbage *context) {
		if garbage.parent != nil && garbage.parent.trigger != nil {
			garbage.parent.trigger.Unsubscribe(garbage.trigger)
		}
	})
	return child
}

// New returns a new Context.
func New() Context {
	return &context{
		parent:  nil,
		fields:  nilPairs,
		trigger: newTrigger(),
	}
}

func (c *context) GetValue(key interface{}) (value interface{}, present bool) {
	return c.fields.Get(key)
}

func (c *context) WithFields(fields Fields) Context {
	var new Context = c
	for k, v := range fields {
		new = new.WithValue(k, v)
	}
	return new
}

func (c *context) WithValue(key interface{}, value interface{}) Context {
	parent := c
	child := parent.newChild()
	child.fields = parent.fields.append(key, value)
	return child
}

func (c *context) Cancel() {
	if c.timer != nil {
		c.timer.Stop()
	}
	c.trigger.Fire(ErrCanceled)
}

func (c *context) Error() error {
	<-c.Terminated()
	return c.trigger.Get()
}

func (c *context) Terminated() <-chan struct{} {
	return c.trigger.Fired()
}

func (c *context) WithTimeout(timeout time.Duration) Context {
	deadline := time.Now().Add(timeout)
	parent := c
	child := parent.newChild()
	child.deadline = deadline
	if !parent.deadline.IsZero() && deadline.After(parent.deadline) {
		child.deadline = parent.deadline
	} else {
		child.timer = time.AfterFunc(timeout, func() { child.trigger.Fire(ErrTimeout) })
	}
	return child
}

func (c *context) TimeRemaining() (remaining time.Duration, timeoutPresent bool) {
	if c.deadline.IsZero() {
		return
	}
	if c.trigger.Get() != nil {
		return
	}
	timeoutPresent = true
	now := time.Now()
	if now.Before(c.deadline) {
		remaining = c.deadline.Sub(now)
	}
	return
}
