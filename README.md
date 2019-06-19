# Overview

A simple, but powerfully-configurable way to retry things that may fail.

```go
err := retry.How(&retry.With{Times: 3, WaitFor: 10*time.Seconds}).This( func()error {
	// Retry something!
	return errors.New("boom")
})
if err != nil {
	// log the last error encountered
	log.Error(err.Last())
}


// Re-usable
threeTimes := retry.With{Times: 3, WaitFor: 10*time.Second}

err = retry.How(&threeTimes).This(func()error {
	// Retry something!
	return nil
})
if err == nil {
	log.Println("no errors, hurray!")
}
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