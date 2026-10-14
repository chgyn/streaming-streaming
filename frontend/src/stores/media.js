import { defineStore } from 'pinia'
import { ref } from 'vue'
import { tmdbService } from '../services/api'

export const useMediaStore = defineStore('media', () => {
  const trendingItems = ref([])
  const searchResults = ref([])
  const currentMediaType = ref('movie')
  const isLoading = ref(false)
  const errorMessage = ref(null)
  const healthStats = ref(null)

  async function fetchTrending(mediaType = currentMediaType.value, page = 1) {
    isLoading.value = true
    errorMessage.value = null
    currentMediaType.value = mediaType
    try {
      const data = await tmdbService.getTrending(mediaType, 'week', page)
      trendingItems.value = data.results || []
    } catch (err) {
      errorMessage.value = err.message
      trendingItems.value = []
    } finally {
      isLoading.value = false
    }
  }

  async function search(query, type = 'multi', page = 1) {
    if (!query || query.trim().length < 2) return
    isLoading.value = true
    errorMessage.value = null
    try {
      const data = await tmdbService.search(query.trim(), type, page)
      searchResults.value = data.results || []
    } catch (err) {
      errorMessage.value = err.message
      searchResults.value = []
    } finally {
      isLoading.value = false
    }
  }

  async function fetchHealth() {
    try {
      healthStats.value = await tmdbService.getHealth()
    } catch (err) {
      console.warn('Falha ao obter status de saúde:', err)
    }
  }

  return {
    trendingItems,
    searchResults,
    currentMediaType,
    isLoading,
    errorMessage,
    healthStats,
    fetchTrending,
    search,
    fetchHealth
  }
})
