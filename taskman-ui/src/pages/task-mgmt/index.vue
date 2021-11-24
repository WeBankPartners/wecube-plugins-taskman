<template>
  <div>
    <Button @click="backToTask" icon="ios-undo-outline" style="margin-bottom: 8px">{{ $t('back_to_template') }}</Button>
    <div style="width: 84%;margin: 0 auto;">
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
            <Panel :name="dataIndex + ''" :key="dataIndex">
              <template v-if="dataIndex === 0">
                <Tag style="font-size:14px" type="border" size="medium" color="blue"
                  >{{ $t('request_id') }}:{{ data.requestId }}</Tag
                >
                <Tag style="font-size:14px" type="border" size="medium" color="orange"
                  >{{ $t('request_name') }}:{{ data.requestName }}</Tag
                >
                <Tag style="font-size:14px" type="border" size="medium" color="green"
                  >{{ $t('template') }}:{{ data.requestTemplate }}</Tag
                >
                <Tag style="font-size:14px" type="border" size="medium" color="warning"
                  >{{ $t('reporter') }}:{{ data.reporter }}</Tag
                >
                <Tag style="font-size:14px" type="border" size="medium" color="cyan"
                  >{{ $t('report_time') }}:{{ data.reportTime }}</Tag
                >
                <Tag style="font-size:14px" type="border" size="medium" color="blue"
                  >{{ $t('expected_completion_time') }}:{{ data.expectTime }}</Tag
                >
              </template>
              <template v-else-if="dataIndex < dataInfo.length - 1">
                <Tag style="font-size:14px" type="border" size="medium" color="primary"
                  >{{ $t('task_name') }}:{{ data.taskName }}</Tag
                >
                <Tag style="font-size:14px" type="border" size="medium" color="warning"
                  >{{ $t('handler') }}:{{ data.handler }}</Tag
                >
                <Tag style="font-size:14px" type="border" size="medium" color="cyan"
                  >{{ $t('handle_time') }}:{{ data.handleTime }}</Tag
                >
                <Tag style="font-size:14px" type="border" size="medium" color="blue"
                  >{{ $t('expire_time') }}:{{ data.expireTime }}</Tag
                >
              </template>
              <template v-else>
                <Tag style="font-size:14px" type="border" size="medium" color="primary"
                  >{{ $t('task_name') }}:{{ data.taskName }}</Tag
                >
                <template v-if="data.status === 'done'">
                  <Tag style="font-size:14px" type="border" size="medium" color="warning"
                    >{{ $t('handler') }}:{{ data.handler }}</Tag
                  >
                  <Tag style="font-size:14px" type="border" size="medium" color="cyan"
                    >{{ $t('handle_time') }}:{{ data.handleTime }}</Tag
                  >
                  <Tag style="font-size:14px" type="border" size="medium" color="blue"
                    >{{ $t('expire_time') }}:{{ data.expireTime }}</Tag
                  >
                </template>
              </template>
              <p slot="content">
                <template v-if="dataIndex === dataInfo.length - 1 && !enforceDisable">
                  <Upload
                    :action="uploadUrl"
                    :show-upload-list="false"
                    with-credentials
                    style="display:inline-block;margin-left:32px"
                    :headers="headers"
                    :on-success="uploadSucess"
                    :on-error="uploadFailed"
                  >
                    <Button size="small">{{ $t('upload_attachment') }}</Button>
                  </Upload>
                </template>
                <template v-else>
                  <div style="display:inline-block;width:30px"></div>
                </template>
                <Tag
                  type="border"
                  v-for="file in data.attachFiles"
                  :closable="dataIndex === dataInfo.length - 1 && !enforceDisable"
                  checkable
                  :key="file.id"
                  @on-close="removeFile(file)"
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
                      <FormItem :label="$t('process_result')" v-if="data.nextOption.length !== 0">
                        <Select v-model="data.choseOption" :disabled="!data.editable || enforceDisable">
                          <Option v-for="option in data.nextOption" :value="option" :key="option">{{ option }}</Option>
                        </Select>
                      </FormItem>
                      <FormItem :label="$t('process_comments')">
                        <Input :disabled="!data.editable || enforceDisable" v-model="data.comment" type="textarea" />
                      </FormItem>
                    </Form>
                    <div style="text-align:center">
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
    removeFile (file) {
      this.$Modal.confirm({
        title: this.$t('confirm_to_delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const { statusCode } = await deleteAttach(file.id)
          if (statusCode === 'OK') {
            this.getTaskDetail()
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
            let content = JSON.stringify(response.data)
            let fileName = `${file.name}`
            let blob = new Blob([content])
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
    uploadFailed (val) {
      this.$Notice.error({
        title: 'Error',
        desc: val.statusMessage
      })
    },
    async uploadSucess () {
      this.$Notice.success({
        title: 'Successful',
        desc: 'Successful'
      })
      this.getTaskDetail()
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
  height: calc(100vh - 100px);
  margin-top: 24px;
  overflow: auto;
}
</style>
