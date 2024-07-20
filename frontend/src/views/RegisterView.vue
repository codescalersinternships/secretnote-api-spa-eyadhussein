<script lang="ts">
import { RouterLink } from 'vue-router'
import apiClient from '@/APIClient'
import FormInput from '@/components/FormInput.vue'
export default {
  name: 'RegisterView',
  data() {
    return {
      username: '',
      email: '',
      password: '',
      password_confirmation: ''
    }
  },
  methods: {
    async register() {
      try {
        await apiClient.post('/auth/register', {
          username: this.username,
          email: this.email,
          password: this.password,
          password_confirmation: this.password_confirmation
        })
      } catch (error) {
        console.error('error during registration', error)
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
    <div class="w-full max-w-sm">
      <form @submit.prevent="register" class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
        <FormInput
          id="username"
          label="Username"
          placeholder="username"
          @update-value="updateValue"
        />

        <FormInput
          id="email"
          label="Email"
          placeholder="email@example.com"
          @update-value="updateValue"
        />

        <FormInput
          id="password"
          label="Password"
          type="password"
          placeholder="*******"
          @update-value="updateValue"
        />
        <FormInput
          id="password_confirmation"
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
            Register
          </button>
          <RouterLink
            class="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800"
            to="/login"
          >
            Already have an account?
          </RouterLink>
        </div>
      </form>
    </div>
  </div>
</template>
