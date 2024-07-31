<template>
  <div>
    <Header>
      <div class="menus">
        <Menu mode="horizontal" theme="dark" :active-name="activeName" @on-select="changeMenu">
          <div>
            <img @click="goHome" src="../../images/taskman.png" alt="LOGO" class="img-logo" />
          </div>
          <div>
            <MenuItem v-for="menu in menus" :name="menu.name" :key="menu.path">
              {{ menu.display }}
            </MenuItem>
          </div>
        </Menu>
      </div>
      <div class="header-right_container">
        <div class="profile">
          <Dropdown style="cursor: pointer" trigger="click">
            <img class="p-icon" src="../../images/icon/icon_usr.png" width="12" height="12" />
            <span>{{ username }}</span>
            <Icon type="ios-arrow-down"></Icon>
            <Badge :count="pendingCount" @click="userMgmt"></Badge>
            <DropdownMenu slot="list">
              <!-- <DropdownItem name="logout" to="/login">
                <a @click="showChangePassword" style="width: 100%; display: block">
                  {{ $t('change_password') }}
                </a>
              </DropdownItem> -->
              <DropdownItem name="userApply">
                <a @click="roleApply" style="width: 100%; display: block">
                  {{ $t('tw_apply_roles') }}
                </a>
              </DropdownItem>
              <DropdownItem name="userMgmt">
                <a @click="userMgmt" style="width: 100%; display: block">
                  {{ $t('tw_user_mgmt') }}
                  <Badge :count="pendingCount" @click="userMgmt" style="top: -2px"></Badge>
                </a>
              </DropdownItem>
              <DropdownItem name="logout" to="/login">
                <a @click="logout" style="width: 100%; display: block">
                  {{ $t('logout') }}
                </a>
              </DropdownItem>
            </DropdownMenu>
          </Dropdown>
        </div>
        <div class="language">
          <Dropdown>
            <a href="javascript:void(0)">
              <img
                class="p-icon"
                v-if="currentLanguage === 'English'"
                src="../../images/icon/icon_lan_EN.png"
                width="12"
                height="12"
              />
              <img class="p-icon" v-else src="../../images/icon/icon_lan_CN.png" width="12" height="12" />
              {{ currentLanguage }}
              <Icon type="ios-arrow-down"></Icon>
            </a>
            <DropdownMenu slot="list">
              <DropdownItem v-for="(item, key) in language" :key="item.id" @click.native="changeLanguage(key)">
                {{ item }}
              </DropdownItem>
            </DropdownMenu>
          </Dropdown>
        </div>
        <div class="language" @click="changeDocs">
          <img class="p-icon" src="../../images/icon/icon_help.png" width="12" height="12" />
          {{ $t('help_docs') }}
        </div>
        <!-- <div class="version">{{ version }}</div> -->
      </div>
    </Header>
    <UserMgmt ref="userMgmtRef"></UserMgmt>
    <RoleApply ref="roleApplyRef"></RoleApply>
  </div>
</template>
<script>
import { getProcessableList } from '@/api/server.js'
import UserMgmt from './user-mgmt.vue'
import RoleApply from './role-apply.vue'
import { clearAllCookie } from '@/pages/util/cookie'
import Vue from 'vue'
export default {
  data () {
    return {
      activeName: '/taskman/workbench',
      username: '',
      currentLanguage: '',
      language: {
        'zh-CN': '简体中文',
        'en-US': 'English'
      },
      menus: [
        {
          display: this.$t('tw_workbench'),
          name: '/taskman/workbench',
          path: '/taskman/workbench'
        },
        {
          display: this.$t('tw_template_group_mgmt'),
          name: '/taskman/template-group',
          path: '/taskman/template-group'
        },
        {
          display: this.$t('tw_template_mgmt'),
          name: '/taskman/template-mgmt',
          path: '/taskman/template-mgmt'
        }
      ],
      needLoad: true,
      version: '',
      changePassword: false,
      pendingCount: 0, // 待审批数量
      timer: null
    }
  },
  watch: {
    $lang: async function (lang) {
      window.location.reload()
    },
    $route: {
      handler (val) {
        const activeMenu = this.menus.find(menu => val.name.startsWith(menu.path))
        if (activeMenu) {
          this.activeName = activeMenu.name
        }
      },
      immediate: true
    }
  },
  async created () {
    this.getLocalLang()
    this.username = window.localStorage.getItem('username')
    this.$bus.$on('fetchApplyCount', () => {
      this.getPendingCount()
    })
  },
  mounted () {
    this.getPendingCount()
    this.timer = setInterval(() => {
      this.getPendingCount()
    }, 5 * 60 * 1000)
  },
  destroyed () {
    if (this.timer) {
      clearInterval(this.timer)
    }
  },
  methods: {
    changeMenu (path) {
      this.$router.push({ path: path })
    },
    goHome () {
      this.$router.push({ path: '/taskman/workbench' })
    },
    userMgmt () {
      this.$refs.userMgmtRef.openModal()
    },
    roleApply () {
      this.$refs.roleApplyRef.openModal()
    },
    logout () {
      clearInterval(this.timer)
      localStorage.clear()
      clearAllCookie()
      window.location.href = window.location.origin + window.location.pathname + '#/login'
    },
    changeLanguage (lan) {
      Vue.config.lang = lan
      this.currentLanguage = this.language[lan]
      localStorage.setItem('lang', lan)
      window.location.reload()
    },
    getLocalLang () {
      let currentLangKey = localStorage.getItem('lang') || navigator.language
      const lang = this.language[currentLangKey] || 'English'
      this.currentLanguage = lang
    },
    changeDocs () {
      window.open('http://webankpartners.gitee.io/wecube-docs')
    },
    async getPendingCount () {
      const params = {
        filters: [
          {
            name: 'status',
            operator: 'in',
            value: ['init']
          }
        ],
        paging: true,
        pageable: {
          startIndex: 0,
          pageSize: 10000
        },
        sorting: [
          {
            asc: false,
            field: 'createdTime'
          }
        ]
      }
      const { status, data } = await getProcessableList(params)
      if (status === 'OK') {
        this.pendingCount = data.pageInfo.totalRows
      }
    }
  },
  components: {
    UserMgmt,
    RoleApply
  }
}
</script>

<style lang="scss" scoped>
.img-logo {
  height: 20px;
  margin: 0 4px 6px 0;
  vertical-align: middle;
  cursor: pointer;
}
.ivu-layout-header {
  padding: 0 20px;
}
.header {
  display: flex;
  .ivu-layout-header {
    height: 50px;
    line-height: 50px;
    background: linear-gradient(90deg, #8bb8fa 0%, #e1ecfb 100%);
  }
  a {
    color: #404144;
  }
  .menus {
    display: inline-block;
    .ivu-menu-horizontal {
      height: 50px;
      line-height: 50px;
      display: flex;
      .ivu-menu-submenu {
        padding: 0 8px;
        font-size: 15px;
        color: #404144;
      }
      .ivu-menu-item {
        font-size: 15px;
        color: #404144;
      }
    }
    .ivu-menu-dark {
      background: transparent;
    }
    .ivu-menu-dark.ivu-menu-horizontal .ivu-menu-submenu {
      color: #404144;
    }
    .ivu-menu-item-active,
    .ivu-menu-item:hover {
      color: #116ef9 !important;
    }
    .ivu-menu-dark.ivu-menu-horizontal .ivu-menu-submenu-active,
    .ivu-menu-dark.ivu-menu-horizontal .ivu-menu-submenu:hover {
      color: #116ef9;
    }
    // .ivu-menu-drop-list {
    //   .ivu-menu-item-active,
    //   .ivu-menu-item:hover {
    //     color: black;
    //   }
    // }
  }
  .header-right_container {
    position: absolute;
    right: 20px;
    top: 0;
    .language,
    .help,
    .version,
    .profile {
      float: right;
      display: inline-block;
      vertical-align: middle;
      margin-left: 20px;
      cursor: pointer;
    }
    .version {
      color: #404144;
    }

    .p-icon {
      margin-right: 6px;
    }

    .ivu-dropdown-rel {
      display: flex;
      align-items: center;
      a {
        display: flex;
        align-items: center;
      }
    }
  }
}
</style>
