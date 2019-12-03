// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package action

import (
	"bubble/env"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	goftp "github.com/jlaffaye/ftp"
)

// FtpFactory struct.
type FtpFactory struct {
	address  string
	username string
	password string
}

// Validate whether the configure is correct.
func (f *FtpFactory) Validate(conf env.IAny) error {
	if !conf.IsMap() {
		return errors.New("ftp configure format is incorrect")
	}

	m := conf.Map()

	a, ok := m["address"]
	if !ok {
		return errors.New("there is no \"address\" in ftp configuration")
	}
	f.address = a.String()

	u, ok := m["username"]
	if !ok {
		return errors.New("there is no \"username\" in ftp configuration")
	}
	f.username = u.String()

	p, ok := m["password"]
	if !ok {
		return errors.New("there is no \"password\" in ftp configuration")
	}
	f.password = p.String()

	conn, err := goftp.Dial(f.address, goftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}
	defer conn.Quit()

	err = conn.Login(f.username, f.password)
	if err != nil {
		return err
	}
	conn.Logout()

	return nil
}

// Create ftp action.
func (f *FtpFactory) Create() IAction {
	return &ftp{address: f.address, username: f.username, password: f.password}
}

// --- Action ---

type ftp struct {
	Action
	address  string
	username string
	password string
}

func (f *ftp) Execute(script env.IAny, target string, env env.IEnv, log ILog) chan bool {
	success := make(chan bool, 1)

	if !script.IsArr() {
		f.error = errors.New("ftp command format is incorrect")
		success <- false
		return success
	}

	conn, err := goftp.Dial(f.address, goftp.DialWithTimeout(5*time.Second))
	if err != nil {
		f.error = err
		success <- false
		return success
	}
	defer conn.Quit()

	err = conn.Login(f.username, f.password)
	if err != nil {
		f.error = err
		success <- false
		return success
	}
	defer conn.Logout()

	arr := script.Array()
	for _, code := range arr {
		log.Debugf("-- ftp [%s]\n", code.String())
		args := strings.Split(env.Format(code), " ")
		if len(args) > 1 {
			if err = f.upload(conn, args[0], args[1:]...); err != nil {
				f.error = err
				success <- false
				return success
			}
		}
	}

	f.error = nil
	success <- true
	return success
}

func (f *ftp) upload(conn *goftp.ServerConn, rootPath string, files ...string) error {
	for _, v := range files {
		file, err := os.Open(path.Join(f.Cwd(), v))
		if err != nil {
			return err
		}

		stat, err := file.Stat()
		if err != nil {
			return err
		}

		if stat.IsDir() {
			fs, err := file.Readdir(-1)
			if err != nil {
				return err
			}

			for _, i := range fs {
				err = f.upload(conn, rootPath, v+"/"+i.Name())
				if err != nil {
					return err
				}
			}
		} else {
			p := path.Join(rootPath, v)
			dir := filepath.Dir(p)
			err = conn.ChangeDir(dir)
			if err != nil {
				if err = conn.MakeDir(dir); err != nil {
					return err
				}
			}
			conn.ChangeDir("/")
			if err = conn.Stor(p, file); err != nil {
				return err
			}
		}
	}

	return nil
}
