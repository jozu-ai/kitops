import type { DirectiveBinding } from 'vue'

// Common tracking function to keep it DRY
function trackEvent(action, binding: DirectiveBinding) {
  // Category, action and label are needed, the rest is custom, per implementation
  const { category, label, ...values } = binding

  return (event: InputEvent | Event) => {

    const parameters = {
      event_category: category,
      event_label: label,
      value: event?.target?.value || values?.value,
      ...values,
    }

    // Call the gtag function to track the event
    // window.gtag('event', action, parameters);
  }
}

export default {
  // When mounted, add the bindings...
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    if (typeof binding.value !== 'object') {
      console.error('Directive value must be an object')
      return;
    }

    // Attach the given event/action
    const action = binding.value.action || 'click'
    el.addEventListener(action, trackEvent(action, binding.value));
  },

  beforeUnmount(el: HTMLElement, binding: DirectiveBinding) {
    // Remove the listener when destroying the directive instance
    const action = binding.value.action || 'click'
    el.removeEventListener(action, trackEvent(action, binding.value));
  }
};
