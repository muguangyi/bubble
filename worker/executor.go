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
	"os"
	"path"
	"path/filepath"
)

// NewExecutor create a new IExecutor with parameters.
func NewExecutor(proxy iserver.IServiceProxy, worker IWorker, uid uint64, disk string, runner IRunner, ctx ICtx) IExecutor {
	return &executor{share: share{proxy: proxy, uid: uid, disk: disk}, worker: worker, runner: runner, ctx: ctx}
}

type executor struct {
	share
	worker   IWorker
	length   int64
	checksum string
	f        *os.File
	chunks   []int64
	runner   IRunner
	ctx      ICtx
}

func (e *executor) Execute() {
	e.worker.Progress(e.runner.Name(), e.ctx.Master(), e.ctx.UID(), []byte(""))

	if e.disk == "" {
		e.runner.Execute(e.ctx)
	} else {
		e.proxy.AsyncCall("BeforeSend", e.worker.UID(), e.uid, e.disk)
	}
}

func (e *executor) Cancel() {
	e.runner.Cancel(e.uid)
}

func (e *executor) BeforeReceive(length, count int64, checksum string) {
	log.Debugf("Start to receive length [%d], count [%d] and checksum [%s].\n", length, count, checksum)

	wp := e.workPath()
	_, err := os.Stat(wp)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(wp, os.ModePerm)
		if err != nil {
			log.Errorf("Failed to create dir [%s] with err: [%s]!", wp, err.Error())
		}
	}

	e.length = length
	e.checksum = checksum
	e.f, _ = os.OpenFile(e.workFilePath(), os.O_CREATE|os.O_RDWR, os.ModePerm)
	e.chunks = make([]int64, count)
	for i := range e.chunks {
		e.chunks[i] = -1
	}

	log.Debugf("Trigger Send to start transfer for [%d].\n", e.uid)
	e.proxy.AsyncCall("Send", e.uid)
}

func (e *executor) Receive(index int64, data []byte) {
	log.Debugf("Receiving index [%d] and data size [%d].\n", index, len(data))

	size, err := e.f.WriteAt(data, int64(index*CHUNKSIZE))
	if err != nil || size != len(data) {
		log.Errorf("Receiving data index [%d] failed because of err or write size [%d] is not equals to [%d].\n", index, size, len(data))
		return
	}

	e.chunks[index] = int64(size)
}

func (e *executor) AfterReceive() {
	log.Debug("Wait until all chunks have been received.\n")

	// Wait until all chunks have been received.
	// TODO: add timeout support.
	for complete := true; complete == false; complete = true {
		for i := range e.chunks {
			if i == -1 {
				complete = false
				break
			}
		}
	}

	e.f.Close()

	log.Debug("Verify checksum is correct.\n")
	checksum := util.CalcFileChecksum(e.workFilePath())
	if checksum != e.checksum {
		// TODO: handle error.
		return
	}

	// Uncompress file
	log.Debugf("Unzip file [%s].\n", e.workFilePath())
	reader, err := zipper.OpenReader(e.workFilePath())
	if err != nil {
		// TODO: handle error.
	}

	for _, f := range reader.File {
		target := path.Join(e.cwd(), f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(target, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
			// TODO: handle error.
		}

		outFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			// TODO: handle error.
		}

		rc, err := f.Open()
		if err != nil {
			// TODO: handle error.
		}

		log.Debugf("Unzip inner file [%s].\n", target)
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
	}
	reader.Close()

	// Clean the temp folder.
	e.Clean()

	log.Debug("Start to execute command.\n")
	e.runner.Execute(e.ctx)
}
