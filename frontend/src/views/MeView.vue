<template>
  <div>
    <h1 class="text-2xl font-bold text-blue-300 mb-6">Mon profil</h1>
    <div v-if="user" class="bg-gray-800 border border-gray-700 rounded-2xl shadow-xl px-8 py-8 max-w-lg">
      <div class="flex items-center space-x-6 mb-4">
        <img v-if="user.Picture || user.picture" :src="user.Picture || user.picture" class="rounded-full shadow-lg border-4 border-blue-500" width="90" height="90" alt="Avatar" />
        <div>
          <div class="text-lg font-semibold">{{ user.Name || user.name }}</div>
          <div class="text-blue-300 text-sm font-mono">{{ user.Email || user.email }}</div>
          <div class="text-sm text-gray-400 mt-2">
            Rôles :
            <span class="ml-1 font-mono text-xs bg-blue-700 rounded px-2 py-1" v-for="role in userRoles" :key="role">{{ role }}</span>
          </div>
        </div>
      </div>
      <div class="grid grid-cols-2 gap-2 mt-2 text-gray-300">
        <div class="text-sm">ID (base) :</div>
        <div class="text-sm font-mono text-white">{{ user.ID || user.id }}</div>
        <div class="text-sm">Google ID :</div>
        <div class="text-xs font-mono text-blue-300 break-all">{{ user.GoogleID || user.google_id }}</div>
      </div>
      <div class="flex gap-4 mt-8">
        <button @click="openPopup('organisateur')" class="bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-2 rounded-lg shadow">Devenir organisateur</button>
        <button @click="openPopup('admin')" class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg shadow">Devenir admin</button>
      </div>
    </div>
    <div v-else class="text-gray-300 mt-12 text-lg">
      <em>Chargement…</em>
    </div>

    <div v-if="popup" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div class="bg-gray-800 p-8 rounded-2xl shadow-xl w-full max-w-xs text-center">
        <div class="mb-4 text-lg text-white">Confirmer : ajouter le rôle <b class="text-blue-300">{{ popup }}</b> à votre profil ?</div>
        <div class="flex gap-4 justify-center">
          <button @click="updateRole" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded shadow">Oui</button>
          <button @click="closePopup" class="bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded shadow">Annuler</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  data() {
    return {
      user: null,
      popup: null,    
    }
  },
  computed: {
    userRoles() {
      if (!this.user || !this.user.Roles) return ['user']
      try {
        return JSON.parse(this.user.Roles)
      } catch {
        return ['user']
      }
    }
  },
  async mounted() {
    const token = localStorage.getItem('token')
    const res = await axios.get('http://localhost:8080/me', {
      headers: { Authorization: `Bearer ${token}` }
    })
    this.user = res.data
  },
  methods: {
    openPopup(role) { this.popup = role },
    closePopup() { this.popup = null },
    async updateRole() {
      const token = localStorage.getItem('token')
      await axios.post('http://localhost:8080/me/role', { role: this.popup }, {
        headers: { Authorization: `Bearer ${token}` }
      })
      const res = await axios.get('http://localhost:8080/me', {
        headers: { Authorization: `Bearer ${token}` }
      })
      this.user = res.data
      this.closePopup()
    }
  }
}
</script>
