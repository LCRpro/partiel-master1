<template>
  <div v-if="conf" class="bg-gray-800 rounded-2xl shadow-xl p-8 max-w-xl mx-auto">
    <div class="flex items-center space-x-6 mb-5">

      <div>
        <h1 class="text-2xl font-bold text-white">{{ conf.Title }}</h1>
        <div class="text-blue-300 text-lg mb-1">{{ dateFormat(conf.StartTime) }} → {{ dateFormat(conf.EndTime) }}</div>
        <div class="text-sm text-gray-300">Salle <b class="text-blue-400">{{ conf.Room }}</b></div>
        <div class="mt-2 text-gray-400 font-mono text-xs">Conférencier : {{ conf.SpeakerName }}</div>
      </div>
    </div>
    <div class="mb-6 text-gray-200">
      <div class="font-bold mb-1">Résumé :</div>
      <div>{{ conf.Description }}</div>
    </div>

    <button
      @click="register"
      :disabled="registered"
      class="bg-blue-600 hover:bg-blue-700 text-white font-semibold px-6 py-2 rounded-lg shadow transition disabled:opacity-50"
    >{{ registered ? 'Déjà inscrit' : "S'inscrire à cette conférence" }}</button>
    <div v-if="message" class="mt-4 text-green-400">{{ message }}</div>
    <div v-if="error" class="mt-4 text-red-400">{{ error }}</div>
  </div>
  <div v-else class="text-gray-300">Chargement…</div>
</template>

<script>
import axios from 'axios'
export default {
  data() {
    return {
      conf: null,
      registered: false,
      message: "",
      error: ""
    }
  },
  async mounted() {
    const id = this.$route.params.id
    const token = localStorage.getItem('token')
    const res = await axios.get('http://localhost:8080/conferences', {
      headers: { Authorization: `Bearer ${token}` }
    })
    this.conf = res.data.find(c => (c.ID || c.id) == id)

    const myRes = await axios.get('http://localhost:8080/me/conferences', {
      headers: { Authorization: `Bearer ${token}` }
    })
    this.registered = !!myRes.data.find(c => (c.ID || c.id) == id)
  },
  methods: {
    async register() {
      this.error = this.message = ""
      const token = localStorage.getItem('token')
      try {
        const id = this.$route.params.id
        await axios.post(`http://localhost:8080/conferences/${id}/join`, {}, {
          headers: { Authorization: `Bearer ${token}` }
        })
        this.registered = true
        this.message = "Inscription réussie !"
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
