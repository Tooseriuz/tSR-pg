<script setup lang="ts">
import type { components } from '../../../shared/generated/frontend'

const health: components['schemas']['HealthResponse'] = {
  status: 'ok',
}
</script>

<template>
  <div class="shell">
    <main class="panel">
      <p class="eyebrow">Placeholder Monorepo</p>
      <h1>Nuxt frontend scaffold</h1>
      <p class="copy">
        This app is intentionally minimal. Tailwind and shadcn/ui structure can
        be layered in later without changing the repo shape.
      </p>
      <div class="actions">
        <span class="badge">Nuxt</span>
        <span class="badge">Tailwind-ready</span>
        <span class="badge">shadcn-ready</span>
        <span class="badge">API status: {{ health.status }}</span>
      </div>
    </main>
  </div>
</template>
