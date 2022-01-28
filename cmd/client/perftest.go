package cmd

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
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var perftestCmdOpts struct {
	NumRequests int
	Concurrent  int
	KeepAlive   bool
	Headers     []string
	NoGzip      bool
	SecureTLS   bool
}

// perftestCmd represents the perftest command
var perftestCmd = &cobra.Command{
	Use:   "perftest",
	Short: "Execute performance tests over servers, applications or services",
	Run: func(cmd *cobra.Command, args []string) {
		testing(args)
	},
}

func init() {
	rootCmd.AddCommand(perftestCmd)

	var numRequests int
	webNumRequests := os.Getenv("WEB_NUM_REQUESTS")
	if webNumRequests == "" {
		numRequests = 1
	}
	var concurrent int
	webConcurrent := os.Getenv("WEB_CONCURRENT")
	if webConcurrent == "" {
		concurrent = 1
	}
	var keepAlive bool
	webKeepAlive := os.Getenv("WEB_KEEP_ALIVE")
	if webKeepAlive == "" {
		keepAlive = false
	}
	var headers []string
	webHeaders := os.Getenv("WEB_HEADERS")
	if webHeaders == "" {
		headers = nil
	}
	var noGzip bool
	webNoGzip := os.Getenv("WEB_NO_GZIP")
	if webNoGzip == "" {
		noGzip = false
	}
	var secureTLS bool
	webSecureTLS := os.Getenv("WEB_SECURE_TLS")
	if webSecureTLS == "" {
		secureTLS = false
	}
	perftestCmd.PersistentFlags().IntVar(&perftestCmdOpts.NumRequests, "num-requests", numRequests, "Number of requests to make")
	perftestCmd.PersistentFlags().IntVar(&perftestCmdOpts.Concurrent, "concurrent", concurrent, "Number of concurrent connections to make")
	perftestCmd.PersistentFlags().BoolVar(&perftestCmdOpts.KeepAlive, "keep-alive", keepAlive, "Use keep alive connection")
	perftestCmd.PersistentFlags().StringArrayVar(&perftestCmdOpts.Headers, "header", headers, "Header to include in request (can be used multiple times)")
	perftestCmd.PersistentFlags().BoolVar(&perftestCmdOpts.KeepAlive, "no-gzip", noGzip, "Disable gzip accept encoding")
	perftestCmd.PersistentFlags().BoolVar(&perftestCmdOpts.KeepAlive, "secure-tls", secureTLS, "Validate TLS/SSL certificates")
}

type result struct {
	duration   time.Duration
	statusCode int
	bytesRead  int
	err        error
}

type Summary struct {
	numRequests          int
	totalRequestDuration time.Duration
	avgRequestDuration   time.Duration
	duration             time.Duration
	numSuccesses         int
	numFailures          int
	numUnavailables      int
	requestsPerSecond    float64
	totalBytesRead       int
}

var requestChan chan *http.Request
var resultChan chan *result
var summaryChan chan *Summary
var client *http.Client

func doRequests() {
	for request := range requestChan {
		startTime := time.Now()
		response, err := client.Do(request)
		if err != nil {
			resultChan <- &result{err: err}
			continue

		}
		bytesRead, err := io.Copy(ioutil.Discard, response.Body)
		if err != nil {
			resultChan <- &result{err: err}
			continue
		}

		resultChan <- &result{duration: time.Since(startTime), statusCode: response.StatusCode, bytesRead: int(bytesRead)}
	}
}

func generateRequests(target string, headers []string, numRequests int) {
	request, err := http.NewRequest("GET", target, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create HTTP request - %v\n", err)
		os.Exit(1)
	}

	if !perftestCmdOpts.NoGzip {
		request.Header.Add("Accept-Encoding", "gzip")
	}

	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "Invalid header - %s\n", h)
			os.Exit(1)
		}
		request.Header.Add(parts[0], parts[1])
	}

	for i := 0; i < numRequests; i++ {
		requestChan <- request
	}
	close(requestChan)
}

func summarizeResults(numRequests int, startTime time.Time) {
	summary := new(Summary)

	for i := 0; i < numRequests; i++ {
		result := <-resultChan
		summary.numRequests++
		if result.err != nil {
			summary.numUnavailables++
		} else if result.statusCode >= 400 {
			summary.numFailures++
		} else {
			summary.numSuccesses++
			summary.totalRequestDuration += result.duration
			summary.totalBytesRead += result.bytesRead
		}
	}

	summary.duration = time.Since(startTime)
	if 0 < summary.numSuccesses {
		summary.avgRequestDuration = time.Duration(int64(summary.totalRequestDuration) / int64(summary.numSuccesses))
	}
	summary.requestsPerSecond = float64(summary.numSuccesses) / summary.duration.Seconds()
	summaryChan <- summary
}

func testing(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Requires one target URL")
		return
	}
	target := args[0]
	fmt.Println("generating web load for", target)
	requestChan = make(chan *http.Request)
	resultChan = make(chan *result)
	summaryChan = make(chan *Summary)
	transport := &http.Transport{
		DisableKeepAlives:   !perftestCmdOpts.KeepAlive,
		MaxIdleConnsPerHost: perftestCmdOpts.Concurrent,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: !perftestCmdOpts.SecureTLS},
		DisableCompression:  true,
	}
	client = &http.Client{Transport: transport}

	startTime := time.Now()

	for i := 0; i < perftestCmdOpts.Concurrent; i++ {
		go doRequests()
	}
	go generateRequests(target, perftestCmdOpts.Headers, perftestCmdOpts.NumRequests)
	go summarizeResults(perftestCmdOpts.NumRequests, startTime)

	summary := <-summaryChan

	fmt.Printf("# Requests: %v\n", summary.numRequests)
	fmt.Printf("# Successes: %v\n", summary.numSuccesses)
	fmt.Printf("# Failures: %v\n", summary.numFailures)
	fmt.Printf("# Unavailable: %v\n", summary.numUnavailables)
	fmt.Printf("Duration: %v\n", summary.duration)
	fmt.Printf("Average Request Duration: %v\n", summary.avgRequestDuration)
	fmt.Printf("Requests Per Second: %f\n", summary.requestsPerSecond)
	fmt.Printf("Bytes Received (excluding headers): %d\n", summary.totalBytesRead)
}
