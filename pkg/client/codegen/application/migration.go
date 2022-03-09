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
	"database/sql"
	"io/ioutil"
	"path/filepath"
	"strings"

	cliLogger "github.com/bhojpur/web/pkg/client/logger"
	"github.com/bhojpur/web/pkg/client/utils"
)

var SQL utils.DocValue
var SQLMode utils.DocValue
var SQLModePath utils.DocValue

var (
	SQLModeUp   = "up"
	SQLModeDown = "down"
)

func (c *Container) Migration(args []string) {
	c.initUserOption()
	db, err := sql.Open(c.UserOption.Driver, c.UserOption.Dsn)
	if err != nil {
		cliLogger.Log.Fatalf("Could not connect to '%s' database using '%s': %s", c.UserOption.Driver, c.UserOption.Dsn, err)
		return
	}
	defer db.Close()
	switch SQLMode.String() {
	case SQLModeUp:
		doByMode(db, "up.sql")
	case SQLModeDown:
		doByMode(db, "down.sql")
	default:
		doBySqlFile(db)
	}
}

func doBySqlFile(db *sql.DB) {
	fileName := SQL.String()
	if !utils.IsExist(fileName) {
		cliLogger.Log.Fatalf("sql mode path not exist, path %s", SQL.String())
	}
	doDb(db, fileName)
}

func doByMode(db *sql.DB, suffix string) {
	pathName := SQLModePath.String()
	if !utils.IsExist(pathName) {
		cliLogger.Log.Fatalf("sql mode path not exist, path %s", SQLModePath.String())
	}

	rd, err := ioutil.ReadDir(pathName)
	if err != nil {
		cliLogger.Log.Fatalf("read dir err, path %s, err %s", pathName, err)
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			if !strings.HasSuffix(fi.Name(), suffix) {
				continue
			}
			doDb(db, filepath.Join(pathName, fi.Name()))
		}
	}
}

func doDb(db *sql.DB, filePath string) {
	absFile, _ := filepath.Abs(filePath)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		cliLogger.Log.Errorf("read file err %s, abs file %s", err, absFile)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		cliLogger.Log.Errorf("db exec err %s", err)
	}
	cliLogger.Log.Infof("db exec info %s", filePath)
}
