<script setup>
import { computed } from 'vue'
import { Star, Play } from 'lucide-vue-next'
import { usePlayerStore } from '../stores/player'

const props = defineProps({
  media: {
    type: Object,
    required: true
  }
})

const playerStore = usePlayerStore()

const title = computed(() => props.media.title || props.media.name || 'Sem título')
const releaseYear = computed(() => {
  const date = props.media.release_date || props.media.first_air_date
  return date ? date.substring(0, 4) : ''
})

const posterUrl = computed(() => {
  if (props.media.poster_path) {
    return `https://image.tmdb.org/t/p/w342${props.media.poster_path}`
  }
  return 'data:image/svg+xml;charset=UTF-8,%3Csvg%20width%3D%22200%22%20height%3D%22300%22%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%3E%3Crect%20width%3D%22100%25%22%20height%3D%22100%25%22%20fill%3D%22%231e293b%22%2F%3E%3Ctext%20x%3D%2250%25%22%20y%3D%2250%25%22%20fill%3D%22%2394a3b8%22%20dominant-baseline%3D%22middle%22%20text-anchor%3D%22middle%22%3ESem%20Poster%3C%2Ftext%3E%3C%2Fsvg%3E'
})

const rating = computed(() => {
  return props.media.vote_average ? props.media.vote_average.toFixed(1) : 'N/A'
})

function handleClick() {
  playerStore.openPlayer(props.media)
}
</script>

<template>
  <div
    class="group relative glass-card rounded-2xl overflow-hidden cursor-pointer flex flex-col"
    @click="handleClick"
  >
    <!-- Poster Container com Overlay -->
    <div class="relative w-full aspect-[2/3] bg-slate-900 overflow-hidden">
      <img
        :src="posterUrl"
        :alt="title"
        loading="lazy"
        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
      />
      <!-- Hover Play Overlay -->
      <div class="absolute inset-0 bg-slate-950/60 opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-center justify-center">
        <div class="w-12 h-12 rounded-full bg-sky-400 text-slate-950 flex items-center justify-center shadow-lg shadow-sky-400/40 transform scale-75 group-hover:scale-100 transition-transform duration-300">
          <Play class="w-6 h-6 fill-slate-950 ml-0.5" />
        </div>
      </div>
      <!-- Rating Badge -->
      <div class="absolute top-2.5 right-2.5 px-2 py-0.5 rounded-lg bg-slate-950/80 backdrop-blur-md border border-slate-700/50 flex items-center gap-1 text-xs font-bold text-amber-400">
        <Star class="w-3 h-3 fill-amber-400" />
        <span>{{ rating }}</span>
      </div>
    </div>

    <!-- Info Body -->
    <div class="p-3.5 flex-1 flex flex-col justify-between">
      <h3 class="font-heading font-semibold text-sm line-clamp-1 group-hover:text-sky-300 transition-colors">
        {{ title }}
      </h3>
      <div class="mt-1 flex items-center justify-between text-xs text-slate-400 font-medium">
        <span>{{ releaseYear }}</span>
        <span class="uppercase tracking-wider text-[10px] px-1.5 py-0.5 rounded bg-slate-800/80 text-slate-300">
          {{ media.media_type || 'filme' }}
        </span>
      </div>
    </div>
  </div>
</template>
