<template>
  <div style="margin: 24px">
    <Button @click="backToTemplate" icon="ios-undo-outline" style="margin-bottom: 8px">{{
      $t('back_to_template')
    }}</Button>
    <Steps :current="currentStep">
      <Step :title="$t('basic_information_settings')" icon="ios-add-circle"></Step>
      <Step :title="$t('form_item_selection')" icon="md-apps"></Step>
      <Step :title="$t('request_form_settings')" icon="md-cog"></Step>
      <Step :title="$t('task_form_settings')" icon="ios-settings"></Step>
    </Steps>
    <div v-if="currentStep !== -1" style="margin-top:48px;">
      <BasicInfo
        @basicInfoNextStep="basicInfoNextStep"
        :requestTemplateId="requestTemplateId"
        v-if="currentStep === 0"
      ></BasicInfo>
      <FormSelect
        @formSelectNextStep="formSelectNextStep"
        :requestTemplateId="requestTemplateId"
        v-if="currentStep === 1"
      ></FormSelect>
      <RequestForm
        :requestTemplateId="requestTemplateId"
        @formSelectNextStep="formSelectNextStep"
        v-if="currentStep === 2"
      ></RequestForm>
      <TaskForm :requestTemplateId="requestTemplateId" v-if="currentStep === 3"></TaskForm>
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
      currentStep: -1,
      requestTemplateId: ''
    }
  },
  mounted () {
    if (this.$route.query.requestTemplateId !== '') {
      this.requestTemplateId = this.$route.query.requestTemplateId
    }
    this.currentStep = 0
  },
  methods: {
    backToTemplate () {
      this.$router.push({ path: '/' })
    },
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
