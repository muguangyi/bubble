<template>
    <div>
        <b-row style="margin-top: 10px">
            <b-col cols="1"><strong>Interval</strong></b-col>
            <b-col>
                <b-dropdown class="bb-dropdown" v-bind:text="getSelectedType()">
                    <b-dropdown-item v-on:click="selectType(4+1)">{{values[4]}}</b-dropdown-item>
                    <b-dropdown-item v-on:click="selectType(3+1)">{{values[3]}}</b-dropdown-item>
                    <b-dropdown-item v-on:click="selectType(2+1)">{{values[2]}}</b-dropdown-item>
                    <b-dropdown-item v-on:click="selectType(1+1)">{{values[1]}}</b-dropdown-item>
                    <b-dropdown-item v-on:click="selectType(0+1)">{{values[0]}}</b-dropdown-item>
                </b-dropdown>
            </b-col>
            <b-col>
                <b-button v-on:click="addCron()">Add</b-button>
            </b-col>
        </b-row>
        <div class="bb-divider"/>
        <b-row v-for="i in job.items" :key="i">
            <b-col cols="1"><strong>Interval</strong></b-col>
            <b-col><em>{{values[i.type-1]}}</em></b-col>
            <b-col><b-button v-on:click="removeCron(i.id)">Delete</b-button></b-col>
        </b-row>
    </div>
</template>

<script>
import {CronJob} from "./cron.js"

export default {
    name: "Cron",
    props: {
        target: String
    },
    data() {
        return {
            values: ["Quaterhourly", "Hourly", "Daily", "Weekly", "Monthly"],
            type: 5,
            job: new CronJob(this.target),
        };
    },
    mounted() {
    },
    methods: {
        getSelectedType: function() {
            return this.values[this.type-1];
        },
        selectType: function(type) {
            this.type = type;
        },
        addCron: function() {
            this.job.add(this.type);
        },
        removeCron: function(id) {
            this.job.remove(id);
        }
    }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.bb-dropdown {
    width: 200px;
}

.bb-divider {
    margin-top: 10px;
    margin-bottom: 10px;
    height: 1px;
    width: 100%;
    border-top: 1px solid lightgray;
}
</style>
