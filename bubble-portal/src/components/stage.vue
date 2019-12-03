<template>
    <div>
        <div v-if="!job.isLoading()" style="margin-top: 10px">
            <b-row class="bb-row" v-for="r in job.runners" :key="r.id">
                <b-col cols="1" md="auto">
                    <b-badge
                        class="bb-badge"
                        v-bind:variant="getStatusBVariant(r.status)"
                    >{{getStatusString(r.status)}}</b-badge>
                </b-col>
                <b-col>
                    <b-list-group horizontal>
                        <b-list-group-item
                            class="py-0 d-flex justify-content-between bb-list-group-item"
                            v-for="c in r.cmds"
                            v-bind:key="c.name"
                            v-bind:variant="getStatusBVariant(c.status)"
                            v-on:click="showConsoleWindow(r.id, c.index)"
                            button
                            style="max-width: 10rem;"
                        >
                            {{c.alias}}
                            <b-badge class="bb-badge-pill" v-if="c.measure > 0" variant="info" pill>{{formatMeasure(c.measure)}}</b-badge>
                        </b-list-group-item>
                    </b-list-group>
                </b-col>
                <b-col cols="1" md="right">
                    <b-button v-if="r.cancelable()" v-on:click="r.cancel()">Cancel</b-button>
                </b-col>
                <div class="bb-divider">
            </b-row>
            <b-pagination v-model="job.index" v-on:change="job.navigate()" :per-page="job.perPage" :total-rows="rows" size="sm"></b-pagination>
        </div>
        <div v-else class="d-flex align-items-center">
            <strong>Loading...</strong>
            <b-spinner label="Spinning"></b-spinner>
        </div>
        <b-modal id="console-window" size="xl" header-bg-variant="primary" scrollable hide-footer hide-header-close>
            <template v-slot:modal-header="{ close }">
                <!-- Emulate built in modal header close button action -->
                <b-button size="sm" variant="primary" v-on:click="hideConsoleWindow('console-window')">Close</b-button>
                <h5>Console</h5>
            </template>
            <div v-if="cmd !== null && !cmd.isLoading()">
                <a v-if="!cmd.logIsFull" href="#" v-on:click="cmd.questFull()">...</a>
                <pre class="text-white bg-dark">{{cmd.logs}}</pre>
            </div>
            <div v-else class="d-flex align-items-center">
                <strong>Loading...</strong>
                <b-spinner label="Spinning"></b-spinner>
            </div>
        </b-modal>
    </div>
</template>

<script>
import {Job} from "./job.js"

export default {
    name: "Stage",
    props: {
        target: String
    },
    data() {
        return {
            job: new Job(),
            cmd: null,
        };
    },
    computed: {
        rows() {
            return this.job.total;
        }
    },
    mounted() {
        this.job.bind(this.target);
    },
    beforeDestroy() {
        this.job.destroy()
    },
    watch: {
        target: "updateTarget"
    },
    methods: {
        updateTarget: function() {
            this.job.bind(this.target);
        },
        showConsoleWindow: function(runner, index) {
            this.cmd = this.job.getRunner(runner).getCmd(index);
            this.cmd.load();
            this.$bvModal.show('console-window');
        },
        hideConsoleWindow: function() {
            this.$bvModal.hide('console-window');
            this.cmd.reset();
            this.cmd = null;
        },
        formatMeasure: function(measure) {
            var h = parseInt(measure / 3600);
            measure = measure % 3600;
            var min = parseInt(measure / 60);
            var sec = measure % 60;
            return (h < 10 ? "0" : "") + h + ":" + (min < 10 ? "0" : "") + min + ":" + (sec < 10 ? "0" : "") + sec;
        },
        getStatusString: function(status) {
            switch (status) {
                case 0:
                case 3:
                    return "pending";
                case 1:
                    return "success";
                case 2:
                    return "progress";
                case 4:
                    return "failed";
                case 5:
                    return "canceled";
                case 6:
                    return "interrupt";
                default:
                    return "";
            }
        },
        getStatusBVariant: function(status) {
            switch (status) {
                case 0:
                case 6:
                    return "dark";
                case 1:
                    return "success";
                case 2:
                    return "primary";
                case 3:
                    return "warning";
                case 4:
                    return "danger";
                case 5:
                    return "secondary";
                default:
                    return "";
            }
        },
    }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.bb-divider {
    margin-top: 10px;
    margin-bottom: 10px;
    height: 1px;
    width: 100%;
    border-top: 1px solid lightgray;
}

.bb-row {
    margin-left: 0px;
    margin-right: 0px;
}

.bb-list-group-item {
    border-style: solid;
    border-width: 0 1px 4px 1px;
    padding: 2px 2px 2px 10px;
}

.bb-badge {
    width: 70px;
}

.bb-badge-pill {
    border-radius: 0 0 10px 10px;
}
</style>
