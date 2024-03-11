<template>
  <div>
    <Button @click="backToTemplate" icon="ios-undo-outline" style="margin-bottom: 8px">{{
      $t('back_to_template')
    }}</Button>
    <Row type="flex">
      <Col span="17" offset="2">
        <Steps :current="currentStep">
          <Step :title="$t('basic_information_settings')"></Step>
          <Step :title="$t('request_form_settings')"></Step>
          <Step :title="$t('approval_form_settings')"></Step>
          <Step :title="$t('task_form_settings')"></Step>
        </Steps>
      </Col>
      <!-- <Col span="3">
        <Button
          @click="submitTemplate"
          :disabled="currentStep !== 3 || $parent.isCheck === 'Y'"
          size="small"
          type="primary"
          >{{ $t('submit_for_review') }}</Button
        >
      </Col> -->
    </Row>
    <div></div>
    <div v-if="currentStep !== -1" style="margin-top:16px;">
      <BasicInfo @gotoStep="gotoStep" :requestTemplateId="requestTemplateId" v-if="currentStep === 0"></BasicInfo>
      <RequestForm @gotoStep="gotoStep" :requestTemplateId="requestTemplateId" v-if="currentStep === 1"></RequestForm>
      <ApprovalForm @gotoStep="gotoStep" :requestTemplateId="requestTemplateId" v-if="currentStep === 2"></ApprovalForm>
      <TaskForm @gotoStep="gotoStep" :requestTemplateId="requestTemplateId" v-if="currentStep === 3"></TaskForm>
    </div>
  </div>
</template>

<script>
import ApprovalForm from './approval-form'
import RequestForm from './request-form'
import BasicInfo from './basic-info'
import TaskForm from './task-form'
import { submitTemplate } from '@/api/server.js'
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
    async submitTemplate () {
      this.$Modal.confirm({
        title: `${this.$t('submit_for_review')}`,
        content: `${this.$t('submit_for_review_tip')}`,
        'z-index': 1000000,
        okText: this.$t('confirm'),
        cancelText: this.$t('cancel'),
        onOk: async () => {
          let data = {
            requestTemplateId: this.requestTemplateId,
            status: 'created',
            targetStatus: 'pending',
            reason: '{}'
          }
          const { statusCode } = await submitTemplate(data)
          if (statusCode === 'OK') {
            this.$Notice.success({
              title: this.$t('successful'),
              desc: this.$t('successful')
            })
            this.$router.push({
              path: '/taskman/template-mgmt',
              query: {
                status: 'pending'
              }
            })
          }
        },
        onCancel: () => {}
      })
    },
    changeStep (val) {
      this.currentStep = val
    },
    backToTemplate () {
      this.$Modal.confirm({
        title: `${this.$t('back_to_template')}`,
        content: `${this.$t('back_to_template_tip')}`,
        'z-index': 1000000,
        okText: this.$t('confirm'),
        cancelText: this.$t('cancel'),
        onOk: async () => {
          this.$router.push({ path: '/taskman/template-mgmt' })
        },
        onCancel: () => {}
      })
    },
    gotoStep (tmpId, stepDirection) {
      console.log(77, stepDirection)
      this.requestTemplateId = tmpId
      if (stepDirection === 'forward') {
        this.currentStep++
      } else if (stepDirection === 'backward') {
        this.currentStep--
      }
    }
  },
  components: {
    ApprovalForm,
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
