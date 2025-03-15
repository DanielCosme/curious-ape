<script setup lang="ts">
  import { useAuthStore } from "@/stores/AuthStore.ts";
  import { onBeforeMount } from "vue";

  const authStore = useAuthStore();

  onBeforeMount(() => {
    authStore.checkLogin()
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
  </main>
</template>

<style scoped>
</style>
