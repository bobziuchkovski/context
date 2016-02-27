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

/*
Package context provides a key/value storage with cancelation signaling.

Basics

Context provides both key/value storage as well as cancelation signaling.
New, empty contexts may be created via the New function:

	ctx = New()

A new context contains no key/value pairs and has no expiry associated with it.
Child contexts may be created with key/value pairs using Context.WithValue and
Context.WithFields.  Expiring children may be created using
Context.WithTimeout.

Cancelation

Contexts may be canceled either explicitly, using Context.Cancel, or implicitly
after a timeout using Context.WithTimeout.  When a Context is canceled, all
child Contexts are canceled at the same time, using the same cancelation reason
(either ErrCanceled or ErrTimeout).

A subscriber may wait for cancellation using Context.Terminated or
Context.Error.  Context.Terminated returns a channel that is closed when the
context is canceled.  Context.Error blocks until the context is canceled and
returns the cancelation reason.

	ctx := New()
	expiring := ctx.WithTimeout(time.Second)

	// This will block for a second
	<- expiring.Terminated()
*/
package context
