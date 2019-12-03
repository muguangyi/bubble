// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

// IMaster interface.
type IMaster interface {
	// Create a Job by name.
	Create(job string) error

	// Delete a Job by name.
	Delete(job string) error

	// Get a Job by name.
	Get(job string) (IJob, error)

	// List all Jobs.
	List() []IJob

	// Filter a Worker which satisfy the commands requirements.
	Select(cmds []ICommand) IWorker

	// Workers returns all workers.
	Workers() []IWorker
}
