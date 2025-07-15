<template>
  <div class="max-w-xl mt-12 mx-auto bg-gray-800 rounded-2xl shadow-xl p-8">
    <h1 class="text-3xl font-extrabold text-blue-300 mb-6 tracking-tight">NWS Conférence</h1>

    <div v-if="user" class="flex items-center space-x-5">
      <img
        v-if="user.Picture || user.picture"
        :src="user.Picture || user.picture"
        class="rounded-full border-4 border-blue-500 shadow-lg"
        width="80"
        height="80"
        alt="Avatar"
      />
      <div>
        <p class="text-lg text-white mb-1">Vous êtes connecté en tant que <b class="text-blue-300">{{ user.Name || user.Name }}</b></p>
        <span class="text-gray-400 font-mono text-sm">{{ user.Email || user.Email }}</span>
      </div>
    </div>

    <div v-else class="text-gray-300 text-lg">
      <p>Vous n'êtes pas connecté.</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return { user: null }
  },
  mounted() {
    const userData = localStorage.getItem('user')
    if (userData) this.user = JSON.parse(userData)
  }
}
</script>
