<script setup>
import { ref } from 'vue'
import { Play, Radio, Tv } from 'lucide-vue-next'
import { usePlayerStore } from '../stores/player'

const props = defineProps({
  channel: {
    type: Object,
    required: true
  }
})

const playerStore = usePlayerStore()
const imgError = ref(false)

function handleClick() {
  playerStore.openChannel(props.channel)
}

function handleImgError() {
  imgError.value = true
}
</script>

<template>
  <div
    class="group relative glass-card rounded-2xl overflow-hidden cursor-pointer flex flex-col border border-slate-800/80 hover:border-sky-500/50 transition-all duration-300 hover:shadow-xl hover:shadow-sky-500/10"
    @click="handleClick"
  >
    <!-- Logo & Visual Container -->
    <div class="relative w-full aspect-[16/10] bg-gradient-to-b from-slate-900 to-slate-950 flex items-center justify-center p-6 overflow-hidden">
      <!-- Glow decorativo sutil no fundo -->
      <div class="absolute inset-0 bg-radial from-sky-500/5 to-transparent pointer-events-none"></div>

      <!-- Imagem do Logo com fallback -->
      <img
        v-if="channel.logo_url && !imgError"
        :src="channel.logo_url"
        :alt="channel.name"
        loading="lazy"
        referrerpolicy="no-referrer"
        class="max-h-16 max-w-[80%] object-contain drop-shadow-md transition-transform duration-500 group-hover:scale-110"
        @error="handleImgError"
      />
      <div v-else class="flex flex-col items-center justify-center text-slate-500 gap-1.5">
        <Tv class="w-10 h-10 text-slate-600" />
        <span class="text-xs font-semibold text-slate-400 text-center line-clamp-1">{{ channel.name }}</span>
      </div>

      <!-- Hover Play Overlay -->
      <div class="absolute inset-0 bg-slate-950/70 opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-center justify-center">
        <div class="w-12 h-12 rounded-full bg-sky-400 text-slate-950 flex items-center justify-center shadow-lg shadow-sky-400/40 transform scale-75 group-hover:scale-100 transition-transform duration-300">
          <Play class="w-6 h-6 fill-slate-950 ml-0.5" />
        </div>
      </div>

      <!-- Badge AO VIVO -->
      <div class="absolute top-2.5 left-2.5 inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-rose-500/20 border border-rose-500/40 text-[10px] font-extrabold text-rose-400 tracking-wider uppercase">
        <span class="w-1.5 h-1.5 rounded-full bg-rose-400 animate-pulse"></span>
        <span>Ao Vivo</span>
      </div>

      <!-- Badge de Categoria -->
      <div v-if="channel.category" class="absolute top-2.5 right-2.5 px-2 py-0.5 rounded-lg bg-slate-900/80 backdrop-blur-md border border-slate-700/60 text-[10px] font-medium text-slate-300 truncate max-w-[120px]">
        {{ channel.category }}
      </div>
    </div>

    <!-- Info Footer -->
    <div class="p-3.5 bg-slate-950/60 border-t border-slate-800/60 flex items-center justify-between gap-2">
      <div class="truncate">
        <h3 class="font-heading font-semibold text-sm text-slate-100 group-hover:text-sky-300 transition-colors truncate">
          {{ channel.name }}
        </h3>
        <p class="text-[11px] text-slate-400 truncate">
          {{ channel.category || 'Canal de TV' }}
        </p>
      </div>
      <div class="w-7 h-7 rounded-lg bg-slate-900 border border-slate-800 flex items-center justify-center text-slate-400 group-hover:text-sky-400 group-hover:border-sky-500/40 transition-colors flex-shrink-0">
        <Radio class="w-3.5 h-3.5" />
      </div>
    </div>
  </div>
</template>
