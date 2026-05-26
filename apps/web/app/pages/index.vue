<script setup lang="ts">
import { Mail } from '@lucide/vue'
import githubIcon from '../assets/icons/github.svg'

const isHeroCollapsed = ref(false)
const heroRef = ref<HTMLElement | null>(null)

useHead({
  link: [
    { rel: 'icon', href: '/favicon.ico' },
  ],
})

onMounted(() => {
  const hero = heroRef.value

  if (!hero) {
    return
  }

  const observer = new IntersectionObserver(
    ([entry]) => {
      if (!entry) {
        return
      }

      isHeroCollapsed.value = !entry.isIntersecting
    },
    {
      threshold: 0.18,
    },
  )

  observer.observe(hero)

  onUnmounted(() => {
    observer.disconnect()
  })
})
</script>

<template>
  <div
    class="min-h-screen scroll-smooth bg-background text-foreground "
  >
    <section
      id="top"
      ref="heroRef"
      class="flex min-h-[100dvh] flex-col p-6 pt-20 sm:p-8 sm:pt-24"
    >
      <main class="grid flex-1 place-content-center gap-5 text-center">
        <h1
          class="brand-wordmark m-0 flex select-none items-center justify-center whitespace-nowrap text-5xl font-black leading-none tracking-normal text-foreground sm:text-7xl lg:text-9xl"
        >
          <span>too</span>
          <span class="mx-3 text-primary sm:mx-5 lg:mx-7">&lt;</span>
          <span class="text-accent">&gt;</span>
          <span class="ml-3 sm:ml-5 lg:ml-7">seriuz</span>
        </h1>
        <p
          class="m-0 px-2.5 font-mono text-left select-none text-xs font-semibold tracking-[0.3em] text-muted-foreground sm:text-sm"
          aria-label="Wasuphon Naksut"
        >
          <span aria-hidden="true">&gt;</span><span class="typing-name" aria-hidden="true">Wasuphon Naksut</span>
        </p>
      </main>

      <footer class="flex justify-center gap-2 pb-2 pt-6">
        <a
          class="inline-flex select-none size-10 items-center justify-center rounded-md border border-border bg-background text-foreground transition hover:bg-surface active:-translate-y-[1px]"
          href="mailto:hello@tooseriuz.com"
          aria-label="Email hello@tooseriuz.com"
        >
          <Mail class="size-5" aria-hidden="true" />
        </a>

        <a
          class="inline-flex select-none size-10 items-center justify-center rounded-md border border-border bg-background text-foreground transition hover:bg-surface active:-translate-y-[1px]"
          href="/cv.pdf"
          target="_blank"
          rel="noreferrer"
          aria-label="View CV"
        >
          <span class="font-mono text-xs font-bold tracking-normal" aria-hidden="true">CV</span>
        </a>

        <a
          class="inline-flex select-none size-10 items-center justify-center rounded-md border border-border bg-background text-foreground transition hover:bg-surface active:-translate-y-[1px]"
          href="https://github.com/tooseriuz"
          target="_blank"
          rel="noreferrer"
          aria-label="GitHub profile"
        >
          <img class="size-5" :src="githubIcon" alt="" aria-hidden="true">
        </a>
      </footer>
    </section>
  </div>
</template>
