// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package worker

import (
	log "github.com/cihub/seelog"
	"github.com/giant-tech/go-service/framework/iserver"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

const (
	// WORKFOLDER defines the working folder.
	WORKFOLDER string = ".prov"
	// WORKEFILE defines the compress target file name.
	WORKEFILE string = "streamdisk.zip"
	// CHUNKSIZE defines the file chunk size (30k).
	CHUNKSIZE int64 = 30 * 1024
)

type share struct {
	proxy iserver.IServiceProxy
	uid   uint64
	disk  string
}

func (s *share) Clean() {
	wp := s.workPath()
	_, err := os.Stat(wp)
	if err == nil {
		err = os.RemoveAll(wp)
		if err != nil {
			log.Errorf("Failed to clean dir [%s] with err: [%s]!", wp, err.Error())
		}
	}
}

func (s *share) cwd() string {
	ext, _ := os.Executable()
	return path.Join(filepath.Dir(ext), "jobs", strconv.FormatUint(s.uid, 16))
}

func (s *share) workPath() string {
	return path.Join(s.cwd(), WORKFOLDER)
}

func (s *share) workFilePath() string {
	return path.Join(s.workPath(), WORKEFILE)
}
