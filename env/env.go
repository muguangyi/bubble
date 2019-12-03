// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package env

import (
	"errors"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func NewEnv() IEnv {
	return &env{vars: make(map[string]IAny), funcs: make(map[string]MethodFunc)}
}

type env struct {
	vars  map[string]IAny
	funcs map[string]MethodFunc
}

func (e *env) Format(code IAny) string {
	l := NewLexer()
	p := NewParser()
	l.Parse(code.ToString(), p)

	format := ""
	n := p.Result()
	for n != nil {
		format += n.Execute(e)
		n = n.Next()
	}

	return format
}

func (e *env) Get(name string) (IAny, error) {
	if strings.HasPrefix(name, "$") {
		name = name[1:]
	}

	v, ok := e.vars[strings.ToLower(name)]
	if !ok {
		// Try get variable from system.
		return GetSysvar(name)
	}

	return v, nil
}

func (e *env) Set(name string, value IAny) error {
	e.vars[strings.ToLower(name)] = value
	return nil
}

func (e *env) GetFunc(name string) (MethodFunc, error) {
	if strings.HasPrefix(name, "$") {
		name = name[1:]
	}

	f, ok := e.funcs[strings.ToLower(name)]
	if !ok {
		var err error
		if f, err = GetSysfunc(name); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (e *env) SetFunc(name string, f MethodFunc) error {
	e.funcs[strings.ToLower(name)] = f
	return nil
}

func (e *env) FromBytes(bytes []byte) error {
	if bytes == nil {
		return nil
	}

	a := NewAny(nil)
	if err := a.FromBytes(bytes); err != nil {
		return err
	}

	if !a.IsMap() {
		return errors.New("The env data is not a map!")
	}

	e.vars = a.Map()
	return nil
}

func (e *env) ToBytes() ([]byte, error) {
	data := make(map[string]interface{})
	for k, v := range e.vars {
		data[k] = v.(*any).source
	}

	return yaml.Marshal(data)
}
