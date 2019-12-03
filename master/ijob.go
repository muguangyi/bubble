// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

import (
	"bubble/cron"
)

// IJob is the interface for Job.
type IJob interface {
	// ID returns Job unique id.
	ID() uint64

	// Name returns the Job name.
	Name() string

	// Trigger the Job.
	Trigger() error

	// Cancel the target Runner of the Job.
	Cancel(runner uint64) error

	// Script returns the Job script code.
	Script() ([]byte, error)

	// SetScript update bytes to the Job script code.
	SetScript(bytes []byte) error

	// Triggers returns all triggers of the Job.
	Triggers() ([]cron.ITrigger, error)

	// AddTrigger create a new trigger with interval type.
	AddTrigger(interval cron.Type) (cron.ITrigger, error)

	// RemoveTrigger delete target trigger by id.
	RemoveTrigger(id uint64) error

	// Destroy delete Job data.
	Destroy() error

	// Runners returns all IRunner of this Job.
	Runners() []IRunner

	// GetRunner return the target Runner by id.
	GetRunner(runner uint64) (IRunner, error)
}
