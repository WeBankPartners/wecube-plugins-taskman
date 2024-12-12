<template>
  <div class="taskman-wrap">
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  mounted () {
    // 获取body下所有的div元素
    setTimeout(() => {
      var body = document.querySelector('body')
      body.classList.add('taskman-wrap')
    }, 300)
  }
}
</script>

<style>

</style>
