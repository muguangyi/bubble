// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package action

import (
	"bubble/env"
)

// IAction is Action interface.
type IAction interface {
	// Init is used for initializing the Action.
	Init(uid uint64, env env.IEnv)

	// Getcwd returns the current working directory.
	Cwd() string

	// Execute the Action in a coroutine.
	Execute(script env.IAny, target string, env env.IEnv, log ILog) chan bool

	// Cancel the Action execution.
	Cancel() error

	// Error returns the error message string if exists.
	Error() string
}
