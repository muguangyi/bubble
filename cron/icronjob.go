// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cron

// ICronJob interface.
type ICronJob interface {
	// Repeat returns if the Job need to repeat.
	Repeat() bool

	// Execute the job logic.
	Execute()

	// FromBytes imports []byte into the Job.
	FromBytes(bytes []byte)

	// ToBytes returns the Job data to []byte.
	ToBytes() []byte
}
