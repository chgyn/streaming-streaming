import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Accept': 'application/json',
    'Content-Type': 'application/json'
  }
})

// Interceptor de resposta para tratamento uniforme de erros
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    const errorMsg = error.response?.data?.error || error.message || 'Erro na comunicação com o servidor'
    return Promise.reject(new Error(errorMsg))
  }
)

export const tmdbService = {
  getTrending: (mediaType = 'movie', timeWindow = 'week', page = 1) =>
    api.get('/trending', { params: { media_type: mediaType, time_window: timeWindow, page } }),

  search: (query, type = 'multi', page = 1) =>
    api.get('/search', { params: { q: query, type, page } }),

  getDetails: (mediaType, id) =>
    api.get(`/details/${mediaType}/${id}`),

  getHealth: () =>
    api.get('/health')
}

export const playerService = {
  resolve: (type, tmdbId, season = null, episode = null, channelId = null) => {
    const params = { type }
    if (type === 'channel') {
      params.channel_id = channelId
    } else {
      params.tmdb_id = tmdbId
      if (type === 'tv') {
        if (season != null) params.season = season
        if (episode != null) params.episode = episode
      }
    }
    return api.get('/player/resolve', { params })
  }
}

export const channelService = {
  getChannels: (category = '', q = '') => {
    const params = {}
    if (category && category !== 'Todos') params.category = category
    if (q) params.q = q
    return api.get('/channels', { params })
  },

  getCategories: () =>
    api.get('/channels/categories')
}

export default api

