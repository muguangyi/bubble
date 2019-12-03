// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

import (
	"bubble/def"
	"bubble/env"
	"path"
	"strconv"
	"time"
)

// NewCommand method create a new command by runner and index.
func NewCommand(runner *runner, index int) ICommand {
	c := &command{
		runner:      runner,
		index:       index,
		name:        "unknown",
		alias:       "",
		disk:        "",
		variables:   env.NewAny(nil),
		when:        "success",
		where:       -1,
		target:      "",
		prefer:      "",
		status:      def.NOTSTART,
		beginStamp:  -1,
		finishStamp: -1,
	}
	c.payloader = newPayloader(c.LogFilePath())

	return c
}

type command struct {
	runner      *runner
	index       int
	name        string
	alias       string
	disk        string
	script      env.IAny
	variables   env.IAny
	when        string
	where       int
	target      string
	prefer      string
	group       *group
	status      def.STATUS
	beginStamp  int64
	finishStamp int64
	payloader   *payloader
}

type commandStat struct {
	Status     def.STATUS `json:"status"`
	BeginTime  int64      `json:"begin"`
	FinishTime int64      `json:"finish"`
}

// --- ICommand ---

func (c *command) Index() int {
	return c.index
}

func (c *command) Name() string {
	return c.name
}

func (c *command) Alias() string {
	if c.alias == "" {
		return c.name
	}

	return c.alias
}

func (c *command) Disk() string {
	return c.disk
}

func (c *command) Script() env.IAny {
	return c.script
}

func (c *command) Variables() env.IAny {
	return c.variables
}

func (c *command) When() WHEN {
	switch c.when {
	case "always":
		return ALWAYS
	case "failure":
		return FAILURE
	default:
		return SUCCESS
	}
}

func (c *command) Where() int {
	return c.where
}

func (c *command) Target() string {
	return c.target
}

func (c *command) Prefer() string {
	return c.prefer
}

func (c *command) Status() def.STATUS {
	return c.status
}

func (c *command) Measure() int64 {
	if c.beginStamp == -1 {
		return -1
	}

	if c.finishStamp == -1 {
		return time.Now().Unix() - c.beginStamp
	}

	return c.finishStamp - c.beginStamp
}

func (c *command) Logs(full bool) ([]byte, bool, error) {
	return c.payloader.Bytes(full)
}

func (c *command) Notify(status def.STATUS, payload []byte) error {
	c.payloader.Write(payload)

	switch c.status = status; status {
	case def.ONGOING:
		if c.beginStamp == -1 {
			c.beginStamp = time.Now().Unix()
		}
	case def.SUCCESS, def.FAILURE, def.CANCEL, def.INTERRUPT:
		c.finishStamp = time.Now().Unix()
		c.payloader.Flush()
	}

	return nil
}

// --- Inner ---

func (c *command) LogFilePath() string {
	name := "." + strconv.Itoa(c.index) + ".log"
	return path.Join(c.runner.Dir(), name)
}

type group struct {
	cmds   []ICommand
	worker IWorker
}
