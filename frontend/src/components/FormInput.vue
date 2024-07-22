<script lang="ts">
export default {
  name: 'FormInput',
  props: {
    id: {
      type: String,
      required: true
    },
    type: {
      type: String,
      default: 'text'
    },
    label: {
      type: String,
      required: true
    },
    value: {
      type: String,
      default: ''
    },
    placeholder: {
      type: String,
      default: ''
    },
    error: {
      type: String,
      default: ''
    }
  },
  emits: ['update-value'],
  methods: {
    onInput(event: Event) {
      const target = event.target as HTMLInputElement
      this.$emit('update-value', this.id, target.value)
    }
  }
}
</script>

<template>
  <div class="mb-4">
    <label :for="id" class="block text-gray-700 text-sm font-bold mb-2">{{ label }}</label>
    <input
      :id="id"
      :type="type"
      :value="value"
      @input="onInput"
      class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
      :placeholder="placeholder"
    />
    <span v-if="error" class="text-red-500 text-xs italic mt-2">{{ error }}</span>
  </div>
</template>
