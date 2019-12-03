// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package worker

import (
	"bubble/env"
)

// IRunner is the interface for an Action runner.
type IRunner interface {
	// Name returns the runner name.
	Name() string

	// Conf returns the runner conf data.
	Conf() env.IAny

	// Validate the target Action configuration.
	Validate(conf env.IAny) error

	// Execute the target Action with scripts.
	Execute(ctx ICtx)

	// Cancel the target job
	Cancel(uid uint64)

	// Workload returns the Runner workload.
	Workload() int
}
