<template>
    <div>
        <b-row>
            <b-col>
                <button-group v-if="editable === false">
                    <button type="button" class="btn btn-primary" v-on:click="edit">Edit</button>
                </button-group>
                <button-group v-else>
                    <button type="button" class="btn btn-warning" v-on:click="submit">Submit</button>
                    <button type="button" class="btn btn-secondary" v-on:click="cancel">Cancel</button>
                </button-group>
            </b-col>
        </b-row>
        <b-row>
            <b-col v-if="editable === false"><pre class="text-white bg-dark">{{code}}</pre></b-col>
            <b-col v-else><b-form-textarea class="text-white bg-dark" rows="5" max-rows="30" v-model="code"></b-form-textarea></b-col>
        </b-row>
    </div>
</template>

<script>
export default {
    name: "Script",
    props: {
        target: String
    },
    data() {
        return {
            code: "",
            editable: false,
        };
    },
    mounted() {
        this.$axios
            .get(this.HOME + "jobs/" + this.target + "/script")
            .then(response => {
                this.code = atob(response.data.data);
            });
    },
    methods: {
        edit: function() {
            this.editable = true;
        },
        submit: function() {
            this.$axios
                .post(this.HOME + "jobs/" + this.target + "/script", btoa(this.code))
                .then(response => {
                    if (response.data.status == 0) {
                        this.editable = false;
                    }
                });
        },
        cancel: function() {
            this.editable = false;
        },
    }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
