<template>
  <div id="app">
    <div class="app-content-container">
      <BackTop :height="100" :bottom="100">
        <div class="w-back-top">
          <img :src="fesUpIcon" style="width: 24px;"/>

        </div>
      </BackTop>
      <router-view />
    </div>
  </div>
</template>
<script>
import fesUpIcon from '@/images/icon/fes_up.svg'
export default {
  data () {
    return {
      fesUpIcon
    }
  },
  methods: {
    jump (path) {
      this.$router.push({ path: path })
    }
  }
}
</script>
<style lang="scss">
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #0F1222; // Fes字体颜色规范
  height: 100%;
  min-width: 1280px;

  .app-content-container {
    height: 100%;

    .w-back-top {
      width: 44px;
      height: 44px;
      background-color: #fff;
      border-radius: 22px;
      box-shadow: 0 2px 8px 0 rgba(15, 18, 34, 0.1);
      display: flex;
      align-items: center;
      justify-content: center;
      &:hover {
        background-color: #f5f8ff;
      }
    }
  }
}
#nav {
  padding: 30px;
  a {
    font-weight: bold;
    color: #0F1222; // Fes字体颜色规范
    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
.ivu-layout,
.ivu-layout-sider {
  height: 100%;
}
.spin-icon-load {
  animation: ani-demo-spin 1s linear infinite;
}
// .ivu-form-item {
//   margin-bottom: 8px;
// }
body {
  height: 100%;
  overflow: auto !important;
}
html {
  height: 100%;
}
// .ivu-table-fixed-body {
//   height: auto !important;
// }
</style>
