<template>
  <div>
    <Button @click="backToTask" icon="ios-undo-outline" style="margin-bottom: 8px">{{ $t('back_to_template') }}</Button>
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
              <div v-else>
                <Tag style="font-size: 14px" type="border" size="medium" color="primary"
                  >{{ $t('task_name') }}:{{ data.taskName }}</Tag
                >
                <Tag v-if="data.status === 'done'" style="font-size: 14px" type="border" size="medium" color="warning"
                  >{{ $t('handler') }}:{{ data.handler }}</Tag
                >
                <Tag v-if="data.status === 'done'" style="font-size: 14px" type="border" size="medium" color="warning"
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
                <template v-if="dataIndex === dataInfo.length - 1 && !enforceDisable">
                  <Upload
                    :action="uploadUrl"
                    :before-upload="handleUpload"
                    :show-upload-list="false"
                    with-credentials
                    style="display: inline-block; margin-left: 32px"
                    :headers="headers"
                    :on-success="res => uploadSucess(res, dataIndex)"
                    :on-error="uploadFailed"
                  >
                    <Button size="small" type="success">{{ $t('upload_attachment') }}</Button>
                  </Upload>
                </template>
                <template v-else>
                  <div style="display: inline-block; width: 30px"></div>
                </template>
                <Tag
                  type="border"
                  v-for="file in data.attachFiles"
                  :closable="dataIndex === dataInfo.length - 1 && !enforceDisable"
                  checkable
                  :key="file.id"
                  @on-close="removeFile(file, dataIndex)"
                  @on-change="downloadFile(file)"
                  color="primary"
                  >{{ file.name }}</Tag
                >
                <Tabs :value="data.activeTab">
                  <template v-for="form in data.formData">
                    <TabPane :label="form.itemGroupName" :name="form.itemGroup" :key="form.itemGroup">
                      <div :style="{ 'max-height': MODALHEIGHT + 'px', overflow: 'auto' }">
                        <TaskData
                          :data="form"
                          :isDisabled="!data.editable"
                          :requestId="requestId"
                          :enforceDisable="enforceDisable"
                        ></TaskData>
                      </div>
                    </TabPane>
                  </template>
                </Tabs>
                <span>
                  <div v-if="dataIndex !== 0 || data.requestId === ''">
                    <Form :label-width="80" style="margin: 16px 0">
                      <FormItem v-if="data.requestId === ''" :label="$t('task') + $t('description')">
                        <Input disabled v-model="data.description" type="textarea" />
                      </FormItem>
                      <FormItem :label="$t('process_result')" v-if="data.nextOption && data.nextOption.length !== 0">
                        <span slot="label">
                          {{ $t('process_result') }}
                          <span style="color: #ed4014"> * </span>
                        </span>
                        <Select v-model="data.choseOption" :disabled="!data.editable || enforceDisable">
                          <Option v-for="option in data.nextOption" :value="option" :key="option">{{ option }}</Option>
                        </Select>
                      </FormItem>
                      <FormItem :label="$t('process_comments')">
                        <Input :disabled="!data.editable || enforceDisable" v-model="data.comment" type="textarea" />
                      </FormItem>
                    </Form>
                    <div style="text-align: center">
                      <Button v-if="data.editable" :disabled="enforceDisable" @click="saveTaskData" type="info">{{
                        $t('save')
                      }}</Button>
                      <Button v-if="data.editable" :disabled="enforceDisable" @click="commitTaskData" type="primary">{{
                        $t('commit')
                      }}</Button>
                    </div>
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
import { getTaskDetail, saveTaskData, commitTaskData, deleteAttach } from '@/api/server.js'
import TaskData from './task-data'
import { getCookie } from '@/pages/util/cookie'
import axios from 'axios'
export default {
  name: '',
  data () {
    return {
      MODALHEIGHT: 500,
      activeStep: 0,
      taskId: '',
      requestId: '',
      enforceDisable: false,
      timeStep: [],
      openPanel: '',
      dataInfo: [],
      uploadUrl: '',
      headers: {}
    }
  },
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 400
    this.taskId = this.$route.query.taskId
    this.enforceDisable = this.$route.query.enforceDisable === 'Y'
    this.uploadUrl = `/taskman/api/v1/task/attach-file/upload/${this.taskId}`
    const accessToken = getCookie('accessToken')
    this.headers = {
      Authorization: 'Bearer ' + accessToken
    }
    this.getTaskDetail()
  },
  methods: {
    handleUpload (file) {
      this.$Message.info(this.$t('upload_tip'))
      return true
    },
    removeFile (file, index) {
      this.$Modal.confirm({
        title: this.$t('confirm_to_delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const { statusCode, data } = await deleteAttach(file.id)
          if (statusCode === 'OK') {
            this.dataInfo[index].attachFiles = data
          }
        },
        onCancel: () => {}
      })
    },
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
            // let content = JSON.stringify(response.data)
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
    uploadFailed (val, response) {
      console.log(val)
      this.$Notice.error({
        title: 'Error',
        desc: response.statusMessage
      })
    },
    async uploadSucess (item, index) {
      this.$Notice.success({
        title: 'Successful',
        desc: 'Successful'
      })
      this.dataInfo[index].attachFiles = item.data
    },
    backToTask () {
      this.$router.push({ path: '/taskman/task-mgmt' })
    },
    success () {
      this.$Notice.success({
        title: this.$t('successful'),
        desc: this.$t('successful')
      })
    },
    async getTaskDetail () {
      const { statusCode, data } = await getTaskDetail(this.taskId)
      if (statusCode === 'OK') {
        this.requestId = data.request
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
    },
    async saveTaskData () {
      const taskData = this.dataInfo.find(d => d.editable === true)
      const result = this.paramsCheck(taskData)
      if (result) {
        const { statusCode } = await saveTaskData(this.taskId, taskData)
        if (statusCode === 'OK') {
          this.success()
        }
      } else {
        this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('required_tip')
        })
      }
    },
    async commitTaskData () {
      const taskData = this.dataInfo.find(d => d.editable === true)
      const { statusCode } = await commitTaskData(this.taskId, taskData)
      if (statusCode === 'OK') {
        this.success()
        this.$router.push({ path: '/taskman/task-mgmt' })
      }
    },
    paramsCheck (taskData) {
      let result = true
      taskData.formData.forEach(requestData => {
        let requiredName = []
        requestData.title.forEach(t => {
          if (t.required === 'yes') {
            requiredName.push(t.name)
          }
        })
        requestData.value.forEach(v => {
          requiredName.forEach(key => {
            let val = v.entityData[key]
            if (Array.isArray(val)) {
              if (val.length === 0) {
                result = false
              }
            } else {
              if (val === '') {
                result = false
              }
            }
          })
        })
      })
      return result
    }
  },
  components: {
    TaskData
  }
}
</script>

<style scoped lang="scss">
.task-form {
  height: calc(100vh - 200px);
  margin-top: 24px;
  overflow: auto;
}
.history-comment {
  display: inline-block;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  word-break: break-all;
  white-space: nowrap;
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
