package context

import (
	"github.com/bhojpur/web/pkg/context"
)

// BhojpurOutput does work for sending response header.
type BhojpurOutput context.BhojpurOutput

// NewOutput returns new BhojpurOutput.
// it contains nothing now.
func NewOutput() *BhojpurOutput {
	return (*BhojpurOutput)(context.NewOutput())
}

// Reset init BhojpurOutput
func (output *BhojpurOutput) Reset(ctx *Context) {
	(*context.BhojpurOutput)(output).Reset((*context.Context)(ctx))
}

// Header sets response header item string via given key.
func (output *BhojpurOutput) Header(key, val string) {
	(*context.BhojpurOutput)(output).Header(key, val)
}

// Body sets response body content.
// if EnableGzip, compress content string.
// it sends out response body directly.
func (output *BhojpurOutput) Body(content []byte) error {
	return (*context.BhojpurOutput)(output).Body(content)
}

// Cookie sets cookie value via given key.
// others are ordered as cookie's max age time, path,domain, secure and httponly.
func (output *BhojpurOutput) Cookie(name string, value string, others ...interface{}) {
	(*context.BhojpurOutput)(output).Cookie(name, value, others)
}

// JSON writes json to response body.
// if encoding is true, it converts utf-8 to \u0000 type.
func (output *BhojpurOutput) JSON(data interface{}, hasIndent bool, encoding bool) error {
	return (*context.BhojpurOutput)(output).JSON(data, hasIndent, encoding)
}

// YAML writes yaml to response body.
func (output *BhojpurOutput) YAML(data interface{}) error {
	return (*context.BhojpurOutput)(output).YAML(data)
}

// JSONP writes jsonp to response body.
func (output *BhojpurOutput) JSONP(data interface{}, hasIndent bool) error {
	return (*context.BhojpurOutput)(output).JSONP(data, hasIndent)
}

// XML writes xml string to response body.
func (output *BhojpurOutput) XML(data interface{}, hasIndent bool) error {
	return (*context.BhojpurOutput)(output).XML(data, hasIndent)
}

// ServeFormatted serve YAML, XML OR JSON, depending on the value of the Accept header
func (output *BhojpurOutput) ServeFormatted(data interface{}, hasIndent bool, hasEncode ...bool) {
	(*context.BhojpurOutput)(output).ServeFormatted(data, hasIndent, hasEncode...)
}

// Download forces response for download file.
// it prepares the download response header automatically.
func (output *BhojpurOutput) Download(file string, filename ...string) {
	(*context.BhojpurOutput)(output).Download(file, filename...)
}

// ContentType sets the content type from ext string.
// MIME type is given in mime package.
func (output *BhojpurOutput) ContentType(ext string) {
	(*context.BhojpurOutput)(output).ContentType(ext)
}

// SetStatus sets response status code.
// It writes response header directly.
func (output *BhojpurOutput) SetStatus(status int) {
	(*context.BhojpurOutput)(output).SetStatus(status)
}

// IsCachable returns boolean of this request is cached.
// HTTP 304 means cached.
func (output *BhojpurOutput) IsCachable() bool {
	return (*context.BhojpurOutput)(output).IsCachable()
}

// IsEmpty returns boolean of this request is empty.
// HTTP 201，204 and 304 means empty.
func (output *BhojpurOutput) IsEmpty() bool {
	return (*context.BhojpurOutput)(output).IsEmpty()
}

// IsOk returns boolean of this request runs well.
// HTTP 200 means ok.
func (output *BhojpurOutput) IsOk() bool {
	return (*context.BhojpurOutput)(output).IsOk()
}

// IsSuccessful returns boolean of this request runs successfully.
// HTTP 2xx means ok.
func (output *BhojpurOutput) IsSuccessful() bool {
	return (*context.BhojpurOutput)(output).IsSuccessful()
}

// IsRedirect returns boolean of this request is redirection header.
// HTTP 301,302,307 means redirection.
func (output *BhojpurOutput) IsRedirect() bool {
	return (*context.BhojpurOutput)(output).IsRedirect()
}

// IsForbidden returns boolean of this request is forbidden.
// HTTP 403 means forbidden.
func (output *BhojpurOutput) IsForbidden() bool {
	return (*context.BhojpurOutput)(output).IsForbidden()
}

// IsNotFound returns boolean of this request is not found.
// HTTP 404 means not found.
func (output *BhojpurOutput) IsNotFound() bool {
	return (*context.BhojpurOutput)(output).IsNotFound()
}

// IsClientError returns boolean of this request client sends error data.
// HTTP 4xx means client error.
func (output *BhojpurOutput) IsClientError() bool {
	return (*context.BhojpurOutput)(output).IsClientError()
}

// IsServerError returns boolean of this server handler errors.
// HTTP 5xx means server internal error.
func (output *BhojpurOutput) IsServerError() bool {
	return (*context.BhojpurOutput)(output).IsServerError()
}

// Session sets session item value with given key.
func (output *BhojpurOutput) Session(name interface{}, value interface{}) {
	(*context.BhojpurOutput)(output).Session(name, value)
}
