<template>
    <div id="app">
        <AppTopbar/>
        <AppSidebar/>
        <router-view class="content"/>
        <AppFooter/>
    </div>
</template>

<script>
import AppSidebar from './components/AppSidebar.vue';
import AppTopbar from './components/AppTopbar.vue';
import AppFooter from './components/AppFooter.vue';

export default {
    name: 'App',
    components: {
        AppSidebar,
        AppTopbar,
        AppFooter
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

:root {
    --background-color: rgb(25, 25, 25);
    --text-color: rgb(211, 211, 211);
    --secondary-text-color: rgb(118, 118, 118);
    --accent-color: #0c8dd8;
    --hover-background-color: rgb(32, 32, 32);
    --progress-bar-background-color: rgba(255, 255, 255, 0.043);
    --game-box-info-background-color: rgb(30, 30, 30);
    --outline: #383838;
    --hover-sometimes: rgb(29, 29, 29);
}

html, body {
    height: 100%;
    width: 100%;
}

body {
    background-color: #000000;
    color: white;
    overflow-x: hidden; /* Prevent horizontal scrolling */
}

::-webkit-scrollbar {
    width: 10px;
    height: 10px;
}

::-webkit-scrollbar-track {
    background: var(--background-color);
}

::-webkit-scrollbar-thumb {
    background-color: var(--outline);
    border-radius: 8px;
    border: 2px solid var(--background-color);
}

::-webkit-scrollbar-thumb:hover {
    background-color: var(--accent-color);
}

#app {
    height: 100%;
    background-color: rgb(25, 24, 24);
    width: 100%;
}

.content {
    position: absolute;
    top: 40px;
    left: 200px;
    width: calc(100% - 200px);
    background-color: rgb(25, 24, 24);
}
</style>