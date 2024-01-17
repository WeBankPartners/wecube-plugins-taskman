<template>
  <div class="workbench-publish-create">
    <Row class="back-header">
      <Icon size="26" type="md-arrow-back" style="cursor:pointer" @click="$router.back()" />
      <span class="name">
        {{ `${templateName || ''}` }}<Tag size="medium">{{ version }}</Tag>
      </span>
    </Row>
    <Row class="w-header">
      <Col span="19" class="steps">
        <!--请求进度-->
        <span class="title">{{ $t('tw_request_progress') }}：</span>
        <Steps :current="0" style="max-width:600px;">
          <Step v-for="(i, index) in progressList" :key="index" :content="i.name">
            <template #icon>
              <Icon style="font-weight:bold" size="24" :type="i.icon" :color="i.color" />
            </template>
            <div class="role" slot="content">
              <Tooltip :content="i.name">
                <div class="word-eclipse">{{ i.name }}</div>
              </Tooltip>
              <span>{{ i.handler }}</span>
            </div>
          </Step>
        </Steps>
      </Col>
      <Col span="5" class="btn-group">
        <!--保存草稿-->
        <Button :disabled="!requestData.length" @click="handleDraft(false)" style="margin-right:10px;">{{
          $t('tw_save_draft')
        }}</Button>
        <!--提交-->
        <Button :disabled="!requestData.length" type="primary" @click="handlePublish">{{ $t('tw_commit') }}</Button>
      </Col>
    </Row>
    <div style="display:flex;" class="content">
      <div style="width:calc(100% - 420px)" class="split-line">
        <Form :model="form" label-position="right" :label-width="120">
          <template>
            <!--请求信息-->
            <HeaderTitle :title="$t('tw_request_title')">
              <!--请求名-->
              <FormItem :label="$t('request_name')" required>
                <Input
                  v-model="form.name"
                  :maxlength="70"
                  show-word-limit
                  :placeholder="$t('request_name')"
                  style="width:100%;"
                />
              </FormItem>
              <!--请求描述-->
              <FormItem :label="$t('tw_publish_des')">
                <Input
                  v-model="form.description"
                  type="textarea"
                  :maxlength="200"
                  show-word-limit
                  :placeholder="$t('tw_publish_des')"
                  style="width:100%;"
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
              <!--附件-->
              <FormItem :label="$t('tw_attach')">
                <UploadFile :id="requestId" :files="attachFiles" type="request" :formDisable="formDisable"></UploadFile>
              </FormItem>
            </HeaderTitle>
            <!--目标对象-->
            <HeaderTitle :title="$t('tw_publish_object')">
              <!--选择目标对象-->
              <FormItem :label="$t('tw_choose_object')" required>
                <Select v-model="form.rootEntityId" :disabled="formDisable" clearable filterable style="width:300px;">
                  <Option v-for="item in rootEntityOptions" :value="item.guid" :key="item.guid">{{
                    item.key_name
                  }}</Option>
                </Select>
              </FormItem>
              <FormItem v-if="requestData.length" :label="$t('tw_selected')">
                <EntityTable ref="entityTable" :data="requestData" :requestId="requestId" :isAdd="true"></EntityTable>
              </FormItem>
            </HeaderTitle>
          </template>
        </Form>
      </div>
      <!--编排流程-->
      <div style="width:420px;padding-left:20px;">
        <StaticFlow :requestTemplate="requestTemplate"></StaticFlow>
      </div>
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
import { deepClone } from '@/pages/util/index'
import {
  getCreateInfo,
  getProgressInfo,
  getPublishInfo,
  getRootEntity,
  getEntityData,
  savePublishData,
  getRequestInfo,
  updateRequestStatus
} from '@/api/server'
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
    DataBind,
    UploadFile
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
      // actionName: this.$route.query.action, // 1发布,2请求(3问题,4事件,5变更)
      templateName: '',
      version: '', // 模板版本号
      formDisable: this.$route.query.isCheck === 'Y', // 查看标志
      jumpFrom: this.$route.query.jumpFrom, // 入口tab标记
      requestTemplate: this.$route.query.requestTemplate,
      requestId: this.$route.query.requestId,
      // procDefId: '',
      // procDefKey: '',
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
      attachFiles: [] // 上传附件
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
  async created () {
    // 新建发布无requestid, 调用创建接口
    if (this.requestId) {
      this.getPublishInfo()
      // 获取附件和操作单元
      this.getRequestInfo()
    } else {
      await this.getCreateInfo()
    }
    // 获取请求进度
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
        this.form.name = (data.name && data.name.substr(0, 70)) || ''
        // this.form.expectTime = dayjs()
        //   .add(data.expireDay || 0, 'day')
        //   .format('YYYY-MM-DD HH:mm:ss')
        this.form.expectTime = data.expectTime || ''
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
              item.name = this.$t('tw_commit_request') // 提交请求
              break
            case 'requestPending':
              item.name = this.$t('tw_request_pending') // 请求定版
              break
            case 'requestComplete':
              item.name = this.$t('tw_request_complete') // 请求完成
              break
            case 'autoExit':
              item.name = this.$t('status_faulted') // 自动退出
              this.errorNode = item.node
              break
            case 'internallyTerminated':
              item.name = this.$t('status_termination') // 手动终止
              this.errorNode = item.node
              break
            default:
              item.name = item.node
              break
          }
          if (item.handler === 'autoNode') {
            item.handler = this.$t('tw_auto_tag') // 自动节点
            this.errorNode = item.name
          }
        })
      }
    },
    async getRequestInfo () {
      const { statusCode, data } = await getRequestInfo(this.requestId)
      if (statusCode === 'OK') {
        this.attachFiles = data.attachFiles
        this.form.rootEntityId = data.cache
      }
    },
    // 获取发布信息
    async getPublishInfo () {
      const { statusCode, data } = await getPublishInfo(this.requestId)
      if (statusCode === 'OK') {
        const { name, description, expectTime } = data.request || {}
        this.form = Object.assign({}, this.form, {
          name,
          description,
          expectTime
        })
        this.templateName = data.request.templateName
        this.version = data.request.version
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
        // 新建和重新发起（不包括去发起），默认值逻辑
        if (this.jumpFrom !== 'my_drafts') {
          requestData.forEach(i => {
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
        this.requestData = requestData
      }
    },
    // 保存草稿
    async handleDraft (noJump) {
      if (!this.form.name) {
        this.$Message.warning(this.$t('request_name') + this.$t('can_not_be_empty'))
        return
      }
      if (!this.form.rootEntityId) {
        this.$Message.warning(this.$t('root_entity') + this.$t('can_not_be_empty'))
        return
      }
      // 提取表格勾选的数据
      const requestData = deepClone(this.$refs.entityTable && this.$refs.entityTable.requestData)
      this.form.data =
        requestData.map(item => {
          let refKeys = []
          item.title.forEach(t => {
            if (t.elementType === 'select') {
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
      // 获取发布目标对象entityName
      const item = this.rootEntityOptions.find(item => item.guid === this.form.rootEntityId) || {}
      this.form.entityName = item.key_name
      // 必填项校验提示
      if (!this.requiredCheck(this.form.data)) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('required_tip')
        })
      }
      // 表格至少勾选一条数据校验
      const tabName = this.$refs.entityTable.activeTab
      if (!this.noChooseCheck(this.form.data)) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: `【${tabName}】${this.$t('tw_table_noChoose_tips')}`
        })
      }
      const { statusCode } = await savePublishData(this.requestId, this.form)
      if (statusCode === 'OK') {
        if (noJump) {
          return statusCode
        } else {
          // if (this.jumpFrom === 'my_submit') {
          //   const rollback = this.$route.query.rollback
          //   this.$router.push({
          //     path: `/taskman/workbench?tabName=submit&actionName=${this.actionName}&rollback=${rollback}`
          //   })
          // } else {
          //   this.$router.push({ path: `/taskman/workbench?tabName=draft&actionName=${this.actionName}` })
          // }
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
              this.$router.push({ path: `/taskman/workbench?tabName=submit&actionName=${this.actionName}` })
            }
          }
        },
        onCancel: () => {}
      })
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
    }
  }
}
</script>

<style lang="scss" scoped>
.workbench-publish-create {
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
        max-width: 180px;
        text-overflow: ellipsis;
        overflow: hidden;
        white-space: nowrap;
      }
    }
    .btn-group {
      display: flex;
      justify-content: flex-end;
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
.workbench-publish-create {
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
}
</style>
