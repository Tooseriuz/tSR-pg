<script setup lang="ts">
import { ChevronDown, ChevronRight, Mail } from '@lucide/vue'
import githubIcon from '../assets/icons/github.svg'
import type { components } from '../../types/dto/openapi'
import type { Journey, JourneyMonth, JourneyPoint, JourneyYear } from '../types/journey'

type ApiJourney = components['schemas']['Journey']

const isHeroCollapsed = ref(false)
const hasJourneyIntroAnimated = ref(false)
const brandMarkRef = ref<HTMLElement | null>(null)
const journeysSectionRef = ref<HTMLElement | null>(null)
const journeyRailRef = ref<HTMLElement | null>(null)
const selectedJourneyYear = ref<JourneyYear | null>(null)
const selectedJourneyPoint = ref<JourneyPoint | null>(null)
const focusedJourneyCardIndex = ref(0)
let brandMarkObserver: IntersectionObserver | null = null
let journeyIntroObserver: IntersectionObserver | null = null

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
      id: journey.id,
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
const canShowPreviousJourneys = computed(() =>
  visibleJourneys.value.length > 1 && focusedJourneyCardIndex.value > 0,
)
const canShowMoreJourneys = computed(() =>
  visibleJourneys.value.length > 1 && focusedJourneyCardIndex.value < visibleJourneys.value.length - 1,
)

function updateJourneyRailOverflow() {
  const rail = journeyRailRef.value

  if (!rail || visibleJourneys.value.length <= 1) {
    focusedJourneyCardIndex.value = 0
    return
  }

  focusedJourneyCardIndex.value = Math.min(
    visibleJourneys.value.length - 1,
    Math.max(0, Math.round(rail.scrollLeft / getJourneyRailScrollDistance(rail))),
  )
}

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

function getJourneyRailScrollDistance(rail: HTMLElement) {
  const firstCard = rail.firstElementChild

  if (!firstCard) {
    return rail.clientWidth
  }

  const gap = Number.parseFloat(window.getComputedStyle(rail).columnGap)

  return firstCard.getBoundingClientRect().width + (Number.isNaN(gap) ? 0 : gap)
}

function showMoreJourneys() {
  const rail = journeyRailRef.value

  if (!rail) {
    return
  }

  rail.scrollTo({
    left: getJourneyRailScrollDistance(rail) * (focusedJourneyCardIndex.value + 1),
    behavior: 'smooth',
  })
}

function showPreviousJourneys() {
  const rail = journeyRailRef.value

  if (!rail) {
    return
  }

  rail.scrollTo({
    left: getJourneyRailScrollDistance(rail) * (focusedJourneyCardIndex.value - 1),
    behavior: 'smooth',
  })
}

function getJourneyCardRevealClass(index: number) {
  const delayClasses = ['delay-150', 'delay-200', 'delay-300', 'delay-[400ms]', 'delay-500']

  return [
    hasJourneyIntroAnimated.value
      ? 'translate-y-0 opacity-100'
      : 'translate-y-10 opacity-0',
    delayClasses[Math.min(index, delayClasses.length - 1)],
  ]
}

function scrollToJourneys(behavior: ScrollBehavior = 'smooth') {
  document.getElementById('journeys')?.scrollIntoView({ behavior })
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
  if (window.location.hash === '#journeys') {
    requestAnimationFrame(() => scrollToJourneys('auto'))
  }

  const brandMark = brandMarkRef.value

  if (brandMark) {
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
  }

  const journeysSection = journeysSectionRef.value

  if (journeysSection) {
    journeyIntroObserver = new IntersectionObserver(
      ([entry]) => {
        if (!entry?.isIntersecting) {
          return
        }

        hasJourneyIntroAnimated.value = true
        journeyIntroObserver?.disconnect()
        journeyIntroObserver = null
      },
      {
        threshold: 0.28,
      },
    )

    journeyIntroObserver.observe(journeysSection)
  }

  updateJourneyRailOverflow()
  window.addEventListener('resize', updateJourneyRailOverflow)
})

onUnmounted(() => {
  brandMarkObserver?.disconnect()
  journeyIntroObserver?.disconnect()
  window.removeEventListener('resize', updateJourneyRailOverflow)
})

watch(visibleJourneys, async () => {
  await nextTick()

  if (journeyRailRef.value) {
    journeyRailRef.value.scrollLeft = 0
  }

  updateJourneyRailOverflow()
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
          <HeroTypewriter />
        </div>
      </main>
      <a
        href="#journeys"
        class="scroll-cue absolute bottom-24 left-1/2 grid size-12 place-items-center rounded-full border border-border bg-background text-foreground shadow-soft transition hover:bg-surface active:-translate-y-[1px]"
        aria-label="Scroll to journeys"
        @click.prevent="scrollToJourneys()"
      >
        <ChevronDown class="size-5" aria-hidden="true" />
      </a>
    </section>

    <section
      id="journeys"
      ref="journeysSectionRef"
      class="relative isolate grid min-h-[100dvh] content-center overflow-hidden border-t border-border bg-background py-20 pl-6 pr-0 sm:pl-8 lg:pl-[20%]"
    >
      <div class="grid w-full gap-12 lg:grid-cols-[minmax(18rem,28rem)_minmax(0,1fr)] lg:items-center">
        <div class="grid gap-8 pr-6 sm:pr-8 lg:pr-0 transition duration-700 ease-out motion-reduce:translate-x-0 motion-reduce:opacity-100 motion-reduce:transition-none"
          :class="hasJourneyIntroAnimated ? 'translate-x-0 opacity-100' : '-translate-x-12 opacity-0'">
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

          <JourneyTimeline
            :has-journeys="hasJourneys"
            :items="journeyTimeline"
            :selected-year="selectedJourneyYear"
            :selected-point="selectedJourneyPoint"
            @select-year="selectJourneyYear"
            @select-point="selectJourneyPoint"
          />
        </div>

        <div
          v-if="hasJourneys"
          class="relative -ml-6 overflow-hidden sm:ml-0"
        >
          <button
            v-if="canShowPreviousJourneys"
            type="button"
            class="absolute inset-y-0 left-0 z-[1] grid w-14 place-items-center text-foreground outline-none transition hover:text-primary sm:w-32 sm:bg-gradient-to-r sm:from-background sm:via-background/90 sm:to-background/10"
            aria-label="Show previous journey cards"
            @click="showPreviousJourneys"
          >
            <span class="inline-flex items-center gap-2 rounded-md border border-border bg-background px-3 py-2 font-mono text-xs font-bold shadow-soft transition hover:border-primary">
              <ChevronRight class="size-5 rotate-180" aria-hidden="true" />
            </span>
          </button>

          <div
            ref="journeyRailRef"
            class="flex snap-x gap-4 overflow-x-auto px-[11vw] pb-2 [scrollbar-width:none] sm:px-16 [&::-webkit-scrollbar]:hidden"
            aria-label="Journey cards"
            @scroll="updateJourneyRailOverflow"
          >
            <JourneyCard
              v-for="(journey, index) in visibleJourneys"
              :key="`${selectedJourneyYear ?? selectedJourneyPoint?.year ?? 'highlight'}-${journey.id}`"
              class="w-[78vw] shrink-0 snap-center transition duration-700 ease-out motion-reduce:translate-y-0 motion-reduce:opacity-100 motion-reduce:transition-none sm:w-72 sm:snap-start lg:w-80"
              :class="getJourneyCardRevealClass(index)"
              :journey="journey"
            />
          </div>

          <button
            v-if="canShowMoreJourneys"
            type="button"
            class="absolute inset-y-0 right-0 grid w-14 place-items-center text-foreground outline-none transition hover:text-primary sm:w-32 sm:bg-gradient-to-l sm:from-background sm:via-background/90 sm:to-background/10"
            aria-label="Show more journey cards"
            @click="showMoreJourneys"
          >
            <span class="inline-flex items-center gap-2 rounded-md border border-border bg-background px-3 py-2 font-mono text-xs font-bold shadow-soft transition hover:border-primary">
              <ChevronRight class="size-5" aria-hidden="true" />
            </span>
          </button>
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
