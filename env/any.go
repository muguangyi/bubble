// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package env

import (
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func Load(file string) (IAny, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	s := NewAny(nil)
	s.FromBytes(bytes)

	return s, nil
}

func NewAny(source interface{}) IAny {
	s := &any{source: source, t: UNKNOWN}
	s.init()

	return s
}

type TYPE uint

const (
	UNKNOWN TYPE = iota
	NIL
	BOOL
	INT
	INT8
	INT16
	INT32
	INT64
	UINT
	UINT8
	UINT16
	UINT32
	UINT64
	STRING
	FLOAT32
	FLOAT64
	ARRAY
	MAP
)

type any struct {
	source interface{}
	bytes  []byte
	t      TYPE
}

func (a *any) IsNil() bool {
	return a.source == nil
}

func (a *any) IsString() bool {
	return a.t == STRING
}

func (a *any) IsArr() bool {
	return a.t == ARRAY
}

func (a *any) IsMap() bool {
	return a.t == MAP
}

func (a *any) Bool() bool {
	return a.source.(bool)
}

func (a *any) String() string {
	return a.source.(string)
}

func (a *any) Int() int {
	switch a.t {
	case INT:
		return a.source.(int)
	case INT8:
		return int(a.source.(int8))
	case INT16:
		return int(a.source.(int16))
	case INT32:
		return int(a.source.(int32))
	case INT64:
		return int(a.source.(int64))
	case UINT:
		return int(a.source.(uint))
	case UINT8:
		return int(a.source.(uint8))
	case UINT16:
		return int(a.source.(uint16))
	case UINT32:
		return int(a.source.(uint32))
	case UINT64:
		return int(a.source.(uint64))
	case FLOAT32:
		return int(a.source.(float32))
	case FLOAT64:
		return int(a.source.(float64))
	case STRING:
		i, _ := strconv.ParseInt(a.source.(string), 10, 64)
		return int(i)
	}

	panic("Can't convert to int!")
}

func (a *any) Uint() uint {
	switch a.t {
	case INT:
		return uint(a.source.(int))
	case INT8:
		return uint(a.source.(int8))
	case INT16:
		return uint(a.source.(int16))
	case INT32:
		return uint(a.source.(int32))
	case INT64:
		return uint(a.source.(int64))
	case UINT:
		return a.source.(uint)
	case UINT8:
		return uint(a.source.(uint8))
	case UINT16:
		return uint(a.source.(uint16))
	case UINT32:
		return uint(a.source.(uint32))
	case UINT64:
		return uint(a.source.(uint64))
	case FLOAT32:
		return uint(a.source.(float32))
	case FLOAT64:
		return uint(a.source.(float64))
	case STRING:
		ui, _ := strconv.ParseUint(a.source.(string), 10, 64)
		return uint(ui)
	}

	panic("Can't convert to uint!")
}

func (a *any) Float() float64 {
	switch a.t {
	case INT:
		return float64(a.source.(int))
	case INT8:
		return float64(a.source.(int8))
	case INT16:
		return float64(a.source.(int16))
	case INT32:
		return float64(a.source.(int32))
	case INT64:
		return float64(a.source.(int64))
	case UINT:
		return float64(a.source.(uint))
	case UINT8:
		return float64(a.source.(uint8))
	case UINT16:
		return float64(a.source.(uint16))
	case UINT32:
		return float64(a.source.(uint32))
	case UINT64:
		return float64(a.source.(uint64))
	case FLOAT32:
		return float64(a.source.(float32))
	case FLOAT64:
		return a.source.(float64)
	case STRING:
		f, _ := strconv.ParseFloat(a.source.(string), 64)
		return f
	}

	panic("Can't convert to uint!")
}

func (a *any) Array() []IAny {
	if a.t == ARRAY {
		arr := a.source.([]interface{})
		ret := make([]IAny, len(arr))
		for i, v := range arr {
			ret[i] = NewAny(v)
		}

		return ret
	}

	return nil
}

func (a *any) Map() map[string]IAny {
	if a.t == MAP {
		dict := a.source.(map[interface{}]interface{})
		ret := make(map[string]IAny)
		for k, v := range dict {
			ret[k.(string)] = NewAny(v)
		}

		return ret
	}

	return nil
}

func (a *any) FromBytes(bytes []byte) error {
	a.bytes = bytes

	err := yaml.Unmarshal(bytes, &a.source)
	if err != nil {
		return err
	}

	a.init()
	return nil
}

func (a *any) ToBytes() ([]byte, error) {
	var err error
	if a.bytes == nil {
		a.bytes, err = yaml.Marshal(a.source)
	}

	return a.bytes, err
}

func (a *any) ToString() string {
	switch a.t {
	case BOOL:
		return strconv.FormatBool(a.source.(bool))
	case STRING:
		return a.source.(string)
	case INT:
		return strconv.FormatInt(int64(a.source.(int)), 10)
	case INT8:
		return strconv.FormatInt(int64(a.source.(int8)), 10)
	case INT16:
		return strconv.FormatInt(int64(a.source.(int16)), 10)
	case INT32:
		return strconv.FormatInt(int64(a.source.(int32)), 10)
	case INT64:
		return strconv.FormatInt(a.source.(int64), 10)
	case UINT:
		return strconv.FormatUint(uint64(a.source.(uint)), 10)
	case UINT8:
		return strconv.FormatUint(uint64(a.source.(uint8)), 10)
	case UINT16:
		return strconv.FormatUint(uint64(a.source.(uint16)), 10)
	case UINT32:
		return strconv.FormatUint(uint64(a.source.(uint32)), 10)
	case UINT64:
		return strconv.FormatUint(a.source.(uint64), 10)
	case FLOAT32:
		return strconv.FormatFloat(float64(a.source.(float32)), 'f', 5, 32)
	case FLOAT64:
		return strconv.FormatFloat(a.source.(float64), 'f', 5, 64)
	case ARRAY:
		{
			strs := make([]string, 0)
			arr := a.Array()
			for _, i := range arr {
				strs = append(strs, i.ToString())
			}

			return "[" + strings.Join(strs, ",") + "]"
		}
	case MAP:
		{
			strs := make([]string, 0)
			m := a.Map()
			for k, v := range m {
				strs = append(strs, k+":"+v.ToString())
			}

			return "{" + strings.Join(strs, ",") + "}"
		}
	}

	return ""
}

func (a *any) init() {
	if a.source == nil {
		a.t = NIL
	}

	switch reflect.ValueOf(a.source).Kind() {
	case reflect.Bool:
		a.t = BOOL
	case reflect.String:
		a.t = STRING
	case reflect.Int:
		a.t = INT
	case reflect.Int8:
		a.t = INT8
	case reflect.Int16:
		a.t = INT16
	case reflect.Int32:
		a.t = INT32
	case reflect.Int64:
		a.t = INT64
	case reflect.Uint:
		a.t = UINT
	case reflect.Uint8:
		a.t = UINT8
	case reflect.Uint16:
		a.t = UINT16
	case reflect.Uint32:
		a.t = UINT32
	case reflect.Uint64:
		a.t = UINT64
	case reflect.Float32:
		a.t = FLOAT32
	case reflect.Float64:
		a.t = FLOAT64
	case reflect.Array, reflect.Slice:
		a.t = ARRAY
	case reflect.Map:
		a.t = MAP
	}
}
