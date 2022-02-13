package multiplexer

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
	"io"
	"io/ioutil"
	"testing"
)

func BenchmarkPing(b *testing.B) {
	client, server := testClientServer()
	defer func() {
		client.Close()
		server.Close()
	}()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rtt, err := client.Ping()
		if err != nil {
			b.Fatalf("err: %v", err)
		}
		if rtt == 0 {
			b.Fatalf("bad: %v", rtt)
		}
	}
}

func BenchmarkAccept(b *testing.B) {
	client, server := testClientServer()
	defer func() {
		client.Close()
		server.Close()
	}()

	doneCh := make(chan struct{})
	b.ReportAllocs()
	b.ResetTimer()

	go func() {
		defer close(doneCh)

		for i := 0; i < b.N; i++ {
			stream, err := server.AcceptStream()
			if err != nil {
				return
			}
			stream.Close()
		}
	}()

	for i := 0; i < b.N; i++ {
		stream, err := client.Open()
		if err != nil {
			b.Fatalf("err: %v", err)
		}
		stream.Close()
	}
	<-doneCh
}

func BenchmarkSendRecv32(b *testing.B) {
	const payloadSize = 32
	benchmarkSendRecv(b, payloadSize, payloadSize)
}

func BenchmarkSendRecv64(b *testing.B) {
	const payloadSize = 64
	benchmarkSendRecv(b, payloadSize, payloadSize)
}

func BenchmarkSendRecv128(b *testing.B) {
	const payloadSize = 128
	benchmarkSendRecv(b, payloadSize, payloadSize)
}

func BenchmarkSendRecv256(b *testing.B) {
	const payloadSize = 256
	benchmarkSendRecv(b, payloadSize, payloadSize)
}

func BenchmarkSendRecv512(b *testing.B) {
	const payloadSize = 512
	benchmarkSendRecv(b, payloadSize, payloadSize)
}

func BenchmarkSendRecv1024(b *testing.B) {
	const payloadSize = 1024
	benchmarkSendRecv(b, payloadSize, payloadSize)
}

func BenchmarkSendRecv2048(b *testing.B) {
	const payloadSize = 2048
	benchmarkSendRecv(b, payloadSize, payloadSize)
}

func BenchmarkSendRecv4096(b *testing.B) {
	const payloadSize = 4096
	benchmarkSendRecv(b, payloadSize, payloadSize)
}

func BenchmarkSendRecvLarge(b *testing.B) {
	const sendSize = 512 * 1024 * 1024 //512 MB
	const recvSize = 4 * 1024          //4 KB
	benchmarkSendRecv(b, sendSize, recvSize)
}

func benchmarkSendRecv(b *testing.B, sendSize, recvSize int) {
	client, server := testClientServer()
	defer func() {
		client.Close()
		server.Close()
	}()

	sendBuf := make([]byte, sendSize)
	recvBuf := make([]byte, recvSize)
	doneCh := make(chan struct{})

	b.SetBytes(int64(sendSize))
	b.ReportAllocs()
	b.ResetTimer()

	go func() {
		defer close(doneCh)

		stream, err := server.AcceptStream()
		if err != nil {
			return
		}
		defer stream.Close()

		switch {
		case sendSize == recvSize:
			for i := 0; i < b.N; i++ {
				if _, err := stream.Read(recvBuf); err != nil {
					b.Fatalf("err: %v", err)
				}
			}

		case recvSize > sendSize:
			b.Fatalf("bad test case; recvSize was: %d and sendSize was: %d, but recvSize must be <= sendSize!", recvSize, sendSize)

		default:
			chunks := sendSize / recvSize
			for i := 0; i < b.N; i++ {
				for j := 0; j < chunks; j++ {
					if _, err := stream.Read(recvBuf); err != nil {
						b.Fatalf("err: %v", err)
					}
				}
			}
		}
	}()

	stream, err := client.Open()
	if err != nil {
		b.Fatalf("err: %v", err)
	}
	defer stream.Close()

	for i := 0; i < b.N; i++ {
		if _, err := stream.Write(sendBuf); err != nil {
			b.Fatalf("err: %v", err)
		}
	}
	<-doneCh
}

func BenchmarkSendRecvParallel32(b *testing.B) {
	const payloadSize = 32
	benchmarkSendRecvParallel(b, payloadSize)
}

func BenchmarkSendRecvParallel64(b *testing.B) {
	const payloadSize = 64
	benchmarkSendRecvParallel(b, payloadSize)
}

func BenchmarkSendRecvParallel128(b *testing.B) {
	const payloadSize = 128
	benchmarkSendRecvParallel(b, payloadSize)
}

func BenchmarkSendRecvParallel256(b *testing.B) {
	const payloadSize = 256
	benchmarkSendRecvParallel(b, payloadSize)
}

func BenchmarkSendRecvParallel512(b *testing.B) {
	const payloadSize = 512
	benchmarkSendRecvParallel(b, payloadSize)
}

func BenchmarkSendRecvParallel1024(b *testing.B) {
	const payloadSize = 1024
	benchmarkSendRecvParallel(b, payloadSize)
}

func BenchmarkSendRecvParallel2048(b *testing.B) {
	const payloadSize = 2048
	benchmarkSendRecvParallel(b, payloadSize)
}

func BenchmarkSendRecvParallel4096(b *testing.B) {
	const payloadSize = 4096
	benchmarkSendRecvParallel(b, payloadSize)
}

func benchmarkSendRecvParallel(b *testing.B, sendSize int) {
	client, server := testClientServer()
	defer func() {
		client.Close()
		server.Close()
	}()

	sendBuf := make([]byte, sendSize)
	discarder := ioutil.Discard.(io.ReaderFrom)
	b.SetBytes(int64(sendSize))
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		doneCh := make(chan struct{})

		go func() {
			defer close(doneCh)

			stream, err := server.AcceptStream()
			if err != nil {
				return
			}
			defer stream.Close()

			if _, err := discarder.ReadFrom(stream); err != nil {
				b.Fatalf("err: %v", err)
			}
		}()

		stream, err := client.Open()
		if err != nil {
			b.Fatalf("err: %v", err)
		}

		for pb.Next() {
			if _, err := stream.Write(sendBuf); err != nil {
				b.Fatalf("err: %v", err)
			}
		}

		stream.Close()
		<-doneCh
	})
}
