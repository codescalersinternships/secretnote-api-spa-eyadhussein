<script lang="ts">
import { defineComponent, type PropType } from 'vue'
import type Note from '@/types/Note'

export default defineComponent({
  name: 'NoteCard',
  props: {
    note: {
      type: Object as PropType<Note>,
      required: true
    }
  },
  setup() {
    const formatDate = (dateString: string): string => {
      const date = new Date(dateString)
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    return {
      formatDate,
      currentPath: window.location.origin
    }
  }
})
</script>

<template>
  <div class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
    <h3 class="text-xl font-bold mb-2">{{ note.title }}</h3>
    <p class="text-gray-700 mb-4">{{ note.content }}</p>
    <p class="text-gray-500 text-sm">
      <span class="font-bold">Current Views:</span> {{ note.current_views }}
    </p>
    <p class="text-gray-500 text-sm">
      <span class="font-bold">Max Views:</span> {{ note.max_views }}
    </p>
    <p class="text-gray-500 text-sm">
      <span class="font-bold">Expires At:</span> {{ formatDate(note.expires_at) }}
    </p>
    <p class="text-gray-500 text-sm">
      <span class="font-bold">Unique URL:</span>
      <RouterLink
        class="font-medium text-blue-600 dark:text-blue-500 hover:underline"
        :to="`notes/${note.id}`"
        >{{ `${currentPath}/notes/${note.id}` }}</RouterLink
      >
    </p>
  </div>
</template>
