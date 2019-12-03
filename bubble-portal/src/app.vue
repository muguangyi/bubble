<template>
    <b-container fluid id="app">
        <Main />
    </b-container>
</template>

<script>
import Main from "./components/main.vue";

export default {
    name: "App",
    components: {
        Main
    }
};
</script>

<style>
#app {
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    flex: 1 0 auto;
}

html, body {
  height: 100%;
}

.modal-backdrop {
    opacity: .5;
}

.modal-header {
    padding: 0em 1em;
}

.modal-body {
    padding: 2px;
}

pre {
    padding: 1em;
}
</style>
