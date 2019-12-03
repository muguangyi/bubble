// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

import (
	"bubble/def"
	"bubble/env"
	"encoding/binary"
	"fmt"

	log "github.com/cihub/seelog"
	"github.com/giant-tech/go-service/framework/iserver"
)

// NewWorker method create an IWorker by master id and worker id.
func NewWorker(master, id uint64) IWorker {
	proxy := iserver.GetServiceProxyMgr().GetServiceByID(id)
	if proxy == nil {
		return nil
	}

	return &worker{
		master:  master,
		proxy:   proxy,
		actions: make(map[string]IAction),
	}
}

type worker struct {
	master   uint64
	proxy    iserver.IServiceProxy
	actions  map[string]IAction
	workload int
}

// --- IWorker ---

func (w *worker) ID() uint64 {
	return w.proxy.GetSID()
}

func (w *worker) Bind(supports map[string][]byte) error {
	for k, bytes := range supports {
		a, err := NewAction(w, k, bytes)
		if err != nil {
			return err
		}

		w.actions[k] = a
	}

	return nil
}

func (w *worker) Satisfy(command ICommand) bool {
	action, ok := w.actions[command.Name()]
	if !ok {
		return false
	}

	if command.Target() != "" {
		contains := false
		for _, t := range action.Target() {
			if t.String() == command.Target() {
				contains = true
				break
			}
		}

		if !contains {
			return false
		}
	}

	if command.Prefer() != "" {
		contains := false
		for _, p := range action.Prefer() {
			if p.String() == command.Prefer() {
				contains = true
				break
			}
		}

		if !contains {
			return false
		}
	}

	return true
}

func (w *worker) Workload() int {
	return w.workload
}

func (w *worker) Get(name string) IAction {
	a, ok := w.actions[name]
	if !ok {
		return nil
	}

	return a
}

func (w *worker) Finish(action string, runner uint64, success bool, env env.IEnv) error {
	a, ok := w.actions[action]
	if !ok {
		log.Errorf("There is no target Action [%s]!", action)
		return fmt.Errorf("there is no target Action [%s]", action)
	}

	return a.Finish(runner, success, env)
}

func (w *worker) Progress(action string, runner uint64, payload []byte) error {
	a, ok := w.actions[action]
	if !ok {
		log.Errorf("There is no target Action [%s]!", action)
		return fmt.Errorf("there is no target Action [%s]", action)
	}

	return a.Progress(runner, payload)
}

func (w *worker) Broadcast(t def.TYPE, payload []byte) {
	switch t {
	case def.WORKLOAD:
		w.workload = int(binary.BigEndian.Uint32(payload))
	}
}

func (w *worker) Clean(runner uint64) error {
	return w.proxy.AsyncCall("Clean", runner)
}

func (w *worker) Destroy() {
	for _, a := range w.actions {
		a.Destroy()
	}
	w.actions = nil
}

// --- Inner ---

func (w *worker) Execute(action string, runner, lastWorker uint64, disk string, script, variables []byte, target string, env env.IEnv) error {
	envData, err := env.ToBytes()
	if err != nil {
		return err
	}

	return w.proxy.AsyncCall("Execute", action, w.master, lastWorker, runner, disk, script, variables, target, envData)
}

func (w *worker) Cancel(action string, runner uint64) error {
	return w.proxy.AsyncCall("Cancel", action, runner)
}
