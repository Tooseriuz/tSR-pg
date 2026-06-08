<script setup lang="ts">
import { ChevronLeft, ChevronRight } from '@lucide/vue'
import type { JourneyMonth, JourneyPoint, JourneyTimelineItem, JourneyYear } from '../types/journey'

defineProps<{
  hasJourneys: boolean
  items: JourneyTimelineItem[]
  selectedYear: JourneyYear | null
  selectedPoint: JourneyPoint | null
}>()

const emit = defineEmits<{
  selectYear: [year: JourneyYear]
  selectPoint: [point: JourneyPoint]
}>()

const timelineRef = ref<HTMLElement | null>(null)
const monthOrder: JourneyMonth[] = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']

function getMonthOffset(months: JourneyMonth[], month: JourneyMonth) {
  const orderedMonths = [...months].sort((previousMonth, nextMonth) =>
    monthOrder.indexOf(previousMonth) - monthOrder.indexOf(nextMonth),
  )
  const monthIndex = orderedMonths.indexOf(month)

  if (monthIndex < 0) {
    return '50%'
  }

  return `${((monthIndex + 1) / (orderedMonths.length + 1)) * 100}%`
}

function moveTimeline(direction: 'previous' | 'next') {
  const timeline = timelineRef.value

  if (!timeline) {
    return
  }

  timeline.scrollBy({
    left: direction === 'next' ? timeline.clientWidth * 0.72 : timeline.clientWidth * -0.72,
    behavior: 'smooth',
  })
}
</script>

<template>
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
          @click="moveTimeline('previous')"
        >
          <ChevronLeft class="size-4" aria-hidden="true" />
        </button>
        <button
          type="button"
          class="inline-flex size-9 items-center justify-center rounded-md border border-border bg-background text-foreground transition hover:bg-surface active:-translate-y-[1px]"
          aria-label="Show next journey years"
          @click="moveTimeline('next')"
        >
          <ChevronRight class="size-4" aria-hidden="true" />
        </button>
      </div>
    </div>

    <div
      v-if="hasJourneys"
      ref="timelineRef"
      class="overflow-x-auto pb-3 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
    >
      <div class="flex min-w-max items-start" role="list" aria-label="Journey timeline">
        <div
          v-for="(item, index) in items"
          :key="item.year"
          class="relative grid w-44 shrink-0 gap-4 pr-6 text-left after:absolute after:left-6 after:top-3 after:h-px after:w-[calc(100%-1.5rem)] sm:w-56"
          :class="index === items.length - 1 ? 'after:bg-border/50' : 'after:bg-border'"
        >
          <button
            type="button"
            class="group relative z-[1] grid w-fit justify-items-start gap-3 text-left outline-none"
            :aria-pressed="selectedYear === item.year"
            @click="emit('selectYear', item.year)"
          >
            <span
              class="relative size-6 rounded-full border bg-background transition group-hover:border-primary group-hover:bg-primary group-hover:text-background group-active:scale-95"
              :class="selectedYear === item.year ? 'border-primary bg-primary text-background' : 'border-border text-foreground'"
              aria-hidden="true"
            >
              <span class="absolute left-1/2 top-1/2 size-2 -translate-x-1/2 -translate-y-1/2 rounded-full bg-current" />
            </span>
            <span
              class="font-mono text-sm font-bold transition"
              :class="selectedYear === item.year ? 'text-primary' : 'text-muted-foreground group-hover:text-foreground'"
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
              class="group/month absolute top-1/2 grid -translate-x-1/2 -translate-y-1/2 place-items-center rounded-full outline-none transition active:scale-95"
              :style="{ left: getMonthOffset(item.months, month) }"
              :aria-label="`Show ${month} ${item.year} journeys`"
              :aria-pressed="selectedPoint?.year === item.year && selectedPoint.month === month"
              @click="emit('selectPoint', { year: item.year, month })"
            >
              <span
                class="size-2 rounded-full border transition group-hover/month:size-3 group-hover/month:border-primary group-hover/month:bg-primary"
                :class="selectedPoint?.year === item.year && selectedPoint.month === month ? 'size-3 border-primary bg-primary' : 'border-muted-foreground/50 bg-background'"
                aria-hidden="true"
              />
              <span
                v-if="selectedPoint?.year === item.year && selectedPoint.month === month"
                class="absolute left-1/2 top-6 -translate-x-1/2 rounded-sm bg-background font-mono text-[10px] font-bold text-primary shadow-soft"
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
</template>
