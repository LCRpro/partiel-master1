<template>
  <div v-if="conf" class="bg-gray-800 rounded-2xl shadow-2xl p-8 max-w-2xl mx-auto my-10 border border-blue-700">
    <div class="flex flex-col sm:flex-row items-start sm:items-center gap-6 mb-6">
      <div class="flex-1">
        <h1 class="text-3xl font-black text-white mb-1">{{ conf.Title }}</h1>
        <div class="text-blue-300 text-xl font-mono mb-2">
          {{ dateFormat(conf.StartTime) }} → {{ dateFormat(conf.EndTime) }}
        </div>
        <div class="flex gap-4 items-center text-gray-400 mb-2 text-sm">
          <span>Salle <span class="text-blue-400 font-bold text-lg">{{ conf.Room }}</span></span>
          <span class="pl-3 border-l border-gray-600 font-mono text-xs">Conférencier : <span class="text-green-400 font-semibold">{{ conf.SpeakerName }}</span></span>
        </div>
      </div>
      <div class="bg-gray-900 px-4 py-3 rounded-xl shadow border border-gray-700 text-center min-w-[120px]">
        <div class="text-xs text-gray-400 font-mono">Participants</div>
        <div class="text-yellow-300 text-2xl font-black">
          {{ participants.length }} / 20
        </div>
        <div v-if="isOwner" class="mt-2">
          <button
            @click="showUsers = true"
            class="text-blue-300 hover:text-white font-semibold text-xs underline transition"
          >Voir la liste</button>
        </div>
      </div>
    </div>

    <div class="mb-6 text-gray-200 bg-gray-900 p-4 rounded-xl border border-gray-700 shadow-inner">
      <div class="font-bold text-lg mb-1 text-blue-200">Résumé :</div>
      <div class="leading-relaxed">{{ conf.Description }}</div>
    </div>

    <button
      v-if="!isOwner"
      @click="register"
      :disabled="registered || full"
      class="bg-blue-600 hover:bg-blue-700 text-white font-semibold px-6 py-2 rounded-lg shadow transition disabled:opacity-50"
    >
      {{ registered ? 'Déjà inscrit' : (full ? 'Complet' : "S'inscrire à cette conférence") }}
    </button>
    <div v-if="message" class="mt-4 text-green-400">{{ message }}</div>
    <div v-if="error" class="mt-4 text-red-400">{{ error }}</div>

    <div v-if="showUsers" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-60">
      <div class="bg-gray-900 border border-blue-700 rounded-2xl shadow-2xl p-6 w-full max-w-xs text-left relative">
        <button class="absolute top-3 right-4 text-2xl text-gray-500 hover:text-white" @click="showUsers = false">&times;</button>
        <div class="mb-3 text-lg text-blue-300 font-bold">Participants inscrits</div>
        <ul>
          <li v-for="u in participants" :key="u.ID || u.id" class="py-1 px-2 rounded text-gray-200 font-mono text-sm border-b border-gray-700 last:border-b-0">
            {{ u.Name || u.name || u.Email || u.email }}
          </li>
        </ul>
        <div v-if="participants.length === 0" class="text-gray-400 text-sm mt-3">Aucun inscrit.</div>
      </div>
    </div>
  </div>
  <div v-else class="text-gray-300 text-lg text-center my-10">Chargement…</div>
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


