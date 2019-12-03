// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

import (
	"bubble/env"
	"errors"
	"fmt"
)

// Parse Job scripts to command sequence.
func Parse(runner *runner, bytes []byte) ([]ICommand, error) {
	script := env.NewAny(nil)
	err := script.FromBytes(bytes)
	if err != nil {
		return nil, err
	}

	if !script.IsArr() {
		return nil, errors.New("job script is not an array")
	}

	arr := script.Array()
	cmds := make([]ICommand, len(arr))
	for i, c := range arr {
		if !c.IsMap() {
			return nil, fmt.Errorf("command [%d] format is incorrect", i)
		}

		cmd := NewCommand(runner, i).(*command)
		cmds[i] = cmd

		detail := c.Map()
		for k, v := range detail {
			switch k {
			case "action":
				cmd.name = v.String()
			case "alias":
				cmd.alias = v.String()
			case "disk":
				cmd.disk = v.String()
			case "script":
				cmd.script = v
			case "variables":
				cmd.variables = v
			case "when":
				cmd.when = v.String()
			case "where":
				{
					if v.IsNil() {
						cmd.where = 0
					} else {
						cmd.where = v.Int()
					}
				}
			case "target":
				cmd.target = v.String()
			case "prefer":
				cmd.prefer = v.String()
			}
		}

		where := cmd.Where()
		if where == -1 || i == 0 {
			// Anywhere.
			cmd.group = &group{cmds: make([]ICommand, 0)}
		} else if where == 0 {
			// The previous item.
			cmd.group = cmds[i-1].(*command).group
		} else {
			// The target index item.
			cmd.group = cmds[where-1].(*command).group
		}
		cmd.group.cmds = append(cmd.group.cmds, cmd)
	}

	return cmds, nil
}
