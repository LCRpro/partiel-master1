<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-2xl font-bold text-blue-300">Liste des conférences</h1>
      <button
        v-if="isSpeaker"
        @click="openModal"
        class="bg-blue-500 hover:bg-blue-600 text-white font-semibold px-4 py-2 rounded-lg shadow transition"
      >
        + Nouvelle conférence
      </button>
    </div>

    <div class="flex items-center mb-8 gap-2">
      <label class="text-gray-400 font-medium mr-4">Filtrer par jour :</label>
      <button
        v-for="d in days"
        :key="d.value"
        @click="selectedDay = d.value"
        :class="[
          'px-4 py-1 rounded font-mono transition border',
          selectedDay === d.value
            ? 'bg-blue-600 text-white border-blue-800'
            : 'bg-gray-800 text-gray-300 border-gray-700 hover:bg-blue-900 hover:text-white'
        ]"
      >{{ d.label }}</button>
    </div>

    <div v-if="loading" class="text-gray-300">Chargement…</div>
    <div v-else>
      <div v-if="Array.isArray(filteredConfs) && filteredConfs.length === 0" class="text-gray-300">
        Aucune conférence programmée pour ce jour.
      </div>
      <div class="grid md:grid-cols-2 gap-8">
        <div v-for="room in 10" :key="room" class="mb-8">
          <h2 class="text-blue-400 text-xl font-bold mb-2">Salle {{ room }}</h2>
          <table class="w-full bg-gray-800 rounded-lg overflow-hidden shadow mb-4">
            <thead>
              <tr class="bg-gray-700 text-gray-200">
                <th class="px-3 py-2">Titre</th>
                <th class="px-3 py-2">Début</th>
                <th class="px-3 py-2">Fin</th>
                <th class="px-3 py-2">Conférencier</th>
                <th class="px-3 py-2">Participants</th>
                <th class="px-3 py-2">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="conf in filteredConfs.filter(c => c.Room == room)"
                :key="conf.ID || conf.id"
                :class="rowClass(conf)"
              >
                <td class="px-3 py-2 font-semibold">
                  <span :class="{ 'line-through text-gray-400': isPast(conf) }">{{ conf.Title }}</span>
                </td>
                <td class="px-3 py-2 font-mono">
                  <span :class="{ 'line-through text-gray-400': isPast(conf) }">{{ dateFormat(conf.StartTime) }}</span>
                </td>
                <td class="px-3 py-2 font-mono">
                  <span :class="{ 'line-through text-gray-400': isPast(conf) }">{{ dateFormat(conf.EndTime) }}</span>
                </td>
                <td class="px-3 py-2">
                  <span :class="{ 'line-through text-gray-400': isPast(conf) }">{{ conf.SpeakerName }}</span>
                </td>
                <td class="px-3 py-2 text-yellow-400 font-mono">
                  <span :class="{ 'line-through text-gray-400': isPast(conf) }">{{ conf.ParticipantCount || 0 }} / 20</span>
                </td>
                <td class="px-3 py-2">
                  <router-link
                    :to="`/conferences/${conf.ID || conf.id}`"
                    class="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded text-xs"
                  >Détails</router-link>
                </td>
              </tr>
              <tr v-if="filteredConfs.filter(c => c.Room == room).length === 0">
                <td colspan="6" class="px-3 py-3 text-gray-500 text-center">
                  Aucune conférence dans cette salle pour ce jour.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <div v-if="modal" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-40">
      <div class="bg-gray-900 p-8 rounded-xl w-full max-w-lg shadow-2xl relative">
        <button class="absolute right-2 top-2 text-gray-400 text-xl" @click="closeModal">&times;</button>
        <h2 class="text-xl font-bold text-blue-400 mb-4">Nouvelle conférence</h2>
        <form @submit.prevent="createConference" class="space-y-4">
          <input
            v-model="form.Title"
            placeholder="Titre"
            required
            class="w-full px-3 py-2 rounded border border-blue-700 bg-gray-800 text-white focus:outline-none focus:ring-2 focus:ring-blue-400 transition"
          />
          <textarea
            v-model="form.Description"
            placeholder="Description"
            required
            class="w-full px-3 py-2 rounded border border-blue-700 bg-gray-800 text-white focus:outline-none focus:ring-2 focus:ring-blue-400 transition"
          ></textarea>
          <div class="flex gap-2 items-center">
            <label class="text-gray-300">Jour :</label>
            <select v-model="selectedDay" class="px-2 py-2 rounded border border-green-700 bg-gray-800 text-white focus:outline-none focus:ring-2 focus:ring-green-400 transition" required>
              <option v-for="(day, idx) in days" :key="idx" :value="day.value">{{ day.label }}</option>
            </select>
            <label class="ml-4 text-gray-300">Salle :</label>
            <select v-model="form.Room" class="px-2 py-2 rounded border border-green-700 bg-gray-800 text-white focus:outline-none focus:ring-2 focus:ring-green-400 transition" required>
              <option v-for="n in 10" :key="n" :value="n">{{ n }}</option>
            </select>
          </div>
          <div>
            <div class="mb-1 text-gray-300">Créneaux horaires disponibles (1h)</div>
            <div class="grid grid-cols-3 gap-2">
              <button
                v-for="(slot, idx) in availableSlots"
                :key="slot.value"
                type="button"
                :class="[
                  'px-3 py-2 rounded border text-sm font-mono transition',
                  slot.disabled
                    ? 'bg-gray-800 border-gray-800 text-gray-500 cursor-not-allowed'
                    : selectedSlots.includes(idx)
                      ? 'bg-blue-500 border-blue-700 text-white font-bold'
                      : 'bg-gray-700 border-gray-600 text-gray-200 hover:bg-blue-800 hover:text-white'
                ]"
                :disabled="slot.disabled"
                @click="toggleSlot(idx)"
              >
                {{ slot.label }}
              </button>
            </div>
          </div>
          <div>
            <label class="text-gray-400">Organisateur :</label>
            <span class="ml-2 text-blue-400 font-semibold">{{ connectedName }}</span>
          </div>
          <button class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded w-full mt-2" type="submit">Créer</button>
        </form>
        <div v-if="error" class="text-red-500 mt-2">{{ error }}</div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  data() {
    return {
      confs: [],
      loading: true,
      modal: false,
      error: null,
      form: {
        Title: "",
        Description: "",
        Room: 1,
        StartTime: "",
        EndTime: ""
      },
      days: [
        { label: "Vendredi 18/07", value: "2025-07-18" },
        { label: "Samedi 19/07", value: "2025-07-19" },
        { label: "Dimanche 20/07", value: "2025-07-20" },
      ],
      selectedDay: "2025-07-18",
      selectedSlots: []
    }
  },
  computed: {
    isSpeaker() {
      const u = localStorage.getItem('user')
      if (!u) return false
      try {
        const roles = JSON.parse(JSON.parse(u).Roles)
        return roles.includes("conferencier")
      } catch { return false }
    },
    connectedName() {
      const u = localStorage.getItem('user')
      if (!u) return ''
      try {
        return JSON.parse(u).SpeakerName || JSON.parse(u).Name || JSON.parse(u).name
      } catch { return '' }
    },
    filteredConfs() {
      return this.confs.filter(c => c.StartTime && c.StartTime.substr(0, 10) === this.selectedDay)
    },
    availableSlots() {
      const slots = []
      const room = Number(this.form.Room)
      const dayConfs = this.filteredConfs.filter(
        c => c.Room == room
      )
      const maxed = dayConfs.length >= 10
      for (let h = 9; h < 19; h++) {
        const start = `${this.selectedDay}T${String(h).padStart(2, "0")}:00`
        const end = `${this.selectedDay}T${String(h + 1).padStart(2, "0")}:00`
        let disabled = maxed
        for (const conf of dayConfs) {
          if (
            (start.replace("T", " ") >= conf.StartTime && start.replace("T", " ") < conf.EndTime) ||
            (end.replace("T", " ") > conf.StartTime && end.replace("T", " ") <= conf.EndTime)
          ) {
            disabled = true
            break
          }
        }
        slots.push({
          label: `${h}h - ${h + 1}h`,
          value: { start, end },
          disabled
        })
      }
      return slots
    }
  },
  watch: {
    selectedDay() {
      this.selectedSlots = []
    },
    'form.Room'() {
      this.selectedSlots = []
    }
  },
  methods: {
    isPast(conf) {
      return new Date(conf.EndTime.replace(" ", "T")) < new Date()
    },
    rowClass(conf) {
      return this.isPast(conf)
        ? "bg-gray-900 opacity-60"
        : "hover:bg-gray-700 transition"
    },
    async load() {
      this.loading = true
      const token = localStorage.getItem('token')
      const res = await axios.get('http://localhost:8080/conferences', {
        headers: { Authorization: `Bearer ${token}` }
      })
      this.confs = res.data || []
      this.loading = false
    },
    openModal() {
      this.modal = true
      this.error = null
      this.form = {
        Title: "", Description: "",
        Room: 1, StartTime: "", EndTime: ""
      }
      this.selectedSlots = []
    },
    closeModal() { this.modal = false },
    toggleSlot(idx) {
      if (this.availableSlots[idx].disabled) return
      this.selectedSlots = [idx]
    },
    async createConference() {
      this.error = null
      if (this.selectedSlots.length !== 1) {
        this.error = "Sélectionne un créneau d'1h"
        return
      }
      const slot = this.availableSlots[this.selectedSlots[0]]
      const startSlot = slot.value.start.replace("T", " ")
      const endSlot = slot.value.end.replace("T", " ")
      if (startSlot < "2025-07-18 09:00" || endSlot > "2025-07-20 19:00") {
        this.error = "Les conférences doivent avoir lieu entre le 18/07/2025 9h et le 20/07/2025 19h."
        return
      }
      const toDb = {
        ...this.form,
        StartTime: startSlot,
        EndTime: endSlot,
        Room: Number(this.form.Room)
      }
      try {
        const token = localStorage.getItem('token')
        await axios.post('http://localhost:8080/conferences', toDb, {
          headers: { Authorization: `Bearer ${token}` }
        })
        this.closeModal()
        this.load()
      } catch (e) {
        this.error = e.response?.data?.error || "Erreur création"
      }
    },
    dateFormat(dt) {
      if (!dt) return ""
      const d = new Date(dt.replace(" ", "T"))
      return d.toLocaleString('fr-FR', { dateStyle: 'short', timeStyle: 'short' })
    }
  },
  mounted() { this.load() }
}
</script>
