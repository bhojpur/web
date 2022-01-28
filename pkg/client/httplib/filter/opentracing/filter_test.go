package opentracing

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/web/pkg/client/httplib"
)

func TestFilterChainBuilder_FilterChain(t *testing.T) {
	next := func(ctx context.Context, req *httplib.BhojpurHTTPRequest) (*http.Response, error) {
		time.Sleep(100 * time.Millisecond)
		return &http.Response{
			StatusCode: 404,
		}, errors.New("hello")
	}
	builder := &FilterChainBuilder{}
	filter := builder.FilterChain(next)
	req := httplib.Get("https://github.com/notifications?query=repo%3Abhojpur%2Fweb")
	resp, err := filter(context.Background(), req)
	assert.NotNil(t, resp)
	assert.NotNil(t, err)
}
