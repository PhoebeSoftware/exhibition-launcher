import {createRouter, createWebHistory} from 'vue-router'
import ExplorePage from './components/ExplorePage.vue'
import LibraryPage from './components/LibraryPage.vue'

const routes = [
    {path: '/', redirect: '/explore'},
    {path: '/explore', component: ExplorePage},
    {path: '/library', component: LibraryPage}
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router