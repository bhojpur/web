# Bhojpur Web - Application Service Framework

The `Bhojpur Web` is an enterprise grade, distributed `applications`/`services` framework, and
a client/server engine used by the [Bhojpur.NET Platform](https://github.com/bhojpur/platform)
for secure applications and/or services delivery using HTTP(S) protocols. It is pre-integrated
with the [Bhojpur IAM](https://github.com/bhojpur/iam) for enable identity and access management.

## Server Engine

It is used as a primary HTTP(S) server engine within the [Bhojpur.NET Platform](https://github.com/bhojpur/platform) ecosystem to host a wide range of web-enabled applications or services. It complies fully with the HTTP1.1, HTTP/2.0 protocol standards. It can function as a `WebAssembly` hosting environment as well.

## Client Engine

Just like the cURL utility, it could be utilized as an HTTP/S client software by application
software testing tools. For example, web performance testing tools benefit from the statistics
framework.

### WebAssembly Application

The `client-side` framework features development of web applications based on `WebAssembly` in Go.
You can compile a custom developed web application using the following commands.

```bash
$ GOARCH=wasm GOOS=js go build -o web/app.wasm
$ go build -o myapp
```

**NOTE** that the build output is explicitly set to `web/app.wasm`. The reason why it is that way:
the HTTP `Handler` associated from the **client-side** engine expects it to be a static resource
located at the `/web/app.wasm` path.

## Reverse Proxy

It is used as a primary Reverse Proxy server within the [Bhojpur.NET Platform](https://github.com/bhojpur/platform) ecosystem to route HTTP traffic securely among applications.

## Application Generators

Using template files, it can automatically generate web application for multiple languages.

```bash
$ webctl generate --pkg testdata views/... --o views.go
```

## Application Performance Testing

The Bhojpur Web is used as a distributed application load testing tool. We benchmark HTTP servers and applications using different features built into the framework.

### Load Testing Usage

```bash
$ webctl perftest [options] URL
```

Application Options:
    --num-requests  Number of requests to make (1)
    --concurrent    Number of concurrent connections to make (1)
    --keep-alive    Use keep alive connection
    --no-gzip       Disable gzip accept encoding
    --secure-tls    Validate TLS/SSL certificates

#### For Example

```bash
$ webctl perftest --num-requests 100 --concurrent 4 https://www.bhojpur.net
```

Then, it would display something like this

```bash
    # Requests: 100
    # Successes: 100
    # Failures: 0
    # Unavailable: 0
    Duration: 1.719238256s
    Average Request Duration: 13.575435ms
```
