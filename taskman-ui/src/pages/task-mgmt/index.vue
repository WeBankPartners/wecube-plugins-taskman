<template>
  <div>
    <div style="width:50%;margin: 0 auto;">
      <Steps :current="1">
        <template v-for="(step, stepIndex) in timeStep">
          <Step :title="step.name" :key="stepIndex"></Step>
        </template>
      </Steps>
    </div>
    <div>
      <template v-for="(data, dataIndex) in dataInfo">
        <div :key="dataIndex">
          <Tag>{{ $t('reporter') }}:{{ data.reporter }}</Tag>
          <Tag>{{ $t('report_time') }}:{{ data.reportTime }}</Tag>
        </div>
        <Tabs :key="dataIndex" :value="data.activeTab" @on-click="changeTab">
          <template v-for="entity in requestData">
            <TabPane :label="entity.entity" :name="entity.entity" :key="entity.entity">
              <DataMgmt ref="dataMgmt" @getEntityData="getEntityData"></DataMgmt>
            </TabPane>
          </template>
        </Tabs>
      </template>
    </div>
  </div>
</template>

<script>
import { getTaskDetail } from '@/api/server.js'
export default {
  name: '',
  data () {
    return {
      taskId: '',
      timeStep: [],
      dataInfo: []
    }
  },
  mounted () {
    if (this.$route.query.taskId !== '') {
      this.taskId = this.$route.query.taskId
      this.getTaskDetail()
    }
  },
  methods: {
    async getTaskDetail () {
      const { statusCode, data } = await getTaskDetail('61692c6fd4d6f603')
      if (statusCode === 'OK') {
        this.timeStep = data.timeStep
        this.dataInfo = data.data.map(d => {
          d.activeTab = d.formData[0]
        })
        console.log(data)
      }
    }
  },
  components: {}
}
</script>

<style scoped lang="scss"></style>
