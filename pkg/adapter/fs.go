package adapter

import (
	"net/http"
	"path/filepath"

	web "github.com/bhojpur/web/pkg/engine"
)

type FileSystem web.FileSystem

func (d FileSystem) Open(name string) (http.File, error) {
	return (web.FileSystem)(d).Open(name)
}

// Walk walks the file tree rooted at root in filesystem, calling walkFn for each file or
// directory in the tree, including root. All errors that arise visiting files
// and directories are filtered by walkFn.
func Walk(fs http.FileSystem, root string, walkFn filepath.WalkFunc) error {
	return web.Walk(fs, root, walkFn)
}
