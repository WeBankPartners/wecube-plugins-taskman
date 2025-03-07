<template>
  <div id="workbench">
    <div :style="workbenchStyle">
      <transition name="fade" mode="out-in">
        <router-view></router-view>
      </transition>
      <BaseMenu :menuList="menuList">
        <template slot="header">
          <img @click="handleGoHome" style="width:23px;height:23px;margin-right:10px;" src="@/images/menu_desk.png" />
          <span @click="handleGoHome">{{ $t('tw_workbench') }}</span>
        </template>
      </BaseMenu>
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
      expand: true,
      menuList: [
        {
          title: this.$t('tw_publish'),
          icon: 'md-person-add',
          name: '1',
          children: [
            { title: this.$t('tw_new'), path: '/taskman/workbench/template?type=1', name: '1-1' },
            { title: this.$t('tw_history'), path: '/taskman/workbench/publishHistory', name: '1-2' }
          ]
        },
        {
          title: this.$t('tw_request'),
          icon: 'ios-send',
          name: '2',
          children: [
            { title: this.$t('tw_new'), path: '/taskman/workbench/template?type=2', name: '2-1' },
            { title: this.$t('tw_history'), path: '/taskman/workbench/requestHistory', name: '2-2' }
          ]
        },
        {
          title: this.$t('tw_question'),
          icon: 'md-help-circle',
          name: '3',
          children: [
            { title: this.$t('tw_new'), path: '/taskman/workbench/template?type=3', name: '3-1' },
            { title: this.$t('tw_history'), path: '/taskman/workbench/problemHistory', name: '3-2' }
          ]
        },
        {
          title: this.$t('tw_event'),
          icon: 'md-pulse',
          name: '4',
          children: [
            { title: this.$t('tw_new'), path: '/taskman/workbench/template?type=4', name: '4-1' },
            { title: this.$t('tw_history'), path: '/taskman/workbench/eventHistory', name: '4-2' }
          ]
        },
        {
          title: this.$t('fork'),
          icon: 'md-git-merge',
          name: '5',
          children: [
            { title: this.$t('tw_new'), path: '/taskman/workbench/template?type=5', name: '5-1' },
            { title: this.$t('tw_history'), path: '/taskman/workbench/changeHistory', name: '5-2' }
          ]
        }
      ]
    }
  },
  computed: {
    workbenchStyle () {
      return {
        paddingLeft: this.expand ? '140px' : '0px'
      }
    }
  },
  created () {
    this.$bus.$on('expand-menu', val => {
      this.expand = val
    })
  },
  methods: {
    handleGoHome () {
      this.$router.push({
        path: '/taskman/workbench'
      })
    }
  }
}
</script>
<style lang="scss">
#workbench {
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
