<template>
  <div class="workbench-header-tag">
    <Row v-if="showHeader" class="title" :gutter="10">
      <Col :span="3" class="line">{{ $t('handler_role') }}</Col>
      <Col :span="2" class="line">{{ $t('handler') }}</Col>
      <Col :span="2" class="line">{{ $t('t_action') }}</Col>
      <Col :span="4" class="line">{{ $t('handle_time') }}</Col>
      <Col :span="3" class="line">{{ $t('tw_assume') }}</Col>
      <Col :span="5" class="line">{{ $t('tw_note') }}</Col>
      <Col :span="5" class="line">{{ $t('tw_attach') }}</Col>
    </Row>
    <Row class="content" :gutter="10">
      <Col :span="3" class="line">{{ data.handleRoleName || '-' }}</Col>
      <Col :span="2" class="line">{{ data.handler || '-' }}</Col>
      <Col :span="2" class="line">{{ data.choseOption || operation }}</Col>
      <Col :span="4" class="line">{{ data.handleTime }}</Col>
      <Col :span="3" class="line">{{ getDiffTime }}</Col>
      <Col :span="5" class="line"
        ><div class="text-overflow">{{ data.comment }}</div></Col
      >
      <Col :span="5" class="line">
        <div v-for="file in data.attachFiles" style="display:inline-block;" :key="file.id">
          <Tag type="border" :closable="false" checkable @on-change="downloadFile(file)" color="primary">{{
            file.name
          }}</Tag>
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
    getDiffTime () {
      const newDate = dayjs(this.data.handleTime)
      const oldDate = dayjs(this.data.createTime)
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
    .text-overflow {
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }
  }
  .line {
    line-height: 32px;
  }
}
</style>
