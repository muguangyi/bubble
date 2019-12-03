// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package worker

import (
	"bubble/cron"
	"bubble/def"
	"bubble/env"
	"encoding/binary"
	"encoding/json"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"

	log "github.com/cihub/seelog"
	"github.com/giant-tech/go-service/framework/idata"
	"github.com/giant-tech/go-service/framework/iserver"
	"github.com/giant-tech/go-service/framework/service"
)

const (
	// WorkerConfigFile defines the worker configure file full name.
	WorkerConfigFile string = "./worker.yml"
	// CRONFILE defines the bubble cron job file name.
	CRONFILE string = ".bubble.crons"
)

// Worker type.
type Worker struct {
	service.BaseService
	mastersLocker sync.Mutex
	masters       map[uint64]iserver.IServiceProxy
	runners       map[string]IRunner
	providers     map[uint64]IProvider
	executors     map[uint64]IExecutor
	cron          cron.ICron
}

// OnInit method initialize the Worker.
func (w *Worker) OnInit() error {
	w.masters = make(map[uint64]iserver.IServiceProxy)
	w.runners = make(map[string]IRunner)
	w.providers = make(map[uint64]IProvider)
	w.executors = make(map[uint64]IExecutor)

	all, err := env.Load(WorkerConfigFile)
	if err != nil {
		return err
	}

	if !all.IsMap() {
		return errors.New("worker configure file format is incorrect")
	}

	data := all.Map()
	for k, cf := range data {
		runner := NewRunner(k, w)
		err = runner.Validate(cf)
		if err != nil {
			return err
		}

		w.runners[k] = runner
	}

	w.cron = cron.NewCron(func() cron.ICronJob { return &clean{w: w} }, path.Join(w.dir(), CRONFILE))
	w.cron.StartAll()

	return nil
}

// OnDestroy method.
func (w *Worker) OnDestroy() {
}

// OnConnected method.
func (w *Worker) OnConnected(info []*idata.ServiceInfo) {
	for _, i := range info {
		if i.Type == def.MasterService {
			proxy := iserver.GetServiceProxyMgr().GetServiceByID(i.ServiceID)
			supports := make(map[string][]byte)
			for k, r := range w.runners {
				supports[k], _ = r.Conf().ToBytes()
			}
			proxy.AsyncCall("Register", w.GetSID(), supports)

			w.mastersLocker.Lock()
			{
				w.masters[i.ServiceID] = proxy
			}
			w.mastersLocker.Unlock()
			log.Debugf("Worker [%d] register to Master.", w.GetSID())
		}
	}
}

// OnDisconnected method.
func (w *Worker) OnDisconnected(info []*idata.ServiceInfo) {
	for _, i := range info {
		if i.Type == def.MasterService {
			w.mastersLocker.Lock()
			{
				delete(w.masters, i.ServiceID)
			}
			w.mastersLocker.Unlock()
		}
	}
}

// OnTick method.
func (w *Worker) OnTick() {
	if len(w.masters) == 0 {
		return
	}

	workload := 0
	for _, r := range w.runners {
		workload += r.Workload()
	}
	payload := make([]byte, 8)
	binary.BigEndian.PutUint32(payload, uint32(workload))

	w.Broadcast(def.WORKLOAD, payload)
}

// --- RPC ---

// RPCExecute will trigger target action with master id, last provider, uid, script, variables, target and env parameters.
func (w *Worker) RPCExecute(action string, master, provider, uid uint64, disk string, script, variables []byte, target string, envData []byte) {
	log.Debugf("Trigger Action [%s] execution in target [%s] of Instance [%d].\n", action, target, uid)

	e := env.NewEnv()
	e.FromBytes(envData)

	s := env.NewAny(nil)
	s.FromBytes(script)

	vars := env.NewAny(nil)
	vars.FromBytes(variables)

	r, ok := w.runners[action]
	if !ok {
		log.Errorf("There is no action [%] in this Worker!\n", action)
		w.Finish(action, master, uid, false, e)
		return
	}

	// Reset disk parameter if last worker is the current worker.
	if w.GetSID() == provider {
		disk = ""
	}

	proxy := iserver.GetServiceProxyMgr().GetServiceByID(provider)
	executor := NewExecutor(proxy, w, uid, disk, r, NewCtx(master, uid, s, vars, target, e))
	w.executors[uid] = executor
	go executor.Execute()
}

// RPCCancel will cancel target action for Instance.
func (w *Worker) RPCCancel(action string, uid uint64) {
	log.Debugf("Cancel Action [%s] for Instance [%d].\n", action, uid)

	ector, ok := w.executors[uid]
	if ok {
		ector.Cancel()
	}
}

// RPCClean will clean local data for target uid task.
func (w *Worker) RPCClean(uid uint64) {
	ts, _ := w.cron.Triggers()
	for _, t := range ts {
		// If there is a same trigger already, ignore the task.
		if t.Job().(*clean).target == uid {
			return
		}
	}

	// Create a new trigger.
	t, err := w.cron.Add(cron.HOURLY)
	if err != nil {
		return
	}

	t.Job().(*clean).target = uid
	t.Start()
	w.cron.Flush()
}

// RPCBeforeSend handle the pre transfer disk request.
func (w *Worker) RPCBeforeSend(worker, uid uint64, disk string) {
	log.Debugf("RPCBeforeSend to worker [%d], uid [%d] and disk [%s].\n", worker, uid, disk)

	_, ok := w.providers[uid]
	if ok {
		return
	}

	proxy := iserver.GetServiceProxyMgr().GetServiceByID(worker)
	provider := NewProvider(proxy, uid, disk)
	w.providers[uid] = provider
	go provider.BeforeSend()
}

// RPCBeforeReceive handle the pre transfer disk response.
func (w *Worker) RPCBeforeReceive(uid uint64, length, chunks int64, checksum string) {
	log.Debugf("RPCBeforeReceive to uid [%d], length [%d], chunks [%d] and checksum [%s].\n", uid, length, chunks, checksum)

	executor, ok := w.executors[uid]
	if ok {
		executor.BeforeReceive(length, chunks, checksum)
	}
}

// RPCSend handle the transfer disk request.
func (w *Worker) RPCSend(uid uint64) {
	log.Debugf("RPCSend to uid [%d].\n", uid)

	provider, ok := w.providers[uid]
	if !ok {
		return
	}

	go provider.Send()
}

// RPCReceive handle the transfer disk progress.
func (w *Worker) RPCReceive(uid uint64, index int64, data []byte) {
	log.Debugf("RPCReceive to uid [%d], index [%d] and data length [%d].\n", uid, index, len(data))

	executor, ok := w.executors[uid]
	if ok {
		go executor.Receive(index, data)
	}
}

// RPCAfterReceive handle the post transfer disk response.
func (w *Worker) RPCAfterReceive(uid uint64) {
	log.Debugf("RPCAfterReceive to uid [%d].\n", uid)

	executor, ok := w.executors[uid]
	if ok {
		go executor.AfterReceive()
		delete(w.executors, uid)
	}
}

// --- IWorker ---

// UID returns the worker unique id.
func (w *Worker) UID() uint64 {
	return w.GetSID()
}

// Finish method notify the Master to finish the target action with payload data.
func (w *Worker) Finish(action string, master, uid uint64, success bool, env env.IEnv) {
	proxy, ok := w.masters[master]
	if !ok {
		log.Errorf("There is no Master [%d] to finish!", master)
		return
	}

	log.Debugf("Finish Instance [%d] with result [%t] to Master [%d].\n", uid, success, master)
	ebytes, err := env.ToBytes()
	if err != nil {
		log.Error(err)
		return
	}

	proxy.AsyncCall("OnFinish", w.GetSID(), action, uid, success, ebytes)
}

// Progress method notify the Master the target action progress.
func (w *Worker) Progress(action string, master, uid uint64, payload []byte) {
	proxy, ok := w.masters[master]
	if !ok {
		log.Errorf("There is no Master [%d] to notify!", master)
		return
	}

	log.Debugf("Progress Instance [%d] to Master [%d].\n", uid, master)
	proxy.AsyncCall("OnProgress", w.GetSID(), action, uid, payload)
}

// Broadcast method broadcast data to all Masters.
func (w *Worker) Broadcast(t def.TYPE, payload []byte) {
	w.mastersLocker.Lock()
	defer w.mastersLocker.Unlock()

	for _, proxy := range w.masters {
		proxy.AsyncCall("OnBroadcast", w.GetSID(), t, payload)
	}
}

// --- Inner ---

func (w *Worker) dir() string {
	ext, _ := os.Executable()
	return filepath.Dir(ext)
}

// --- ICronJob ---

type clean struct {
	w      *Worker
	target uint64
}

func (c *clean) Repeat() bool {
	return false
}

func (c *clean) Execute() {
	dir := path.Join(c.w.dir(), "jobs", strconv.FormatUint(c.target, 16))
	_, err := os.Stat(dir)
	if err == nil {
		log.Debugf("Try to clean dir: [%s].", dir)
		err = os.RemoveAll(dir)
		if err != nil {
			log.Errorf("Failed to clean dir [%s] with err: [%s]!", dir, err.Error())
		}
	}
}

func (c *clean) FromBytes(bytes []byte) {
	json.Unmarshal(bytes, &c.target)
}

func (c *clean) ToBytes() []byte {
	bytes, _ := json.Marshal(c.target)
	return bytes
}
