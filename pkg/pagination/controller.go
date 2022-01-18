package pagination

import (
	"github.com/bhojpur/web/pkg/context"
	"github.com/bhojpur/web/pkg/core/utils/pagination"
)

// SetPaginator Instantiates a Paginator and assigns it to context.Input.Data("paginator").
func SetPaginator(context *context.Context, per int, nums int64) (paginator *pagination.Paginator) {
	paginator = pagination.NewPaginator(context.Request, per, nums)
	context.Input.SetData("paginator", &paginator)
	return
}
