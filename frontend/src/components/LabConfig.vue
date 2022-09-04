<script lang="ts" setup>
import type { main as goModels } from '@wails/go/models'
import { ChangeJupyterTemplate, ChangeLabDir, ChangeLyXTemplate } from '@wails/go/main/App'
defineProps<{
  config: goModels.LabManConfig
}>()
</script>

<template>
  <div class="py-2 px-4">
    <div>
      <span>Labs directory:</span>
      {{ config.labDir }}
      <button class="text-slate-400" @click="ChangeLabDir">
        <i class="i-mdi-folder-arrow-left-right" />
        Change
      </button>
    </div>
    <div>
      <span>Jupyter template:</span>
      {{ config.jupyterTemplate }}
      <button class="text-slate-400" @click="ChangeJupyterTemplate">
        <i class="i-mdi-file-arrow-left-right" />
        Change
      </button>
    </div>
    <div>
      <span>LyX template:</span>
      {{ config.lyxTemplate }}
      <button class="text-slate-400" @click="ChangeLyXTemplate">
        <i class="i-mdi-file-arrow-left-right" />
        Change
      </button>
    </div>
  </div>
</template>
