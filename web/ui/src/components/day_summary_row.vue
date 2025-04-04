<script setup lang="ts">
import {type DaySummary, HabitState, HabitType} from "@/stores/day_summary_store.ts";
import {post_ape} from "@/api/fetch.ts";
import {ref} from 'vue'

interface Props {
  day: DaySummary;
}

let props = defineProps<Props>()
let ds = ref(props.day)

async function sync() {
  let res = await post_ape(`/api/v1/days/sync?day=${ds.value.date}`, null)
  ds.value = await res.json()
}

async function update_habit(ht: HabitType, hs: HabitState) {
  let res = await post_ape(`/api/v1/habits/update?date=${ds.value.date}&type=${ht}&state=${hs}`, null)
  ds.value = await res.json()
}

</script>

<template>
  <tr>
    <td>{{ ds.day }}</td>
    <td>
      <span class="habit-done" v-if="ds.wake_up_habit.state === HabitState.Done"><strong>O</strong></span>
      <span class="habit-not-done" v-else-if="ds.wake_up_habit.state === HabitState.NotDone">X</span>
      <span v-else-if="ds.wake_up_habit.state === HabitState.NoInfo">_</span>
      <button @click="update_habit(ds.wake_up_habit.type, HabitState.Done)" class="habit-button">Y</button>
      <button @click="update_habit(ds.wake_up_habit.type, HabitState.NotDone)" class="habit-button">N</button>
      <span>{{ ds.wake_up_detail }}</span>
    </td>
    <td>
      <span class="habit-done" v-if="ds.fitness_habit.state === HabitState.Done"><strong>O</strong></span>
      <span class="habit-not-done" v-else-if="ds.fitness_habit.state === HabitState.NotDone">X</span>
      <span v-else-if="ds.fitness_habit.state === HabitState.NoInfo">_</span>
      <button @click="update_habit(ds.fitness_habit.type, HabitState.Done)" class="habit-button">Y</button>
      <button @click="update_habit(ds.fitness_habit.type, HabitState.NotDone)" class="habit-button">N</button>
      <span>{{ ds.fitness_detail }}</span>
    </td>
    <td>
      <span class="habit-done" v-if="ds.work_habit.state === HabitState.Done"><strong>O</strong></span>
      <span class="habit-not-done" v-else-if="ds.work_habit.state === HabitState.NotDone">X</span>
      <span v-else-if="ds.work_habit.state === HabitState.NoInfo">_</span>
      <button @click="update_habit(ds.work_habit.type, HabitState.Done)" class="habit-button">Y</button>
      <button @click="update_habit(ds.work_habit.type, HabitState.NotDone)" class="habit-button">N</button>
      <span>{{ ds.work_detail }}</span>
    </td>
    <td>
      <span class="habit-done" v-if="ds.eat_habit.state === HabitState.Done"><strong>O</strong></span>
      <span class="habit-not-done" v-else-if="ds.eat_habit.state === HabitState.NotDone">X</span>
      <span v-else-if="ds.eat_habit.state === HabitState.NoInfo">_</span>
      <button @click="update_habit(ds.eat_habit.type, HabitState.Done)" class="habit-button">Y</button>
      <button @click="update_habit(ds.eat_habit.type, HabitState.NotDone)" class="habit-button">N</button>
    </td>
    <td>
      <button class="button-primary" @click="sync">Sync</button>
    </td>
    <td>{{ ds.score }}</td>
  </tr>
</template>

<style scoped>
.habit-done {
  color: forestgreen;
}
.habit-not-done {
  color: red;
}
span {
  padding: 0 0.3rem;
}
a {
  padding: 0 0.2rem;
}
.habit-button {
  padding: 4px;
  height: 20px;
  line-height: 0;
  margin: 1px;
}
.habit-button:hover {
  background-color: rgba(115,255,236,0.25);
}
</style>
