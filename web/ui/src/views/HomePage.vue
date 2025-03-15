<script setup lang="ts">
import { useDaysStore } from "@/stores/DaysStore.ts";
import { onMounted } from "vue";

const daysStore = useDaysStore();

onMounted(() => {
  daysStore.fetchDays()
})
</script>

<template>
  <h2>Days</h2>
  <pre>{{ daysStore.days }}</pre>
</template>

<style scoped>

</style>