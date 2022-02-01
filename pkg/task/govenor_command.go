package task

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
	"context"
	"fmt"
	"html/template"

	"github.com/pkg/errors"

	"github.com/bhojpur/web/pkg/core/admin"
)

type listTaskCommand struct {
}

func (l *listTaskCommand) Execute(params ...interface{}) *admin.Result {
	resultList := make([][]string, 0, len(globalTaskManager.adminTaskList))
	for tname, tk := range globalTaskManager.adminTaskList {
		result := []string{
			template.HTMLEscapeString(tname),
			template.HTMLEscapeString(tk.GetSpec(nil)),
			template.HTMLEscapeString(tk.GetStatus(nil)),
			template.HTMLEscapeString(tk.GetPrev(context.Background()).String()),
		}
		resultList = append(resultList, result)
	}

	return &admin.Result{
		Status:  200,
		Content: resultList,
	}
}

type runTaskCommand struct {
}

func (r *runTaskCommand) Execute(params ...interface{}) *admin.Result {
	if len(params) == 0 {
		return &admin.Result{
			Status: 400,
			Error:  errors.New("task name not passed"),
		}
	}

	tn, ok := params[0].(string)

	if !ok {
		return &admin.Result{
			Status: 400,
			Error:  errors.New("parameter is invalid"),
		}
	}

	if t, ok := globalTaskManager.adminTaskList[tn]; ok {
		err := t.Run(context.Background())
		if err != nil {
			return &admin.Result{
				Status: 500,
				Error:  err,
			}
		}
		return &admin.Result{
			Status:  200,
			Content: t.GetStatus(context.Background()),
		}
	} else {
		return &admin.Result{
			Status: 400,
			Error:  errors.New(fmt.Sprintf("task with name %s not found", tn)),
		}
	}

}

func registerCommands() {
	admin.RegisterCommand("task", "list", &listTaskCommand{})
	admin.RegisterCommand("task", "run", &runTaskCommand{})
}
