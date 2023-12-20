<template>
  <div class="workbench-request-create">
    <Row class="back-header">
      <Icon size="26" type="md-arrow-back" style="cursor:pointer" @click="$router.back()" />
      <span v-if="isAdd" class="name">
        {{ `${templateName}` }}<Tag size="medium">{{ version }}</Tag>
      </span>
      <span v-else class="name">
        {{ `${detailInfo.name}` }}
      </span>
    </Row>
    <Row class="w-header">
      <Col span="19" class="steps">
        <span class="title">请求进度：</span>
        <Steps :current="0" style="max-width:600px;">
          <Step v-for="(i, index) in progressList" :key="index" :content="i.name">
            <template #icon>
              <Icon size="24" :type="i.icon" :color="i.color" />
            </template>
            <div class="role" slot="content">
              <div class="word-eclipse">{{ i.name }}</div>
              <span>{{ i.handler }}</span>
            </div>
          </Step>
        </Steps>
        <div v-if="errorNode" style="margin:0 0 10px 20px;max-width:500px;">
          <Alert v-if="errorNode === 'autoExit'" type="error">
            请求已终止，因为编排流程根据判断条件走到了退出节点，请复制请求ID，发送给wecube平台管理员处理
          </Alert>
          <Alert v-else-if="errorNode === 'internallyTerminated'" type="error">
            请求已终止，因为编排执行已被管理员手动终止，请复制请求ID，发送给wecube平台管理员处理
          </Alert>
          <Alert v-else type="error">
            {{ errorNode }}节点报错，流程已暂停，请复制请求ID，发送给wecube平台管理员处理
          </Alert>
        </div>
      </Col>
      <Col span="5" class="btn-group">
        <Button v-if="isAdd" @click="handleDraft(false)" style="margin-right:10px;">保存草稿</Button>
        <Button v-if="isAdd" type="primary" @click="handlePublish">提交</Button>
        <Button v-if="jumpFrom === 'submit' && detailInfo.status === 'Pending'" type="error" @click="handleRecall"
          >撤回</Button
        >
      </Col>
    </Row>
    <div style="display:flex;" class="content">
      <div style="width:calc(100% - 420px)" class="split-line">
        <Form :model="form" label-position="right" :label-width="120">
          <template v-if="isAdd">
            <HeaderTitle title="发布信息">
              <FormItem label="请求名称" required>
                <Input v-model="form.name" :maxlength="50" show-word-limit placeholder="请输入" style="width:100%;" />
              </FormItem>
              <FormItem label="发布描述">
                <Input
                  v-model="form.description"
                  type="textarea"
                  :maxlength="100"
                  show-word-limit
                  placeholder="请输入"
                  style="width:100%;"
                />
              </FormItem>
              <FormItem :label="$t('expected_completion_time')">
                <DatePicker
                  type="datetime"
                  :value="form.expectTime"
                  @on-change="
                    val => {
                      form.expectTime = val
                    }
                  "
                  placeholder="Select date"
                  :options="{
                    disabledDate(date) {
                      return date && date.valueOf() < Date.now() - 86400000
                    }
                  }"
                  style="width:400px;"
                ></DatePicker>
              </FormItem>
            </HeaderTitle>
            <HeaderTitle title="发布目标对象">
              <FormItem label="选择操作单元" required>
                <Select v-model="form.rootEntityId" :disabled="formDisable" clearable filterable style="width:300px;">
                  <Option v-for="item in rootEntityOptions" :value="item.guid" :key="item.guid">{{
                    item.key_name
                  }}</Option>
                </Select>
              </FormItem>
              <FormItem v-if="requestData.length" label="已选择">
                <EntityTable
                  ref="entityTable"
                  :data="requestData"
                  :requestId="requestId"
                  :isAdd="isAdd"
                  type="request"
                ></EntityTable>
              </FormItem>
            </HeaderTitle>
          </template>
          <template v-else>
            <HeaderTitle title="请求信息">
              <Row :gutter="20">
                <Col :span="3">请求ID：</Col>
                <Col :span="9">{{ detailInfo.id }}</Col>
                <Col :span="3">请求类型：</Col>
                <Col :span="9">{{ { 0: '请求', 1: '发布' }[detailInfo.requestType] }}</Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="3">创建时间：</Col>
                <Col :span="9">{{ detailInfo.createdTime }}</Col>
                <Col :span="3">期望完成时间：</Col>
                <Col :span="9">{{ detailInfo.expectTime }}</Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="3">请求进度：</Col>
                <Col :span="9">
                  <Progress :percent="detailInfo.progress" style="width:150px;" />
                </Col>
                <Col :span="3">请求状态：</Col>
                <Col :span="9">{{ getStatusName(detailInfo.status) }}</Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="3">当前节点：</Col>
                <Col :span="9">{{
                  {
                    waitCommit: '等待提交',
                    sendRequest: '提起请求',
                    requestPending: '请求定版',
                    requestComplete: '请求完成'
                  }[detailInfo.curNode] || detailInfo.curNode
                }}</Col>
                <Col :span="3">当前处理人：</Col>
                <Col :span="9">{{ detailInfo.handler }}</Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="3">创建人：</Col>
                <Col :span="9">{{ detailInfo.createdBy }}</Col>
                <Col :span="3">创建人角色：</Col>
                <Col :span="9">{{ detailInfo.role }}</Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="3">使用模板：</Col>
                <Col :span="9"
                  >{{ detailInfo.templateName }}<Tag>{{ version }}</Tag></Col
                >
                <Col :span="3">模板组：</Col>
                <Col :span="9">{{ detailInfo.templateGroupName }}</Col>
              </Row>
              <Row style="margin-top:10px;" :gutter="20">
                <Col :span="3">请求描述：</Col>
                <Col :span="9">{{ detailInfo.description }}</Col>
              </Row>
            </HeaderTitle>
            <!--处理历史-->
            <HeaderTitle title="处理历史">
              <Collapse v-model="activeStep" simple style="margin-top:30px;">
                <Panel v-for="(data, index) in historyData" :name="index + ''" :key="index">
                  <!--提交请求-->
                  <template v-if="index === 0">
                    <div style="display:flex;align-items:center;width:100%;">
                      <div style="width:70px;font-weight:bold;">提交请求</div>
                      <div style="width:calc(100% - 70px)">
                        <HeaderTag :showHeader="true" :data="data" operation="提交审批"></HeaderTag>
                      </div>
                    </div>
                    <div slot="content" class="history">
                      <FormItem label="选择操作单元" required>
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
                      <div style="width:70px;font-weight:bold;">请求定版</div>
                      <div style="width:calc(100% - 70px)">
                        <HeaderTag :data="data" operation="确认定版"></HeaderTag>
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
                      <div style="width:70px;font-weight:bold;">{{ data.taskName }}</div>
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
                        <div style="text-align: center">
                          <Button v-if="data.editable" :disabled="enforceDisable" @click="saveTaskData" type="info">{{
                            $t('save')
                          }}</Button>
                          <Button
                            v-if="data.editable"
                            :disabled="
                              enforceDisable ||
                                (data.nextOption && data.nextOption.length !== 0 && data.choseOption === '')
                            "
                            @click="commitTaskData"
                            type="primary"
                            >{{ $t('commit') }}</Button
                          >
                        </div>
                      </div>
                    </div>
                  </template>
                </Panel>
              </Collapse>
            </HeaderTitle>
            <!--当前处理-任务审批-->
            <HeaderTitle
              v-if="isHandle && detailInfo.status === 'InProgress'"
              title="当前处理"
              :subTitle="handleData.taskName"
            >
              <Steps :current="1" direction="vertical">
                <Step>
                  <div slot="title" class="task-step">
                    <div>第一步：填写任务表单</div>
                    <div>（请按格式填写任务表单）</div>
                  </div>
                  <div slot="content" style="padding:20px 0px;">
                    <EntityTable
                      :data="handleData.formData"
                      :requestId="requestId"
                      :formDisable="!handleData.editable || enforceDisable"
                    ></EntityTable>
                  </div>
                </Step>
                <Step>
                  <div slot="title" class="task-step">
                    <div>第二步：填写反馈结果</div>
                    <div>（默认为同意，支持修改，支持附件上传）</div>
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
                        >{{ $t('commit') }}</Button
                      >
                    </div>
                  </div>
                </Step>
              </Steps>
            </HeaderTitle>
            <!--当前处理-请求定版-->
            <HeaderTitle v-if="isHandle && detailInfo.status === 'Pending'" title="当前处理" subTitle="请求定版">
              <Steps :current="1" direction="vertical">
                <Step>
                  <div slot="title" class="task-step">
                    <div>第一步：表单填写审核</div>
                    <div>（确认表单填写格式、数据正确不可修改）</div>
                  </div>
                  <div slot="content" style="padding:20px 0px;">
                    <FormItem label="选择操作单元" required>
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
                    <div>第二步：编排节点使用数据确认</div>
                    <div>（确认编排的自动节点绑定的对象可以修改）</div>
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
      <!--编排流程-->
      <div style="width:420px;padding-left:20px;">
        <StaticFlow
          v-if="isAdd || ['Draft', 'Pending'].includes(detailInfo.status)"
          :requestTemplate="requestTemplate"
        ></StaticFlow>
        <DynamicFlow v-else :flowId="detailInfo.procInstanceId"></DynamicFlow>
      </div>
    </div>
  </div>
</template>

<script>
import HeaderTitle from '../components/header-title.vue'
import HeaderTag from '../components/header-tag.vue'
import StaticFlow from '../components/flow/static-flow.vue'
import DynamicFlow from '../components/flow/dynamic-flow.vue'
import EntityTable from '../components/entity-table.vue'
import DataBind from '../components/data-bind.vue'
import {
  getCreateInfo,
  getProgressInfo,
  getRootEntity,
  getEntityData,
  savePublishData,
  getPublishInfo,
  getRequestInfo,
  updateRequestStatus,
  saveTaskData,
  commitTaskData,
  recallRequest
} from '@/api/server'
import dayjs from 'dayjs'
const statusIcon = {
  1: 'md-pin', // 进行中
  2: 'md-radio-button-on', // 未开始
  3: 'ios-checkmark-circle-outline', // 已完成
  4: 'md-close-circle', // 节点失败(包含超时)
  5: 'md-exit', // 自动退出
  6: 'md-exit' // 手动终止
}
const statusColor = {
  1: '#ffa500',
  2: '#8189a5',
  3: '#19be6b',
  4: '#ed4014',
  5: '#ed4014',
  6: '#ed4014'
}
export default {
  components: {
    HeaderTitle,
    HeaderTag,
    StaticFlow,
    DynamicFlow,
    EntityTable,
    DataBind
  },
  data () {
    return {
      actionName: '2', // 1发布,2请求(3问题,4事件,5变更)
      templateName: '',
      version: '', // 模板版本号
      enforceDisable: this.$route.query.enforceDisable === 'Y',
      isAdd: this.$route.query.isAdd === 'Y',
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
      progressList: [],
      requestData: [], // 发布目标对象表格数据
      historyData: [], // 处理历史数据
      handleData: {},
      activeStep: '', // 处理历史当前展开
      errorNode: ''
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
  watch: {
    'form.rootEntityId' (val) {
      if (val) {
        this.getEntityData()
      } else {
        this.requestData = []
      }
    }
  },
  async mounted () {
    // 路由有requestId详情接口，无requestId创建接口
    if (this.requestId) {
      await this.getPublishInfo()
      this.getRequestInfo()
    } else {
      await this.getCreateInfo()
    }
    this.getProgressInfo()
    this.getEntity()
  },
  methods: {
    async getCreateInfo () {
      const params = {
        requestTemplate: this.requestTemplate,
        role: this.$route.query.role
      }
      const { statusCode, data } = await getCreateInfo(params)
      if (statusCode === 'OK') {
        this.requestId = data.id
        this.templateName = data.requestTemplateName
        this.version = data.templateVersion
        this.form.name = data.name
        this.form.expectTime = dayjs()
          .add(data.expireDay || 0, 'day')
          .format('YYYY-MM-DD HH:mm:ss')
      }
    },
    // 获取请求进度
    async getProgressInfo () {
      const params = {
        templateId: this.requestTemplate,
        requestId: this.requestId
      }
      const { statusCode, data } = await getProgressInfo(params)
      if (statusCode === 'OK') {
        this.progressList = data
        this.progressList.forEach(item => {
          item.icon = statusIcon[item.status]
          item.color = statusColor[item.status]
          switch (item.node) {
            case 'sendRequest':
              item.name = '提起请求'
              break
            case 'requestPending':
              item.name = '请求定版'
              break
            case 'requestComplete':
              item.name = '请求完成'
              break
            case 'autoExit':
              item.name = '自动退出'
              this.errorNode = item.node
              break
            case 'internallyTerminated':
              item.name = '手动终止'
              this.errorNode = item.node
              break
            default:
              item.name = item.node
              break
          }
          if (item.handler === 'autoNode') {
            item.handler = '自动节点'
            this.errorNode = item.name
          }
        })
      }
    },
    async getRequestInfo () {
      const { statusCode, data } = await getRequestInfo(this.requestId)
      if (statusCode === 'OK') {
        // this.attachFiles = data.attachFiles
        this.form.rootEntityId = data.cache
      }
    },
    // 获取发布信息
    async getPublishInfo () {
      const { statusCode, data } = await getPublishInfo(this.requestId)
      if (statusCode === 'OK') {
        this.detailInfo = data.request || {}
        this.templateName = data.request.templateName
        this.version = data.request.version
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
    // 操作目标对象下拉值
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
    // 获取目标对象对应数据
    async getEntityData () {
      let params = {
        params: {
          requestId: this.requestId,
          rootEntityId: this.form.rootEntityId
        }
      }
      const { statusCode, data } = await getEntityData(params)
      if (statusCode === 'OK') {
        this.requestData = data.data
      }
    },
    // 保存草稿
    async handleDraft (noJump) {
      if (!this.form.rootEntityId) {
        this.$Message.warning(this.$t('root_entity') + this.$t('can_not_be_empty'))
        return
      }
      // 提取表格勾选的数据
      this.form.data =
        (this.$refs.entityTable &&
          this.$refs.entityTable.requestData.map(item => {
            if (Array.isArray(item.value)) {
              item.value = item.value.filter(j => {
                return j.entityData._checked
              })
            }
            return item
          })) ||
        []
      const item = this.rootEntityOptions.find(item => item.guid === this.form.rootEntityId) || {}
      this.form.entityName = item.key_name
      if (this.requestDataCheck()) {
        const { statusCode } = await savePublishData(this.requestId, this.form)
        if (statusCode === 'OK') {
          if (noJump) {
            return statusCode
          } else {
            this.$router.push({ path: `/taskman/workbench?tabName=draft&actionName=${this.actionName}` })
          }
        }
      } else {
        this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('required_tip')
        })
      }
    },
    // 发布
    async handlePublish () {
      this.$Modal.confirm({
        title: this.$t('confirm') + this.$t('commit'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const draftResult = await this.handleDraft(true)
          if (draftResult === 'OK') {
            const { statusCode } = await updateRequestStatus(this.requestId, 'Pending')
            if (statusCode === 'OK') {
              this.$router.push({ path: `/taskman/workbench?tabName=submit&actionName=${this.actionName}` })
            }
          }
        },
        onCancel: () => {}
      })
    },
    // 校验表格数据必填项
    requestDataCheck () {
      let result = true
      this.form.data.forEach(requestData => {
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
              }
            } else {
              if (val === '') {
                result = false
              }
            }
          })
        })
      })
      return result
    },
    // 任务审批保存
    async saveTaskData () {
      // const taskData = this.historyData.find(d => d.editable === true)
      const taskData = this.handleData
      const result = this.paramsCheck(taskData)
      if (result) {
        const { statusCode } = await saveTaskData(taskData.taskId, taskData)
        if (statusCode === 'OK') {
          this.$Notice.success({
            title: this.$t('successful'),
            desc: this.$t('successful')
          })
        }
      } else {
        this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('required_tip')
        })
      }
    },
    // 任务审批提交
    async commitTaskData () {
      // const taskData = this.historyData.find(d => d.editable === true)
      const taskData = this.handleData
      const { statusCode } = await commitTaskData(taskData.taskId, taskData)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.$router.push({ path: `/taskman/workbench?tabName=hasProcessed&actionName=${this.actionName}` })
      }
    },
    paramsCheck (taskData) {
      let result = true
      taskData.formData.forEach(requestData => {
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
              }
            } else {
              if (val === '') {
                result = false
              }
            }
          })
        })
      })
      return result
    },
    // 表格操作撤回
    async handleRecall () {
      this.$Modal.confirm({
        title: this.$t('confirm') + '撤回',
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
.workbench-request-create {
  .back-header {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    .name {
      font-size: 16px;
      margin-left: 20px;
      display: flex;
      align-items: center;
    }
  }
  .w-header {
    padding-bottom: 15px;
    margin-bottom: 20px;
    border-bottom: 2px dashed #d7dadc;
    .steps {
      display: flex;
      align-items: center;
      .title {
        font-size: 14px;
        font-weight: 500;
        margin-right: 20px;
      }
      .role {
        display: flex;
        flex-direction: column;
      }
      .word-eclipse {
        max-width: 100px;
        text-overflow: ellipsis;
        overflow: hidden;
        white-space: nowrap;
      }
    }
    .btn-group {
      text-align: right;
    }
  }
  .content {
    min-height: 500px;
    .split-line {
      border-right: 2px dashed #d7dadc;
      padding-right: 20px;
    }
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
}
</style>
<style lang="scss">
.workbench-request-create {
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
    margin-bottom: 25px;
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
  .ivu-icon {
    font-weight: bold;
  }
}
</style>
