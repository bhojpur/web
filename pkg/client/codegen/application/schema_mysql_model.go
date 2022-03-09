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
	"errors"

	cliLogger "github.com/bhojpur/web/pkg/client/logger"
)

type TableSchema struct {
	TableName              string
	ColumnName             string
	IsNullable             string
	DataType               string
	CharacterMaximumLength sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericScale           sql.NullInt64
	ColumnType             string
	ColumnKey              string
	Comment                string
}

type TableSchemas []TableSchema

func (tableSchemas TableSchemas) ToTableMap() (resp map[string]ModelInfos) {

	resp = make(map[string]ModelInfos)
	for _, value := range tableSchemas {
		if _, ok := resp[value.TableName]; !ok {
			resp[value.TableName] = make(ModelInfos, 0)
		}

		modelInfos := resp[value.TableName]
		inputType, goType, err := value.ToGoType()
		if err != nil {
			cliLogger.Log.Fatalf("parse go type err %s", err)
			return
		}

		modelInfo := ModelInfo{
			Name:      value.ColumnName,
			InputType: inputType,
			GoType:    goType,
			Comment:   value.Comment,
		}

		if value.ColumnKey == "PRI" {
			modelInfo.Orm = "pk"
		}
		resp[value.TableName] = append(modelInfos, modelInfo)
	}
	return
}

// GetGoDataType maps an SQL data type to Golang data type
func (col TableSchema) ToGoType() (inputType string, goType string, err error) {
	switch col.DataType {
	case "char", "varchar", "enum", "set", "text", "longtext", "mediumtext", "tinytext":
		goType = "string"
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		goType = "[]byte"
	case "date", "time", "datetime", "timestamp":
		goType, inputType = "time.Time", "dateTime"
	case "tinyint", "smallint", "int", "mediumint":
		goType = "int"
	case "bit", "bigint":
		goType = "int64"
	case "float", "decimal", "double":
		goType = "float64"
	}
	if goType == "" {
		err = errors.New("No compatible datatype (" + col.DataType + ", CamelName: " + col.ColumnName + ")  found")
	}
	return
}
