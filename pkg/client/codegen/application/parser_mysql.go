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
	"fmt"

	"github.com/bhojpur/web/pkg/client/codegen/utils"
	cliLogger "github.com/bhojpur/web/pkg/client/logger"
)

type MysqlParser struct {
	userOption UserOption
	tmplOption TmplOption
}

func (m *MysqlParser) RegisterOption(userOption UserOption, tmplOption TmplOption) {
	m.userOption = userOption
	m.tmplOption = tmplOption

}

func (*MysqlParser) Parse(descriptor Descriptor) {

}

func (m *MysqlParser) GetRenderInfos(descriptor Descriptor) (output []RenderInfo) {
	tableSchemas, err := m.getTableSchemas()
	if err != nil {
		cliLogger.Log.Fatalf("get table schemas err %s", err)
	}
	models := tableSchemas.ToTableMap()

	output = make([]RenderInfo, 0)
	// model table name, model table schema
	for modelName, content := range models {
		output = append(output, RenderInfo{
			Module:     descriptor.Module,
			ModelName:  modelName,
			Content:    content,
			Option:     m.userOption,
			Descriptor: descriptor,
			TmplPath:   m.tmplOption.RenderPath,
		})
	}
	return
}

func (t *MysqlParser) Unregister() {

}

func (m *MysqlParser) getTableSchemas() (resp TableSchemas, err error) {
	dsn, err := utils.ParseDSN(m.userOption.Dsn)
	if err != nil {
		cliLogger.Log.Fatalf("parse dsn err %s", err)
		return
	}

	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/information_schema", dsn.User, dsn.Passwd, dsn.Addr))
	if err != nil {
		cliLogger.Log.Fatalf("Could not connect to mysql database using '%s': %s", m.userOption.Dsn, err)
		return
	}
	defer conn.Close()

	q := `SELECT TABLE_NAME, COLUMN_NAME, IS_NULLABLE, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, 
NUMERIC_PRECISION, NUMERIC_SCALE,COLUMN_TYPE,COLUMN_KEY,COLUMN_COMMENT 
FROM COLUMNS WHERE TABLE_SCHEMA = ?  ORDER BY TABLE_NAME, ORDINAL_POSITION`
	rows, err := conn.Query(q, dsn.DBName)
	if err != nil {
		return nil, err
	}
	columns := make(TableSchemas, 0)
	for rows.Next() {
		cs := TableSchema{}
		err := rows.Scan(&cs.TableName, &cs.ColumnName, &cs.IsNullable, &cs.DataType,
			&cs.CharacterMaximumLength, &cs.NumericPrecision, &cs.NumericScale,
			&cs.ColumnType, &cs.ColumnKey, &cs.Comment)
		if err != nil {
			return nil, err
		}
		columns = append(columns, cs)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return columns, nil
}
