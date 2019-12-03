import axios from "axios"
import { setInterval, clearInterval } from "timers";

const HOME = "/api/v1/";

const ONGOING = 2;
const PENDING = 3;

export class Job {
    constructor() {
        this.target = null;
        this.runners = [];
        this.runnerMap = new Map();
        this.index = 1;
        this.total = 0;
        this.perPage = 0;
        this.timer = null;
        this.questing = false;
    }
    bind(target) {
        if (target == null || this.target == target) {
            return;
        }

        this.target = target;
        this.runners = [];
        this.runnerMap = new Map();
        this.index = 1;
        this.total = 0;
        this.questing = true;

        if (this.timer == null) {
            this.timer = setInterval(() => { this._refresh(); }, 3000)
        }
    }
    destroy() {
        if (this.timer != null) {
            clearInterval(this.timer);
            this.timer = null;
        }
    }
    navigate() {
        this._refresh();
    }
    getRunner(runner) {
        return this.runnerMap.get(runner);
    }
    isLoading() {
        return this.questing;
    }
    _refresh() {
        axios.get(HOME + "jobs/" + this.target + "/list/" + (this.index - 1))
            .then(response => {
                if (response.data.status == 0) {
                    var stat = response.data.data;
                    this.total = stat.total;
                    this.perPage = stat.perpage;
                    this.runners.length = 0;
                    for (var i = 0; i < stat.runners.length; ++i) {
                        var d = stat.runners[i];
                        var r = null;
                        if (this.runnerMap.has(d.id)) {
                            r = this.runnerMap.get(d.id);
                            r.refresh(d)
                        } else {
                            r = new Runner(this, d)
                            this.runnerMap.set(d.id, r)
                        }
                        this.runners.push(r)
                    }
                }
            })
            .finally(() => {
                this.questing = false;
            })
    }
}

class Runner {
    constructor(job, data) {
        this.job = job
        this.id = data.id;
        this.status = data.status;
        this.cmds = []

        var list = data.cmds
        for (var i = 0; i < list.length; ++i) {
            var c = new Cmd(this, list[i])
            this.cmds.push(c)
        }
    }
    getCmd(step) {
        if (step < 0 || step >= this.cmds.length) {
            return null;
        }

        return this.cmds[step]
    }
    refresh(data) {
        this.status = data.status;

        var list = data.cmds;
        for (var i = 0; i < list.length; ++i) {
            this.cmds[i].refresh(list[i]);
        }
    }
    cancelable() {
        return (this.status == ONGOING || this.status.PENDING);
    }
    cancel() {
        axios.get(HOME + "jobs/" + this.job.target + "/cancel/" + this.id);
    }
}

class Cmd {
    constructor(runner, data) {
        this.runner = runner;
        this.index = data.index;
        this.name = data.name;
        this.alias = data.alias.length >= 15 ? data.alias.substring(0, 14) + "..." : data.alias;
        this.status = data.status;
        this.measure = data.measure;
        this.logFull = false;
        this.logs = null;
        this.logIsFull = false;
        this.questing = true;
    }
    refresh(data) {
        this.status = data.status;
        this.measure = data.measure;
    }
    questFull() {
        this.logFull = true;
        this._questLogs();
    }
    load() {
        this.questing = true;
        if (this.logs == null || this.status == PENDING || this.status == ONGOING) {
            this._questLogs();
        }
    }
    reset() {
        this.questing = false;
    }
    isLoading() {
        return this.logs === null;
    }
    _questLogs() {
        axios.get(HOME + "jobs/" + this.runner.job.target + "/log/" + this.runner.id + "/" + this.index + "/" + this.logFull)
            .then(response => {
                if (response.data.status == 0) {
                    var result = response.data.data;
                    this.logIsFull = result.full;
                    this.logs = atob(result.log);
                }

                if (this.questing && (this.status == PENDING || this.status == ONGOING)) {
                    setTimeout(() => { this._questLogs(); }, 1000)
                }
            });
    }
}