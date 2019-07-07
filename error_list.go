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

import "strings"

type errorList struct {
	recordedErrors []error
}

func newErrorList() *errorList {
	return &errorList{
		recordedErrors: make([]error, 0),
	}
}

func (e *errorList) Errors() []error {
	return e.recordedErrors
}

// Error is a default message. Override to make your own.
func (e *errorList) Error() string {
	errList := make([]string, len(e.recordedErrors))
	for i, err := range e.recordedErrors {
		errList[i] = err.Error()
	}
	return "retries exceeded; encountered errors: " + strings.Join(errList, ", ")
}

// Last gets the final error returned by the retry, this is special because it's usually the one you want
func (e *errorList) Last() error {
	if len(e.recordedErrors) == 0 {
		return nil
	}
	return e.recordedErrors[len(e.recordedErrors)-1]
}

// Append adds an error to the list
func (e *errorList) Append(err error) {
	e.recordedErrors = append(e.recordedErrors, err)
}
