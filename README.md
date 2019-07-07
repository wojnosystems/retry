# Overview

A simple, but powerfully-configurable way to retry things that may fail.

## Linear, max tries-gated

The Retry.MaxAttempts creates a retry service that aborts once so many attempts have occurred. It has a linear time back-off and will wait 10 seconds after each failure in which a retry should follow (it will not trigger a wait state if the function will not perform an additional try).

```go
err := retry.How(&retry.MaxAttempts{
        Times: 3, 
        WaitFor: 10*time.Seconds,
	}.New()).This( func(controller retry.Controller)error {
	// Retry something!
	return errors.New("boom")
})
if err != nil {
	// log the last error encountered
	log.Error(err.Last())
	
	// log all errors using library's default message, which joins all errors together. This is sloppy and not intended for full production use, but it's your party.
	log.Error(err)
	
	// or loop through the slice of errors
	allErrors := err.Errors()
	for i := range allErrors {
		anError := allErrors[i]
		_ = anError // do something with it
	}
}
```

## Re-usable configuration

This allows you to create the configuration once, and re-use it for multiple instances. Calling .New() creates a new counter so that there is no shared state.

```go
// Re-usable
threeTimes := retry.MaxAttempts{
	Times: 3, 
	WaitFor: 10*time.Second,
}

err = retry.How(threeTimes.New()).This(func(controller retry.Controller)error {
	// Retry something!
	return nil
})
if err == nil {
	log.Println("no errors, hurray!")
}
```

## Hard/Fast fail

If you need to stop retrying for whatever reason, just tell controller to Abort.

```go
err = retry.How(threeTimes.New()).This(func(controller retry.Controller)error {
	// Only try once
	controller.Abort()
	return errors.New("just once")
})
if err != nil {
	// Only the "just once" error will appear and it will appear ONLY one time
	log.Error(err.Last())
}
```

## Exponential (Base2)

Linear waits are fine for many applications, but exponential waiting is very common requests for retry architectures.

```go
// Tries at most 7 times, with back-offs:
// 1st error: 10 seconds (total: 10 seconds)
// 2nd error: 20 seconds (total: 30 seconds)
// 3rd error: 40 seconds (total: 70 seconds)
// 4th error: 80 seconds (total: 150 seconds)
// 5th error: 160 seconds (total: 310 seconds)
// 6th error: 200 seconds (total: 510 seconds)
// 7th error: 200 seconds (total: 710 seconds)
// 
base2 := retry.ExpBase2{
	Times: 7, // maximum number of attempts to make, after which an error will be returned
	Scaling: 10*time.Second, // The starting point for our exponential back waits
	MaxAttemptWaitTime: 200*time.Second, // no single attempt may wait longer than 200 seconds after a failure
}
sevenErrorsReturned := retry.How(base2.New()).This(func(controller retry.Controller)error {

	return errors.New("boom!")
})
```

## Exponential With Context

The exponential configuration controls the time waiting for EACH attempt, but if you want to have a global timeout, send in your context configured with a global deadline.

```go
// Tries as many times as needed, but does not exceed 5 minutes
// 1st error: 10 seconds (total: 10 seconds)
// 2nd error: 20 seconds (total: 30 seconds)
// 3rd error: 40 seconds (total: 70 seconds/1.16 minutes)
// 4th error: 80 seconds (total: 150 seconds/2.5 minutes)
// 5th error: never occurs, context causes an abort
withTimeout := retry.ExpBase2{
	Scaling: 10*time.Second,
	// No need for MaxAttemptWaitTime
}
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
fourErrorsReturned := retry.How(withTimeout.NewWithContext(ctx)).This(func(controller retry.Controller)error {
	return errors.New("boom")
})
```

# Why is this so complicated

This module doesn't just implement a very basic Retry system, it also provides and implements a very extensible interface, too. You can build your own retry controllers (With) and your code shouldn't have to change except what you toss into How(). Neat, eh?

# Copyright

Copyright © 2019 Chris Wojno. All rights reserved.

No Warranties. Use this software at your own risk.

# License

Copyright 2019 Chris Wojno

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions: The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.