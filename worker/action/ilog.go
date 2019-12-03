// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package action

import (
	"io"
)

// ILog interface.
type ILog interface {
	Info(v ...interface{})
	Infof(format string, params ...interface{})
	Debug(v ...interface{})
	Debugf(format string, params ...interface{})
	Warn(v ...interface{})
	Warnf(format string, params ...interface{})
	Error(v ...interface{})
	Errorf(format string, params ...interface{})
	Critical(v ...interface{})
	Criticalf(format string, params ...interface{})
	Std() io.Writer
}
