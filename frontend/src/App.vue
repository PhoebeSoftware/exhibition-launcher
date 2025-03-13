<template>
  <div id="app">
    <AppTopbar/>
    <AppSidebar/>
    <router-view class="content"/>
  </div>
</template>

<script>
import AppSidebar from './components/AppSidebar.vue';
import AppTopbar from './components/AppTopbar.vue';

export default {
  name: 'App',
  components: {
    AppSidebar,
    AppTopbar
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  font-family: "Montserrat", serif;
  font-weight: 300;
}
html, body {
  height: 100%;
  width: 100%;
}

body {
  background-color: #000000;
  color: white;
}

#app {
  height: 100%;
  width: 100%;
}

.content {
  position: absolute;
  top: 40px;
  left: 200px;
  width: calc(100% - 200px);
  height: calc(100% - 40px);
  background-color: rgb(25, 24, 24);
}
</style>