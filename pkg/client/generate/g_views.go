package generate

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
	"os"
	"path"

	cliLogger "github.com/bhojpur/web/pkg/client/logger"
	"github.com/bhojpur/web/pkg/client/logger/colors"
	"github.com/bhojpur/web/pkg/client/utils"
)

// recipe
// admin/recipe
func GenerateView(viewpath, currpath string) {
	w := colors.NewColorWriter(os.Stdout)

	cliLogger.Log.Info("Generating view...")

	absViewPath := path.Join(currpath, "views", viewpath)
	err := os.MkdirAll(absViewPath, os.ModePerm)
	if err != nil {
		cliLogger.Log.Fatalf("Could not create '%s' view: %s", viewpath, err)
	}

	cfile := path.Join(absViewPath, "index.tpl")
	if f, err := os.OpenFile(cfile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
		defer utils.CloseFile(f)
		f.WriteString(cfile)
		fmt.Fprintf(w, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", cfile, "\x1b[0m")
	} else {
		cliLogger.Log.Fatalf("Could not create view file: %s", err)
	}

	cfile = path.Join(absViewPath, "show.tpl")
	if f, err := os.OpenFile(cfile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
		defer utils.CloseFile(f)
		f.WriteString(cfile)
		fmt.Fprintf(w, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", cfile, "\x1b[0m")
	} else {
		cliLogger.Log.Fatalf("Could not create view file: %s", err)
	}

	cfile = path.Join(absViewPath, "create.tpl")
	if f, err := os.OpenFile(cfile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
		defer utils.CloseFile(f)
		f.WriteString(cfile)
		fmt.Fprintf(w, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", cfile, "\x1b[0m")
	} else {
		cliLogger.Log.Fatalf("Could not create view file: %s", err)
	}

	cfile = path.Join(absViewPath, "edit.tpl")
	if f, err := os.OpenFile(cfile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
		defer utils.CloseFile(f)
		f.WriteString(cfile)
		fmt.Fprintf(w, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", cfile, "\x1b[0m")
	} else {
		cliLogger.Log.Fatalf("Could not create view file: %s", err)
	}
}
