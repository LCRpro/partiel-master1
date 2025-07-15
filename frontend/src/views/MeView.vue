<template>
  <div>
    <h2>Mon profil</h2>
    <div v-if="user">
      <img
        v-if="user.Picture || user.picture"
        :src="user.Picture || user.picture"
        width="100"
        class="rounded-full shadow border-2 border-gray-200 mb-3"
        alt="Avatar"
      />
      <p><b>ID (base) :</b> {{ user.ID || user.id }}</p>
      <p><b>GoogleID :</b> {{ user.GoogleID || user.google_id }}</p>
      <p><b>Email :</b> {{ user.Email || user.email }}</p>
      <p><b>Nom :</b> {{ user.Name || user.name }}</p>
    </div>
    <div v-else>
      <em>Chargement...</em>
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
