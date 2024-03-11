<template>
  <div class="taskman-workbench-data-card">
    <Card v-for="(i, index) in cardList" :key="index" border :style="activeStyles(i)" class="card">
      <div class="content" @click="handleTabChange(i)">
        <div class="content-left">
          <Icon :type="i.icon" :color="i.color" size="28"></Icon>
          <span style="margin-left:10px;">{{ `${i.label}` }}</span>
        </div>
        <div class="content-right">
          <span class="number">{{ i[$parent.actionName] || 0 }}</span>
          <span v-if="i.type === 'pending' && getPendingNum($parent.actionName) > 0" class="badge">{{
            getPendingNum($parent.actionName)
          }}</span>
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
      cardList: [
        {
          type: 'pending',
          label: this.$t('tw_pending'),
          icon: 'ios-alert',
          color: '#ed4014',
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
          color: '#1990ff',
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
          color: '#19be6b',
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
          color: '#b886f8',
          '1': 0,
          '2': 0,
          '3': 0,
          '4': 0,
          '5': 0
        },
        {
          type: 'collect',
          label: this.$t('tw_collect'),
          icon: 'ios-star',
          color: '#ff9900',
          '1': 0,
          '2': 0,
          '3': 0,
          '4': 0,
          '5': 0
        }
      ],
      pendingNumObj: {
        '1': [], // 发布
        '2': [], // 请求
        '3': [], // 问题
        '4': [], // 事件
        '5': [] // 变更
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
      const params = {
        params: {
          scene: Number(this.$parent.actionName)
        }
      }
      const { statusCode, data } = await overviewData(params)
      if (statusCode === 'OK') {
        for (let key in data) {
          this.cardList.forEach(item => {
            if (item.type === key) {
              item[this.$parent.actionName] = data[key]
            }
          })
        }
        this.pendingNumObj[this.$parent.actionName] = [
          data.pendingTask,
          data.pendingApprove,
          data.pendingCheck,
          data.pendingConfirm
        ]
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
