<template>
  <div id="taskman">
    <div class="header">
      <Header></Header>
    </div>
    <div class="taskman-content-container">
      <Breadcrumb :style="setBreadcrumbStyle" v-if="isShowBreadcrum">
        <BreadcrumbItem>
          <a @click="homePageClickHandler">{{ $t('tw_home') }}</a>
        </BreadcrumbItem>
        <BreadcrumbItem v-for="(item, index) in breadcrumbList" :key="index">
          {{ item }}
        </BreadcrumbItem>
      </Breadcrumb>
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
      expandSideMenu: false,
      allMenusAry: [],
      childBreadcrumb: '',
      breadcrumbList: [],
      isSetting: this.$route.path.startsWith('/setting')
    }
  },
  computed: {
    setBreadcrumbStyle () {
      // 给侧边菜单栏适配样式
      return {
        margin: this.expandSideMenu ? '10px 0 10px 140px' : '10px 0'
      }
    }
  },
  watch: {
    $route: {
      handler (val) {
        this.breadcrumbList = []
        this.setBreadcrumb()
      },
      immediate: true
    }
  },
  created () {
    this.$bus.$on('expand-menu', val => {
      this.expandSideMenu = val
    })
  },
  methods: {
    setBreadcrumb () {
      this.isShowBreadcrum = !(this.$route.path === '/homepage' || this.$route.path === '/404')
      let lang = localStorage.getItem('lang') || navigator.language
      let langKey = lang === 'zh-CN' ? 'zh' : 'en'
      const routeList = this.$route.matched || []
      // const curBreadcrumb = this.$route.meta && this.$route.meta[langKey] || ''
      // this.breadcrumbList.push(curBreadcrumb)
      routeList.forEach(item => {
        if (item.meta && item.meta[langKey]) {
          this.breadcrumbList.push(item.meta[langKey])
        }
      })
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
