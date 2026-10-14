<script setup>
import { onMounted, ref, computed } from 'vue'
import { useMediaStore } from '../stores/media'
import MediaCard from '../components/MediaCard.vue'
import { Film, Tv, Flame, Sparkles, AlertCircle, Cpu } from 'lucide-vue-next'

const mediaStore = useMediaStore()

onMounted(() => {
  mediaStore.fetchTrending('movie')
  mediaStore.fetchHealth()
})

const activeList = computed(() => {
  if (mediaStore.searchResults.length > 0) {
    return mediaStore.searchResults
  }
  return mediaStore.trendingItems
})

const listTitle = computed(() => {
  if (mediaStore.searchResults.length > 0) {
    return 'Resultados da Busca'
  }
  return mediaStore.currentMediaType === 'movie' ? 'Filmes em Alta na Semana' : 'Séries em Alta na Semana'
})

function switchCategory(type) {
  mediaStore.searchResults = []
  mediaStore.fetchTrending(type)
}
</script>

<template>
  <main class="max-w-7xl mx-auto px-6 py-8">
    <!-- Hero Section -->
    <section class="text-center mb-10">
      <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-sky-500/10 border border-sky-500/20 text-xs font-semibold text-sky-400 mb-4">
        <Sparkles class="w-3.5 h-3.5" />
        <span>BFF Autocontido em Go + Embedded SPA</span>
      </div>
      <h1 class="font-heading text-4xl sm:text-5xl font-extrabold tracking-tight mb-3">
        Catálogo & Streaming Local
      </h1>
      <p class="text-slate-400 max-w-xl mx-auto text-base">
        Proxy com cache em memória para metadados do TMDB e geração unificada de URLs do player.
      </p>
    </section>

    <!-- Header da Lista & Filtros -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
      <div class="flex items-center gap-2">
        <Flame class="w-5 h-5 text-amber-400" />
        <h2 class="font-heading font-bold text-xl text-slate-100">
          {{ listTitle }}
        </h2>
      </div>

      <div class="flex items-center gap-1.5 p-1 rounded-xl bg-slate-900 border border-slate-800">
        <button
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-all"
          :class="mediaStore.currentMediaType === 'movie' && mediaStore.searchResults.length === 0 ? 'bg-sky-500 text-slate-950 shadow-md shadow-sky-500/20' : 'text-slate-400 hover:text-white'"
          @click="switchCategory('movie')"
        >
          <Film class="w-3.5 h-3.5" />
          <span>Filmes</span>
        </button>

        <button
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-all"
          :class="mediaStore.currentMediaType === 'tv' && mediaStore.searchResults.length === 0 ? 'bg-sky-500 text-slate-950 shadow-md shadow-sky-500/20' : 'text-slate-400 hover:text-white'"
          @click="switchCategory('tv')"
        >
          <Tv class="w-3.5 h-3.5" />
          <span>Séries</span>
        </button>
      </div>
    </div>

    <!-- Mensagem de Erro / Aviso se TMDB Key não estiver presente -->
    <div
      v-if="mediaStore.errorMessage"
      class="mb-8 p-4 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 flex items-start gap-3"
    >
      <AlertCircle class="w-5 h-5 flex-shrink-0 mt-0.5 text-rose-400" />
      <div>
        <p class="font-semibold text-sm">Aviso de Consulta</p>
        <p class="text-xs text-rose-400/90 mt-0.5">{{ mediaStore.errorMessage }}</p>
      </div>
    </div>

    <!-- Grid de Cards -->
    <div v-if="mediaStore.isLoading" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-5">
      <div v-for="i in 12" :key="i" class="animate-pulse rounded-2xl bg-slate-900 border border-slate-800 aspect-[2/3]"></div>
    </div>

    <div
      v-else-if="activeList.length > 0"
      class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-5"
    >
      <MediaCard v-for="item in activeList" :key="item.id" :media="item" />
    </div>

    <div v-else-if="!mediaStore.errorMessage" class="text-center py-16 text-slate-500 text-sm">
      Nenhum título encontrado para os filtros selecionados.
    </div>

    <!-- Footer com Monitor de Performance (< 25 MB RAM) -->
    <footer v-if="mediaStore.healthStats" class="mt-16 pt-6 border-t border-slate-800/60 flex flex-wrap items-center justify-between gap-4 text-xs text-slate-400">
      <div class="flex items-center gap-2">
        <Cpu class="w-4 h-4 text-sky-400" />
        <span><strong>Memória Go Alocada:</strong> {{ mediaStore.healthStats.memory.alloc_mb.toFixed(2) }} MB (Heap) / {{ mediaStore.healthStats.memory.sys_mb.toFixed(2) }} MB (Sys)</span>
      </div>
      <div>
        <span><strong>Uptime do Servidor:</strong> {{ mediaStore.healthStats.uptime }}</span>
      </div>
      <div>
        <span><strong>Cache em Memória:</strong> {{ mediaStore.healthStats.cache.item_count }} itens armazenados</span>
      </div>
    </footer>
  </main>
</template>
