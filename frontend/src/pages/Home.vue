<script setup>
import { onMounted, ref, computed } from 'vue'
import { useMediaStore } from '../stores/media'
import { useChannelStore } from '../stores/channels'
import MediaCard from '../components/MediaCard.vue'
import ChannelCard from '../components/ChannelCard.vue'
import { Film, Tv, Radio, Flame, Sparkles, AlertCircle, Cpu, Search, Layers } from 'lucide-vue-next'

const mediaStore = useMediaStore()
const channelStore = useChannelStore()

// Aba ativa: 'movie', 'tv', 'channel'
const activeTab = ref('movie')

onMounted(() => {
  mediaStore.fetchTrending('movie')
  mediaStore.fetchHealth()
})

function switchTab(type) {
  activeTab.value = type
  if (type === 'channel') {
    if (channelStore.channels.length === 0) {
      channelStore.fetchChannels()
      channelStore.fetchCategories()
    }
  } else {
    mediaStore.searchResults = []
    mediaStore.fetchTrending(type)
  }
}

const activeMediaList = computed(() => {
  if (mediaStore.searchResults.length > 0) {
    return mediaStore.searchResults
  }
  return mediaStore.trendingItems
})

const listTitle = computed(() => {
  if (activeTab.value === 'channel') {
    if (channelStore.selectedCategory && channelStore.selectedCategory !== 'Todos') {
      return `Canais Ao Vivo — ${channelStore.selectedCategory}`
    }
    return 'Canais de TV Ao Vivo'
  }
  if (mediaStore.searchResults.length > 0) {
    return 'Resultados da Busca'
  }
  return activeTab.value === 'movie' ? 'Filmes em Alta na Semana' : 'Séries em Alta na Semana'
})
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
        Filmes, séries do TMDB e canais de TV ao vivo com player unificado e cache em memória.
      </p>
    </section>

    <!-- Header da Lista & Seletor de Tipo de Mídia -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
      <div class="flex items-center gap-2">
        <Flame v-if="activeTab !== 'channel'" class="w-5 h-5 text-amber-400" />
        <Radio v-else class="w-5 h-5 text-rose-500 animate-pulse" />
        <h2 class="font-heading font-bold text-xl text-slate-100">
          {{ listTitle }}
        </h2>
      </div>

      <!-- Abas de Navegação -->
      <div class="flex items-center gap-1.5 p-1 rounded-xl bg-slate-900 border border-slate-800">
        <button
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-all"
          :class="activeTab === 'movie' && mediaStore.searchResults.length === 0 ? 'bg-sky-500 text-slate-950 shadow-md shadow-sky-500/20' : 'text-slate-400 hover:text-white'"
          @click="switchTab('movie')"
        >
          <Film class="w-3.5 h-3.5" />
          <span>Filmes</span>
        </button>

        <button
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-all"
          :class="activeTab === 'tv' && mediaStore.searchResults.length === 0 ? 'bg-sky-500 text-slate-950 shadow-md shadow-sky-500/20' : 'text-slate-400 hover:text-white'"
          @click="switchTab('tv')"
        >
          <Tv class="w-3.5 h-3.5" />
          <span>Séries</span>
        </button>

        <button
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-all"
          :class="activeTab === 'channel' ? 'bg-rose-500 text-white shadow-md shadow-rose-500/20' : 'text-slate-400 hover:text-white'"
          @click="switchTab('channel')"
        >
          <Radio class="w-3.5 h-3.5" />
          <span class="flex items-center gap-1">
            <span>Canais Ao Vivo</span>
            <span class="w-1.5 h-1.5 rounded-full bg-white animate-pulse"></span>
          </span>
        </button>
      </div>
    </div>

    <!-- SEÇÃO 1: CANAIS AO VIVO -->
    <div v-if="activeTab === 'channel'" class="space-y-6">
      <!-- Barra de Filtros de Canais (Busca e Categorias) -->
      <div class="flex flex-col md:flex-row items-stretch md:items-center justify-between gap-4 p-3 rounded-2xl bg-slate-900/60 border border-slate-800/80 backdrop-blur-md">
        <!-- Campo de Busca em Tempo Real no Catálogo de Canais -->
        <div class="relative flex-1 max-w-md">
          <Search class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
          <input
            v-model="channelStore.searchQuery"
            type="text"
            placeholder="Filtrar canal (ex: Globo, SporTV, ESPN, HBO)..."
            class="w-full bg-slate-950/80 text-xs sm:text-sm border border-slate-800 rounded-xl pl-10 pr-4 py-2 text-slate-100 placeholder-slate-400 focus:outline-none focus:border-rose-400 focus:ring-2 focus:ring-rose-400/20 transition-all"
          />
        </div>

        <!-- Pílulas de Categorias -->
        <div class="flex items-center gap-2 overflow-x-auto pb-1 md:pb-0 scrollbar-none">
          <button
            v-for="cat in channelStore.categories"
            :key="cat"
            class="px-3 py-1.5 rounded-xl text-xs font-medium whitespace-nowrap transition-all"
            :class="channelStore.selectedCategory === cat ? 'bg-sky-500 text-slate-950 font-bold shadow-md shadow-sky-500/20' : 'bg-slate-950/60 text-slate-400 hover:text-white border border-slate-800'"
            @click="channelStore.setCategory(cat)"
          >
            {{ cat }}
          </button>
        </div>
      </div>

      <!-- Aviso de Erro de Canais -->
      <div
        v-if="channelStore.errorMessage"
        class="p-4 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 flex items-start gap-3"
      >
        <AlertCircle class="w-5 h-5 flex-shrink-0 mt-0.5 text-rose-400" />
        <div>
          <p class="font-semibold text-sm">Falha ao carregar canais</p>
          <p class="text-xs text-rose-400/90 mt-0.5">{{ channelStore.errorMessage }}</p>
        </div>
      </div>

      <!-- Skeleton Loading de Canais -->
      <div v-if="channelStore.isLoading" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-5">
        <div v-for="i in 12" :key="i" class="animate-pulse rounded-2xl bg-slate-900 border border-slate-800 aspect-[16/10]"></div>
      </div>

      <!-- Grid de Canais -->
      <div
        v-else-if="channelStore.filteredChannels.length > 0"
        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-5"
      >
        <ChannelCard v-for="ch in channelStore.filteredChannels" :key="ch.id" :channel="ch" />
      </div>

      <!-- Estado Vazio para Canais -->
      <div v-else-if="!channelStore.errorMessage" class="text-center py-16 text-slate-500 text-sm">
        Nenhum canal encontrado para a pesquisa ou categoria selecionada.
      </div>
    </div>

    <!-- SEÇÃO 2: FILMES E SÉRIES (TMDB) -->
    <div v-else class="space-y-6">
      <!-- Mensagem de Erro / Aviso se TMDB Key não estiver presente -->
      <div
        v-if="mediaStore.errorMessage"
        class="p-4 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 flex items-start gap-3"
      >
        <AlertCircle class="w-5 h-5 flex-shrink-0 mt-0.5 text-rose-400" />
        <div>
          <p class="font-semibold text-sm">Aviso de Consulta</p>
          <p class="text-xs text-rose-400/90 mt-0.5">{{ mediaStore.errorMessage }}</p>
        </div>
      </div>

      <!-- Grid de Skeletons -->
      <div v-if="mediaStore.isLoading" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-5">
        <div v-for="i in 12" :key="i" class="animate-pulse rounded-2xl bg-slate-900 border border-slate-800 aspect-[2/3]"></div>
      </div>

      <!-- Grid de Filmes / Séries -->
      <div
        v-else-if="activeMediaList.length > 0"
        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-5"
      >
        <MediaCard v-for="item in activeMediaList" :key="item.id" :media="item" />
      </div>

      <!-- Estado Vazio -->
      <div v-else-if="!mediaStore.errorMessage" class="text-center py-16 text-slate-500 text-sm">
        Nenhum título encontrado para os filtros selecionados.
      </div>
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
