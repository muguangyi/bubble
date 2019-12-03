// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

import (
	"bubble/def"
	"bubble/env"
)

// NewCtx method create a new ctx by runner.
func NewCtx(runner IRunner) ICtx {
	c := &ctx{runner: runner, Cmd: nil, Result: make(chan def.STATUS, 1), env: env.NewEnv()}
	c.env.Set("_INSTANCE", env.NewAny(runner.ID())) // Set "_INSTANCE" variable.

	return c
}

type ctx struct {
	runner IRunner
	Cmd    ICommand
	Result chan def.STATUS
	env    env.IEnv
}

func (c *ctx) ID() uint64 {
	return c.runner.ID()
}

func (c *ctx) LastWorker() uint64 {
	if c.Cmd == nil || c.Cmd.Index() == 0 {
		return 0
	}

	cmd := c.runner.Commands()[c.Cmd.Index()-1].(*command)
	return cmd.group.worker.ID()
}

func (c *ctx) Disk() string {
	return c.Cmd.Disk()
}

func (c *ctx) Script() []byte {
	if c.Cmd == nil {
		return nil
	}

	bytes, err := c.Cmd.Script().ToBytes()
	if err != nil {
		return nil
	}

	return bytes
}

func (c *ctx) Variables() []byte {
	if c.Cmd == nil {
		return nil
	}

	bytes, err := c.Cmd.Variables().ToBytes()
	if err != nil {
		return nil
	}

	return bytes
}

func (c *ctx) Target() string {
	if c.Cmd == nil {
		return ""
	}

	return c.Cmd.Target()
}

func (c *ctx) Env() env.IEnv {
	return c.env
}

func (c *ctx) Notify(status def.STATUS, payload []byte) {
	if c.Cmd != nil {
		c.Cmd.Notify(status, payload)
	}
}

func (c *ctx) SetResult(result def.STATUS, env env.IEnv) {
	c.Notify(result, nil)
	c.env = env
	c.Result <- result
}
