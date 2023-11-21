<template>
  <div id="taskman">
    <div class="content-container" :style="workbenchStyle">
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
      parentBreadcrumb: '',
      childBreadcrumb: '',
      isSetting: this.$route.path.startsWith('/setting'),
      expand: true
    }
  },
  computed: {
    workbenchStyle () {
      return {
        paddingLeft: this.expand ? '180px' : '0px'
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
.content-container {
  height: calc(100% - 50px);
}

.ivu-breadcrumb {
  color: #515a6e;
}
</style>
