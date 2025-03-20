import { defineStore } from 'pinia'

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
                // TODO: create an endpoint on the backend to make this more efficient.
                //      The endpoint would return maybe the current user.
                let res = await fetch('http://localhost:4000/api/v1', {
                   credentials: 'include'
                })
                if (!res.ok) {
                    if (res.status === 401) {
                        this.isAuthenticated = false
                        console.log(res.statusText);
                        return
                    }
                    console.log(res.statusText);
                }
                if (res.status === 200) {
                    this.isAuthenticated = true
                    console.log("User Authenticated");
                }
            },
        }
    }
)
