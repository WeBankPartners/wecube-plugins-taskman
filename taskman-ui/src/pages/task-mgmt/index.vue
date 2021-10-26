<template>
  <div style="width:60%;margin: 0 auto;">
    <div>
      <Steps :current="1">
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
              <Tag>{{ $t('request_name') }}:{{ data.requestName }}</Tag>
              <Tag>{{ $t('reporter') }}:{{ data.reporter }}</Tag>
              <Tag>{{ $t('report_time') }}:{{ data.reportTime }}</Tag>
            </template>
            <template v-else>
              <Tag>{{ $t('handler') }}:{{ data.requestName }}</Tag>
              <Tag>{{ $t('handle_time') }}:{{ data.reportTime }}</Tag>
            </template>
            <p slot="content">
              <Tabs :value="data.activeTab">
                <template v-for="form in data.formData">
                  <TabPane :label="form.itemGroup" :name="form.itemGroup" :key="form.itemGroup">
                    <TaskData :data="form" :isDisabled="!data.editable" :enforceDisable="enforceDisable"></TaskData>
                  </TabPane>
                </template>
              </Tabs>
              <span>
                <div v-if="dataIndex !== 0">
                  <Form :label-width="80" style="margin: 16px 0">
                    <FormItem :label="$t('approval_result')" v-if="data.nextOption.length !== 0">
                      <Select v-model="data.choseOption" :disabled="!data.editable">
                        <Option v-for="option in data.nextOption" :value="option" :key="option">{{ option }}</Option>
                      </Select>
                    </FormItem>
                    <FormItem :label="$t('approval_comments')">
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
</template>

<script>
import { getTaskDetail, saveTaskData, commitTaskData } from '@/api/server.js'
import TaskData from './task-data'
export default {
  name: '',
  data () {
    return {
      taskId: '',
      enforceDisable: false,
      timeStep: [],
      openPanel: '',
      dataInfo: []
    }
  },
  mounted () {
    this.taskId = this.$route.query.taskId
    this.enforceDisable = this.$route.query.enforceDisable
    this.getTaskDetail()
  },
  methods: {
    success () {
      this.$Notice.success({
        title: this.$t('successful'),
        desc: this.$t('successful')
      })
    },
    async getTaskDetail () {
      const { statusCode, data } = await getTaskDetail(this.taskId)
      if (statusCode === 'OK') {
        this.timeStep = data.timeStep
        this.dataInfo = data.data.map(d => {
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
      const taskData = this.dataInfo.find(d => d.editable === true)
      const commitData = {
        comment: taskData.comment,
        choseOption: taskData.choseOption
      }
      const { statusCode } = await commitTaskData(this.taskId, commitData)
      if (statusCode === 'OK') {
        this.success()
        this.$router.push({ path: '/task' })
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
  overflow: auto;
}
</style>
