// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

import (
	"bubble/def"
	"bubble/env"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/giant-tech/go-service/framework/idata"
	"github.com/giant-tech/go-service/framework/service"
)

const (
	// MasterConfigFile defines the configure file full path.
	MasterConfigFile string = "./master.yml"
)

// Master type.
type Master struct {
	service.BaseService
	workers map[uint64]IWorker
	jobs    map[string]IJob
	web     IWeb
}

// OnInit method.
func (m *Master) OnInit() error {
	m.workers = make(map[uint64]IWorker)
	m.loadJobs()

	// Load configure file.
	conf, err := env.Load(MasterConfigFile)
	if err != nil {
		return err
	}

	if !conf.IsMap() {
		return errors.New("master configure file format is incorrect")
	}

	all := conf.Map()

	conf, ok := all["web"]
	if !ok {
		return errors.New("not setting \"web\" for Master configure")
	}

	// Run http service.
	m.web = NewWeb(m, conf)
	return m.web.Serve()
}

// OnDestroy method.
func (m *Master) OnDestroy() {
	m.web.Close()
}

// OnConnected method.
func (m *Master) OnConnected(info []*idata.ServiceInfo) {
	for _, i := range info {
		if i.Type == def.WorkerService {
			worker := NewWorker(m.GetSID(), i.ServiceID)
			if worker != nil {
				// TODO: Add lock
				m.workers[i.ServiceID] = worker
			}
		}
	}
}

// OnDisconnected method.
func (m *Master) OnDisconnected(info []*idata.ServiceInfo) {
	for _, i := range info {
		if i.Type == def.WorkerService {
			worker, ok := m.workers[i.ServiceID]
			if ok {
				worker.Destroy()

				// TODO: Add lock
				delete(m.workers, i.ServiceID)
			}
		}
	}
}

// OnTick method.
func (m *Master) OnTick() {
}

// --- RPC ---

// RPCRegister register a Worker to the Master.
func (m *Master) RPCRegister(id uint64, supports map[string][]byte) {
	log.Infof("Master receive Worker [%d] register.", id)
	worker, ok := m.workers[id]
	if !ok {
		return
	}

	err := worker.Bind(supports)
	if err != nil {
		log.Error(err)
	}
}

// RPCOnFinish receive the finish status from Worker.
func (m *Master) RPCOnFinish(worker uint64, action string, runner uint64, success bool, envData []byte) {
	w, ok := m.workers[worker]
	if !ok {
		// TODO: Log error
		return
	}

	e := env.NewEnv()
	if err := e.FromBytes(envData); err != nil {
		// TODO: Log error
		return
	}

	w.Finish(action, runner, success, e)
}

// RPCOnProgress receive the progress data from Worker.
func (m *Master) RPCOnProgress(worker uint64, action string, runner uint64, payload []byte) {
	w, ok := m.workers[worker]
	if !ok {
		// TODO: Log error
		return
	}

	w.Progress(action, runner, payload)
}

// RPCOnBroadcast receive data from Worker.
func (m *Master) RPCOnBroadcast(worker uint64, t def.TYPE, payload []byte) {
	w, ok := m.workers[worker]
	if !ok {
		// TODO: Log error
		return
	}

	w.Broadcast(t, payload)
}

// --- IMaster ---

// Create method.
func (m *Master) Create(job string) error {
	_, ok := m.jobs[job]
	if ok {
		return fmt.Errorf("can't create Job [%s] since it's already exist", job)
	}

	uid, err := def.NextUid()
	if err != nil {
		return err
	}

	j, err := NewJob(m, uid, job)
	if err != nil {
		return err
	}

	m.jobs[job] = j

	return nil
}

// Delete method.
func (m *Master) Delete(job string) error {
	j, ok := m.jobs[job]
	if !ok {
		return fmt.Errorf("can't delete Job [%s] since it's not exist", job)
	}

	j.Destroy()
	delete(m.jobs, job)

	return nil
}

// Get method.
func (m *Master) Get(job string) (IJob, error) {
	j, ok := m.jobs[job]
	if !ok {
		return nil, fmt.Errorf("can't find Job [%s] since it's already exist", job)
	}

	return j, nil
}

// List method.
func (m *Master) List() []IJob {
	jobs := make([]IJob, 0)
	for _, j := range m.jobs {
		jobs = append(jobs, j)
	}

	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].ID() < jobs[j].ID()
	})

	return jobs
}

// Select method.
func (m *Master) Select(cmds []ICommand) IWorker {
	var worker IWorker
	for _, w := range m.workers {
		satisfy := true
		for _, cmd := range cmds {
			if !w.Satisfy(cmd) {
				satisfy = false
				break
			}
		}

		if satisfy && (worker == nil || w.Workload() < worker.Workload()) {
			worker = w
		}
	}

	return worker
}

// Workers method.
func (m *Master) Workers() []IWorker {
	workers := make([]IWorker, 0)
	for _, w := range m.workers {
		workers = append(workers, w)
	}

	sort.Slice(workers, func(i, j int) bool {
		return workers[i].ID() < workers[j].ID()
	})

	return workers
}

// --- Inner ---

func (m *Master) loadJobs() error {
	m.jobs = make(map[string]IJob)

	ext, _ := os.Executable()
	dir, err := os.Open(path.Join(filepath.Dir(ext), "jobs"))
	stat, err := dir.Stat()
	if err != nil {
		return err
	}

	if stat.IsDir() {
		fs, err := dir.Readdir(-1)
		if err != nil {
			return err
		}

		for _, f := range fs {
			if f.IsDir() {
				parts := strings.Split(f.Name(), "@")
				name := parts[0]
				uid, _ := strconv.ParseUint(parts[1], 16, 64)
				j, err := NewJob(m, uid, name)
				if err != nil {
					return err
				}

				m.jobs[j.Name()] = j
			}
		}
	}

	return nil
}
