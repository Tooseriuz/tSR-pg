import type { components } from '../../../types/dto/openapi'

type UploadJourneyImagesResponse = components['schemas']['UploadJourneyImagesResponse']

export default defineEventHandler(async (event) => {
  const { apiBaseUrl } = useRuntimeConfig()
  const parts = await readMultipartFormData(event)
  const cookie = getHeader(event, 'cookie')
  const formData = new FormData()

  for (const part of parts ?? []) {
    if (!part.filename || part.name !== 'images') {
      continue
    }

    formData.append('images', new Blob([part.data], { type: part.type }), part.filename)
  }

  return $fetch<UploadJourneyImagesResponse>('/journey/images', {
    baseURL: apiBaseUrl,
    method: 'POST',
    body: formData,
    headers: {
      cookie: cookie ?? '',
    },
  })
})
