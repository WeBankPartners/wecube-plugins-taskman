<template>
  <div class="template-management-index">
    <div class="header-step">
      <Button @click="backToTemplate" icon="ios-undo-outline" class="button">{{ $t('back_to_template') }}</Button>
      <Steps :current="currentStep" class="step">
        <Step :title="$t('basic_information_settings')"></Step>
        <Step :title="$t('form')"></Step>
        <Step :title="$t('tw_approval')"></Step>
        <Step :title="$t('task')"></Step>
      </Steps>
    </div>
    <div></div>
    <div v-if="currentStep !== -1" style="margin-top:10px;">
      <BasicInfo
        @gotoStep="gotoStep"
        :isCheck="isCheck"
        :requestTemplateId="requestTemplateId"
        v-if="currentStep === 0"
      ></BasicInfo>
      <RequestForm
        @gotoStep="gotoStep"
        :isCheck="isCheck"
        :requestTemplateId="requestTemplateId"
        v-if="currentStep === 1"
      ></RequestForm>
      <ApprovalForm
        @gotoStep="gotoStep"
        :isCheck="isCheck"
        :requestTemplateId="requestTemplateId"
        v-if="currentStep === 2"
      ></ApprovalForm>
      <TaskForm
        @gotoStep="gotoStep"
        :isCheck="isCheck"
        :requestTemplateId="requestTemplateId"
        v-if="currentStep === 3"
      ></TaskForm>
    </div>
  </div>
</template>

<script>
import ApprovalForm from './approval-form'
import RequestForm from './request-form'
import BasicInfo from './basic-info'
import TaskForm from './task-form'
export default {
  name: '',
  data () {
    return {
      currentStep: -1,
      requestTemplateId: '',
      isCheck: 'N',
      parentStatus: ''
    }
  },
  mounted () {
    if (this.$route.query.requestTemplateId !== '') {
      this.requestTemplateId = this.$route.query.requestTemplateId
    }
    this.isCheck = this.$route.query.isCheck
    this.parentStatus = this.$route.query.parentStatus
    this.currentStep = 0
  },
  methods: {
    changeStep (val) {
      this.currentStep = val
    },
    backToTemplate () {
      if (this.isCheck === 'Y') {
        this.$router.push({
          path: '/template-mgmt',
          query: {
            status: this.parentStatus,
            needCache: 'yes'
          }
        })
      } else {
        this.$Modal.confirm({
          title: `${this.$t('back_to_template')}`,
          content: `${this.$t('back_to_template_tip')}`,
          'z-index': 1000000,
          okText: this.$t('confirm'),
          cancelText: this.$t('cancel'),
          onOk: async () => {
            this.$router.push({
              path: '/template-mgmt',
              query: {
                status: 'created',
                needCache: 'yes'
              }
            })
          },
          onCancel: () => {}
        })
      }
    },
    gotoStep (tmpId, stepDirection) {
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

<style lang="scss">
.template-management-index {
  .header-step {
    display: flex;
    align-items: center;
    .button {
      width: fit-content;
    }
    .step {
      flex: 1;
      padding: 0 200px;
    }
  }
  .header-icon {
    float: right;
    margin: 3px 40px 0 0 !important;
  }
}
</style>
