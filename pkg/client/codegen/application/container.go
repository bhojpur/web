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
	"io/ioutil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/bhojpur/web/pkg/client/codegen/git"
	"github.com/bhojpur/web/pkg/client/codegen/system"
	cliLogger "github.com/bhojpur/web/pkg/client/logger"
	"github.com/bhojpur/web/pkg/client/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
)

var GitRemotePath utils.DocValue

const MDateFormat = "20060102_150405"

var DefaultBhojpurWeb = &Container{
	BhojpurWebFile: system.CurrentDir + "/bhojpur.toml",
	TimestampFile:  system.CurrentDir + "/.bhojpur.timestamp",
	GoModFile:      system.CurrentDir + "/go.mod",
	UserOption: UserOption{
		Debug:         false,
		ContextDebug:  false,
		Dsn:           "",
		Driver:        "mysql",
		ProType:       "default",
		ApiPrefix:     "/api",
		EnableModule:  nil,
		Models:        make(map[string]TextModel),
		GitRemotePath: "",
		Branch:        "master",
		EnableFormat:  true,
		SourceGen:     "text",
		EnableGitPull: true,
		Path: map[string]string{
			"bhojpur": ".",
		},
		EnableGomod:    true,
		RefreshGitTime: 24 * 3600,
		Extend:         nil,
	},
	GenerateTime:     time.Now().Format(MDateFormat),
	GenerateTimeUnix: time.Now().Unix(),
	TmplOption:       TmplOption{},
	CurPath:          system.CurrentDir,
	EnableModules:    make(map[string]interface{}), // get the user configuration, get the enable module result
	FunctionOnce:     make(map[string]sync.Once),   // get the tmpl configuration, get the function once result
}

func (c *Container) Run() {
	// init git refresh cache time
	c.initTimestamp()
	c.initUserOption()
	c.initTemplateOption()
	c.initParser()
	c.initRender()
	c.flushTimestamp()
}

func (c *Container) initUserOption() {
	if !utils.IsExist(c.BhojpurWebFile) {
		cliLogger.Log.Fatalf("Bhojpur Web config does not exist, bhojpur json path: %s", c.BhojpurWebFile)
		return
	}
	viper.SetConfigFile(c.BhojpurWebFile)
	err := viper.ReadInConfig()
	if err != nil {
		cliLogger.Log.Fatalf("read bhojpur config content, err: %s", err.Error())
		return
	}
	err = viper.Unmarshal(&c.UserOption)
	if err != nil {
		cliLogger.Log.Fatalf("Bhojpur Web config unmarshal error, err: %s", err.Error())
		return
	}
	if c.UserOption.Debug {
		viper.Debug()
	}

	if c.UserOption.EnableGomod {
		if !utils.IsExist(c.GoModFile) {
			cliLogger.Log.Fatalf("go mod does not exist, please create go mod file")
			return
		}
	}

	for _, value := range c.UserOption.EnableModule {
		c.EnableModules[value] = struct{}{}
	}

	if len(c.EnableModules) == 0 {
		c.EnableModules["*"] = struct{}{}
	}

	if c.UserOption.Debug {
		fmt.Println("c.modules", c.EnableModules)
	}
}

func (c *Container) initTemplateOption() {
	c.GetLocalPath()
	if c.UserOption.EnableGitPull && (c.GenerateTimeUnix-c.Timestamp.GitCacheLastRefresh > c.UserOption.RefreshGitTime) {
		err := git.CloneORPullRepo(c.UserOption.GitRemotePath, c.UserOption.GitLocalPath)
		if err != nil {
			cliLogger.Log.Fatalf("Bhojpur Web git clone or pull repo error, err: %s", err)
			return
		}
		c.Timestamp.GitCacheLastRefresh = c.GenerateTimeUnix
	}

	tree, err := toml.LoadFile(c.UserOption.GitLocalPath + "/" + c.UserOption.ProType + "/bhojpur.toml")

	if err != nil {
		cliLogger.Log.Fatalf("bhojpur tmpl exec error, err: %s", err)
		return
	}
	err = tree.Unmarshal(&c.TmplOption)
	if err != nil {
		cliLogger.Log.Fatalf("bhojpur tmpl parse error, err: %s", err)
		return
	}

	if c.UserOption.Debug {
		spew.Dump("tmpl", c.TmplOption)
	}

	for _, value := range c.TmplOption.Descriptor {
		if value.Once {
			c.FunctionOnce[value.SrcName] = sync.Once{}
		}
	}
}

func (c *Container) initParser() {
	driver, flag := ParserDriver[c.UserOption.SourceGen]
	if !flag {
		cliLogger.Log.Fatalf("parse driver not exit, source gen %s", c.UserOption.SourceGen)
	}
	driver.RegisterOption(c.UserOption, c.TmplOption)
	c.Parser = driver
}

func (c *Container) initRender() {
	for _, desc := range c.TmplOption.Descriptor {
		_, allFlag := c.EnableModules["*"]
		_, moduleFlag := c.EnableModules[desc.Module]
		if !allFlag && !moduleFlag {
			continue
		}

		models := c.Parser.GetRenderInfos(desc)
		// model table name, model table schema
		for _, m := range models {
			// some render exec once
			syncOnce, flag := c.FunctionOnce[desc.SrcName]
			if flag {
				syncOnce.Do(func() {
					c.renderModel(m)
				})
				continue
			}
			c.renderModel(m)
		}
	}
}

func (c *Container) renderModel(m RenderInfo) {
	// todo optimize
	m.GenerateTime = c.GenerateTime
	render := NewRender(m)
	render.Exec(m.Descriptor.SrcName)
	if render.Descriptor.IsExistScript() {
		err := render.Descriptor.ExecScript(c.CurPath)
		if err != nil {
			cliLogger.Log.Fatalf("Bhojpur Web exec shell error, err: %s", err)
		}
	}
}

func (c *Container) initTimestamp() {
	if utils.IsExist(c.TimestampFile) {
		tree, err := toml.LoadFile(c.TimestampFile)
		if err != nil {
			cliLogger.Log.Fatalf("Bhojpur Web timestamp tmpl exec error, err: %s", err)
			return
		}
		err = tree.Unmarshal(&c.Timestamp)
		if err != nil {
			cliLogger.Log.Fatalf("Bhojpur Web timestamp tmpl parse error, err: %s", err)
			return
		}
	}
	c.Timestamp.Generate = c.GenerateTimeUnix
}

func (c *Container) flushTimestamp() {
	tomlByte, err := toml.Marshal(c.Timestamp)
	if err != nil {
		cliLogger.Log.Fatalf("marshal timestamp tmpl parse error, err: %s", err)
	}
	err = ioutil.WriteFile(c.TimestampFile, tomlByte, 0644)
	if err != nil {
		cliLogger.Log.Fatalf("flush timestamp tmpl parse error, err: %s", err)
	}
}

func (c *Container) InitToml() {
	if exist := utils.IsExist(c.BhojpurWebFile); exist {
		cliLogger.Log.Fatalf("file bhojpur.toml already exists")
	}
	sourceFile := c.UserOption.GitLocalPath + "/bhojpur.toml"
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		cliLogger.Log.Fatalf("read bhojpur.toml file err, %s", err.Error())
		return
	}
	err = ioutil.WriteFile(c.BhojpurWebFile, input, 0644)
	if err != nil {
		cliLogger.Log.Fatalf("create bhojpur.toml file err, %s", err.Error())
		return
	}
	cliLogger.Log.Success("Successfully created file bhojpur.toml")
}

//form https://github.com/bhojpur/web.git
//get bhojpur/web
func (c *Container) GetLocalPath() {
	if c.UserOption.GitLocalPath != "" {
		return
	}
	if GitRemotePath != "" {
		c.UserOption.GitRemotePath = GitRemotePath.String()
	}
	if c.UserOption.GitRemotePath == "" {
		c.UserOption.GitRemotePath = "https://github.com/bhojpur/web.git"
	}
	parse, err := url.Parse(c.UserOption.GitRemotePath)
	if err != nil {
		cliLogger.Log.Fatalf("git GitRemotePath err, %s", err.Error())
		return
	}
	s := parse.Path
	s = strings.TrimRight(s, ".git")
	c.UserOption.GitLocalPath = system.BhojpurHome + s
}
