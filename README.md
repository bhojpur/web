# Bhojpur Web - Service Engine

The Bhojpur Web is a service engine used by the Bhojpur.NET Platform for secure application delivery.

## Web Application Load Testing

It is a software tool for web load testing and benchmarking HTTP servers and applications.

### Load Testing Usage

    web load [options] URL

Application Options:

    --num-requests= Number of requests to make (1)
    --concurrent=   Number of concurrent connections to make (1)
    --keep-alive    Use keep alive connection
    --no-gzip       Disable gzip accept encoding
    --secure-tls    Validate TLS/SSL certificates

For Example:

    $ web load --num-requests 100 --concurrent 4 https://www.bhojpur.net
    # Requests: 100
    # Successes: 100
    # Failures: 0
    # Unavailable: 0
    Duration: 1.719238256s
    Average Request Duration: 13.575435ms
