<script lang="ts" setup>
import { DeleteLab, ExploreLab } from '@wails/go/main/App'
const props = defineProps<{ name: string }>()

const isDev = import.meta.env.DEV
</script>

<template>
  <div class="p-4 my-2 bg-slate-800" border="1 rounded-md slate-900" shadow="md">
    <div class="p-2">
      {{ name }}
    </div>
    <div class="flex">
      <button class="text-green-500" @click="ExploreLab(name)">
        <i class="i-mdi-folder-open" />
        Explore
      </button>
      <div class="flex-1" />
      <button v-if="isDev" class="text-red-600" @click="DeleteLab(name)">
        <i class="i-mdi-delete" />
        Delete Lab
      </button>
    </div>
  </div>
</template>

<style>

</style>
