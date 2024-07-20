<script lang="ts">
import apiClient from '@/APIClient'
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import NoteCard from '@/components/NoteCard.vue'

import type Note from '@/types/Note'

export default {
  name: 'NoteListView',

  components: {
    NoteCard
  },

  setup() {
    const notes = ref<Note[]>([])
    const err = ref('')
    const router = useRouter()

    const fetchNotes = async () => {
      try {
        const response = await apiClient.get('users/notes')
        notes.value = response.data
      } catch (error) {
        console.error('Error fetching notes', error)
        err.value = 'Error fetching notes too many requests'
      }
    }

    const handleLogout = async () => {
      try {
        await apiClient.post('auth/logout')
        router.push('/login')
      } catch (error) {
        console.error('Error logging out', error)
        err.value = 'Error logging out'
      }
    }

    const goToCreateNote = () => {
      router.push('/notes/create')
    }

    onMounted(fetchNotes)

    return {
      notes,
      err,
      handleLogout,
      goToCreateNote
    }
  }
}
</script>

<template>
  <div class="flex flex-col items-center mt-10">
    <div class="w-full max-w-2xl">
      <div class="flex justify-between mb-4">
        <button
          @click="goToCreateNote"
          class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
        >
          Create Note
        </button>
        <button
          @click="handleLogout"
          class="bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
        >
          Logout
        </button>
      </div>
      <h2 class="text-2xl mb-4">Your Notes</h2>
      <div
        v-if="err"
        class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative"
        role="alert"
      >
        <span class="block sm:inline">{{ err }}</span>
      </div>

      <div v-else>
        <div v-if="notes.length === 0" class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
          <p class="text-gray-700">You have no notes.</p>
        </div>
        <div v-else class="space-y-4">
          <NoteCard v-for="note in notes" :key="note.id" :note="note" />
        </div>
      </div>
    </div>
  </div>
</template>
