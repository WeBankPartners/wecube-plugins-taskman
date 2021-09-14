<template>
  <div style="margin: 24px">
    <Steps :current="currentStep">
      <Step title="基础信息设置" icon="ios-add-circle"></Step>
      <Step title="表单项选择" icon="md-apps"></Step>
      <Step title="请求表单设置" icon="md-cog"></Step>
      <Step title="任务表单设置" icon="ios-settings"></Step>
    </Steps>
    <div style="margin-top:48px;">
      <BasicInfo @basicInfoNextStep="basicInfoNextStep" v-if="currentStep === 0"></BasicInfo>
      <FormSelect
        @formSelectNextStep="formSelectNextStep"
        :requestTemplateId="requestTemplateId"
        v-if="currentStep === 1"
      ></FormSelect>
      <RequestForm v-if="currentStep === 2"></RequestForm>
      <TaskForm v-if="currentStep === 3"></TaskForm>
    </div>
  </div>
</template>

<script>
import FormSelect from './form-select'
import RequestForm from './request-form'
import BasicInfo from './basic-info'
import TaskForm from './task-form'
export default {
  name: '',
  data () {
    return {
      currentStep: 1,
      requestTemplateId: ''
    }
  },
  mounted () {},
  methods: {
    basicInfoNextStep (data) {
      this.requestTemplateId = data.id
      this.currentStep++
    },
    formSelectNextStep () {
      this.currentStep++
    }
  },
  components: {
    FormSelect,
    RequestForm,
    BasicInfo,
    TaskForm
  }
}
</script>

<style scoped lang="scss">
.header-icon {
  float: right;
  margin: 3px 40px 0 0 !important;
}
</style>
