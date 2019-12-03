// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package worker

import (
	"fmt"
	"io"

	log "github.com/cihub/seelog"
)

type logger struct {
	runner *runner
	ctx    ICtx
}

func (l *logger) Info(v ...interface{}) {
	log.Info(v...)
	l.Notify([]byte(fmt.Sprint(v...)))
}

func (l *logger) Infof(format string, params ...interface{}) {
	log.Infof(format, params...)
	l.Notify([]byte(fmt.Sprintf(format, params...)))
}

func (l *logger) Debug(v ...interface{}) {
	log.Debug(v...)
	l.Notify([]byte(fmt.Sprint(v...)))
}

func (l *logger) Debugf(format string, params ...interface{}) {
	log.Debugf(format, params...)
	l.Notify([]byte(fmt.Sprintf(format, params...)))
}

func (l *logger) Warn(v ...interface{}) {
	log.Warn(v...)
	l.Notify([]byte(fmt.Sprint(v...)))
}

func (l *logger) Warnf(format string, params ...interface{}) {
	log.Warnf(format, params...)
	l.Notify([]byte(fmt.Sprintf(format, params...)))
}

func (l *logger) Error(v ...interface{}) {
	log.Error(v...)
	l.Notify([]byte(fmt.Sprint(v...)))
}

func (l *logger) Errorf(format string, params ...interface{}) {
	log.Errorf(format, params...)
	l.Notify([]byte(fmt.Sprintf(format, params...)))
}

func (l *logger) Critical(v ...interface{}) {
	log.Critical(v...)
	l.Notify([]byte(fmt.Sprint(v...)))
}

func (l *logger) Criticalf(format string, params ...interface{}) {
	log.Criticalf(format, params...)
	l.Notify([]byte(fmt.Sprintf(format, params...)))
}

func (l *logger) Std() io.Writer {
	return l
}

func (l *logger) Write(bytes []byte) (int, error) {
	fmt.Print(string(bytes))
	l.Notify(bytes)
	return len(bytes), nil
}

// --- Inner ---

func (l *logger) Notify(bytes []byte) {
	l.runner.worker.Progress(l.runner.name, l.ctx.Master(), l.ctx.UID(), bytes)
}
