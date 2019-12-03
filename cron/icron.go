// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cron

type Type int

const (
	// QUARTERHOURLY defines 15 minutes.
	QUARTERHOURLY Type = 1
	// HOURLY defines 1 hour.
	HOURLY Type = 2
	// DAILY defines 1 day.
	DAILY Type = 3
	// WEEKLY defines 1 week.
	WEEKLY Type = 4
	// MONTHLY defines 1 month (30 days).
	MONTHLY Type = 5
)

// ICron interface.
type ICron interface {
	// StartAll will start all triggers of the Cron.
	StartAll()

	// Flush the stats into file.
	Flush()

	// Add a new Trigger by type.
	Add(t Type) (ITrigger, error)

	// Remove a target Trigger by id.
	Remove(id uint64) (ITrigger, error)

	// Triggers returns all Triggers in this Cron.
	Triggers() ([]ITrigger, error)

	// Destroy all Triggers.
	Destroy()
}
