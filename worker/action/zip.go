// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// `zip` Action could compress target files/directory into target zip file.
//
// ```yaml
// -
//  action: zip
//  script:
//   - ./a.zip ./any/folder ./any/file
// ```

package action

import (
	zipper "archive/zip"
	"bubble/env"
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

// ZipFactory struct.
type ZipFactory struct {
}

// Validate do nothing.
func (f *ZipFactory) Validate(conf env.IAny) error {
	return nil
}

// Create zip action.
func (f *ZipFactory) Create() IAction {
	return &zip{}
}

// --- Action ---

type zip struct {
	Action
}

func (z *zip) Execute(script env.IAny, target string, env env.IEnv, log ILog) chan bool {
	success := make(chan bool, 1)

	if !script.IsArr() {
		z.error = errors.New("zip command format is incorrect")
		success <- false
		return success
	}

	arr := script.Array()
	for _, code := range arr {
		log.Debugf("-- zip [%s]\n", code.ToString())
		args := strings.Split(env.Format(code), " ")
		if len(args) > 1 {
			if err := z.compress(args[0], args[1:]); err != nil {
				z.error = err
				success <- false
				return success
			}
		}
	}

	z.error = nil
	success <- true
	return success
}

func (z *zip) compress(target string, files []string) error {
	filePath := path.Join(z.Cwd(), target)
	if _, err := os.Stat(filePath); err == nil {
		// Delete file if exist.
		err = os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	zipFile, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer zipFile.Close()

	writer := zipper.NewWriter(zipFile)
	defer writer.Close()

	for _, file := range files {
		z.addFileToZip(writer, file)
	}

	return nil
}

func (z *zip) addFileToZip(writer *zipper.Writer, file string) error {
	filePath := path.Join(z.Cwd(), file)
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil && os.IsNotExist(err) {
		// File or directory isn't exist.
		return err
	}

	if stat.IsDir() {
		fs, err := f.Readdir(-1)
		if err != nil {
			return err
		}

		for _, i := range fs {
			err = z.addFileToZip(writer, file+"/"+i.Name())
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zipper.FileInfoHeader(stat)
		if err != nil {
			return err
		}

		header.Name = file
		header.Method = zipper.Store

		w, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(w, f)
		if err != nil {
			return err
		}
	}

	return nil
}
