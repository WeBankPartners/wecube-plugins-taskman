<template>
  <Row type="flex">
    <Col span="20" offset="2">
      <Tabs :value="activeTab" @on-click="changTab">
        <TabPane label="1.信息表单" name="msgForm">
          <RequestFormMsg @gotoNextStep="gotoNextStep" ref="msgFormRef"></RequestFormMsg>
        </TabPane>
        <TabPane label="2.数据表单" name="dataForm">
          <RequestFormData
            @gotoNextStep="gotoNextStep"
            ref="dataFormRef"
            :requestTemplateId="requestTemplateId"
          ></RequestFormData>
        </TabPane>
      </Tabs>
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
  props: ['requestTemplateId'],
  mounted () {
    this.changTab(this.activeTab)
  },
  methods: {
    changTab (tabName) {
      this.activeTab = tabName
      this.$refs[`${this.activeTab}Ref`].loadPage(this.requestTemplateId)
    },
    gotoNextStep () {
      this.$emit('gotoNextStep', this.requestTemplateId)
    }
  },
  components: {
    RequestFormMsg,
    RequestFormData
  }
}
</script>
<style>
.ivu-input[disabled],
fieldset[disabled] .ivu-input {
  color: #757575 !important;
}
.ivu-select-input[disabled] {
  color: #757575 !important;
  -webkit-text-fill-color: #757575 !important;
}
.ivu-select-disabled .ivu-select-selection {
  color: #757575 !important;
}
</style>
<style scoped lang="scss">
.active-zone {
  color: #338cf0;
}
.ivu-form-item {
  margin-bottom: 8px;
}
.list-group-item- {
  display: inline-block;
  margin: 2px 0;
}
</style>
