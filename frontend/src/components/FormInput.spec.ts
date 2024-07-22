import { describe, it, expect, beforeEach } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import FormInput from '@/components/FormInput.vue'

describe('FormInput.vue', () => {
  let wrapper: ReturnType<typeof shallowMount>

  beforeEach(() => {
    wrapper = shallowMount(FormInput, {
      props: {
        id: 'username',
        label: 'Username',
        value: '',
        placeholder: 'Enter username',
        error: ''
      }
    })
  })

  it('renders the label correctly', () => {
    expect(wrapper.find('label').text()).toBe('Username')
  })

  it('emits an update-value event on input', async () => {
    const input = wrapper.find('input')
    await input.setValue('new value')
    expect(wrapper.emitted('update-value')).toBeTruthy()
    expect(wrapper.emitted('update-value')![0]).toEqual(['username', 'new value'])
  })

  it('displays an error message when error prop is passed', async () => {
    await wrapper.setProps({ error: 'This field is required' })
    expect(wrapper.find('span').text()).toBe('This field is required')
  })
})
