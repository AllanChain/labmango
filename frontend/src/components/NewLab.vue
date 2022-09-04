<script lang="ts" setup>
import { CreateLab } from '@wails/go/main/App'
import { ref } from 'vue'

const labName = ref('')

const createLab = () => {
  CreateLab(labName.value)
  labName.value = ''
}
</script>

<template>
  <div class="py-2">
    <input v-model="labName" class="bg-slate-800" @keydown.enter="createLab">
    <button class="text-green-500" @click="createLab">
      <i class="i-mdi-plus" />
      Create New Lab
    </button>
  </div>
</template>
