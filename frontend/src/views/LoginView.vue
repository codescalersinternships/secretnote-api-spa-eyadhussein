<script lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import apiClient from '@/APIClient'
import FormInput from '@/components/FormInput.vue'

export default {
  name: 'RegisterView',

  components: {
    FormInput
  },

  setup() {
    const router = useRouter()
    const username = ref('')
    const password = ref('')

    const login = async () => {
      try {
        await apiClient.post('/auth/login', {
          username: username.value,
          password: password.value
        })
        router.push('/notes')
      } catch (error) {
        console.error('Error during logging in', error)
      }
    }

    const updateValue = (inputName: string, value: string) => {
      if (inputName === 'username') {
        username.value = value
      } else if (inputName === 'password') {
        password.value = value
      }
    }

    return {
      username,
      password,
      login,
      updateValue
    }
  }
}
</script>

<template>
  <div class="flex justify-center items-center mt-52">
    <div class="w-full max-w-xs">
      <form @submit.prevent="login" class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
        <FormInput
          id="username"
          label="Username"
          placeholder="username"
          @update-value="updateValue"
        />

        <FormInput
          id="password"
          label="Password"
          type="password"
          placeholder="*******"
          @update-value="updateValue"
        />

        <div class="flex items-center justify-between">
          <button
            class="bg-blue-600 hover:bg-black text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
            type="submit"
          >
            Sign In
          </button>
          <RouterLink
            class="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800"
            to="/register"
          >
            Forgot Password?
          </RouterLink>
        </div>
      </form>
    </div>
  </div>
</template>
