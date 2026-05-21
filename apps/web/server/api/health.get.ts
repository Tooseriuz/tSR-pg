import type { components } from '../../types/dto/openapi'

type HealthResponse = components['schemas']['HealthResponse']

export default defineEventHandler(async () => {
  const { apiBaseUrl } = useRuntimeConfig()

  return $fetch<HealthResponse>('/health', {
    baseURL: apiBaseUrl,
  })
})
