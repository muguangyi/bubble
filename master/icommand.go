// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

import (
	"bubble/def"
	"bubble/env"
)

// WHEN redefines int8 as status type.
type WHEN int8

const (
	// ALWAYS defines the action could be triggered without any pre-condition.
	ALWAYS WHEN = 0
	// SUCCESS defines the action should be triggered when the pre action is success.
	SUCCESS WHEN = 1
	// FAILURE defines the action should be triggered when the pre actions is failed.
	FAILURE WHEN = -1
)

// ICommand is the interface for Command in Job.
type ICommand interface {
	// Index returns the Command index.
	Index() int

	// Name returns the Action type.
	Name() string

	// Alias returns the Action alias name.
	Alias() string

	// Disk returns the streaming disk path.
	Disk() string

	// Script returns the script object of the command.
	Script() env.IAny

	// Variables returns the Command variables.
	Variables() env.IAny

	// When returns the trigger condition.
	When() WHEN

	// Where the command should be executed. -1 indicates anywhere.
	Where() int

	// Target returns the sub info of the command.
	Target() string

	// Prefer returns the prefer ability of the Action.
	Prefer() string

	// Status returns the command status.
	Status() def.STATUS

	// Total cost in seconds.
	Measure() int64

	// Logs return all log data of the Command.
	Logs(full bool) ([]byte, bool, error)

	// Notify Command status.
	Notify(status def.STATUS, payload []byte) error
}
