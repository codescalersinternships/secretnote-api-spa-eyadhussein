<script lang="ts">
import { RouterLink } from 'vue-router'
import apiClient from '@/APIClient'
import FormInput from '@/components/FormInput.vue'
export default {
  name: 'RegisterView',
  data() {
    return {
      username: '',
      password: ''
    }
  },
  methods: {
    async login() {
      try {
        await apiClient.post('/auth/login', {
          username: this.username,
          password: this.password
        })
      } catch (error) {
        console.error('error during logging in', error)
      }
    },
    updateValue(inputName: string, value: string) {
      if (inputName in this) {
        ;(this as any)[inputName] = value
      }
    }
  },
  components: {
    FormInput,
    RouterLink
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
