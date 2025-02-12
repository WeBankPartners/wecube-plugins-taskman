<template>
  <div class="workbench-publish-create">
    <Row class="back-header">
      <Icon size="22" type="md-arrow-back" class="icon" @click="handleToHome" />
      <span class="name">
        {{ `${detail.templateName || detail.requestTemplateName || ''}` }}
        <Tag size="medium">{{ detail.version || detail.templateVersion || '' }}</Tag>
      </span>
    </Row>
    <div class="w-header">
      <div class="w-header-progress">
        <!--请求进度-->
        <BaseProgress ref="progress"></BaseProgress>
      </div>
      <div class="w-header-btn">
        <!--保存草稿-->
        <Button @click="handleDraft(false)" style="margin-right:10px;">{{ $t('tw_save_draft') }}</Button>
        <!--提交-->
        <Button type="primary" @click="handlePublish">{{ $t('tw_commit') }}</Button>
      </div>
    </div>
    <div class="content">
      <Form :model="form" label-position="right" :label-width="120" style="width:100%;">
        <!--基础信息-->
        <BaseHeaderTitle :title="$t('tw_request_title')">
          <!--请求名-->
          <FormItem :label="$t('request_name')" required>
            <Input
              v-model="form.name"
              :maxlength="70"
              show-word-limit
              :placeholder="$t('request_name')"
              style="width:60%;"
            />
          </FormItem>
          <!--期望完成时间-->
          <FormItem :label="$t('expected_completion_time')" required>
            <DatePicker
              type="datetime"
              :value="form.expectTime"
              @on-change="
                val => {
                  form.expectTime = val
                }
              "
              :editable="false"
              :placeholder="$t('tw_please_select')"
              :options="{
                disabledDate(date) {
                  return date && date.valueOf() < Date.now() - 86400000
                }
              }"
              style="width:400px;"
              :clearable="false"
            ></DatePicker>
          </FormItem>
          <!--请求描述-->
          <FormItem :label="$t('tw_publish_des')">
            <Input
              v-model="form.description"
              type="textarea"
              :maxlength="200"
              show-word-limit
              :placeholder="$t('tw_publish_des')"
              style="width:60%;"
            />
          </FormItem>
          <!--关联单-->
          <FormItem :label="$t('tw_ref')">
            <Select v-model="form.refType" @on-change="handleRefTypeChange" style="width:100px;">
              <Option v-for="item in refTypeOptions" :value="item.value" :key="item.value">{{ item.label }}</Option>
            </Select>
            <Select
              ref="refSelect"
              v-model="form.refId"
              filterable
              :remote-method="() => {}"
              @on-query-change="remoteRefData"
              :loading="refLoading"
              clearable
              style="width:calc(60% - 105px);"
              @on-open-change="handleRefOpenChange"
            >
              <Option v-for="item in refOptions" :value="item.id" :key="item.id">{{
                `【${item.id}】${item.name}`
              }}</Option>
            </Select>
          </FormItem>
          <!--附件-->
          <FormItem :label="$t('tw_attach')">
            <UploadFile :id="requestId" :files="attachFiles" type="request"></UploadFile>
          </FormItem>
        </BaseHeaderTitle>
        <!--表单详情-->
        <BaseHeaderTitle :title="$t('tw_form_detail')">
          <div class="request-form">
            <template v-if="Object.keys(form.customForm.value).length > 0">
              <Divider style="margin: 0 0 30px 0" orientation="left">
                <span class="sub-header">{{ $t('tw_information_form') }}</span>
              </Divider>
              <CustomForm
                ref="customForm"
                :options="form.customForm.title"
                v-model="form.customForm.value"
                :requestId="requestId"
              ></CustomForm>
            </template>
            <Divider style="margin: 20px 0 30px 0" orientation="left">
              <span class="sub-header">{{ $t('tw_data_form') }}</span>
            </Divider>
            <!--选择目标对象-->
            <FormItem
              v-if="detail.associationWorkflow"
              :label="$t('tw_choose_object')"
              :label-width="lang === 'zh-CN' ? 110 : 150"
              required
            >
              <Select v-model="form.rootEntityId" clearable filterable style="width:300px;">
                <Option v-for="item in rootEntityOptions" :value="item.guid" :key="item.guid">{{
                  item.key_name
                }}</Option>
              </Select>
            </FormItem>
            <!--数据表单-->
            <EntityTable
              ref="entityTable"
              :data="requestData"
              :originRequestData="originRequestData"
              :requestId="requestId"
              isAdd
              autoAddRow
            ></EntityTable>
            <div v-if="noRequestForm" class="no-data">{{ $t('tw_no_formConfig') }}</div>
          </div>
        </BaseHeaderTitle>
        <!--审批流程-->
        <BaseHeaderTitle v-if="approvalList.length > 0" :title="$t('tw_approval_step')">
          <div class="step-wrap">
            <div v-for="(i, index) in approvalList" :key="index" class="step-item">
              <div class="step-item-left">
                <div class="circle">{{ index + 1 }}</div>
                <div v-if="index + 1 !== approvalList.length" class="line" />
              </div>
              <div class="step-item-content">
                <div class="title">
                  {{ i.name }}
                  <Tag color="default">{{ approvalTypeName[i.handleMode] }}</Tag>
                  <span style="color:#808695;">{{ i.description }}</span>
                  <Icon size="24" color="#00CB91" type="md-time" />
                  <span style="color:#00CB91;">{{ `${i.expireDay}${$t('day')}` }}</span>
                </div>
                <!--单人自定义、协同、并行-->
                <div v-if="['custom', 'any', 'all'].includes(i.handleMode)" class="step-background">
                  <div v-for="(j, idx) in i.handleTemplates" :key="idx" class="form-item">
                    <!--审批角色设置(模板指定/提交人指定)-->
                    <FormItem label=" " required :label-width="15">
                      <Select
                        v-if="['custom', 'template'].includes(j.assign)"
                        v-model="j.role"
                        filterable
                        :disabled="j.assign === 'template' ? true : false"
                        :placeholder="$t('tw_handleRole_placeholder')"
                        style="width:300px;"
                        @on-change="j.handler = ''"
                      >
                        <Option v-for="i in userRoleList" :key="i.id" :value="i.id">{{ i.displayName }}</Option>
                      </Select>
                    </FormItem>
                    <!--审批人设置(模板指定/建议/提交人指定/建议/组内系统分配/组内主动认领)-->
                    <FormItem label=" " required :label-width="15" style="margin-left:20px;">
                      <Input
                        v-if="['template', 'template_suggest'].includes(j.handlerType)"
                        v-model="j.handler"
                        disabled
                        :placeholder="$t('tw_handler_placeholder')"
                        style="width:300px;"
                      />
                      <Select
                        v-else-if="['custom', 'custom_suggest'].includes(j.handlerType)"
                        v-model="j.handler"
                        filterable
                        :placeholder="$t('tw_handler_placeholder')"
                        style="width:300px;"
                        @on-open-change="getHandlerByRole(i, j)"
                      >
                        <Option v-for="i in j.handlerList" :key="i.id" :value="i.id">{{ i.displayName }}</Option>
                      </Select>
                      <!--组内系统分配-->
                      <Input
                        v-else-if="j.handlerType === 'system'"
                        :value="$t('tw_group_assign')"
                        disabled
                        style="width:300px;"
                      />
                      <!--组内主动认领-->
                      <Input
                        v-else-if="j.handlerType === 'claim'"
                        :value="$t('tw_group_claim')"
                        disabled
                        style="width:300px;"
                      />
                    </FormItem>
                  </div>
                </div>
                <!--提交人角色管理员-->
                <div v-else-if="i.handleMode === 'admin'" class="step-background">
                  <div class="form-item">
                    <FormItem label=" " required :label-width="15">
                      <Select
                        v-model="i.handleTemplates[0].role"
                        disabled
                        :placeholder="$t('tw_handleRole_placeholder')"
                        style="width:300px;"
                      >
                        <Option v-for="i in userRoleList" :key="i.id" :value="i.id">{{ i.displayName }}</Option>
                      </Select>
                    </FormItem>
                    <FormItem label=" " required :label-width="15" style="margin-left:20px;">
                      <Input
                        v-model="i.handleTemplates[0].handler"
                        disabled
                        :placeholder="$t('tw_handler_placeholder')"
                        style="width:300px;"
                      />
                    </FormItem>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </BaseHeaderTitle>
        <!--任务流程-->
        <BaseHeaderTitle v-if="taskList.length > 0" :title="$t('tw_task_step')">
          <div class="step-wrap">
            <div v-for="(i, index) in taskList" :key="index" class="step-item">
              <div class="step-item-left">
                <div class="circle">{{ index + 1 }}</div>
                <div v-if="index + 1 !== taskList.length" class="line" />
              </div>
              <div class="step-item-content">
                <div class="title">
                  {{ i.name }}
                  <Tag v-if="approvalTypeName[i.handleMode]" color="default">{{ approvalTypeName[i.handleMode] }}</Tag>
                  <Tag :color="i.nodeDefId ? 'gold' : 'purple'">{{
                    i.nodeDefId ? $t('tw_workflow_task') : $t('tw_custom_task')
                  }}</Tag>
                  <span style="color:#808695;">{{ i.description }}</span>
                  <Icon size="24" color="#00CB91" type="md-time" />
                  <span style="color:#00CB91;">{{ `${i.expireDay}${$t('day')}` }}</span>
                </div>
                <!--单人自定义-->
                <div v-if="['custom', 'any', 'all'].includes(i.handleMode)" class="step-background">
                  <div v-for="(j, idx) in i.handleTemplates" :key="idx" class="form-item">
                    <!--审批角色设置-->
                    <FormItem label=" " required :label-width="15">
                      <Select
                        v-if="['custom', 'template'].includes(j.assign)"
                        v-model="j.role"
                        filterable
                        :disabled="j.assign === 'template' ? true : false"
                        :placeholder="$t('tw_handleRole_placeholder')"
                        style="width:300px;"
                        @on-change="j.handler = ''"
                      >
                        <Option v-for="i in userRoleList" :key="i.id" :value="i.id">{{ i.displayName }}</Option>
                      </Select>
                    </FormItem>
                    <!--审批人设置(模板指定/建议/提交人指定/建议/组内系统分配/组内主动认领)-->
                    <FormItem label=" " required :label-width="15" style="margin-left:20px;">
                      <Input
                        v-if="['template', 'template_suggest'].includes(j.handlerType)"
                        v-model="j.handler"
                        disabled
                        :placeholder="$t('tw_handler_placeholder')"
                        style="width:300px;"
                      />
                      <Select
                        v-else-if="['custom', 'custom_suggest'].includes(j.handlerType)"
                        v-model="j.handler"
                        filterable
                        :placeholder="$t('tw_handler_placeholder')"
                        style="width:300px;"
                        @on-open-change="getHandlerByRole(i, j)"
                      >
                        <Option v-for="i in j.handlerList" :key="i.id" :value="i.id">{{ i.displayName }}</Option>
                      </Select>
                      <!--组内系统分配-->
                      <Input
                        v-else-if="j.handlerType === 'system'"
                        :value="$t('tw_group_assign')"
                        disabled
                        style="width:300px;"
                      />
                      <!--组内主动认领-->
                      <Input
                        v-else-if="j.handlerType === 'claim'"
                        :value="$t('tw_group_claim')"
                        disabled
                        style="width:300px;"
                      />
                    </FormItem>
                  </div>
                </div>
                <!--提交人角色管理员-->
                <div v-else-if="i.handleMode === 'admin'" class="step-background">
                  <div class="form-item">
                    <FormItem label=" " required :label-width="15">
                      <Select
                        v-model="i.handleTemplates[0].role"
                        disabled
                        :placeholder="$t('tw_handleRole_placeholder')"
                        style="width:300px;"
                      >
                        <Option v-for="i in userRoleList" :key="i.id" :value="i.id">{{ i.displayName }}</Option>
                      </Select>
                    </FormItem>
                    <FormItem label=" " required :label-width="15" style="margin-left:20px;">
                      <Input
                        v-model="i.handleTemplates[0].handler"
                        disabled
                        :placeholder="$t('tw_handler_placeholder')"
                        style="width:300px;"
                      />
                    </FormItem>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </BaseHeaderTitle>
      </Form>
    </div>
    <!--编排流程图-->
    <template v-if="detail.associationWorkflow">
      <div class="expand-btn" :style="{ right: flowVisible ? '445px' : '0px' }" @click="flowVisible = !flowVisible">
        <Icon v-if="flowVisible" type="ios-arrow-dropright-circle" :size="28" />
        <Icon v-else type="ios-arrow-dropleft-circle" :size="28" />
      </div>
      <!--使用v-show会导致流程图加载不出来-->
      <div v-if="flowVisible" class="flow-expand">
        <StaticFlow :requestTemplate="requestTemplate" :flowId="detail.procId"></StaticFlow>
      </div>
    </template>
  </div>
</template>

<script>
import HeaderTag from './header-tag.vue'
import StaticFlow from './flow/static-flow.vue'
import DynamicFlow from './flow/dynamic-flow.vue'
import EntityTable from './entity-table.vue'
import DataBind from './data-bind.vue'
import UploadFile from './upload.vue'
import CustomForm from './custom-form.vue'
import BaseProgress from './base-progress.vue'
import { debounce, deepClone } from '@/pages/util'
import { requiredCheck, noChooseCheck, approvalCheck } from '../util'
import {
  getCreateInfo,
  getPublishInfo,
  getRootEntity,
  getEntityData,
  savePublishData,
  updateRequestStatus,
  getTaskConfig,
  getUserRoles,
  getHandlerRoles,
  getAdminUserByRole,
  getAssociateList
} from '@/api/server'
import dayjs from 'dayjs'
export default {
  components: {
    HeaderTag,
    StaticFlow,
    DynamicFlow,
    EntityTable,
    DataBind,
    UploadFile,
    CustomForm,
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
      jumpFrom: this.$route.query.jumpFrom || '', // 首页tabName
      type: this.$route.query.type, // 首页类型type
      requestTemplate: this.$route.query.requestTemplate,
      requestId: this.$route.query.requestId,
      role: this.$route.query.role,
      form: {
        name: '',
        description: '',
        expectTime: '',
        rootEntityId: '', // 目标对象
        customForm: {
          title: [],
          value: {}
        }, // 自定义表单
        data: [], // 数据表单
        approvalList: [], // 审批和任务列表
        refId: '', // 关联单Id
        refType: 1 // 关联单类型
      },
      refLoading: false,
      refOptions: [],
      detail: {},
      initExpectTime: '', // 记录初始期望完成时间
      expireDay: '',
      rootEntityOptions: [],
      requestData: [], // 数据表单
      originRequestData: [], // 数据表单初始值
      attachFiles: [], // 上传附件
      flowVisible: false,
      approvalList: [], // 审批流程
      taskList: [], // 任务流程
      approvalTypeName: {
        custom: this.$t('tw_onlyOne'), // 单人
        any: this.$t('tw_anyWidth'), // 协同
        all: this.$t('tw_allWidth'), // 并行
        admin: this.$t('tw_roleAdmin'), // 提交人角色管理员
        auto: this.$t('tw_autoWith') // 自动通过
      },
      userRoleList: [], // 用户角色列表
      noRequestForm: false, // 请求表单为空标识
      lang: window.localStorage.getItem('lang') || 'zh-CN',
      refTypeOptions: [
        { value: 1, label: this.$t('tw_publish') },
        { value: 2, label: this.$t('tw_request') },
        { value: 3, label: this.$t('tw_question') },
        { value: 4, label: this.$t('tw_event') },
        { value: 5, label: this.$t('fork') }
      ]
    }
  },
  watch: {
    'form.rootEntityId' (val) {
      if (val) {
        this.requestData = [] // 清空数据，解决数据缓存下拉框不回显问题
        this.getEntityData()
      } else {
        this.noRequestForm = false
        this.requestData = []
      }
    }
  },
  async created () {
    // 有id调用详情接口，无id调用创建接口
    if (this.requestId) {
      await this.getRequestInfoNew()
    } else {
      await this.getCreateInfo()
    }
    // 获取审批和任务流程
    if (this.form.approvalList.length === 0) {
      this.getApprovalAndTaskList()
    }
    // 获取请求进度
    this.$nextTick(() => {
      this.$refs.progress.initData(this.requestId)
    })
  },
  methods: {
    handleRefTypeChange () {
      this.form.refId = ''
      this.$refs.refSelect.query = ''
    },
    handleRefOpenChange (flag) {
      if (flag) {
        this.refOptions = []
        this.remoteRefData(this.form.refId)
      }
    },
    // 获取关联单下拉列表
    remoteRefData: debounce(async function (query) {
      // const cur = dayjs().format('YYYY-MM-DD')
      // const pre = dayjs()
      //   .subtract(3, 'month')
      //   .format('YYYY-MM-DD')
      const params = {
        action: this.form.refType, // 所有
        // reportTimeStart: pre + ' 00:00:00',
        // reportTimeEnd: cur + ' 23:59:59',
        query: query || '',
        startIndex: 0,
        pageSize: 50
      }
      this.refLoading = true
      const { statusCode, data } = await getAssociateList(params)
      this.refLoading = false
      if (statusCode === 'OK') {
        this.refOptions = data.contents || []
      }
    }, 500),
    handleToHome () {
      if (this.$route.query.jumpFrom) {
        this.$router.push({
          path: `/taskman/workbench?tabName=${this.jumpFrom}&actionName=${this.actionName}&${
            this.jumpFrom === 'submit' ? 'rollback' : 'type'
          }=${this.type}&needCache=yes`
        })
      } else {
        const pathMap = {
          1: 'publishHistory',
          2: 'requestHistory',
          3: 'problemHistory',
          4: 'eventHistory',
          5: 'changeHistory'
        }
        this.$router.push({
          path: `/taskman/workbench/${pathMap[this.actionName]}?&needCache=yes`
        })
      }
    },
    async getCreateInfo () {
      const params = {
        requestTemplate: this.requestTemplate,
        role: this.$route.query.role
      }
      const { statusCode, data } = await getCreateInfo(params)
      if (statusCode === 'OK') {
        this.detail = data || {}
        this.requestId = data.id
        // 解决刷新页面，会一直创建请求的问题
        this.$router.replace({
          path: this.$router.path,
          query: {
            requestTemplate: this.$route.query.requestTemplate,
            requestId: this.requestId,
            role: this.$route.query.role
          }
        })
        this.form.name = (data.name && data.name.substr(0, 70)) || ''
        this.form.expectTime = data.expectTime || ''
        this.form.customForm = {
          title: data.customForm.title || [],
          value: data.customForm.value || {}
        }
        this.form.customForm.title.forEach(item => {
          // 默认清空标志为false,赋值默认值
          if (item.defaultClear === 'no') {
            this.$set(this.form.customForm.value, item.name, item.defaultValue || '')
          } else {
            this.$set(this.form.customForm.value, item.name, '')
          }
        })
        this.expireDay = data.expireDay
        this.initExpectTime = this.form.expectTime
        // 模板未关联编排
        if (!this.detail.associationWorkflow) {
          this.getEntityData()
        } else {
          // 获取目标对象
          this.getEntity()
        }
      }
    },
    // 获取请求详情
    async getRequestInfoNew () {
      const params = {
        params: {
          requestId: this.requestId,
          taskId: ''
        }
      }
      const { statusCode, data } = await getPublishInfo(params)
      if (statusCode === 'OK') {
        this.detail = data.request || {}
        const { name, description, rootEntityId, expireDay, customForm, attachFiles, refId, refType } =
          data.request || {}
        this.form.name = (name && name.substr(0, 70)) || ''
        this.form.description = description
        this.form.rootEntityId = rootEntityId
        this.form.refId = refId
        this.form.refType = refType || 1
        this.attachFiles = attachFiles
        this.role = data.request.role
        // 初始化customForm
        this.form.customForm = {
          title: customForm.title || [],
          value: customForm.value || {}
        }
        this.form.customForm.title.forEach(item => {
          if (!this.form.customForm.value.hasOwnProperty(item.name)) {
            // 默认清空标志为false,赋值默认值
            if (item.defaultClear === 'no') {
              this.$set(this.form.customForm.value, item.name, item.defaultValue)
            } else {
              this.$set(this.form.customForm.value, item.name, '')
            }
          }
        })
        // 初始化审批和任务流程
        this.form.approvalList = data.approvalList || []
        const approvalList = data.approvalList || []
        this.approvalList = (approvalList && approvalList.filter(i => i.type === 'approve')) || []
        this.taskList = (approvalList && approvalList.filter(i => i.type === 'implement')) || []
        if (approvalList.length > 0) {
          // 角色下拉列表初始化
          const { statusCode, data } = await getUserRoles()
          if (statusCode === 'OK') {
            this.userRoleList = data
          }
          // 用户下拉列表初始化
          this.taskList.forEach(i => {
            i.handleTemplates &&
              i.handleTemplates.forEach(j => {
                if (j.handler) {
                  this.$set(j, 'handlerList', [{ displayName: j.handler, id: j.handler }])
                }
              })
          })
          // 用户下拉列表初始化
          this.approvalList.forEach(i => {
            i.handleTemplates &&
              i.handleTemplates.forEach(j => {
                if (j.handler) {
                  this.$set(j, 'handlerList', [{ displayName: j.handler, id: j.handler }])
                }
              })
          })
        }
        // 初始化期望完成时间
        this.expireDay = expireDay
        this.form.expectTime = dayjs()
          .add(this.expireDay || 0, 'day')
          .format('YYYY-MM-DD HH:mm:ss')
        this.initExpectTime = this.form.expectTime
        // 模板未关联编排
        if (!this.detail.associationWorkflow) {
          this.getEntityData()
        } else {
          // 获取目标对象
          this.getEntity()
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
        this.requestData = data.data || []
        this.originRequestData = deepClone(this.requestData)
        // 新建操作，默认值赋值逻辑
        if (this.jumpFrom !== 'draft') {
          this.requestData.forEach(i => {
            i.value.forEach(v => {
              i.title.forEach(t => {
                // 默认清空标志为false, 且初始值为空，赋值默认值
                if (t.defaultClear === 'no' && !v.entityData[t.name]) {
                  v.entityData[t.name] = t.defaultValue || ''
                }
                // if (t.defaultClear === 'yes' && !Array.isArray(v.entityData[t.name])) {
                //   v.entityData[t.name] = ''
                // }
              })
            })
          })
        }
        if (this.requestData.length === 0) {
          this.noRequestForm = true
        } else {
          this.noRequestForm = false
        }
      }
    },
    // 获取审批流程列表
    getApprovalConfigList () {
      return new Promise(async resolve => {
        // 任务类型: check 定版, approve 审批, implement 执行类型, confirm 请求确认
        const { statusCode, data } = await getTaskConfig(this.requestTemplate, 'approve')
        if (statusCode === 'OK') {
          resolve(data || [])
        } else {
          resolve([])
        }
      })
    },
    // 获取任务流程列表
    getTaskConfigList () {
      return new Promise(async resolve => {
        const { statusCode, data } = await getTaskConfig(this.requestTemplate, 'implement')
        if (statusCode === 'OK') {
          resolve(data || [])
        } else {
          resolve([])
        }
      })
    },
    getApprovalAndTaskList () {
      Promise.all([this.getApprovalConfigList(), this.getTaskConfigList()]).then(async response => {
        // 获取所有角色列表
        const { statusCode, data } = await getUserRoles()
        if (statusCode === 'OK') {
          this.userRoleList = data
        }
        for (let i of response[0]) {
          // 获取提交人角色管理员
          if (i.handleMode === 'admin') {
            const { statusCode, data } = await getAdminUserByRole(this.role)
            const obj = {
              role: this.role,
              handler: ''
            }
            if (statusCode === 'OK') {
              obj.handler = (data && data[0]) || ''
            }
            this.$set(i, 'handleTemplates', [obj])
          }
        }
        for (let i of response[1]) {
          // 获取提交人角色管理员
          if (i.handleMode === 'admin') {
            const { statusCode, data } = await getAdminUserByRole(this.role)
            const obj = {
              role: this.role,
              handler: ''
            }
            if (statusCode === 'OK') {
              obj.handler = (data && data[0]) || ''
            }
            this.$set(i, 'handleTemplates', [obj])
          }
        }
        this.approvalList = response[0]
        this.taskList = response[1]
      })
    },
    // 获取角色对应的处理人
    async getHandlerByRole (group, item) {
      const params = {
        params: {
          roles: item.role
        }
      }
      const { statusCode, data } = await getHandlerRoles(params)
      if (statusCode === 'OK') {
        let handlerList = data.map(d => {
          return {
            displayName: d,
            id: d
          }
        })
        // 同组下相同角色不能选择同一个处理人
        group.handleTemplates.map(obj => {
          if (obj.handler && item.handler !== obj.handler && item.role === obj.role) {
            handlerList = handlerList.filter(i => i.id !== obj.handler)
          }
        })
        this.$set(item, 'handlerList', handlerList)
      }
    },
    // 保存草稿
    async handleDraft (noJump) {
      // 请求表单
      const requestData = deepClone((this.$refs.entityTable && this.$refs.entityTable.requestData) || [])
      this.form.data =
        requestData.map(item => {
          let refKeys = [] // 引用类型
          let sensitiveKeys = [] // 敏感字段类型
          item.title.forEach(t => {
            if (t.elementType === 'select' || t.elementType === 'wecmdbEntity') {
              refKeys.push(t.name)
            }
            if (t.cmdbAttr) {
              const cmdbAttr = JSON.parse(t.cmdbAttr)
              if (cmdbAttr.sensitive === 'yes') {
                sensitiveKeys.push(t.name)
              }
            }
          })
          if (Array.isArray(item.value)) {
            item.value.forEach(v => {
              refKeys.forEach(ref => {
                delete v.entityData[ref + 'Options']
              })
              // 前端添加一行的数据，删除相应属性
              if (v.addFlag) {
                delete v.addFlag
              }
              // 删除表单隐藏属性, 并清空值
              for (const key in v.entityData) {
                if (v.entityData[key + 'Hidden']) {
                  v.entityData[key] = ''
                }
                delete v.entityData[key + 'Hidden']
              }
            })
          }
          // 敏感字段值不变, 删除属性，不传给后台
          const originData = this.originRequestData.find(i => i.entity === item.entity || i.itemGroup === item.itemGroup)
          sensitiveKeys.forEach(key => {
            for (let origin of (originData.value || [])) {
              for (let current of (item.value || [])) {
                if (origin.id === current.id && origin.entityData[key] === current.entityData[key]) {
                  // delete current.entityData[key]
                }
              }
            }
          })
          return item
        }) || []
      // 名称必填校验
      if (!this.form.name) {
        this.$Message.warning(this.$t('request_name') + this.$t('can_not_be_empty'))
        return
      }
      // 操作目标对象必填校验
      if (!this.form.rootEntityId && this.detail.associationWorkflow) {
        this.$Message.warning(this.$t('root_entity') + this.$t('can_not_be_empty'))
        return
      }
      // 信息表单必填校验
      if (!this.customFormValid()) {
        return this.$Message.warning(this.$t('tw_infoForm_valid'))
      }
      // 数据表单必填项-校验提示
      if (!requiredCheck(this.form.data, this.$refs.entityTable)) {
        const tabName = this.$refs.entityTable.activeTab
        return this.$Message.warning(`【${tabName}】${this.$t('required_tip')}`)
      }
      // 数据表单至少勾选一条数据校验
      if (!noChooseCheck(this.form.data, this.$refs.entityTable)) {
        const tabName = this.$refs.entityTable.activeTab
        return this.$Message.warning(`【${tabName}】${this.$t('tw_table_noChoose_tips')}`)
      }
      // 审批流程角色和用户必填校验
      if (!approvalCheck(this.approvalList)) {
        return this.$Message.warning(this.$t('tw_approvalStep_valid'))
      }
      // 任务流程角色和用户必填校验
      if (!approvalCheck(this.taskList)) {
        return this.$Message.warning(this.$t('tw_taskStep_valid'))
      }
      // 信息表单
      const customTitles = (this.$refs.customForm && this.$refs.customForm.formOptions) || []
      customTitles.forEach(t => {
        if (t.hidden) {
          // 清空隐藏表单的值
          this.form.customForm.value[t.name] = ''
        }
      })
      // 审批列表
      this.form.approvalList = [...this.approvalList, ...this.taskList]
      // 发布目标对象名称
      const item = this.rootEntityOptions.find(item => item.guid === this.form.rootEntityId) || {}
      this.form.entityName = item.key_name
      // 如果期望时间没有手动更改，提交时自动获取最新的
      if (this.initExpectTime === this.form.expectTime) {
        this.form.expectTime = dayjs()
          .add(this.expireDay || 0, 'day')
          .format('YYYY-MM-DD HH:mm:ss')
        this.initExpectTime = this.form.expectTime
      }
      const { statusCode } = await savePublishData(this.requestId, this.form)
      if (statusCode === 'OK') {
        if (noJump) {
          return statusCode
        } else {
          this.$Notice.success({
            title: this.$t('successful'),
            desc: this.$t('successful')
          })
        }
      }
    },
    // 发布
    async handlePublish () {
      this.$Modal.confirm({
        title: this.$t('tw_confirm_commit'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const draftResult = await this.handleDraft(true)
          if (draftResult === 'OK') {
            const { statusCode } = await updateRequestStatus(this.requestId, 'Pending')
            if (statusCode === 'OK') {
              this.$Notice.success({
                title: this.$t('successful'),
                desc: this.$t('successful')
              })
              this.$router.push({ path: `/taskman/workbench?tabName=submit&actionName=${this.actionName}` })
            }
          }
        },
        onCancel: () => {}
      })
    },
    customFormValid () {
      let result = true
      let requiredName = []
      this.form.customForm.title.forEach(t => {
        if (t.required === 'yes') {
          requiredName.push(t.name)
        }
      })
      requiredName.forEach(key => {
        let val = this.form.customForm.value[key]
        if (Array.isArray(val)) {
          if (val.length === 0) {
            result = false
          }
        } else {
          if (val === '' || val === undefined) {
            result = false
          }
        }
      })
      return result
    }
  }
}
</script>

<style lang="scss" scoped>
.workbench-publish-create {
  position: relative;
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
      background: #5384ff;
    }
    .name {
      font-size: 16px;
      margin-left: 16px;
      display: flex;
      align-items: center;
    }
  }
  .w-header {
    padding-bottom: 10px;
    margin-bottom: 20px;
    border-bottom: 2px dashed #d7dadc;
    display: flex;
    &-progress {
      flex: 1;
    }
    &-btn {
      width: 220px;
      text-align: right;
    }
  }
  .content {
    display: flex;
    min-height: 500px;
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
  }
  .footer-btn {
    padding: 10px 0 30px 0;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  .step-wrap {
    .step-item {
      display: flex;
      height: auto;
      .step-item-left {
        width: 40px;
        display: flex;
        flex-direction: column;
        align-items: center;
        .circle {
          text-align: center;
          line-height: 20px;
          width: 20px;
          height: 20px;
          border-radius: 20px;
          background-color: #e1e9f0;
          color: #9da7b3;
          font-size: 12px;
        }
        .line {
          height: calc(100% - 20px);
          width: 1px;
          background-color: #e1e9f0;
        }
      }
      .step-item-content {
        flex: 1;
        padding-bottom: 20px;
        margin-top: -5px;
        .type-tag {
          padding: 5px 8px;
          height: 20px;
          line-height: 20px;
          color: #606975;
          background-color: #ebf0f5;
          border-radius: 4px;
          font-size: 12px;
        }
        .time {
          font-weight: bold;
          margin-left: 2px;
        }
        .form-item {
          padding-top: 10px;
          display: flex;
        }
      }
    }
    .step-background {
      background: #f0faff;
      padding: 5px 20px;
      width: 690px;
      margin-top: 10px;
      border-radius: 4px;
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
.workbench-publish-create {
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
