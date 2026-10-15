import { defineStore } from 'pinia'
import { ref } from 'vue'
import { playerService } from '../services/api'

export const usePlayerStore = defineStore('player', () => {
  const isOpen = ref(false)
  const currentTitle = ref('')
  const currentMedia = ref(null)
  const embedUrl = ref('')
  const season = ref(1)
  const episode = ref(1)
  const isLoading = ref(false)
  const errorMessage = ref(null)

  async function openPlayer(media) {
    currentMedia.value = media
    currentTitle.value = media.title || media.name || 'Assistir'
    season.value = 1
    episode.value = 1
    isOpen.value = true
    await updateEmbedUrl()
  }

  async function openChannel(channel) {
    currentMedia.value = { ...channel, media_type: 'channel' }
    currentTitle.value = channel.name || 'Canal Ao Vivo'
    season.value = 1
    episode.value = 1
    isOpen.value = true
    await updateEmbedUrl()
  }

  async function updateEmbedUrl() {
    if (!currentMedia.value) return
    isLoading.value = true
    errorMessage.value = null

    const type = currentMedia.value.media_type || 'movie'
    try {
      const data = await playerService.resolve(
        type,
        type === 'channel' ? null : currentMedia.value.id,
        type === 'tv' ? season.value : null,
        type === 'tv' ? episode.value : null,
        type === 'channel' ? currentMedia.value.id : null
      )
      embedUrl.value = data.embed_url
    } catch (err) {
      errorMessage.value = err.message
      embedUrl.value = ''
    } finally {
      isLoading.value = false
    }
  }

  function closePlayer() {
    isOpen.value = false
    embedUrl.value = ''
    currentMedia.value = null
  }

  return {
    isOpen,
    currentTitle,
    currentMedia,
    embedUrl,
    season,
    episode,
    isLoading,
    errorMessage,
    openPlayer,
    openChannel,
    updateEmbedUrl,
    closePlayer
  }
})
