// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package master

import (
	"bubble/cron"
	"bubble/def"
	"bubble/env"
	mweb "bubble/master/web"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// NewWeb method create a new IWeb by master and configure.
func NewWeb(master IMaster, conf env.IAny) IWeb {
	if !conf.IsMap() {
		return nil
	}

	m := conf.Map()

	port, ok := m["port"]
	if !ok {
		return nil
	}

	root, ok := m["root"]
	if !ok {
		return nil
	}

	index, ok := m["index"]
	if !ok {
		return nil
	}

	w := &web{master: master, port: port.Int()}
	w.bind(mweb.NewWebApi())
	w.bind(mweb.NewPortal(root.String(), index.String()))

	return w
}

const (
	runnersPerPage int = 20
)

type web struct {
	master   IMaster
	port     int
	controls []mweb.IWebControl
	router   *mux.Router
}

func (w *web) bind(control mweb.IWebControl) {
	w.controls = append(w.controls, control)
}

func (w *web) Serve() error {
	w.router = mux.NewRouter()
	for _, c := range w.controls {
		c.Init(w)
	}

	go http.ListenAndServe(fmt.Sprintf(":%d", w.port), w.router)

	return nil
}

func (w *web) Close() {

}

func (w *web) HandleFunc(path string, f func(http.ResponseWriter, *http.Request), method string) {
	w.router.Path(path).HandlerFunc(f).Methods(method)
}

func (w *web) HandleStatic(path string, f func(http.ResponseWriter, *http.Request)) {
	w.router.PathPrefix(path).HandlerFunc(f)
}

func (w *web) Create(job string) error {
	return w.master.Create(job)
}

func (w *web) Delete(job string) error {
	return w.master.Delete(job)
}

func (w *web) List() ([]string, error) {
	jobs := w.master.List()
	names := make([]string, len(jobs))
	for i, j := range jobs {
		names[i] = j.Name()
	}

	return names, nil
}

func (w *web) JobScript(job string) (string, error) {
	j, err := w.master.Get(job)
	if err != nil {
		return "", err
	}

	script, err := j.Script()
	if err != nil {
		return "", nil
	}

	return base64.StdEncoding.EncodeToString(script), nil
}

func (w *web) JobSetScript(job string, script string) error {
	j, err := w.master.Get(job)
	if err != nil {
		return err
	}

	bytes, err := base64.StdEncoding.DecodeString(script)
	if err != nil {
		return err
	}

	return j.SetScript(bytes)
}

func (w *web) JobAddCron(job string, cronType int) (json.RawMessage, error) {
	j, err := w.master.Get(job)
	if err != nil {
		return nil, err
	}

	t, err := j.AddTrigger(cron.Type(cronType))
	if err != nil {
		return nil, err
	}

	return json.Marshal(&triggerData{
		ID:   strconv.FormatUint(t.Id(), 16),
		Type: int(t.Type()),
	})
}

func (w *web) JobRemoveCron(job string, id uint64) error {
	j, err := w.master.Get(job)
	if err != nil {
		return err
	}

	return j.RemoveTrigger(id)
}

func (w *web) JobListCrons(job string) (json.RawMessage, error) {
	j, err := w.master.Get(job)
	if err != nil {
		return nil, err
	}

	trs, err := j.Triggers()
	if err != nil {
		return nil, err
	}

	types := make([]*triggerData, len(trs))
	for i, c := range trs {
		types[i] = &triggerData{
			ID:   strconv.FormatUint(c.Id(), 16),
			Type: int(c.Type()),
		}
	}

	return json.Marshal(types)
}

func (w *web) JobTrigger(job string) error {
	j, err := w.master.Get(job)
	if err != nil {
		return err
	}

	return j.Trigger()
}

func (w *web) JobCancel(job string, runner uint64) error {
	j, err := w.master.Get(job)
	if err != nil {
		return err
	}

	return j.Cancel(runner)
}

func (w *web) JobList(job string, index int) (json.RawMessage, error) {
	j, err := w.master.Get(job)
	if err != nil {
		return nil, err
	}

	runners := j.Runners()

	jobStat := &jobStatus{
		Index:   index,
		Runners: make([]*runnerStatus, 0),
		Total:   len(runners),
		PerPage: runnersPerPage,
	}
	for i := index * runnersPerPage; i < len(runners) && i < (index+1)*runnersPerPage; i++ {
		r := runners[i]
		rs := &runnerStatus{ID: strconv.FormatUint(r.ID(), 16), Status: def.NOTSTART}

		commands := r.Commands()
		rs.Cmds = make([]*cmdStatus, len(commands))
		for j, c := range commands {
			rs.Cmds[j] = &cmdStatus{
				Index:   c.Index(),
				Name:    c.Name(),
				Alias:   c.Alias(),
				Status:  c.Status(),
				Measure: c.Measure(),
			}

			if c.Status() > rs.Status {
				rs.Status = c.Status()
			}
		}

		jobStat.Runners = append(jobStat.Runners, rs)
	}

	return json.Marshal(jobStat)
}

func (w *web) JobLogRunnerIndex(job string, runner uint64, index int, full bool) (json.RawMessage, error) {
	j, err := w.master.Get(job)
	if err != nil {
		return nil, err
	}

	r, err := j.GetRunner(runner)
	if err != nil {
		return nil, err
	}

	cmds := r.Commands()
	if index < 0 || index >= len(cmds) {
		return nil, fmt.Errorf("runner index [%d] is out of range", index)
	}

	c := cmds[index]
	bytes, full, err := c.Logs(full)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&logData{
		Full: full,
		Log:  base64.StdEncoding.EncodeToString(bytes), // Use Base64 to encode log string.
	})
}

func (w *web) Monitor() (json.RawMessage, error) {
	workers := w.master.Workers()

	stats := make([]*workerStatus, len(workers))
	for i, w := range workers {
		stats[i] = &workerStatus{
			ID:       strconv.FormatUint(w.ID(), 16),
			Workload: w.Workload(),
		}
	}

	return json.Marshal(stats)
}

type cmdStatus struct {
	Index   int        `json:"index"`
	Name    string     `json:"name"`
	Alias   string     `json:"alias"`
	Status  def.STATUS `json:"status"`
	Measure int64      `json:"measure"`
}

type runnerStatus struct {
	ID     string       `json:"id"`
	Status def.STATUS   `json:"status"`
	Cmds   []*cmdStatus `json:"cmds"`
}

type jobStatus struct {
	Index   int             `json:"index"`
	Runners []*runnerStatus `json:"runners"`
	Total   int             `json:"total"`
	PerPage int             `json:"perpage"`
}

type workerStatus struct {
	ID       string `json:"id"`
	Workload int    `json:"workload"`
}

type triggerData struct {
	ID   string `json:"id"`
	Type int    `json:"type"`
}

type logData struct {
	Full bool   `json:"full"`
	Log  string `json:"log"`
}
