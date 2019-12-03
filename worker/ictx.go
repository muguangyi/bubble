// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package worker

import (
	"bubble/env"
)

// ICtx presents a Job execution info.
type ICtx interface {
	// Master returns Master id.
	Master() uint64

	// UID returns the unique task id.
	UID() uint64

	// Script return Job script data.
	Script() env.IAny

	// Variables return script variables.
	Variables() env.IAny

	// Target return Job running target.
	Target() string

	// Env return Job running env.
	Env() env.IEnv
}
