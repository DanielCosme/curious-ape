import { defineStore } from 'pinia'
import { get_ape } from '@/api/fetch.ts'

export const useDaysStore = defineStore(
    'DaysStore',
    {
        state: () => ({
            days: []
        }),
        // Computed
        getters: {
        },
        // Methods
        actions: {
            async fetchDays() {
                this.days = await get_ape("http://localhost:4000/api/v1")
                   .then(res => res.json())
                  .catch(err => console.log(err))
            }
        }
    }
)
