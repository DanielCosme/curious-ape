<script setup lang="ts">

// TODO: how to get ISO Date from JS Date.
// `${date1.getFullYear()}-${date1.getMonth()+1}-${date1.getDate()};

import { get_ape } from "@/api/fetch.ts";
import { onMounted } from "vue";
import Day_summary_row from "@/components/day_summary_row.vue";
import { days_summary } from "@/stores/day_summary_store.ts";


async function fetchData() {
  let res = await get_ape("/api/v1")
  days_summary.value = await res.json()
}

onMounted( () => {
  fetchData()
})
</script>

<template>
  <h2>{{ days_summary.month }}</h2>
  <table class="u-full-width">
    <thead>
      <tr>
        <th>Date</th>
        <th>Wake Up</th>
        <th>Fitness</th>
        <th>Work</th>
        <th>Eat Clean</th>
        <th></th> <!-- Sync button column -->
        <th>Score</th>
      </tr>
    </thead>
    <tbody>
      <day_summary_row v-for="day in days_summary.days" :key="day.key" :day="day"/>
    </tbody>
  </table>
</template>

<style scoped>

</style>