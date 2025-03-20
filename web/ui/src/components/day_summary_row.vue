<script setup lang="ts">

import {type DaySummary, HabitState} from "@/stores/day_summary_store.ts";
import {post_ape} from "@/api/fetch.ts";
import { ref } from 'vue'

interface Props {
  day: DaySummary;
}

let props = defineProps<Props>()
let ds = ref(props.day)

async function sync() {
  let res = await post_ape(`http://localhost:4000/api/v1/sync?day=${ds.value.date}`, null)
  ds.value = await res.json()
  console.log(ds.value)
}
</script>

<template>
  <tr>
    <td>{{ ds.day }}</td>
    <td>
      <span v-if="ds.wake_up.state === HabitState.Done">O</span>
      <span v-else-if="ds.wake_up.state === HabitState.NotDone">X</span>
      <span v-else-if="ds.wake_up.state === HabitState.NoInfo">_</span>
      <a href="#">Y</a>
      <a href="#">N</a>
    </td>
    <td>
      <span v-if="ds.fitness.state === HabitState.Done">O</span>
      <span v-else-if="ds.fitness.state === HabitState.NotDone">X</span>
      <span v-else-if="ds.fitness.state === HabitState.NoInfo">_</span>
      <a href="#">Y</a>
      <a href="#">N</a>
    </td>
    <td>
      <span v-if="ds.work.state === HabitState.Done">O</span>
      <span v-else-if="ds.work.state === HabitState.NotDone">X</span>
      <span v-else-if="ds.work.state === HabitState.NoInfo">_</span>
      <a href="#">Y</a>
      <a href="#">N</a>
    </td>
    <td>
      <span v-if="ds.eat.state === HabitState.Done">O</span>
      <span v-else-if="ds.eat.state === HabitState.NotDone">X</span>
      <span v-else-if="ds.eat.state === HabitState.NoInfo">_</span>
      <a href="#">Y</a>
      <a href="#">N</a>
    </td>
    <td>
      <button class="button-primary" @click="sync">Sync</button>
    </td>
    <td>2</td>
  </tr>
</template>

<style scoped>
span {
  padding: 0 0.2rem;
}
a {
  padding: 0 0.2rem;
}
</style>

// Emit event from here
// The parent component will go the Get request and refresh the child component?
