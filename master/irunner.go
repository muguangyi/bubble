// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

// IRunner presents an executation of a Job.
type IRunner interface {
	// ID returns Runner unique id.
	ID() uint64

	// Execute the Runner.
	Execute() error

	// Cancel the executation of the Runner.
	Cancel() error

	// Commands returns all ICommand of the Runner.
	Commands() []ICommand
}
