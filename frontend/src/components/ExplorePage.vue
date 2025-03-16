<template>
    <div class="page">
      <h1>Explore</h1>
      <p>Discover new content here.</p>
      <button @click="real()">bruh button</button>
    </div>

</template>

<script setup>
import {RealDebridClient} from "../../bindings/derpy-launcher072/torrent/realdebrid/index.js";

function real() {
    RealDebridClient.GetTrafic().then((result) => {
        console.log(result)
    })
}
</script>

<script>
export default {
    name: 'ExplorePage'
}

</script>

<style scoped>
.page {
    padding: 20px;
}

h1 {
    margin-bottom: 15px;
}
</style>