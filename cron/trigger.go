// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cron

import (
	"encoding/json"
	"time"

	log "github.com/cihub/seelog"
)

func newTrigger(cron *cron, id uint64, t Type, lastStamp int64) *trigger {
	return &trigger{cron: cron, job: cron.f(), id: id, t: t, lastStamp: lastStamp, timer: nil}
}

type trigger struct {
	cron      *cron
	job       ICronJob
	id        uint64
	t         Type
	lastStamp int64
	timer     *time.Timer
}

type triggerStat struct {
	T         Type            `json:"type"`
	LastStamp int64           `json:"last"`
	Payload   json.RawMessage `json:"payload"`
}

func (t *trigger) Id() uint64 {
	return t.id
}

func (t *trigger) Type() Type {
	return t.t
}

func (t *trigger) Job() ICronJob {
	return t.job
}

func (t *trigger) Start() {
	if t.timer != nil {
		return
	}

	elapsedTime := time.Now().Unix() - t.lastStamp
	d := t.interval()
	d -= time.Duration(elapsedTime) * time.Second
	t.process(d)
}

func (t *trigger) Stop() {
	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}
}

func (t *trigger) interval() time.Duration {
	var d time.Duration
	switch t.t {
	case QUARTERHOURLY:
		d = 15 * time.Minute
	case HOURLY:
		d = 1 * time.Hour
	case DAILY:
		d = 24 * time.Hour
	case WEEKLY:
		d = 7 * 24 * time.Hour
	case MONTHLY:
		d = 30 * 24 * time.Hour
	}

	return d
}

func (t *trigger) process(duration time.Duration) {
	// If the duration is less than 0, then start the trigger after 1 second.
	if duration <= 0 {
		duration = 1 * time.Second
	}
	log.Infof("A trigger of type [%d] will be start after %d mins.\n", t.t, duration/time.Minute)

	t.timer = time.NewTimer(duration)
	go func() {
		<-t.timer.C

		t.lastStamp = time.Now().Unix()
		t.job.Execute()
		t.Stop()

		if t.job.Repeat() {
			t.process(t.interval())
		} else {
			t.cron.Remove(t.id)
		}
		t.cron.Flush()
	}()
}
