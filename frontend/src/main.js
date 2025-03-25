import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

import { AddToLibrary } from '../bindings/derpy-launcher072/library/library'
import {RealDebridClient} from "../bindings/derpy-launcher072/torrent/realdebrid/index.js";

//AddToLibrary(69)

const app = createApp(App)
app.use(router)
app.mount('#app')