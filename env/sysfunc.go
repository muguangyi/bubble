// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package env

import (
	"fmt"
	"strings"
)

func GetSysfunc(name string) (MethodFunc, error) {
	f, ok := sysfuncs[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("Method [%s] is not exist!", name)
	}

	return f, nil
}

var (
	sysfuncs map[string]MethodFunc = make(map[string]MethodFunc)
)

func regMethod(name string, f MethodFunc) {
	sysfuncs[strings.ToLower(name)] = f
}
