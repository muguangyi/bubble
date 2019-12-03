// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bubble/def"
	"bubble/master"

	"github.com/giant-tech/go-service/framework/app"
	"github.com/giant-tech/go-service/framework/service"
)

const (
	ServiceName       string = "Master"
	ServiceConfigPath string = "./master.toml"
)

func main() {
	service.RegService(def.MasterService, ServiceName, &master.Master{})

	app.Run(ServiceConfigPath)
}
