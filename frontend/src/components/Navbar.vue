<script setup>
import { ref } from 'vue'
import { useMediaStore } from '../stores/media'
import { Play, Search, Activity } from 'lucide-vue-next'

const mediaStore = useMediaStore()
const searchQuery = ref('')

function handleSearch() {
  if (searchQuery.value.trim().length >= 2) {
    mediaStore.search(searchQuery.value.trim())
  }
}
</script>

<template>
  <header class="sticky top-0 z-40 w-full glass-panel border-b border-slate-800/80 px-6 py-4 flex items-center justify-between">
    <div class="flex items-center gap-3">
      <div class="w-9 h-9 rounded-xl bg-gradient-to-tr from-sky-400 to-indigo-500 flex items-center justify-center shadow-lg shadow-sky-500/20 text-slate-950 font-bold">
        <Play class="w-5 h-5 fill-slate-950 ml-0.5" />
      </div>
      <span class="font-heading font-extrabold text-xl tracking-tight bg-gradient-to-r from-sky-400 via-indigo-300 to-indigo-400 bg-clip-text text-transparent">
        Stream Hub
      </span>
    </div>

    <!-- Barra de busca central -->
    <div class="relative w-full max-w-md hidden md:block">
      <Search class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Buscar filmes, séries por título..."
        class="w-full bg-slate-900/90 text-sm border border-slate-700/60 rounded-xl pl-10 pr-4 py-2.5 text-slate-100 placeholder-slate-400 focus:outline-none focus:border-sky-400 focus:ring-2 focus:ring-sky-400/20 transition-all"
        @keydown.enter="handleSearch"
      />
    </div>

    <!-- Badge de Status do BFF -->
    <div class="flex items-center gap-2">
      <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-xs font-semibold text-emerald-400">
        <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
        <span>Go BFF Ativo</span>
      </div>
    </div>
  </header>
</template>
