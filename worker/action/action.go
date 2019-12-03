// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package action

import (
	"bubble/env"
	"os"
	"path"
	"path/filepath"
	"strconv"

	log "github.com/cihub/seelog"
)

// Action struct.
type Action struct {
	cwd   string
	error error
}

// Init is used for initializing the Action.
func (a *Action) Init(uid uint64, env env.IEnv) {
	a.error = nil

	ext, _ := os.Executable()
	a.cwd = path.Join(filepath.Dir(ext), "jobs", strconv.FormatUint(uid, 16))
	_, err := os.Stat(a.cwd)
	if err != nil && os.IsNotExist(err) {
		log.Debugf("Try to create dir: [%s].", a.cwd)
		err = os.MkdirAll(a.cwd, os.ModePerm)
		if err != nil {
			log.Errorf("Failed to create dir [%s] with err: [%s]!", a.cwd, err.Error())
		}
	}

	env.SetFunc("_SIZEOF", a.SizeOf)
}

// Cwd returns the current working directory.
func (a *Action) Cwd() string {
	return a.cwd
}

// Cancel the Action execution.
func (a *Action) Cancel() error {
	return nil
}

func (a *Action) Error() string {
	if a.error == nil {
		return ""
	}

	return a.error.Error()
}

// SizeOf calculate the target file size.
func (a *Action) SizeOf(args ...env.IAny) (env.IAny, error) {
	if len(args) == 0 {
		return env.NewAny(0), nil
	}

	file := args[0].ToString()
	log.Infof("Sizeof file [%s]\n", file)
	return env.NewAny(a.calcSize(file)), nil
}

func (a *Action) calcSize(p string) int64 {
	filePath := path.Join(a.Cwd(), p)
	f, err := os.Open(filePath)
	if err != nil {
		return 0
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil && os.IsNotExist(err) {
		// File or directory isn't exist.
		return 0
	}

	if stat.IsDir() {
		fs, err := f.Readdir(-1)
		if err != nil {
			return 0
		}

		var total int64
		for _, i := range fs {
			total += a.calcSize(p + "/" + i.Name())
		}

		return total
	}

	return stat.Size()
}
