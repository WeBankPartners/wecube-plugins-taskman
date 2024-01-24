<template>
  <div>
    <Button @click="backToRequest" icon="ios-undo-outline" style="margin-bottom: 8px">{{
      $t('back_to_template')
    }}</Button>
    <div style="width: 84%; margin: 0 auto">
      <div>
        <Steps :current="activeStep">
          <template v-for="(step, stepIndex) in timeStep">
            <Step :title="step.name" :key="stepIndex"></Step>
          </template>
        </Steps>
      </div>
      <div class="task-form">
        <Collapse simple v-model="openPanel">
          <template v-for="(data, dataIndex) in dataInfo">
            <Panel :name="dataIndex + ''" :key="dataIndex" :if-history="data.isHistory ? 'history' : 'current'">
              <!--请求-->
              <div v-if="!data.taskId">
                <Tag style="font-size: 14px" type="border" size="medium" color="blue"
                  >{{ $t('request_id') }}:{{ data.requestId }}</Tag
                >
                <Tag style="font-size: 14px" type="border" size="medium" color="orange"
                  >{{ $t('request_name') }}:{{ data.requestName }}</Tag
                >
                <Tag style="font-size: 14px" type="border" size="medium" color="green"
                  >{{ $t('template') }}:{{ data.requestTemplate }}</Tag
                >
                <Tag style="font-size: 14px" type="border" size="medium" color="warning"
                  >{{ $t('reporter') }}:{{ data.reporter }}</Tag
                >
                <Tag style="font-size: 14px" type="border" size="medium" color="cyan"
                  >{{ $t('report_time') }}:{{ data.reportTime }}</Tag
                >
                <Tag style="font-size: 14px" type="border" size="medium" color="blue"
                  >{{ $t('expected_completion_time') }}:{{ data.expectTime }}</Tag
                >
              </div>
              <!--任务-->
              <div v-else>
                <Tag style="font-size: 14px" type="border" size="medium" color="primary"
                  >{{ $t('task_name') }}:{{ data.taskName }}</Tag
                >
                <Tag style="font-size: 14px" type="border" size="medium" color="warning"
                  >{{ $t('handler') }}:{{ data.handler }}</Tag
                >
                <Tag style="font-size: 14px" type="border" size="medium" color="warning"
                  >{{ $t('handler_role') }}:{{ data.handleRoleName }}</Tag
                >
                <Tag v-if="data.status === 'done'" style="font-size: 14px" type="border" size="medium" color="cyan"
                  >{{ $t('handle_time') }}:{{ data.handleTime }}</Tag
                >
                <Tag style="font-size: 14px" type="border" size="medium" color="cyan"
                  >{{ $t('report_time') }}:{{ data.reportTime }}</Tag
                >
                <Tag style="font-size: 14px" type="border" size="medium" color="blue"
                  >{{ $t('expire_time') }}:{{ data.expireTime }}</Tag
                >
                <Tag v-if="data.isHistory" style="font-size: 14px;" type="border" size="medium" color="magenta"
                  ><span class="history-comment">{{ $t('process_comments') }}: {{ data.comment }}</span></Tag
                >
              </div>
              <p slot="content">
                <Tag
                  type="border"
                  v-for="file in data.attachFiles"
                  x
                  checkable
                  :key="file.id"
                  @on-change="downloadFile(file)"
                  color="primary"
                  >{{ file.name }}</Tag
                >
                <Tabs :value="data.activeTab">
                  <template v-for="form in data.formData">
                    <TabPane :label="form.itemGroup" :name="form.itemGroup" :key="form.itemGroup">
                      <RequestCheckData :data="form" :requestId="requestId"></RequestCheckData>
                    </TabPane>
                  </template>
                </Tabs>
                <span>
                  <div v-if="dataIndex !== 0">
                    <Form :label-width="80" style="margin: 16px 0">
                      <FormItem :label="$t('process_result')" v-if="data.nextOption && data.nextOption.length !== 0">
                        <span slot="label">
                          {{ $t('process_result') }}
                          <span style="color: #ed4014"> * </span>
                        </span>
                        <Select v-model="data.choseOption" disabled>
                          <Option v-for="option in data.nextOption" :value="option" :key="option">{{ option }}</Option>
                        </Select>
                      </FormItem>
                      <FormItem :label="$t('process_comments')">
                        <Input disabled v-model="data.comment" type="textarea" />
                      </FormItem>
                    </Form>
                  </div>
                </span>
              </p>
            </Panel>
          </template>
        </Collapse>
      </div>
    </div>
  </div>
</template>

<script>
import { getRequestDetail } from '@/api/server.js'
import { getCookie } from '@/pages/util/cookie'
import axios from 'axios'
import RequestCheckData from './request-check-data'
export default {
  name: '',
  data () {
    return {
      activeStep: 0,
      taskId: '',
      requestId: '',
      timeStep: [],
      openPanel: '',
      dataInfo: [],
      uploadUrl: '',
      headers: {}
    }
  },
  mounted () {
    this.requestId = this.$route.query.requestId
    this.jumpFrom = this.$route.query.jumpFrom
    this.uploadUrl = `/taskman/api/v1/task/attach-file/upload/${this.taskId}`
    const accessToken = getCookie('accessToken')
    this.headers = {
      Authorization: 'Bearer ' + accessToken
    }
    this.getRequestDetail()
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
          this.isExport = false
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
        .catch(() => {
          this.$Message.warning('Error')
        })
    },
    backToRequest () {
      this.$router.push({ path: '/taskman/request-mgmt', query: { activeTab: this.jumpFrom } })
    },
    success () {
      this.$Notice.success({
        title: this.$t('successful'),
        desc: this.$t('successful')
      })
    },
    async getRequestDetail () {
      const { statusCode, data } = await getRequestDetail(this.requestId)
      if (statusCode === 'OK') {
        // this.requestId = data.request
        this.timeStep = data.timeStep
        this.activeStep = this.timeStep.findIndex(t => t.active === true)
        this.dataInfo = data.data.map(d => {
          this.requestId = d.requestId
          d.activeTab = ''
          if (d.formData.length > 0) {
            d.activeTab = d.formData[0].itemGroup
          }
          return d
        })

        this.openPanel = data.data.length - 1 + ''
      }
    }
  },
  components: {
    RequestCheckData
  }
}
</script>

<style scoped lang="scss">
.task-form {
  height: calc(100vh - 100px);
  margin-top: 24px;
  overflow: auto;
}
</style>

<style scoped>
.task-form >>> .ivu-collapse-header {
  height: auto !important;
  display: flex;
  align-items: center;
}
.task-form >>> .ivu-collapse-item[if-history='history'] {
  background: #ddd;
}
.task-form >>> .ivu-collapse-item[if-history='history'] .ivu-collapse-content {
  background: #ddd;
}
</style>
