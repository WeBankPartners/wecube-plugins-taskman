<!--当前处理-->
<template>
  <div class="workbench-current-handle">
    <!--定版-->
    <HeaderTitle v-if="detail.status === 'Pending'" :title="$t('tw_cur_handle')">
      <div class="sub-title" slot="sub-title">
        <Tag class="tag" :color="handleTypeColor[handleData.type]">{{ $t('tw_request_pending') }}</Tag>
      </div>
      <div class="step-item">
        <div class="step-item-left">
          <div class="circle">1</div>
          <div class="line" />
        </div>
        <div class="step-item-content">
          <div class="title">
            {{ $t('tw_pending_step1') }}
            <span>{{ $t('tw_pending_step1_tips') }}</span>
          </div>
          <div class="content"></div>
        </div>
      </div>
      <div class="step-item">
        <div class="step-item-left">
          <div class="circle">2</div>
        </div>
        <div class="step-item-content">
          <div class="title">
            {{ $t('tw_pending_step2') }}
            <span>{{ $t('tw_pending_step2_tips') }}</span>
          </div>
          <div class="content">
            <DataBind
              :isHandle="isHandle"
              :requestTemplate="requestTemplate"
              :requestId="requestId"
              :formDisable="detail.status !== 'Pending'"
              :actionName="actionName"
            ></DataBind>
          </div>
        </div>
      </div>
    </HeaderTitle>
    <!--审批和任务-->
    <HeaderTitle v-else-if="['InApproval', 'InProgress'].includes(detail.status)" :title="$t('tw_cur_handle')">
      <div class="sub-title" slot="sub-title">
        <Tag>{{ approvalTypeName[handleData.handleMode] || '' }}</Tag>
        <Tag style="margin-left:5px;" :color="handleTypeColor[handleData.type]">{{
          `${{
            approve: '审批',
            implement_custom: '自定义任务',
            implement_process: '编排任务'
          }[handleData.type] || '-'}：${handleData.name}`
        }}</Tag>
        <span class="description">{{ $t('description') }}：{{ handleData.description || '-' }}</span>
      </div>
      <div class="step-item">
        <div class="step-item-left">
          <div class="circle">1</div>
          <div class="line" />
        </div>
        <div class="step-item-content">
          <div class="title">
            {{ $t('tw_approval_step1') }}
            <span>{{ $t('tw_approval_step1_tips') }}</span>
          </div>
          <div class="content">
            <EntityTable ref="entityTable" :data="handleData.formData" :requestId="requestId"></EntityTable>
            <div v-if="handleData.formData && handleData.formData.length === 0" class="no-data">
              暂未配置表单
            </div>
          </div>
        </div>
      </div>
      <div class="step-item">
        <div class="step-item-left">
          <div class="circle">2</div>
        </div>
        <div class="step-item-content">
          <div class="title">
            {{ $t('tw_approval_step2') }}
            <span>{{ $t('tw_approval_step2_tips') }}</span>
          </div>
          <div class="content">
            <Form :label-width="80" label-position="left">
              <!--处理结果-审批类型-->
              <FormItem v-if="handleData.type === 'approve'" required label="操作">
                <Select v-model="taskForm.choseOption" @on-change="handleChoseOptionChange">
                  <Option v-for="(item, index) in approvalNextOptions" :value="item.value" :key="index">{{
                    item.label
                  }}</Option>
                </Select>
              </FormItem>
              <!--处理结果-编排任务-->
              <FormItem
                v-if="
                  handleData.type === 'implement_process' && handleData.nextOptions && handleData.nextOptions.length > 0
                "
                required
                label="操作"
              >
                <Select v-model="taskForm.choseOption">
                  <Option v-for="option in handleData.nextOptions" :value="option" :key="option">{{ option }}</Option>
                </Select>
              </FormItem>
              <!--完成状态(只有任务有)-->
              <FormItem v-if="['implement_custom', 'implement_process'].includes(handleData.type)" label="完成状态">
                <Select v-model="taskForm.handleStatus">
                  <Option v-for="(item, index) in taskStatusList" :value="item.value" :key="index">{{
                    item.label
                  }}</Option>
                </Select>
              </FormItem>
              <!--处理意见-->
              <FormItem
                :label="$t('process_comments')"
                :required="handleData.type === 'approve' && taskForm.choseOption === 'redraw'"
              >
                <Input v-model="taskForm.comment" type="textarea" :maxlength="200" show-word-limit />
              </FormItem>
              <FormItem :label="$t('tw_attach')">
                <UploadFile
                  :id="handleData.id"
                  :taskHandleId="taskHandleId"
                  :files="taskForm.attachFiles"
                  type="task"
                ></UploadFile>
              </FormItem>
            </Form>
            <div v-if="handleData.editable" style="text-align: center">
              <!-- <Button @click="saveTaskData" type="info">{{ $t('save') }}</Button> -->
              <Button :disabled="commitTaskDisabled" @click="commitTaskData" type="primary">{{
                $t('tw_commit')
              }}</Button>
            </div>
          </div>
        </div>
      </div>
    </HeaderTitle>
    <!--确认-->
    <HeaderTitle v-else-if="detail.status === 'Confirm'" :title="$t('tw_cur_handle')">
      <div class="sub-title" slot="sub-title">
        <Tag class="tag" :color="handleTypeColor[handleData.type]">{{ $t('tw_request_confirm') }}</Tag>
      </div>
      <div style="padding:20px 16px;">
        <Form :label-width="80" label-position="left">
          <FormItem label="状态">
            <Select v-model="confirmRequestForm.completeStatus" @on-change="confirmRequestForm.markTaskId = []">
              <Option v-for="(i, index) in completeStatusList" :value="i.value" :key="index">{{ i.label }}</Option>
            </Select>
          </FormItem>
          <template v-if="confirmRequestForm.completeStatus === 'uncompleted'">
            <FormItem label="未完成任务节点">
              <Select
                v-model="confirmRequestForm.markTaskId"
                multiple
                clearable
                @on-open-change="
                  flag => {
                    flag && geTaskTagList()
                  }
                "
              >
                <Option v-for="i in taskTagList" :value="i.id" :key="i.id">{{ i.name }}</Option>
              </Select>
            </FormItem>
          </template>
          <FormItem :label="$t('process_comments')">
            <Input v-model="confirmRequestForm.notes" type="textarea" :maxlength="200" show-word-limit />
          </FormItem>
          <FormItem :label="$t('tw_attach')">
            <UploadFile :id="handleData.id" :taskHandleId="taskHandleId" type="task"></UploadFile>
          </FormItem>
        </Form>
        <div style="text-align:center;">
          <Button @click="confirmRequest" type="primary">{{ $t('tw_commit') }}</Button>
        </div>
      </div>
    </HeaderTitle>
  </div>
</template>

<script>
import HeaderTitle from '../../components/header-title.vue'
import EntityTable from '../../components/entity-table.vue'
import DataBind from '../../components/data-bind.vue'
import UploadFile from '../../components/upload.vue'
import { deepClone } from '@/pages/util/index'
import { commitTaskData, saveTaskData, geTaskTagList, confirmRequest } from '@/api/server'
import { requiredCheck, noChooseCheck } from '../../util'
export default {
  components: {
    HeaderTitle,
    EntityTable,
    DataBind,
    UploadFile
  },
  props: {
    detail: {
      type: Object,
      default: () => {}
    },
    handleData: {
      type: Object,
      default: () => ({ taskHandleList: [] })
    },
    // 1发布,2请求(3问题,4事件,5变更)
    actionName: {
      type: String,
      default: '1'
    }
  },
  data () {
    return {
      isHandle: this.$route.query.isHandle === 'Y', // 处理标志
      requestTemplate: this.$route.query.requestTemplate, // 请求模板ID
      requestId: this.$route.query.requestId, // 请求ID
      taskHandleId: this.$route.query.taskHandleId,
      // 任务和审批表单
      taskForm: {
        comment: '', // 处理意见
        choseOption: '', // 处理结果
        handleStatus: '', // 处理状态
        attachFiles: []
      },
      // 请求确认表单
      confirmRequestForm: {
        markTaskId: [], // 关注任务ID
        completeStatus: 'complete', // 请求完成状态complete、uncompleted
        notes: '',
        attachFiles: []
      },
      taskTagList: [], // 任务节点列表
      completeStatusList: [
        {
          label: this.$t('tw_completed'),
          value: 'complete'
        },
        {
          label: '未完成/部分完成',
          value: 'uncompleted'
        }
      ],
      approvalNextOptions: [
        {
          label: '拒绝',
          value: 'deny'
        },
        {
          label: '同意',
          value: 'approve'
        },
        {
          label: '退回',
          value: 'redraw'
        }
      ],
      taskStatusList: [
        {
          label: this.$t('tw_completed'),
          value: 'complete'
        },
        {
          label: this.$t('tw_incomplete'),
          value: 'uncompleted'
        }
      ],
      approvalTypeName: {
        custom: '单人',
        any: '协同',
        all: '并行',
        admin: '提交人角色管理员',
        auto: '自动通过'
      },
      handleTypeColor: {
        check: '#ffa2d3',
        approve: '#2d8cf0',
        implement_process: '#cba43f',
        implement_custom: '#b886f8',
        confirm: '#19be6b'
      },
      hasRunOnce: false
    }
  },
  computed: {
    commitTaskDisabled () {
      const approveFlag = this.handleData.type === 'approve' && !this.taskForm.choseOption
      const processFlag =
        this.handleData.type === 'implement_process' &&
        this.handleData.nextOptions &&
        this.handleData.nextOptions.length > 0 &&
        !this.taskForm.choseOption
      const commentFlag =
        this.handleData.type === 'approve' && this.taskForm.choseOption === 'redraw' && !this.taskForm.comment
      if (approveFlag || processFlag || commentFlag) {
        return true
      } else {
        return false
      }
    }
  },
  watch: {
    handleData: {
      handler (val) {
        if (val && this.taskHandleId) {
          if (this.hasRunOnce) return
          this.hasRunOnce = true
          const list = val.taskHandleList || []
          list.forEach(item => {
            if (item.id === this.taskHandleId) {
              if (['InApproval', 'InProgress'].includes(this.detail.status)) {
                this.taskForm.comment = item.resultDesc
                // 通过createdTime===updatedTime判断首次编辑时，给默认值
                if (val.type === 'approve' && item.createdTime === item.updatedTime) {
                  this.taskForm.choseOption = 'approve'
                } else {
                  this.taskForm.choseOption = item.handleResult
                }
                if (item.createdTime === item.updatedTime) {
                  this.taskForm.handleStatus = 'complete'
                } else {
                  this.taskForm.handleStatus = item.handleStatus
                }
                this.taskForm.attachFiles = item.attachFiles
              }
            }
          })
        }
      },
      deep: true
    }
  },
  methods: {
    handleChoseOptionChange (val) {
      // 给退回操作默认处理意见
      this.taskForm.comment = ''
      if (val === 'redraw') {
        this.taskForm.comment = '退回'
      }
    },
    // 任务审批保存
    async saveTaskData () {
      // 提取表格勾选的数据
      const requestData = deepClone((this.$refs.entityTable && this.$refs.entityTable.requestData) || [])
      this.handleData.formData =
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
            })
          }
          return item
        }) || []
      // 必填项校验提示
      if (!requiredCheck(this.handleData.formData, this.$refs.entityTable)) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('required_tip')
        })
      }
      // 表格至少勾选一条数据校验
      if (!noChooseCheck(this.handleData.formData, this.$refs.entityTable)) {
        const tabName = this.$refs.entityTable.activeTab
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: `【${tabName}】${this.$t('tw_table_noChoose_tips')}`
        })
      }
      const params = {
        formData: this.handleData.formData,
        comment: this.taskForm.comment,
        choseOption: this.taskForm.choseOption,
        handleStatus: this.taskForm.handleStatus,
        taskHandleId: this.taskHandleId
      }
      const { statusCode } = await saveTaskData(this.handleData.id, params)
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
      const requestData = deepClone((this.$refs.entityTable && this.$refs.entityTable.requestData) || [])
      this.handleData.formData =
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
            })
          }
          return item
        }) || []
      // 必填项校验提示
      if (!requiredCheck(this.handleData.formData, this.$refs.entityTable)) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('required_tip')
        })
      }
      // 表格至少勾选一条数据校验
      if (!noChooseCheck(this.handleData.formData, this.$refs.entityTable)) {
        const tabName = this.$refs.entityTable.activeTab
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: `【${tabName}】${this.$t('tw_table_noChoose_tips')}`
        })
      }
      const params = {
        formData: this.handleData.formData,
        comment: this.taskForm.comment,
        choseOption: this.taskForm.choseOption,
        handleStatus: this.taskForm.handleStatus,
        taskHandleId: this.taskHandleId
      }
      const { statusCode } = await commitTaskData(this.handleData.id, params)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.$router.push({
          path: `/taskman/workbench?tabName=hasProcessed&actionName=${this.actionName}&type=${
            this.detail.status === 'InProgress' ? 2 : 3
          }`
        })
      }
    },
    // 获取关注的任务列表
    async geTaskTagList () {
      const { statusCode, data } = await geTaskTagList(this.requestId)
      if (statusCode === 'OK') {
        this.taskTagList = data || []
      }
    },
    // 请求确认提交
    async confirmRequest () {
      const params = {
        id: this.requestId,
        taskId: this.handleData.id,
        markTaskId: this.confirmRequestForm.markTaskId,
        completeStatus: this.confirmRequestForm.completeStatus,
        notes: this.confirmRequestForm.notes
      }
      const { statusCode } = await confirmRequest(params)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.$router.push({ path: `/taskman/workbench?tabName=hasProcessed&actionName=${this.actionName}&type=4` })
      }
    }
  }
}
</script>
<style lang="scss">
.workbench-current-handle {
  margin-top: 10px;
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
      .title {
        span {
          color: #808695;
        }
      }
      .content {
        padding: 20px 0;
        .no-data {
          height: 60px;
          line-height: 60px;
          color: #515a6e;
        }
      }
    }
  }
  .sub-title {
    width: calc(100% - 150px);
    font-size: 14px;
    margin-left: 5px;
    display: flex;
    align-items: center;
    .description {
      margin-left: 10px;
      color: #808695;
      max-width: calc(100% - 300px);
      display: inline-block;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }
  .ivu-tag {
    line-height: 24px !important;
    padding: 0px 5px !important;
  }
}
</style>
