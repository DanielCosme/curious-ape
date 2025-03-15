import IntegrationsPage from '@/views/IntegrationsPage.vue'
import LoginPage from '@/views/LoginPage.vue'
import HomePage from '@/views/HomePage.vue'

export const routes = [
    {
        path: "/",
        component: HomePage
    },
    {
        path: '/login',
        component: LoginPage
    },
    {
        path: '/integrations',
        component: IntegrationsPage
    }
];
