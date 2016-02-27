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

// Using a linked list is bad for lookup efficiency, but I'm not sure n will be
// large enough to matter for most use cases.  This could be an immutable hash
// trie or red-black tree or something if needed.
type pairs struct {
	prev  *pairs
	key   interface{}
	value interface{}
}

func (l *pairs) append(key interface{}, value interface{}) *pairs {
	return &pairs{
		prev:  l,
		key:   key,
		value: value,
	}
}

func (l *pairs) Get(key interface{}) (value interface{}, present bool) {
	switch {
	case l == nil:
		// No-op
	case l.key == key:
		value, present = l.value, true
	default:
		value, present = l.prev.Get(key)
	}
	return
}
