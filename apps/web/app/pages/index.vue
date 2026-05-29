<script setup lang="ts">
import { ArrowUpRight, CalendarDays, ChevronDown, ChevronLeft, ChevronRight, Mail, MapPin } from '@lucide/vue'
import githubIcon from '../assets/icons/github.svg'
import type { components } from '../../types/dto/openapi'

type JourneyYear = `${number}`
type JourneyMonth = 'Jan' | 'Feb' | 'Mar' | 'Apr' | 'May' | 'Jun' | 'Jul' | 'Aug' | 'Sep' | 'Oct' | 'Nov' | 'Dec'
type ApiJourney = components['schemas']['Journey']

interface JourneyPoint {
  year: JourneyYear
  month: JourneyMonth
}

interface Journey {
  year: JourneyYear
  month: JourneyMonth
  topic: string
  place: string
  date: string
  thumbnail: string
}

const isHeroCollapsed = ref(false)
const brandMarkRef = ref<HTMLElement | null>(null)
const journeyTimelineRef = ref<HTMLElement | null>(null)
const selectedJourneyYear = ref<JourneyYear | null>(null)
const selectedJourneyPoint = ref<JourneyPoint | null>(null)
let brandMarkObserver: IntersectionObserver | null = null

const monthOrder: JourneyMonth[] = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
const monthFormatter = new Intl.DateTimeFormat('en', { month: 'short', timeZone: 'UTC' })

const { data: apiJourneys } = await useFetch<ApiJourney[]>('/api/journeys', {
  default: () => [],
})

const journeys = computed<Journey[]>(() =>
  apiJourneys.value.map((journey) => {
    const timestamp = new Date(`${journey.timestamp}T00:00:00.000Z`)
    const year = String(timestamp.getUTCFullYear()) as JourneyYear
    const month = monthFormatter.format(timestamp) as JourneyMonth

    return {
      year,
      month,
      topic: journey.name,
      place: journey.location,
      date: `${month} ${year}`,
      thumbnail: journey.thumbnail,
    }
  }),
)

const journeyTimeline = computed(() =>
  [...new Set(journeys.value.map(journey => journey.year))]
    .sort((previousYear, nextYear) => Number(previousYear) - Number(nextYear))
    .map(year => ({
      year,
      months: journeys.value
        .filter(journey => journey.year === year)
        .map(journey => journey.month)
        .filter((month, index, months) => months.indexOf(month) === index),
    })),
)

function getJourneyMonthOffset(month: JourneyMonth) {
  const monthIndex = monthOrder.indexOf(month)

  if (monthIndex < 0 || monthOrder.length <= 1) {
    return '50%'
  }

  return `${16 + (monthIndex / (monthOrder.length - 1)) * 68}%`
}

const visibleJourneys = computed(() => {
  if (selectedJourneyPoint.value) {
    return journeys.value.filter(journey =>
      journey.year === selectedJourneyPoint.value?.year
      && journey.month === selectedJourneyPoint.value.month,
    )
  }

  if (selectedJourneyYear.value) {
    return journeys.value.filter(journey => journey.year === selectedJourneyYear.value)
  }

  return journeys.value.slice(0, 3)
})

const hasJourneys = computed(() => journeys.value.length > 0)

function selectJourneyYear(year: JourneyYear) {
  selectedJourneyYear.value = selectedJourneyYear.value === year ? null : year
  selectedJourneyPoint.value = null
}

function selectJourneyPoint(point: JourneyPoint) {
  const isSelected = selectedJourneyPoint.value?.year === point.year
    && selectedJourneyPoint.value.month === point.month

  selectedJourneyPoint.value = isSelected ? null : point
  selectedJourneyYear.value = null
}

function moveJourneyTimeline(direction: 'previous' | 'next') {
  const timeline = journeyTimelineRef.value

  if (!timeline) {
    return
  }

  timeline.scrollBy({
    left: direction === 'next' ? timeline.clientWidth * 0.72 : timeline.clientWidth * -0.72,
    behavior: 'smooth',
  })
}

useHead({
  title: 'tooseriuz - software engineer',
  meta: [
    { name: 'description', content: 'tooseriuz, Wasuphon Naksut, software engineer' },
  ],
  link: [
    { rel: 'icon', href: '/favicon.ico' },
  ],
})

onMounted(() => {
  const brandMark = brandMarkRef.value

  if (!brandMark) {
    return
  }

  brandMarkObserver = new IntersectionObserver(
    ([entry]) => {
      if (!entry) {
        return
      }

      isHeroCollapsed.value = !entry.isIntersecting
    },
    {
      threshold: 0,
    },
  )

  brandMarkObserver.observe(brandMark)
})

onUnmounted(() => {
  brandMarkObserver?.disconnect()
})
</script>

<template>
  <div
    class="min-h-screen scroll-smooth bg-background text-foreground "
  >
    <a
      href="#top"
      class="brand-wordmark fixed left-5 top-5 z-20 m-0 flex select-none items-center whitespace-nowrap text-xl font-black leading-none tracking-normal text-foreground transition duration-300 sm:text-2xl"
      :class="isHeroCollapsed ? 'translate-y-0 opacity-100' : 'pointer-events-none -translate-y-3 opacity-0'"
      aria-label="Back to top"
    >
      <span>too</span>
      <span class="mx-1.5 text-primary">&lt;</span>
      <span class="text-accent">&gt;</span>
      <span class="ml-1.5">seriuz</span>
    </a>

    <nav
      class="fixed bottom-5 left-1/2 z-20 flex -translate-x-1/2 gap-2 rounded-lg border border-border bg-background/88 p-1.5 shadow-soft backdrop-blur-md"
      aria-label="Contact links"
    >
      <a
        class="inline-flex select-none size-10 items-center justify-center rounded-md text-foreground transition hover:bg-surface active:-translate-y-[1px]"
        href="mailto:hello@tooseriuz.com"
        aria-label="Email hello@tooseriuz.com"
      >
        <Mail class="size-5" aria-hidden="true" />
      </a>

      <a
        class="inline-flex select-none size-10 items-center justify-center rounded-md text-foreground transition hover:bg-surface active:-translate-y-[1px]"
        href="/cv.pdf"
        target="_blank"
        rel="noreferrer"
        aria-label="View CV"
      >
        <span class="font-mono text-xs font-bold tracking-normal" aria-hidden="true">CV</span>
      </a>

      <a
        class="inline-flex select-none size-10 items-center justify-center rounded-md text-foreground transition hover:bg-surface active:-translate-y-[1px]"
        href="https://github.com/tooseriuz"
        target="_blank"
        rel="noreferrer"
        aria-label="GitHub profile"
      >
        <img class="size-5" :src="githubIcon" alt="" aria-hidden="true">
      </a>
    </nav>

    <section
      id="top"
      class="relative flex min-h-[100dvh] flex-col p-6 pt-20 sm:p-8 sm:pt-24"
    >
      <main class="grid flex-1 place-content-center justify-items-center gap-5 text-center">
        <div class="grid w-fit justify-items-start gap-5">
          <h1
            ref="brandMarkRef"
            class="brand-wordmark hero-wordmark m-0 flex select-none items-center justify-center whitespace-nowrap text-5xl font-black leading-none tracking-normal text-foreground sm:text-7xl lg:text-9xl"
          >
            <span>too</span>
            <span class="mx-3 text-primary sm:mx-5 lg:mx-7">&lt;</span>
            <span class="text-accent">&gt;</span>
            <span class="ml-3 sm:ml-5 lg:ml-7">seriuz</span>
          </h1>
          <p
            class="hero-typewriter m-0 grid select-none gap-1 text-left font-mono text-xs font-semibold tracking-[0.3em] text-muted-foreground sm:text-sm"
            aria-label="Wasuphon Naksut, software engineer"
          >
            <span class="terminal-line" aria-hidden="true">
              <span class="typing-prompt-name text-accent">&gt;</span><span class="typing-name">Wasuphon Naksut</span>
            </span>
            <span class="terminal-line" aria-hidden="true">
              <span class="typing-prompt-role text-accent">&gt;</span><span class="typing-role">software engineer</span>
            </span>
          </p>
        </div>
      </main>
      <a
        href="#journeys"
        class="scroll-cue absolute bottom-24 left-1/2 grid size-12 place-items-center rounded-full border border-border bg-background text-foreground shadow-soft transition hover:bg-surface active:-translate-y-[1px]"
        aria-label="Scroll to journeys"
      >
        <ChevronDown class="size-5" aria-hidden="true" />
      </a>
    </section>

    <section
      id="journeys"
      class="relative isolate grid min-h-[100dvh] content-center overflow-hidden border-t border-border bg-background px-6 py-20 sm:px-8 lg:px-12"
    >
      <div
        class="pointer-events-none absolute -right-36 -top-28 -z-10 h-36 w-[34rem] rotate-[-18deg] bg-primary/15 sm:-right-28 sm:w-[46rem]"
        aria-hidden="true"
      />
      <div
        class="pointer-events-none absolute -right-40 -top-8 -z-10 h-2 w-[35rem] rotate-[-18deg] bg-primary/45 sm:-right-28 sm:w-[48rem]"
        aria-hidden="true"
      />
      <div
        class="pointer-events-none absolute -right-44 top-20 -z-10 h-4 w-[42rem] rotate-[-18deg] bg-primary/35 sm:-right-32"
        aria-hidden="true"
      />
      <div
        class="pointer-events-none absolute -right-52 top-36 -z-10 h-2 w-[28rem] rotate-[-18deg] bg-primary/25 sm:-right-40 sm:w-[38rem]"
        aria-hidden="true"
      />
      <div
        class="pointer-events-none absolute -right-24 top-48 -z-10 h-1.5 w-[18rem] rotate-[-18deg] bg-primary/40 sm:w-[26rem]"
        aria-hidden="true"
      />
      <div
        class="pointer-events-none absolute -bottom-24 -left-40 -z-10 h-32 w-[32rem] rotate-[-18deg] bg-primary/10 sm:w-[44rem]"
        aria-hidden="true"
      />
      <div
        class="pointer-events-none absolute -bottom-14 -left-44 -z-10 h-2 w-[38rem] rotate-[-18deg] bg-primary/20 sm:w-[50rem]"
        aria-hidden="true"
      />
      <div
        class="pointer-events-none absolute -bottom-5 -left-32 -z-10 h-3 w-[34rem] rotate-[-18deg] bg-primary/30"
        aria-hidden="true"
      />
      <div
        class="pointer-events-none absolute bottom-12 -left-28 -z-10 h-1.5 w-[26rem] rotate-[-18deg] bg-primary/45"
        aria-hidden="true"
      />
      <div
        class="pointer-events-none absolute bottom-24 -left-44 -z-10 h-2 w-[18rem] rotate-[-18deg] bg-primary/25 sm:w-[30rem]"
        aria-hidden="true"
      />

      <div class="mx-auto grid w-full max-w-7xl gap-12 lg:grid-cols-[0.82fr_1.18fr] lg:items-center">
        <div class="grid gap-8">
          <div class="grid gap-4">
            <p class="m-0 font-mono text-xs font-bold uppercase tracking-[0.32em] text-accent">
              journeys
            </p>
            <h2 class="m-0 max-w-[10ch] text-5xl font-black leading-none tracking-normal text-foreground sm:text-6xl lg:text-7xl">
              What made me `ME`
            </h2>
            <p class="m-0 max-w-[54ch] text-base leading-7 text-muted-foreground">
              A compact diary of my life from what I remember until now.
            </p>
          </div>

          <div class="grid gap-4">
            <div class="flex items-center justify-between gap-3">
              <div class="h-px flex-1 bg-border" aria-hidden="true" />
              <div
                v-if="hasJourneys"
                class="flex gap-2"
              >
                <button
                  type="button"
                  class="inline-flex size-9 items-center justify-center rounded-md border border-border bg-background text-foreground transition hover:bg-surface active:-translate-y-[1px]"
                  aria-label="Show previous journey years"
                  @click="moveJourneyTimeline('previous')"
                >
                  <ChevronLeft class="size-4" aria-hidden="true" />
                </button>
                <button
                  type="button"
                  class="inline-flex size-9 items-center justify-center rounded-md border border-border bg-background text-foreground transition hover:bg-surface active:-translate-y-[1px]"
                  aria-label="Show next journey years"
                  @click="moveJourneyTimeline('next')"
                >
                  <ChevronRight class="size-4" aria-hidden="true" />
                </button>
              </div>
            </div>

            <div
              v-if="hasJourneys"
              ref="journeyTimelineRef"
              class="overflow-x-auto pb-3 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
            >
              <div class="flex min-w-max items-start" role="list" aria-label="Journey timeline">
                <div
                  v-for="(item, index) in journeyTimeline"
                  :key="item.year"
                  class="relative grid w-44 shrink-0 gap-4 pr-6 text-left after:absolute after:left-6 after:top-3 after:h-px after:w-[calc(100%-1.5rem)] sm:w-56"
                  :class="index === journeyTimeline.length - 1 ? 'after:bg-border/50' : 'after:bg-border'"
                >
                  <button
                    type="button"
                    class="group relative z-[1] grid w-fit justify-items-start gap-3 text-left outline-none"
                    :aria-pressed="selectedJourneyYear === item.year"
                    @click="selectJourneyYear(item.year)"
                  >
                    <span
                      class="relative size-6 rounded-full border bg-background transition group-hover:border-primary group-hover:bg-primary group-hover:text-background group-active:scale-95"
                      :class="selectedJourneyYear === item.year ? 'border-primary bg-primary text-background' : 'border-border text-foreground'"
                      aria-hidden="true"
                    >
                      <span class="absolute left-1/2 top-1/2 size-2 -translate-x-1/2 -translate-y-1/2 rounded-full bg-current" />
                    </span>
                    <span
                      class="font-mono text-sm font-bold transition"
                      :class="selectedJourneyYear === item.year ? 'text-primary' : 'text-muted-foreground group-hover:text-foreground'"
                    >
                      {{ item.year }}
                    </span>
                  </button>

                  <div
                    class="absolute left-6 right-6 top-3 h-px"
                    aria-label="Journey months"
                  >
                    <button
                      v-for="month in item.months"
                      :key="`${item.year}-${month}`"
                      type="button"
                      class="group/month absolute top-1/2 grid size-6 -translate-x-1/2 -translate-y-1/2 place-items-center rounded-full outline-none transition active:scale-95"
                      :style="{ left: getJourneyMonthOffset(month) }"
                      :aria-label="`Show ${month} ${item.year} journeys`"
                      :aria-pressed="selectedJourneyPoint?.year === item.year && selectedJourneyPoint.month === month"
                      @click="selectJourneyPoint({ year: item.year, month })"
                    >
                      <span
                        class="size-2 rounded-full border transition group-hover/month:size-3 group-hover/month:border-primary group-hover/month:bg-primary"
                        :class="selectedJourneyPoint?.year === item.year && selectedJourneyPoint.month === month ? 'size-3 border-primary bg-primary' : 'border-muted-foreground/50 bg-background'"
                        aria-hidden="true"
                      />
                      <span
                        v-if="selectedJourneyPoint?.year === item.year && selectedJourneyPoint.month === month"
                        class="absolute left-1/2 top-6 -translate-x-1/2 rounded-sm bg-background px-1.5 py-0.5 font-mono text-[10px] font-bold text-primary shadow-soft"
                        aria-hidden="true"
                      >
                        {{ month }}
                      </span>
                      <span class="sr-only">{{ month }}</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div
          v-if="hasJourneys"
          class="grid gap-4 sm:grid-cols-2"
        >
          <article
            v-for="(journey, index) in visibleJourneys"
            :key="`${selectedJourneyYear ?? 'highlight'}-${journey.year}-${journey.topic}`"
            class="group overflow-hidden rounded-md border border-border bg-background shadow-soft transition duration-300 hover:-translate-y-1 hover:border-primary/40"
            :class="[
              index === 0 ? 'sm:row-span-2' : '',
              visibleJourneys.length === 3 && index === 2 ? 'sm:col-start-2' : '',
            ]"
          >
            <div
              class="relative overflow-hidden"
              :class="index === 0 ? 'aspect-[4/5]' : 'aspect-[4/3]'"
            >
              <img
                class="size-full object-cover grayscale transition duration-500 group-hover:scale-[1.04] group-hover:grayscale-0"
                :src="journey.thumbnail"
                :alt="journey.topic"
              >
              <div class="absolute inset-x-0 top-0 flex items-center justify-between p-3">
                <span class="rounded-sm bg-background/90 px-2 py-1 font-mono text-[11px] font-bold text-foreground">
                  {{ journey.year }}
                </span>
                <ArrowUpRight class="size-4 text-background drop-shadow" aria-hidden="true" />
              </div>
            </div>

            <div class="grid gap-4 p-4 sm:p-5">
              <h3 class="m-0 text-xl font-black leading-tight tracking-normal text-foreground">
                {{ journey.topic }}
              </h3>
              <dl class="m-0 grid gap-2 font-mono text-xs font-semibold text-muted-foreground">
                <div class="flex items-center gap-2">
                  <CalendarDays class="size-4 text-accent" aria-hidden="true" />
                  <dt class="sr-only">Date</dt>
                  <dd class="m-0">
                    {{ journey.date }}
                  </dd>
                </div>
                <div class="flex items-center gap-2">
                  <MapPin class="size-4 text-primary" aria-hidden="true" />
                  <dt class="sr-only">Place</dt>
                  <dd class="m-0">
                    {{ journey.place }}
                  </dd>
                </div>
              </dl>
            </div>
          </article>
        </div>
        <p
          v-else
          class="m-0 rounded-md border border-dashed border-border bg-background/88 p-6 text-base font-semibold leading-7 text-muted-foreground shadow-soft sm:p-8"
        >
          no journey for now, the author is adventuring!
        </p>
      </div>
    </section>
  </div>
</template>
