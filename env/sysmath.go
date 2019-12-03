// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package env

import (
	"fmt"
	"math"
)

func init() {
	regMethod("_ADD", _Add)
	regMethod("_SUB", _Sub)
	regMethod("_MUL", _Mul)
	regMethod("_DIV", _Div)
	regMethod("_ROUND", _Round)
	regMethod("_CEIL", _Ceil)
	regMethod("_FLOOR", _Floor)
}

func _Add(args ...IAny) (IAny, error) {
	var total float64 = 0
	for _, a := range args {
		total += a.Float()
	}

	return NewAny(total), nil
}

func _Sub(args ...IAny) (IAny, error) {
	if len(args) == 0 {
		return NewAny(0), nil
	}

	var total float64 = args[0].Float()
	for i := 1; i < len(args); i++ {
		total -= args[i].Float()
	}

	return NewAny(total), nil
}

func _Mul(args ...IAny) (IAny, error) {
	var total float64 = 1
	for _, a := range args {
		total *= a.Float()
	}

	return NewAny(total), nil
}

func _Div(args ...IAny) (IAny, error) {
	if len(args) == 0 {
		return NewAny(0), nil
	}

	var total float64 = args[0].Float()
	for i := 1; i < len(args); i++ {
		r := args[i].Float()
		if r == 0 {
			return NewAny(0), fmt.Errorf("Div by zero for arg [%f]\n", r)
		}

		total /= r
	}

	return NewAny(total), nil
}

func _Round(args ...IAny) (IAny, error) {
	if len(args) == 0 {
		return NewAny(0), nil
	}

	return NewAny(int(math.Round(args[0].Float()))), nil
}

func _Ceil(args ...IAny) (IAny, error) {
	if len(args) == 0 {
		return NewAny(0), nil
	}

	return NewAny(int(math.Ceil(args[0].Float()))), nil
}

func _Floor(args ...IAny) (IAny, error) {
	if len(args) == 0 {
		return NewAny(0), nil
	}

	return NewAny(int(math.Floor(args[0].Float()))), nil
}
