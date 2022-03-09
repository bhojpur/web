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
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bhojpur/web/pkg/client/codegen/command"
	"github.com/bhojpur/web/pkg/client/codegen/render"
	"github.com/bhojpur/web/pkg/client/codegen/system"
	cliLogger "github.com/bhojpur/web/pkg/client/logger"
	"github.com/bhojpur/web/pkg/client/utils"
	"github.com/bhojpur/web/pkg/template"
)

// store all data
type Container struct {
	BhojpurWebFile   string                 // bhojpur toml
	TimestampFile    string                 // store ts file
	GoModFile        string                 // go mod file
	UserOption       UserOption             // user option
	TmplOption       TmplOption             // tmpl option
	CurPath          string                 // user current path
	EnableModules    map[string]interface{} // bhojpur web provider a collection of module
	FunctionOnce     map[string]sync.Once   // exec function once
	Timestamp        Timestamp
	GenerateTime     string
	GenerateTimeUnix int64
	Parser           Parser
}

// user option
type UserOption struct {
	Debug          bool                 `json:"debug"`
	ContextDebug   bool                 `json:"contextDebug"`
	Dsn            string               `json:"dsn"`
	Driver         string               `json:"driver"`
	ProType        string               `json:"proType"`
	ApiPrefix      string               `json:"apiPrefix"`
	EnableModule   []string             `json:"enableModule"`
	Models         map[string]TextModel `json:"models"`
	GitRemotePath  string               `json:"gitRemotePath"`
	Branch         string               `json:"branch"`
	GitLocalPath   string               `json:"gitLocalPath"`
	EnableFormat   bool                 `json:"enableFormat"`
	SourceGen      string               `json:"sourceGen"`
	EnableGitPull  bool                 `json:"enbaleGitPull"`
	Path           map[string]string    `json:"path"`
	EnableGomod    bool                 `json:"enableGomod"`
	RefreshGitTime int64                `json:"refreshGitTime"`
	Extend         map[string]string    `json:"extend"` // extend user data
}

// tmpl option
type TmplOption struct {
	RenderPath string `toml:"renderPath"`
	Descriptor []Descriptor
}

type Descriptor struct {
	Module  string `toml:"module"`
	SrcName string `toml:"srcName"`
	DstPath string `toml:"dstPath"`
	Once    bool   `toml:"once"`
	Script  string `toml:"script"`
}

func (descriptor Descriptor) Parse(modelName string, paths map[string]string) (newDescriptor Descriptor, ctx template.Context) {
	var (
		err             error
		relativeDstPath string
		absFile         string
		relPath         string
	)

	newDescriptor = descriptor
	render := render.NewRender("")
	ctx = make(template.Context)
	for key, value := range paths {
		absFile, err = filepath.Abs(value)
		if err != nil {
			cliLogger.Log.Fatalf("absolute path error %s from key %s and value %s", err, key, value)
		}
		relPath, err = filepath.Rel(system.CurrentDir, absFile)
		if err != nil {
			cliLogger.Log.Fatalf("Could not get the relative path: %s", err)
		}
		// user input path
		ctx["path"+utils.CamelCase(key)] = value
		// relativePath
		ctx["pathRel"+utils.CamelCase(key)] = relPath
	}
	ctx["modelName"] = lowerFirst(utils.CamelString(modelName))
	relativeDstPath, err = render.TemplateFromString(descriptor.DstPath).Execute(ctx)
	if err != nil {
		cliLogger.Log.Fatalf("bhojpur tmpl exec error, err: %s", err)
		return
	}

	newDescriptor.DstPath, err = filepath.Abs(relativeDstPath)
	if err != nil {
		cliLogger.Log.Fatalf("absolute path error %s from flush file %s", err, relativeDstPath)
	}

	newDescriptor.Script, err = render.TemplateFromString(descriptor.Script).Execute(ctx)
	if err != nil {
		cliLogger.Log.Fatalf("parse script %s, error %s", descriptor.Script, err)
	}
	return
}

func (descriptor Descriptor) IsExistScript() bool {
	return descriptor.Script != ""
}

func (d Descriptor) ExecScript(path string) (err error) {
	arr := strings.Split(d.Script, " ")
	if len(arr) == 0 {
		return
	}

	stdout, stderr, err := command.ExecCmdDir(path, arr[0], arr[1:]...)
	if err != nil {
		return concatenateError(err, stderr)
	}

	cliLogger.Log.Info(stdout)
	return nil
}

type Timestamp struct {
	GitCacheLastRefresh int64 `toml:"gitCacheLastRefresh"`
	Generate            int64 `toml:"generate"`
}

func concatenateError(err error, stderr string) error {
	if len(stderr) == 0 {
		return err
	}
	return fmt.Errorf("%v: %s", err, stderr)
}
