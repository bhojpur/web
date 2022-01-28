package adapter

import (
	"github.com/bhojpur/web/pkg/adapter/context"
	beecontext "github.com/bhojpur/web/pkg/context"

	web "github.com/bhojpur/web/pkg/engine"
)

// Tree has three elements: FixRouter/wildcard/leaves
// fixRouter stores Fixed Router
// wildcard stores params
// leaves store the endpoint information
type Tree web.Tree

// NewTree return a new Tree
func NewTree() *Tree {
	return (*Tree)(web.NewTree())
}

// AddTree will add tree to the exist Tree
// prefix should has no params
func (t *Tree) AddTree(prefix string, tree *Tree) {
	(*web.Tree)(t).AddTree(prefix, (*web.Tree)(tree))
}

// AddRouter call addseg function
func (t *Tree) AddRouter(pattern string, runObject interface{}) {
	(*web.Tree)(t).AddRouter(pattern, runObject)
}

// Match router to runObject & params
func (t *Tree) Match(pattern string, ctx *context.Context) (runObject interface{}) {
	return (*web.Tree)(t).Match(pattern, (*beecontext.Context)(ctx))
}