<template>
    <div>
        <b-row>
            <b-col>
                <b-form-input v-model="form.jobname" placeholder="Enter job name..."></b-form-input>
                <b-button v-on:click="onCreate" variant="danger">Create</b-button>
            </b-col>
        </b-row>
    </div>
</template>

<script>
export default {
    name: "Create",
    data() {
        return {
            form: {
                jobname: ""
            }
        };
    },
    methods: {
        onCreate() {
            this.$axios
                .get(this.HOME + "jobs/create/" + this.form.jobname)
                .then(response => {
                    if (response.data.status == 0) {
                        this.$emit("on-create")
                    }
                })
                .catch(err => {
                    this.form.jobname = "";
                    alert(err);
                });
        }
    }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
