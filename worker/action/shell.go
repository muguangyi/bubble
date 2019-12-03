// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// `shell` Action could handle shell commands.
//
// ```yaml
// -
//  action: shell
//  script:
//   - mkdir ...
//   - echo ...
// ```

package action

import (
	"bubble/env"
	"errors"
	"os/exec"
)

// ShellFactory struct.
type ShellFactory struct {
}

// Validate do nothing.
func (f *ShellFactory) Validate(conf env.IAny) error {
	return nil
}

// Create shell action.
func (f *ShellFactory) Create() IAction {
	return &shell{}
}

// -- Action --

type shell struct {
	Action
	cmd *exec.Cmd
}

func (s *shell) Execute(script env.IAny, target string, env env.IEnv, log ILog) chan bool {
	success := make(chan bool, 1)

	if !script.IsArr() {
		s.error = errors.New("shell command format is incorrect")
		success <- false
	} else {
		arr := script.Array()
		for _, code := range arr {
			v := env.Format(code)
			log.Debugf("-- Shell [%s]\n", v)

			s.cmd = exec.Command("sh", "-c", v)
			s.cmd.Dir = s.Cwd()
			s.cmd.Stdout = log.Std()
			s.cmd.Stderr = log.Std()
			err := s.cmd.Run()
			if err != nil {
				s.error = err
				break
			}
		}

		success <- (s.error == nil)
	}

	return success
}

func (s *shell) Cancel() error {
	if s.cmd != nil {
		return s.cmd.Process.Kill()
	}

	return nil
}
