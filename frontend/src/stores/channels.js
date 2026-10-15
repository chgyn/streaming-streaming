import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { channelService } from '../services/api'

export const useChannelStore = defineStore('channels', () => {
  const channels = ref([])
  const categories = ref([])
  const selectedCategory = ref('Todos')
  const searchQuery = ref('')
  const isLoading = ref(false)
  const errorMessage = ref(null)

  const filteredChannels = computed(() => {
    let result = channels.value
    if (selectedCategory.value && selectedCategory.value !== 'Todos') {
      result = result.filter(ch => ch.category === selectedCategory.value)
    }
    if (searchQuery.value && searchQuery.value.trim().length > 0) {
      const q = searchQuery.value.trim().toLowerCase()
      result = result.filter(ch =>
        (ch.name && ch.name.toLowerCase().includes(q)) ||
        (ch.category && ch.category.toLowerCase().includes(q)) ||
        (ch.id && ch.id.toLowerCase().includes(q))
      )
    }
    return result
  })

  async function fetchChannels(category = '', query = '') {
    isLoading.value = true
    errorMessage.value = null
    try {
      const res = await channelService.getChannels(category, query)
      channels.value = res.data || []
    } catch (err) {
      errorMessage.value = err.message
      channels.value = []
    } finally {
      isLoading.value = false
    }
  }

  async function fetchCategories() {
    try {
      const res = await channelService.getCategories()
      categories.value = ['Todos', ...(res.data || [])]
    } catch (err) {
      console.warn('Falha ao carregar categorias de canais:', err)
      categories.value = ['Todos']
    }
  }

  function setCategory(cat) {
    selectedCategory.value = cat
  }

  function setSearchQuery(q) {
    searchQuery.value = q
  }

  return {
    channels,
    categories,
    selectedCategory,
    searchQuery,
    isLoading,
    errorMessage,
    filteredChannels,
    fetchChannels,
    fetchCategories,
    setCategory,
    setSearchQuery
  }
})
