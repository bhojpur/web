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
	cliLogger "github.com/bhojpur/web/pkg/client/logger"
)

type TextModel struct {
	Names    []string
	Orms     []string
	Comments []string
	Extends  []string
}

func (content TextModel) ToModelInfos() (output []ModelInfo) {
	namesLen := len(content.Names)
	ormsLen := len(content.Orms)
	commentsLen := len(content.Comments)
	if namesLen != ormsLen && namesLen != commentsLen {
		cliLogger.Log.Fatalf("length error, namesLen is %d, ormsLen is %d, commentsLen is %d", namesLen, ormsLen, commentsLen)
	}
	extendLen := len(content.Extends)
	if extendLen != 0 && extendLen != namesLen {
		cliLogger.Log.Fatalf("extend length error, namesLen is %d, extendsLen is %d", namesLen, extendLen)
	}

	output = make([]ModelInfo, 0)
	for i, name := range content.Names {
		comment := content.Comments[i]
		if comment == "" {
			comment = name
		}
		inputType, goType, mysqlType, ormTag := getModelType(content.Orms[i])

		m := ModelInfo{
			Name:      name,
			InputType: inputType,
			GoType:    goType,
			Orm:       ormTag,
			Comment:   comment,
			MysqlType: mysqlType,
			Extend:    "",
		}
		// extend value
		if extendLen != 0 {
			m.Extend = content.Extends[i]
		}
		output = append(output, m)
	}
	return
}
