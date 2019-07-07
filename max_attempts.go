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

import "time"

// MaxAttempts will try a function X number of times, waiting WaitFor in between each failed attempt
type MaxAttempts struct {
	// MaxAttempts is the maximum number of retries to allow before returning a failure
	Times uint
	// WaitFor is the time to wait between failures
	WaitFor time.Duration
}

// New creates a new MaxAttempts. New is needed to create a counter state required for this invocation
func (l MaxAttempts) New() Service {
	return &maxExponentialService{
		config: Exponential{
			Times:   l.Times,
			YOffset: l.WaitFor,
			Base:    0,
		},
	}
}
