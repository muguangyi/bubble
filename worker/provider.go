// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package worker

import (
	zipper "archive/zip"
	"bubble/util"
	log "github.com/cihub/seelog"
	"github.com/giant-tech/go-service/framework/iserver"
	"io"
	"math"
	"os"
	"path"
)

// NewProvider method create a new ITransfer by uid and disk.
func NewProvider(proxy iserver.IServiceProxy, uid uint64, disk string) IProvider {
	return &provider{share: share{proxy: proxy, uid: uid, disk: disk}}
}

type provider struct {
	share
}

func (p *provider) BeforeSend() {
	// Create working folder
	log.Debugf("Try to create [%s] working folder.", p.workPath())

	wp := p.workPath()
	_, err := os.Stat(wp)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(wp, os.ModePerm)
		if err != nil {
			log.Errorf("Failed to create dir [%s] with err: [%s]!", wp, err.Error())
		}
	}

	log.Debugf("Try to compress to [%s].", p.workFilePath())
	target := p.workFilePath()
	err = p.compress(target)
	if err != nil {
		log.Errorf("Compress to [%s] failed.", p.workFilePath())
		return
	}

	stat, err := os.Stat(target)
	if err != nil && os.IsNotExist(err) {
		log.Errorf("File [%s] isn't exist.", p.workFilePath())
		return
	}

	fileLength := stat.Size()
	chunks := int64(math.Ceil(float64(fileLength) / float64(CHUNKSIZE)))
	checksum := util.CalcFileChecksum(target)

	log.Debugf("Trigger BeforeReceive: %d, %d, %s.", fileLength, chunks, checksum)
	p.proxy.AsyncCall("BeforeReceive", p.uid, fileLength, chunks, checksum)
}

func (p *provider) Send() {
	log.Debugf("Start Send from uid [%d].\n", p.uid)

	f, err := os.Open(p.workFilePath())
	if err != nil {
		log.Errorf("File [%s] isn't exist.", p.workFilePath())
		return
	}

	stat, _ := f.Stat()
	fileLength := stat.Size()
	chunks := int(math.Ceil(float64(fileLength) / float64(CHUNKSIZE)))
	for i := 0; i < chunks; i++ {
		data := make([]byte, CHUNKSIZE)
		size, _ := f.ReadAt(data, CHUNKSIZE*int64(i))
		log.Debugf("Trigger Receive: uid [%d], index [%d], and data length [%d].\n", p.uid, i, size)
		p.proxy.AsyncCall("Receive", p.uid, int64(i), data[0:size])
	}
	f.Close()

	log.Debugf("Trigger AfterReceive: uid [%d].\n", p.uid)
	p.proxy.AsyncCall("AfterReceive", p.uid)

	// Clean the temp folder.
	p.Clean()
}

func (p *provider) compress(filePath string) error {
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

	return p.addFileToZip(writer, p.disk)
}

func (p *provider) addFileToZip(writer *zipper.Writer, file string) error {
	filePath := path.Join(p.cwd(), file)
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
			err = p.addFileToZip(writer, file+"/"+i.Name())
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
