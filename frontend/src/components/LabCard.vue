<script lang="ts" setup>
import { DeleteLab, EditReport, ExploreLab, LaunchLab } from '@wails/go/main/App'
const props = defineProps<{ name: string; jlabRunning: boolean }>()

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
      <button class="text-blue-400" @click="EditReport(name)">
        <i class="i-mdi-file-edit" />
        Edit Report
      </button>
      <button v-if="!jlabRunning" class="text-orange-500" @click="LaunchLab(name)">
        <i class="i-mdi-rocket-launch" />
        JupyterLab
      </button>
      <button v-else disabled>
        <i class="i-mdi-cog-play" />
        JLab Running
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
