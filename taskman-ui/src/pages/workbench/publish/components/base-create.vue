<template>
  <div class="workbench-publish-create">
    <Row class="back-header">
      <Icon size="22" type="md-arrow-back" class="icon" @click="handleToHome" />
      <span class="name">
        {{ `${detail.templateName || detail.requestTemplateName || ''}` }}
        <Tag size="medium">{{ detail.version || detail.templateVersion || '' }}</Tag>
      </span>
    </Row>
    <Row class="w-header">
      <Col span="19">
        <!--请求进度-->
        <BaseProgress ref="progress"></BaseProgress>
      </Col>
      <Col span="5" class="btn-group">
        <!--保存草稿-->
        <Button @click="handleDraft(false)" style="margin-right:10px;">{{ $t('tw_save_draft') }}</Button>
        <!--提交-->
        <Button type="primary" @click="handlePublish">{{ $t('tw_commit') }}</Button>
      </Col>
    </Row>
    <div class="content">
      <Form :model="form" label-position="right" :label-width="120" style="width:100%;">
        <!--请求信息-->
        <HeaderTitle :title="$t('tw_request_title')">
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
          <!--附件-->
          <FormItem :label="$t('tw_attach')">
            <UploadFile :id="requestId" :files="attachFiles" type="request" :formDisable="formDisable"></UploadFile>
          </FormItem>
          <!--自定义信息表单-->
          <CustomForm
            v-model="form.customForm.value"
            :options="form.customForm.title"
            :requestId="requestId"
          ></CustomForm>
        </HeaderTitle>
        <!--请求表单-->
        <HeaderTitle title="请求表单">
          <!--选择目标对象-->
          <FormItem v-if="detail.associationWorkflow" :label="$t('tw_choose_object')" required>
            <Select v-model="form.rootEntityId" :disabled="formDisable" clearable filterable style="width:300px;">
              <Option v-for="item in rootEntityOptions" :value="item.guid" :key="item.guid">{{ item.key_name }}</Option>
            </Select>
          </FormItem>
          <EntityTable
            ref="entityTable"
            :data="requestData"
            :requestId="requestId"
            :type="actionName"
            isAdd
            autoAddRow
            style="width:calc(100% - 20px);margin-left:16px;"
          ></EntityTable>
          <div v-if="noRequestForm" class="no-data">暂未配置表单</div>
        </HeaderTitle>
        <!--审批流程-->
        <HeaderTitle v-if="approvalList.length > 0" title="审批流程">
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
                  <Icon size="24" color="#19be6b" type="md-time" />
                  <span style="color:#19be6b;">{{ i.expireDay }}天</span>
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
                        placeholder="请选择处理角色"
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
                        placeholder="请选择处理人"
                        style="width:300px;"
                      />
                      <Select
                        v-else-if="['custom', 'custom_suggest'].includes(j.handlerType)"
                        v-model="j.handler"
                        filterable
                        placeholder="请选择处理人"
                        style="width:300px;"
                        @on-open-change="getHandlerByRole(i, j)"
                      >
                        <Option v-for="i in j.handlerList" :key="i.id" :value="i.id">{{ i.displayName }}</Option>
                      </Select>
                      <Input
                        v-else-if="j.handlerType === 'system'"
                        value="组内系统分配"
                        disabled
                        style="width:300px;"
                      />
                      <Input v-else-if="j.handlerType === 'claim'" value="组内主动认领" disabled style="width:300px;" />
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
                        placeholder="请选择处理角色"
                        style="width:300px;"
                      >
                        <Option v-for="i in userRoleList" :key="i.id" :value="i.id">{{ i.displayName }}</Option>
                      </Select>
                    </FormItem>
                    <FormItem label=" " required :label-width="15" style="margin-left:20px;">
                      <Input
                        v-model="i.handleTemplates[0].handler"
                        disabled
                        placeholder="请选择处理人"
                        style="width:300px;"
                      />
                    </FormItem>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </HeaderTitle>
        <!--任务流程-->
        <HeaderTitle v-if="taskList.length > 0" title="任务流程">
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
                  <Tag :color="i.nodeDefId ? 'gold' : 'purple'">{{ i.nodeDefId ? '编排人工任务' : '自定义任务' }}</Tag>
                  <span style="color:#808695;">{{ i.description }}</span>
                  <Icon size="24" color="#19be6b" type="md-time" />
                  <span style="color:#19be6b;">{{ i.expireDay }}天</span>
                </div>
                <!--单人自定义-->
                <div v-if="i.handleMode === 'custom'" class="step-background">
                  <div v-for="(j, idx) in i.handleTemplates" :key="idx" class="form-item">
                    <!--审批角色设置-->
                    <FormItem label=" " required :label-width="15">
                      <Select
                        v-if="['custom', 'template'].includes(j.assign)"
                        v-model="j.role"
                        filterable
                        :disabled="j.assign === 'template' ? true : false"
                        placeholder="请选择处理角色"
                        style="width:300px;"
                        @on-change="j.handler = ''"
                      >
                        <Option v-for="i in userRoleList" :key="i.id" :value="i.id">{{ i.displayName }}</Option>
                      </Select>
                    </FormItem>
                    <!--审批人设置-->
                    <FormItem label=" " required :label-width="15" style="margin-left:20px;">
                      <Input
                        v-if="['template', 'template_suggest'].includes(j.handlerType)"
                        v-model="j.handler"
                        disabled
                        placeholder="请选择处理人"
                        style="width:300px;"
                      />
                      <Select
                        v-else-if="['custom', 'custom_suggest'].includes(j.handlerType)"
                        v-model="j.handler"
                        filterable
                        placeholder="请选择处理人"
                        style="width:300px;"
                        @on-open-change="getHandlerByRole(i, j)"
                      >
                        <Option v-for="i in j.handlerList" :key="i.id" :value="i.id">{{ i.displayName }}</Option>
                      </Select>
                      <Input
                        v-else-if="j.handlerType === 'system'"
                        value="组内系统分配"
                        disabled
                        style="width:300px;"
                      />
                      <Input v-else-if="j.handlerType === 'claim'" value="组内主动认领" disabled style="width:300px;" />
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
                        placeholder="请选择处理角色"
                        style="width:300px;"
                      >
                        <Option v-for="i in userRoleList" :key="i.id" :value="i.id">{{ i.displayName }}</Option>
                      </Select>
                    </FormItem>
                    <FormItem label=" " required :label-width="15" style="margin-left:20px;">
                      <Input
                        v-model="i.handleTemplates[0].handler"
                        disabled
                        placeholder="请选择处理人"
                        style="width:300px;"
                      />
                    </FormItem>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </HeaderTitle>
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
        <StaticFlow :requestTemplate="requestTemplate"></StaticFlow>
      </div>
    </template>
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
import CustomForm from '../../components/custom-form.vue'
import BaseProgress from './progress.vue'
import { deepClone } from '@/pages/util/index'
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
  getAdminUserByRole
} from '@/api/server'
import dayjs from 'dayjs'
export default {
  components: {
    HeaderTitle,
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
      formDisable: this.$route.query.isCheck === 'Y', // 查看标志
      jumpFrom: this.$route.query.jumpFrom, // 首页tabName
      type: this.$route.query.type, // 首页类型type
      requestTemplate: this.$route.query.requestTemplate,
      requestId: this.$route.query.requestId,
      role: this.$route.query.role,
      // procDefId: '',
      // procDefKey: '',
      form: {
        name: '',
        description: '',
        expectTime: '',
        rootEntityId: '', // 目标对象
        customForm: {
          title: [],
          value: {}
        }, // 自定义表单
        data: [],
        approvalList: [] // 审批列表
      },
      detail: {},
      initExpectTime: '', // 记录初始期望完成时间
      expireDay: '',
      rootEntityOptions: [],
      requestData: [], // 发布目标对象表格数据
      attachFiles: [], // 上传附件
      flowVisible: false,
      approvalList: [], // 审批流程
      taskList: [], // 任务流程
      approvalTypeName: {
        custom: '单人',
        any: '协同',
        all: '并行',
        admin: '提交人角色管理员',
        auto: '自动通过'
      },
      userRoleList: [], // 用户角色列表
      noRequestForm: false // 请求表单为空标识
    }
  },
  watch: {
    'form.rootEntityId' (val) {
      if (val) {
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
    handleToHome () {
      this.$router.push({
        path: `/taskman/workbench?tabName=${this.jumpFrom}&actionName=${this.actionName}&${
          this.jumpFrom === 'submit' ? 'rollback' : 'type'
        }=${this.type}`
      })
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
        this.form.name = (data.name && data.name.substr(0, 70)) || ''
        this.form.expectTime = data.expectTime || ''
        this.form.customForm = {
          title: data.customForm.title || [],
          value: data.customForm.value || {}
        }
        this.form.customForm.title.forEach(item => {
          this.form.customForm.value[item.name] = ''
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
      const { statusCode, data } = await getPublishInfo(this.requestId)
      if (statusCode === 'OK') {
        this.detail = data.request || {}
        const { name, description, rootEntityId, expireDay, customForm, attachFiles } = data.request || {}
        this.form.name = (name && name.substr(0, 70)) || ''
        this.form.description = description
        this.form.rootEntityId = rootEntityId
        this.attachFiles = attachFiles
        this.role = data.request.role
        // 初始化customForm
        this.form.customForm = {
          title: customForm.title || [],
          value: customForm.value || {}
        }
        this.form.customForm.title.forEach(item => {
          if (!this.form.customForm.value.hasOwnProperty(item.name)) {
            this.form.customForm.value[item.name] = ''
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
        this.requestData = data.data || []
        // 新建操作，默认值赋值逻辑
        if (this.jumpFrom !== 'draft') {
          this.requestData.forEach(i => {
            i.value.forEach(v => {
              i.title.forEach(t => {
                // 默认清空标志为false, 且初始值为空，赋值默认值
                if (t.defaultClear === 'no' && !v.entityData[t.name]) {
                  v.entityData[t.name] = t.defaultValue || ''
                }
                if (t.defaultClear === 'yes' && !Array.isArray(v.entityData[t.name])) {
                  v.entityData[t.name] = ''
                }
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
          let refKeys = []
          item.title.forEach(t => {
            if (t.elementType === 'select' || t.elementType === 'wecmdbEntity') {
              refKeys.push(t.name)
            }
          })
          if (Array.isArray(item.value)) {
            // 删除多余的属性
            item.value.forEach(v => {
              refKeys.forEach(ref => {
                delete v.entityData[ref + 'Options']
              })
              // 前端添加一行数据，删除多余属性
              if (v.addFlag) {
                delete v.addFlag
                v.id = ''
                v.entityData._id = ''
              }
            })
          }
          return item
        }) || []
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
      // 请求名称必填校验
      if (!this.form.name) {
        this.$Message.warning(this.$t('request_name') + this.$t('can_not_be_empty'))
        return
      }
      // 操作目标对象必填校验
      if (!this.form.rootEntityId && this.detail.associationWorkflow) {
        this.$Message.warning(this.$t('root_entity') + this.$t('can_not_be_empty'))
        return
      }
      // 请求信息自定义表单必填校验
      if (!this.customFormValid()) {
        return this.$Message.warning('请求信息必填项为空')
      }
      // 请求表单必填项-校验提示
      if (!this.requiredCheck(this.form.data)) {
        const tabName = this.$refs.entityTable.activeTab
        return this.$Message.warning(`【${tabName}】${this.$t('required_tip')}`)
      }
      // 审批任务流程角色和用户必填校验
      if (!this.approvalCheck(this.approvalList)) {
        return this.$Message.warning('审批流程处理角色和处理人必填')
      }
      if (!this.approvalCheck(this.taskList)) {
        return this.$Message.warning('任务流程处理角色和处理人必填')
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
    // 请求表单数据必填项校验
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
    },
    approvalCheck (data) {
      let result = true
      data.forEach(i => {
        if (i.handleTemplates && i.handleTemplates.length > 0) {
          i.handleTemplates.forEach(j => {
            if (!j.role || (!j.handler && !['system', 'claim'].includes(j.handlerType))) {
              result = false
            }
          })
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
    display: flex;
    min-height: 500px;
    .no-data {
      padding-left: 20px;
      height: 60px;
      line-height: 60px;
      color: #515a6e;
    }
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
