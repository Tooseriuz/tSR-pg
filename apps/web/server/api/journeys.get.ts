import type { components } from '../../types/dto/openapi'

type JourneysResponse = components['schemas']['JourneysResponse']

export default defineEventHandler(async () => {
  const { apiBaseUrl } = useRuntimeConfig()

  return $fetch<JourneysResponse>('/journeys', {
    baseURL: apiBaseUrl,
  })
})
