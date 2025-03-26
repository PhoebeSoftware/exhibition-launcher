import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

//AddToLibrary(69)

const app = createApp(App)
app.use(router)
app.mount('#app')