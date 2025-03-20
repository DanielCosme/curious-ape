import { useAuthStore } from "@/stores/AuthStore.ts";
import { useRouter } from "vue-router";

export async function get_ape(url: string) :Promise<Response> {
    return run(url, null, "get")
}

export async function post_ape(url: string, body: any) : Promise<Response> {
    return run(url, body, "post")
}

async function run(url: string, body: any, method: string) : Promise<Response> {
    const router = useRouter();
    const authStore = useAuthStore();
    let res = await fetch(url, {
        credentials: 'include',
        body: body,
        method: method,
    })
    if (!res.ok && res.status === 401) {
        authStore.isAuthenticated = false;
        log("Not Authenticated");
        await router.push('/login');
    }
    return res
}

function log(...data: any[]) {
    console.log("fetch_ape: ", data)
}
