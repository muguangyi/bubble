// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cron

// ITrigger interface.
type ITrigger interface {
	// Id returns the Trigger id.
	Id() uint64

	// Type returns the Trigger type.
	Type() Type

	// Job returns the inner ICronJob instance.
	Job() ICronJob

	// Start the trigger.
	Start()

	// Stop the trigger.
	Stop()
}
