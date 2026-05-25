<script setup lang="ts">
import { ArrowLeft } from 'lucide-vue-next'
import { blogTopics } from '../../data/blogTopics'

const route = useRoute()
const topic = computed(() =>
  blogTopics.find((item) => item.slug === route.params.slug),
)

if (!topic.value) {
  throw createError({
    statusCode: 404,
    statusMessage: 'Note topic not found',
  })
}
</script>

<template>
  <main class="code-pattern min-h-[100dvh] bg-background px-6 py-8 text-foreground sm:px-8 sm:py-10">
    <div class="mx-auto max-w-4xl">
      <NuxtLink
        class="inline-flex h-10 items-center gap-2 rounded-md border border-border bg-background px-4 font-mono text-xs font-semibold uppercase tracking-[0.14em] text-foreground no-underline transition hover:bg-surface active:-translate-y-[1px]"
        to="/#blog-topics"
      >
        <ArrowLeft class="size-4" aria-hidden="true" />
        Notes
      </NuxtLink>

      <article class="pt-16 sm:pt-20">
        <p class="m-0 font-mono text-[0.72rem] font-semibold uppercase tracking-[0.18em] text-muted-foreground">
          topic placeholder
        </p>
        <h1 class="mt-4 text-5xl font-semibold leading-none tracking-normal text-foreground sm:text-7xl">
          {{ topic?.name }}
        </h1>
        <p class="mt-8 max-w-2xl text-lg leading-8 text-muted-foreground">
          {{ topic?.description }}
        </p>

        <div class="mt-14 grid gap-8 border-t border-border pt-10 lg:grid-cols-[0.7fr_1.3fr]">
          <aside>
            <p class="m-0 font-mono text-[0.72rem] font-semibold uppercase tracking-[0.18em] text-muted-foreground">
              draft state
            </p>
            <p class="mt-4 text-sm leading-6 text-muted-foreground">
              Placeholder page only. Final posts can replace this structure with real sections, images, code notes, and references.
            </p>
          </aside>

          <section class="space-y-8">
            <div>
              <h2 class="text-2xl font-semibold leading-tight tracking-normal text-foreground">
                Overview
              </h2>
              <p class="mt-4 text-base leading-7 text-muted-foreground">
                {{ topic?.detail }}
              </p>
            </div>

            <div>
              <h2 class="text-2xl font-semibold leading-tight tracking-normal text-foreground">
                More Details
              </h2>
              <p class="mt-4 text-base leading-7 text-muted-foreground">
                This placeholder area is reserved for deeper notes, supporting examples, implementation choices, and links to future writing.
              </p>
            </div>
          </section>
        </div>
      </article>
    </div>
  </main>
</template>
