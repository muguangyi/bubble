// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package env

import (
	"testing"
	"time"
)

func TestFuncName(t *testing.T) {
	_, err := GetSysfunc("_DATE")
	if err != nil {
		t.Fail()
	}

	_, err = GetSysfunc("_date")
	if err != nil {
		t.Fail()
	}

	_, err = GetSysfunc("_Date")
	if err != nil {
		t.Fail()
	}

	_, err = GetSysfunc("Date")
	if err == nil {
		t.Fail()
	}
}

func TestDate(t *testing.T) {
	f, err := GetSysfunc("_DATE")
	if err != nil {
		t.Error(err)
	}

	ret, err := f()
	if err != nil {
		t.Error(err)
	}

	if time.Now().Format("2006-01-02") != ret.ToString() {
		t.Fail()
	}
}

func TestAddInt(t *testing.T) {
	f, err := GetSysfunc("_add")
	if err != nil {
		t.Error(err)
	}

	ret, err := f(NewAny(1), NewAny(2), NewAny(3))
	if err != nil {
		t.Error(err)
	}

	if ret.Float() != 6 {
		t.Fail()
	}
}

func TestAddFloat(t *testing.T) {
	f, err := GetSysfunc("_add")
	if err != nil {
		t.Error(err)
	}

	ret, err := f(NewAny(1.1), NewAny(2), NewAny(3.4))
	if err != nil {
		t.Error(err)
	}

	if ret.Float() != 6.5 {
		t.Fail()
	}

	if ret.ToString() != "6.50000" {
		t.Fail()
	}
}

func TestRound(t *testing.T) {
	f, err := GetSysfunc("_round")
	if err != nil {
		t.Error(err)
	}

	ret, err := f(NewAny(1.53213))
	if err != nil {
		t.Error(err)
	}

	if ret.Int() != 2 {
		t.Fail()
	}
}
