<template>
  <div class="workbench-menu" :style="{ width: expand ? '180px' : '0px' }">
    <div v-show="expand" style="height:100%;">
      <div class="home">
        <img style="width:23px;height:23px;margin-right:10px;" src="@/images/menu_desk.png" />
        工作台
      </div>
      <Menu
        @on-select="handleSelectMenu"
        theme="dark"
        active-name="1-1"
        :open-names="['1']"
        style="width:180px;height:100%;"
      >
        <Submenu v-for="(i, index) in menuList" :key="index" :name="i.name">
          <template #title>
            <div class="menu-item">
              <img :src="i.icon" />
              {{ i.title }}
            </div>
          </template>
          <MenuItem v-for="(j, idx) in i.children" :key="idx" :name="j.name" :to="j.path">{{ j.title }}</MenuItem>
        </Submenu>
      </Menu>
    </div>
    <!-- <div v-else class="small-menu">
      <div class="small-menu-item">
        <img style="width:23px;height:23px;" src="@/images/menu_desk.png" />
        <span>工作台</span>
      </div>
      <div v-for="(i, index) in menuList" :key="index" class="small-menu-item">
        <img :src="i.icon" />
        <span>{{ i.title }}</span>
      </div>
    </div> -->
    <div class="expand" :style="{ left: expand ? '180px' : '0px' }">
      <Icon v-if="expand" @click="handleExpand" type="ios-arrow-dropleft" size="28" />
      <Icon v-else @click="handleExpand" type="ios-arrow-dropright" size="28" />
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
      expand: true,
      activeName: '1-1',
      menuList: [
        {
          title: '发布',
          icon: require('@/images/menu_publish.png'),
          name: '1',
          children: [
            { title: '新建发布', path: '/taskman/workbench/template', name: '1-1' },
            { title: '发布历史', path: '/taskman/workbench/publishHistory', name: '1-2' }
          ]
        },
        {
          title: '请求',
          icon: require('@/images/menu_request.png'),
          name: '2',
          children: [
            { title: '新建请求', path: '/taskman/workbench/template', name: '2-1' },
            { title: '请求历史', path: '/taskman/workbench/requestHistory', name: '2-2' }
          ]
        }
      ]
    }
  },
  methods: {
    handleExpand () {
      this.expand = !this.expand
      this.$bus.$emit('expand-menu', this.expand)
    },
    handleSelectMenu (name) {
      console.log('111111', name)
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
  }
}
</style>
<style lang="scss" scoped>
.workbench-menu {
  position: fixed;
  top: 0;
  left: 0;
  height: 100%;
  z-index: 9999;
  .home {
    display: flex;
    align-items: center;
    padding: 20px 20px 10px 20px;
    width: 180px;
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
