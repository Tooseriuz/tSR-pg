<script setup lang="ts">
import { Github } from 'lucide-vue-next'
</script>

<template>
  <div
    class="code-pattern flex min-h-screen flex-col bg-background p-6 text-foreground sm:p-8"
  >
    <main class="grid flex-1 place-content-center gap-3 text-center sm:gap-4" aria-label="too seriuz landing page">
      <p class="m-0 text-[0.82rem] font-semibold uppercase tracking-[0.16em] text-muted-foreground">
        software engineer
      </p>
      <h1
        class="m-0 bg-gradient-to-r from-foreground via-primary to-accent bg-clip-text text-5xl font-bold leading-[0.9] tracking-normal text-transparent sm:text-7xl lg:text-9xl"
      >
        too&lt;&gt;seriuz
      </h1>
    </main>

    <footer class="flex justify-center pb-2 pt-6">
      <UiButton
        as="a"
        variant="outline"
        size="icon"
        href="https://github.com/tooseriuz"
        target="_blank"
        rel="noreferrer"
        aria-label="GitHub profile"
      >
        <Github aria-hidden="true" />
      </UiButton>
    </footer>
  </div>
</template>
