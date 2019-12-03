<template>
    <div style="margin-top: 10px">
        <b-card class="border-success mb-3" v-for="w in workers" :key="w.id" v-bind:header="w.id" style="max-width: 10rem;">
            <b-card-text>
                Workload: {{w.workload}}
            </b-card-text>
        </b-card>
    </div>
</template>

<script>
import { setInterval, clearInterval } from 'timers';

export default {
    name: "Monitor",
    data() {
        return {
           workers: [],
           timer: null,
        };
    },
    mounted() {
        this.refresh();
        this.timer = setInterval(() => {
            this.refresh();
        }, 3000);
    },
    beforeDestroy() {
        clearInterval(this.timer);
    },
    methods: {
        refresh: function() {
            this.$axios.get(this.HOME+"workers/monitor")
                .then(response => {
                    if (response.data.status == 0) {
                        this.workers = response.data.data;
                    }
                })
        }
    }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
