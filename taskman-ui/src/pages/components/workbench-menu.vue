<template>
  <div
    class="workbench-menu"
    :style="{
      width: expand ? '140px' : '0px',
      top: scrollTop > 50 ? '0px' : 50 - scrollTop + 'px'
    }"
  >
    <div v-show="expand" style="height:100%;">
      <div class="home" @click="handleGoHome">
        <img style="width:23px;height:23px;margin-right:10px;" src="@/images/menu_desk.png" />
        {{ $t('tw_workbench') }}
      </div>
      <Menu
        @on-select="handleSelectMenu"
        theme="dark"
        :active-name="activeName"
        :open-names="openNames"
        style="width:140px;height:100%;"
      >
        <Submenu v-for="(i, index) in menuList" :key="index" :name="i.name">
          <template #title>
            <div class="menu-item">
              <img v-if="i.img" :src="i.img" />
              <Icon v-if="i.icon" :type="i.icon" :size="22" style="margin-right:10px;" color="#fff" />
              {{ i.title }}
            </div>
          </template>
          <MenuItem v-for="(j, idx) in i.children" :key="idx" :name="j.name" :to="j.path" :replace="false">{{
            j.title
          }}</MenuItem>
        </Submenu>
      </Menu>
    </div>
    <div class="expand" :style="{ left: expand ? '140px' : '0px' }">
      <Icon v-if="expand" @click="handleExpand" type="ios-arrow-dropleft" size="28" />
      <Icon v-else @click="handleExpand" type="ios-arrow-dropright" size="28" />
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
      scrollTop: 0,
      expand: true,
      activeName: '',
      openNames: [],
      menuList: [
        {
          title: this.$t('tw_publish'),
          icon: 'md-person-add',
          name: '1',
          children: [
            { title: this.$t('tw_new_publish'), path: '/taskman/workbench/template?type=1', name: '1-1' },
            { title: this.$t('tw_publish_history'), path: '/taskman/workbench/publishHistory', name: '1-2' }
          ]
        },
        {
          title: this.$t('tw_request'),
          icon: 'ios-send',
          name: '2',
          children: [
            { title: this.$t('tw_new_request'), path: '/taskman/workbench/template?type=2', name: '2-1' },
            { title: this.$t('tw_request_history'), path: '/taskman/workbench/requestHistory', name: '2-2' }
          ]
        },
        {
          title: '问题',
          icon: 'md-help-circle',
          name: '3',
          children: [
            { title: '新建问题', path: '/taskman/workbench/template?type=3', name: '3-1' },
            { title: '问题历史', path: '/taskman/workbench/problemHistory', name: '3-2' }
          ]
        },
        {
          title: '事件',
          icon: 'md-pulse',
          name: '4',
          children: [
            { title: '新建事件', path: '/taskman/workbench/template?type=4', name: '4-1' },
            { title: '事件历史', path: '/taskman/workbench/eventHistory', name: '4-2' }
          ]
        },
        {
          title: '变更',
          icon: 'md-git-merge',
          name: '5',
          children: [
            { title: '新建变更', path: '/taskman/workbench/template?type=5', name: '5-1' },
            { title: '变更历史', path: '/taskman/workbench/changeHistory', name: '5-2' }
          ]
        }
      ]
    }
  },
  created () {
    this.menuList.forEach(i => {
      for (let j of i.children) {
        if (j.path === this.$route.fullPath) {
          this.activeName = j.name
          this.openNames = [i.name]
        }
      }
    })
  },
  mounted () {
    if (this.$eventBusP) {
      this.$eventBusP.$emit('expand-menu', this.expand)
    } else {
      this.$bus.$emit('expand-menu', this.expand)
    }
    window.addEventListener('scroll', this.getScrollTop)
  },
  beforeDestroy () {
    if (this.$eventBusP) {
      this.$eventBusP.$emit('expand-menu', false)
    } else {
      this.$bus.$emit('expand-menu', false)
    }
    window.removeEventListener('scroll', this.getScrollTop)
  },
  methods: {
    getScrollTop () {
      this.scrollTop = document.documentElement.scrollTop || document.body.scrollTop
    },
    handleExpand () {
      this.expand = !this.expand
      if (this.$eventBusP) {
        this.$eventBusP.$emit('expand-menu', this.expand)
      } else {
        this.$bus.$emit('expand-menu', this.expand)
      }
    },
    handleSelectMenu (name) {
      this.activeName = name
    },
    handleGoHome () {
      this.$router.push({
        path: '/taskman/workbench'
      })
    }
  }
}
</script>

<style lang="scss">
.workbench-menu {
  .ivu-menu-dark {
    background: #001529;
  }
  .ivu-menu-dark.ivu-menu-vertical .ivu-menu-opened .ivu-menu-submenu-title {
    background: #10192b;
  }
  .ivu-menu-dark.ivu-menu-vertical .ivu-menu-item,
  .ivu-menu-dark.ivu-menu-vertical .ivu-menu-submenu-title {
    background: #10192b;
    padding: 10px;
  }
  .ivu-menu-item {
    padding-left: 32px !important;
  }
}
</style>
<style lang="scss" scoped>
.workbench-menu {
  position: fixed;
  left: 0;
  height: 100%;
  z-index: 100;
  .home {
    display: flex;
    align-items: center;
    padding: 20px 12px 10px 12px;
    width: 140px;
    background: #002140;
    color: #fff;
    font-size: 14px;
    cursor: pointer;
  }
  .menu-item {
    display: flex;
    align-items: center;
    img {
      width: 23px;
      height: 23px;
      margin-right: 10px;
    }
  }
  .small-menu {
    width: 70px;
    height: 100%;
    background: #10192b;
    &-item {
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      padding-top: 20px;
      cursor: pointer;
      img {
        width: 23px;
        height: 23px;
      }
      span {
        font-size: 14px;
        color: #fff;
        font-weight: bold;
      }
    }
  }
  .expand {
    position: absolute;
    top: calc(50% - 14px);
    cursor: pointer;
  }
}
</style>
