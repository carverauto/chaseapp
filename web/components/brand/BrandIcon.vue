<template>
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 120">
    <defs>
      <linearGradient
        id="linear-gradient"
        x1="-730.45"
        y1="-181.06"
        x2="-731.15"
        y2="-182.04"
        gradientTransform="translate(73144 21862) scale(100 120)"
        gradientUnits="userSpaceOnUse"
      >
        <stop offset="0" :stop-color="primary" />
        <stop offset="1" :stop-color="secondary" />
      </linearGradient>
    </defs>
    <path style="fill: url(#linear-gradient);" d="M20,10,48.4,41.6H104a15.94,15.94,0,0,1,16,15.8v56.8A15.88,15.88,0,0,1,104,130H36a15.94,15.94,0,0,1-16-15.8ZM90.3,67H49.7A4.69,4.69,0,0,0,45,71.7V117l14.3-13.2h31A4.69,4.69,0,0,0,95,99.1V71.6A4.67,4.67,0,0,0,90.3,67Z" transform="translate(-20 -10)" />
  </svg>
</template>

<script>
export default {
  props: {
    primary: {
      type: String,
      required: false,
      default: '#FB605A',
    },
    secondary: {
      type: String,
      required: false,
      default: '#974BD3',
    },
  },
}
</script>
