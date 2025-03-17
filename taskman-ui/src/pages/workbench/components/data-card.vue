<template>
  <div ref="dataCardWrap" class="taskman-workbench-data-card">
    <Card v-for="(i, index) in cardList" :key="index" border :style="activeStyles(i)">
      <div class="content" @click="handleTabChange(i)">
        <div class="w-header">
          <Icon :type="i.icon" :color="i.color" :size="i.size"></Icon>
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
            <span class="number">{{ i[j.type] || '0' }}</span>
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
      cardList: [
        {
          type: 'myPending',
          label: this.$t('tw_my_pending'),
          icon: 'ios-person',
          size: '28',
          color: '#FF4D4F',
          total: 0,
          '1': 0,
          '2': 0,
          '3': 0,
          '4': 0,
          '5': 0
        },
        {
          type: 'pending',
          label: this.$t('tw_group_pending'),
          icon: 'ios-people',
          size: '28',
          color: '#F29360',
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
          color: '#00CB91',
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
          label: this.$t('tw_publish') // 发布
        },
        {
          type: '2',
          label: this.$t('tw_request') // 请求
        },
        {
          type: '3',
          label: this.$t('tw_question') // 问题
        },
        {
          type: '4',
          label: this.$t('tw_event') // 事件
        },
        {
          type: '5',
          label: this.$t('fork') // 变更
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
          marginRight: i.type === 'hasProcessed' ? '0px' : '10px',
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
    // 设置定时器，每分钟刷新本人处理/本组处理数量，刷新列表数据
    this.interval = setInterval(() => {
      this.getData(false, true)
      if (['myPending', 'pending'].includes(this.active)) {
        this.$parent && this.$parent.getList()
      }
    }, 60 * 1000)
  },
  destroyed () {
    if (this.interval) {
      clearInterval(this.interval)
    }
  },
  methods: {
    async getData (init = false, interval = false) {
      // ini为true，初始化拉取所有数据，后续拉取特定场景下的数据
      const params = {
        tab: init ? 'all' : this.active,
        queryTimeStart: this.$parent.queryTime[0] && this.$parent.queryTime[0] + ' 00:00:00',
        queryTimeEnd: this.$parent.queryTime[1] && this.$parent.queryTime[1] + ' 23:59:59'
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
                for (let key in data[item.type]) {
                  item[sceneMapValue[key]] = data[item.type][key]
                }
              }
            })
          }
        } else {
          // 面板数据
          this.cardList.forEach(item => {
            // 刷新当前tab页签数量
            if (item.type === this.active) {
              for (let key in data[item.type]) {
                item[sceneMapValue[key]] = data[item.type][key]
              }
            }
            // 定时器刷新本人处理/本组处理数量
            if (interval && ['myPending', 'pending'].includes(item.type)) {
              for (let key in data[item.type]) {
                item[sceneMapValue[key]] = data[item.type][key]
              }
            }
          })
        }
        // tab页签数据
        if (['myPending', 'pending'].includes(this.active)) {
          for (let key in this.pendingNumObj) {
            for (let type in this.pendingNumObj[key]) {
              this.pendingNumObj[key][type] = data[this.active][`${sceneMapWord[key]}${type}`]
            }
          }
        }
      }
    },
    handleTabChange: debounce(function (item, subType) {
      this.active = item.type
      this.action = subType || this.action
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
        padding: 0 5px;
        span {
          font-size: 12px;
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
