<template>
  <div id="taskman">
    <div class="taskman-content-container" :style="workbenchStyle">
      <transition name="fade" mode="out-in">
        <router-view></router-view>
      </transition>
      <WorkbenchMenu></WorkbenchMenu>
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
      isShowBreadcrum: true,
      allMenusAry: [],
      childBreadcrumb: '',
      isSetting: this.$route.path.startsWith('/setting'),
      expand: true
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
  computed: {
    setBreadcrumbStyle () {
      // 给侧边菜单栏适配样式
      return {
        margin: this.expandSideMenu ? '10px 0 10px 140px' : '10px 0'
      }
    },
    workbenchStyle () {
      return {
        paddingLeft: this.expand ? '140px' : '0px'
      }
    }
  },
  mounted () {
    if (this.$eventBusP) {
      this.$eventBusP.$on('expand-menu', val => {
        this.expand = val
      })
    } else {
      this.$bus.$on('expand-menu', val => {
        this.expand = val
      })
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
