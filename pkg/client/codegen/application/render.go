package application

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
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/bhojpur/web/pkg/client/codegen/render"
	"github.com/bhojpur/web/pkg/template"
	"github.com/davecgh/go-spew/spew"

	"github.com/bhojpur/web/pkg/client/codegen/system"
	cliLogger "github.com/bhojpur/web/pkg/client/logger"
)

// render
type RenderFile struct {
	*render.Render
	Context      template.Context
	GenerateTime string
	Option       UserOption
	ModelName    string
	PackageName  string
	FlushFile    string
	PkgPath      string
	TmplPath     string
	Descriptor   Descriptor
}

func NewRender(m RenderInfo) *RenderFile {
	var (
		pathCtx       template.Context
		newDescriptor Descriptor
	)

	// parse descriptor, get flush file path, bhojpur path, etc...
	newDescriptor, pathCtx = m.Descriptor.Parse(m.ModelName, m.Option.Path)

	obj := &RenderFile{
		Context:      make(template.Context),
		Option:       m.Option,
		ModelName:    m.ModelName,
		GenerateTime: m.GenerateTime,
		Descriptor:   newDescriptor,
	}

	obj.FlushFile = newDescriptor.DstPath

	// new render
	obj.Render = render.NewRender(path.Join(obj.Option.GitLocalPath, obj.Option.ProType, m.TmplPath))

	filePath := path.Dir(obj.FlushFile)
	err := createPath(filePath)
	if err != nil {
		cliLogger.Log.Fatalf("Could not create the controllers directory: %s", err)
	}
	// get go package path
	obj.PkgPath = getPackagePath()

	relativePath, err := filepath.Rel(system.CurrentDir, obj.FlushFile)
	if err != nil {
		cliLogger.Log.Fatalf("Could not get the relative path: %s", err)
	}

	modelSchemas := m.Content.ToModelSchemas()
	camelPrimaryKey := modelSchemas.GetPrimaryKey()
	importMaps := make(map[string]struct{})
	if modelSchemas.IsExistTime() {
		importMaps["time"] = struct{}{}
	}
	obj.PackageName = filepath.Base(filepath.Dir(relativePath))
	cliLogger.Log.Infof("Using '%s' as name", obj.ModelName)

	cliLogger.Log.Infof("Using '%s' as package name from %s", obj.ModelName, obj.PackageName)

	// package
	obj.SetContext("packageName", obj.PackageName)
	obj.SetContext("packageImports", importMaps)

	// todo optimize
	// todo Set the bhojpur directory, should recalculate the package
	if pathCtx["pathRelBhojpur"] == "." {
		obj.SetContext("packagePath", obj.PkgPath)
	} else {
		obj.SetContext("packagePath", obj.PkgPath+"/"+pathCtx["pathRelBhojpur"].(string))
	}

	obj.SetContext("packageMod", obj.PkgPath)

	obj.SetContext("modelSchemas", modelSchemas)
	obj.SetContext("modelPrimaryKey", camelPrimaryKey)

	for key, value := range pathCtx {
		obj.SetContext(key, value)
	}

	obj.SetContext("apiPrefix", obj.Option.ApiPrefix)
	obj.SetContext("generateTime", obj.GenerateTime)

	if obj.Option.ContextDebug {
		spew.Dump(obj.Context)
	}
	return obj
}

func (r *RenderFile) SetContext(key string, value interface{}) {
	r.Context[key] = value
}

func (r *RenderFile) Exec(name string) {
	var (
		buf string
		err error
	)
	buf, err = r.Render.Template(name).Execute(r.Context)
	if err != nil {
		cliLogger.Log.Fatalf("Could not create the %s render tmpl: %s", name, err)
		return
	}
	_, err = os.Stat(r.Descriptor.DstPath)
	var orgContent []byte
	if err == nil {
		if org, err := os.OpenFile(r.Descriptor.DstPath, os.O_RDONLY, 0666); err == nil {
			orgContent, _ = ioutil.ReadAll(org)
			org.Close()
		} else {
			cliLogger.Log.Infof("file err %s", err)
		}
	}
	// Replace or create when content changes
	output := []byte(buf)
	ext := filepath.Ext(r.FlushFile)
	if r.Option.EnableFormat && ext == ".go" {
		// format code
		var bts []byte
		bts, err = format.Source([]byte(buf))
		if err != nil {
			cliLogger.Log.Warnf("format buf error %s", err.Error())
		}
		output = bts
	}

	if FileContentChange(orgContent, output, GetSeg(ext)) {
		err = r.write(r.FlushFile, output)
		if err != nil {
			cliLogger.Log.Fatalf("Could not create file: %s", err)
			return
		}
		cliLogger.Log.Infof("create file '%s' from %s", r.FlushFile, r.PackageName)
	}
}
