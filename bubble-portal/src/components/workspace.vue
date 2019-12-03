<template>
    <div>
        <b-row>
            <b-col cols="2" md="auto">
                <div class="btn-group-vertical text-left" data-toggle="buttons" style="margin-top: 5px">
                    <button type="button" class="btn btn-info" v-on:click="createJob">New Job...</button>
                    <div class="btn-group" role="group" aria-label="Button group with nested dropdown" v-for="job in jobs" :key="job">
                        <button type="button" class="btn btn-outline-info" v-on:click="viewJob(job, $event)">{{job}}</button>
                        <b-dropdown variant="outline-info">
                            <b-dropdown-item v-on:click="triggerJob(job, $event)" variant="success">Trigger</b-dropdown-item>
                            <b-dropdown-divider></b-dropdown-divider>
                            <b-dropdown-item v-on:click="settingJob(job, $event)" variant="warning">Setting</b-dropdown-item>
                            <b-dropdown-divider></b-dropdown-divider>
                            <b-dropdown-item v-on:click="deleteJob(job, $event)" variant="danger">Delete</b-dropdown-item>
                        </b-dropdown>
                    </div>
                </div>
            </b-col>
            <b-col>
                <div v-if="op === 'CREATE'">
                    <Create v-on:on-create="onCreateJob" />
                </div>
                <div v-else-if="op === 'STAGE'">
                    <Stage v-bind:target="target" />
                </div>
                <div v-else-if="op === 'SETTING'">
                    <Setting v-bind:target="target" />
                </div>
                <div v-else>LOADING...</div>
            </b-col>
        </b-row>
    </div>
</template>

<script>
import Create from "./create.vue";
import Stage from "./stage.vue";
import Setting from "./setting.vue";

export default {
    name: "Workspace",
    components: {
        Create,
        Stage,
        Setting
    },
    data() {
        return {
            op: "LOADING",
            target: "",
            jobs: []
        };
    },
    mounted() {
        this.update();
    },
    methods: {
        update: function() {
            this.$axios.get(this.HOME + "jobs/list").then(response => {
                this.jobs = response.data.data;
                if (this.jobs.length > 0) {
                    this.op = "STAGE";
                    this.target = this.jobs[0];
                } else {
                    this.op = "CREATE";
                }
            });
        },
        createJob: function(evt) {
            evt.stopPropagation();
            this.op = "CREATE";
        },
        onCreateJob: function() {
            this.update();
        },
        deleteJob: function(job, evt) {
            evt.stopPropagation();
            this.$axios.delete(this.HOME + "jobs/delete/" + job).then(response => {
                if (response.data.status == 0) {
                    this.update();
                }
            })
        },
        viewJob: function(job, evt) {
            evt.stopPropagation();
            this.op = "STAGE";
            this.target = job;
        },
        settingJob: function(job, evt) {
            evt.stopPropagation();
            this.op = "SETTING";
            this.target = job;
        },
        triggerJob: function(job, evt) {
            evt.stopPropagation();
            this.$axios.get(this.HOME + "jobs/" + job + "/trigger");
        }
    }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.job {
    height: 30px;
}

.job-step {
    width: 100px;
}
</style>
