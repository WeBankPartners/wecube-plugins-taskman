<template>
  <div style="margin: 24px">
    <Button @click="backToTemplate" icon="ios-undo-outline" style="margin-bottom: 8px">{{
      $t('back_to_template')
    }}</Button>

    <template v-if="currentStep === -1">
      <TemplateSelect @choiceTemp="choiceTemp"></TemplateSelect>
    </template>
    <template v-else>
      <div style="width: 84%;margin: 0 auto;">
        <Steps :current="currentStep">
          <Step icon="md-apps">
            <span slot="title" @click="changeStep(0)">{{ $t('basic_information_settings') }}</span>
          </Step>
          <Step icon="md-cog">
            <span slot="title" @click="changeStep(1)">{{ $t('form_settings') }}</span>
          </Step>
          <Step icon="ios-settings" v-if="['group_handle'].includes(jumpFrom)">
            <span slot="title" @click="changeStep(2)">{{ $t('data_binding') }}</span>
          </Step>
        </Steps>
      </div>
      <div v-if="currentStep !== -1" style="margin-top:48px;">
        <BasicForm @basicForm="basicForm" v-if="currentStep === 0"></BasicForm>
        <DataCrud @nextStep="nextStep" v-if="currentStep === 1"></DataCrud>
        <DataBind v-if="currentStep === 2"></DataBind>
      </div>
    </template>
  </div>
</template>

<script>
import TemplateSelect from './template-select'
import BasicForm from './basic-form'
import DataCrud from './data-crud'
import DataBind from './data-bind'
export default {
  name: '',
  data () {
    return {
      currentStep: -1,
      isAdd: true,
      isHandle: false, // 处理标志
      formDisable: false, // 查看标志
      jumpFrom: '', // 入口tab标记
      requestTemplate: '',
      procDefId: '',
      procDefKey: '',
      requestId: ''
    }
  },
  mounted () {
    this.jumpFrom = this.$route.query.jumpFrom
    this.requestTemplate = this.$route.query.requestTemplate
    this.requestId = this.$route.query.requestId
    this.isAdd = this.$route.query.isAdd === 'Y'
    this.isHandle = this.$route.query.isHandle === 'Y'
    this.formDisable = this.$route.query.isCheck === 'Y'
    this.currentStep = this.isAdd ? -1 : 0
  },
  methods: {
    changeStep (val) {
      if (this.requestId === '') {
        return
      }
      this.currentStep = val
    },
    backToTemplate () {
      this.$router.push({ path: '/taskman/request-mgmt', query: { activeTab: this.jumpFrom } })
    },
    choiceTemp (data) {
      this.procDefId = data.procDefId
      this.procDefKey = data.procDefKey
      this.requestTemplate = data.id
      this.nextStep()
    },
    basicForm (id) {
      this.requestId = id
      this.nextStep()
    },
    nextStep () {
      this.currentStep++
    }
  },
  components: {
    TemplateSelect,
    BasicForm,
    DataCrud,
    DataBind
  }
}
</script>

<style scoped lang="scss">
.header-icon {
  float: right;
  margin: 3px 40px 0 0 !important;
}
</style>
