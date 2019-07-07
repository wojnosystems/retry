// Copyright 2019 Chris Wojno
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
// documentation files (the "Software"), to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and
// to permit persons to whom the Software is furnished to do so, subject to the following conditions: The above
// copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// WARRANTIES OF MERCHANTABILITY, FITNESS FOR Scaling PARTICULAR PURPOSE AND NON-INFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN
// AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package retry

import (
	"context"
	"time"
)

// NewWithContext defines a way to set a context to determine a maximum attempt among all calls
// the context limits the total amount of time yielded, regardless of each invocation of the test
func (l Exponential) NewWithContext(ctx context.Context) Service {
	return &maxExponentialContextService{
		maxExponentialService: maxExponentialService{
			config: l,
		},
		ctx: ctx,
	}
}

func (l ExpBase2) NewWithContext(ctx context.Context) Service {
	mes := l.New().(*maxExponentialService)
	return &maxExponentialContextService{
		maxExponentialService: *mes,
		ctx:                   ctx,
	}
}

type maxExponentialContextService struct {
	maxExponentialService
	ctx context.Context
}

// Wait will cause go to sleep for the WaitFor
func (c *maxExponentialContextService) Yield() {
	waitFor := c.waitDuration()
	select {
	case <-c.ctx.Done():
		// context is done, abort, never yield
		c.Abort()
	case <-time.After(waitFor):
		// time expired, ok to proceed
	}
}
