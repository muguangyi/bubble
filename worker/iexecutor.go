// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package worker

// IExecutor interface.
type IExecutor interface {
	Execute()
	Cancel()
	BeforeReceive(length, count int64, checksum string)
	Receive(index int64, data []byte)
	AfterReceive()
}
