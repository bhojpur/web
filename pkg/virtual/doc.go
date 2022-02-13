// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package virtual

// It implements a general purpose virtual hosting stream-multiplexing protocol.
// It allows clients applications to multiplex any io.ReadWriteCloser (like a
// net.Conn) into multiple, independent full-duplex streams.
//
// It is a useful protocol for any two communicating processes. It is an excellent
// base protocol for implementing lightweight RPC. It eliminates the need for
// custom async/pipeling code from your peers in order to support multiple
// simultaneous inflight requests between peers. For the same reason, it also
// eliminates the need to build connection pools for your clients. It enables
// servers to initiate streams to clients without building any NAT traversal.
// It can also yield performance improvements (especially latency) for protocols
// that require rapidly opening many concurrent connections.
//
// Here's an example client which responds to simple JSON requests from a server.
//
//     conn, _ := net.Dial("tcp", "bhojpur.net:1234")
//     sess := virtual.StreamClient(conn)
//     for {
//         stream, _ := sess.Accept()
//         go func(str net.Conn) {
//             defer str.Close()
//             var req Request
//             json.NewDecoder(str).Decode(&req)
//             response := handleRequest(&req)
//             json.NewEncoder(str).Encode(response)
//         }(stream)
//     }
//
// May be the client wants to make a request to the server instead of just
// responding. This is easy as well:
//
//     stream, _ := sess.Open()
//     req := Request{
//         Query: "What is the meaning of Thekua, the cookies and everything?",
//     }
//     json.NewEncoder(stream).Encode(&req)
//     var resp Response
//     json.dec.Decode(&resp)
//     if resp.Answer != "42" {
//         panic("wrong answer to the ultimate question!")
//     }
//
// It defines the following terms for further clarity:
//
// A "Transport" is an underlying stream (typically TCP) that is multiplexed by
// sending frames between virtual peers over this transport.
//
// A "Stream" is any of the full-duplex byte-streams multiplexed over the transport
//
// A "Session" is two peers running the multiplexing protocol over a single transport
//
// Our design is influenced heavily by the framing layer of HTTP2 and SPDY. However,
// instead of being specialized in a higher-level protocol, we designed it in a
// protocol agnostic way with simplicity and speed in mind. More advanced features
// are left to higher-level libraries and protocols.
//
// Virtual Hosting API is designed to make it seamless to integrate into existing
// programs. virtual.Session implements the net.Listener interface and virtual.Stream
// implements net.Conn.
//
// It ships with two wrappers that add commonly used functionality. The first is a
// TypedStreamSession which allows a client application to open streams with a type
// identifier so that the remote peer can identify the protocol that will be
// communicated on that stream.
//
// The second wrapper is a simple Heartbeat, which issues a callback to the
// application informing it of round-trip latency and heartbeat failure.
