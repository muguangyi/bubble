// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

import (
	gobytes "bytes"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

// newPayloader method create a new payloader by path.
func newPayloader(path string) *payloader {
	p := &payloader{path: path, flag: UNFLUSHED}

	_, err := os.Stat(path)
	if err == nil {
		p.flag = FLUSHED
	}

	return p
}

const (
	// UNFLUSHED state.
	UNFLUSHED byte = 0
	// FLUSHING state.
	FLUSHING byte = 1
	// FLUSHED state.
	FLUSHED byte = 2
	// SEGLENGTH state.
	SEGLENGTH int = 10 * 1024
)

type payloader struct {
	path   string
	buf    gobytes.Buffer
	locker sync.Mutex
	flag   byte
}

func (p *payloader) Write(bytes []byte) (int, error) {
	if bytes == nil || len(bytes) == 0 {
		return 0, nil
	}

	p.locker.Lock()
	defer p.locker.Unlock()

	return p.buf.Write(bytes)
}

func (p *payloader) Bytes(full bool) ([]byte, bool, error) {
	if p.flag == FLUSHED {
		p.flag = FLUSHING
		go func() {
			bytes, err := ioutil.ReadFile(p.path)
			if err != nil {
				p.flag = FLUSHED
				return
			}

			p.Write(bytes)
			p.flag = UNFLUSHED
		}()
	}

	if p.flag != UNFLUSHED {
		return nil, full, errors.New("payload is not loaded")
	}

	all := p.buf.Bytes()
	if full || len(all) <= SEGLENGTH {
		return all, true, nil
	}

	part := all[len(all)-SEGLENGTH:]
	i := gobytes.IndexByte(part, '\n')
	if i < 0 || i >= len(part)-1 {
		return part, false, nil
	}

	return part[i+1:], false, nil
}

func (p *payloader) Flush() {
	if p.flag == UNFLUSHED {
		p.flag = FLUSHING

		go func() {
			p.locker.Lock()
			defer p.locker.Unlock()

			ioutil.WriteFile(p.path, p.buf.Bytes(), os.ModePerm)
			p.buf.Reset()
			p.flag = FLUSHED
		}()
	}
}
