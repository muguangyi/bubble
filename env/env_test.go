// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package env

import (
	"testing"
	"time"
)

func TestFormatVariable(t *testing.T) {
	e := NewEnv()
	e.Set("x", NewAny(10))

	ret := e.Format(NewAny("$x"))
	if ret != "10" {
		t.Fail()
	}
}

func TestVariableWithBracket(t *testing.T) {
	e := NewEnv()
	e.Set("x", NewAny(10))

	ret := e.Format(NewAny("($x)"))
	if ret != "(10)" {
		t.Logf("Expect [(10)], but actual [%s]\n", ret)
		t.Fail()
	}
}

func TestFormatDateFunc(t *testing.T) {
	e := NewEnv()
	ret := e.Format(NewAny("$_DATE()"))

	if time.Now().Format("2006-01-02") != ret {
		t.Logf("Actual result: %s\n", ret)
		t.Fail()
	}
}

func TestFormatDateFuncWithFormat(t *testing.T) {
	e := NewEnv()
	ret := e.Format(NewAny("$_DATE(2006-01-02 15:04:05)"))

	expect := time.Now().Format("2006-01-02 15:04:05")
	if expect != ret {
		t.Logf("Expect [%s], but actual [%s]\n", expect, ret)
		t.Fail()
	}
}

func TestFormatDateFuncWithSpace(t *testing.T) {
	e := NewEnv()
	ret := e.Format(NewAny("$_DATE(   2006-01-02)"))

	expect := time.Now().Format("2006-01-02")
	if expect != ret {
		t.Logf("Expect [%s], but actual [%s]\n", expect, ret)
		t.Fail()
	}
}

func TestFuncWithVariable(t *testing.T) {
	e := NewEnv()
	e.Set("format", NewAny("2006-01-02 15:04:05"))
	ret := e.Format(NewAny("$_date($format)"))

	expect := time.Now().Format("2006-01-02 15:04:05")
	if expect != ret {
		t.Logf("Expect [%s], but actual [%s]\n", expect, ret)
		t.Fail()
	}
}

func TestFuncWithBracketInfo(t *testing.T) {
	e := NewEnv()

	ret := e.Format(NewAny("$_DATE().zip (please)"))
	expect := time.Now().Format("2006-01-02") + ".zip (please)"

	if expect != ret {
		t.Logf("Expect [%s], but actual [%s]\n", expect, ret)
		t.Fail()
	}
}

func TestNestedFunc(t *testing.T) {
	e := NewEnv()
	e.Set("VAR", NewAny(2))

	ret := e.Format(NewAny("$_ROUND($_Mul(5,$_Add(1,2),$VAR))"))

	if ret != "30" {
		t.Logf("Expect 30, but actual [%s]\n", ret)
		t.Fail()
	}
}

func TestSequenceFunc(t *testing.T) {
	e := NewEnv()
	e.Set("VAR", NewAny(2))

	ret := e.Format(NewAny("$_ROUND($_Add(1,2)), $_ROUND($_SUB(5, 4))"))

	if ret != "3, 1" {
		t.Logf("Expect [3, 1], but actual [%s]\n", ret)
		t.Fail()
	}
}

func TestDiv(t *testing.T) {
	e := NewEnv()
	e.SetFunc("_SIZEOF", func(args ...IAny) (IAny, error) { return NewAny(1024 * 1024), nil })

	ret := e.Format(NewAny("is $_ROUND($_DIV($_SIZEOF(./README.md), 1024, 1024)) m"))

	if ret != "is 1 m" {
		t.Logf("Expect [is 1 m], but actual [%s]\n", ret)
		t.Fail()
	}
}

func TestValueFuncAsOneParameter(t *testing.T) {
	e := NewEnv()
	e.SetFunc("_SIZEOF", func(args ...IAny) (IAny, error) {
		if args[0].ToString() != "./README3" {
			t.Logf("Expect [./README3], but actual [%s]\n", args[0].ToString())
			t.Fail()
		}

		return NewAny(1024), nil
	})

	ret := e.Format(NewAny("$_SIZEOF(./README$_ROUND($_ADD(1,2)))"))
	if ret != "1024" {
		t.Logf("Expect [1024], but actual [%s]\n", ret)
		t.Fail()
	}
}
