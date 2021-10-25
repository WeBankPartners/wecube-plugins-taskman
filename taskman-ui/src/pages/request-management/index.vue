<template>
  <div style="margin: 24px">
    <Button @click="backToTemplate" icon="ios-undo-outline" style="margin-bottom: 8px">{{
      $t('back_to_template')
    }}</Button>
    <template v-if="currentStep === -1">
      <TemplateSelect @choiceTemp="choiceTemp"></TemplateSelect>
    </template>
    <template v-else>
      <Steps :current="currentStep">
        <!-- <Step icon="ios-add-circle">
          <span slot="title" @click="changeStep(0)">{{ $t('select_template') }}</span>
        </Step> -->
        <Step icon="md-apps">
          <span slot="title" @click="changeStep(0)">{{ $t('basic_information_settings') }}</span>
        </Step>
        <Step icon="md-cog">
          <span slot="title" @click="changeStep(1)">{{ $t('data_management') }}</span>
        </Step>
        <Step icon="ios-settings">
          <span slot="title" @click="changeStep(2)">{{ $t('data_binding') }}</span>
        </Step>
      </Steps>
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
import { confirmTemplate } from '@/api/server.js'
export default {
  name: '',
  data () {
    return {
      currentStep: -1,
      requestTemplate: '',
      procDefId: '',
      procDefKey: '',
      requestId: ''
    }
  },
  mounted () {
    if (this.$route.query.requestTemplate !== '') {
      this.requestTemplate = this.$route.query.requestTemplate
    }
    this.currentStep = -1
  },
  methods: {
    async confirmTemplate () {
      const { statusCode } = await confirmTemplate(this.requestTemplate)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
      }
    },
    changeStep (val) {
      if (this.requestId === '') {
        return
      }
      this.currentStep = val
    },
    backToTemplate () {
      this.$router.push({ path: '/' })
    },
    choiceTemp (data) {
      this.procDefId = data.procDefId
      this.procDefKey = data.procDefKey
      this.requestTemplate = data.id
      this.nextStep()
    },
    basicForm (data) {
      this.requestId = data.id
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
