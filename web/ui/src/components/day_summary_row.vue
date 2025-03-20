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
  let res = await post_ape(`http://localhost:4000/api/v1/days/sync?day=${ds.value.date}`, null)
  ds.value = await res.json()
}

async function update_habit(ht: HabitType, hs: HabitState) {
  let res = await post_ape(`http://localhost:4000/api/v1/habits/update?date=${ds.value.date}&type=${ht}&state=${hs}`, null)
  ds.value = await res.json()
}

</script>

<template>
  <tr>
    <td>{{ ds.day }}</td>
    <td>
      <span v-if="ds.wake_up.state === HabitState.Done">O</span>
      <span v-else-if="ds.wake_up.state === HabitState.NotDone">X</span>
      <span v-else-if="ds.wake_up.state === HabitState.NoInfo">_</span>
      <button @click="update_habit(ds.wake_up.type, HabitState.Done)" class="habit-button">Y</button>
      <button @click="update_habit(ds.wake_up.type, HabitState.NotDone)" class="habit-button">N</button>
    </td>
    <td>
      <span v-if="ds.fitness.state === HabitState.Done">O</span>
      <span v-else-if="ds.fitness.state === HabitState.NotDone">X</span>
      <span v-else-if="ds.fitness.state === HabitState.NoInfo">_</span>
      <button @click="update_habit(ds.fitness.type, HabitState.Done)" class="habit-button">Y</button>
      <button @click="update_habit(ds.fitness.type, HabitState.NotDone)" class="habit-button">N</button>
    </td>
    <td>
      <span v-if="ds.work.state === HabitState.Done">O</span>
      <span v-else-if="ds.work.state === HabitState.NotDone">X</span>
      <span v-else-if="ds.work.state === HabitState.NoInfo">_</span>
      <button @click="update_habit(ds.work.type, HabitState.Done)" class="habit-button">Y</button>
      <button @click="update_habit(ds.work.type, HabitState.NotDone)" class="habit-button">N</button>
    </td>
    <td>
      <span v-if="ds.eat.state === HabitState.Done">O</span>
      <span v-else-if="ds.eat.state === HabitState.NotDone">X</span>
      <span v-else-if="ds.eat.state === HabitState.NoInfo">_</span>
      <button @click="update_habit(ds.eat.type, HabitState.Done)" class="habit-button">Y</button>
      <button @click="update_habit(ds.eat.type, HabitState.NotDone)" class="habit-button">N</button>
    </td>
    <td>
      <button class="button-primary" @click="sync">Sync</button>
    </td>
    <td>2</td>
  </tr>
</template>

<style scoped>
span {
  padding: 0 0.3rem;
}
a {
  padding: 0 0.2rem;
}
.habit-button {
  padding: 4px;
  margin-bottom: 0;
  height: 20px;
  line-height: 0;
}
.habit-button:hover {
  background-color: rgba(115,255,236,0.25);
}
</style>

// Emit event from here
// The parent component will go the Get request and refresh the child component?
