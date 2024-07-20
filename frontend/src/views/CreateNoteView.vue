<script lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import apiClient from '@/APIClient'
import FormInput from '@/components/FormInput.vue'
import { AxiosError } from 'axios'

export default {
  name: 'CreateNoteView',

  components: {
    FormInput
  },

  setup() {
    const router = useRouter()
    const title = ref('')
    const content = ref('')
    const maxViews = ref(1)
    const expiresAt = ref('')

    const createNote = async () => {
      try {
        const formattedDate = formatDate(expiresAt.value)
        await apiClient.post('/notes', {
          title: title.value,
          content: content.value,
          max_views: maxViews.value,
          expires_at: formattedDate
        })
        router.push('/notes')
      } catch (error) {
        if (error instanceof AxiosError) {
          console.error('Error during creating note', error.response?.data)
        } else {
          console.error('Error during creating note', error)
        }
      }
    }

    const formatDate = (dateString: string): string => {
      const date = new Date(dateString)
      return date.toISOString()
    }

    const updateValue = (inputName: string, value: string) => {
      if (inputName === 'maxViews') {
        maxViews.value = parseInt(value)
      } else if (inputName === 'title') {
        title.value = value
      } else if (inputName === 'content') {
        content.value = value
      } else if (inputName === 'expiresAt') {
        expiresAt.value = value
      }
    }

    return {
      title,
      content,
      maxViews,
      expiresAt,
      createNote,
      updateValue
    }
  }
}
</script>

<template>
  <div class="flex justify-center items-center mt-52">
    <div class="w-full max-w-md">
      <form @submit.prevent="createNote" class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
        <h2 class="text-2xl mb-4">Create Note</h2>

        <FormInput
          id="title"
          label="Title"
          placeholder="Enter note title"
          @update-value="updateValue"
        />

        <div class="mb-4">
          <label for="content" class="block text-gray-700 text-sm font-bold mb-2">Content</label>
          <textarea
            id="content"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            v-model="content"
          ></textarea>
        </div>

        <FormInput
          id="maxViews"
          label="Max Views"
          placeholder="Enter maximum views"
          type="number"
          @update-value="updateValue"
        />

        <FormInput
          id="expiresAt"
          label="Expiration Date"
          placeholder="Enter expiration date"
          type="date"
          @update-value="updateValue"
        />

        <div class="flex items-center justify-between">
          <button
            class="bg-blue-600 hover:bg-black text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
            type="submit"
          >
            Create Note
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
