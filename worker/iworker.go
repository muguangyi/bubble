// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package worker

import (
	"bubble/def"
	"bubble/env"
)

// IWorker is the interface to notify status to Master.
type IWorker interface {
	// UID returns the worker unique id.
	UID() uint64

	// Finish action to Master.
	Finish(action string, master, uid uint64, success bool, env env.IEnv)

	// Progress action info to Master.
	Progress(action string, master, uid uint64, payload []byte)

	// Broadcast data to all connected Masters.
	Broadcast(t def.TYPE, payload []byte)
}
