import axios from "axios"

const HOME = "/api/v1/";

export class CronJob {
    constructor(job) {
        this.job = job;
        this.items = [];
        axios.get(HOME + "jobs/" + this.job + "/crons/list")
            .then(response => {
                if (response.data.status == 0) {
                    var data = response.data.data;
                    for (var i = 0; i < data.length; ++i) {
                        this.items.push(new Item(data[i]));
                    }
                }
            })
    }
    add(type) {
        axios.get(HOME + "jobs/" + this.job + "/crons/add/" + type)
            .then(response => {
                if (response.data.status == 0) {
                    this.items.push(new Item(response.data.data));
                }
            })
    }
    remove(id) {
        axios.delete(HOME + "jobs/" + this.job + "/crons/remove/" + id)
            .then(response => {
                if (response.data.status == 0) {
                    for (var i = 0; i < this.items.length; ++i) {
                        if (this.items[i].id == id) {
                            this.items.splice(i, 1)
                            break;
                        }
                    }
                }
            })
    }
}

class Item {
    constructor(data) {
        this.id = data.id
        this.type = data.type
    }
}