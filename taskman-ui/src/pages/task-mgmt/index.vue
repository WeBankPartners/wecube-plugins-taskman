<template>
  <div style="margin: 24px">
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
                <Tag style="font-size:14px" type="border" size="medium" color="primary"
                  >{{ $t('request_name') }}:{{ data.requestName }}</Tag
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
                  >{{ $t('handler') }}:{{ data.reporter }}</Tag
                >
                <Tag style="font-size:14px" type="border" size="medium" color="cyan"
                  >{{ $t('handle_time') }}:{{ data.reportTime }}</Tag
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
                    >{{ $t('handler') }}:{{ data.reporter }}</Tag
                  >
                  <Tag style="font-size:14px" type="border" size="medium" color="cyan"
                    >{{ $t('handle_time') }}:{{ data.reportTime }}</Tag
                  >
                  <Tag style="font-size:14px" type="border" size="medium" color="blue"
                    >{{ $t('expire_time') }}:{{ data.expireTime }}</Tag
                  >
                </template>
              </template>
              <p slot="content">
                <Tabs :value="data.activeTab">
                  <template v-for="form in data.formData">
                    <TabPane :label="form.itemGroup" :name="form.itemGroup" :key="form.itemGroup">
                      <TaskData
                        :data="form"
                        :isDisabled="!data.editable"
                        :requestId="requestId"
                        :enforceDisable="enforceDisable"
                      ></TaskData>
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
                    <Button v-if="data.editable" :disabled="enforceDisable" @click="saveTaskData" type="info">{{
                      $t('save')
                    }}</Button>
                    <Button v-if="data.editable" :disabled="enforceDisable" @click="commitTaskData" type="primary">{{
                      $t('commit')
                    }}</Button>
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
import { getTaskDetail, saveTaskData, commitTaskData } from '@/api/server.js'
import TaskData from './task-data'
export default {
  name: '',
  data () {
    return {
      activeStep: 0,
      taskId: '',
      requestId: '',
      enforceDisable: false,
      timeStep: [],
      openPanel: '',
      dataInfo: []
    }
  },
  mounted () {
    this.taskId = this.$route.query.taskId
    this.enforceDisable = this.$route.query.enforceDisable === 'Y'
    this.getTaskDetail()
  },
  methods: {
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
      const taskData = this.dataInfo.find(d => d.editable === true).formData
      const { statusCode } = await saveTaskData(this.taskId, taskData)
      if (statusCode === 'OK') {
        this.success()
      }
    },
    async commitTaskData () {
      await this.saveTaskData()
      const taskData = this.dataInfo.find(d => d.editable === true)
      const commitData = {
        comment: taskData.comment,
        choseOption: taskData.choseOption
      }
      const { statusCode } = await commitTaskData(this.taskId, commitData)
      if (statusCode === 'OK') {
        this.success()
        this.$router.push({ path: '/taskman/task-mgmt' })
      }
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
