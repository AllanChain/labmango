<script lang="ts" setup>
import { ExploreLab } from '@wails/go/main/App'
const props = defineProps<{ name: string }>()

const exploreLab = () => {
  ExploreLab(props.name)
}
</script>

<template>
  <div class="p-4 bg-slate-800" border="1 rounded-md slate-900" shadow="md">
    <div class="p-2">
      {{ name }}
    </div>
    <button class="text-green-500" @click="exploreLab">
      <i class="i-mdi-folder-open" />
      Explore
    </button>
  </div>
</template>

<style>

</style>
