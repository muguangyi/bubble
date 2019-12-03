// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

import (
	"bubble/def"
	"bubble/env"
)

// IWorker interface.
type IWorker interface {
	// ID returns the Worker id.
	ID() uint64

	// Bind target Worker info to local Worker proxy.
	Bind(supports map[string][]byte) error

	// Satisfy returns whether could support the target command.
	Satisfy(command ICommand) bool

	// Workload returns the Worker running command quantity.
	Workload() int

	// Get target Action.
	Get(name string) IAction

	// Finish action with related parameters.
	Finish(action string, runner uint64, success bool, env env.IEnv) error

	// Notify action with related parameters.
	Progress(action string, runner uint64, payload []byte) error

	// Broadcast handles data from corresponding Worker.
	Broadcast(t def.TYPE, payload []byte)

	// Clean the runner data on the Worker.
	Clean(runner uint64) error

	// Destroy the Worker.
	Destroy()
}
