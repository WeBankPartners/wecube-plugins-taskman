export default {
  data () {
    return {
      baseSearch: {
        id: {
          key: 'id',
          placeholder: 'id',
          component: 'input'
        },
        name: {
          key: 'name',
          placeholder: '名称',
          component: 'input'
        },
        status: {
          key: 'status',
          placeholder: '状态',
          component: 'select',
          multiple: true,
          list: [
            { label: this.$t('status_pending'), value: 'Pending' },
            { label: this.$t('status_inProgress'), value: 'InProgress' },
            { label: this.$t('status_inProgress_faulted'), value: 'InProgress(Faulted)' },
            { label: this.$t('status_termination'), value: 'Termination' },
            { label: this.$t('status_complete'), value: 'Completed' },
            { label: this.$t('status_inProgress_timeouted'), value: 'InProgress(Timeouted)' },
            { label: this.$t('status_faulted'), value: 'Faulted' },
            { label: this.$t('status_draft'), value: 'Draft' }
          ]
        },
        handler: {
          key: 'handler',
          placeholder: '处理人',
          component: 'select',
          multiple: true,
          list: []
        },
        createdBy: {
          key: 'createdBy',
          placeholder: '创建人',
          component: 'select',
          multiple: true,
          list: []
        },
        templateId: {
          key: 'templateId',
          placeholder: '模板',
          multiple: true,
          component: 'select',
          list: []
        },
        procDefName: {
          key: 'procDefName',
          placeholder: '使用编排',
          multiple: true,
          component: 'select',
          list: []
        },
        operatorObjType: {
          key: 'operatorObjType',
          placeholder: '操作对象类型',
          multiple: true,
          component: 'select',
          list: []
        }
      },
      // pending待处理,hasProcessed已处理,submit我提交的,draft我的暂存,collect收藏
      pendingTaskSearch: [],
      pendingSearch: [],
      hasProcessedTaskSearch: [],
      hasProcessedSearch: [],
      submitSearch: [],
      draftSearch: []
    }
  },
  mounted () {
    // 待处理-任务处理
    this.pendingTaskSearch = [
      this.baseSearch.id,
      {
        key: 'taskName',
        placeholder: '任务名称',
        component: 'input'
      },
      this.baseSearch.name,
      this.baseSearch.status,
      this.baseSearch.handler,
      this.baseSearch.createdBy,
      this.baseSearch.templateId,
      this.baseSearch.procDefName,
      this.baseSearch.operatorObjType,
      {
        key: 'taskExpectTime',
        label: '任务期望完成时间',
        dateType: 4,
        labelWidth: 140,
        component: 'custom-time'
      },
      {
        key: 'taskReportTime',
        label: '任务提交时间',
        dateType: 4,
        labelWidth: 110,
        component: 'custom-time'
      }
    ]

    // 待处理-请求定版
    this.pendingSearch = [
      this.baseSearch.id,
      this.baseSearch.name,
      this.baseSearch.status,
      this.baseSearch.handler,
      this.baseSearch.createdBy,
      this.baseSearch.templateId,
      this.baseSearch.procDefName,
      this.baseSearch.operatorObjType,
      {
        key: 'expectTime',
        label: '请求期望完成时间',
        dateType: 4,
        labelWidth: 140,
        component: 'custom-time'
      },
      {
        key: 'reportTime',
        label: '任务提交时间',
        dateType: 4,
        labelWidth: 110,
        component: 'custom-time'
      }
    ]

    // 已处理-任务处理
    this.hasProcessedTaskSearch = [
      this.baseSearch.id,
      {
        key: 'taskName',
        placeholder: '任务名称',
        component: 'input'
      },
      this.baseSearch.name,
      this.baseSearch.status,
      this.baseSearch.handler,
      this.baseSearch.createdBy,
      this.baseSearch.templateId,
      this.baseSearch.procDefName,
      this.baseSearch.operatorObjType,
      {
        key: 'taskExpectTime',
        label: '任务期望完成时间',
        dateType: 4,
        labelWidth: 140,
        component: 'custom-time'
      },
      {
        key: 'taskReportTime',
        label: '任务提交时间',
        dateType: 4,
        labelWidth: 110,
        component: 'custom-time'
      },
      {
        key: 'taskApprovalTime',
        label: '处理时间',
        dateType: 4,
        labelWidth: 85,
        component: 'custom-time'
      }
    ]

    // 已处理-请求定版
    this.hasProcessedSearch = [
      this.baseSearch.id,
      this.baseSearch.name,
      this.baseSearch.status,
      this.baseSearch.handler,
      this.baseSearch.createdBy,
      this.baseSearch.templateId,
      this.baseSearch.procDefName,
      this.baseSearch.operatorObjType,
      {
        key: 'expectTime',
        label: '请求期望完成时间',
        dateType: 4,
        labelWidth: 140,
        component: 'custom-time'
      },
      {
        key: 'reportTime',
        label: '任务提交时间',
        dateType: 4,
        labelWidth: 110,
        component: 'custom-time'
      },
      {
        key: 'approvalTime',
        label: '处理时间',
        dateType: 4,
        labelWidth: 85,
        component: 'custom-time'
      }
    ]

    // 我提交的
    this.submitSearch = [
      this.baseSearch.id,
      this.baseSearch.name,
      this.baseSearch.status,
      this.baseSearch.handler,
      this.baseSearch.createdBy,
      this.baseSearch.templateId,
      this.baseSearch.procDefName,
      this.baseSearch.operatorObjType,
      {
        key: 'expectTime',
        label: '期望完成时间',
        dateType: 4,
        labelWidth: 110,
        component: 'custom-time'
      },
      {
        key: 'reportTime',
        label: '请求提交时间',
        dateType: 4,
        labelWidth: 110,
        component: 'custom-time'
      }
    ]

    // 我暂存的
    this.draftSearch = [
      this.baseSearch.id,
      this.baseSearch.name,
      this.baseSearch.status,
      this.baseSearch.handler,
      this.baseSearch.createdBy,
      this.baseSearch.templateId,
      this.baseSearch.procDefName,
      this.baseSearch.operatorObjType,
      {
        key: 'createdTime',
        label: '创建时间',
        dateType: 4,
        labelWidth: 85,
        component: 'custom-time'
      },
      {
        key: 'updatedTime',
        label: '更新时间',
        dateType: 4,
        labelWidth: 85,
        component: 'custom-time'
      },
      {
        key: 'expectTime',
        label: '期望完成时间',
        dateType: 4,
        labelWidth: 110,
        component: 'custom-time'
      }
    ]
  }
}
