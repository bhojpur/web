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
	"reflect"
	"time"
)

// get table alias.
func getDbAlias(name string) *alias {
	if al, ok := dataBaseCache.get(name); ok {
		return al
	}
	panic(fmt.Errorf("unknown DataBase alias name %s", name))
}

// get pk column info.
func getExistPk(mi *modelInfo, ind reflect.Value) (column string, value interface{}, exist bool) {
	fi := mi.fields.pk

	v := ind.FieldByIndex(fi.fieldIndex)
	if fi.fieldType&IsPositiveIntegerField > 0 {
		vu := v.Uint()
		exist = vu > 0
		value = vu
	} else if fi.fieldType&IsIntegerField > 0 {
		vu := v.Int()
		exist = true
		value = vu
	} else if fi.fieldType&IsRelField > 0 {
		_, value, exist = getExistPk(fi.relModelInfo, reflect.Indirect(v))
	} else {
		vu := v.String()
		exist = vu != ""
		value = vu
	}

	column = fi.column
	return
}

// get fields description as flatted string.
func getFlatParams(fi *fieldInfo, args []interface{}, tz *time.Location) (params []interface{}) {

outFor:
	for _, arg := range args {
		val := reflect.ValueOf(arg)

		if arg == nil {
			params = append(params, arg)
			continue
		}

		kind := val.Kind()
		if kind == reflect.Ptr {
			val = val.Elem()
			kind = val.Kind()
			arg = val.Interface()
		}

		switch kind {
		case reflect.String:
			v := val.String()
			if fi != nil {
				if fi.fieldType == TypeTimeField || fi.fieldType == TypeDateField || fi.fieldType == TypeDateTimeField {
					var t time.Time
					var err error
					if len(v) >= 19 {
						s := v[:19]
						t, err = time.ParseInLocation(formatDateTime, s, DefaultTimeLoc)
					} else if len(v) >= 10 {
						s := v
						if len(v) > 10 {
							s = v[:10]
						}
						t, err = time.ParseInLocation(formatDate, s, tz)
					} else {
						s := v
						if len(s) > 8 {
							s = v[:8]
						}
						t, err = time.ParseInLocation(formatTime, s, tz)
					}
					if err == nil {
						if fi.fieldType == TypeDateField {
							v = t.In(tz).Format(formatDate)
						} else if fi.fieldType == TypeDateTimeField {
							v = t.In(tz).Format(formatDateTime)
						} else {
							v = t.In(tz).Format(formatTime)
						}
					}
				}
			}
			arg = v
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			arg = val.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			arg = val.Uint()
		case reflect.Float32:
			arg, _ = StrTo(ToStr(arg)).Float64()
		case reflect.Float64:
			arg = val.Float()
		case reflect.Bool:
			arg = val.Bool()
		case reflect.Slice, reflect.Array:
			if _, ok := arg.([]byte); ok {
				continue outFor
			}

			var args []interface{}
			for i := 0; i < val.Len(); i++ {
				v := val.Index(i)

				var vu interface{}
				if v.CanInterface() {
					vu = v.Interface()
				}

				if vu == nil {
					continue
				}

				args = append(args, vu)
			}

			if len(args) > 0 {
				p := getFlatParams(fi, args, tz)
				params = append(params, p...)
			}
			continue outFor
		case reflect.Struct:
			if v, ok := arg.(time.Time); ok {
				if fi != nil && fi.fieldType == TypeDateField {
					arg = v.In(tz).Format(formatDate)
				} else if fi != nil && fi.fieldType == TypeDateTimeField {
					arg = v.In(tz).Format(formatDateTime)
				} else if fi != nil && fi.fieldType == TypeTimeField {
					arg = v.In(tz).Format(formatTime)
				} else {
					arg = v.In(tz).Format(formatDateTime)
				}
			} else {
				typ := val.Type()
				name := getFullName(typ)
				var value interface{}
				if mmi, ok := modelCache.getByFullName(name); ok {
					if _, vu, exist := getExistPk(mmi, val); exist {
						value = vu
					}
				}
				arg = value

				if arg == nil {
					panic(fmt.Errorf("need a valid args value, unknown table or value `%s`", name))
				}
			}
		}

		params = append(params, arg)
	}
	return
}
