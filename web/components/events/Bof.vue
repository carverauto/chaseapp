<template>
  <ul class="-mb-8">
    <li>
      <div class="relative pb-8">
        <span class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200" aria-hidden="true" />
        <div class="relative flex space-x-3">
          <div>
            <span class="h-8 w-8 rounded-full bg-gray-400 flex items-center justify-center ring-8 ring-white">
              <img class="h-8 w-8 rounded-full bg-gray-400 flex items-center justify-center ring-8 ring-white" src="bof.png" alt="Birds of a Feather">
            </span>
          </div>
          <div class="min-w-0 flex-1 pt-1.5 flex justify-between space-x-4">
            <div>
              <p class="text-sm text-gray-500 truncate"><a href="#" class="font-medium text-gray-900">{{ event.payload.text }}</a></p>
            </div>
            <div class="text-right text-sm whitespace-nowrap text-gray-500">
              {{ $dayjs(event.createdAt).fromNow() }}
            </div>
          </div>
        </div>
      </div>
    </li>
  </ul>
</template>

<script>
export default {
  name: 'Bof',
  props: {
    event: {
      type: Object,
      required: true,
    },
  },
}
</script>
