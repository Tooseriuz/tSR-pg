<script setup lang="ts">
import { Github } from 'lucide-vue-next'

const isHeroCollapsed = ref(false)
const heroRef = ref<HTMLElement | null>(null)

const blogTopics = [
  {
    name: 'Architecture Notes',
    description:
      'Short field notes on service boundaries, generated contracts, and the tradeoffs that shape maintainable systems.',
  },
  {
    name: 'Frontend Journal',
    description:
      'Practical writing on Nuxt interfaces, component structure, and the small details that make product screens feel composed.',
  },
  {
    name: 'Infrastructure Logs',
    description:
      'A working record of Cloud Run, Terraform, release paths, and the operational choices behind production software.',
  },
  {
    name: 'Engineering Essays',
    description:
      'Clear reflections on debugging, code review, team habits, and the discipline of building software with fewer surprises.',
  },
]

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
    class="code-pattern min-h-screen scroll-smooth bg-background text-foreground"
  >
    <header
      class="fixed inset-x-0 top-0 z-50 flex h-16 items-center justify-between border-b border-border bg-background/90 px-6 backdrop-blur sm:px-8"
    >
      <a
        class="font-mono text-sm font-bold tracking-normal text-foreground transition duration-300"
        :class="isHeroCollapsed ? 'opacity-100 translate-y-0' : 'pointer-events-none -translate-y-1 opacity-0'"
        href="#top"
        aria-label="Back to top"
      >
        too&lt;&gt;seriuz
      </a>

      <a
        class="inline-flex size-10 items-center justify-center rounded-md border border-border bg-background text-foreground transition hover:bg-surface"
        href="https://github.com/tooseriuz"
        target="_blank"
        rel="noreferrer"
        aria-label="GitHub profile"
      >
        <Github class="size-5" aria-hidden="true" />
      </a>
    </header>

    <section
      id="top"
      ref="heroRef"
      class="flex min-h-screen flex-col p-6 pt-20 sm:p-8 sm:pt-24"
      aria-label="too seriuz landing page"
    >
      <main class="grid flex-1 place-content-center gap-4 text-center">
        <p class="m-0 font-mono text-[0.76rem] font-semibold uppercase tracking-[0.18em] text-muted-foreground">
          software engineer
        </p>
        <h1
          class="m-0 bg-gradient-to-r from-foreground via-primary to-accent bg-clip-text text-5xl font-bold leading-[0.9] tracking-normal text-transparent sm:text-7xl lg:text-9xl"
        >
          too&lt;&gt;seriuz
        </h1>
      </main>

      <footer class="flex justify-center pb-2 pt-6">
        <a
          class="inline-flex h-10 items-center justify-center rounded-md border border-border bg-background px-4 font-mono text-xs font-semibold uppercase tracking-[0.14em] text-foreground transition hover:bg-surface"
          href="#blog-topics"
          aria-label="Scroll to blog topics"
        >
          read notes
        </a>
      </footer>
    </section>

    <section
      id="blog-topics"
      class="border-t border-border bg-surface px-6 py-24 sm:px-8 sm:py-32"
      aria-label="Blog topic placeholders"
    >
      <div class="mx-auto grid max-w-6xl gap-10 lg:grid-cols-[0.8fr_1.2fr] lg:gap-16">
        <div class="max-w-xl">
          <p class="m-0 font-mono text-[0.72rem] font-semibold uppercase tracking-[0.18em] text-muted-foreground">
            blog index
          </p>
          <h2 class="mt-4 font-serif text-4xl leading-tight tracking-[-0.025em] text-foreground sm:text-5xl">
            Working topics for future writing.
          </h2>
          <p class="mt-5 text-base leading-7 text-muted-foreground">
            Placeholder cards for the themes I want to publish around: practical systems, focused interfaces, and notes from real delivery work.
          </p>
        </div>

        <div class="grid gap-4 sm:grid-cols-2">
          <article
            v-for="topic in blogTopics"
            :key="topic.name"
            class="min-h-56 rounded-lg border border-border bg-background p-6 transition duration-200 hover:shadow-soft"
          >
            <p class="m-0 font-mono text-[0.7rem] font-semibold uppercase tracking-[0.16em] text-muted-foreground">
              placeholder
            </p>
            <h3 class="mt-6 text-2xl font-semibold leading-tight tracking-normal text-foreground">
              {{ topic.name }}
            </h3>
            <p class="mt-4 text-sm leading-6 text-muted-foreground">
              {{ topic.description }}
            </p>
          </article>
        </div>
      </div>
    </section>
  </div>
</template>
