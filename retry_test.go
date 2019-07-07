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
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	caseCounter := make(map[string]int)

	cases := map[string]struct {
		test     func(controller ServiceController) error
		cfg      MaxAttempts
		expected Errorer
	}{
		"always fail 3 times": {
			test: func(controller ServiceController) error {
				return errors.New("boom")
			},
			cfg: MaxAttempts{
				Times:   3,
				WaitFor: 1 * time.Microsecond,
			},
			expected: func(controller ServiceController) *errorList {
				el := newErrorList()
				el.Append(errors.New("boom"))
				el.Append(errors.New("boom"))
				el.Append(errors.New("boom"))
				return el
			}(nil),
		},
		"error twice, then succeed": {
			test: func(controller ServiceController) error {
				var c int
				var ok bool
				if c, ok = caseCounter["error twice, then succeed"]; !ok {
					caseCounter["error twice, then succeed"] = 1
				} else {
					caseCounter["error twice, then succeed"] = caseCounter["error twice, then succeed"] + 1
				}
				if c < 2 {
					return errors.New("boom")
				}
				return nil
			},
			cfg: MaxAttempts{
				Times:   3,
				WaitFor: 1 * time.Microsecond,
			},
			expected: nil,
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actualError := How(c.cfg.New()).This(c.test)
			if c.expected != nil {
				if actualError.Error() != c.expected.Error() {
					t.Errorf(`errors did not match; expected "%s", got: "%s"`, c.expected.Error(), actualError.Error())
				}
				if len(actualError.Errors()) != len(c.expected.Errors()) {
					t.Errorf(`number of errors did not match; expected "%d", got: "%d"`, len(c.expected.Errors()), len(actualError.Errors()))
				}
				if actualError.Last().Error() != c.expected.Last().Error() {
					t.Errorf(`last error did not match; expected "%s", got: "%s"`, c.expected.Last().Error(), actualError.Last().Error())
				}
			} else {
				if actualError != nil {
					t.Errorf("expected error to be nil")
				}
			}
		})
	}
}

// TestMaxExponentialService_Abort_Attempts ensures that abort prevents calling the methods after abort is called
func TestRetry_Abort(t *testing.T) {
	expectedError := errors.New("boom")
	ma := MaxAttempts{Times: 3, WaitFor: 10 * time.Second}
	errList := How(ma.New()).This(func(controller ServiceController) error {
		controller.Abort()
		return expectedError
	})

	if errList == nil {
		t.Error("expected an error list")
	} else {
		if errList.Last() != expectedError {
			t.Errorf(`expected error: "%s" but got: "%s"`, expectedError.Error(), errList.Last().Error())
		}
	}
}

// TestMaxAttemptsService_Abort_NoWait ensures that Yield is not called after Abort is called
func TestRetry_Abort_NoYield(t *testing.T) {
	ma := MaxAttempts{Times: 3, WaitFor: 10 * time.Second}
	startTime := time.Now()
	expectedEndTime := startTime.Add(ma.WaitFor)
	_ = How(ma.New()).This(func(controller ServiceController) error {
		controller.Abort()
		return nil
	})

	endTime := time.Now()

	// The test waited after abort, it should NOT do this
	if expectedEndTime.Before(endTime) || expectedEndTime.Equal(endTime) {
		t.Error("retry waited after abort was called and should not have")
	}
}
