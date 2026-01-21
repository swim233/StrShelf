<script setup lang="ts">
const props = defineProps(['icon'])
</script>

<template>
  <div class="icon">
    <div v-if="icon === 'Success'">
      <svg
        color="green"
        xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink"
        viewBox="0 0 32 32"
      >
        <path
          d="M14 21.414l-5-5.001L10.413 15L14 18.586L21.585 11L23 12.415l-9 8.999z"
          fill="currentColor"
        ></path>
        <path
          d="M16 2a14 14 0 1 0 14 14A14 14 0 0 0 16 2zm0 26a12 12 0 1 1 12-12a12 12 0 0 1-12 12z"
          fill="currentColor"
        ></path>
      </svg>
    </div>
    <div v-if="icon === 'Error'">
      <svg
        color="red"
        xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink"
        viewBox="0 0 32 32"
      >
        <path
          d="M16 2C8.2 2 2 8.2 2 16s6.2 14 14 14s14-6.2 14-14S23.8 2 16 2zm0 26C9.4 28 4 22.6 4 16S9.4 4 16 4s12 5.4 12 12s-5.4 12-12 12z"
          fill="currentColor"
        ></path>
        <path
          d="M21.4 23L16 17.6L10.6 23L9 21.4l5.4-5.4L9 10.6L10.6 9l5.4 5.4L21.4 9l1.6 1.6l-5.4 5.4l5.4 5.4z"
          fill="currentColor"
        ></path>
      </svg>
    </div>
    <div v-if="icon === 'Info'">
      <svg
        color="rgb(23, 80, 193)"
        xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink"
        viewBox="0 0 32 32"
      >
        <path d="M17 22v-8h-4v2h2v6h-3v2h8v-2h-3z" fill="currentColor"></path>
        <path
          d="M16 8a1.5 1.5 0 1 0 1.5 1.5A1.5 1.5 0 0 0 16 8z"
          fill="currentColor"
        ></path>
        <path
          d="M16 30a14 14 0 1 1 14-14a14 14 0 0 1-14 14zm0-26a12 12 0 1 0 12 12A12 12 0 0 0 16 4z"
          fill="currentColor"
        ></path>
      </svg>
    </div>
  </div>
</template>

<style lang="css">
.icon {
  width: 24px;
  height: 24px;
  display: inline-block;
  position: relative;
  top: 4px;
}
</style>
