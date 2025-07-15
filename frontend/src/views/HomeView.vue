<template>
  <div class="max-w-xl mt-16 mx-auto bg-gradient-to-br from-gray-900 via-gray-800 to-blue-950 rounded-3xl shadow-2xl p-10 border border-blue-900/40">
    <div class="flex items-center justify-between mb-5">
      <h1 class="text-3xl font-extrabold text-blue-300 tracking-tight drop-shadow">NWS Conférence</h1>
      <span class="bg-blue-800/80 text-blue-200 px-4 py-1.5 rounded-full font-semibold text-xs shadow">
        Salon du 18/07 au 20/07
      </span>
    </div>
    <div v-if="user" class="flex items-center space-x-6 mt-5">
      <img
        v-if="user.Picture || user.picture"
        :src="user.Picture || user.picture"
        class="rounded-full border-4 border-blue-500 shadow-lg ring-2 ring-blue-900"
        width="85"
        height="85"
        alt="Avatar"
      />
      <div>
        <p class="text-lg text-white mb-1">
          <span class="opacity-60 font-medium">Connecté en tant que</span>
          <b class="text-blue-300 ml-1">{{ user.Name || user.name }}</b>
        </p>
        <span class="text-gray-400 font-mono text-sm">{{ user.Email || user.email }}</span>
      </div>
    </div>
    <div v-else class="text-gray-300 text-lg mt-8 text-center">
      <p>
        <span class="inline-block bg-gray-700/60 px-4 py-2 rounded-xl">
          Vous n'êtes pas connecté.
        </span>
      </p>
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
