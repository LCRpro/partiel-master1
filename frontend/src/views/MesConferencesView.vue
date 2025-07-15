<template>
  <div>
    <h1 class="text-2xl font-bold text-blue-300 mb-8">Mes conférences</h1>

    <div v-if="loading" class="text-gray-400">Chargement…</div>

    <div v-else-if="!isSpeaker" class="text-red-500 font-bold text-lg">
      Accès interdit. Seuls les conférenciers peuvent accéder à cette page.
    </div>

    <div v-else>
      <div v-if="myConfs.length === 0" class="text-gray-400 text-lg">
        Aucune conférence créée.
      </div>
      <table v-else class="w-full bg-gray-800 rounded-xl overflow-hidden shadow-lg">
        <thead>
          <tr class="bg-gray-700 text-gray-100">
            <th class="px-4 py-3 text-left">Titre</th>
            <th class="px-3 py-3 text-left">Jour</th>
            <th class="px-3 py-3 text-left">Salle</th>
            <th class="px-3 py-3 text-left">Participants</th>
            <th class="px-3 py-3 text-left">Remplissage</th>
            <th class="px-3 py-3">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="conf in myConfs" :key="conf.ID || conf.id" class="hover:bg-gray-700 transition">
            <td class="px-4 py-3 font-semibold text-white">{{ conf.Title }}</td>
            <td class="px-3 py-3 text-blue-300 font-mono">{{ dateFormat(conf.StartTime) }}</td>
            <td class="px-3 py-3 text-blue-400 font-bold">{{ conf.Room }}</td>
            <td class="px-3 py-3 text-yellow-400 font-mono">{{ conf.ParticipantCount || 0 }} / 20</td>
            <td class="px-3 py-3">
              <div class="flex items-center gap-2">
                <div class="w-32 h-2 bg-gray-600 rounded">
                  <div
                    class="h-2 rounded bg-green-400"
                    :style="{ width: percent(conf) + '%' }"
                  ></div>
                </div>
                <span class="text-xs text-gray-300">{{ percent(conf) }}%</span>
              </div>
            </td>
            <td class="px-3 py-3 flex gap-2">
              <router-link
                :to="`/conferences/${conf.ID || conf.id}`"
                class="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded text-xs"
              >Détails</router-link>
              <button
                @click="openEdit(conf)"
                class="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-1 rounded text-xs"
              >Modifier</button>
              <button
                @click="confirmDelete(conf)"
                class="bg-red-600 hover:bg-red-700 text-white px-3 py-1 rounded text-xs"
              >Supprimer</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="modalEdit" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div class="bg-gray-900 p-8 rounded-xl w-full max-w-lg shadow-2xl relative">
        <button class="absolute right-4 top-3 text-gray-400 text-xl" @click="closeEdit">&times;</button>
        <h2 class="text-xl font-bold text-blue-400 mb-4">Modifier la conférence</h2>
        <form @submit.prevent="saveEdit" class="space-y-4">
          <input
            v-model="editForm.Title"
            placeholder="Titre"
            required
            class="w-full px-3 py-2 rounded border border-blue-700 bg-gray-800 text-white"
          />
          <textarea
            v-model="editForm.Description"
            placeholder="Description"
            required
            class="w-full px-3 py-2 rounded border border-blue-700 bg-gray-800 text-white"
          ></textarea>
          <button class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded w-full mt-2" type="submit">
            Enregistrer
          </button>
        </form>
        <div v-if="editError" class="text-red-500 mt-2">{{ editError }}</div>
      </div>
    </div>

    <div v-if="modalDelete" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-60">
      <div class="bg-gray-900 p-8 rounded-xl shadow-xl w-full max-w-xs text-center">
        <div class="text-lg text-white mb-6">
          Supprimer la conférence<br>
          <span class="text-blue-300 font-bold">{{ deleteTarget?.Title }}</span> ?
        </div>
        <div class="flex gap-4 justify-center">
          <button @click="doDelete" class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded shadow">Oui</button>
          <button @click="closeDelete" class="bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded shadow">Annuler</button>
        </div>
        <div v-if="deleteError" class="text-red-400 mt-3">{{ deleteError }}</div>
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


