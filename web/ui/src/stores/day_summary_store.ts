import {ref} from "vue";

export interface DayPayload {
    month: string;
    days: DaySummary[];
}

export interface DaySummary {
    key: string
    date: string
    day: string
    wake_up: HabitSummary
    fitness: HabitSummary
    work: HabitSummary
    eat: HabitSummary
    score: number
}

export interface HabitSummary {
    state: HabitState
    type: HabitType
}

export enum HabitState {
    Done = "done",
    NotDone = "not_done",
    NoInfo = "no_info",
}

export enum HabitType {
    WakeUp = "wake_up",
    Fitness = "fitness",
    DeepWork = "deep_work",
    Food = "food",
}

export let days_summary = ref<DayPayload>({
    month: "Undefined",
    days: []
})
