<script setup lang="ts">
  import { useAuthStore } from "@/stores/AuthStore.ts";
  import { useGlobalStore } from '@/stores/global_store.ts';
  import { onBeforeMount, onMounted, watch } from "vue";
  import { useRouter } from 'vue-router'

  const router = useRouter()
  const authStore = useAuthStore();
  const globalStore = useGlobalStore();

  onBeforeMount(() => {
    authStore.checkLogin()
  })
  onMounted(() => {
    globalStore.fetchVersion()
  })
  watch(() => authStore.isAuthenticated, async (isAuthenticated) => {
    if (!isAuthenticated) {
      authStore.isAuthenticated = false
      await router.push('/login')
    }
  })
</script>

<template>
  <header>
  </header>
  <main>
    <nav>
      <h1>Curious Ape</h1>
      <template v-if="authStore.isAuthenticated">
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/integrations">Integrations</RouterLink>
      </template>
    </nav>
    <RouterView/>
    <button @click="authStore.logout" v-if="authStore.isAuthenticated">Logout</button>
    <footer>
      <p>{{ globalStore.info.version }}</p>
    </footer>
  </main>
</template>
<style scoped>
nav a {
  margin-right: 1rem;
}
</style>
