import type { components } from '../../types/dto/openapi'

type AdminVerifyRequest = components['schemas']['AdminVerifyRequest']

export default defineEventHandler(async (event) => {
  const { apiBaseUrl } = useRuntimeConfig()
  const body = await readBody<AdminVerifyRequest>(event)
  const response = await $fetch.raw('/admin-verify', {
    baseURL: apiBaseUrl,
    method: 'POST',
    body,
  })

  const setCookie = response.headers.get('set-cookie')
  if (setCookie) {
    setResponseHeader(event, 'Set-Cookie', setCookie)
  }

  return { ok: true }
})
