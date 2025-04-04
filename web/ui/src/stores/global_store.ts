import { defineStore } from 'pinia'
import {get_ape} from "@/api/fetch.ts";

export const useGlobalStore = defineStore(
    "global_store",
    {
        state: () => ({
            info: {
                version: "",
            }
        }),
        getters: {},
        actions: {
            async fetchVersion() {
                this.info = await get_ape("/api/v1/version")
                    .then(res => res.json())
                    .catch(err => console.log(err));
            }
        }
    }
)
