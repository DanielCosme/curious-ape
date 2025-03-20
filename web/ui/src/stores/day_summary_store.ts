import {ref} from "vue";

export interface DayPayload {
    month: string;
    days: DaySummary[];
}

export interface DaySummary {
    key: string
    date: string
    day: string
    wake_up_habit: HabitSummary
    fitness_habit: HabitSummary
    work_habit: HabitSummary
    eat_habit: HabitSummary
    wake_up_detail: string
    fitness_detail: string
    work_detail: string
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
