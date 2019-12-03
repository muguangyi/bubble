// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package def

import (
	"github.com/giant-tech/go-service/framework/idata"
	"github.com/sony/sonyflake"
)

const (
	MasterService idata.ServiceType = 101
	WorkerService idata.ServiceType = 102
)

// NextUid generates a new unique ID of uint64 type.
func NextUid() (uint64, error) {
	return uid.NextID()
}

var uid *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	uid = sonyflake.NewSonyflake(st)
	if uid == nil {
		panic("sonyflake not created")
	}
}
