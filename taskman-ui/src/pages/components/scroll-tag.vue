<template>
  <div class="taskman-scroll-column">
    <div v-for="(i, index) in list" :key="index" class="tag-111">
      {{ i }}
    </div>
  </div>
</template>

<script>
export default {
  props: {
    list: {
      type: Array,
      default: () => []
    }
  }
}
</script>

<style lang="scss" scoped>
.taskman-scroll-column {
  max-height: 66px;
  overflow-y: scroll;
  padding-right: 3px;
  .tag-111 {
    padding: 4px 5px;
    border-radius: 3px;
    background: #2d8cf0;
    color: #fff;
    margin-top: 5px;
  }
}
::-webkit-scrollbar {
  width: 5px;
  height: 5px;
}
::-webkit-scrollbar-thumb {
  background-color: #c1c1c1;
  border-radius: 16px;
}
::-webkit-scrollbar-track {
  background-color: transparent;
  border-radius: 16px;
}
</style>
