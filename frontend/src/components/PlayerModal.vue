<script setup>
import { computed } from 'vue'
import { usePlayerStore } from '../stores/player'
import { X, Tv, Film, Loader2, AlertTriangle, ExternalLink } from 'lucide-vue-next'

const playerStore = usePlayerStore()

const isTv = computed(() => {
  return playerStore.currentMedia?.media_type === 'tv'
})

const unreleasedInfo = computed(() => {
  const dateStr = playerStore.currentMedia?.release_date || playerStore.currentMedia?.first_air_date
  if (!dateStr) return null
  const releaseDate = new Date(dateStr)
  if (releaseDate > new Date()) {
    return new Intl.DateTimeFormat('pt-BR', { dateStyle: 'long' }).format(releaseDate)
  }
  return null
})
</script>

<template>
  <div
    v-if="playerStore.isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-6 bg-slate-950/85 backdrop-blur-xl transition-all"
    @click.self="playerStore.closePlayer"
  >
    <div class="relative w-full max-w-4xl bg-slate-900 border border-slate-700/80 rounded-2xl shadow-2xl overflow-hidden flex flex-col">
      <!-- Modal Header -->
      <div class="px-5 py-4 flex items-center justify-between border-b border-slate-800">
        <div class="flex items-center gap-2.5">
          <component :is="isTv ? Tv : Film" class="w-5 h-5 text-sky-400" />
          <h2 class="font-heading font-bold text-lg text-slate-100 truncate max-w-lg">
            {{ playerStore.currentTitle }}
          </h2>
        </div>
        <div class="flex items-center gap-2">
          <a
            v-if="playerStore.embedUrl"
            :href="playerStore.embedUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="p-1.5 rounded-lg text-slate-400 hover:text-sky-300 hover:bg-slate-800 transition-colors"
            title="Abrir player em nova aba"
          >
            <ExternalLink class="w-4 h-4" />
          </a>
          <button
            class="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
            @click="playerStore.closePlayer"
          >
            <X class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- Aviso para títulos não lançados -->
      <div v-if="unreleasedInfo" class="px-5 py-3 bg-amber-500/10 border-b border-amber-500/20 text-amber-300 text-xs flex items-center gap-2.5">
        <AlertTriangle class="w-4 h-4 text-amber-400 flex-shrink-0" />
        <span>
          <strong>Lançamento Futuro:</strong> Estreia prevista para <strong>{{ unreleasedInfo }}</strong>. O servidor do player retornará <strong>"404 | Not Found"</strong> pois este título ainda não foi lançado nos cinemas ou no streaming.
        </span>
      </div>

      <!-- Iframe Container -->
      <div class="relative w-full aspect-video bg-black flex items-center justify-center">
        <div v-if="playerStore.isLoading" class="absolute inset-0 flex items-center justify-center bg-black/60 z-10">
          <Loader2 class="w-8 h-8 text-sky-400 animate-spin" />
        </div>

        <iframe
          v-if="playerStore.embedUrl"
          :src="playerStore.embedUrl"
          class="w-full h-full border-0"
          allowfullscreen
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
        ></iframe>

        <div v-else-if="playerStore.errorMessage" class="text-rose-400 text-sm px-6 text-center">
          {{ playerStore.errorMessage }}
        </div>
      </div>

      <!-- Episode & Season Selector for TV series -->
      <div v-if="isTv" class="px-5 py-3.5 bg-slate-950/60 border-t border-slate-800 flex flex-wrap items-center gap-4 text-sm">
        <div class="flex items-center gap-2">
          <label class="text-slate-400 font-medium">Temporada:</label>
          <input
            v-model.number="playerStore.season"
            type="number"
            min="1"
            class="w-16 bg-slate-800 border border-slate-700 rounded-lg px-2.5 py-1 text-center font-bold text-sky-400 focus:outline-none focus:border-sky-400"
            @change="playerStore.updateEmbedUrl"
          />
        </div>

        <div class="flex items-center gap-2">
          <label class="text-slate-400 font-medium">Episódio:</label>
          <input
            v-model.number="playerStore.episode"
            type="number"
            min="1"
            class="w-16 bg-slate-800 border border-slate-700 rounded-lg px-2.5 py-1 text-center font-bold text-sky-400 focus:outline-none focus:border-sky-400"
            @change="playerStore.updateEmbedUrl"
          />
        </div>

        <button
          class="ml-auto px-4 py-1.5 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 font-semibold text-xs transition-colors"
          @click="playerStore.updateEmbedUrl"
        >
          Carregar Episódio
        </button>
      </div>
    </div>
  </div>
</template>
