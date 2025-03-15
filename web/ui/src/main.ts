import './assets/normalize.css'
import './assets/skeleton.css'

import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import { createPinia } from 'pinia'

import App from './App.vue'
import { routes } from './router.ts'
import { useAuthStore } from "@/stores/AuthStore.ts";

const app = createApp(App)
const router = createRouter({
    history: createWebHashHistory(),
    routes
})
const pinia = createPinia()

app.config.errorHandler = (err, vm, info) => {
    console.log("APE_ERR: ", err, info)
}

app.use(router)
app.use(pinia)

const authStore = useAuthStore();
router.beforeEach((to, from) => {
    if (!authStore.isAuthenticated && to.path !== "/login") {
        return { path: "/login" }
    }
    if (to.path === "/login" && authStore.isAuthenticated) {
        return { path: "/" }
    }
})

app.mount('#app')

// TODO: Implement a global store cleanup/purge.