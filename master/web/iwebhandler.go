// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package web

import (
	"encoding/json"
	"net/http"
)

// IWebHandler presents a web handler.
type IWebHandler interface {
	// HandleFunc registers request relative path with handle func.
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request), method string)

	// HandleStatic registers request relative path files with handle func.
	HandleStatic(path string, f func(http.ResponseWriter, *http.Request))

	// Create a Job by name.
	Create(job string) error

	// Delete a Job by name.
	Delete(job string) error

	// List all Job names.
	List() ([]string, error)

	// JobScript returns the Job script code.
	JobScript(job string) (string, error)

	// JobSetScript update target Job script code.
	JobSetScript(job string, script string) error

	// JobAddCron add a cron with type.
	JobAddCron(job string, cronType int) (json.RawMessage, error)

	// JobRemoveCron remove a cron at index.
	JobRemoveCron(job string, id uint64) error

	// JobListCrons list all crons of the Job.
	JobListCrons(job string) (json.RawMessage, error)

	// JobTrigger to trigger the target Job with script data.
	JobTrigger(job string) error

	// JobCancel to cancel the target Job.
	JobCancel(job string, runner uint64) error

	// JobList list runner info of page index of the target Job.
	JobList(job string, index int) (json.RawMessage, error)

	// JobLogRunnerIndex quest target runner index detail log info.
	JobLogRunnerIndex(job string, runner uint64, index int, full bool) (json.RawMessage, error)

	// Monitor is tracking all Worker status.
	Monitor() (json.RawMessage, error)
}
