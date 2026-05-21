<script setup lang="ts">
import type { components } from './types/generated/openapi'

type HealthResponse = components['schemas']['HealthResponse']

const { data: health, error, pending } = await useFetch<HealthResponse>('/api/health')
</script>

<template>
  <div class="shell">
    <main class="panel">
      <p class="eyebrow">tSR-pg Monorepo</p>
      <h1>Frontend health check</h1>
      <p class="copy">
        This Nuxt page reads the current API health endpoint through the
        generated OpenAPI response type.
      </p>
      <div class="status-card">
        <span class="status-dot" :class="{ online: health?.status === 'ok' }" />
        <div>
          <p class="status-label">API status</p>
          <p class="status-value">
            <span v-if="pending">Checking...</span>
            <span v-else-if="error">Unavailable</span>
            <span v-else>{{ health?.status }}</span>
          </p>
        </div>
      </div>
    </main>
  </div>
</template>
