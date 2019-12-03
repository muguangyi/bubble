// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// runner carry a Job script and create a sub-folder under Job.

package master

import (
	"bubble/def"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	log "github.com/cihub/seelog"
)

// NewRunner create a new IRunner by id and job.
func NewRunner(id uint64, job *job) IRunner {
	r := &runner{id: id, job: job}

	dir := r.Dir()
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil
		}
	}

	var bytes []byte
	scriptFilePath := r.ScriptPath()
	_, err = os.Stat(scriptFilePath)
	if err != nil && os.IsNotExist(err) {
		bytes, err = job.script.ToBytes()
		if err != nil {
			return nil
		}

		if err = ioutil.WriteFile(scriptFilePath, bytes, os.ModePerm); err != nil {
			return nil
		}
	} else if bytes, err = ioutil.ReadFile(scriptFilePath); err != nil {
		return nil
	}

	r.cmds, err = Parse(r, bytes)
	if err != nil {
		return nil
	}

	// Load status.
	statusFilePath := r.StatusPath()
	_, err = os.Stat(statusFilePath)
	if err == nil {
		bytes, err := ioutil.ReadFile(statusFilePath)
		if err != nil {
			log.Errorf("Read Job [%s] Runner [%s] status file failed!", r.job.name, r.id)
			return nil
		}

		var stats []*commandStat
		err = json.Unmarshal(bytes, &stats)
		if err != nil {
			log.Errorf("Unmarshal Job [%s] Runner [%s] status failed!", r.job.name, r.id)
			return nil
		}

		for i, s := range stats {
			cmd := r.cmds[i].(*command)
			cmd.status = s.Status
			cmd.beginStamp = s.BeginTime
			cmd.finishStamp = s.FinishTime

			// Make sure finishStamp is valid.
			if cmd.finishStamp == -1 {
				switch cmd.status {
				case def.SUCCESS, def.FAILURE, def.CANCEL, def.INTERRUPT:
					cmd.finishStamp = cmd.beginStamp
				}
			}
		}
	}

	return r
}

const (
	// STATUSFILE defines the file name.
	STATUSFILE string = ".bubble.stat"
)

type runner struct {
	id   uint64
	job  *job
	cmds []ICommand
}

func (r *runner) ID() uint64 {
	return r.id
}

func (r *runner) Execute() error {
	go func() {
		log.Infof("Job [%s] is executing.\n", r.job.Name())

		status := def.SUCCESS
		ctx := NewCtx(r).(*ctx)
		for i := 0; i < len(r.cmds) && status != def.INTERRUPT; i++ {
			cmd := r.cmds[i].(*command)
			ctx.Cmd = cmd
			if cmd.group.worker == nil {
				// Find proper Worker and wait 1 min for time out if can't find.
				var worker IWorker
				timeOut := 1 * time.Minute
				for {
					worker = r.job.master.Select(cmd.group.cmds)
					if worker != nil || timeOut <= 0 {
						break
					}
					time.Sleep(time.Second)
					timeOut -= time.Second
				}

				// No proper Worker.
				if worker == nil {
					log.Error("There is no suitable worker!\n")
					status = def.FAILURE
					cmd.Notify(def.FAILURE, nil)
					break
				}

				defer worker.Clean(r.id) // Clean the runner data on the candidated Worker.
				cmd.group.worker = worker
			}

			when := cmd.When()
			if status != def.CANCEL &&
				(when == ALWAYS || (when == SUCCESS && status == def.SUCCESS) || (when == FAILURE && status == def.FAILURE)) {
				action := cmd.group.worker.Get(cmd.Name())
				if action != nil {
					log.Infof("Action [%s] start to execute.\n", cmd.Name())
					action.Execute(ctx)
					status = <-ctx.Result
				} else {
					log.Errorf("There is no Action [%s] in Worker!\n", cmd.Name())
				}
			}
		}

		// Save status.
		stats := make([]*commandStat, len(r.cmds))
		for i, c := range r.cmds {
			cmd := c.(*command)
			stats[i] = &commandStat{
				Status:     cmd.status,
				BeginTime:  cmd.beginStamp,
				FinishTime: cmd.finishStamp,
			}
		}
		bytes, err := json.Marshal(stats)
		if err != nil {
			log.Errorf("Marshal Job [%s] Runner [%s] status failed!", r.job.name, r.id)
			return
		}

		err = ioutil.WriteFile(r.StatusPath(), bytes, os.ModePerm)
		if err != nil {
			log.Errorf("Write Job [%s] Runner [%s] status file failed!", r.job.name, r.id)
		}

		log.Debugf("Job [%s] has been completed!\n", r.job.Name())
	}()

	return nil
}

func (r *runner) Cancel() error {
	for _, c := range r.cmds {
		cmd := c.(*command)
		if cmd.group.worker != nil {
			action := cmd.group.worker.Get(cmd.Name())
			if action != nil {
				action.Cancel(r.id)
			}
		}
	}

	return nil
}

func (r *runner) Commands() []ICommand {
	return r.cmds
}

// --- Inner ---

// Dir returns the Runner working directory.
func (r *runner) Dir() string {
	return path.Join(r.job.Dir(), strconv.FormatUint(r.id, 16))
}

func (r *runner) ScriptPath() string {
	return path.Join(r.Dir(), BUBBLEFILE)
}

func (r *runner) StatusPath() string {
	return path.Join(r.Dir(), STATUSFILE)
}
