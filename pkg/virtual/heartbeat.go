package virtual

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

import (
	"encoding/binary"
	"io"
	"math/rand"
	"net"
	"time"
)

const (
	defaultHeartbeatInterval             = 10 * time.Second
	defaultHeartbeatTolerance            = 15 * time.Second
	defaultStreamType         StreamType = 0xFFFFFFFF
)

type HeartbeatConfig struct {
	Interval  time.Duration
	Tolerance time.Duration
	Type      StreamType
}

func NewHeartbeatConfig() *HeartbeatConfig {
	return &HeartbeatConfig{
		Interval:  defaultHeartbeatInterval,
		Tolerance: defaultHeartbeatTolerance,
		Type:      defaultStreamType,
	}
}

type Heartbeat struct {
	TypedStreamSession
	config HeartbeatConfig
	closed chan int
	cb     func(time.Duration)
}

func NewHeartbeat(sess TypedStreamSession, cb func(time.Duration), config *HeartbeatConfig) *Heartbeat {
	if config == nil {
		config = NewHeartbeatConfig()
	}
	return &Heartbeat{
		TypedStreamSession: sess,
		config:             *config,
		closed:             make(chan int, 1),
		cb:                 cb,
	}
}

func (h *Heartbeat) Accept() (net.Conn, error) {
	return h.AcceptTypedStream()
}

func (h *Heartbeat) AcceptStream() (Stream, error) {
	return h.TypedStreamSession.AcceptTypedStream()
}

func (h *Heartbeat) Close() error {
	select {
	case h.closed <- 1:
	default:
	}
	return h.TypedStreamSession.Close()
}

func (h *Heartbeat) AcceptTypedStream() (TypedStream, error) {
	for {
		str, err := h.TypedStreamSession.AcceptTypedStream()
		if err != nil {
			return nil, err
		}
		if str.StreamType() != h.config.Type {
			return str, nil
		}
		go h.responder(str)
	}
}

func (h *Heartbeat) Start() {
	mark := make(chan time.Duration)
	go h.requester(mark)
	go h.check(mark)
}

func (h *Heartbeat) check(mark chan time.Duration) {
	t := time.NewTimer(h.config.Interval + h.config.Tolerance)
	for {
		select {
		case <-t.C:
			// timed out waiting for a response!
			h.cb(0)

		case dur := <-mark:
			h.cb(dur)
			t.Reset(h.config.Interval + h.config.Tolerance)

		case <-h.closed:
			return
		}
	}
}

func (h *Heartbeat) requester(mark chan time.Duration) {
	// make random number generator
	r := rand.New(rand.NewSource(time.Now().Unix()))

	// open a new stream for the heartbeat
	stream, err := h.OpenTypedStream(h.config.Type)
	if err != nil {
		return
	}

	// send heartbeats and then check that we got them back
	for {
		time.Sleep(h.config.Interval)
		start := time.Now()
		// assign a new random value to echo
		id := uint32(r.Int31())
		if err := binary.Write(stream, binary.BigEndian, id); err != nil {
			return
		}
		var respId uint32
		if err := binary.Read(stream, binary.BigEndian, &respId); err != nil {
			return
		}
		if id != respId {
			return
		}
		// record the time
		mark <- time.Since(start)
	}
}

func (h *Heartbeat) responder(s Stream) {
	// read the next heartbeat id and respond
	buf := make([]byte, 4)
	for {
		_, err := io.ReadFull(s, buf)
		if err != nil {
			return
		}
		_, err = s.Write(buf)
		if err != nil {
			return
		}
	}
}
