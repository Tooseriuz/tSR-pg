import type { components } from '../../types/dto/openapi'

type CreateJourneyRequest = components['schemas']['CreateJourneyRequest']
type CreateJourneyResponse = components['schemas']['CreateJourneyResponse']

export default defineEventHandler(async (event) => {
  const { apiBaseUrl } = useRuntimeConfig()
  const body = await readBody<CreateJourneyRequest>(event)
  const cookie = getHeader(event, 'cookie')

  return $fetch<CreateJourneyResponse>('/journey', {
    baseURL: apiBaseUrl,
    method: 'POST',
    body,
    headers: {
      cookie: cookie ?? '',
    },
  })
})
