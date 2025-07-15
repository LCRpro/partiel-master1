<template>
  <div v-if="conf"
    class="bg-gradient-to-br from-blue-950 via-gray-900 to-gray-900 rounded-3xl shadow-2xl p-10 max-w-2xl mx-auto my-12 border border-blue-800/60">
    <div class="flex flex-col sm:flex-row items-start sm:items-center gap-7 mb-8">
      <div class="flex-1 min-w-0">
        <h1 class="text-3xl sm:text-4xl font-black text-white mb-2 tracking-tight drop-shadow">
          {{ conf.Title }}
        </h1>
        <div class="flex flex-col sm:flex-row gap-1 sm:gap-4 text-blue-300 font-mono text-lg mb-2">
          <span>{{ dateFormat(conf.StartTime) }}</span>
          <span class="hidden sm:inline-block text-gray-500">→</span>
          <span>{{ dateFormat(conf.EndTime) }}</span>
        </div>
        <div class="flex gap-4 items-center text-gray-400 mb-1 text-sm">
          <span class="flex items-center gap-2">
            <span>Salle</span>
            <span
              class="inline-block px-3 py-0.5 bg-blue-900/70 border border-blue-700 rounded-lg font-bold text-blue-200 text-base">
              {{ conf.Room }}
            </span>
          </span>
          <span class="pl-4 border-l border-gray-600 font-mono text-xs">
            Conférencier :
            <span class="text-green-400 font-semibold">{{ conf.SpeakerName }}</span>
          </span>
        </div>
      </div>
      <div class="bg-gray-900 px-6 py-4 rounded-2xl shadow-lg border border-blue-900/60 text-center min-w-[120px]">
        <div class="text-xs text-gray-400 font-mono mb-1">Participants</div>
        <div class="text-yellow-300 text-3xl font-black">
          {{ participants.length }} <span class="text-gray-500 text-lg">/ 20</span>
        </div>
        <div v-if="isOwner" class="mt-2">
          <button @click="showUsers = true"
            class="text-blue-300 hover:text-white font-semibold text-xs underline transition">Voir la liste</button>
        </div>
      </div>
    </div>

    <div class="mb-8 text-gray-200 bg-gray-900/90 p-5 rounded-xl border border-blue-900/40 shadow-inner">
      <div class="font-bold text-lg mb-1 text-blue-200">Résumé :</div>
      <div class="leading-relaxed text-base">{{ conf.Description }}</div>
    </div>

    <div class="mb-1">
      <button v-if="!isOwner" @click="register" :disabled="registered || full"
        class="bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white font-semibold px-7 py-2 rounded-xl shadow-lg transition disabled:opacity-40 disabled:cursor-not-allowed">
        {{ registered ? 'Déjà inscrit' : (full ? 'Complet' : "S'inscrire à cette conférence") }}
      </button>
    </div>
    <div v-if="message" class="mt-4 text-green-400 font-semibold">{{ message }}</div>
    <div v-if="error" class="mt-4 text-red-400 font-semibold">{{ error }}</div>

    <div v-if="showUsers" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-70">
      <div class="bg-gray-900 border border-blue-700 rounded-2xl shadow-2xl p-7 w-full max-w-xs text-left relative">
        <button class="absolute top-3 right-4 text-2xl text-gray-500 hover:text-white"
          @click="showUsers = false">&times;</button>
        <div class="mb-3 text-lg text-blue-300 font-bold">Participants inscrits</div>
        <ul>
          <li v-for="u in participants" :key="u.ID || u.id"
            class="py-1.5 px-2 rounded text-gray-100 font-mono text-[15px] border-b border-blue-900 last:border-b-0">
            <span class="text-blue-400 font-bold">{{ u.Name || u.name || u.Email || u.email }}</span>
          </li>
        </ul>
        <div v-if="participants.length === 0" class="text-gray-400 text-sm mt-4">Aucun inscrit.</div>
      </div>
    </div>
  </div>
  <div v-else class="text-gray-400 text-lg text-center my-14 font-mono">Chargement…</div>
</template>

<script>
import axios from 'axios'
export default {
  data() {
    return {
      conf: null,
      registered: false,
      message: "",
      error: "",
      participants: [],
      showUsers: false
    }
  },
  computed: {
    isOwner() {
      try {
        const u = JSON.parse(localStorage.getItem('user'))
        return u && this.conf && (u.ID == this.conf.OrganizerID || u.id == this.conf.OrganizerID)
      } catch {
        return false
      }
    },
    full() {
      return this.participants.length >= 20
    }
  },
  async mounted() {
    const id = this.$route.params.id
    const token = localStorage.getItem('token')
    const res = await axios.get('http://localhost:8080/conferences', {
      headers: { Authorization: `Bearer ${token}` }
    })
    this.conf = res.data.find(c => (c.ID || c.id) == id)

    try {
      const pres = await axios.get(`http://localhost:8080/conferences/${id}/participants`, {
        headers: { Authorization: `Bearer ${token}` }
      })
      this.participants = pres.data || []
    } catch {
      this.participants = []
    }

    try {
      const myRes = await axios.get('http://localhost:8080/me/conferences', {
        headers: { Authorization: `Bearer ${token}` }
      })
      this.registered = !!myRes.data.find(c => (c.ID || c.id) == id)
    } catch { this.registered = false }
  },
  methods: {
    async register() {
      this.error = this.message = ""
      if (this.registered || this.full) return
      const token = localStorage.getItem('token')
      try {
        const id = this.$route.params.id
        await axios.post(`http://localhost:8080/conferences/${id}/join`, {}, {
          headers: { Authorization: `Bearer ${token}` }
        })
        this.registered = true
        this.message = "Inscription réussie !"
        const pres = await axios.get(`http://localhost:8080/conferences/${id}/participants`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        this.participants = pres.data || []
      } catch (e) {
        this.error = e.response?.data?.error || "Erreur lors de l'inscription"
      }
    },
    dateFormat(dt) {
      if (!dt) return ""
      const d = new Date(dt.replace(" ", "T"))
      return d.toLocaleString('fr-FR', { dateStyle: 'short', timeStyle: 'short' })
    }
  }
}
</script>
