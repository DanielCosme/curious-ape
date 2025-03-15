<script setup lang="ts">
  import { useAuthStore } from "@/stores/AuthStore.ts";
  import { onMounted } from "vue";

  const authStore = useAuthStore();

  onMounted(() => {
    authStore.checkLogin()
  })
</script>

<template>
  <header>
  </header>
  <main>
    <nav>
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
