<template>
<div class="discord-banner p-6 my-12 flex gap-6">
  <div class="grow-0 shrink-0 basis-6">
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
      <path fill="none" d="M0 0h24v24H0z"></path><path d="M21 10.063V4C21 3.44772 20.5523 3 20 3H19C17.0214 4.97864 13.3027 6.08728 11 6.61281V17.3872C13.3027 17.9127 17.0214 19.0214 19 21H20C20.5523 21 21 20.5523 21 20V13.937C21.8626 13.715 22.5 12.9319 22.5 12 22.5 11.0681 21.8626 10.285 21 10.063ZM5 7C3.89543 7 3 7.89543 3 9V15C3 16.1046 3.89543 17 5 17H6L7 22H9V7H5Z"></path>
    </svg>
  </div>

  <div class="flex-1 space-y-4">
    <strong>Need help getting your ML projects to production?</strong>

    <div>Talk with our team about how KitOps can help</div>

    <div>
      <button class="kit-button kit-button-small kit-button-cornflower">Join our Discord</button>
    </div>
  </div>
</div>
</template>

<style>
.discord-banner {
  background: var(--vp-custom-block-tip-bg);
}
</style>
