<script setup lang="ts">
import { ref, onMounted } from "vue";
import { get_ape } from '@/api/fetch.ts'
import integration_item from '@/components/integration_item.vue'

let data = ref({
  integrations: []
})
const getIntegrations = async function() {
  let res = await get_ape("/api/v1/integrations")
  data.value.integrations = await res.json()
}
onMounted(async () => {
  await getIntegrations()
  console.log(data.value.integrations)
})
</script>

<template>
  <h2>Integrations</h2>
  <integration_item v-for="item in data.integrations" :item="item"/>
</template>

<style scoped>
</style>
