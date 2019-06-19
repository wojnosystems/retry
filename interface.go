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

type Retrier interface {
	// This attempts to retry the function passed into it.
	// It will return nil if no error occurred, or an Errorer if at least 1 failure occurred.
	// A failure occurs when the function passed as a parameter returns a non-nil value
	// @param test the function to test for failure. Return an error to indicate a failure, nil to indicate success
	// @return nil if the function (eventually) succeeded, an Errorer containing a reference to the errors if not
	This(test func() error) Errorer
}

type Controller interface {
	// ShouldTry returns an indication that another retry should be attempted
	// Calling this method should NOT be interpreted as though a retry has occurred. Use Wait for that.
	ShouldTry() bool

	// WaitIfShouldTry should trigger some sort of wait state (unless ShouldTry returns false) that causes the CPU to
	// become free for a duration of time. When WaitIfShouldTry is called, it indicates that a retry was just
	// previously attempted, so you can use this to track how many retries have occurred
	WaitIfShouldTry()

	// New creates a clone of the controller ready for use with a new retry attempt
	New() Controller
}

type Errorer interface {
	// Errors gets every error returned by each try
	Errors() []error

	// Error returns this error and all nested errors as a string, for easy debugging
	Error() string

	// Last gets the very last error message generated by a failure, or nil if no error was encountered
	Last() error
}
