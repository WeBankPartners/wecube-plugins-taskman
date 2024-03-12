<template>
  <div class="taskman-workbench-data-card">
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
          <!--发布-->
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
              <span v-if="i.type === 'pending' && getPendingNum(j.type) > 0" class="badge">{{
                getPendingNum(j.type)
              }}</span>
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
      active: '',
      action: '',
      cardList: [
        {
          type: 'pending',
          label: this.$t('tw_pending'),
          icon: 'ios-alert',
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
    async getData (init = false) {
      // ini为true，初始化拉取所有数据，后续拉取特定场景下的数据
      const params = {
        scene: init ? 0 : Number(this.action)
      }
      const { statusCode, data } = await overviewData(params)
      if (statusCode === 'OK') {
        for (let key in data) {
          this.cardList.forEach(item => {
            if (item.type === key) {
              if (init) {
                data[key].forEach((number, index) => {
                  item[String(index + 1)] = number || 0
                })
              } else {
                item[this.action] = data[key][0]
              }
            }
          })
        }
        if (init) {
          data.pendingTask.forEach((_, index) => {
            this.pendingNumObj[String(index + 1)] = [
              data.pendingTask[index] || 0,
              data.pendingApprove[index] || 0,
              data.pendingCheck[index] || 0,
              data.pendingConfirm[index] || 0
            ]
          })
        } else {
          this.pendingNumObj[this.action] = [
            data.pendingTask[0] || 0,
            data.pendingApprove[0] || 0,
            data.pendingCheck[0] || 0,
            data.pendingConfirm[0] || 0
          ]
        }
      }
    },
    handleTabChange: debounce(function (item, subType) {
      this.active = item.type
      this.action = subType || '1'
      this.$emit('fetchData', item.type, this.action)
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
      span {
        color: rgba(16, 16, 16, 1);
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
