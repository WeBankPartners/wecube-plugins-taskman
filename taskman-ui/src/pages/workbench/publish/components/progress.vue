<!--请求进度-->
<template>
  <div class="workbench-base-progress">
    <!--请求进度-->
    <div class="steps">
      <div class="title">{{ $t('tw_request_progress') }}：</div>
      <Steps :current="0" :style="{ width: progress.requestProgress.length * 170 + 'px' }">
        <Step v-for="(i, index) in progress.requestProgress" :key="index" :content="i.name">
          <template #icon>
            <Icon style="font-weight:bold" size="22" :type="i.icon" :color="i.color" />
            <span v-if="i.node === 'task' && status" @click="handleExpand(i.node)" class="expand-btn">
              {{ taskExpand ? '收起' : '展开' }}
            </span>
            <span v-if="i.node === 'approval' && status" @click="handleExpand(i.node)" class="expand-btn">
              {{ approvalExpand ? '收起' : '展开' }}
            </span>
          </template>
          <div class="role" slot="content">
            <Tooltip :content="i.name">
              <div class="word-eclipse">{{ i.name }}</div>
            </Tooltip>
            <div class="word-eclipse" style="margin-top:-5px;">
              <span v-if="i.role">{{ i.role }} /</span>
              <span>{{ i.handler }}</span>
            </div>
          </div>
        </Step>
      </Steps>
      <div v-if="errorNode" class="error-node">
        <Alert v-if="errorNode === 'autoExit'" show-icon type="error">
          {{ $t('tw_auto_exit_tips') }}
        </Alert>
        <Alert v-else-if="errorNode === 'internallyTerminated'" show-icon type="error">
          {{ $t('tw_terminate_tips') }}
        </Alert>
        <Alert v-else show-icon type="error"> {{ errorNode }}{{ $t('tw_tag_error_tips') }} </Alert>
      </div>
    </div>
    <!--审批进度-->
    <div v-if="approvalExpand" class="steps" style="margin-top:15px;">
      <span class="title">审批进度：</span>
      <Steps :current="0" :style="{ width: progress.approvalProgress.length * 170 + 'px' }">
        <Step v-for="(i, index) in progress.approvalProgress" :key="index" :content="i.name">
          <template #icon>
            <Icon style="font-weight:bold" size="22" :type="i.icon" :color="i.color" />
          </template>
          <div class="role" slot="content">
            <Tooltip :content="i.name">
              <div class="word-eclipse">{{ i.name }}</div>
            </Tooltip>
            <span class="mode">{{ approvalTypeName[i.approveType] || '' }}</span>
            <div v-for="(j, index) in i.taskHandleList" :key="index" class="word-eclipse">
              <span>{{ j.role || '-' }} /</span>
              <span>{{ j.handler || handlerType[j.handlerType] || '-' }}</span>
            </div>
          </div>
        </Step>
      </Steps>
    </div>
    <!--任务进度-->
    <div v-if="taskExpand" class="steps" style="margin-top:15px;">
      <span class="title">任务进度：</span>
      <Steps :current="0" :style="{ width: progress.taskProgress.length * 170 + 'px' }">
        <Step v-for="(i, index) in progress.taskProgress" :key="index" :content="i.name">
          <template #icon>
            <Icon style="font-weight:bold" size="22" :type="i.icon" :color="i.color" />
          </template>
          <div class="role" slot="content">
            <Tooltip :content="i.name">
              <div class="word-eclipse">{{ i.name }}</div>
            </Tooltip>
            <span class="mode">{{ approvalTypeName[i.approveType] || '' }}</span>
            <div v-for="(j, index) in i.taskHandleList" :key="index" class="word-eclipse">
              <template v-if="j.handler === 'autoNode'">
                <span>{{ $t('tw_auto_tag') }}</span>
              </template>
              <template v-else>
                <span>{{ j.role || '-' }} /</span>
                <span>{{ j.handler || handlerType[j.handlerType] || '-' }}</span>
              </template>
            </div>
          </div>
        </Step>
      </Steps>
    </div>
  </div>
</template>

<script>
import { getProgressInfo } from '@/api/server'
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
  props: {
    status: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      progress: {
        approvalProgress: [],
        requestProgress: [],
        taskProgress: []
      },
      approvalExpand: false,
      taskExpand: false,
      errorNode: '',
      approvalTypeName: {
        custom: '单人',
        any: '协同',
        all: '并行',
        admin: '提交人角色管理员',
        auto: '自动通过'
      },
      handlerType: {
        template: '模板指定',
        template_suggest: '模板建议',
        custom: '提交人指定',
        custom_suggest: '提交人建议',
        system: '组内系统分配',
        claim: '组内主动认领'
      }
    }
  },
  watch: {
    status (val) {
      if (val === 'InApproval') {
        this.approvalExpand = true
      } else if (val === 'InProgress') {
        this.taskExpand = true
      }
    }
  },
  methods: {
    // 获取请求进度
    async initData (requestId) {
      const params = {
        params: {
          requestId: requestId
        }
      }
      const { statusCode, data } = await getProgressInfo(params)
      if (statusCode === 'OK') {
        const { approvalProgress, requestProgress, taskProgress } = data
        this.progress.requestProgress = requestProgress || [] // 请求进度
        this.progress.approvalProgress = approvalProgress || [] // 审批进度
        this.progress.taskProgress = taskProgress || [] // 任务进度
        // 请求进度节点处理
        this.progress.requestProgress.forEach(item => {
          item.icon = statusIcon[item.status]
          item.color = statusColor[item.status]
          switch (item.node) {
            case 'submit':
              item.name = this.$t('tw_commit_request') // 提交请求
              break
            case 'check':
              item.name = this.$t('tw_request_pending') // 请求定版
              break
            case 'approval':
              item.name = '审批' // 审批
              item.handler = `${this.progress.approvalProgress.length}个节点`
              break
            case 'task':
              item.name = '任务' // 任务
              // 过滤掉自动节点
              const noAutoTagList = this.progress.taskProgress.filter(i => i.nodeType !== 'auto') || []
              item.handler = `${noAutoTagList.length}个节点`
              break
            case 'confirm':
              item.name = '请求确认' // 请求确认
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
        this.progress.approvalProgress.forEach(item => {
          item.icon = statusIcon[item.status]
          item.color = statusColor[item.status]
          item.name = item.node
        })
        this.progress.taskProgress.forEach(item => {
          item.icon = statusIcon[item.status]
          item.color = statusColor[item.status]
          switch (item.node) {
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
          if (item.nodeType === 'auto') {
            this.errorNode = item.name
          }
        })
      }
    },
    handleExpand (node) {
      if (node === 'approval') {
        this.approvalExpand = !this.approvalExpand
      } else {
        this.taskExpand = !this.taskExpand
      }
    }
  }
}
</script>
<style lang="scss">
.workbench-base-progress {
  .ivu-steps-content {
    padding: 5px !important;
    font-size: 12px;
    color: #3d3c38 !important;
  }
  .ivu-steps-item {
    display: inline-block;
    position: relative;
    vertical-align: top;
    flex: 1;
    overflow: hidden;
    width: 170px;
  }
  .ivu-alert.ivu-alert-with-icon {
    padding: 8px 5px 8px 38px;
  }
  .steps .ivu-steps .ivu-steps-tail > i {
    height: 3px;
    background: #8189a5;
  }
  .steps {
    display: flex;
    align-items: flex-start;
    .title {
      display: inline-block;
      width: 80px;
      font-size: 14px;
      font-weight: 500;
      margin-top: 3px;
    }
    .error-node {
      flex: 1;
      margin: 3px 0 0 -90px;
      max-width: 550px;
    }
    .mode {
      font-size: 12px;
      color: #2d8cf0;
      display: inline-block;
      margin-top: -5px;
    }
    .role {
      display: flex;
      flex-direction: column;
    }
    .word-eclipse {
      max-width: 200px;
      text-overflow: ellipsis;
      overflow: hidden;
      white-space: normal;
    }
    .expand-btn {
      font-size: 14px;
      color: #2d8cf0 !important;
      cursor: pointer;
    }
  }
}
</style>
