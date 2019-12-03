// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package env

import (
	"time"
)

func init() {
	regMethod("_DATE", _Date)
}

func _Date(args ...IAny) (IAny, error) {
	var format string = "2006-01-02"
	if len(args) > 0 {
		format = args[0].ToString()
	}
	return NewAny(time.Now().Format(format)), nil
}
