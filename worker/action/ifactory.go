// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package action

import (
	"bubble/env"
)

// IFactory is an Action factory to validate the Action
// configuration and create the target Action instance.
type IFactory interface {
	// Validate the configuration env for the Action.
	Validate(conf env.IAny) error

	// Create the target Action instance.
	Create() IAction
}
