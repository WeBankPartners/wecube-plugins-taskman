<template>
  <div id="taskman">
    <div class="header">
      <Header></Header>
    </div>
    <div class="taskman-content-container">
      <transition name="fade" mode="out-in">
        <router-view></router-view>
      </transition>
    </div>
  </div>
</template>

<script>
import Header from '@/pages/components/header'
export default {
  components: {
    Header
  },
  data () {
    return {
      isShowBreadcrum: true,
      allMenusAry: [],
      childBreadcrumb: '',
      isSetting: this.$route.path.startsWith('/setting')
    }
  },
  watch: {
    $route: {
      handler (val) {
        this.setBreadcrumb()
      },
      immediate: true
    }
  },
  methods: {
    setBreadcrumb () {
      this.childBreadcrumb = '453264'
    },
    homePageClickHandler () {
      window.needReLoad = false
      this.$router.push('/taskman/workbench')
    }
  }
}
</script>

<style lang="scss" scoped>
#taskman {
  height: 100%;
}
.header {
  width: 100%;
  background-color: #515a6e;
  display: block;
}
.taskman-content-container {
  height: calc(100% - 50px);
  padding: 10px 20px;
}
</style>
<style lang="scss">
#taskman {
  .ivu-breadcrumb {
    color: #515a6e;
  }
  .ivu-tooltip {
    width: auto !important;
  }
  label {
    margin-bottom: 0px !important; /*解决监控插件label样式全局覆盖问题*/
  }
  .ivu-table-tip span {
    position: absolute;
    left: 50%;
    transform: translate(-50%, -50%);
    pointer-events: none;
  }
  .ivu-table-wrapper {
    z-index: 100;
  }
  .ivu-tag {
    display: inline-block;
    line-height: 16px;
    height: auto;
    padding: 5px 6px;
  }
}
</style>
