// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
)

const (
	VERSION string = "v1"
	BASEURL string = "/api/" + VERSION + "/"
)

func NewWebApi() IWebControl {
	return &webapi{}
}

type webapi struct {
	handler IWebHandler
}

type result struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

func (c *webapi) Init(handler IWebHandler) {
	c.handler = handler
	c.handler.HandleFunc(BASEURL+"jobs/list", c.handleJobsList, "GET")
	c.handler.HandleFunc(BASEURL+"jobs/create/{job}", c.handleJobsCreate, "GET")
	c.handler.HandleFunc(BASEURL+"jobs/delete/{job}", c.handleJobsDelete, "DELETE")
	c.handler.HandleFunc(BASEURL+"jobs/{job}/script", c.handleJobsJobScript, "GET")
	c.handler.HandleFunc(BASEURL+"jobs/{job}/script", c.handleJobsJobScript, "POST")
	c.handler.HandleFunc(BASEURL+"jobs/{job}/crons/add/{cron}", c.handleJobsJobAddCron, "GET")
	c.handler.HandleFunc(BASEURL+"jobs/{job}/crons/remove/{id}", c.handleJobsJobRemoveCron, "DELETE")
	c.handler.HandleFunc(BASEURL+"jobs/{job}/crons/list", c.handleJobsJobListCrons, "GET")
	c.handler.HandleFunc(BASEURL+"jobs/{job}/trigger", c.handleJobsJobTrigger, "GET")
	c.handler.HandleFunc(BASEURL+"jobs/{job}/list/{index}", c.handleJobsJobList, "GET")
	c.handler.HandleFunc(BASEURL+"jobs/{job}/cancel/{runner}", c.handleJobsJobCancelRunner, "GET")
	c.handler.HandleFunc(BASEURL+"jobs/{job}/log/{runner}/{index}/{full}", c.handleJobsJobLogRunnerIndex, "GET")
	c.handler.HandleFunc(BASEURL+"workers/monitor", c.handleWorkersMonitor, "GET")
}

func (c *webapi) handleJobsList(w http.ResponseWriter, req *http.Request) {
	ret := &result{Status: 0}
	defer json.NewEncoder(w).Encode(&ret)

	list, err := c.handler.List()
	if err != nil {
		ret.Status = -1
		ret.Data = err.Error()
	} else {
		log.Debugf("Handle list Jobs [%d].\n", len(list))
		ret.Data = list
	}
}

func (c *webapi) handleJobsCreate(w http.ResponseWriter, req *http.Request) {
	ret := &result{Status: 0}
	defer json.NewEncoder(w).Encode(&ret)

	params := mux.Vars(req)
	job := params["job"]
	log.Debugf("Handle creating Job [%s].\n", job)

	err := c.handler.Create(job)
	if err != nil {
		ret.Status = -1
		ret.Data = err.Error()
	}
}

func (c *webapi) handleJobsDelete(w http.ResponseWriter, req *http.Request) {
	ret := &result{Status: 0}
	defer json.NewEncoder(w).Encode(&ret)

	params := mux.Vars(req)
	job := params["job"]
	log.Debugf("Handle deleting Job [%s].\n", job)

	err := c.handler.Delete(job)
	if err != nil {
		ret.Status = -1
		ret.Data = err.Error()
	}
}

func (c *webapi) handleJobsJobScript(w http.ResponseWriter, req *http.Request) {
	ret := &result{Status: 0}
	defer json.NewEncoder(w).Encode(&ret)

	params := mux.Vars(req)
	job := params["job"]
	log.Debugf("Handle setting Job [%s].\n", job)

	switch req.Method {
	case "GET":
		bytes, err := c.handler.JobScript(job)
		if err != nil {
			ret.Status = -1
			ret.Data = err.Error()
		} else {
			ret.Data = string(bytes)
		}
	case "POST":
		bytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			ret.Status = -1
			ret.Data = err.Error()
		} else if err = c.handler.JobSetScript(job, string(bytes)); err != nil {
			ret.Status = -1
			ret.Data = err.Error()
		}
	}
}

func (c *webapi) handleJobsJobAddCron(w http.ResponseWriter, req *http.Request) {
	ret := &result{Status: 0}
	defer json.NewEncoder(w).Encode(&ret)

	params := mux.Vars(req)
	job := params["job"]
	cron, _ := strconv.Atoi(params["cron"])
	log.Debugf("Handle adding cron [%d] of Job [%s].\n", cron, job)

	data, err := c.handler.JobAddCron(job, cron)
	if err != nil {
		ret.Status = -1
		ret.Data = err.Error()
	} else {
		ret.Data = data
	}
}

func (c *webapi) handleJobsJobRemoveCron(w http.ResponseWriter, req *http.Request) {
	ret := &result{Status: 0}
	defer json.NewEncoder(w).Encode(&ret)

	params := mux.Vars(req)
	job := params["job"]
	id, _ := strconv.ParseUint(params["id"], 16, 64)
	log.Debugf("Handle removing cron id [%d] of Job [%s].\n", id, job)

	err := c.handler.JobRemoveCron(job, id)
	if err != nil {
		ret.Status = -1
		ret.Data = err.Error()
	}
}

func (c *webapi) handleJobsJobListCrons(w http.ResponseWriter, req *http.Request) {
	ret := &result{Status: 0}
	defer json.NewEncoder(w).Encode(&ret)

	params := mux.Vars(req)
	job := params["job"]
	log.Debugf("Handle listing crons of Job [%s].\n", job)

	crons, err := c.handler.JobListCrons(job)
	if err != nil {
		ret.Status = -1
		ret.Data = err.Error()
	} else {
		ret.Data = crons
	}
}

func (c *webapi) handleJobsJobTrigger(w http.ResponseWriter, req *http.Request) {
	ret := &result{Status: 0}
	defer json.NewEncoder(w).Encode(&ret)

	params := mux.Vars(req)
	job := params["job"]
	log.Debugf("Handle scheduling Job [%s].\n", job)

	if err := c.handler.JobTrigger(job); err != nil {
		ret.Status = -1
		ret.Data = err.Error()
	}
}

func (c *webapi) handleJobsJobList(w http.ResponseWriter, req *http.Request) {
	ret := &result{Status: 0}
	defer json.NewEncoder(w).Encode(&ret)

	params := mux.Vars(req)
	job := params["job"]
	index, _ := strconv.Atoi(params["index"])
	log.Debugf("Handle list Job [%s] with index [%d].\n", job, index)

	list, err := c.handler.JobList(job, index)
	if err != nil {
		ret.Status = -1
		ret.Data = err.Error()
	} else {
		ret.Data = list
	}
}

func (c *webapi) handleJobsJobCancelRunner(w http.ResponseWriter, req *http.Request) {
	ret := &result{Status: 0}
	defer json.NewEncoder(w).Encode(&ret)

	params := mux.Vars(req)
	job := params["job"]
	runner, _ := strconv.ParseUint(params["runner"], 16, 64)
	log.Infof("Handle canceling Job [%s] Runner [%d].\n", job, runner)

	if err := c.handler.JobCancel(job, runner); err != nil {
		ret.Status = -1
		ret.Data = err.Error()
	}
}

func (c *webapi) handleJobsJobLogRunnerIndex(w http.ResponseWriter, req *http.Request) {
	ret := &result{Status: 0}
	defer json.NewEncoder(w).Encode(&ret)

	params := mux.Vars(req)
	job := params["job"]
	runner, _ := strconv.ParseUint(params["runner"], 16, 64)
	index, _ := strconv.Atoi(params["index"])
	full, _ := strconv.ParseBool(params["full"])
	log.Debugf("Handle log Job [%s] Runner [%d] Index [%d] with [%t].\n", job, runner, index, full)

	l, err := c.handler.JobLogRunnerIndex(job, runner, index, full)
	if err != nil {
		ret.Status = -1
		ret.Data = err.Error()
	} else {
		ret.Data = l
	}
}

func (c *webapi) handleWorkersMonitor(w http.ResponseWriter, req *http.Request) {
	ret := &result{Status: 0}
	defer json.NewEncoder(w).Encode(&ret)

	data, err := c.handler.Monitor()
	if err != nil {
		ret.Status = -1
		ret.Data = err.Error()
	} else {
		ret.Data = data
	}
}
