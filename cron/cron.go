// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cron

import (
	"bubble/def"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	log "github.com/cihub/seelog"
)

// JobFunc to create a target ICronJob instance.
type JobFunc func() ICronJob

// NewCron create a new ICron.
func NewCron(f JobFunc, file string) ICron {
	c := &cron{f: f, file: file}
	c.load()

	return c
}

type cron struct {
	f             JobFunc
	file          string
	triggerLocker sync.Mutex
	triggers      map[uint64]*trigger
}

func (c *cron) StartAll() {
	for _, t := range c.triggers {
		t.Start()
	}
}

func (c *cron) Flush() {
	stats := make(map[uint64]*triggerStat)
	for k, t := range c.triggers {
		stats[k] = &triggerStat{
			T:         t.t,
			LastStamp: t.lastStamp,
			Payload:   t.job.ToBytes(),
		}
	}

	bytes, err := json.Marshal(stats)
	if err != nil {
		log.Error(err)
		return
	}

	if err = ioutil.WriteFile(c.file, bytes, os.ModePerm); err != nil {
		log.Error(err)
	}
}

func (c *cron) Add(t Type) (ITrigger, error) {
	id, _ := def.NextUid()
	tr := newTrigger(c, id, t, time.Now().Unix())

	c.triggerLocker.Lock()
	defer c.triggerLocker.Unlock()
	c.triggers[id] = tr

	return tr, nil
}

func (c *cron) Remove(id uint64) (ITrigger, error) {
	t, ok := c.triggers[id]
	if !ok {
		return nil, fmt.Errorf("trigger [%d] is not exist", id)
	}

	c.triggerLocker.Lock()
	defer c.triggerLocker.Unlock()
	delete(c.triggers, id)

	return t, nil
}

func (c *cron) Triggers() ([]ITrigger, error) {
	arr := make([]ITrigger, 0)
	for _, t := range c.triggers {
		arr = append(arr, t)
	}

	return arr, nil
}

func (c *cron) Destroy() {
	for _, t := range c.triggers {
		t.Stop()
	}
}

func (c *cron) load() {
	c.triggers = make(map[uint64]*trigger)
	_, err := os.Stat(c.file)
	if err == nil {
		bytes, err := ioutil.ReadFile(c.file)
		if err != nil {
			return
		}

		var stats map[uint64]*triggerStat
		if err = json.Unmarshal(bytes, &stats); err != nil {
			return
		}

		for id, data := range stats {
			tr := newTrigger(c, id, data.T, data.LastStamp)
			tr.job.FromBytes(data.Payload)
			c.triggers[id] = tr
		}
	}
}
