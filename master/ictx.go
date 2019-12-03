// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

import (
	"bubble/def"
	"bubble/env"
)

// ICtx interface.
type ICtx interface {
	// ID return the ctx id.
	ID() uint64

	// LastWorker return the worker service ID where the prev action executed.
	LastWorker() uint64

	// Disk returns the code disk for share.
	Disk() string

	// Script return the code script to execute.
	Script() []byte

	// Variables return the code variables.
	Variables() []byte

	// Target return the target info.
	Target() string

	// Env return the current running IEnv.
	Env() env.IEnv

	// Notify the status with payload data.
	Notify(status def.STATUS, payload []byte)

	// SetResult to finish the ICtx execution.
	SetResult(result def.STATUS, env env.IEnv)
}
