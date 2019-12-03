// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package worker

import (
	"bubble/env"
)

// NewCtx method create a new ICtx by parameters.
func NewCtx(master, uid uint64, script, variables env.IAny, target string, env env.IEnv) ICtx {
	return &ctx{master: master, uid: uid, script: script, variables: variables, target: target, env: env}
}

type ctx struct {
	master    uint64
	uid       uint64
	script    env.IAny
	variables env.IAny
	target    string
	env       env.IEnv
}

func (c *ctx) Master() uint64 {
	return c.master
}

func (c *ctx) UID() uint64 {
	return c.uid
}

func (c *ctx) Script() env.IAny {
	return c.script
}

func (c *ctx) Variables() env.IAny {
	return c.variables
}

func (c *ctx) Target() string {
	return c.target
}

func (c *ctx) Env() env.IEnv {
	return c.env
}
