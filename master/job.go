// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

import (
	"bubble/cron"
	"bubble/def"
	"bubble/env"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"sync"

	log "github.com/cihub/seelog"
)

// NewJob method create an IJob by master, id and name.
func NewJob(master IMaster, id uint64, name string) (IJob, error) {
	j := &job{master: master, id: id, name: name, runners: make(map[uint64]IRunner)}
	if err := j.init(); err != nil {
		return nil, err
	}

	return j, nil
}

const (
	// CRONFILE defines the cron job file name.
	CRONFILE string = ".bubble.crons"
	// BUBBLEFILE defines the bubble script file name.
	BUBBLEFILE string = ".bubble.yml"
	// SCRIPT defines the bubble default script.
	SCRIPT string = `# .bubble.yml
-
 action: shell
 script:
 - echo Hello Bubble!
`
)

type job struct {
	master  IMaster
	id      uint64
	name    string
	script  env.IAny
	locker  sync.Mutex
	runners map[uint64]IRunner
	cron    cron.ICron
}

// --- IJob ---

func (j *job) ID() uint64 {
	return j.id
}

func (j *job) Name() string {
	return j.name
}

func (j *job) Trigger() error {
	j.locker.Lock()
	defer j.locker.Unlock()

	uid, err := def.NextUid()
	if err != nil {
		return err
	}

	r := NewRunner(uid, j)
	j.runners[r.ID()] = r
	return r.Execute()
}

func (j *job) Cancel(runner uint64) error {
	r, ok := j.runners[runner]
	if !ok {
		return fmt.Errorf("job [%s] Runner [%d] is not exist", j.name, runner)
	}

	return r.Cancel()
}

func (j *job) Script() ([]byte, error) {
	return j.script.ToBytes()
}

func (j *job) SetScript(bytes []byte) error {
	err := j.script.FromBytes(bytes)
	if err != nil {
		return err
	}

	go func() {
		fp := path.Join(j.Dir(), BUBBLEFILE)
		if err = ioutil.WriteFile(fp, bytes, os.ModePerm); err != nil {
			log.Error(err)
		}
	}()

	return nil
}

func (j *job) Triggers() ([]cron.ITrigger, error) {
	return j.cron.Triggers()
}

func (j *job) AddTrigger(interval cron.Type) (cron.ITrigger, error) {
	t, err := j.cron.Add(interval)
	if err == nil {
		t.Start()
		j.cron.Flush()
	}

	return t, err
}

func (j *job) RemoveTrigger(id uint64) error {
	t, err := j.cron.Remove(id)
	if err == nil {
		t.Stop()
		j.cron.Flush()
	}

	return err
}

func (j *job) Destroy() error {
	// Stop all crons.
	j.cron.Destroy()

	// Clean Job persistent data.
	dir := j.Dir()
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("can't destroy Job [%s] since it's not exist", j.name)
	}

	return os.RemoveAll(dir)
}

func (j *job) Runners() []IRunner {
	rs := make([]IRunner, 0)
	for _, r := range j.runners {
		rs = append(rs, r)
	}

	sort.Slice(rs, func(i, j int) bool {
		return rs[i].ID() > rs[j].ID()
	})

	return rs
}

func (j *job) GetRunner(runner uint64) (IRunner, error) {
	r, ok := j.runners[runner]
	if !ok {
		return nil, fmt.Errorf("runner [%d] is not exist", runner)
	}

	return r, nil
}

// --- ICronJob ---

func (j *job) Repeat() bool {
	return true
}

func (j *job) Execute() {
	j.Trigger()
}

func (j *job) FromBytes(bytes []byte) {

}

func (j *job) ToBytes() []byte {
	return nil
}

// --- Inner ---

func (j *job) Dir() string {
	ext, _ := os.Executable()
	return path.Join(filepath.Dir(ext), "jobs", j.name+"@"+strconv.FormatUint(j.id, 16))
}

func (j *job) init() error {
	// Make Job folder if it's not exist.
	dir := j.Dir()
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Load Job script. If it's not exist, create one with
	// default initial script code.
	var bytes []byte
	fp := path.Join(dir, BUBBLEFILE)
	_, err = os.Stat(fp)
	if err != nil && os.IsNotExist(err) {
		bytes = []byte(SCRIPT)
		if err = ioutil.WriteFile(fp, bytes, os.ModePerm); err != nil {
			return err
		}
	} else if bytes, err = ioutil.ReadFile(fp); err != nil {
		return err
	}

	j.script = env.NewAny(nil)
	err = j.script.FromBytes(bytes)
	if err != nil {
		return err
	}

	// Load all runners under the Job by walking all folders.
	f, _ := os.Open(dir)
	fs, err := f.Readdir(-1)
	if err != nil {
		return err
	}

	j.locker.Lock()
	defer j.locker.Unlock()

	for _, child := range fs {
		cp := path.Join(dir, child.Name())
		stat, err := os.Stat(cp)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			id, _ := strconv.ParseUint(child.Name(), 16, 64)
			r := NewRunner(id, j)
			j.runners[r.ID()] = r
		}
	}

	// Load all crons and start.
	j.cron = cron.NewCron(func() cron.ICronJob { return j }, path.Join(dir, CRONFILE))
	j.cron.StartAll()

	return nil
}
