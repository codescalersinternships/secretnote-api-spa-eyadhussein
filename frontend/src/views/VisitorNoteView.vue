<script lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import apiClient from '@/APIClient'
import type Note from '@/types/Note'
import { AxiosError } from 'axios'

export default {
  name: 'VisitorNoteView',

  setup() {
    const route = useRoute()
    const note = ref<Note | null>(null)
    const noteId = route.params['id']
    const error = ref('')

    const fetchNote = async () => {
      try {
        const response = await apiClient.get<Note>(`/notes/${noteId}`)
        note.value = response.data
      } catch (err) {
        if (err instanceof AxiosError) {
          error.value = err.response?.data.error || 'Error fetching note'
        } else {
          error.value = 'Error fetching note'
        }
      }
    }

    onMounted(fetchNote)

    return {
      note,
      error
    }
  }
}
</script>

<template>
  <div class="flex justify-center items-center mt-10">
    <div class="w-full max-w-2xl">
      <h2 class="text-2xl mb-4">Your Notes</h2>
      <div
        v-if="error"
        class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4"
        role="alert"
      >
        <span class="block sm:inline">{{ error }}</span>
      </div>
      <div v-else-if="note" class="space-y-4">
        <div class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
          <h3 class="text-xl font-bold mb-2">{{ note.title }}</h3>
          <p class="text-gray-700 mb-4">{{ note.content }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
