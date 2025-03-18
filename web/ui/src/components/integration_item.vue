<script setup lang="ts">
  // TODO: Make this component cached with a <KeepAlive>
  import { onMounted, ref } from 'vue'
  import {get_ape} from "@/api/fetch.ts";
  import IntegrationsPage from "@/views/IntegrationsPage.vue";

  interface IntegrationItem {
    name: string
    isConnected: boolean
    Info: string[]
    Problem?: string
  }

  interface Props {
    item: IntegrationItem
  }

  let props = defineProps<Props>()

  let refs = ref({
    integration: {
      name: props.item.name,
      isConnected: props.item.isConnected,
      Info: props.item.Info,
      Problem: props.item.Problem,
    }
  })

  const getIntegration = async (name: string) => {
    let res = await get_ape(`http://localhost:4000/api/v1/integrations/${name.toLowerCase()}`)
    refs.value.integration = await res.json()
    console.log(refs.value.integration)
  }

  onMounted(async () => {
    await getIntegration(props.item.name)
  })
</script>

<template>
  <article>
    <header>
      <h3>{{ refs.integration.name }}</h3>
      <p>
        <strong>State:</strong> {{ refs.integration.isConnected ? 'Connected' : 'Disconnected' }}
      </p>
      <ul v-if="refs.integration.isConnected" v-for="detail in refs.integration.Info">
        <li>{{ detail }}</li>
      </ul>
      <p v-else-if="refs.integration.Problem !== ''"> <strong>ERROR: </strong>{{ refs.integration.Problem }}</p>
    </header>
  </article>
</template>

<style scoped>

</style>
