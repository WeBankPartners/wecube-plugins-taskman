<template>
  <div class="taskman-workbench-data-card">
    <Card v-for="(i, index) in cardList" :key="index" border :style="activeStyles(i)">
      <div class="content" @click="handleTabChange(i)">
        <div class="w-header">
          <Icon :type="i.icon" :color="i.color" :size="i.size"></Icon>
          <!-- <div v-if="i.type === 'pending'" class="group-btn">
            <span
              @click="handlePendGroupChange('my')"
              :style="{ color: pendingType === 'my' ? 'rgb(229, 158, 45)' : '#17233d' }"
            >
              本人处理
            </span>
            <span> / </span>
            <span
              @click="handlePendGroupChange('group')"
              :style="{ color: pendingType === 'group' ? 'rgb(229, 158, 45)' : '#17233d' }"
            >
              本组处理
            </span>
          </div> -->
          <span style="margin-left:5px;">
            {{ `${i.label}` }}
            <!-- <span class="total">{{`(${i.total || 0})`}}</span> -->
          </span>
        </div>
        <div class="data">
          <div
            v-for="j in actionList"
            :key="j.type"
            class="list"
            :style="actionStyles(i, j.type)"
            @click="
              e => {
                e.stopPropagation()
                handleTabChange(i, j.type)
              }
            "
          >
            <span class="number">
              {{ i[j.type] || '0' }}
              <!-- <span v-if="i.type === 'pending' && getPendingNum(j.type) > 0" class="badge">{{
                getPendingNum(j.type)
              }}</span> -->
            </span>
            <span>{{ j.label }}</span>
          </div>
        </div>
      </div>
    </Card>
  </div>
</template>

<script>
import { overviewData } from '@/api/server'
import { debounce } from '@/pages/util'
export default {
  props: {
    initTab: {
      type: String,
      default: ''
    },
    initAction: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      active: '', // pending(myPending本人处理/pending本组处理),hasProcessed已处理,submit我提交的,draft我的暂存
      action: '',
      // pendingType: 'my',
      cardList: [
        {
          type: 'myPending',
          label: '本人处理',
          icon: 'md-person',
          size: '28',
          color: '#ed4014',
          total: 0,
          '1': 0,
          '2': 0,
          '3': 0,
          '4': 0,
          '5': 0
        },
        {
          type: 'pending',
          label: '本组处理',
          icon: 'md-people',
          size: '28',
          color: '#ed4014',
          total: 0,
          '1': 0,
          '2': 0,
          '3': 0,
          '4': 0,
          '5': 0
        },
        {
          type: 'draft',
          label: this.$t('tw_draft'),
          icon: 'ios-archive',
          size: '28',
          color: '#b886f8',
          total: 0,
          '1': 0,
          '2': 0,
          '3': 0,
          '4': 0,
          '5': 0
        },
        {
          type: 'submit',
          label: this.$t('tw_submit'),
          icon: 'ios-send',
          size: '28',
          color: '#19be6b',
          total: 0,
          '1': 0,
          '2': 0,
          '3': 0,
          '4': 0,
          '5': 0
        },
        {
          type: 'hasProcessed',
          label: this.$t('tw_hasProcessed'),
          icon: 'ios-checkmark-circle',
          size: '28',
          color: '#1990ff',
          total: 0,
          '1': 0,
          '2': 0,
          '3': 0,
          '4': 0,
          '5': 0
        }
      ],
      actionList: [
        {
          type: '1',
          label: this.$t('tw_publish')
        },
        {
          type: '2',
          label: this.$t('tw_request')
        },
        {
          type: '3',
          label: '问题'
        },
        {
          type: '4',
          label: '事件'
        },
        {
          type: '5',
          label: '变更'
        }
      ],
      pendingNumObj: {
        '1': { Approve: 0, Task: 0, Check: 0, Confirm: 0 }, // 发布
        '2': { Approve: 0, Task: 0, Check: 0, Confirm: 0 }, // 请求
        '3': { Approve: 0, Task: 0, Check: 0, Confirm: 0 }, // 问题
        '4': { Approve: 0, Task: 0, Check: 0, Confirm: 0 }, // 事件
        '5': { Approve: 0, Task: 0, Check: 0, Confirm: 0 } // 变更
      },
      interval: null // 工作台每秒轮询一次
    }
  },
  computed: {
    activeStyles () {
      return function (i) {
        return {
          width: '100%',
          height: '105px',
          marginRight: i.type === 'hasProcessed' ? '0px' : '15px',
          cursor: 'pointer',
          borderTop: i.type === this.active ? '4px solid #e59e2d' : ''
        }
      }
    },
    actionStyles () {
      return function (i, val) {
        return {
          color: this.action === val && i.type === this.active ? '#e59e2d' : '#17233d'
        }
      }
    }
    // getPendingNum () {
    //   return function (type) {
    //     return this.pendingNumObj[type].reduce((sum, cur) => {
    //       return Number(sum) + Number(cur || 0)
    //     }, 0)
    //   }
    // }
  },
  watch: {
    initAction: {
      handler (val) {
        if (val) {
          this.action = val
          if (this.action && this.active) {
            this.$emit('initFetch', this.active, this.action)
          }
        }
      },
      immediate: true
    },
    initTab: {
      handler (val) {
        if (val) {
          this.active = val
          if (this.action && this.active) {
            this.$emit('initFetch', this.active, this.action)
          }
        }
      },
      immediate: true
    }
  },
  mounted () {
    this.interval = setInterval(() => {
      this.getData(false, true)
    }, 60 * 1000)
  },
  destroyed () {
    if (this.interval) {
      clearInterval(this.interval)
    }
  },
  methods: {
    // handlePendGroupChange (val) {
    //   this.pendingType = val
    // },
    async getData (init = false, interval = false) {
      // ini为true，初始化拉取所有数据，后续拉取特定场景下的数据
      const params = {
        tab: init ? 'all' : this.active
      }
      // 设置每分钟轮询查询本人处理数据
      if (interval) {
        params.tab = 'myPending'
      }
      const sceneMapValue = {
        Release: '1',
        Request: '2',
        Problem: '3',
        Event: '4',
        Change: '5'
      }
      const sceneMapWord = {
        '1': 'Release',
        '2': 'Request',
        '3': 'Problem',
        '4': 'Event',
        '5': 'Change'
      }
      const { statusCode, data } = await overviewData(params)
      if (statusCode === 'OK') {
        if (init) {
          // 初始化所有数据
          for (let tabName in data) {
            this.cardList.forEach(item => {
              if (item.type === tabName) {
                // 面板数据
                if (item.type === 'myPending') {
                  for (let key in data['myPending']) {
                    item[sceneMapValue[key]] = data['myPending'][key]
                  }
                } else {
                  for (let key in data[tabName]) {
                    item[sceneMapValue[key]] = data[tabName][key]
                  }
                }
                // tab页签数据(本人处理/本组处理)
                if (['myPending', 'pending'].includes(item.type)) {
                  for (let key in this.pendingNumObj) {
                    for (let type in this.pendingNumObj[key]) {
                      this.pendingNumObj[key][type] = data[item.type][`${sceneMapWord[key]}${type}`]
                    }
                  }
                }
              }
            })
          }
        } else {
          // 面板数据
          this.cardList.forEach(item => {
            if (item.type === this.active) {
              if (item.type === 'myPending') {
                for (let key in data['myPending']) {
                  item[sceneMapValue[key]] = data['myPending'][key]
                }
              } else {
                for (let key in data[this.active]) {
                  item[sceneMapValue[key]] = data[this.active][key]
                }
              }
            }
          })
          // tab页签数据
          if (['myPending', 'pending'].includes(this.active)) {
            for (let key in this.pendingNumObj) {
              for (let type in this.pendingNumObj[key]) {
                this.pendingNumObj[key][type] = data[this.active][`${sceneMapWord[key]}${type}`]
              }
            }
          }
        }
      }
    },
    handleTabChange: debounce(function (item, subType) {
      this.active = item.type
      this.action = subType || '1'
      this.$emit('fetchData', this.active, this.action)
    }, 300)
  }
}
</script>

<style lang="scss">
.taskman-workbench-data-card .ivu-card-body {
  padding: 10px;
}
</style>
<style lang="scss" scoped>
.taskman-workbench-data-card {
  display: flex;
  .content {
    .w-header {
      display: flex;
      align-items: center;
      .group-btn {
        margin-left: 5px;
      }
      span {
        color: #17233d;
        font-size: 14px;
        font-family: PingFangSC-regular;
        font-weight: bold;
      }
      .total {
        margin-left: 0px;
        font-size: 14px;
        color: #17233d;
      }
    }
    .data {
      display: flex;
      justify-content: space-around;
      align-items: center;
      margin-top: 10px;
      .list {
        position: relative;
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 0 10px;
        span {
          font-size: 13px;
          font-family: PingFangSC-regular;
        }
        .number {
          position: relative;
          font-weight: bold;
          font-size: 14px;
          .badge {
            position: absolute;
            top: -5px;
            right: -25px;
            font-size: 10px;
            background-color: #f56c6c;
            border-radius: 10px;
            color: #fff;
          }
        }
      }
    }
  }
}
</style>
