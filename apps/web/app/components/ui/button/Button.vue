<script setup lang="ts">
import { computed, useAttrs } from 'vue'
import { cn } from '~/lib/utils'
import { buttonVariants, type ButtonVariants } from '~/lib/variants/button'

const props = withDefaults(
  defineProps<{
    as?: string
    variant?: ButtonVariants['variant']
    size?: ButtonVariants['size']
  }>(),
  {
    as: 'button',
    variant: 'default',
    size: 'default',
  },
)

const attrs = useAttrs()

const buttonClass = computed(() =>
  cn(buttonVariants({ variant: props.variant, size: props.size }), attrs.class as string),
)
</script>

<template>
  <component :is="as" v-bind="attrs" :class="buttonClass">
    <slot />
  </component>
</template>
