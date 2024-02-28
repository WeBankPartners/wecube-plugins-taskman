<template>
  <div class="workbench-publish-detail">
    <Row class="back-header">
      <Icon size="22" type="md-arrow-back" class="icon" @click="$router.back()" />
      <span class="name">
        {{ `${detailInfo.name || ''}` }}
      </span>
    </Row>
    <Row class="w-header">
      <Col span="19">
        <!--请求进度-->
        <BaseProgress ref="progress"></BaseProgress>
      </Col>
      <Col span="5" class="btn-group">
        <!--撤回-->
        <Button v-if="jumpFrom === 'submit' && detailInfo.status === 'Pending'" type="error" @click="handleRecall">{{
          $t('tw_recall')
        }}</Button>
      </Col>
    </Row>
    <div style="display:flex;" class="content">
      <div style="width:100%;">
        <Form :model="form" label-position="right" :label-width="120">
          <template>
            <!--请求信息-->
            <HeaderTitle :title="$t('tw_request_title')">
              <Row :gutter="20">
                <Col :span="3">{{ $t('request_id') }}：</Col>
                <Col :span="9">{{ detailInfo.id || '--' }}</Col>
                <Col :span="3">{{ $t('tw_request_type') }}：</Col>
                <Col :span="9">{{ { 0: $t('tw_request'), 1: $t('tw_publish') }[detailInfo.requestType] || '--' }}</Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="3">{{ $t('createdTime') }}：</Col>
                <Col :span="9">{{ detailInfo.createdTime || '--' }}</Col>
                <Col :span="3">{{ $t('expected_completion_time') }}：</Col>
                <Col :span="9">{{ detailInfo.expectTime || '--' }}</Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="3">{{ $t('tw_request_progress') }}：</Col>
                <Col :span="9">
                  <Progress :percent="detailInfo.progress" style="width:150px;" />
                </Col>
                <Col :span="3">{{ $t('tw_request_status') }}：</Col>
                <Col :span="9">{{ getStatusName(detailInfo.status) || '--' }}</Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <!--当前节点-->
                <Col :span="3">{{ $t('tw_cur_tag') }}：</Col>
                <Col :span="9">{{
                  {
                    waitCommit: $t('tw_wait_commit'),
                    sendRequest: $t('tw_commit_request'),
                    requestPending: $t('tw_request_pending'),
                    requestComplete: $t('tw_request_complete')
                  }[detailInfo.curNode] ||
                    detailInfo.curNode ||
                    '--'
                }}</Col>
                <Col :span="3"
                  >{{ detailInfo.status === 'Draft' ? $t('tw_pending_handler') : $t('tw_cur_handler') }}：</Col
                >
                <Col :span="9">{{ detailInfo.handler || '--' }}</Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="3">{{ $t('createdBy') }}：</Col>
                <Col :span="9">{{ detailInfo.createdBy || '--' }}</Col>
                <Col :span="3">{{ $t('tw_creatby_role') }}：</Col>
                <Col :span="9">{{ detailInfo.role || '--' }}</Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="3">{{ $t('tw_use_template') }}：</Col>
                <Col :span="9"
                  >{{ detailInfo.templateName }}<Tag>{{ detailInfo.version }}</Tag></Col
                >
                <Col :span="3">{{ $t('tm_template_group') }}：</Col>
                <Col :span="9">{{ detailInfo.templateGroupName || '--' }}</Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="3">{{ $t('tw_request_des') }}：</Col>
                <Col :span="9">{{ detailInfo.description || '--' }}</Col>
              </Row>
            </HeaderTitle>
            <!--处理历史-->
            <HeaderTitle :title="$t('tw_handle_history')">
              <Collapse v-model="activeStep" simple style="margin-top:30px;">
                <Panel v-for="(data, index) in historyData" :name="index + ''" :key="index">
                  <!--提交请求-->
                  <template v-if="index === 0">
                    <div style="display:flex;align-items:center;width:100%;">
                      <div style="width:70px;font-weight:bold;line-height:22px;">{{ $t('tw_submit_request') }}</div>
                      <div style="width:calc(100% - 70px)">
                        <HeaderTag :showHeader="true" :data="data" :operation="$t('tw_submit_approval')"></HeaderTag>
                      </div>
                    </div>
                    <div slot="content" class="history">
                      <FormItem :label="$t('tw_choose_object')" required>
                        <Select v-model="form.rootEntityId" :disabled="true" clearable filterable style="width:300px;">
                          <Option v-for="item in rootEntityOptions" :value="item.guid" :key="item.guid">{{
                            item.key_name
                          }}</Option>
                        </Select>
                      </FormItem>
                      <EntityTable :data="data.formData" :requestId="requestId" :formDisable="true"></EntityTable>
                    </div>
                  </template>
                  <!--请求定版-->
                  <template v-else-if="index === 1">
                    <div style="display:flex;align-items:center;width:100%;">
                      <div style="width:70px;font-weight:bold;line-height:22px;">{{ $t('tw_request_pending') }}</div>
                      <div style="width:calc(100% - 70px)">
                        <HeaderTag :data="data" :operation="$t('final_version')"></HeaderTag>
                      </div>
                    </div>
                    <div slot="content" class="history">
                      <DataBind
                        :isHandle="isHandle"
                        :requestTemplate="requestTemplate"
                        :requestId="requestId"
                        :formDisable="true"
                        :showBtn="false"
                      ></DataBind>
                    </div>
                  </template>
                  <!--任务审批-->
                  <template v-else>
                    <div style="display:flex;align-items:center;width:100%;">
                      <div style="width:70px;font-weight:bold;line-height:22px;">{{ data.taskName }}</div>
                      <div style="width:calc(100% - 70px)">
                        <HeaderTag :data="data"></HeaderTag>
                      </div>
                    </div>
                    <div slot="content" class="history">
                      <EntityTable
                        :data="data.formData"
                        :requestId="requestId"
                        :formDisable="!data.editable || enforceDisable"
                      ></EntityTable>
                      <div>
                        <Form :label-width="80" style="margin: 16px 0">
                          <FormItem v-if="data.requestId === ''" :label="$t('task') + $t('description')">
                            <Input disabled v-model="data.description" type="textarea" />
                          </FormItem>
                          <FormItem
                            :label="$t('process_result')"
                            v-if="data.nextOption && data.nextOption.length !== 0"
                          >
                            <span slot="label">
                              {{ $t('process_result') }}
                              <span style="color: #ed4014"> * </span>
                            </span>
                            <Select v-model="data.choseOption" :disabled="!data.editable || enforceDisable">
                              <Option v-for="option in data.nextOption" :value="option" :key="option">{{
                                option
                              }}</Option>
                            </Select>
                          </FormItem>
                          <FormItem :label="$t('process_comments')">
                            <Input
                              :disabled="!data.editable || enforceDisable"
                              v-model="data.comment"
                              type="textarea"
                            />
                          </FormItem>
                        </Form>
                      </div>
                    </div>
                  </template>
                </Panel>
              </Collapse>
            </HeaderTitle>
            <!--当前处理-任务审批-->
            <HeaderTitle
              v-if="isHandle && detailInfo.status === 'InProgress'"
              :title="$t('tw_cur_handle')"
              :subTitle="handleData.taskName"
            >
              <Steps :current="1" direction="vertical">
                <Step>
                  <div slot="title" class="task-step">
                    <div>{{ $t('tw_approval_step1') }}</div>
                    <div>{{ $t('tw_approval_step1_tips') }}</div>
                  </div>
                  <div slot="content" style="padding:20px 0px;">
                    <EntityTable
                      ref="entityTable"
                      :data="handleData.formData"
                      :requestId="requestId"
                      :formDisable="!handleData.editable || enforceDisable"
                      :isAddRow="true"
                    ></EntityTable>
                  </div>
                </Step>
                <Step>
                  <div slot="title" class="task-step">
                    <div>{{ $t('tw_approval_step2') }}</div>
                    <div>{{ $t('tw_approval_step2_tips') }}</div>
                  </div>
                  <div slot="content" style="padding:20px 0px;">
                    <Form :label-width="80" style="margin: 16px 0">
                      <FormItem v-if="handleData.requestId === ''" :label="$t('task') + $t('description')">
                        <Input disabled v-model="handleData.description" type="textarea" />
                      </FormItem>
                      <FormItem
                        :label="$t('process_result')"
                        v-if="handleData.nextOption && handleData.nextOption.length !== 0"
                      >
                        <span slot="label">
                          {{ $t('process_result') }}
                          <span style="color: #ed4014"> * </span>
                        </span>
                        <Select v-model="handleData.choseOption" :disabled="!handleData.editable || enforceDisable">
                          <Option v-for="option in handleData.nextOption" :value="option" :key="option">{{
                            option
                          }}</Option>
                        </Select>
                      </FormItem>
                      <FormItem :label="$t('process_comments')">
                        <Input
                          :disabled="!handleData.editable || enforceDisable"
                          v-model="handleData.comment"
                          type="textarea"
                        />
                      </FormItem>
                      <FormItem :label="$t('tw_attach')">
                        <UploadFile
                          :id="handleData.taskId"
                          :files="handleData.attachFiles"
                          :formDisable="enforceDisable"
                          type="task"
                        ></UploadFile>
                      </FormItem>
                    </Form>
                    <div style="text-align: center">
                      <Button v-if="handleData.editable" :disabled="enforceDisable" @click="saveTaskData" type="info">{{
                        $t('save')
                      }}</Button>
                      <Button
                        v-if="handleData.editable"
                        :disabled="
                          enforceDisable ||
                            (handleData.nextOption &&
                              handleData.nextOption.length !== 0 &&
                              handleData.choseOption === '')
                        "
                        @click="commitTaskData"
                        type="primary"
                        >{{ $t('tw_commit') }}</Button
                      >
                    </div>
                  </div>
                </Step>
              </Steps>
            </HeaderTitle>
            <!--当前处理-请求定版-->
            <HeaderTitle
              v-if="isHandle && detailInfo.status === 'Pending'"
              :title="$t('tw_cur_handle')"
              :subTitle="$t('tw_request_pending')"
            >
              <Steps :current="1" direction="vertical">
                <Step>
                  <div slot="title" class="task-step">
                    <div>{{ $t('tw_pending_step1') }}</div>
                    <div>{{ $t('tw_pending_step1_tips') }}</div>
                  </div>
                  <div slot="content" style="padding:20px 0px;">
                    <FormItem :label="$t('tw_choose_object')" required>
                      <Select v-model="form.rootEntityId" :disabled="true" clearable filterable style="width:300px;">
                        <Option v-for="item in rootEntityOptions" :value="item.guid" :key="item.guid">{{
                          item.key_name
                        }}</Option>
                      </Select>
                    </FormItem>
                    <EntityTable
                      :data="historyData[0].formData"
                      :requestId="requestId"
                      :formDisable="true"
                    ></EntityTable>
                  </div>
                </Step>
                <Step>
                  <div slot="title" class="task-step">
                    <div>{{ $t('tw_pending_step2') }}</div>
                    <div>{{ $t('tw_pending_step2_tips') }}</div>
                  </div>
                  <div slot="content" style="padding:20px 0px;">
                    <DataBind
                      :isHandle="isHandle"
                      :requestTemplate="requestTemplate"
                      :requestId="requestId"
                      :formDisable="formDisable || detailInfo.status !== 'Pending'"
                      :actionName="actionName"
                    ></DataBind>
                  </div>
                </Step>
              </Steps>
            </HeaderTitle>
          </template>
        </Form>
      </div>
    </div>
    <!--编排流程-->
    <div class="expand-btn" :style="{ right: flowVisible ? '440px' : '0px' }" @click="flowVisible = !flowVisible">
      <Icon v-if="flowVisible" type="ios-arrow-dropright-circle" :size="28" />
      <Icon v-else type="ios-arrow-dropleft-circle" :size="28" />
    </div>
    <div v-show="flowVisible" class="flow-expand">
      <StaticFlow
        v-if="['Draft', 'Pending'].includes(detailInfo.status)"
        :requestTemplate="requestTemplate"
      ></StaticFlow>
      <DynamicFlow v-else :flowId="detailInfo.procInstanceId"></DynamicFlow>
    </div>
  </div>
</template>

<script>
import HeaderTitle from '../../components/header-title.vue'
import HeaderTag from '../../components/header-tag.vue'
import StaticFlow from '../../components/flow/static-flow.vue'
import DynamicFlow from '../../components/flow/dynamic-flow.vue'
import EntityTable from '../../components/entity-table.vue'
import DataBind from '../../components/data-bind.vue'
import UploadFile from '../../components/upload.vue'
import BaseProgress from '../../components/base-progress.vue'
import { deepClone } from '@/pages/util/index'
import {
  getRootEntity,
  getEntityData,
  getPublishInfo,
  getRequestInfo,
  saveTaskData,
  commitTaskData,
  recallRequest
} from '@/api/server'
export default {
  components: {
    HeaderTitle,
    HeaderTag,
    StaticFlow,
    DynamicFlow,
    EntityTable,
    DataBind,
    UploadFile,
    BaseProgress
  },
  props: {
    // 1发布,2请求(3问题,4事件,5变更)
    actionName: {
      type: String,
      default: '1'
    }
  },
  data () {
    return {
      enforceDisable: this.$route.query.enforceDisable === 'Y',
      isHandle: this.$route.query.isHandle === 'Y', // 处理标志(用于请求定版)
      formDisable: this.$route.query.isCheck === 'Y', // 查看标志
      jumpFrom: this.$route.query.jumpFrom, // 入口tab标记
      requestTemplate: this.$route.query.requestTemplate,
      requestId: this.$route.query.requestId,
      // procDefId: '',
      // procDefKey: '',
      detailInfo: {}, // 详情信息
      form: {
        name: '',
        description: '',
        expectTime: '',
        rootEntityId: '', // 目标对象
        data: []
      },
      rootEntityOptions: [],
      requestData: [], // 发布目标对象表格数据
      historyData: [], // 处理历史数据
      handleData: {},
      attachFiles: [], // 上传附件
      activeStep: '', // 处理历史当前展开
      errorNode: '',
      flowVisible: false
    }
  },
  computed: {
    getStatusName () {
      return function (val) {
        const list = [
          { label: this.$t('status_pending'), value: 'Pending', color: '#b886f8' },
          { label: this.$t('status_inProgress'), value: 'InProgress', color: '#1990ff' },
          { label: this.$t('status_inProgress_faulted'), value: 'InProgress(Faulted)', color: '#f26161' },
          { label: this.$t('status_termination'), value: 'Termination', color: '#e29836' },
          { label: this.$t('status_complete'), value: 'Completed', color: '#7ac756' },
          { label: this.$t('status_inProgress_timeouted'), value: 'InProgress(Timeouted)', color: '#f26161' },
          { label: this.$t('status_faulted'), value: 'Faulted', color: '#e29836' },
          { label: this.$t('status_draft'), value: 'Draft', color: '#808695' }
        ]
        const item = list.find(i => i.value === val) || {}
        return item.label
      }
    }
  },
  async created () {
    // 获取请求进度
    this.$nextTick(() => {
      this.$refs.progress.initData(this.requestTemplate, this.requestId)
    })
    // 获取详情信息
    this.getPublishInfo()
    // 获取附件和操作单元
    this.getRequestInfo()
    this.getEntity()
  },
  methods: {
    // 获取附件和操作单元
    async getRequestInfo () {
      const { statusCode, data } = await getRequestInfo(this.requestId)
      if (statusCode === 'OK') {
        this.attachFiles = data.attachFiles
        this.form.rootEntityId = data.cache
      }
    },
    // 获取详情数据
    async getPublishInfo () {
      const { statusCode, data } = await getPublishInfo(this.requestId)
      if (statusCode === 'OK') {
        this.detailInfo = data.request || {}
        const { name, description, expectTime } = this.detailInfo
        this.form = Object.assign({}, this.form, {
          name,
          description,
          expectTime
        })
        this.historyData = data.data || []
        // 获取请求定版-当前处理数据
        if (this.detailInfo.status === 'Pending') {
          if (this.historyData && this.historyData.length > 1) {
            this.handleData = this.historyData[1]
            this.historyData.splice(1, 1)
          }
          // 任务处理-当前处理数据
        } else if (this.detailInfo.status === 'InProgress') {
          const index = this.historyData.findIndex(item => item.editable)
          this.handleData = this.historyData[index]
          this.historyData.splice(index, 1)
        }
        if (this.isHandle === false) {
          this.activeStep = this.historyData.length - 1 + ''
        }
      }
    },
    // 操作目标对象
    async getEntity () {
      let params = {
        params: {
          requestId: this.requestId
        }
      }
      const { statusCode, data } = await getRootEntity(params)
      if (statusCode === 'OK') {
        this.rootEntityOptions = data.data || []
        // this.form.rootEntityId = this.rootEntityOptions[0] && this.rootEntityOptions[0].guid
      }
    },
    // 获取目标对象对应表单配置
    async getEntityData () {
      let params = {
        params: {
          requestId: this.requestId,
          rootEntityId: this.form.rootEntityId
        }
      }
      const { statusCode, data } = await getEntityData(params)
      if (statusCode === 'OK') {
        const requestData = data.data || []
        this.requestData = requestData
      }
    },
    // 校验表格数据必填项
    requiredCheck (data) {
      let tabIndex = ''
      let result = true
      data.forEach((requestData, index) => {
        let requiredName = []
        requestData.title.forEach(t => {
          if (t.required === 'yes') {
            requiredName.push(t.name)
          }
        })
        requestData.value.forEach(v => {
          requiredName.forEach(key => {
            let val = v.entityData[key]
            if (Array.isArray(val)) {
              if (val.length === 0) {
                result = false
                if (tabIndex === '') {
                  tabIndex = index
                }
              }
            } else {
              if (val === '' || val === undefined) {
                result = false
                if (tabIndex === '') {
                  tabIndex = index
                }
              }
            }
          })
        })
      })
      this.$refs.entityTable.validTable(tabIndex)
      return result
    },
    noChooseCheck (data) {
      let tabIndex = ''
      let result = true
      data.forEach((requestData, index) => {
        if (requestData.value && requestData.value.length === 0) {
          tabIndex = index
          result = false
        }
      })
      this.$refs.entityTable.validTable(tabIndex)
      return result
    },
    // 任务审批保存
    async saveTaskData () {
      // 提取表格勾选的数据
      const requestData = deepClone(this.$refs.entityTable && this.$refs.entityTable.requestData)
      this.handleData.formData =
        requestData.map(item => {
          let refKeys = []
          item.title.forEach(t => {
            if (t.elementType === 'select' || t.elementType === 'wecmdbEntity') {
              refKeys.push(t.name)
            }
          })
          if (Array.isArray(item.value)) {
            item.value = item.value.filter(j => {
              return j.entityData._checked
            })
            // 删除多余的属性
            item.value.forEach(v => {
              delete v.entityData._checked
              refKeys.forEach(ref => {
                delete v.entityData[ref + 'Options']
              })
            })
          }
          return item
        }) || []
      // 必填项校验提示
      if (!this.requiredCheck(this.handleData.formData)) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('required_tip')
        })
      }
      // 表格至少勾选一条数据校验
      const tabName = this.$refs.entityTable.activeTab
      if (!this.noChooseCheck(this.handleData.formData)) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: `【${tabName}】${this.$t('tw_table_noChoose_tips')}`
        })
      }
      const { statusCode } = await saveTaskData(this.handleData.taskId, this.handleData)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
      }
    },
    // 任务审批提交
    async commitTaskData () {
      // 提取表格勾选的数据
      const requestData = deepClone(this.$refs.entityTable && this.$refs.entityTable.requestData)
      this.handleData.formData =
        requestData.map(item => {
          let refKeys = []
          item.title.forEach(t => {
            if (t.elementType === 'select' || t.elementType === 'wecmdbEntity') {
              refKeys.push(t.name)
            }
          })
          if (Array.isArray(item.value)) {
            item.value = item.value.filter(j => {
              return j.entityData._checked
            })
            // 删除多余的属性
            item.value.forEach(v => {
              delete v.entityData._checked
              refKeys.forEach(ref => {
                delete v.entityData[ref + 'Options']
              })
            })
          }
          return item
        }) || []
      // 必填项校验提示
      if (!this.requiredCheck(this.handleData.formData)) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('required_tip')
        })
      }
      // 表格至少勾选一条数据校验
      const tabName = this.$refs.entityTable.activeTab
      if (!this.noChooseCheck(this.handleData.formData)) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: `【${tabName}】${this.$t('tw_table_noChoose_tips')}`
        })
      }
      const { statusCode } = await commitTaskData(this.handleData.taskId, this.handleData)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.$router.push({ path: `/taskman/workbench?tabName=hasProcessed&actionName=${this.actionName}&type=2` })
      }
    },
    // 撤回
    async handleRecall () {
      this.$Modal.confirm({
        title: this.$t('tw_confirm') + this.$t('tw_recall'),
        content: this.$t('tw_recall_tips'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const { statusCode } = await recallRequest(this.requestId)
          if (statusCode === 'OK') {
            this.$Notice.success({
              title: this.$t('successful'),
              desc: this.$t('successful')
            })
            this.$router.push({ path: `/taskman/workbench?tabName=submit&actionName=${this.actionName}` })
          }
        },
        onCancel: () => {}
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.workbench-publish-detail {
  .back-header {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    .icon {
      cursor: pointer;
      width: 28px;
      height: 24px;
      color: #fff;
      border-radius: 2px;
      background: #2d8cf0;
    }
    .name {
      font-size: 16px;
      margin-left: 16px;
      display: flex;
      align-items: center;
    }
  }
  .w-header {
    padding-bottom: 15px;
    margin-bottom: 20px;
    border-bottom: 2px dashed #d7dadc;
    .btn-group {
      display: flex;
      justify-content: flex-end;
    }
  }
  .content {
    min-height: 500px;
    .history {
      padding: 20px;
      border: 1px dashed #d7dadc;
      margin-top: 20px;
    }
    .task-step {
      display: flex;
      div:first-child {
        color: #515a6e;
      }
      div:last-child {
        font-weight: 400;
        font-size: 12px;
        color: #515a6e;
      }
    }
  }
  .expand-btn {
    position: fixed;
    top: calc(50% - 14px);
    cursor: pointer;
  }
  .flow-expand {
    height: 100%;
    position: fixed;
    right: 0px;
    top: 0px;
    background: #f8f8f9;
    padding: 20px;
  }
}
</style>
<style lang="scss">
.workbench-publish-detail {
  .ivu-steps-content {
    padding-left: 0px !important;
    padding-top: 5px;
    font-size: 12px;
    color: #3d3c38 !important;
  }
  .steps .ivu-steps .ivu-steps-tail > i {
    height: 3px;
    background: #8189a5;
  }
  .ivu-form-item {
    margin-bottom: 15px !important;
  }
  .ivu-form-item-content {
    line-height: 20px;
  }
  .required {
    color: red;
  }
  .ivu-collapse {
    border: none !important;
  }
  .ivu-collapse > .ivu-collapse-item {
    border-top: none;
  }
  .ivu-collapse-header {
    height: auto !important;
    display: flex;
    align-items: center;
    padding-left: 0px !important;
  }
  .ivu-collapse-content {
    padding: 0 0 0 16px;
  }
}
</style>
