<template>
  <div>
    <h1 class="text-2xl font-bold text-blue-300 mb-6">Mon profil</h1>
    <div v-if="user" class="bg-gray-800 border border-gray-700 rounded-2xl shadow-xl px-8 py-8 max-w-lg">
      <div class="flex items-center space-x-6 mb-4">
        <img
          v-if="user.Picture || user.picture"
          :src="user.Picture || user.picture"
          class="rounded-full shadow-lg border-4 border-blue-500"
          width="90"
          height="90"
          alt="Avatar"
        />
        <div>
          <div class="text-lg font-semibold">{{ user.Name || user.name }}</div>
          <div class="text-blue-300 text-sm font-mono">{{ user.Email || user.email }}</div>
        </div>
      </div>
      <div class="grid grid-cols-2 gap-2 mt-2 text-gray-300">
        <div class="text-sm">ID (base) :</div>
        <div class="text-sm font-mono text-white">{{ user.ID || user.id }}</div>
        <div class="text-sm">Google ID :</div>
        <div class="text-xs font-mono text-blue-300 break-all">{{ user.GoogleID || user.google_id }}</div>
      </div>
    </div>
    <div v-else class="text-gray-300 mt-12 text-lg">
      <em>Chargement…</em>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  data() {
    return { user: null }
  },
  async mounted() {
    const token = localStorage.getItem('token')
    const res = await axios.get('http://localhost:8080/me', {
      headers: { Authorization: `Bearer ${token}` }
    })
    this.user = res.data
  }
}
</script>
