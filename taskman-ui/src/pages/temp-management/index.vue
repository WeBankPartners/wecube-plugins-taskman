<template>
  <div>
    <Button @click="backToTemplate" icon="ios-undo-outline" style="margin-bottom: 8px">{{
      $t('back_to_template')
    }}</Button>
    <Row type="flex">
      <Col span="20" offset="1">
        <Steps :current="currentStep">
          <Step icon="ios-add-circle">
            <span slot="title" @click="changeStep(0)">{{ $t('basic_information_settings') }}</span>
          </Step>
          <Step icon="md-apps">
            <span slot="title" @click="changeStep(1)">{{ $t('form_item_selection') }}</span>
          </Step>
          <Step icon="md-cog">
            <span slot="title" @click="changeStep(2)">{{ $t('request_form_settings') }}</span>
          </Step>
          <Step icon="ios-settings">
            <span slot="title" @click="changeStep(3)">{{ $t('task_form_settings') }}</span>
          </Step>
        </Steps>
      </Col>
      <Col span="3">
        <Button
          @click="confirmTemplate"
          :disabled="currentStep !== 3 || $parent.isCheck === 'Y'"
          size="small"
          type="primary"
          >{{ $t('publish_template') }}</Button
        >
      </Col>
    </Row>
    <div v-if="currentStep !== -1" style="margin-top:16px;">
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
import { confirmTemplate } from '@/api/server.js'
export default {
  name: '',
  data () {
    return {
      currentStep: -1,
      requestTemplateId: '',
      isCheck: 'N'
    }
  },
  mounted () {
    if (this.$route.query.requestTemplateId !== '') {
      this.requestTemplateId = this.$route.query.requestTemplateId
    }
    this.isCheck = this.$route.query.isCheck
    this.currentStep = 0
  },
  methods: {
    async confirmTemplate () {
      const { statusCode } = await confirmTemplate(this.requestTemplateId)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.$router.push({ path: '/taskman/template-mgmt' })
      }
    },
    changeStep (val) {
      this.currentStep = val
    },
    backToTemplate () {
      this.$router.push({ path: '/taskman/template-mgmt' })
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
