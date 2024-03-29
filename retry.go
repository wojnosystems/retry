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

type basic struct {
	svc Service
}

// How creates a new retry service
func How(svc Service) Retrier {
	return &basic{
		svc: svc,
	}
}

// This invokes the developer's method to retry
func (b *basic) This(test func(controller ServiceController) error) Errorer {
	var errorList ErrorAppender
	// Retry until we should not
	for true {
		// Perform the action under test, this is the thing the developer would like to retry
		err := test(b.svc.Controller())
		// Notify our service that the try/retry has occurred
		b.svc.NotifyRetry()
		if err != nil {
			if errorList == nil {
				// factory a new error list, if not yet created (lazy-create)
				errorList = b.svc.NewErrorList()
			}
			// Got an error, record it
			errorList.Append(err)
			// Wait, but only if we should try again
			if b.svc.ShouldTry() {
				b.svc.Yield()
			} else {
				return errorList
			}
		} else {
			// success, no need to retry
			return nil
		}
	}
	// should be unreachable, but go is complaining that there is no return
	return nil
}
