// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package worker

import (
	"bubble/env"
	"bubble/worker/action"
	"fmt"
	"sync"

	log "github.com/cihub/seelog"
)

// NewRunner method create a IRunner by name and worker.
func NewRunner(name string, worker IWorker) IRunner {
	return &runner{
		name:    name,
		worker:  worker,
		factory: create(name),
		procs:   make(map[uint64]action.IAction),
	}
}

type runner struct {
	name        string
	conf        env.IAny
	worker      IWorker
	factory     action.IFactory
	procsLocker sync.Mutex
	procs       map[uint64]action.IAction
}

// --- IRunner ---

// Name returns the runner name.
func (r *runner) Name() string {
	return r.name
}

func (r *runner) Conf() env.IAny {
	return r.conf
}

func (r *runner) Validate(conf env.IAny) error {
	err := r.factory.Validate(conf)
	if err != nil {
		return err
	}

	r.conf = conf
	return nil
}

func (r *runner) Execute(ctx ICtx) {
	a, err := r.queue(ctx)
	if err != nil {
		log.Error(err)
		r.worker.Finish(r.name, ctx.Master(), ctx.UID(), false, ctx.Env())
	} else {
		log.Infof("Execute proc [%d] in target [%s].\n", ctx.UID(), ctx.Target())

		// Initialize variables for Action scope.
		e := ctx.Env()
		vars := ctx.Variables()
		if !vars.IsNil() {
			m := vars.Map()
			for k, v := range m {
				e.Set(k, env.NewAny(e.Format(v)))
			}
		}

		// Execute the Action for the Context.
		success := <-a.Execute(ctx.Script(), ctx.Target(), e, &logger{runner: r, ctx: ctx})

		// Finish Action execution to Master.
		r.worker.Finish(r.name, ctx.Master(), ctx.UID(), success, e)

		r.procsLocker.Lock()
		defer r.procsLocker.Unlock()
		delete(r.procs, ctx.UID())
	}
}

func (r *runner) Cancel(uid uint64) {
	a, ok := r.procs[uid]
	if ok {
		a.Cancel()

		r.procsLocker.Lock()
		defer r.procsLocker.Unlock()
		delete(r.procs, uid)
	}
}

func (r *runner) Workload() int {
	r.procsLocker.Lock()
	defer r.procsLocker.Unlock()

	return len(r.procs)
}

// --- Inner ---

func (r *runner) queue(ctx ICtx) (action.IAction, error) {
	r.procsLocker.Lock()
	defer r.procsLocker.Unlock()

	_, ok := r.procs[ctx.UID()]
	if ok {
		return nil, fmt.Errorf("duplicated proc [%d]", ctx.UID())
	}

	a := r.factory.Create()
	r.procs[ctx.UID()] = a
	a.Init(ctx.UID(), ctx.Env())

	return a, nil
}

// --- Global ---

type maker func() action.IFactory

var (
	factory = make(map[string]maker)
)

func register(name string, m maker) {
	factory[name] = m
}

func create(name string) action.IFactory {
	m, ok := factory[name]
	if !ok {
		return nil
	}

	return m()
}

func init() {
	register("shell", func() action.IFactory { return &action.ShellFactory{} })
	register("unity", func() action.IFactory { return &action.UnityFactory{} })
	register("zip", func() action.IFactory { return &action.ZipFactory{} })
	register("ftp", func() action.IFactory { return &action.FtpFactory{} })
	register("email", func() action.IFactory { return &action.EmailFactory{} })
}
