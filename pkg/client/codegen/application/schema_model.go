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
	"strings"

	"github.com/bhojpur/web/pkg/client/utils"
)

// parse get the model info
type ModelInfo struct {
	Name      string `json:"name"`      // mysql name
	InputType string `json:"inputType"` // user input type
	MysqlType string `json:"mysqlType"` // mysql type
	GoType    string `json:"goType"`    // go type
	Orm       string `json:"orm"`       // orm tag
	Comment   string `json:"comment"`   // mysql comment
	Extend    string `json:"extend"`    // user extend info
}

func (m ModelInfo) GetColumnKey() (columnKey string) {
	if m.InputType == "auto" || m.Orm == "pk" {
		columnKey = "PRI"
	}
	return
}

func (m ModelInfo) IsPrimaryKey() (flag bool) {
	if m.Orm == "auto" || m.Orm == "pk" || strings.ToLower(m.Name) == "id" {
		flag = true
	}
	return
}

type ModelInfos []ModelInfo

// to render model schemas
func (modelInfos ModelInfos) ToModelSchemas() (output ModelSchemas) {
	output = make(ModelSchemas, 0)
	for i, value := range modelInfos {
		if i == 0 && !value.IsPrimaryKey() {
			inputType, goType, mysqlType, ormTag := getModelType("auto")
			output = append(output, &ModelSchema{
				Name:      "id",
				InputType: inputType,
				ColumnKey: "PRI",
				Comment:   "ID",
				MysqlType: mysqlType,
				GoType:    goType,
				OrmTag:    ormTag,
				Extend:    value.Extend,
			})
		}

		modelSchema := &ModelSchema{
			Name:      value.Name,
			InputType: value.InputType,
			ColumnKey: value.GetColumnKey(),
			MysqlType: value.MysqlType,
			Comment:   value.Comment,
			GoType:    value.GoType,
			OrmTag:    value.Orm,
		}
		output = append(output, modelSchema)
	}
	return
}

type ModelSchema struct {
	Name      string // column name
	InputType string // user input type
	MysqlType string // mysql type
	ColumnKey string // PRI
	Comment   string // comment
	GoType    string // go type
	OrmTag    string // orm tag
	Extend    string
}

type ModelSchemas []*ModelSchema

func (m ModelSchemas) IsExistTime() bool {
	for _, value := range m {
		if value.InputType == "datetime" {
			return true
		}
	}
	return false
}

func (m ModelSchemas) GetPrimaryKey() string {
	camelPrimaryKey := ""
	for _, value := range m {
		if value.ColumnKey == "PRI" {
			camelPrimaryKey = utils.CamelString(value.Name)
		}
	}
	return camelPrimaryKey
}
