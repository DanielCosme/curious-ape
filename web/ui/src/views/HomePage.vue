<script setup lang="ts">
import { useDaysStore } from "@/stores/DaysStore.ts";
import { onMounted } from "vue";


interface DaySummary {
  Date: string // TODO: Maybe I can parse this as a JS Date
  Wake: Habit  // TODO Maybe I can separate this into a different API Call.
  Fitness: Habit
  Work: Habit
  Eat: Habit
  Score: number
}

interface Habit {

}

</script>

<template>
  <h1>Days</h1>

  <h2>January</h2>
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

      <!-- TODO: This <tr> needs to be it's own component -->
      <tr>
        <td>Monday, 01</td>
        <td>
          <span>O</span>
          <a href="#">Y</a>
          <a href="#">N</a>
          <span>Detail</span>
        </td>
        <td>
          <span>O</span>
          <a href="#">Y</a>
          <a href="#">N</a>
          <span>Detail</span>
        </td>
        <td>
          <span>X</span>
          <a href="#">Y</a>
          <a href="#">N</a>
          <span>Detail</span>
        </td>
        <td>
          <span>X</span>
          <a href="#">Y</a>
          <a href="#">N</a>
          <span>Detail</span>
        </td>
        <td>
          <button class="button-primary">Sync</button>
        </td>
        <td>2</td>
      </tr>
    </tbody>
  </table>
</template>

<style scoped>

</style>