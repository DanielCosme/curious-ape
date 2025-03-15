<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from "@/stores/AuthStore.ts";
import { useRouter } from "vue-router";

let authStore = useAuthStore();
let username = ref('')
let password = ref('')
let router = useRouter()

const handleLogin = async function() {
  let result = await fetch('http://localhost:4000/api/v1/login', {
    method: 'POST',
    credentials: 'include',
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    // TODO: make this payload JSON on the server.
    body: new URLSearchParams({
      "username": username.value,
      "password": password.value
    })
  })
  if (!result.ok) {
    authStore.isAuthenticated = false
    console.log(result.statusText)
    return
  }
  if (result.status === 200) {
    authStore.isAuthenticated = true
    console.log("User Authenticated");
    router.push('/');
  }
  username.value = ''
  password.value = ''
}
</script>

<template>
  <h1>Login</h1>
  <form @submit.prevent="handleLogin" novalidate>
    <div>
      <label>Username</label>
      <input v-model="username" type="text" placeholder=""></input>
    </div>
    <div>
      <label>Password</label>
      <input v-model="password" type="password" placeholder=""></input>
    </div>
    <div>
      <input type="submit" value="Login"></input>
    </div>
  </form>
</template>

<style scoped>
</style>
