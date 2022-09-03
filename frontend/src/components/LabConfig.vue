<script lang="ts" setup>
import type { main as goModels } from '@wails/go/models'
import { ChangeLabDir } from '@wails/go/main/App'
defineProps<{
  config: goModels.LabManConfig
}>()
</script>

<template>
  <div class="py-2 px-4">
    <span>Labs directory:</span>
    {{ config.labDir }}
    <button class="text-slate-400" @click="ChangeLabDir">
      <i class="i-mdi-folder-arrow-left-right" />
      Change
    </button>
  </div>
</template>
