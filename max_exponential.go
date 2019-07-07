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
	"math"
	"time"
)

// Exponential performs retries and waits according to an exponential equation: f(x) = a*B^x + y where B is the base, x is the number of retries (starting at 0) and y is the offset time (cannot be less than 0)
type Exponential struct {
	// Times is the maximum number of times to execute the back-off
	Times uint

	// Base is the base of the exponent
	Base time.Duration

	// YOffset is the y-component of the exponential equation
	YOffset time.Duration

	// Scaling is multiplication factor, this controls the units of the exponentiation
	Scaling time.Duration

	// MaxAttemptWaitTime The maximum amount of time to wait for a particular attempt (does not account for total time), regardless of the exponential equation (leave as 0 to ignore)
	MaxAttemptWaitTime time.Duration
}

func (l Exponential) New() Service {
	return &maxExponentialService{
		config: l,
	}
}

// ExpBase2 configures the Exponential but with base 2 instead of an arbitrary base.
type ExpBase2 struct {
	// Times is the maximum number of times to execute the back-off
	Times uint

	// YOffset is the y-component of the exponential equation
	YOffset time.Duration

	// Scaling is multiplication factor, this controls the units of the exponentiation
	Scaling time.Duration

	// MaxAttemptWaitTime The maximum amount of time to wait for a particular attempt (does not account for total time), regardless of the exponential equation (leave as 0 to ignore)
	MaxAttemptWaitTime time.Duration
}

func (l ExpBase2) New() Service {
	return &maxExponentialService{
		config: Exponential{
			Times:              l.Times,
			YOffset:            l.YOffset,
			Scaling:            l.Scaling,
			MaxAttemptWaitTime: l.MaxAttemptWaitTime,
			Base:               2,
		},
	}
}

type maxExponentialService struct {
	config     Exponential
	triesSoFar uint
}

// ShouldTry will execute unless all of our retries allotted have failed
func (c *maxExponentialService) ShouldTry() bool {
	return c.triesSoFar < c.config.Times
}

// Wait will cause go to sleep for the WaitFor
func (c *maxExponentialService) Yield() {
	time.Sleep(c.waitDuration())
}

func (c maxExponentialService) waitDuration() time.Duration {
	// implement the equation: a*B^x + y, and constrain to the bounds, if necessary
	var waitFor time.Duration
	switch c.config.Base {
	case 0:
		// do nothing, no exponent, constant function
	case 1:
		// constant with scaling
		waitFor = 1
	case 2:
		// cheat for base 2
		waitFor = 1 << (c.triesSoFar - 1)
	default:
		// Only perform power if necessary for other bases
		waitFor = time.Duration(math.Pow(float64(c.config.Base), math.Max(0.0, float64(c.triesSoFar-1))))
	}

	// scaling factor
	waitFor = waitFor * c.config.Scaling

	// add the YOffset
	waitFor = waitFor + c.config.YOffset

	if c.config.MaxAttemptWaitTime != 0 && waitFor > c.config.MaxAttemptWaitTime {
		// Constrain the wait time
		waitFor = c.config.MaxAttemptWaitTime
	}
	return waitFor
}

// Returns the svc for the service so that the developer can svc it
func (c *maxExponentialService) Controller() ServiceController {
	return c
}

// Wait will cause go to sleep for the WaitFor
func (c *maxExponentialService) Abort() {
	c.triesSoFar = c.config.Times
}

// Wait will cause go to sleep for the WaitFor
func (c *maxExponentialService) NotifyRetry() {
	c.triesSoFar++
}

// NewErrorList creates the default error list
func (c *maxExponentialService) NewErrorList() ErrorAppender {
	return newErrorList()
}
