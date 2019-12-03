// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

import (
	"bubble/def"
	"bubble/env"
	"fmt"
	"sync"

	log "github.com/cihub/seelog"
)

// NewAction method create a new IAction by worker, name and script data.
func NewAction(worker *worker, name string, bytes []byte) (IAction, error) {
	c := env.NewAny(nil)
	err := c.FromBytes(bytes)
	if err != nil {
		return nil, err
	}

	a := &action{worker: worker, name: name}
	a.procs = make(map[uint64]ICtx)

	if c.IsMap() {
		m := c.Map()

		t, ok := m["target"]
		if ok {
			a.target = t.Array()
		}

		p, ok := m["prefer"]
		if ok {
			a.prefer = p.Array()
		}
	}

	return a, nil
}

type action struct {
	worker      *worker
	name        string
	target      []env.IAny
	prefer      []env.IAny
	procsLocker sync.Mutex
	procs       map[uint64]ICtx
}

func (a *action) Target() []env.IAny {
	return a.target
}

func (a *action) Prefer() []env.IAny {
	return a.prefer
}

func (a *action) Execute(ctx ICtx) {
	err := a.worker.Execute(a.name, ctx.ID(), ctx.LastWorker(), ctx.Disk(), ctx.Script(), ctx.Variables(), ctx.Target(), ctx.Env())
	if err != nil {
		ctx.SetResult(def.FAILURE, ctx.Env())
		return
	}

	a.procsLocker.Lock()
	defer a.procsLocker.Unlock()
	a.procs[ctx.ID()] = ctx
}

func (a *action) Cancel(runner uint64) error {
	err := a.worker.Cancel(a.name, runner)
	if err != nil {
		return err
	}

	proc, ok := a.procs[runner]
	if !ok {
		return fmt.Errorf("runner [%d] is not exist", runner)
	}

	proc.SetResult(def.CANCEL, proc.Env())

	a.procsLocker.Lock()
	defer a.procsLocker.Unlock()
	delete(a.procs, runner)

	return nil
}

func (a *action) Finish(runner uint64, success bool, env env.IEnv) error {
	status := def.SUCCESS
	if !success {
		status = def.FAILURE
	}

	log.Debugf("Action [%s] receive notify for Runner [%d] with Status [%d].\n", a.name, runner, status)

	proc, ok := a.procs[runner]
	if !ok {
		return fmt.Errorf("runner [%d] is not exist", runner)
	}

	proc.SetResult(status, env)

	a.procsLocker.Lock()
	defer a.procsLocker.Unlock()
	delete(a.procs, runner)

	return nil
}

func (a *action) Progress(runner uint64, payload []byte) error {
	log.Debugf("Action [%s] receive progress for Runner [%d].\n", a.name, runner)

	proc, ok := a.procs[runner]
	if !ok {
		return fmt.Errorf("runner [%d] is not exist", runner)
	}

	proc.Notify(def.ONGOING, payload)

	return nil
}

func (a *action) Destroy() {
	a.procsLocker.Lock()
	defer a.procsLocker.Unlock()

	for _, p := range a.procs {
		p.SetResult(def.INTERRUPT, p.Env())
	}
	a.procs = nil
}
