// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"
	"path/filepath"
)

func NewPortal(root string, indexFile string) IWebControl {
	return &portal{root: root, indexFile: indexFile}
}

type portal struct {
	root      string
	indexFile string
}

func (p *portal) Init(handler IWebHandler) {
	handler.HandleStatic("/css/", p.handleFile)
	handler.HandleStatic("/js/", p.handleFile)
	handler.HandleStatic("/", p.handleFile)
}

func (p *portal) handleFile(w http.ResponseWriter, req *http.Request) {
	target := req.URL.Path[1:]
	if target == "" {
		target = p.indexFile
	}
	target = filepath.Join(p.root, target)
	http.ServeFile(w, req, target)
}
