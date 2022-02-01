# Bhojpur Web - Client HTTPlib

The HTTPlib is a library that helps you to `curl` remote url.

## How to use?

### GET

You can use `GET` to crawl the data.

	import "github.com/bhojpur/web/pkg/client/httplib"

	str, err := httplib.Get("http://bhojpur.net/").String()
	if err != nil {
        	// error
	}
	fmt.Println(str)

### POST

The `POST` data to remote url

	req := httplib.Post("http://bhojpur.net/")
	req.Param("username","bhojpur")
	req.Param("password","123456")
	str, err := req.String()
	if err != nil {
        	// error
	}
	fmt.Println(str)

### Set timeout

The default timeout is `60` seconds, function prototype:

	SetTimeout(connectTimeout, readWriteTimeout time.Duration)

Example:

	// GET
	httplib.Get("http://bhojpur.net/").SetTimeout(100 * time.Second, 30 * time.Second)

	// POST
	httplib.Post("http://bhojpur.net/").SetTimeout(100 * time.Second, 30 * time.Second)

### Debug

If you want to debug the request info, set the debug on

	httplib.Get("http://bhojpur.net/").Debug(true)

### Set HTTP Basic Auth

	str, err := Get("http://bhojpur.net/").SetBasicAuth("user", "passwd").String()
	if err != nil {
        	// error
	}
	fmt.Println(str)

### Set HTTPS

If requested url is using `https`, then you can set the client support TLS:

	httplib.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

More info about the `tls.Config` please visit http://golang.org/pkg/crypto/tls/#Config

### Set HTTP Version

Some servers need to specify the protocol version of HTTP

	httplib.Get("http://bhojpur.net/").SetProtocolVersion("HTTP/1.1")

### Set Cookie

Some HTTP requests need setcookie. So, set it like this:

	cookie := &http.Cookie{}
	cookie.Name = "username"
	cookie.Value  = "bhojpur"
	httplib.Get("http://bhojpur.net/").SetCookie(cookie)

### Upload file

The HTTPlib supports mutil-file upload, use `req.PostFile()`

	req := httplib.Post("http://bhojpur.net/")
	req.Param("username","bhojpur")
	req.PostFile("uploadfile1", "httplib.pdf")
	str, err := req.String()
	if err != nil {
        	// error
	}
	fmt.Println(str)
