<template>
  <div>
    <Button @click="backToRequest" icon="ios-undo-outline" style="margin-bottom: 8px">{{
      $t('back_to_template')
    }}</Button>
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
                      <RequestCheckData
                        :data="form"
                        :isDisabled="!data.editable"
                        :requestId="requestId"
                        :enforceDisable="enforceDisable"
                      ></RequestCheckData>
                    </TabPane>
                  </template>
                </Tabs>
                <span>
                  <div v-if="dataIndex !== 0">
                    <Form :label-width="80" style="margin: 16px 0">
                      <FormItem :label="$t('process_result')" v-if="data.nextOption.length !== 0">
                        <Select v-model="data.choseOption" :disabled="!data.editable || enforceDisable">
                          <Option v-for="option in data.nextOption" :value="option" :key="option">{{ option }}</Option>
                        </Select>
                      </FormItem>
                      <FormItem :label="$t('process_comments')">
                        <Input :disabled="!data.editable || enforceDisable" v-model="data.comment" type="textarea" />
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
import RequestCheckData from './request-check-data'
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
    this.requestId = this.$route.query.requestId
    this.jumpFrom = this.$route.query.jumpFrom
    this.enforceDisable = this.$route.query.enforceDisable === 'Y'
    this.getRequestDetail()
  },
  methods: {
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
