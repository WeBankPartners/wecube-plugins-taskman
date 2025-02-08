import dayjs from 'dayjs'
export default {
  data () {
    return {
      baseSearch: {
        id: {
          key: 'id',
          placeholder: this.$t('tw_request_id'),
          component: 'input'
        },
        name: {
          key: 'name',
          placeholder: this.$t('request_name'),
          component: 'input'
        },
        status: {
          key: 'status',
          placeholder: this.$t('tw_request_status'),
          component: 'tag-select',
          multiple: true,
          list: [
            { label: this.$t('status_draft'), value: 'Draft', color: '#808695' },
            { label: this.$t('status_pending'), value: 'Pending', color: '#b886f8' },
            { label: this.$t('tw_request_confirm'), value: 'Confirm', color: '#b886f8' },
            { label: this.$t('tw_inApproval'), value: 'InApproval', color: '#1990ff' },
            { label: this.$t('status_inProgress'), value: 'InProgress', color: '#1990ff' },
            { label: this.$t('status_complete'), value: 'Completed', color: '#7ac756' },
            { label: this.$t('status_inProgress_faulted'), value: 'InProgress(Faulted)', color: '#FF4D4F' },
            { label: this.$t('status_inProgress_timeouted'), value: 'InProgress(Timeouted)', color: '#FF4D4F' },
            { label: this.$t('tw_stop'), value: 'Stop', color: '#FF4D4F' },
            { label: this.$t('status_faulted'), value: 'Faulted', color: '#e29836' },
            { label: this.$t('status_termination'), value: 'Termination', color: '#e29836' }
          ]
        },
        createdBy: {
          key: 'createdBy',
          placeholder: this.$t('tw_reporter'),
          component: 'select',
          multiple: true,
          list: []
        },
        templateId: {
          key: 'templateId',
          placeholder: this.$t('tw_use_template'),
          multiple: true,
          component: 'select',
          list: []
        },
        procDefName: {
          key: 'procDefName',
          placeholder: this.$t('tw_template_flow'),
          multiple: true,
          component: 'select',
          list: []
        },
        operatorObjType: {
          key: 'operatorObjType',
          placeholder: this.$t('tw_operator_type'),
          multiple: true,
          component: 'select',
          list: []
        },
        requestRefId: {
          key: 'requestRefId',
          placeholder: this.$t('tw_ref_id'),
          nullType: 'no',
          component: 'null-input' // 支持空值搜索
        }
      },
      // 任务工作台
      pendingTaskSearch: [],
      pendingSearch: [],
      hasProcessedTaskSearch: [],
      hasProcessedSearch: [],
      submitSearch: [],
      draftSearch: [],
      initDate: []
    }
  },
  mounted () {
    const cur = dayjs().format('YYYY-MM-DD')
    const pre = dayjs()
      .subtract(3, 'month')
      .format('YYYY-MM-DD')
    this.initDate = [pre, cur]
    this.getFilterOptions()
    // 待处理-任务和审批
    this.pendingTaskSearch = [
      this.baseSearch.id,
      this.baseSearch.name,
      {
        key: 'taskName',
        placeholder: this.$t('task_name'),
        component: 'input'
      },
      {
        key: 'taskHandleUpdatedTime',
        label: this.$t('tw_taskUpdated'),
        dateType: 1,
        initValue: this.initDate,
        labelWidth: 110,
        component: 'custom-time'
      },
      this.baseSearch.status,
      {
        key: 'taskExpectTime',
        label: this.$t('tw_taskEnd'),
        dateType: 4,
        labelWidth: 140,
        component: 'custom-time'
      },
      this.baseSearch.createdBy,
      this.baseSearch.templateId,
      this.baseSearch.procDefName,
      this.baseSearch.operatorObjType,
      this.baseSearch.requestRefId
    ]

    // 待处理-定版和请求确认
    this.pendingSearch = [
      this.baseSearch.id,
      this.baseSearch.name,
      this.baseSearch.status,
      {
        key: 'taskHandleUpdatedTime',
        label: this.$t('tw_taskUpdated'), // 任务更新
        dateType: 1,
        initValue: this.initDate,
        labelWidth: 110,
        component: 'custom-time'
      },
      {
        key: 'taskExpectTime',
        label: this.$t('tw_taskEnd'), // 任务截止
        dateType: 4,
        labelWidth: 140,
        component: 'custom-time'
      },
      this.baseSearch.createdBy,
      this.baseSearch.templateId,
      this.baseSearch.procDefName,
      this.baseSearch.operatorObjType,
      this.baseSearch.requestRefId
    ]

    // 已处理-任务和审批
    this.hasProcessedTaskSearch = [
      this.baseSearch.id,
      this.baseSearch.name,
      {
        key: 'taskName',
        placeholder: this.$t('task_name'),
        component: 'input'
      },
      {
        key: 'taskApprovalTime',
        label: this.$t('handle_time'),
        dateType: 1,
        initValue: this.initDate,
        labelWidth: 85,
        component: 'custom-time'
      },
      this.baseSearch.status,
      this.baseSearch.createdBy,
      this.baseSearch.templateId,
      this.baseSearch.procDefName,
      this.baseSearch.operatorObjType,
      this.baseSearch.requestRefId,
      {
        key: 'taskCreatedTime',
        label: this.$t('tw_taskCreated'), // 任务创建
        dateType: 4,
        labelWidth: 110,
        component: 'custom-time'
      },
      {
        key: 'taskExpectTime',
        label: this.$t('tw_taskEnd'), // 任务截止
        dateType: 4,
        labelWidth: 140,
        component: 'custom-time'
      }
    ]

    // 已处理-定版和请求确认
    this.hasProcessedSearch = [
      this.baseSearch.id,
      this.baseSearch.name,
      this.baseSearch.status,
      {
        key: 'taskApprovalTime',
        label: this.$t('handle_time'),
        dateType: 1,
        initValue: this.initDate,
        labelWidth: 85,
        component: 'custom-time'
      },
      this.baseSearch.createdBy,
      this.baseSearch.templateId,
      this.baseSearch.procDefName,
      this.baseSearch.operatorObjType,
      this.baseSearch.requestRefId,
      {
        key: 'taskCreatedTime',
        label: this.$t('tw_taskCreated'), // 任务创建
        dateType: 4,
        labelWidth: 110,
        component: 'custom-time'
      },
      {
        key: 'taskExpectTime',
        label: this.$t('tw_taskEnd'), // 任务截止
        dateType: 4,
        labelWidth: 140,
        component: 'custom-time'
      }
    ]

    // 我提交的
    this.submitSearch = [
      this.baseSearch.id,
      this.baseSearch.name,
      this.baseSearch.status,
      {
        key: 'reportTime',
        label: this.$t('tw_request_commit_time'),
        dateType: 1,
        initValue: this.initDate,
        labelWidth: 110,
        component: 'custom-time'
      },
      this.baseSearch.createdBy,
      this.baseSearch.templateId,
      this.baseSearch.procDefName,
      this.baseSearch.operatorObjType,
      this.baseSearch.requestRefId,
      {
        key: 'expectTime',
        label: this.$t('tw_expect_time'),
        dateType: 4,
        labelWidth: 110,
        component: 'custom-time'
      }
    ]

    // 我暂存的
    this.draftSearch = [
      this.baseSearch.id,
      this.baseSearch.name,
      // this.baseSearch.status,
      {
        key: 'updatedTime',
        label: this.$t('tw_update_time'),
        dateType: 1,
        initValue: this.initDate,
        labelWidth: 85,
        component: 'custom-time'
      },
      this.baseSearch.createdBy,
      this.baseSearch.templateId,
      this.baseSearch.procDefName,
      this.baseSearch.operatorObjType,
      this.baseSearch.requestRefId,
      {
        key: 'createdTime',
        label: this.$t('tw_created_time'),
        dateType: 4,
        labelWidth: 85,
        component: 'custom-time'
      },
      {
        key: 'expectTime',
        label: this.$t('tw_expect_time'),
        dateType: 4,
        labelWidth: 110,
        component: 'custom-time'
      }
    ]
  },
  methods: {
    // 获取搜索条件的下拉值
    async getFilterOptions () {
      const pre = dayjs()
        .subtract(12, 'month')
        .format('YYYY-MM-DD')
      import('@/api/server').then(async ({ getPlatformFilter }) => {
        const { statusCode, data } = await getPlatformFilter({ startTime: pre })
        if (statusCode === 'OK') {
          const keys = Object.keys(this.baseSearch)
          for (let key of keys) {
            if (key === 'operatorObjType') {
              this.baseSearch[key].list =
                data.operatorObjTypeList &&
                data.operatorObjTypeList.map(item => {
                  return {
                    label: item,
                    value: item
                  }
                })
            } else if (key === 'templateId') {
              // 获取发布模板
              if (this.actionName === '1') {
                this.baseSearch[key].list =
                  data.releaseTemplateList &&
                  data.releaseTemplateList.map(item => {
                    return {
                      label: `${item.templateName}【${item.version}】`,
                      value: item.templateId
                    }
                  })
              // 获取请求模板
              } else if (this.actionName === '2') {
                this.baseSearch[key].list =
                  data.requestTemplateList &&
                  data.requestTemplateList.map(item => {
                    return {
                      label: `${item.templateName}【${item.version}】`,
                      value: item.templateId
                    }
                  })
              // 获取问题模板
              } else if (this.actionName === '3') {
                this.baseSearch[key].list =
                  data.problemTemplateList &&
                  data.problemTemplateList.map(item => {
                    return {
                      label: `${item.templateName}【${item.version}】`,
                      value: item.templateId
                    }
                  })
              // 获取事件模板
              } else if (this.actionName === '4') {
                this.baseSearch[key].list =
                  data.eventTemplateList &&
                  data.eventTemplateList.map(item => {
                    return {
                      label: `${item.templateName}【${item.version}】`,
                      value: item.templateId
                    }
                  })
              // 获取变更模板
              } else if (this.actionName === '5') {
                this.baseSearch[key].list =
                  data.changeTemplateList &&
                  data.changeTemplateList.map(item => {
                    return {
                      label: `${item.templateName}【${item.version}】`,
                      value: item.templateId
                    }
                  })
              // 获取全部模板
              } else {
                this.baseSearch[key].list =
                  data.templateList &&
                  data.templateList.map(item => {
                    return {
                      label: `${item.templateName}【${item.version}】`,
                      value: item.templateId
                    }
                  })
              }
            } else if (key === 'procDefName') {
              this.baseSearch[key].list =
                data.procDefNameList &&
                data.procDefNameList.map(item => {
                  return {
                    label: item,
                    value: item
                  }
                })
            } else if (key === 'createdBy') {
              this.baseSearch[key].list =
                data.createdByList &&
                data.createdByList.map(item => {
                  return {
                    label: item,
                    value: item
                  }
                })
            }
          }
        }
      })
    }
  }
}
