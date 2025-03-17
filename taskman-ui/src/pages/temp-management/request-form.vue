<template>
  <Row type="flex">
    <Col span="24">
      <div>
        <div @click="changTab('msgForm')" :class="activeTab === 'msgForm' ? 'tab-active' : 'tab'">
          1.{{ this.$t('tw_information_form') }}
        </div>
        <div @click="changTab('dataForm')" :class="activeTab === 'dataForm' ? 'tab-active' : 'tab'">
          2.{{ this.$t('tw_data_form') }}
        </div>
      </div>
      <div style="margin-top:8px;">
        <RequestFormMsg
          v-if="activeTab === 'msgForm'"
          @gotoStep="gotoStep"
          @changTab="changTab"
          :isCheck="isCheck"
          ref="msgFormRef"
        ></RequestFormMsg>
        <RequestFormData
          v-if="activeTab === 'dataForm'"
          @gotoStep="gotoStep"
          @changTab="changTab"
          :isCheck="isCheck"
          ref="dataFormRef"
          :requestTemplateId="requestTemplateId"
        ></RequestFormData>
      </div>
    </Col>
  </Row>
</template>

<script>
import RequestFormMsg from './request-form-msg.vue'
import RequestFormData from './request-form-data.vue'
export default {
  name: 'request-form',
  data () {
    return {
      activeTab: 'msgForm'
    }
  },
  props: ['isCheck', 'requestTemplateId'],
  mounted () {
    this.changTab(this.activeTab, true)
  },
  methods: {
    changTab (tabName, firstLoad = false) {
      if (!firstLoad && tabName === this.activeTab) {
        return
      }

      if (this.isCheck === 'Y') {
        this.activeTab = tabName
        this.$nextTick(() => {
          this.$refs[`${this.activeTab}Ref`].loadPage(this.requestTemplateId)
        })
      } else {
        const tabStatus = this.$refs[`${this.activeTab}Ref`].panalStatus()
        if (tabStatus) {
          this.$nextTick(() => {
            this.$refs[`${this.activeTab}Ref`].tabChange()
          })
        } else {
          this.activeTab = tabName
          this.$nextTick(() => {
            this.$refs[`${this.activeTab}Ref`].loadPage(this.requestTemplateId)
          })
        }
      }
    },
    gotoStep (requestTemplateId, stepDirection) {
      this.$emit('gotoStep', requestTemplateId, stepDirection)
    }
  },
  components: {
    RequestFormMsg,
    RequestFormData
  }
}
</script>
<style scoped lang="scss">
.tab {
  display: inline-block;
  height: 100%;
  padding: 8px 16px;
  margin-right: 16px;
  box-sizing: border-box;
  cursor: pointer;
  text-decoration: none;
}
.tab-active {
  @extend .tab;
  color: #5384ff;
  border-bottom: 2px solid #5384ff;
}
</style>
