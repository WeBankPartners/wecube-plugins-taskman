<template>
  <div id="taskman-wrap">
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  mounted () {
    document.body.classList.add('taskman-wrap')
  }
}
</script>

<style>

</style>
