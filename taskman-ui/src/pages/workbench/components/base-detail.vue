<template>
  <div class="workbench-publish-detail">
    <Row class="back-header">
      <Icon size="22" type="md-arrow-back" class="icon" @click="handleToHome" />
      <span class="name">
        {{ `${detail.name || ''}` }}
      </span>
    </Row>
    <Row class="w-header">
      <Col span="24">
        <!--请求进度-->
        <BaseProgress ref="progress" :status="detail.status"></BaseProgress>
      </Col>
    </Row>
    <div class="content">
      <Form :model="form" label-position="right" :label-width="120" style="width:100%;">
        <!--基础信息-->
        <HeaderTitle :title="$t('tw_request_title')">
          <div style="padding-left:16px;">
            <Row :gutter="20">
              <Col :span="12" class="info-item">
                <div class="info-item-label">{{ $t('request_id') }}：</div>
                <div class="info-item-value">{{ detail.id || '-' }}</div>
              </Col>
              <Col :span="12" class="info-item">
                <div class="info-item-label">{{ $t('tw_request_type') }}：</div>
                <div class="info-item-value">
                  {{ typeMap[detail.requestType] || '-' }}
                </div>
              </Col>
            </Row>
            <Row style="margin-top:10px;" :gutter="20">
              <Col :span="12" class="info-item">
                <div class="info-item-label">{{ $t('createdTime') }}：</div>
                <div class="info-item-value">{{ detail.createdTime || '-' }}</div>
              </Col>
              <Col :span="12" class="info-item">
                <div class="info-item-label">{{ $t('expected_completion_time') }}：</div>
                <div class="info-item-value">{{ detail.expectTime || '-' }}</div>
              </Col>
            </Row>
            <Row style="margin-top:10px;" :gutter="20">
              <Col :span="12" class="info-item">
                <div class="info-item-label">{{ $t('tw_request_progress') }}：</div>
                <div class="info-item-value">
                  <Progress :percent="detail.progress" style="width:150px;" />
                </div>
              </Col>
              <Col :span="12" class="info-item">
                <div class="info-item-label">{{ $t('tw_request_status') }}：</div>
                <div class="info-item-value">{{ getStatusName(detail.status) || '-' }}</div>
              </Col>
            </Row>
            <Row style="margin-top:10px;" :gutter="20">
              <!--当前节点-->
              <Col :span="12" class="info-item">
                <div class="info-item-label">{{ $t('tw_cur_tag') }}：</div>
                <div class="info-item-value">
                  {{
                    {
                      waitCommit: $t('tw_wait_commit'),
                      sendRequest: $t('tw_commit_request'),
                      requestPending: $t('tw_request_pending'),
                      requestComplete: $t('tw_request_complete')
                    }[detail.curNode] ||
                      detail.curNode ||
                      '-'
                  }}
                </div>
              </Col>
              <Col :span="12" class="info-item">
                <div class="info-item-label">
                  {{ detail.status === 'Draft' ? $t('tw_pending_handler') : $t('tw_cur_handler') }}：
                </div>
                <div class="info-item-value">{{ formatHandler || '-' }}</div>
              </Col>
            </Row>
            <Row style="margin-top:10px;" :gutter="20">
              <Col :span="12" class="info-item">
                <div class="info-item-label">{{ $t('createdBy') }}：</div>
                <div class="info-item-value">{{ detail.createdBy || '-' }}</div>
              </Col>
              <Col :span="12" class="info-item">
                <div class="info-item-label">{{ $t('tw_creatby_role') }}：</div>
                <div class="info-item-value">{{ detail.roleDisplay || '-' }}</div>
              </Col>
            </Row>
            <Row style="margin-top:10px;" :gutter="20">
              <Col :span="12" class="info-item">
                <div class="info-item-label">{{ $t('tw_use_template') }}：</div>
                <div class="info-item-value">
                  {{ detail.templateName }}<span class="tag">{{ detail.version }}</span>
                </div>
              </Col>
              <Col :span="12" class="info-item">
                <div class="info-item-label">{{ $t('tm_template_group') }}：</div>
                <div class="info-item-value">{{ detail.templateGroupName || '-' }}</div>
              </Col>
            </Row>
            <Row style="margin-top:10px;" :gutter="20">
              <Col :span="12" class="info-item">
                <div class="info-item-label">{{ $t('tw_request_des') }}：</div>
                <div class="info-item-value">{{ detail.description || '-' }}</div>
              </Col>
              <Col :span="12" class="info-item">
                <div class="info-item-label">{{ $t('tw_attach') }}：</div>
                <div class="info-item-value">
                  <UploadFile :id="requestId" :files="attachFiles" type="request" formDisable onlyShowFile />
                  <span v-if="attachFiles.length === 0">-</span>
                </div>
              </Col>
            </Row>
            <Row style="margin-top:10px;" :gutter="20">
              <Col :span="12" class="info-item">
                <div class="info-item-label">{{ $t('tw_ref') }}：</div>
                <div class="info-item-value">
                  {{ detail.refName ? `【${typeMap[detail.refType]}】${detail.refName}【${detail.refId}】` : '-' }}
                </div>
              </Col>
            </Row>
          </div>
        </HeaderTitle>
        <!--表单详情-->
        <HeaderTitle :title="$t('tw_form_detail')">
          <div class="request-form">
            <template v-if="detail.customForm && detail.customForm.value">
              <Divider style="margin: 0 0 30px 0" orientation="left">
                <span class="sub-header">{{ $t('tw_information_form') }}</span>
              </Divider>
              <CustomForm
                v-if="detail.customForm && detail.customForm.value"
                v-model="detail.customForm.value"
                :options="detail.customForm.title"
                :requestId="requestId"
                disabled
              ></CustomForm>
            </template>
            <Divider style="margin: 20px 0 30px 0" orientation="left">
              <span class="sub-header">{{ $t('tw_data_form') }}</span>
            </Divider>
            <FormItem
              v-if="detail.associationWorkflow"
              :label="$t('tw_choose_object')"
              :label-width="lang === 'zh-CN' ? 110 : 150"
              required
            >
              <Select v-model="form.rootEntityId" :disabled="true" clearable filterable style="width:300px;">
                <Option v-for="item in rootEntityOptions" :value="item.guid" :key="item.guid">{{
                  item.key_name
                }}</Option>
              </Select>
            </FormItem>
            <EntityTable v-if="form.data.length" :data="form.data" :requestId="requestId" formDisable></EntityTable>
            <div v-else class="no-data">{{ $t('tw_no_formConfig') }}</div>
          </div>
        </HeaderTitle>
        <!--处理历史-->
        <HeaderTitle :title="$t('tw_handle_history')">
          <Collapse v-model="activeStep" simple style="margin-top:30px;">
            <Panel
              v-for="(data, index) in historyData"
              :name="index + ''"
              :key="index"
              :hide-arrow="['revoke', 'submit'].includes(data.type) ? true : false"
            >
              <!--提交-->
              <template v-if="data.type === 'submit'">
                <div class="custom-panel">
                  <div class="custom-panel-title" style="margin-left:24px;">{{ $t('tw_submit_request') }}</div>
                  <HeaderTag
                    class="custom-panel-header"
                    :showHeader="index === 0 ? true : false"
                    :data="data"
                    :operation="$t('tw_submit_request')"
                  ></HeaderTag>
                </div>
              </template>
              <!--撤回-->
              <template v-if="data.type === 'revoke'">
                <div class="custom-panel">
                  <div class="custom-panel-title" style="margin-left:24px;">{{ $t('tw_recall') }}</div>
                  <HeaderTag class="custom-panel-header" :data="data" :operation="$t('tw_recall')"></HeaderTag>
                </div>
              </template>
              <!--定版-->
              <template v-else-if="data.type === 'check'">
                <div class="custom-panel">
                  <div class="custom-panel-title">{{ $t('tw_request_pending') }}</div>
                  <HeaderTag class="custom-panel-header" :data="data" :operation="$t('final_version')"></HeaderTag>
                </div>
                <div slot="content" class="history">
                  <DataBind
                    :isHandle="isHandle"
                    :requestTemplate="requestTemplate"
                    :requestId="requestId"
                    formDisable
                    :showBtn="false"
                    :formData="detail.formData"
                  ></DataBind>
                </div>
              </template>
              <!--审批和任务-->
              <template v-else-if="['approve', 'implement_process', 'implement_custom'].includes(data.type)">
                <div class="custom-panel">
                  <div class="custom-panel-title">
                    <span class="type">{{
                      `${{
                        approve: $t('tw_approval'),
                        implement_custom: $t('task'),
                        implement_process: $t('task')
                      }[data.type] || '-'}: ${data.name}`
                    }}</span>
                    <span class="mode">{{ approvalTypeName[data.handleMode] || '' }}</span>
                  </div>
                  <HeaderTag class="custom-panel-header" :data="data"></HeaderTag>
                </div>
                <div slot="content" class="history">
                  <!--未开启表单过滤-->
                  <!-- <template v-if="!data.filterFlag">
                    <EntityTable
                      v-if="data.formData && data.formData.length"
                      :data="data.formData"
                      :requestId="requestId"
                      formDisable
                    ></EntityTable>
                    <div v-else class="no-data">
                      {{ $t('tw_no_formConfig') }}
                    </div>
                  </template> -->
                  <Tabs>
                    <TabPane v-for="item in data.taskHandleList" :key="item.id" :label="item.handler" :name="item.id">
                      <!--审批和任务操作选择了【无需处理】不展示表单-->
                      <div v-if="item.handleResult !== 'unrelated'">
                        <EntityTable
                          v-if="item.formData && item.formData.length"
                          :data="item.formData"
                          :requestId="requestId"
                          formDisable
                        ></EntityTable>
                        <div v-else class="no-data">
                          {{ $t('tw_no_formConfig') }}
                        </div>
                      </div>
                      <div v-else class="no-data">
                        {{ '用户选择无需处理,未提交表单' }}
                      </div>
                    </TabPane>
                  </Tabs>
                </div>
              </template>
              <!--确认-->
              <template v-else-if="data.type === 'confirm'">
                <div class="custom-panel">
                  <div class="custom-panel-title">{{ $t('tw_request_confirm') }}</div>
                  <HeaderTag class="custom-panel-header" :data="data" :operation="$t('tw_request_confirm')"></HeaderTag>
                </div>
                <div slot="content" class="history">
                  <Form :label-width="80" label-position="left">
                    <FormItem :label="$t('status')">
                      <Input disabled :value="completeStatus" />
                    </FormItem>
                    <FormItem :label="$t('tw_uncompleted_tag')">
                      <Input disabled :value="uncompletedTasks.join(', ')" />
                    </FormItem>
                  </Form>
                </div>
              </template>
            </Panel>
          </Collapse>
        </HeaderTitle>
        <!--当前处理-->
        <CurrentHandle
          v-if="isHandle && Object.keys(handleData).length > 0"
          :detail="detail"
          :handleData="handleData"
          :actionName="actionName"
        />
      </Form>
    </div>
    <div class="footer-btn">
      <!--撤回-->
      <Button
        v-if="jumpFrom === 'submit' && ['Pending', 'InApproval'].includes(detail.status) && detail.revokeBtn"
        type="error"
        @click="handleRecall"
        >{{ $t('tw_recall') }}</Button
      >
    </div>
    <!--编排流程图-->
    <template v-if="detail.associationWorkflow">
      <div class="expand-btn" :style="{ right: flowVisible ? '440px' : '0px' }" @click="flowVisible = !flowVisible">
        <Icon v-if="flowVisible" type="ios-arrow-dropright-circle" :size="28" />
        <Icon v-else type="ios-arrow-dropleft-circle" :size="28" />
      </div>
      <div v-if="flowVisible" class="flow-expand">
        <StaticFlow v-if="!detail.procInstanceId" :requestTemplate="requestTemplate"></StaticFlow>
        <DynamicFlow v-else :flowId="detail.procInstanceId"></DynamicFlow>
      </div>
    </template>
  </div>
</template>

<script>
import HeaderTitle from './header-title.vue'
import HeaderTag from './header-tag.vue'
import StaticFlow from './flow/static-flow.vue'
import DynamicFlow from './flow/dynamic-flow.vue'
import EntityTable from './entity-table.vue'
import DataBind from './data-bind.vue'
import UploadFile from './upload.vue'
import BaseProgress from './base-progress.vue'
import CustomForm from './custom-form.vue'
import CurrentHandle from './base-handle.vue'
import { getPublishInfo, recallRequest, getRequestHistory } from '@/api/server'
export default {
  components: {
    HeaderTitle,
    HeaderTag,
    StaticFlow,
    DynamicFlow,
    EntityTable,
    DataBind,
    UploadFile,
    BaseProgress,
    CustomForm,
    CurrentHandle
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
      isHandle: this.$route.query.isHandle === 'Y', // 处理标志
      jumpFrom: this.$route.query.jumpFrom || '', // 入口tab标记
      type: this.$route.query.type, // 首页类型type
      requestTemplate: this.$route.query.requestTemplate,
      requestId: this.$route.query.requestId,
      taskId: this.$route.query.taskId,
      taskHandleId: this.$route.query.taskHandleId, // 当前处理人任务ID(协同并行可能有多人)
      detail: {}, // 详情信息
      form: {
        name: '',
        description: '',
        expectTime: '',
        rootEntityId: '', // 目标对象
        data: []
      },
      rootEntityOptions: [],
      historyData: [], // 处理历史数据
      handleData: {}, // 当前处理数据
      activeStep: '', // 处理历史展开项
      attachFiles: [], // 请求附件
      completeStatus: '', // 请求确认-状态
      uncompletedTasks: [], // 请求确认-未完成任务
      flowVisible: false,
      approvalTypeName: {
        custom: this.$t('tw_onlyOne'), // 单人
        any: this.$t('tw_anyWidth'), // 协同
        all: this.$t('tw_allWidth'), // 并行
        admin: this.$t('tw_roleAdmin'), // 提交人角色管理员
        auto: this.$t('tw_autoWith') // 自动通过
      },
      typeMap: {
        1: this.$t('tw_publish'),
        2: this.$t('tw_request'),
        3: this.$t('tw_question'),
        4: this.$t('tw_event'),
        5: this.$t('fork')
      },
      handleTypeColor: {
        check: '#ffa2d3',
        approve: '#2d8cf0',
        implement_process: '#cba43f',
        implement_custom: '#b886f8',
        confirm: '#19be6b'
      },
      lang: window.localStorage.getItem('lang') || 'zh-CN'
    }
  },
  computed: {
    getStatusName () {
      return function (val) {
        const list = [
          { label: this.$t('status_pending'), value: 'Pending', color: '#b886f8' },
          { label: this.$t('tw_inApproval'), value: 'InApproval', color: '#1990ff' },
          { label: this.$t('status_inProgress'), value: 'InProgress', color: '#1990ff' },
          { label: this.$t('tw_request_confirm'), value: 'Confirm', color: '#b886f8' },
          { label: this.$t('status_inProgress_faulted'), value: 'InProgress(Faulted)', color: '#f26161' },
          { label: this.$t('status_termination'), value: 'Termination', color: '#e29836' },
          { label: this.$t('status_complete'), value: 'Completed', color: '#7ac756' },
          { label: this.$t('status_inProgress_timeouted'), value: 'InProgress(Timeouted)', color: '#f26161' },
          { label: this.$t('status_faulted'), value: 'Faulted', color: '#e29836' },
          { label: this.$t('status_draft'), value: 'Draft', color: '#808695' },
          { label: this.$t('tw_stop'), value: 'Stop', color: '#ed4014' }
        ]
        const item = list.find(i => i.value === val) || {}
        return item.label
      }
    },
    formatHandler () {
      let list = (this.detail.handler && this.detail.handler.split(',')) || []
      list = list.filter(i => i)
      return Array.from(new Set(list)).join(', ')
    }
  },
  async created () {
    // 获取请求进度
    this.$nextTick(() => {
      this.$refs.progress.initData(this.requestId)
    })
    // 获取详情信息
    this.getRequestInfoNew()
  },
  methods: {
    handleToHome () {
      if (this.$route.query.type) {
        this.$router.push({
          path: `/taskman/workbench?tabName=${this.jumpFrom}&actionName=${this.actionName}&${
            this.jumpFrom === 'submit' ? 'rollback' : 'type'
          }=${this.type}`
        })
      } else {
        this.$router.back()
      }
    },
    // 获取详情数据
    async getRequestInfoNew () {
      const params = {
        params: {
          requestId: this.requestId,
          taskId: this.$route.query.taskId || ''
        }
      }
      const { statusCode, data } = await getPublishInfo(params)
      if (statusCode === 'OK') {
        this.detail = data.request || {}
        const { name, rootEntityId, operatorObj, description, expectTime, formData, attachFiles } = this.detail
        this.attachFiles = attachFiles
        this.form = Object.assign({}, this.form, {
          name,
          rootEntityId,
          description,
          expectTime,
          data: formData
        })
        this.rootEntityOptions.push({
          guid: rootEntityId,
          key_name: operatorObj
        })
        this.getRequestHistory()
      }
    },
    // 获取请求历史
    async getRequestHistory () {
      const { statusCode, data } = await getRequestHistory(this.requestId)
      if (statusCode === 'OK') {
        this.historyData = data.task || []
        this.historyData.forEach(val => {
          const list = val.taskHandleList || []
          list.forEach(item => {
            if (item.id === this.taskHandleId) {
              // 当前处理，审批和任务表单如果开启了过滤功能，formData需要替换成taskHandleList里面的
              if (['approve', 'implement_process', 'implement_custom'].includes(val.type) && val.filterFlag) {
                val.formData = item.formData
              }
            }
          })
        })
        // 请求确认数据
        const statusMap = {
          complete: this.$t('tw_completed'),
          uncompleted: this.$t('tw_incomplete')
        }
        this.completeStatus = statusMap[data.request.completeStatus] || ''
        this.uncompletedTasks = data.request.uncompletedTasks || []
        // 当前处理【任务、审批、请求确认】
        if (['Pending', 'InProgress', 'InApproval', 'Confirm'].includes(this.detail.status)) {
          // 当前处理数据
          this.handleData = this.historyData.find(item => item.id === this.taskId && item.editable === true) || {}
          // 处理历史列表
          this.historyData = this.historyData.filter(item => item.editable === false)
        }
        if (this.isHandle === false) {
          this.activeStep = this.historyData.length - 1 + ''
        }
      }
    },
    // 撤回
    async handleRecall () {
      this.$Modal.confirm({
        title: this.$t('confirm'),
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
            this.$router.push({ path: `/taskman/workbench?tabName=submit&actionName=${this.actionName}&rollback=3` })
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
  }
  .content {
    min-height: 500px;
    display: flex;
    .sub-header {
      font-size: 14px;
      color: #515a6e;
      font-weight: bold;
    }
    .request-form {
      width: calc(100% - 20px);
      margin: 0 0 12px 16px;
    }
    .no-data {
      padding-left: 20px;
      height: 60px;
      line-height: 60px;
      color: #515a6e;
    }
    .info-item {
      display: flex;
      &-label {
        color: #515a6e;
        width: 140px;
        word-wrap: break-word;
      }
      &-value {
        color: #515a6e;
        width: calc(100% - 100px);
        word-wrap: break-word;
        .tag {
          background: #f7f7f7;
          padding: 4px 8px;
        }
      }
    }
    .custom-panel {
      display: flex;
      align-items: flex-start;
      width: 100%;
      &-title {
        width: 105px;
        font-weight: bold;
        line-height: 20px;
        display: flex;
        flex-direction: column;
        justify-content: center;
        // .type {
        //   font-size: 12px;
        //   display: inline-block;
        //   color: #fff;
        //   padding: 1px 5px;
        //   border-radius: 2px;
        // }
        .mode {
          font-size: 12px;
          display: inline-block;
          color: #19be6b;
          margin-top: 2px;
        }
      }
      &-header {
        width: calc(100% - 70px);
      }
    }
    .history {
      padding: 20px;
      border: 1px dashed #d7dadc;
      margin: 16px 0;
      .no-data {
        height: 60px;
        line-height: 60px;
        color: #515a6e;
      }
    }
  }
  .footer-btn {
    padding: 20px 0 30px 0;
    display: flex;
    justify-content: center;
    align-items: center;
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
    align-items: flex-start;
    padding-left: 0px !important;
    i {
      margin-top: 4px;
      margin-right: 8px !important;
    }
  }
  .ivu-collapse-content-box {
    padding-bottom: 0px;
  }
  .ivu-collapse-content {
    padding: 0 !important;
  }
}
</style>
