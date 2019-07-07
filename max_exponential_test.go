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
	"testing"
	"time"
)

// TestMaxExponential_New tests the wait-time generation for the exponential function
func TestMaxExponential_New(t *testing.T) {
	cases := map[string]struct {
		cfg      Exponential
		expected []time.Duration
	}{
		"linear": {
			cfg: Exponential{
				Base:    0,
				YOffset: 10 * time.Second,
			},
			expected: []time.Duration{10 * time.Second, 10 * time.Second, 10 * time.Second, 10 * time.Second},
		},
		"exponential: 1ns*2^x": {
			cfg: Exponential{
				Base:    1,
				Scaling: 3,
			},
			expected: []time.Duration{3 * time.Nanosecond, 3 * time.Nanosecond, 3 * time.Nanosecond, 3 * time.Nanosecond, 3 * time.Nanosecond},
		},
		"exponential: 1sec*2^x": {
			cfg: Exponential{
				Base:    2,
				Scaling: 1 * time.Second,
			},
			expected: []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 16 * time.Second},
		},
		"exponential: 10sec*2^x": {
			cfg: Exponential{
				Base:    2,
				Scaling: 10 * time.Second,
			},
			expected: []time.Duration{10 * time.Second, 20 * time.Second, 40 * time.Second, 80 * time.Second, 160 * time.Second},
		},
		"exponential: 10sec*3^x": {
			cfg: Exponential{
				Base:    3,
				Scaling: 10 * time.Second,
			},
			expected: []time.Duration{10 * time.Second, 30 * time.Second, 90 * time.Second, 270 * time.Second},
		},
		"exponential: 10sec*2^x + 10sec": {
			cfg: Exponential{
				Base:    2,
				Scaling: 10 * time.Second,
				YOffset: 10 * time.Second,
			},
			expected: []time.Duration{20 * time.Second, 30 * time.Second, 50 * time.Second, 90 * time.Second, 170 * time.Second},
		},
		"exponential: 10sec*2^x + 10sec (limit 50)": {
			cfg: Exponential{
				Base:               2,
				Scaling:            10 * time.Second,
				YOffset:            10 * time.Second,
				MaxAttemptWaitTime: 50 * time.Second,
			},
			expected: []time.Duration{20 * time.Second, 30 * time.Second, 50 * time.Second, 50 * time.Second, 50 * time.Second},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actualTimes := make([]time.Duration, len(c.expected))
			svcIFace := c.cfg.New()
			svc := svcIFace.(*maxExponentialService)
			for i := 0; i < len(c.expected); i++ {
				svc.NotifyRetry()
				actualTimes[i] = svc.waitDuration()
			}
			for i := range c.expected {
				if c.expected[i] != actualTimes[i] {
					t.Errorf(`expected duration: %v but got %v`, c.expected[i], actualTimes[i])
				}
			}
		})
	}
}

func TestExpBase2_New(t *testing.T) {
	base2Svc := ExpBase2{
		Times:              3,
		Scaling:            10 * time.Second,
		MaxAttemptWaitTime: 100 * time.Second,
		YOffset:            5 * time.Second,
	}.New()
	base2 := base2Svc.(*maxExponentialService)
	if base2.config.Base != 2 {
		t.Error("Expected base to be 2, but got ", base2.config.Base)
	}
}
