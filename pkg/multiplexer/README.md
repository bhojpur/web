# Bhojpur Web - Multiplexing Framework

It is a `multiplexing` library that relies on an underlying connection to provide reliability
and ordering, such as: `TCP` or `Unix domain sockets`, and provides stream-oriented multiplexing.

The [Bhojpur Web](https://github.com/bhojpur/web) - Multiplexer features include:

* Bi-directional streams
  * Streams can be opened by either `client` or `server`
  * Useful for NAT traversal
  * Server-side `push` support
* Flow control
  * Avoid starvation
  * Back-pressure to prevent overwhelming a receiver
* Keep Alives
  * Enables persistent connections over a load balancer
* Efficient
  * Enables thousands of logical streams with low overhead

## Simple Usage

```go
import(
    webmux "github.com/bhojpur/web/pkg/multiplexer"
)


func client() {
    // Get a TCP connection
    conn, err := net.Dial(...)
    if err != nil {
        panic(err)
    }

    // Setup client-side of Bhojpur Web multiplexer
    session, err := webmux.Client(conn, nil)
    if err != nil {
        panic(err)
    }

    // Open a new stream
    stream, err := session.Open()
    if err != nil {
        panic(err)
    }

    // Stream implements net.Conn
    stream.Write([]byte("ping"))
}

func server() {
    // Accept a TCP connection
    conn, err := listener.Accept()
    if err != nil {
        panic(err)
    }

    // Setup server-side of Bhojpur Web multiplexer
    session, err := webmux.Server(conn, nil)
    if err != nil {
        panic(err)
    }

    // Accept a stream
    stream, err := session.Accept()
    if err != nil {
        panic(err)
    }

    // Listen for a message
    buf := make([]byte, 4)
    stream.Read(buf)
}

```
