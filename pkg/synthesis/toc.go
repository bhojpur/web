package synthesis

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
	"fmt"
	"io"
	"sort"
	"strings"
)

type assetTree struct {
	Asset    Asset
	Children map[string]*assetTree
}

func newAssetTree() *assetTree {
	tree := &assetTree{}
	tree.Children = make(map[string]*assetTree)
	return tree
}

func (node *assetTree) child(name string) *assetTree {
	rv, ok := node.Children[name]
	if !ok {
		rv = newAssetTree()
		node.Children[name] = rv
	}
	return rv
}

func (root *assetTree) Add(route []string, asset Asset) {
	for _, name := range route {
		root = root.child(name)
	}
	root.Asset = asset
}

func ident(w io.Writer, n int) {
	for i := 0; i < n; i++ {
		w.Write([]byte{'\t'})
	}
}

func (root *assetTree) funcOrNil() string {
	if root.Asset.Func == "" {
		return "nil"
	} else {
		return root.Asset.Func
	}
}

func getFillerSize(tokenIndex int, lengths []int, nident int) int {
	var (
		curlen    int = lengths[tokenIndex]
		maxlen    int = 0
		substart  int = 0
		subend    int = 0
		spacediff int = 0
	)

	if curlen > 0 {
		substart = tokenIndex
		for (substart-1) >= 0 && lengths[substart-1] > 0 {
			substart -= 1
		}

		subend = tokenIndex
		for (subend+1) < len(lengths) && lengths[subend+1] > 0 {
			subend += 1
		}

		var candidate int
		for j := substart; j <= subend; j += 1 {
			candidate = lengths[j]
			if candidate > maxlen {
				maxlen = candidate
			}
		}

		spacediff = maxlen - curlen
	}

	return spacediff
}

func (root *assetTree) writeGoMap(w io.Writer, nident int) {
	fmt.Fprintf(w, "&bhojpurTree{%s, map[string]*bhojpurTree{", root.funcOrNil())

	if len(root.Children) > 0 {
		io.WriteString(w, "\n")

		// Sort to make output stable between invocations
		filenames := make([]string, len(root.Children))
		hasChildren := make(map[string]bool)
		i := 0
		for filename, node := range root.Children {
			filenames[i] = filename
			hasChildren[filename] = len(node.Children) > 0
			i++
		}
		sort.Strings(filenames)

		lengths := make([]int, len(root.Children))
		for i, filename := range filenames {
			if hasChildren[filename] {
				lengths[i] = 0
			} else {
				lengths[i] = len(filename)
			}
		}

		for i, p := range filenames {
			ident(w, nident+1)
			filler := strings.Repeat(" ", getFillerSize(i, lengths, nident))
			fmt.Fprintf(w, `"%s": %s`, p, filler)
			root.Children[p].writeGoMap(w, nident+1)
		}
		ident(w, nident)
	}

	io.WriteString(w, "}}")
	if nident > 0 {
		io.WriteString(w, ",")
	}
	io.WriteString(w, "\n")
}

func (root *assetTree) WriteAsGoMap(w io.Writer) error {
	_, err := fmt.Fprint(w, `type bhojpurTree struct {
	Func     func() (*asset, error)
	Children map[string]*bhojpurTree
}

var _bhojpurTree = `)
	root.writeGoMap(w, 0)
	return err
}

func writeTOCTree(w io.Writer, toc []Asset) error {
	_, err := fmt.Fprintf(w, `// AssetDir returns the file names below a certain
// directory embedded in the file by Bhojpur Webctl.
// For example if you run webctl generate on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then, AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("nonexistent") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bhojpurTree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %%s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %%s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

`)
	if err != nil {
		return err
	}
	tree := newAssetTree()
	for i := range toc {
		pathList := strings.Split(toc[i].Name, "/")
		tree.Add(pathList, toc[i])
	}
	return tree.WriteAsGoMap(w)
}

// writeTOC writes the table of contents file.
func writeTOC(w io.Writer, toc []Asset) error {
	err := writeTOCHeader(w)
	if err != nil {
		return err
	}

	var maxlen = 0
	for i := range toc {
		l := len(toc[i].Name)
		if l > maxlen {
			maxlen = l
		}
	}

	for i := range toc {
		err = writeTOCAsset(w, &toc[i], maxlen)
		if err != nil {
			return err
		}
	}

	return writeTOCFooter(w)
}

// writeTOCHeader writes the table of contents file header.
func writeTOCHeader(w io.Writer) error {
	_, err := fmt.Fprintf(w, `// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bhojpur[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %%s can't read by error: %%v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %%s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bhojpur[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %%s can't read by error: %%v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %%s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bhojpur))
	for name := range _bhojpur {
		names = append(names, name)
	}
	return names
}

// _bhojpur is a table, holding each asset generator, mapped to its name.
var _bhojpur = map[string]func() (*asset, error){
`)
	return err
}

// writeTOCAsset write a TOC entry for the given asset.
func writeTOCAsset(w io.Writer, asset *Asset, maxlen int) error {
	spacediff := maxlen - len(asset.Name)
	filler := strings.Repeat(" ", spacediff)

	_, err := fmt.Fprintf(w, "\t%q: %s%s,\n", asset.Name, filler, asset.Func)
	return err
}

// writeTOCFooter writes the table of contents file footer.
func writeTOCFooter(w io.Writer) error {
	_, err := fmt.Fprintf(w, `}

`)
	return err
}
