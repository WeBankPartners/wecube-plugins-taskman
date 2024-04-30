<template>
  <div class="workbench-header-tag">
    <Row v-if="showHeader" class="title" :gutter="10">
      <Col :span="3" class="line">{{ $t('handler_role') }}</Col>
      <Col :span="2" class="line">{{ $t('handler') }}</Col>
      <Col :span="2" class="line">{{ $t('t_action') }}</Col>
      <Col :span="2" class="line">{{ $t('tw_handleStatus') }}</Col>
      <Col :span="3" class="line">{{ $t('handle_time') }}</Col>
      <Col :span="3" class="line">{{ $t('tw_assume') }}</Col>
      <Col :span="5" class="line">{{ $t('tw_note') }}</Col>
      <Col :span="4" class="line">{{ $t('tw_attach') }}</Col>
    </Row>
    <template v-if="data.taskHandleList && data.taskHandleList.length > 0">
      <Row v-for="i in data.taskHandleList" :key="i.id" class="content" :gutter="10">
        <Col :span="3" class="line">{{ i.role || '-' }}</Col>
        <Col :span="2" class="line">{{ i.handler || '-' }}</Col>
        <Col :span="2" class="line">{{ getOperationName(i) }}</Col>
        <Col :span="2" class="line">{{ getHandleStatus(i) }}</Col>
        <Col :span="3" class="line">{{ i.updatedTime || '-' }}</Col>
        <Col :span="3" class="line">{{ getDiffTime(i) || '-' }}</Col>
        <Col :span="5" class="line">
          <Tooltip max-width="300" :content="i.resultDesc">
            <span class="text-overflow">{{ i.resultDesc || '-' }}</span>
          </Tooltip>
        </Col>
        <Col :span="4" class="line">
          <div v-for="file in i.attachFiles" style="display:inline-block;" :key="file.id">
            <Tag type="border" :closable="false" checkable @on-change="downloadFile(file)" color="primary">
              {{ file.name }}
            </Tag>
          </div>
        </Col>
      </Row>
    </template>
    <Row v-else class="content" :gutter="10">
      <Col :span="3" class="line">{{ data.role || '-' }}</Col>
      <Col :span="2" class="line">{{ data.handler || '-' }}</Col>
      <Col :span="2" class="line">{{ data.choseOption || operation || '-' }}</Col>
      <Col :span="2" class="line">{{ '-' }}</Col>
      <Col :span="3" class="line">{{ data.updatedTime || '-' }}</Col>
      <Col :span="3" class="line">{{ getDiffTime(data) || '-' }}</Col>
      <Col :span="5" class="line">
        <Tooltip max-width="300" :content="data.comment">
          <span style="text-overflow">{{ data.comment }}</span>
        </Tooltip>
      </Col>
      <Col :span="4" class="line">
        <div v-for="file in data.attachFiles" style="display:inline-block;" :key="file.id">
          <Tag type="border" :closable="false" checkable @on-change="downloadFile(file)" color="primary">
            {{ file.name }}
          </Tag>
        </div>
      </Col>
    </Row>
  </div>
</template>

<script>
import axios from 'axios'
import { getCookie } from '@/pages/util/cookie'
import dayjs from 'dayjs'
export default {
  props: {
    data: {
      type: Object,
      default: () => {}
    },
    showHeader: {
      type: Boolean,
      default: false
    },
    operation: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      headers: {}
    }
  },
  computed: {
    getOperationName () {
      return function (i) {
        // 1.编排任务接口下发，2审批和定版退回，前端枚举
        const resultMap = {
          deny: this.$t('tw_reject'), // 拒绝
          approve: this.$t('tw_approve'), // 同意
          redraw: this.$t('tw_send_back'), // 退回
          unrelated: this.$t('tw_unrelated')
        }
        let resultName = ''
        if (['approve', 'check'].includes(this.data.type)) {
          resultName = resultMap[i.handleResult]
        }
        if (this.data.type === 'implement_process') {
          resultName = i.handleResult
        }
        return resultName || this.operation || '-'
      }
    },
    getHandleStatus () {
      return function (i) {
        // 任务处理状态
        const statusMap = {
          complete: this.$t('tw_completed'),
          uncompleted: this.$t('tw_incomplete')
        }
        let statusName = ''
        if (['implement_custom', 'implement_process'].includes(this.data.type)) {
          statusName = statusMap[i.handleStatus]
        }
        return statusName || '-'
      }
    },
    getDiffTime () {
      return function (data) {
        const newDate = dayjs(data.updatedTime)
        const oldDate = dayjs(data.createdTime)
        let subtime = (newDate - oldDate) / 1000
        let days = parseInt(subtime / 86400)
        let hours = parseInt(subtime / 3600) - 24 * days
        let mins = parseInt((subtime % 3600) / 60)
        let secs = parseInt(subtime % 60)
        return (
          (days > 0 ? days + this.$t('tw_day') : '') +
          (hours > 0 ? hours + this.$t('tw_hour') : '') +
          (mins > 0 ? mins + this.$t('tw_minute') : '') +
          (secs > 0 ? secs + this.$t('tw_second') : '')
        )
      }
    }
  },
  mounted () {
    const accessToken = getCookie('accessToken')
    this.headers = {
      Authorization: 'Bearer ' + accessToken
    }
  },
  methods: {
    async downloadFile (file) {
      axios({
        method: 'GET',
        url: `/taskman/api/v1/request/attach-file/download/${file.id}`,
        headers: this.headers,
        responseType: 'blob'
      })
        .then(response => {
          if (response.status < 400) {
            let fileName = `${file.name}`
            let blob = new Blob([response.data])
            if ('msSaveOrOpenBlob' in navigator) {
              window.navigator.msSaveOrOpenBlob(blob, fileName)
            } else {
              if ('download' in document.createElement('a')) {
                // 非IE下载
                let elink = document.createElement('a')
                elink.download = fileName
                elink.style.display = 'none'
                elink.href = URL.createObjectURL(blob)
                document.body.appendChild(elink)
                elink.click()
                URL.revokeObjectURL(elink.href) // 释放URL 对象
                document.body.removeChild(elink)
              } else {
                // IE10+下载
                navigator.msSaveOrOpenBlob(blob, fileName)
              }
            }
          }
        })
        .catch(error => {
          console.log(error)
          this.$Message.warning('Error')
        })
    }
  }
}
</script>

<style lang="scss">
.workbench-header-tag {
  width: calc(100% - 10px);
  margin-left: 10px;
  .title {
    background: #f7f7f7;
    text-align: left;
    margin-top: -36px;
    height: 36px;
    line-height: 36px;
    padding-left: 20px;
    font-size: 12px;
    font-weight: bold;
  }
  .content {
    display: flex;
    align-items: center;
    text-align: left;
    padding-left: 20px;
    background: #f0faff;
    font-size: 12px;
    overflow: hidden;
    word-break: break-all;
    min-height: 50px;
    .text-overflow {
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 1;
      -webkit-box-orient: vertical;
    }
  }
  .line {
    line-height: 32px;
  }
  .ivu-tooltip {
    display: flex;
  }
}
</style>
