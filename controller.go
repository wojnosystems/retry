// Copyright 2019 Chris Wojno
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
// documentation files (the "Software"), to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and
// to permit persons to whom the Software is furnished to do so, subject to the following conditions: The above
// copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN
// AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package retry

import "time"

type With struct {
	// Times is the maximum number of retries to allow before returning a failure
	Times int
	// WaitFor is the time to wait between failures
	WaitFor time.Duration
	// triesSoFar is how many times we've waited between failures & the number of failures that have occurred
	triesSoFar int
}

// ShouldTry will execute unless all of our retries allotted have failed
func (c *With)ShouldTry() bool {
	return c.triesSoFar < c.Times
}

// Wait will cause go to sleep for the WaitFor
func (c *With) WaitIfShouldTry() {
	c.triesSoFar++
	if c.ShouldTry() {
		time.Sleep(c.WaitFor)
	}
}

func (c *With) New() Controller {
	return &With{
		Times: c.Times,
		WaitFor: c.WaitFor,
	}
}