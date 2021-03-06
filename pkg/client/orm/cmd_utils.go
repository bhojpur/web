package orm

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
	"strings"
)

type dbIndex struct {
	Table string
	Name  string
	SQL   string
}

// get database column type string.
func getColumnTyp(al *alias, fi *fieldInfo) (col string) {
	T := al.DbBaser.DbTypes()
	fieldType := fi.fieldType
	fieldSize := fi.size

checkColumn:
	switch fieldType {
	case TypeBooleanField:
		col = T["bool"]
	case TypeVarCharField:
		if al.Driver == DRPostgres && fi.toText {
			col = T["string-text"]
		} else {
			col = fmt.Sprintf(T["string"], fieldSize)
		}
	case TypeCharField:
		col = fmt.Sprintf(T["string-char"], fieldSize)
	case TypeTextField:
		col = T["string-text"]
	case TypeTimeField:
		col = T["time.Time-clock"]
	case TypeDateField:
		col = T["time.Time-date"]
	case TypeDateTimeField:
		// the precision of sqlite is not implemented
		if al.Driver == 2 || fi.timePrecision == nil {
			col = T["time.Time"]
		} else {
			s := T["time.Time-precision"]
			col = fmt.Sprintf(s, *fi.timePrecision)
		}

	case TypeBitField:
		col = T["int8"]
	case TypeSmallIntegerField:
		col = T["int16"]
	case TypeIntegerField:
		col = T["int32"]
	case TypeBigIntegerField:
		if al.Driver == DRSqlite {
			fieldType = TypeIntegerField
			goto checkColumn
		}
		col = T["int64"]
	case TypePositiveBitField:
		col = T["uint8"]
	case TypePositiveSmallIntegerField:
		col = T["uint16"]
	case TypePositiveIntegerField:
		col = T["uint32"]
	case TypePositiveBigIntegerField:
		col = T["uint64"]
	case TypeFloatField:
		col = T["float64"]
	case TypeDecimalField:
		s := T["float64-decimal"]
		if !strings.Contains(s, "%d") {
			col = s
		} else {
			col = fmt.Sprintf(s, fi.digits, fi.decimals)
		}
	case TypeJSONField:
		if al.Driver != DRPostgres {
			fieldType = TypeVarCharField
			goto checkColumn
		}
		col = T["json"]
	case TypeJsonbField:
		if al.Driver != DRPostgres {
			fieldType = TypeVarCharField
			goto checkColumn
		}
		col = T["jsonb"]
	case RelForeignKey, RelOneToOne:
		fieldType = fi.relModelInfo.fields.pk.fieldType
		fieldSize = fi.relModelInfo.fields.pk.size
		goto checkColumn
	}

	return
}

// create alter sql string.
func getColumnAddQuery(al *alias, fi *fieldInfo) string {
	Q := al.DbBaser.TableQuote()
	typ := getColumnTyp(al, fi)

	if !fi.null {
		typ += " " + "NOT NULL"
	}

	return fmt.Sprintf("ALTER TABLE %s%s%s ADD COLUMN %s%s%s %s %s",
		Q, fi.mi.table, Q,
		Q, fi.column, Q,
		typ, getColumnDefault(fi),
	)
}

// Get string value for the attribute "DEFAULT" for the CREATE, ALTER commands
func getColumnDefault(fi *fieldInfo) string {
	var (
		v, t, d string
	)

	// Skip default attribute if field is in relations
	if fi.rel || fi.reverse {
		return v
	}

	t = " DEFAULT '%s' "

	// These defaults will be useful if there no config value orm:"default" and NOT NULL is on
	switch fi.fieldType {
	case TypeTimeField, TypeDateField, TypeDateTimeField, TypeTextField:
		return v

	case TypeBitField, TypeSmallIntegerField, TypeIntegerField,
		TypeBigIntegerField, TypePositiveBitField, TypePositiveSmallIntegerField,
		TypePositiveIntegerField, TypePositiveBigIntegerField, TypeFloatField,
		TypeDecimalField:
		t = " DEFAULT %s "
		d = "0"
	case TypeBooleanField:
		t = " DEFAULT %s "
		d = "FALSE"
	case TypeJSONField, TypeJsonbField:
		d = "{}"
	}

	if fi.colDefault {
		if !fi.initial.Exist() {
			v = fmt.Sprintf(t, "")
		} else {
			v = fmt.Sprintf(t, fi.initial.String())
		}
	} else {
		if !fi.null {
			v = fmt.Sprintf(t, d)
		}
	}

	return v
}
