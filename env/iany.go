// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package env

// IAny interface to present a yaml object.
type IAny interface {
	// IsNil returns whether the value is nil value.
	IsNil() bool

	// IsString check if the value is a string type.
	IsString() bool

	// IsArr returns whether the value is an array.
	IsArr() bool

	// IsMap returns whether the value is a map.
	IsMap() bool

	// Bool converts the value to bool.
	Bool() bool

	// String converts the value to string.
	String() string

	// Int converts the value to int.
	Int() int

	// Uint converts the value to uint.
	Uint() uint

	// Float converts the value to float64.
	Float() float64

	// Array converts the value to []IAny.
	Array() []IAny

	// Map converts the value to map[string]IAny.
	Map() map[string]IAny

	// FromBytes imports []byte into the value.
	FromBytes(bytes []byte) error

	// ToBytes outputs the value into []byte.
	ToBytes() ([]byte, error)

	// ToString try to convert value to string format.
	ToString() string
}
