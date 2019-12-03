<template>
    <div>
        <b-tabs>
            <b-tab title="Script">
                <Script v-bind:target="target" />
            </b-tab>
            <b-tab title="Cron Jobs">
                <Cron v-bind:target="target" />
            </b-tab>
        </b-tabs>
    </div>
</template>

<script>
import Script from "./script.vue"
import Cron from "./cron.vue"

export default {
    name: "Setting",
    components: {
        Script,
        Cron,
    },
    props: {
        target: String
    },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
