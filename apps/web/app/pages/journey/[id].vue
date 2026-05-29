<script setup lang="ts">
import { ArrowLeft } from '@lucide/vue'
import type { components } from '../../../types/dto/openapi'

type JourneyContent = components['schemas']['JourneyContent']

const route = useRoute()
const id = computed(() => String(route.params.id))
const dateFormatter = new Intl.DateTimeFormat('en', {
  day: '2-digit',
  month: 'long',
  year: 'numeric',
  timeZone: 'UTC',
})

const { data: journey, error } = await useFetch<JourneyContent>(() => `/api/journey/${id.value}`)

const happenedOnDate = computed(() => formatDate(journey.value?.timestamp))
const postedOnDate = computed(() => formatDate(journey.value?.created_at))

const renderedContent = computed(() => renderMarkdown(journey.value?.content ?? ''))

useHead(() => ({
  title: journey.value ? `${journey.value.name} - tooseriuz` : 'journey - tooseriuz',
  meta: [
    { name: 'description', content: journey.value?.name ?? 'tooseriuz journey detail' },
  ],
}))

function escapeHtml(value: string) {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;')
}

function formatDate(value?: string) {
  if (!value) {
    return ''
  }

  return dateFormatter.format(new Date(`${value}T00:00:00.000Z`))
}

function renderInlineMarkdown(value: string) {
  return escapeHtml(value)
    .replace(/`([^`]+)`/g, '<code>$1</code>')
    .replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
    .replace(/\*([^*]+)\*/g, '<em>$1</em>')
}

function renderMarkdown(value: string) {
  const blocks: string[] = []
  const lines = value.replace(/\r\n/g, '\n').split('\n')
  let paragraph: string[] = []
  let list: string[] = []
  let code: string[] = []
  let isCodeBlock = false

  function flushParagraph() {
    if (!paragraph.length) {
      return
    }

    blocks.push(`<p>${renderInlineMarkdown(paragraph.join(' '))}</p>`)
    paragraph = []
  }

  function flushList() {
    if (!list.length) {
      return
    }

    blocks.push(`<ul>${list.map(item => `<li>${renderInlineMarkdown(item)}</li>`).join('')}</ul>`)
    list = []
  }

  for (const line of lines) {
    if (line.trim().startsWith('```')) {
      if (isCodeBlock) {
        blocks.push(`<pre><code>${escapeHtml(code.join('\n'))}</code></pre>`)
        code = []
        isCodeBlock = false
      }
      else {
        flushParagraph()
        flushList()
        isCodeBlock = true
      }
      continue
    }

    if (isCodeBlock) {
      code.push(line)
      continue
    }

    const trimmedLine = line.trim()
    if (!trimmedLine) {
      flushParagraph()
      flushList()
      continue
    }

    const heading = /^(#{1,3})\s+(.+)$/.exec(trimmedLine)
    if (heading) {
      flushParagraph()
      flushList()
      const marker = heading[1] ?? '#'
      const text = heading[2] ?? ''
      const level = marker.length + 1
      blocks.push(`<h${level}>${renderInlineMarkdown(text)}</h${level}>`)
      continue
    }

    const listItem = /^[-*]\s+(.+)$/.exec(trimmedLine)
    if (listItem) {
      flushParagraph()
      list.push(listItem[1] ?? '')
      continue
    }

    flushList()
    paragraph.push(trimmedLine)
  }

  flushParagraph()
  flushList()

  if (isCodeBlock && code.length) {
    blocks.push(`<pre><code>${escapeHtml(code.join('\n'))}</code></pre>`)
  }

  return blocks.join('')
}
</script>

<template>
  <main class="min-h-[100dvh] bg-background px-6 py-8 text-foreground sm:px-8 lg:px-12">
    <div class="mx-auto grid w-full max-w-3xl gap-12">
      <NuxtLink
        to="/#journeys"
        class="inline-flex w-fit items-center gap-2 rounded-md border border-border bg-background px-3 py-2 text-sm font-bold text-foreground shadow-soft transition hover:bg-surface active:-translate-y-[1px]"
      >
        <ArrowLeft class="size-4" aria-hidden="true" />
        journeys
      </NuxtLink>

      <article v-if="journey" class="grid gap-10">
        <header class="grid justify-items-center gap-3 border-b border-border pb-10 text-center">
          <h1 class="m-0 max-w-[18ch] text-4xl font-black leading-tight tracking-normal text-foreground sm:text-5xl">
            {{ journey.name }}
          </h1>
          <dl class="m-0 flex flex-wrap justify-center gap-x-5 gap-y-2 font-mono text-xs font-bold uppercase tracking-[0.2em] text-muted-foreground">
            <div class="flex items-center gap-2">
              <dt>happened on</dt>
              <dd class="m-0 text-foreground">{{ happenedOnDate }}</dd>
            </div>
            <div class="flex items-center gap-2">
              <dt>posted on</dt>
              <dd class="m-0 text-foreground">{{ postedOnDate }}</dd>
            </div>
          </dl>
        </header>

        <div class="journal-content" v-html="renderedContent" />
      </article>

      <section
        v-else
        class="grid justify-items-center gap-3 border border-dashed border-border p-8 text-center"
      >
        <h1 class="m-0 text-2xl font-black tracking-normal">
          journey not found
        </h1>
        <p class="m-0 max-w-[44ch] text-base leading-7 text-muted-foreground">
          {{ error ? 'This journal could not be loaded.' : 'There is no journal content for this journey yet.' }}
        </p>
      </section>
    </div>
  </main>
</template>
