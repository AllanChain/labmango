<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { CurrentConfig, ListLabs, RefreshLabs } from '@wails/go/main/App'
import type { main as goModels } from '@wails/go/models'
import { EventsOn } from '@wails/runtime/runtime'
import NewLab from './components/NewLab.vue'
import LabConfig from '@/components/LabConfig.vue'
import LabCard from '@/components/LabCard.vue'

const currentConfig = ref<goModels.LabManConfig | null>(null)
const labs = ref<string[]>([])

onMounted(async () => {
  currentConfig.value = await CurrentConfig()
  labs.value = await ListLabs()
})

EventsOn('config-update', async (config) => {
  currentConfig.value = config
})

EventsOn('labs-refresh', (newLabs) => {
  labs.value = newLabs
})
</script>

<template>
  <main>
    <h1 class="text-center">
      Labs Management
    </h1>
    <div class="text-right mb-2">
      <button class="text-green-500 border-rounded-full" @click="RefreshLabs">
        <i class="i-mdi-refresh" />
        Refresh Labs
      </button>
    </div>
    <NewLab />
    <LabCard v-for="lab in labs" :key="lab" :name="lab" />
    <LabConfig v-if="currentConfig" :config="currentConfig" />
  </main>
</template>

  <style>

  </style>
