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

import (
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	caseCounter := make(map[string]int)

	cases := map[string]struct {
		test     func() error
		control  Controller
		expected Errorer
	}{
		"always fail 3 times": {
			test: func() error {
				return errors.New("boom")
			},
			control: &With{
				Times:   3,
				WaitFor: 1 * time.Microsecond,
			},
			expected: func() *errorList {
				el := newErrorList()
				el.append(errors.New("boom"))
				el.append(errors.New("boom"))
				el.append(errors.New("boom"))
				return el
			}(),
		},
		"error twice, then succeed": {
			test: func() error {
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
			control: &With{
				Times:   3,
				WaitFor: 1 * time.Microsecond,
			},
			expected: nil,
		},
	}

	for caseName, c := range cases {
		actualError := How(c.control).This(c.test)
		if c.expected != nil {
			if actualError.Error() != c.expected.Error() {
				t.Errorf(`%s: errors did not match; expected "%s", got: "%s"`, caseName, c.expected.Error(), actualError.Error())
			}
			if len(actualError.Errors()) != len(c.expected.Errors()) {
				t.Errorf(`%s: number of errors did not match; expected "%d", got: "%d"`, caseName, len(c.expected.Errors()), len(actualError.Errors()))
			}
			if actualError.Last().Error() != c.expected.Last().Error() {
				t.Errorf(`%s: last error did not match; expected "%s", got: "%s"`, caseName, c.expected.Last().Error(), actualError.Last().Error())
			}
		} else {
			if actualError != nil {
				t.Errorf("%s: expected error to be nil", caseName)
			}
		}
	}
}
