<template>
  <div id="app">
    <Button @click="jump('/taskman/template-mgmt')">模板设置</Button>
    <Button @click="jump('/taskman/template-group')">模板组设置</Button>
    <Button @click="jump('/taskman/request-mgmt')">发起请求</Button>
    <Button @click="jump('/taskman/task-mgmt')">任务</Button>
    <Button @click="jump('/taskman/workbench')">个人工作台</Button>
    <div class="app-content-container" :style="showWorkbench ? workbenchStyle : {}">
      <BackTop :height="100" :bottom="100" />
      <WorkbenchMenu v-if="showWorkbench"></WorkbenchMenu>
      <router-view :key="$route.path" />
    </div>
  </div>
</template>
<script>
import WorkbenchMenu from '@/pages/components/workbench-menu.vue'
import Vue from 'vue'
Vue.prototype.$bus = new Vue()
export default {
  components: {
    WorkbenchMenu
  },
  data () {
    return {
      expand: true
    }
  },
  computed: {
    workbenchStyle () {
      return {
        paddingLeft: this.expand ? '200px' : '20px'
      }
    },
    showWorkbench () {
      return this.$route.path.indexOf('workbench') > -1
    }
  },
  mounted () {
    this.$bus.$on('expand-menu', val => {
      this.expand = val
    })
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
  color: #2c3e50;
  height: 100%;
  min-width: 1280px;

  .app-content-container {
    height: 100%;
  }
}
#nav {
  padding: 30px;
  a {
    font-weight: bold;
    color: #2c3e50;
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
.ivu-form-item {
  margin-bottom: 8px;
}
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
