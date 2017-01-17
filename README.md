# svctools-go

A little Go package for decorating mundane service calls with retry, backoff, jitter, and standard instrumentation.

Many parts of this library are protocol agnostic. For example the retry, backoff, jitter, and instrumentation functions can be used with any type of function, let alone underlying protocol or package being used.

### Installation

    go get github.com/buildertools/svctools-go

### Usage

This is a simple library with a few handy functions. Everything is opt-in. Here are a few highlights...

````
// Attempt to fetch http://someawesomeservice.com/v1/whatever for 30 seconds, 
// retrying on IO errors and 500 level status codes. Retry using an exponential
// backoff strategy with a base backoff of 50 millis and between 0 and 50 
// millis of jitter.
r, err := clients.RetryExponential(
	func() {
		return clients.WrapHttpResponseError(http.Get(`http://someawsomeservice.com/v1/whatever`))
	},
	time.Duration(30)*time.Second,
	time.Duration(50)*time.Millisecond,
	time.Duration(50)*time.Millisecond)

// Ping a local TCP socket every second for 30 seconds, fail if ping fails
_, e := clients.RetryPeriodic(
	func() {
		conn, err := net.Dial("tcp", "localhost:3000")
		defer conn.Close()
		if err != nil {
			return nil, NonRetirableError{E:err}
		}
		return nil, RetirableError{E:err}
	},
	time.Duration(30)*time.Second,
	time.Duration(1)*time.Second,
	time.Duration(0))

// Mix and match your own backoff and jitter tooling
r, e := clients.Retry(
	yourRetriableFunction,
	&JitteredBackoff{
		TTL:       time.Duration(1)*time.Second,
		Initial:   time.Duration(10)*time.Millisecond,
		MaxJitter: time.Duration(5)*time.Millisecond,
		Bof:       func(round uint, t time.Duration) time.Duration {
				// your impl
		           },
		Jf:        Jitter
	})

````

A user can provide their own implementation of the PerishableWaiter interface for even more control over the backoff semantics and implementation.


