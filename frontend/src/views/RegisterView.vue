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
    const email = ref('')
    const password = ref('')
    const password_confirmation = ref('')
    const errors = ref({
      username: '',
      email: '',
      password: '',
      password_confirmation: '',
      general: ''
    })

    const validate = () => {
      let valid = true
      if (!username.value) {
        errors.value.username = 'Username is required'
        valid = false
      } else {
        errors.value.username = ''
      }

      if (!email.value) {
        errors.value.email = 'Email is required'
        valid = false
      } else if (!/\S+@\S+\.\S+/.test(email.value)) {
        errors.value.email = 'Email is invalid'
        valid = false
      } else {
        errors.value.email = ''
      }

      if (!password.value) {
        errors.value.password = 'Password is required'
        valid = false
      } else {
        errors.value.password = ''
      }

      if (!password_confirmation.value) {
        errors.value.password_confirmation = 'Password confirmation is required'
        valid = false
      } else if (password_confirmation.value !== password.value) {
        errors.value.password_confirmation = 'Passwords do not match'
        valid = false
      } else {
        errors.value.password_confirmation = ''
      }

      return valid
    }

    const register = async () => {
      if (!validate()) {
        return
      }

      try {
        await apiClient.post('/auth/register', {
          username: username.value,
          email: email.value,
          password: password.value,
          password_confirmation: password_confirmation.value
        })
        router.push('/notes')
      } catch (error) {
        console.error('Error during registration:', error)
        errors.value.general = 'Registration failed'
      }
    }

    const updateValue = (inputName: string, value: string) => {
      if (inputName === 'username') {
        username.value = value
      } else if (inputName === 'email') {
        email.value = value
      } else if (inputName === 'password') {
        password.value = value
      } else if (inputName === 'password_confirmation') {
        password_confirmation.value = value
      }
    }

    return {
      username,
      email,
      password,
      password_confirmation,
      errors,
      register,
      updateValue
    }
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
          :error="errors.username"
          @update-value="updateValue"
        />

        <FormInput
          id="email"
          label="Email"
          placeholder="email@example.com"
          :error="errors.email"
          @update-value="updateValue"
        />

        <FormInput
          id="password"
          label="Password"
          type="password"
          placeholder="*******"
          :error="errors.password"
          @update-value="updateValue"
        />

        <FormInput
          id="password_confirmation"
          label="Confirm Password"
          type="password"
          placeholder="*******"
          :error="errors.password_confirmation"
          @update-value="updateValue"
        />

        <div v-if="errors.general" class="text-red-500 text-xs italic mt-2">
          {{ errors.general }}
        </div>

        <div class="flex items-center justify-between mt-4">
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
