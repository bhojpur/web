package prometheus

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/web/pkg/context"
)

func TestFilterChain(t *testing.T) {
	filter := (&FilterChainBuilder{}).FilterChain(func(ctx *context.Context) {
		// do nothing
		ctx.Input.SetData("invocation", true)
	})

	ctx := context.NewContext()
	r, _ := http.NewRequest("GET", "/prometheus/user", nil)
	w := httptest.NewRecorder()
	ctx.Reset(w, r)
	ctx.Input.SetData("RouterPattern", "my-route")
	filter(ctx)
	assert.True(t, ctx.Input.GetData("invocation").(bool))
}
