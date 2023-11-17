<template>
  <div class="hot-link">
    <span>常用入口：</span>
    <a v-for="(i, index) in linkArr" :key="index" @click="jumpTo(i.path)">{{ i.name }}</a>
  </div>
</template>

<script>
export default {
  data () {
    return {
      linkArr: [
        {
          name: '自定义视图',
          path: '/taskman/task-mgmt'
        },
        {
          name: '自定义视图',
          path: '/taskman/task-mgmt'
        },
        {
          name: '自定义视图',
          path: '/taskman/task-mgmt'
        },
        {
          name: '自定义视图',
          path: '/taskman/task-mgmt'
        }
      ]
    }
  },
  methods: {
    jumpTo (path) {
      this.$router.push(path)
    }
  }
}
</script>

<style lang="scss" scoped>
.hot-link {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  span {
    font-weight: 600;
  }
  a {
    margin-right: 20px;
    font-size: 12px;
  }
}
</style>
