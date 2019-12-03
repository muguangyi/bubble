// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package env

import (
	"regexp"
)

type TokenType int

const (
	TOKEN_ERROR TokenType = iota
	TOKEN_VALUE
	TOKEN_VARIABLE
	TOKEN_BEGIN_METHOD
	TOKEN_END_METHOD
	TOKEN_END_PARAM
)

const (
	PREFIX        = '$'
	LEFT_BRACKET  = '('
	RIGHT_BRACKET = ')'
	COMMA         = ','
	SPACE         = ' '
)

var (
	NameExp = regexp.MustCompile(`[\w]+`)
)
