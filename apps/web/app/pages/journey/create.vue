<script setup lang="ts">
import { ArrowLeft, Check, KeyRound, Send, Star } from '@lucide/vue'
import type { components } from '../../../types/dto/openapi'

type CreateJourneyRequest = components['schemas']['CreateJourneyRequest']
type CreateJourneyResponse = components['schemas']['CreateJourneyResponse']

const router = useRouter()
const contentInput = ref<HTMLTextAreaElement | null>(null)
const tokenInput = ref('')
const name = ref('')
const location = ref('')
const content = ref('')
const highlight = ref(false)
const tokenError = ref('')
const submitError = ref('')
const isVerifying = ref(false)
const isSubmitting = ref(false)
const isVerified = ref(false)
const dateFormatter = new Intl.DateTimeFormat('en', {
  day: '2-digit',
  month: 'long',
  year: 'numeric',
  timeZone: 'UTC',
})
const trimmedName = computed(() => name.value.trim())
const trimmedLocation = computed(() => location.value.trim())
const trimmedContent = computed(() => content.value.trim())
const isNameReady = computed(() => trimmedName.value.length > 0)
const isLocationReady = computed(() => trimmedLocation.value.length > 0)
const isContentReady = computed(() => trimmedContent.value.length > 0)
const canSubmit = computed(() => isVerified.value && isNameReady.value && isLocationReady.value && isContentReady.value && !isSubmitting.value)
const renderedPreview = computed(() => renderMarkdown(trimmedContent.value))
const hasPreview = computed(() => isNameReady.value || isContentReady.value)
const previewDate = computed(() => dateFormatter.format(new Date()))

useHead({
  title: 'create journey - tooseriuz',
  meta: [
    { name: 'description', content: 'Admin journey creation page' },
  ],
})

onMounted(async () => {
  resizeContentInput()

  try {
    await $fetch('/api/admin-verify')
    isVerified.value = true
  }
  catch {
    isVerified.value = false
  }
})

watch(content, async () => {
  await nextTick()
  resizeContentInput()
})

function resizeContentInput() {
  const input = contentInput.value

  if (!input) {
    return
  }

  input.style.height = 'auto'
  input.style.height = `${input.scrollHeight}px`
}

async function verifyToken() {
  tokenError.value = ''
  const token = tokenInput.value.trim()
  if (!token) {
    tokenError.value = 'Admin token is required.'
    return
  }

  isVerifying.value = true
  try {
    await $fetch('/api/admin-verify', {
      method: 'POST',
      body: { token },
    })
    isVerified.value = true
    tokenInput.value = ''
  }
  catch {
    tokenError.value = 'Admin token was rejected.'
  }
  finally {
    isVerifying.value = false
  }
}

async function submitJourney() {
  submitError.value = ''
  const request: CreateJourneyRequest = {
    name: trimmedName.value,
    location: trimmedLocation.value,
    content: trimmedContent.value,
    highlight: highlight.value,
  }

  if (!request.name || !request.location || !request.content) {
    submitError.value = 'Journey name, location, and content are required.'
    return
  }

  isSubmitting.value = true
  try {
    const response = await $fetch<CreateJourneyResponse>('/api/journey', {
      method: 'POST',
      body: request,
    })

    await router.push(`/journey/${response.id}`)
  }
  catch {
    submitError.value = 'Journey could not be submitted.'
  }
  finally {
    isSubmitting.value = false
  }
}

function escapeHtml(value: string) {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;')
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
  <main class="min-h-[100dvh] bg-background px-5 py-6 text-foreground sm:px-8 lg:px-12">
    <div class="mx-auto grid w-full max-w-7xl gap-5">
      <div class="flex flex-wrap items-center justify-between gap-4 border-b border-border pb-5">
        <NuxtLink
          to="/#journeys"
          class="inline-flex w-fit items-center gap-2 rounded-md border border-border bg-background px-3 py-2 text-sm font-bold text-foreground shadow-soft transition hover:bg-surface active:-translate-y-[1px]"
        >
          <ArrowLeft class="size-4" aria-hidden="true" />
          journeys
        </NuxtLink>
      </div>

      <form id="journey-create-form" class="grid gap-5" @submit.prevent="submitJourney">
        <section class="grid gap-5">
          <header class="relative isolate overflow-hidden rounded-lg px-5 py-6 sm:px-6">
            <div class="grid gap-3">
              <p class="m-0 font-mono text-xs font-bold uppercase tracking-[0.32em] text-accent">
                admin editor
              </p>
              <div class="flex flex-wrap items-end justify-between gap-4">
                <h1 class="m-0 max-w-[12ch] text-4xl font-black leading-none tracking-normal text-foreground sm:text-5xl">
                  create journey
                </h1>
              </div>
            </div>
          </header>

          <div class="grid gap-2 rounded-lg border border-border bg-background p-4 shadow-soft">
            <div class="flex items-center justify-between gap-3">
              <label class="text-sm font-bold text-foreground" for="journey-name">journey name</label>
            </div>
            <input
              id="journey-name"
              v-model="name"
              class="min-h-12 rounded-md border border-border bg-surface px-4 text-base leading-6 outline-none transition placeholder:text-muted-foreground focus:border-primary focus:bg-background"
              :disabled="!isVerified || isSubmitting"
            >
          </div>

          <div class="grid gap-2 rounded-lg border border-border bg-background p-4 shadow-soft">
            <div class="flex items-center justify-between gap-3">
              <label class="text-sm font-bold text-foreground" for="journey-location">location</label>
            </div>
            <input
              id="journey-location"
              v-model="location"
              class="min-h-12 rounded-md border border-border bg-surface px-4 text-base leading-6 outline-none transition placeholder:text-muted-foreground focus:border-primary focus:bg-background"
              :disabled="!isVerified || isSubmitting"
            >
          </div>

          <div class="gap-2 overflow-hidden rounded-lg border border-border bg-background p-4 shadow-soft">
            <div class="grid gap-3 pb-2 sm:grid-cols-[1fr_auto] sm:items-center">
              <div class="flex items-center gap-2">
                <label class="text-sm font-bold text-foreground" for="journey-content">content</label>
              </div>
            </div>

            <textarea
              id="journey-content"
              ref="contentInput"
              v-model="content"
              class="min-h-12 w-full overflow-hidden rounded-md border border-border bg-surface px-4 py-4 text-base leading-6 outline-none transition placeholder:text-muted-foreground focus:border-primary focus:bg-background disabled:opacity-55"
              :disabled="!isVerified || isSubmitting"
              @input="resizeContentInput"
            />
          </div>
        </section>

      </form>

      <section class="overflow-hidden rounded-lg border border-border bg-background p-4 shadow-soft" aria-labelledby="journey-preview-title">
        <div class="flex flex-wrap items-center justify-between gap-3 pb-4">
          <p id="journey-preview-title" class="m-0 font-mono text-[11px] font-bold uppercase tracking-[0.24em] text-accent">
            preview
          </p>
        </div>

        <div class="mx-auto grid w-full max-w-3xl gap-12 px-2 py-8 sm:px-4 lg:px-0">
          <article v-if="hasPreview" class="grid gap-10">
            <header class="grid justify-items-center gap-3 border-b border-border pb-10 text-center">
              <h1 class="m-0 max-w-[18ch] text-4xl font-black leading-tight tracking-normal text-foreground sm:text-5xl">
                {{ trimmedName || 'Untitled journey' }}
              </h1>
              <p class="m-0 font-mono text-xs font-bold uppercase tracking-[0.28em] text-muted-foreground">
                {{ previewDate }}
              </p>
            </header>

            <div v-if="isContentReady" class="journal-content" v-html="renderedPreview" />
            <section
              v-else
              class="grid justify-items-center gap-3 border border-dashed border-border p-8 text-center"
            >
              <h2 class="m-0 text-2xl font-black tracking-normal">
                no content yet
              </h2>
              <p class="m-0 max-w-[44ch] text-base leading-7 text-muted-foreground">
                Add content above to preview the journey body.
              </p>
            </section>
          </article>

          <section
            v-else
            class="grid justify-items-center gap-3 border border-dashed border-border p-8 text-center"
          >
            <h1 class="m-0 text-2xl font-black tracking-normal">
              no preview yet
            </h1>
            <p class="m-0 max-w-[44ch] text-base leading-7 text-muted-foreground">
              Start with a title or body content and the submitted journey preview will appear here.
            </p>
          </section>
        </div>

        <div class="grid justify-items-end gap-2 border-t border-border pt-4">
          <p v-if="submitError" class="m-0 w-full rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm font-semibold text-red-700">
            {{ submitError }}
          </p>

          <div class="flex flex-wrap justify-end gap-2">
            <button
              type="button"
              class="inline-flex size-9 items-center justify-center rounded-md border shadow-soft transition active:-translate-y-[1px] disabled:cursor-not-allowed disabled:opacity-55"
              :class="highlight ? 'border-primary bg-primary text-white hover:bg-primary/90' : 'border-border bg-background text-foreground hover:bg-surface'"
              :aria-pressed="highlight"
              aria-label="Toggle highlighted journey"
              :disabled="!isVerified || isSubmitting"
              @click="highlight = !highlight"
            >
              <Star class="size-4" :class="highlight ? 'fill-current' : ''" aria-hidden="true" />
            </button>

            <button
              form="journey-create-form"
              type="submit"
              class="inline-flex w-fit items-center justify-center gap-2 rounded-md border border-primary bg-primary px-3 py-2 text-sm font-bold text-white shadow-soft transition hover:bg-primary/90 active:-translate-y-[1px] disabled:cursor-not-allowed disabled:border-border disabled:bg-background disabled:text-muted-foreground disabled:opacity-100"
              :disabled="!canSubmit"
            >
              <Send class="size-4" aria-hidden="true" />
              {{ isSubmitting ? 'submitting' : 'submit journey' }}
            </button>
          </div>
          <p class="m-0 text-right text-xs font-semibold leading-5 text-muted-foreground">
            {{ isVerified ? 'Title, location, and content are required before submit.' : 'Verify the admin token to unlock writing.' }}
          </p>
        </div>
      </section>
    </div>

    <div
      v-if="!isVerified"
      class="fixed inset-0 grid place-items-center bg-background/88 px-6 backdrop-blur-md"
      role="dialog"
      aria-modal="true"
      aria-labelledby="admin-token-title"
    >
      <form
        class="grid w-full max-w-sm gap-5 rounded-lg border border-border bg-background p-6 shadow-soft"
        @submit.prevent="verifyToken"
      >
        <div class="grid gap-3">
          <div class="grid size-11 place-items-center rounded-md bg-surface text-foreground">
            <KeyRound class="size-5" aria-hidden="true" />
          </div>
          <h2 id="admin-token-title" class="m-0 text-2xl font-black tracking-normal text-foreground">
            admin token
          </h2>
        </div>

        <div class="grid gap-2">
          <label class="text-sm font-bold text-foreground" for="admin-token">token</label>
          <input
            id="admin-token"
            v-model="tokenInput"
            class="min-h-12 rounded-md border border-border bg-surface px-4 text-base font-semibold outline-none transition placeholder:text-muted-foreground focus:border-primary"
            type="password"
            autocomplete="off"
            autofocus
            :disabled="isVerifying"
          >
          <p v-if="tokenError" class="m-0 text-sm font-semibold text-red-700">
            {{ tokenError }}
          </p>
        </div>

        <button
          type="submit"
          class="inline-flex min-h-12 items-center justify-center gap-2 rounded-md border border-primary bg-primary px-5 text-sm font-black text-white shadow-soft transition hover:bg-primary/90 active:-translate-y-[1px] disabled:cursor-not-allowed disabled:border-border disabled:bg-background disabled:text-muted-foreground disabled:opacity-100"
          :disabled="isVerifying"
        >
          <Check class="size-4" aria-hidden="true" />
          {{ isVerifying ? 'checking' : 'verify' }}
        </button>
      </form>
    </div>
  </main>
</template>
