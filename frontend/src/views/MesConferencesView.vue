<template>
  <div class="py-10 px-4 max-w-5xl mx-auto min-h-screen">
    <h1 class="text-3xl font-black text-blue-400 mb-10 tracking-tight drop-shadow">Mes conférences</h1>

    <div v-if="loading" class="text-gray-400 text-lg flex items-center gap-2">
      <svg class="animate-spin w-5 h-5 text-blue-400" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z" />
      </svg>
      Chargement…
    </div>

    <div v-else-if="!isSpeaker" class="text-red-400 font-bold text-lg mt-10 text-center">
      <span class="inline-block bg-red-900/40 px-4 py-3 rounded-xl">
        Accès interdit. Seuls les conférenciers peuvent accéder à cette page.
      </span>
    </div>

    <div v-else>
      <div v-if="myConfs.length === 0" class="text-gray-400 text-lg text-center">
        <span class="inline-block bg-gray-800/80 px-4 py-4 rounded-2xl shadow">
          Aucune conférence créée.
        </span>
      </div>
      <div v-else class="overflow-x-auto">
        <table
          class="w-full rounded-2xl overflow-hidden shadow-2xl bg-gradient-to-br from-gray-800 via-gray-900 to-gray-800">
          <thead>
            <tr class="bg-gradient-to-r from-blue-900 via-gray-800 to-blue-800 text-blue-100">
              <th class="px-5 py-4 text-left font-bold">Titre</th>
              <th class="px-4 py-4 text-left font-semibold">Jour</th>
              <th class="px-4 py-4 text-left font-semibold">Salle</th>
              <th class="px-4 py-4 text-left font-semibold">Participants</th>
              <th class="px-4 py-4 text-left font-semibold">Remplissage</th>
              <th class="px-4 py-4 font-semibold">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="conf in myConfs" :key="conf.ID || conf.id"
              class="hover:bg-blue-900/30 transition-colors border-b border-gray-800/70 last:border-0">
              <td class="px-5 py-4 text-lg text-white font-semibold max-w-xs truncate">
                {{ conf.Title }}
              </td>
              <td class="px-4 py-4 text-blue-200 font-mono">{{ dateFormat(conf.StartTime) }}</td>
              <td class="px-4 py-4 text-blue-400 font-bold">Salle {{ conf.Room }}</td>
              <td class="px-4 py-4 text-yellow-400 font-mono text-base">{{ conf.ParticipantCount || 0 }} / 20</td>
              <td class="px-4 py-4">
                <div class="flex items-center gap-3 min-w-[160px]">
                  <div class="relative w-36 h-2 bg-gray-600 rounded">
                    <div class="h-2 rounded absolute left-0 top-0 transition-all" :class="[
                      percent(conf) >= 100 ? 'bg-green-500' :
                        percent(conf) >= 75 ? 'bg-lime-400' :
                          percent(conf) >= 50 ? 'bg-yellow-400' :
                            percent(conf) >= 25 ? 'bg-orange-400' :
                              'bg-red-500'
                    ]" :style="{ width: percent(conf) + '%' }"></div>
                  </div>
                  <span :class="[
                    'ml-2 text-xs font-bold rounded px-2 py-0.5',
                    percent(conf) >= 100 ? 'bg-green-700 text-white' :
                      percent(conf) >= 75 ? 'bg-lime-700 text-white' :
                        percent(conf) >= 50 ? 'bg-yellow-700 text-yellow-100' :
                          percent(conf) >= 25 ? 'bg-orange-700 text-white' :
                            'bg-red-700 text-white'
                  ]">
                    {{ percent(conf) }}%
                  </span>
                </div>
              </td>
              <td class="px-4 py-4 flex gap-2 flex-wrap items-center">
                <router-link :to="`/conferences/${conf.ID || conf.id}`"
                  class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-1.5 rounded shadow text-xs font-semibold transition"
                  title="Détails">
                  <svg class="inline w-4 h-4 mr-1 -mt-0.5" fill="none" viewBox="0 0 24 24">
                    <path stroke="currentColor" stroke-width="2" d="M15 12l-4-4m0 0l4-4m-4 4h14" />
                  </svg>
                  Détails
                </router-link>
                <button @click="openEdit(conf)"
                  class="bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-1.5 rounded shadow text-xs font-semibold transition"
                  title="Modifier">
                  <svg class="inline w-4 h-4 mr-1 -mt-0.5" fill="none" viewBox="0 0 24 24">
                    <path stroke="currentColor" stroke-width="2"
                      d="M15.232 5.232l3.536 3.536M7.5 18.5l2.964-.741c.295-.074.574-.236.785-.447l8.486-8.486a1.5 1.5 0 000-2.121l-3.536-3.535a1.5 1.5 0 00-2.121 0l-8.486 8.485a1.5 1.5 0 00-.447.786L5.5 16.5a1 1 0 001 1z" />
                  </svg>
                  Modifier
                </button>
                <button @click="confirmDelete(conf)"
                  class="bg-red-600 hover:bg-red-700 text-white px-4 py-1.5 rounded shadow text-xs font-semibold transition"
                  title="Supprimer">
                  <svg class="inline w-4 h-4 mr-1 -mt-0.5" fill="none" viewBox="0 0 24 24">
                    <path stroke="currentColor" stroke-width="2"
                      d="M6 7h12M9 7v10a2 2 0 002 2h2a2 2 0 002-2V7m5 0V5a2 2 0 00-2-2H7a2 2 0 00-2 2v2" />
                  </svg>
                  Supprimer
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="modalEdit" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-70">
      <div class="bg-gray-900 p-10 rounded-2xl w-full max-w-lg shadow-2xl relative">
        <button class="absolute right-5 top-4 text-gray-500 text-2xl hover:text-blue-400"
          @click="closeEdit">&times;</button>
        <h2 class="text-2xl font-bold text-blue-400 mb-6">Modifier la conférence</h2>
        <form @submit.prevent="saveEdit" class="space-y-6">
          <input v-model="editForm.Title" placeholder="Titre" required
            class="w-full px-4 py-2 rounded-xl border border-blue-700 bg-gray-800 text-white font-semibold focus:outline-none focus:ring-2 focus:ring-blue-400" />
          <textarea v-model="editForm.Description" placeholder="Description" required
            class="w-full px-4 py-2 rounded-xl border border-blue-700 bg-gray-800 text-white font-medium focus:outline-none focus:ring-2 focus:ring-blue-400"></textarea>
          <button
            class="bg-green-600 hover:bg-green-700 text-white px-6 py-2 rounded-xl w-full text-lg font-bold mt-2 transition">
            Enregistrer
          </button>
        </form>
        <div v-if="editError" class="text-red-500 mt-4">{{ editError }}</div>
      </div>
    </div>

    <div v-if="modalDelete" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-70">
      <div class="bg-gray-900 p-8 rounded-2xl shadow-2xl w-full max-w-xs text-center border border-red-900/60">
        <div class="text-lg text-white mb-6">
          Supprimer la conférence<br>
          <span class="text-blue-300 font-bold">{{ deleteTarget?.Title }}</span> ?
        </div>
        <div class="flex gap-4 justify-center mb-2">
          <button @click="doDelete"
            class="bg-red-600 hover:bg-red-700 text-white px-6 py-2 rounded-lg font-bold shadow">Oui</button>
          <button @click="closeDelete"
            class="bg-gray-600 hover:bg-gray-700 text-white px-6 py-2 rounded-lg shadow">Annuler</button>
        </div>
        <div v-if="deleteError" class="text-red-400 mt-2">{{ deleteError }}</div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  data() {
    return {
      loading: true,
      isSpeaker: false,
      myConfs: [],
      user: null,
      modalEdit: false,
      editForm: { ID: null, Title: "", Description: "" },
      editError: "",
      modalDelete: false,
      deleteTarget: null,
      deleteError: ""
    }
  },
  async mounted() {
    const userStr = localStorage.getItem('user')
    if (!userStr) {
      this.loading = false
      return
    }
    this.user = JSON.parse(userStr)
    try {
      const roles = JSON.parse(this.user.Roles)
      this.isSpeaker = roles.includes("conferencier")
    } catch {
      this.isSpeaker = false
    }
    if (!this.isSpeaker) {
      this.loading = false
      return
    }
    await this.loadConfs()
    this.loading = false
  },
  methods: {
    async loadConfs() {
      const token = localStorage.getItem('token')
      const res = await axios.get('http://localhost:8080/conferences', {
        headers: { Authorization: `Bearer ${token}` }
      })
      const myId = this.user.ID || this.user.id
      this.myConfs = (res.data || []).filter(c => c.OrganizerID === myId)
    },
    dateFormat(dt) {
      if (!dt) return ""
      const d = new Date(dt.replace(" ", "T"))
      return d.toLocaleString('fr-FR', { dateStyle: 'short', timeStyle: 'short' })
    },
    percent(conf) {
      const count = conf.ParticipantCount || 0
      return Math.round((count / 20) * 100)
    },
    openEdit(conf) {
      this.editError = ""
      this.editForm = {
        ID: conf.ID || conf.id,
        Title: conf.Title,
        Description: conf.Description
      }
      this.modalEdit = true
    },
    closeEdit() {
      this.modalEdit = false
    },
    async saveEdit() {
      this.editError = ""
      const token = localStorage.getItem('token')
      try {
        await axios.patch(`http://localhost:8080/conferences/${this.editForm.ID}`, {
          Title: this.editForm.Title,
          Description: this.editForm.Description
        }, {
          headers: { Authorization: `Bearer ${token}` }
        })
        this.modalEdit = false
        await this.loadConfs()
      } catch (e) {
        this.editError = e.response?.data?.error || "Erreur modification"
      }
    },
    confirmDelete(conf) {
      this.deleteTarget = conf
      this.deleteError = ""
      this.modalDelete = true
    },
    closeDelete() {
      this.modalDelete = false
    },
    async doDelete() {
      this.deleteError = ""
      const token = localStorage.getItem('token')
      try {
        await axios.delete(`http://localhost:8080/conferences/${this.deleteTarget.ID || this.deleteTarget.id}`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        this.modalDelete = false
        await this.loadConfs()
      } catch (e) {
        this.deleteError = e.response?.data?.error || "Erreur suppression"
      }
    }
  }
}
</script>
