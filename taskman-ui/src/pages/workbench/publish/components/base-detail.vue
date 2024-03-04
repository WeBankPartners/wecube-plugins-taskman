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
    <div class="content">
      <Form :model="form" label-position="right" :label-width="120" style="width:100%;">
        <template>
          <!--请求信息-->
          <HeaderTitle :title="$t('tw_request_title')">
            <div style="padding-left:16px;">
              <Row :gutter="20">
                <Col :span="12" class="info-item">
                  <div class="info-item-label">{{ $t('request_id') }}：</div>
                  <div class="info-item-value">{{ detailInfo.id || '--' }}</div>
                </Col>
                <Col :span="12" class="info-item">
                  <div class="info-item-label">{{ $t('tw_request_type') }}：</div>
                  <div class="info-item-value">
                    {{ { 0: $t('tw_request'), 1: $t('tw_publish') }[detailInfo.requestType] || '--' }}
                  </div>
                </Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="12" class="info-item">
                  <div class="info-item-label">{{ $t('createdTime') }}：</div>
                  <div class="info-item-value">{{ detailInfo.createdTime || '--' }}</div>
                </Col>
                <Col :span="12" class="info-item">
                  <div class="info-item-label">{{ $t('expected_completion_time') }}：</div>
                  <div class="info-item-value">{{ detailInfo.expectTime || '--' }}</div>
                </Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="12" class="info-item">
                  <div class="info-item-label">{{ $t('tw_request_progress') }}：</div>
                  <div class="info-item-value">
                    <Progress :percent="detailInfo.progress" style="width:150px;" />
                  </div>
                </Col>
                <Col :span="12" class="info-item">
                  <div class="info-item-label">{{ $t('tw_request_status') }}：</div>
                  <div class="info-item-value">{{ getStatusName(detailInfo.status) || '--' }}</div>
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
                      }[detailInfo.curNode] ||
                        detailInfo.curNode ||
                        '--'
                    }}
                  </div>
                </Col>
                <Col :span="12" class="info-item">
                  <div class="info-item-label">
                    {{ detailInfo.status === 'Draft' ? $t('tw_pending_handler') : $t('tw_cur_handler') }}：
                  </div>
                  <div class="info-item-value">{{ detailInfo.handler || '--' }}</div>
                </Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="12" class="info-item">
                  <div class="info-item-label">{{ $t('createdBy') }}：</div>
                  <div class="info-item-value">{{ detailInfo.createdBy || '--' }}</div>
                </Col>
                <Col :span="12" class="info-item">
                  <div class="info-item-label">{{ $t('tw_creatby_role') }}：</div>
                  <div class="info-item-value">{{ detailInfo.role || '--' }}</div>
                </Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="12" class="info-item">
                  <div class="info-item-label">{{ $t('tw_use_template') }}：</div>
                  <div class="info-item-value">
                    {{ detailInfo.templateName }}<span class="tag">{{ detailInfo.version }}</span>
                  </div>
                </Col>
                <Col :span="12" class="info-item">
                  <div class="info-item-label">{{ $t('tm_template_group') }}：</div>
                  <div class="info-item-value">{{ detailInfo.templateGroupName || '--' }}</div>
                </Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="12" class="info-item">
                  <div class="info-item-label">{{ $t('tw_request_des') }}：</div>
                  <div class="info-item-value">{{ detailInfo.description || '--' }}</div>
                </Col>
              </Row>
              <!--自定义信息表单-->
              <CustomForm
                v-if="detailInfo.customForm"
                v-model="detailInfo.customForm.value"
                :options="detailInfo.customForm.title"
                :requestId="requestId"
                disabled
                :label-width="140"
                labelPosition="left"
                style="margin-top: 10px;"
              ></CustomForm>
            </div>
          </HeaderTitle>
          <!--请求表单-->
          <HeaderTitle title="请求表单">
            <FormItem :label="$t('tw_choose_object')" required>
              <Select v-model="form.rootEntityId" :disabled="true" clearable filterable style="width:300px;">
                <Option v-for="item in rootEntityOptions" :value="item.guid" :key="item.guid">{{
                  item.key_name
                }}</Option>
              </Select>
            </FormItem>
            <EntityTable
              v-if="historyData[0]"
              :data="historyData[0] && historyData[0].formData"
              :requestId="requestId"
              :formDisable="true"
              style="width:calc(100% - 20px);margin-left:16px;"
            ></EntityTable>
          </HeaderTitle>
          <!--处理历史-->
          <HeaderTitle :title="$t('tw_handle_history')">
            <Collapse v-model="activeStep" simple style="margin-top:30px;">
              <Panel
                v-for="(data, index) in historyData"
                :name="index + ''"
                :key="index"
                :hide-arrow="index === 0 ? true : false"
              >
                <!--提交请求-->
                <template v-if="index === 0">
                  <div class="custom-panel">
                    <div class="custom-panel-title" style="margin-left:30px;">{{ $t('tw_submit_request') }}</div>
                    <HeaderTag
                      class="custom-panel-header"
                      :showHeader="true"
                      :data="data"
                      :operation="$t('tw_submit_approval')"
                    ></HeaderTag>
                  </div>
                </template>
                <!--请求定版-->
                <template v-else-if="index === 1">
                  <div class="custom-panel">
                    <div class="custom-panel-title">{{ $t('tw_request_pending') }}</div>
                    <HeaderTag class="custom-panel-header" :data="data" :operation="$t('final_version')"></HeaderTag>
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
                  <div class="custom-panel">
                    <div class="custom-panel-title">{{ data.taskName }}</div>
                    <HeaderTag class="custom-panel-header" :data="data"></HeaderTag>
                  </div>
                  <div slot="content" class="history">
                    <EntityTable
                      :data="data.formData"
                      :requestId="requestId"
                      :formDisable="!data.editable"
                    ></EntityTable>
                    <div>
                      <Form :label-width="80" style="margin: 16px 0">
                        <FormItem v-if="data.requestId === ''" :label="$t('task') + $t('description')">
                          <Input disabled v-model="data.description" type="textarea" />
                        </FormItem>
                        <FormItem :label="$t('process_result')" v-if="data.nextOption && data.nextOption.length !== 0">
                          <span slot="label">
                            {{ $t('process_result') }}
                            <span style="color: #ed4014"> * </span>
                          </span>
                          <Select v-model="data.choseOption" disabled>
                            <Option v-for="option in data.nextOption" :value="option" :key="option">{{
                              option
                            }}</Option>
                          </Select>
                        </FormItem>
                        <FormItem :label="$t('process_comments')">
                          <Input disabled v-model="data.comment" type="textarea" />
                        </FormItem>
                      </Form>
                    </div>
                  </div>
                </template>
              </Panel>
            </Collapse>
          </HeaderTitle>
          <!--当前处理-->
          <CurrentHandle :detailInfo="detailInfo" :handleData="handleData" :actionName="actionName" />
        </template>
      </Form>
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
import BaseProgress from './progress.vue'
import CustomForm from '../../components/custom-form.vue'
import CurrentHandle from './handle.vue'
import { getRootEntity, getPublishInfo, getRequestInfo, recallRequest } from '@/api/server'
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
      historyData: [], // 处理历史数据
      handleData: {},
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
    display: flex;
    .info-item {
      display: flex;
      &-label {
        color: #515a6e;
        width: 140px;
        word-wrap: break-word;
      }
      &-value {
        color: #515a6e;
        width: calc(100% - 120px);
        word-wrap: break-word;
        .tag {
          background: #f7f7f7;
          padding: 4px 8px;
        }
      }
    }
    .custom-panel {
      display: flex;
      align-items: center;
      width: 100%;
      &-title {
        width: 70px;
        font-weight: bold;
        line-height: 22px;
      }
      &-header {
        width: calc(100% - 70px);
      }
    }
    .history {
      padding: 20px;
      border: 1px dashed #d7dadc;
      margin: 16px 0;
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
  .ivu-collapse-content-box {
    padding-bottom: 0px;
  }
  .ivu-collapse-content {
    padding: 0;
  }
}
</style>
