// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package env

// MethodFunc is sys func definition.
type MethodFunc func(args ...IAny) (IAny, error)

// IEnv interface.
type IEnv interface {
	// Format code with variable value if there are variables in code.
	Format(code IAny) string

	// Get target name variable.
	Get(name string) (IAny, error)

	// Set target name variable with value.
	Set(name string, value IAny) error

	// GetFunc try to get target func.
	GetFunc(name string) (MethodFunc, error)

	// SetFunc register a function with name as key.
	SetFunc(name string, f MethodFunc) error

	// FromBytes fill data with []byte.
	FromBytes(bytes []byte) error

	// ToBytes convert env data into []byte.
	ToBytes() ([]byte, error)
}
