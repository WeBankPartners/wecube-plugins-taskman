<template>
  <div class="workbench-data-card">
    <Card v-for="(i, index) in data" :key="index" border :style="activeStyles(i)">
      <div class="content" @click="handleTabChange(i)">
        <div class="w-header">
          <!-- <img :src="i.icon" /> -->
          <Icon :type="i.icon" :color="i.color" size="26"></Icon>
          <span style="margin-left:5px;">{{ `${i.label}(${i.total || 0})` }}</span>
        </div>
        <div class="data">
          <div
            class="list"
            :style="actionStyles(i, '1')"
            @click="
              e => {
                e.stopPropagation()
                handleTabChange(i, '1')
              }
            "
          >
            <span style="font-weight:bold;">{{ i.publishNum || '' }}</span>
            <Badge v-if="i.type === 'pending'" class-name="badge" :count="getPendingNum('1')"></Badge>
            <span>{{ $t('tw_publish') }}</span>
          </div>
          <div
            class="list"
            :style="actionStyles(i, '2')"
            @click="
              e => {
                e.stopPropagation()
                handleTabChange(i, '2')
              }
            "
          >
            <span style="font-weight:bold;">{{ i.requestNum || '' }}</span>
            <Badge v-if="i.type === 'pending'" class-name="badge" :count="getPendingNum('2')"></Badge>
            <span>{{ $t('tw_request') }}</span>
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
          color: '#19be6b',
          requestNum: 0,
          publishNum: 0
        },
        {
          type: 'submit',
          label: this.$t('tw_submit'),
          icon: 'ios-cloud-upload',
          color: '#1990ff',
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
          width: '100%',
          height: '110px',
          marginRight: i.type === 'collect' ? '0px' : '20px',
          cursor: 'pointer',
          borderTop: i.type === this.active ? '4px solid #e59e2d' : ''
        }
      }
    },
    actionStyles () {
      return function (i, val) {
        return {
          color: this.action === val && i.type === this.active ? '#e59e2d' : 'rgba(16, 16, 16, 1)'
        }
      }
    },
    getPendingNum () {
      return function (type) {
        return this.pendingNumObj[type].reduce((sum, cur) => {
          return Number(sum) + Number(cur)
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
    handleTabChange: debounce(function (item, subType) {
      this.active = item.type
      this.action = subType || '1'
      this.$emit('fetchData', item.type, subType)
    }, 300)
  }
}
</script>

<style lang="scss">
.workbench-data-card .ivu-card-body {
  padding: 10px;
}
.workbench-data-card .badge {
  position: absolute;
  top: -28px;
  right: -32px;
  font-size: 10px;
}
</style>
<style lang="scss" scoped>
.workbench-data-card {
  display: flex;
  .content {
    .w-header {
      display: flex;
      align-items: center;
      img {
        width: 24px;
        height: 24px;
        margin-right: 5px;
      }
      span {
        color: rgba(16, 16, 16, 1);
        font-size: 14px;
        font-family: PingFangSC-regular;
        font-weight: bold;
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
          font-size: 14px;
          font-family: PingFangSC-regular;
        }
      }
    }
  }
}
</style>
