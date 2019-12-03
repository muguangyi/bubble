// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package worker

// IProvider interface.
type IProvider interface {
	BeforeSend()
	Send()
}
