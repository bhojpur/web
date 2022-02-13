# Bhojpur Web - Virtual Hosting

It is a very simple library that lets you implement `virtual hosting` features in [Bhojpur Web](https://github.com/bhojpur/web) for different protocols (e.g., HTTP and TLS/SSL). The `virtual` hosting library has a __high-level__ and a __low-level__ interface.

The __high-level__ `virtual` hosting library interface lets you wrap existing `net.Listeners` with the `muxer` objects. You can then `Listen()` on a muxer for a particular virtual host name of interest which will return to you a `net.Listener` for just connections with the virtual hostname of interest.

The __lower-level__ `virtual` hosting library interface are just functions which extract the name/routing information for the given protocol and return an object implementing net.Conn that works as if no bytes had been consumed.

## Simple HTTP Multiplexing

```go
import (
	"net"
	"github.com/bhojpur/web/pkg/virtual"
)

primary, _ := net.Listen("tcp", *listen)

// start multiplexing on the primary TCP listener
mux, _ := virtual.NewHTTPMuxer(primary, muxTimeout)

// listen for connections to different domains
for _, v := range virtualHosts {
	vhost := v

	// vhost.Name is a virtual hostname like "test.bhojpur.net"
	muxListener, _ := mux.Listen(vhost.Name())

	go func(vh virtualHost, ml net.Listener) {
		for {
			conn, _ := ml.Accept()
			go vh.Handle(conn)
		}
	}(vhost, muxListener)
}

for {
	conn, err := mux.NextError()

	switch err.(type) {
	case vhost.BadRequest:
		log.Printf("got a bad request!")
		conn.Write([]byte("bad request"))
	case vhost.NotFound:
		log.Printf("got a connection for an unknown virtual host")
		conn.Write([]byte("virtual host not found"))
	case vhost.Closed:
		log.Printf("closed conn: %s", err)
	default:
		if conn != nil {
			conn.Write([]byte("server error"))
		}
	}

	if conn != nil {
		conn.Close()
	}
}
```

## Low-level Virtual Host API for HTTP

```go
import (
	"net"
	"github.com/bhojpur/web/pkg/virtual"
)

primary, _ := net.Listen("tcp", *listen)

// accept a new connection
conn, _ := primary.Accept()

// parse out the HTTP request and the Host header
if vhostConn, err = virtual.HTTP(conn); err != nil {
	panic("Not a valid HTTP connection!")
}

fmt.Printf("Target Host: ", vhostConn.Host())
// Target Host: bhojpur.net

// vhostConn contains the entire request as if no bytes had been consumed
bytes, _ := ioutil.ReadAll(vhostConn)
fmt.Printf("%s", bytes)
// GET / HTTP/1.1
// Host: bhojpur.net
// User-Agent: ...
// ...
```

### Advanced Introspection

The entire HTTP request headers are available for inspection, in case you want
to multiplex on something besides the Host header:

```go
import (
	"net"
	"github.com/bhojpur/web/pkg/virtual"
)

primary, _ := net.Listen("tcp", *listen)

// accept a new connection
conn, _ := primary.Accept()

// parse out the HTTP request and the Host header
if vhostConn, err = virtual.HTTP(conn); err != nil {
	panic("Not a valid HTTP connection!")
}

httpVersion := virtual.Request.MinorVersion
customRouting := virtual.Request.Header["X-Custom-Routing-Header"]
```

Likewise for `TLS`, you can look at detailed information about the `ClientHello` message:

```go
import (
	"net"
	"github.com/bhojpur/web/pkg/virtual"
)

primary, _ := net.Listen("tcp", *listen)

// accept a new connection
conn, _ := primary.Accept()

if vhostConn, err = virtual.TLS(conn); err != nil {
	panic("Not a valid TLS connection!")
}

cipherSuites := virtual.ClientHelloMsg.CipherSuites
sessionId := virtual.ClientHelloMsg.SessionId
```

#### Memory Reduction with Free

After you are done multiplexing, you probably do not need to inspect the HTTP header
anymore, so you can make it available for garbage collection:

```go
// look up the upstream host
upstreamHost := hostMapping[vhostConn.Host()]

// free up the multiplex data
vhostConn.Free()

// vhostConn.Host() == ""
// vhostConn.Request == nil (HTTP)
// vhostConn.ClientHelloMsg == nil (TLS)
```

## TCP Stream Multiplexing

Here is an example of stream `client`, which responds to simple `JSON` requests
from a stream `server`.

```go
import (
	"net"
	"github.com/bhojpur/web/pkg/virtual"
)

conn, _ := net.Dial("tcp", "bhojpur.net:1234")
sess := virtual.StreamClient(conn)

for {
    stream, _ := sess.Accept()
    go func(str net.Conn) {
        defer str.Close()
        var req Request
        json.NewDecoder(str).Decode(&req)
        response := handleRequest(&req)
        json.NewEncoder(str).Encode(response)
    }(stream)
}
```

May be the stream `client` wants to make a request to the TCP stream `server`
instead of just responding. This is easy as well:

```go
import (
	"net"
	"github.com/bhojpur/web/pkg/virtual"
)

conn, _ := net.Dial("tcp", "bhojpur.net:1234")
sess := virtual.StreamClient(conn)

stream, _ := sess.Open()
req := Request{
    Query: "What is the meaning of Thekua, the cookies and everything?",
}
json.NewEncoder(stream).Encode(&req)
var resp Response
json.dec.Decode(&resp)
if resp.Answer != "42" {
    panic("wrong answer to the ultimate question!")
}
```
