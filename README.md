# go-backpressure

This repository demonstrates `backpressure` implementation with Go buffered channels.

As you may know, more concurrrency does not make programs faster, It is counterintuitive but systems perform better overall when their components limit the amount of work they are willing to perform. Go programs can spawn hundreds, thousands, even tens of thousands of simultaneous goroutines. but the general rule is to use as little concurrency as your program need otherwise your code becomes harder to understand

The `RequestGauge` struct uses a buffered channel to limit the amount of simultaneous requests our application is handling to prevent it from falling behind and becoming overwhelmed.

Each empty spot in our buffered channel represent a running goroutine. Every time a goroutine calls `RequestGauge.run`, either of these scenarios occurs:

- successfully read from the channel thereby removing an element from the channel. It then invokes its callback and writes a token back to the channel when the callback exits, allowing other goroutine to run when they call the `run` function.
- the default case in our select statement runs indicating the buffered channel is empty (there are still running goroutines).

```go
type RequestGauge struct {
    channel chan struct{}
}

func (g *RequestGauge) run(callback func()) error {
    select {
    case <-g.channel:
      callback()
      g.channel <- struct{}{}
      return nil

    default:
      return errors.New("request gauge capacity exceeded")
    }
}
```

The code in `main.go` uses this implementation with the built-in HTTP server

```go
func getHandler(requestGauge *RequestGauge) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    err := requestGauge.run(func() {
      sendJson(w, Payload{Error: "", ServerTime: getServerTimeAfter(requestDuration)})
    })

    if err != nil {
      w.WriteHeader(http.StatusTooManyRequests)
      sendJson(w, Payload{Error: "too many request"})
    }
  }
}
```

## Docker

run this example locally with docker

build docker image

```bash
docker build -t go-backpressure .
```

run docker image

```bash
docker run -it -p 4000:4000 go-backpressure
```

## Contributing

if you have any suggestion or find an issue with this implementation, please help out
