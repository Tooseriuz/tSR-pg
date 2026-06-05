<script setup lang="ts">
import { ArrowUpRight, CalendarDays, MapPin } from '@lucide/vue'
import type { Journey } from '../types/journey'

defineProps<{
  journey: Journey
}>()
</script>

<template>
  <NuxtLink
    :to="`/journey/${journey.id}`"
    class="group overflow-hidden rounded-md border border-border bg-background shadow-soft transition duration-300 hover:-translate-y-1 hover:border-primary/40"
  >
    <div class="relative aspect-[4/3] overflow-hidden">
      <img
        v-if="journey.thumbnail"
        class="size-full object-cover transition duration-500 group-hover:scale-[1.04]"
        :src="journey.thumbnail"
        :alt="journey.topic"
      >
      <div
        v-else
        class="code-pattern size-full bg-surface"
        aria-hidden="true"
      />
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
  </NuxtLink>
</template>
