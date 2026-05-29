export default defineEventHandler(async (event) => {
  const { apiBaseUrl } = useRuntimeConfig()
  const cookie = getHeader(event, 'cookie')

  await $fetch('/admin-verify', {
    baseURL: apiBaseUrl,
    method: 'GET',
    headers: {
      cookie: cookie ?? '',
    },
  })

  return { ok: true }
})
