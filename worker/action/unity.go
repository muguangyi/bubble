// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// `unity` Action could handle unity engine command line process.
//
// ```yaml
// -
//  action: unity
//  script:
//   - -projectPath ... -batchmode -executeMethod ...
//  target: v20184
//  prefer: android
//  where:
// ```

package action

import (
	"bubble/env"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/hpcloud/tail"
)

// UnityFactory struct.
type UnityFactory struct {
	target  []env.IAny
	prefer  []env.IAny
	version map[string]env.IAny
}

// Validate whether the configure is correct.
func (f *UnityFactory) Validate(conf env.IAny) error {
	if !conf.IsMap() {
		return errors.New("unity configure format is incorrect")
	}

	m := conf.Map()

	t, ok := m["target"]
	if !ok {
		return errors.New("not setting \"target\" for unity")
	}
	f.target = t.Array()

	p, ok := m["prefer"]
	if !ok {
		return errors.New("not setting \"prefer\" for unity")
	}
	f.prefer = p.Array()

	vs, ok := m["version"]
	if !ok {
		return errors.New("not setting \"version\" for unity")
	}

	f.version = vs.Map()
	for _, v := range f.version {
		path := v.String()
		if _, err := os.Stat(path); err != nil {
			return err
		}
		f.version[DEFAULT_VERSION] = v
	}

	return nil
}

// Create unity action.
func (f *UnityFactory) Create() IAction {
	return &unity{version: f.version}
}

// -- Action --

type unity struct {
	Action
	version map[string]env.IAny
	cmd     *exec.Cmd
}

func (s *unity) Execute(script env.IAny, target string, env env.IEnv, log ILog) chan bool {
	success := make(chan bool, 1)

	if target == "" {
		target = DEFAULT_VERSION
	}

	u, ok := s.version[target]
	if !ok {
		s.error = fmt.Errorf("there is no target [%s] unity", target)
		success <- false
		return success
	}

	if !script.IsArr() {
		s.error = errors.New("unity script format is incorrect")
		success <- false
		return success
	}

	arr := script.Array()
	for _, code := range arr {
		log.Debugf("-- unity [%s]\n", code.ToString())
		args := strings.Split(env.Format(code), " ")

		// Tail log file.
		s.tail(s.logFilePath(args...), log)

		s.cmd = exec.Command(u.String(), args...)
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
	return success
}

func (s *unity) Cancel() error {
	if s.cmd != nil {
		return s.cmd.Process.Kill()
	}

	return nil
}

func (s *unity) logFilePath(options ...string) string {
	for i, v := range options {
		if strings.ToLower(v) == "-logfile" {
			return path.Join(s.Cwd(), options[i+1])
		}
	}

	if runtime.GOOS == "windows" {
		return path.Join(os.Getenv("LOCALAPPDATA"), "Unity/Editor/Editor.log")
	}

	return path.Join(os.Getenv("HOME"), "Library/Logs/Unity/Editor.log")
}

func (s *unity) tail(file string, log ILog) {
	log.Infof("Tail log file: %s\n", file)
	go func() {
		t, err := tail.TailFile(file, tail.Config{
			ReOpen:    true,
			Follow:    true,
			Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
			MustExist: false,
			Poll:      true,
		})
		if err != nil {
			return
		}

		defer t.Stop()

		for {
			l, ok := <-t.Lines
			if !ok {
				continue
			}
			log.Infof("%s\n", l.Text)
		}
	}()
}

const (
	// DEFAULT_VERSION defines the default version to use.
	DEFAULT_VERSION string = "v?"
)
