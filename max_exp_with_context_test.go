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
	"testing"
	"time"
)

// TestMaxExponentialWithContext_New ensures that the context governs the maximum run time of a retry
func TestMaxExponentialWithContext_New(t *testing.T) {
	cases := map[string]struct {
		cfg         Exponential
		ctxDuration time.Duration
		expected    []time.Duration
	}{
		"linear": {
			cfg: Exponential{
				Base:    0,
				YOffset: 10 * time.Microsecond,
			},
			ctxDuration: 31 * time.Microsecond,
			expected:    []time.Duration{10 * time.Microsecond, 10 * time.Microsecond, 10 * time.Microsecond},
		},
		"exponential: 10us*2^x + 10us (limit 50 with max exec context of 205ms)": {
			cfg: Exponential{
				Base:             2,
				Scaling:          10 * time.Microsecond,
				YOffset:          10 * time.Microsecond,
				AttemptTimeLimit: 50 * time.Microsecond,
			},
			ctxDuration: 205 * time.Microsecond,
			expected:    []time.Duration{20 * time.Microsecond, 30 * time.Microsecond, 50 * time.Microsecond, 50 * time.Microsecond},
		},
		"exponential: 10us*2^x + 10us (limit 50 with max exec context of 105ms)": {
			cfg: Exponential{
				Base:             2,
				Scaling:          10 * time.Microsecond,
				YOffset:          10 * time.Microsecond,
				AttemptTimeLimit: 50 * time.Microsecond,
			},
			ctxDuration: 105 * time.Microsecond,
			expected:    []time.Duration{20 * time.Microsecond, 30 * time.Microsecond, 50 * time.Microsecond},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), c.ctxDuration)
			svcInterface := c.cfg.NewWithContext(ctx)
			svc := svcInterface.(*maxExponentialContextService)
			startAt := time.Now()
			for i := 0; i < len(c.expected); i++ {
				svc.Yield()
			}
			cancel()
			endAt := time.Now()
			expectedEndAt := startAt.Add(c.ctxDuration)
			if expectedEndAt.Add(1 * time.Microsecond).After(endAt) {
				t.Errorf(`exceeded context timeout, over by %v`, endAt.Sub(expectedEndAt))
			}
		})
	}
}
