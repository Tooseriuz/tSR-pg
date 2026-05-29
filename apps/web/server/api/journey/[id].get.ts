import type { components } from '../../../types/dto/openapi'

type JourneyContent = components['schemas']['JourneyContent']

export default defineEventHandler(async (event) => {
  const { apiBaseUrl } = useRuntimeConfig()
  const id = getRouterParam(event, 'id')

  return $fetch<JourneyContent>(`/journey/${id}`, {
    baseURL: apiBaseUrl,
  })
})
