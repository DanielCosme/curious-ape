import { defineStore } from 'pinia'
import { post_ape } from "@/api/fetch.ts";

export const useAuthStore = defineStore(
   'authStore',
    {
        state: () => ({
            isAuthenticated: true,
            user: null
        }),
        getters: {},
        actions: {
            async checkLogin() {
                let res = await fetch('http://localhost:4000/api/v1/user', {
                   credentials: 'include'
                })
                if (!res.ok) {
                    if (res.status === 401) {
                        this.isAuthenticated = false
                        console.log(res.statusText);
                        return
                    }
                }
                if (res.status === 200) {
                    this.isAuthenticated = true
                }
            },
            async logout() {
                let res = await post_ape('http://localhost:4000/api/v1/logout', null)
                if (res.ok) {
                    this.isAuthenticated = false
                }
            }
        }
    }
)
