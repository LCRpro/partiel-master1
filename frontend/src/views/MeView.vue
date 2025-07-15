<template>
  <div class="flex flex-col items-center min-h-screen py-8 bg-gradient-to-b from-gray-900 to-gray-800">
    <div v-if="user" class="bg-gray-800 border border-gray-700 rounded-3xl shadow-2xl px-10 py-10 max-w-xl w-full">
      <div class="flex items-center gap-7 mb-6">
        <img v-if="user.Picture || user.picture" :src="user.Picture || user.picture"
          class="rounded-full shadow-lg border-4 border-blue-400 ring-4 ring-blue-800" width="90" height="90"
          alt="Avatar" />
        <div>
          <div class="flex items-center gap-2">
            <span class="text-2xl font-bold text-white">{{ user.Name || user.name }}</span>
          </div>
          <div class="text-blue-300 text-sm font-mono mt-1">{{ user.Email || user.email }}</div>
          <div class="mt-3 text-sm flex flex-wrap gap-2 items-center">
            <span class="font-bold text-gray-400">Rôles :</span>
            <template v-for="role in userRoles" :key="role">
              <span v-if="role === 'admin'"
                class="bg-gradient-to-r from-red-600 to-red-700 text-white font-bold px-2.5 py-1 rounded shadow text-xs uppercase tracking-wide">admin</span>
              <span v-else-if="role === 'conferencier'"
                class="bg-gradient-to-r from-green-700 to-green-500 text-white font-bold px-2.5 py-1 rounded shadow text-xs uppercase tracking-wide">conférencier</span>
              <span v-else
                class="bg-gradient-to-r from-blue-700 to-blue-500 text-white font-bold px-2 py-1 rounded shadow text-xs uppercase tracking-wide">{{
                role }}</span>
            </template>
          </div>
          <div v-if="userRoles.includes('conferencier') && (user.SpeakerName || user.speaker_name)"
            class="mt-2 text-green-300 text-xs font-mono flex items-center gap-1">
            <span class="font-bold">Nom de conférencier :</span>
            <span class="bg-green-900/70 px-2 py-0.5 rounded font-semibold tracking-wide">
              {{ user.SpeakerName || user.speaker_name }}
            </span>
          </div>
        </div>
      </div>
      <div class="grid grid-cols-2 gap-x-6 gap-y-2 mt-3 text-gray-400 text-[15px]">
        <div>Google ID :</div>
        <div class="font-mono text-blue-300 break-all">{{ user.GoogleID || user.google_id }}</div>
      </div>
      <div class="flex gap-4 mt-8">
        <button @click="openConfPopup"
          class="bg-gradient-to-r from-green-600 to-green-700 hover:from-green-700 hover:to-green-800 text-white px-5 py-2.5 rounded-xl font-semibold shadow transition">Devenir
          conférencier</button>
        <button @click="openPopup('admin')"
          class="bg-gradient-to-r from-red-600 to-red-700 hover:from-red-700 hover:to-red-800 text-white px-5 py-2.5 rounded-xl font-semibold shadow transition">Devenir
          admin</button>
      </div>
    </div>
    <div v-else class="text-gray-300 mt-12 text-lg">
      <em>Chargement…</em>
    </div>

    <div v-if="myConfs.length" class="mt-14 w-full max-w-3xl">
      <h2 class="text-xl font-bold text-blue-300 mb-3 border-b-2 border-blue-800 pb-2 tracking-wide">Programme
        personnalisé
      </h2>
      <div class="space-y-4">
        <div v-for="conf in myConfs" :key="conf.ID || conf.id"
          class="bg-gray-800/95 rounded-2xl px-7 py-5 shadow-md flex flex-col md:flex-row md:items-center gap-4 hover:shadow-blue-800/30 transition-shadow border border-blue-900/40">
          <div class="flex-1 min-w-0">
            <div class="flex gap-2 items-center mb-1">
              <span class="text-lg font-bold text-white">{{ conf.Title }}</span>
              <span v-if="isCreator(conf)"
                class="ml-2 bg-gradient-to-r from-green-700 to-green-500 text-white px-2 py-0.5 rounded font-semibold text-xs shadow">Créateur</span>
              <span v-else
                class="ml-2 bg-gradient-to-r from-gray-700 to-gray-500 text-gray-200 px-2 py-0.5 rounded font-semibold text-xs">Visiteur</span>
            </div>
            <div class="text-blue-400 text-sm font-mono">
              {{ dateFormat(conf.StartTime) }}
              <span class="text-gray-400 font-sans mx-1">•</span>
              Salle <b>{{ conf.Room }}</b>
            </div>
            <div class="text-gray-300 text-sm mt-1">Conférencier : <b>{{ conf.SpeakerName }}</b></div>
          </div>
          <router-link :to="`/conferences/${conf.ID || conf.id}`"
            class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-1.5 rounded-lg text-xs font-bold shadow transition w-full md:w-auto text-center">Détails</router-link>
        </div>
      </div>
    </div>

    <div v-if="popup" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div class="bg-gray-800 p-8 rounded-2xl shadow-xl w-full max-w-xs text-center">
        <div class="mb-4 text-lg text-white">
          Confirmer : ajouter le rôle <b class="text-blue-300">{{ popup }}</b> à votre profil ?
        </div>
        <div class="flex gap-4 justify-center">
          <button @click="updateRole"
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded shadow">Oui</button>
          <button @click="closePopup"
            class="bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded shadow">Annuler</button>
        </div>
      </div>
    </div>

    <div v-if="confPopup" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div class="bg-gray-800 p-8 rounded-2xl shadow-xl w-full max-w-xs text-center">
        <div class="mb-4 text-lg text-white">Choisissez votre nom de conférencier</div>
        <input v-model="speakerName" placeholder="Nom de conférencier"
          class="w-full px-3 py-2 mb-5 rounded border border-green-500 bg-gray-900 text-white focus:outline-none focus:ring-2 focus:ring-green-400 transition" />
        <div class="flex gap-4 justify-center">
          <button @click="becomeSpeaker"
            class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded shadow">Valider</button>
          <button @click="closeConfPopup"
            class="bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded shadow">Annuler</button>
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
      confPopup: false,
      speakerName: "",
      myConfs: []
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

    const [allRes, joinedRes] = await Promise.all([
      axios.get('http://localhost:8080/conferences', { headers: { Authorization: `Bearer ${token}` } }),
      axios.get('http://localhost:8080/me/conferences', { headers: { Authorization: `Bearer ${token}` } })
    ])

    const joined = joinedRes.data || []
    const all = allRes.data || []
    const userId = this.user.ID || this.user.id
    const speakerName = this.user.SpeakerName || this.user.speaker_name

    const mine = all.filter(
      conf =>
        (conf.OrganizerID === userId) ||
        (speakerName && conf.SpeakerName === speakerName)
    )

    const merged = [...joined]
    mine.forEach(c => {
      if (!merged.find(j => (j.ID || j.id) === (c.ID || c.id))) {
        merged.push(c)
      }
    })
    this.myConfs = merged
  },
  methods: {
    isCreator(conf) {
      if (!this.user) return false
      const userId = this.user.ID || this.user.id
      const speakerName = this.user.SpeakerName || this.user.speaker_name
      return (
        (conf.OrganizerID === userId) ||
        (speakerName && conf.SpeakerName === speakerName)
      )
    },
    openPopup(role) { this.popup = role },
    closePopup() { this.popup = null },
    openConfPopup() { this.confPopup = true; this.speakerName = "" },
    closeConfPopup() { this.confPopup = false },
    async becomeSpeaker() {
      if (!this.speakerName) return
      const token = localStorage.getItem('token')
      const res = await axios.post('http://localhost:8080/me/role', {
        role: "conferencier",
        speakerName: this.speakerName
      }, {
        headers: { Authorization: `Bearer ${token}` }
      })
      this.user = res.data.user
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
      this.closeConfPopup()
    },
    async updateRole() {
      const token = localStorage.getItem('token')
      const res = await axios.post('http://localhost:8080/me/role', { role: this.popup }, {
        headers: { Authorization: `Bearer ${token}` }
      })
      this.user = res.data.user
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
      this.closePopup()
    },
    dateFormat(dt) {
      if (!dt) return ""
      const d = new Date(dt.replace(" ", "T"))
      return d.toLocaleString('fr-FR', { dateStyle: 'short', timeStyle: 'short' })
    }
  }
}
</script>
