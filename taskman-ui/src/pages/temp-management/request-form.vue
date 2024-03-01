<template>
  <Row type="flex">
    <Col span="24" style="padding: 0 20px">
      <div style="border-bottom: 1px solid #dcdee2">
        <div @click="changTab('msgForm')" :class="activeTab === 'msgForm' ? 'tab-active' : 'tab'">
          {{ this.$t('1.信息表单') }}
        </div>
        <div @click="changTab('dataForm')" :class="activeTab === 'dataForm' ? 'tab-active' : 'tab'">
          {{ this.$t('2.数据表单') }}
        </div>
      </div>
      <div style="margin-top: 16px;">
        <RequestFormMsg
          v-if="activeTab === 'msgForm'"
          @gotoNextStep="gotoNextStep"
          @changTab="changTab"
          ref="msgFormRef"
        ></RequestFormMsg>
        <RequestFormData
          v-if="activeTab === 'dataForm'"
          @gotoNextStep="gotoNextStep"
          @changTab="changTab"
          ref="dataFormRef"
          :requestTemplateId="requestTemplateId"
        ></RequestFormData>
      </div>
      <!-- <Tabs :value="activeTab" @on-click="changTab">
        <TabPane label="1.信息表单" name="msgForm">
          <RequestFormMsg @gotoNextStep="gotoNextStep" ref="msgFormRef"></RequestFormMsg>
        </TabPane>
        <TabPane label="2.数据表单" name="dataForm" disabled>
          <RequestFormData
            @gotoNextStep="gotoNextStep"
            ref="dataFormRef"
            :requestTemplateId="requestTemplateId"
          ></RequestFormData>
        </TabPane>
      </Tabs> -->
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
    this.changTab(this.activeTab, true)
  },
  methods: {
    changTab (tabName, firstLoad = false) {
      if (!firstLoad && tabName === this.activeTab) {
        return
      }

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

      this.$nextTick(() => {})
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
  color: #2d8cf0;
  border-bottom: 2px solid #2d8cf0;
}
</style>
