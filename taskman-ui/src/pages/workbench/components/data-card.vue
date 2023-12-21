<template>
  <div class="workbench-data-card">
    <Card v-for="(i, index) in data" :key="index" border :style="activeStyles(i)">
      <div class="content" @click="handleTabChange(i)">
        <div class="w-header">
          <img :src="i.icon" />
          <span>{{ `${i.label}(${i.total || ''})` }}</span>
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
          icon: require('@/images/warning.png'),
          requestNum: 0,
          publishNum: 0
        },
        {
          type: 'hasProcessed',
          label: this.$t('tw_hasProcessed'),
          icon: require('@/images/checked.png'),
          requestNum: 0,
          publishNum: 0
        },
        {
          type: 'submit',
          label: this.$t('tw_submit'),
          icon: require('@/images/upload.png'),
          requestNum: 0,
          publishNum: 0
        },
        {
          type: 'draft',
          label: this.$t('tw_draft'),
          icon: require('@/images/save.png'),
          requestNum: 0,
          publishNum: 0
        },
        {
          type: 'collect',
          label: this.$t('tw_collect'),
          icon: require('@/images/star.png'),
          requestNum: 0,
          publishNum: 0
        }
      ]
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
</style>
<style lang="scss" scoped>
.workbench-data-card {
  display: flex;
  .content {
    .w-header {
      display: flex;
      align-items: flex-end;
      img {
        width: 24px;
        height: 24px;
        margin-right: 5px;
      }
      span {
        color: rgba(16, 16, 16, 1);
        font-size: 14px;
        font-family: PingFangSC-regular;
        font-weight: 500;
      }
    }
    .data {
      display: flex;
      justify-content: space-around;
      align-items: center;
      margin-top: 10px;
      .list {
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
