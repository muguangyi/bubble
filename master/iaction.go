// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

import (
	"bubble/env"
)

// IAction to abstract an Action feature on a Worker.
type IAction interface {
	// Target returns the sub infos the Action supports.
	Target() []env.IAny

	// Prefer returns the prefer ability the Action supports.
	Prefer() []env.IAny

	// Execute ICtx.
	Execute(ctx ICtx)

	// Finish the Action with result and env.
	Finish(runner uint64, success bool, env env.IEnv) error

	// Cancel the target job.
	Cancel(runner uint64) error

	// Progress the target job status with payload data.
	Progress(runner uint64, payload []byte) error

	// Destroy the Action.
	Destroy()
}
