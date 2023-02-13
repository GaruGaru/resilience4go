# resilience4go  [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov]

Lightweight fault tolerance library written in Go Inspired by [resilience4j](https://resilience4j.readme.io/docs) and [hystrix](https://github.com/Netflix/Hystrix)

## Usage 

```go
import (
    "github.com/garugaru/resilience4go/circuitbreaker"
    "github.com/garugaru/resilience4go/fallback"
    "github.com/garugaru/resilience4go/resilience"
    "github.com/garugaru/resilience4go/retry"
)

var executor = resilience.New[string](
    // retry 10 times with a 10ms delay
    retry.New[string](retry.NewFixedDelay(10*time.Millisecond), 10),
    // circuit breaker based on error rate for the last 10 requests
    circuitbreaker.New[string](circuitbreaker.NewCountBasedWindow(10), 0.5, 0.5, 100*time.Millisecond),
    // if everything else fails, return a fallback value instead of error
    fallback.New[string]("fallback"),
)

out, err := executor.Execute(func() (string, error) {
        // failing function 
	return callVeryBadExternalService()
})
```

## Features

### Retries 

Retry a failing function a fixed amount of times using customizable delay strategy.
Available strategies are: 

* FixedDelay
* ExponentialBackoff


### Circuit breaker 

A [circuit breaker](https://martinfowler.com/bliki/CircuitBreaker.html) is a pattern used to protect external dependencies from being overwhelmed with requests in case of 
failures.

The circuit breaker takes in account the samples for the last N calls calculating the success / error rate, once the error
rate exceeds a user-defined threshold the requests are discarded for a fixed amount of time giving the external service the 
time to be restored. 
After a customizable amount of time a percentage of the requests will be forwarded to the service in order to check if it's 
still unavailable, once the service operability is restored the circuit will be closed and all the requests will be delivered to 
the service.

### Fallback

Return a fallback value instead of the error. 


[ci-img]: https://github.com/garugaru/resilience4go/actions/workflows/tests.yml/badge.svg
[cov-img]: https://codecov.io/gh/garugaru/resilience4go/branch/master/graph/badge.svg
[ci]: https://github.com/garugaru/resilience4go/actions/workflows/tests.yml
[cov]: https://codecov.io/gh/garugaru/resilience4go