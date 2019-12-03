// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package env

import (
	"fmt"
	"strings"
)

func GetSysvar(name string) (IAny, error) {
	f, ok := sysvars[name]
	if !ok {
		return nil, fmt.Errorf("Variable [%s] is not exist!", name)
	}

	return f()
}

type VarFunc func() (IAny, error)

var (
	sysvars map[string]VarFunc = make(map[string]VarFunc)
)

func regVar(name string, f VarFunc) {
	sysvars[strings.ToLower(name)] = f
}
