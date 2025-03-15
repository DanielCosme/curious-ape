import { useAuthStore } from "@/stores/AuthStore.ts";
import { useRouter } from "vue-router";

export async function get_ape(url: string) :Promise<Response> {
    let router = useRouter();
    const authStore = useAuthStore();
    let res = await fetch(url, {
        credentials: 'include'
    })
    if (!res.ok && res.status === 401) {
        authStore.isAuthenticated = false;
        log("Not Authenticated");
        router.push('/login');
    }
    return res
}

function log(...data: any[]) {
    console.log("fetch_ape: ", data)
}
