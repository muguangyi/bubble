<template>
    <div>
        <b-navbar toggleable="lg" type="dark" variant="primary">
            <b-navbar-brand href="#">
                <img src="/icon/logo-32x32.png" class="d-inline-block align-top" alt="Bubble">
            </b-navbar-brand>
            <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
            <b-collapse id="nav-collapse" is-nav>
                <b-navbar-nav class="ml-left">
                    <b-nav-item v-on:click="goto('WORK')">Workspace</b-nav-item>
                    <b-nav-item v-on:click="goto('MONITOR')">Monitor</b-nav-item>
                </b-navbar-nav>
                <b-navbar-nav class="ml-auto">
                    <b-nav-item href="#">About</b-nav-item>
                </b-navbar-nav>
            </b-collapse>
        </b-navbar>
        <div v-if="section === 'WORK'">
            <Workspace />
        </div>
        <div v-else-if="section === 'MONITOR'">
            <Monitor />
        </div>
    </div>
</template>

<script>
import Workspace from "./workspace.vue";
import Monitor from "./monitor.vue";

export default {
    name: "Main",
    components: {
        Workspace,
        Monitor,
    },
    data() {
        return {
            section: "WORK",
        };
    },
    methods: {
        goto: function(section) {
            this.section = section;
        }
    }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
