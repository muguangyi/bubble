// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package web

// IWebControl presents a web controller.
type IWebControl interface {
	// Init initialize the control with IWebHandler object.
	Init(handler IWebHandler)
}
