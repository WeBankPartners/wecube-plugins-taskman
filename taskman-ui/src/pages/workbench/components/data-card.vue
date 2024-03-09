<template>
  <div class="taskman-workbench-data-card">
    <Card v-for="(i, index) in data" :key="index" border :style="activeStyles(i)" class="card">
      <div class="content" @click="handleTabChange(i)">
        <div class="content-left">
          <Icon :type="i.icon" :color="i.color" size="28"></Icon>
          <span style="margin-left:10px;">{{ `${i.label}` }}</span>
        </div>
        <div class="content-right">
          <span class="number">{{ i.publishNum || '' }}</span>
          <span v-if="i.type === 'pending' && getPendingNum('1') > 0" class="badge">{{ getPendingNum('1') }}</span>
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
      active: '',
      action: '',
      data: [
        {
          type: 'pending',
          label: this.$t('tw_pending'),
          icon: 'ios-alert',
          color: '#ed4014',
          requestNum: 0,
          publishNum: 0
        },
        {
          type: 'hasProcessed',
          label: this.$t('tw_hasProcessed'),
          icon: 'ios-checkmark-circle',
          color: '#1990ff',
          requestNum: 0,
          publishNum: 0
        },
        {
          type: 'submit',
          label: this.$t('tw_submit'),
          icon: 'ios-send',
          color: '#19be6b',
          requestNum: 0,
          publishNum: 0
        },
        {
          type: 'draft',
          label: this.$t('tw_draft'),
          icon: 'ios-archive',
          color: '#b886f8',
          requestNum: 0,
          publishNum: 0
        },
        {
          type: 'collect',
          label: this.$t('tw_collect'),
          icon: 'ios-star',
          color: '#ff9900',
          requestNum: 0,
          publishNum: 0
        }
      ],
      pendingNumObj: {
        '1': [], // 发布待处理数量统计
        '2': [] // 请求待处理数量统计
      }
    }
  },
  computed: {
    activeStyles () {
      return function (i) {
        return {
          marginRight: i.type === 'collect' ? '0px' : '20px',
          borderTop: i.type === this.active ? '4px solid #e59e2d' : ''
        }
      }
    },
    getPendingNum () {
      return function (type) {
        return this.pendingNumObj[type].reduce((sum, cur) => {
          return Number(sum) + Number(cur || 0)
        }, 0)
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
  methods: {
    async getData () {
      const { statusCode, data } = await overviewData()
      if (statusCode === 'OK') {
        for (let key in data) {
          this.data.forEach(i => {
            if (i.type === key) {
              const numArr = (data[key] && data[key].split(';')) || []
              i.publishNum = numArr[0]
              i.requestNum = numArr[1]
              i.total = Number(i.publishNum) + Number(i.requestNum)
            }
          })
        }
        const pendingTaskArr = (data.pendingTask && data.pendingTask.split(';')) || []
        const pendingApprove = (data.pendingApprove && data.pendingApprove.split(';')) || []
        const requestPending = (data.requestPending && data.requestPending.split(';')) || []
        const requestConfirm = (data.requestConfirm && data.requestConfirm.split(';')) || []
        this.pendingNumObj['1'] = [pendingTaskArr[0], pendingApprove[0], requestPending[0], requestConfirm[0]]
        this.pendingNumObj['2'] = [pendingTaskArr[1], pendingApprove[1], requestPending[1], requestConfirm[1]]
      }
    },
    handleTabChange: debounce(function (item) {
      this.active = item.type
      this.$emit('fetchData', item.type)
    }, 300)
  }
}
</script>
<style lang="scss">
.taskman-workbench-data-card {
  .ivu-card-body {
    width: 100%;
    padding: 10px;
  }
  display: flex;
  flex-direction: row;
  margin-top: 20px;
  .card {
    width: 100%;
    height: 80px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;
    .content {
      width: 100%;
      display: flex;
      justify-content: space-between;
      &-left {
        display: flex;
        align-items: center;
        width: calc(100% - 50px);
        color: rgba(16, 16, 16, 1);
        font-size: 15px;
        font-family: PingFangSC-regular;
        font-weight: bold;
      }
      &-right {
        position: relative;
        width: 50px;
        height: 60px;
        margin-right: 30px;
        display: flex;
        flex-direction: column;
        align-items: flex-end;
        justify-content: space-around;
        .name {
          font-size: 14px;
          color: rgba(16, 16, 16, 1);
        }
        .number {
          font-size: 20px;
          color: #e59e2d;
        }
        .badge {
          position: absolute;
          top: 10px;
          right: -20px;
          font-size: 10px;
          background-color: #f56c6c;
          border-radius: 10px;
          color: #fff;
          height: 18px;
          line-height: 18px;
          padding: 0 6px;
          text-align: center;
          white-space: nowrap;
        }
      }
    }
  }
}
</style>
