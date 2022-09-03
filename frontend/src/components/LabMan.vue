<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { CurrentConfig, ListLabs } from '@wails/go/main/App'
import type { main as goModels } from '@wails/go/models'
import { EventsOn } from '@wails/runtime/runtime'
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
  labs.value = await ListLabs()
})
</script>

<template>
  <main>
    <h1 class="text-center">
      Labs Management
    </h1>
    <LabCard v-for="lab in labs" :key="lab" :name="lab" />
    <LabConfig v-if="currentConfig" :config="currentConfig" />
  </main>
</template>

<style>

</style>
