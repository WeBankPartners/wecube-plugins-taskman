import dayjs from 'dayjs'
export default {
  data () {
    return {
      baseColumn: {
        id: {
          title: '请求ID',
          width: 140,
          key: 'id',
          render: (h, params) => {
            return (
              <span
                style="cursor:pointer;"
                onClick={() => {
                  this.handleDbClick(params.row)
                }}
              >
                {params.row.id}
              </span>
            )
          }
        },
        name: {
          title: '请求名称',
          sortable: 'custom',
          minWidth: 250,
          key: 'name',
          render: (h, params) => {
            return (
              <span
                style="cursor:pointer;"
                onClick={() => {
                  this.handleDbClick(params.row)
                }}
              >
                {params.row.name}
              </span>
            )
          }
        },
        status: {
          title: '请求状态',
          sortable: 'custom',
          key: 'status',
          minWidth: 130,
          render: (h, params) => {
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
            const item = list.find(i => i.value === params.row.status)
            return (
              item && (
                <Tag color={item.color}>
                  {// 已处理请求定版的草稿添加被退回说明
                    this.tabName === 'hasProcessed' && this.form.type === 1 && params.row.status === 'Draft'
                      ? `${item.label}(被退回)`
                      : item.label}
                </Tag>
              )
            )
          }
        },
        curNode: {
          title: '当前节点',
          minWidth: 120,
          key: 'curNode',
          render: (h, params) => {
            const map = {
              waitCommit: '等待提交',
              sendRequest: '提起请求',
              requestPending: '请求定版',
              requestComplete: '请求完成',
              Completed: '请求完成'
            }
            return <Tag>{map[params.row.curNode] || params.row.curNode}</Tag>
          }
        },
        handler: {
          title: '当前处理人',
          sortable: 'custom',
          minWidth: 140,
          key: 'handler',
          render: (h, params) => {
            return (
              <div style="display:flex;flex-direction:column">
                <span>{params.row.handler}</span>
                <span>{params.row.handleRole}</span>
              </div>
            )
          }
        },
        progress: {
          title: '进展',
          width: 120,
          key: 'progress',
          render: (h, params) => {
            return (
              <Progress percent={params.row.progress}>
                <span>{params.row.progress + '%'}</span>
              </Progress>
            )
          }
        },
        effectiveDays: {
          renderHeader: () => {
            return <span>{this.form.type === 2 ? '任务停留时长' : '请求停留时长'}</span>
          },
          minWidth: 140,
          key: 'effectiveDays',
          render: (h, params) => {
            const diff = params.row.startTime ? dayjs(new Date()).diff(params.row.startTime, 'day') : 0
            const percent = (diff / params.row.effectiveDays) * 100
            const color = percent > 50 ? (percent > 80 ? '#bd3124' : '#ffbf6b') : '#81b337'
            return (
              <Progress stroke-color={color} percent={percent > 100 ? 100 : percent}>
                <span>{`${diff}日/${params.row.effectiveDays}日`}</span>
              </Progress>
            )
          }
        },
        templateName: {
          title: '使用模板',
          sortable: 'custom',
          minWidth: 200,
          key: 'templateName',
          render: (h, params) => {
            return (
              <span>
                {params.row.templateName}
                <Tag>{params.row.version}</Tag>
              </span>
            )
          }
        },
        procDefName: {
          title: '使用编排',
          sortable: 'custom',
          minWidth: 150,
          key: 'procDefName'
        },
        operatorObjType: {
          title: '操作对象类型',
          resizable: true,
          sortable: 'custom',
          minWidth: 150,
          key: 'operatorObjType',
          render: (h, params) => {
            return params.row.operatorObjType && <Tag>{params.row.operatorObjType}</Tag>
          }
        },
        operatorObj: {
          title: '操作对象',
          resizable: true,
          sortable: 'custom',
          minWidth: 150,
          key: 'operatorObj'
        },
        createdBy: {
          title: '创建人',
          sortable: 'custom',
          minWidth: 140,
          key: 'createdBy',
          render: (h, params) => {
            return (
              <div style="display:flex;flex-direction:column">
                <span>{params.row.createdBy}</span>
                <span>{params.row.role}</span>
              </div>
            )
          }
        },
        action: {
          title: this.$t('t_action'),
          key: 'action',
          minWidth: 160,
          fixed: 'right',
          align: 'center',
          render: (h, params) => {
            return (
              <div>
                {['pending', 'hasProcessed', 'submit'].includes(this.tabName) && params.row.status !== 'Draft' && (
                  <Button
                    size="small"
                    onClick={() => {
                      this.hanldeView(params.row)
                    }}
                    style="margin-right:5px;"
                  >
                    查看
                  </Button>
                )}
                {this.username === params.row.handler &&
                  ['Pending', 'InProgress'].includes(params.row.status) &&
                  this.tabName === 'pending' && (
                  <Button
                    type="warning"
                    size="small"
                    onClick={() => {
                      this.handleEdit(params.row)
                    }}
                  >
                      处理
                  </Button>
                )}
                {!params.row.handler &&
                  ['Pending', 'InProgress'].includes(params.row.status) &&
                  this.tabName === 'pending' && (
                  <Button
                    type="info"
                    size="small"
                    onClick={() => {
                      this.handleTransfer(params.row, 'mark')
                    }}
                  >
                      认领
                  </Button>
                )}
                {params.row.handler &&
                  this.username !== params.row.handler &&
                  ['Pending', 'InProgress'].includes(params.row.status) &&
                  this.tabName === 'pending' && (
                  <Button
                    type="success"
                    size="small"
                    onClick={() => {
                      this.handleTransfer(params.row, 'give')
                    }}
                  >
                      转给我
                  </Button>
                )}
                {['Termination', 'Completed', 'Faulted'].includes(params.row.status) && this.tabName === 'submit' && (
                  <Button
                    type="primary"
                    size="small"
                    onClick={() => {
                      this.handleRepub(params.row)
                    }}
                  >
                    重新发起
                  </Button>
                )}
                {params.row.status === 'Pending' && this.tabName === 'submit' && (
                  <Button
                    type="error"
                    size="small"
                    onClick={() => {
                      this.handleRecall(params.row)
                    }}
                  >
                    撤回
                  </Button>
                )}
                {params.row.status === 'Draft' && (
                  <Button
                    type="success"
                    size="small"
                    onClick={() => {
                      this.hanldeLaunch(params.row)
                    }}
                    style="margin-right:5px;"
                  >
                    去发起
                  </Button>
                )}
                {this.tabName === 'draft' && (
                  <Button
                    type="error"
                    size="small"
                    onClick={() => {
                      this.handleDeleteDraft(params.row)
                    }}
                  >
                    删除
                  </Button>
                )}
              </div>
            )
          }
        }
      },
      // pending待处理,hasProcessed已处理,submit我提交的,draft我的暂存,collect收藏
      pendingTaskColumn: [],
      pendingColumn: [],
      hasProcessedTaskColumn: [],
      hasProcessedColumn: [],
      submitAllColumn: [],
      submitColumn: [],
      draftColumn: []
    }
  },
  mounted () {
    // 待处理-任务处理
    this.pendingTaskColumn = [
      this.baseColumn.id,
      {
        title: '任务名称',
        sortable: 'custom',
        minWidth: 120,
        key: 'taskName',
        render: (h, params) => {
          return (
            <span
              style="cursor:pointer;"
              onClick={() => {
                this.handleDbClick(params.row)
              }}
            >
              {params.row.taskName}
            </span>
          )
        }
      },
      this.baseColumn.name,
      this.baseColumn.status,
      this.baseColumn.curNode,
      this.baseColumn.handler,
      this.baseColumn.progress,
      this.baseColumn.effectiveDays,
      {
        title: '任务期望完成时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'taskExpectTime'
      },
      this.baseColumn.templateName,
      this.baseColumn.procDefName,
      this.baseColumn.operatorObjType,
      this.baseColumn.operatorObj,
      this.baseColumn.createdBy,
      {
        title: '任务提交时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'taskCreatedTime'
      },
      this.baseColumn.action
    ]

    // 待处理-请求定版
    this.pendingColumn = [
      this.baseColumn.id,
      this.baseColumn.name,
      this.baseColumn.status,
      this.baseColumn.curNode,
      this.baseColumn.handler,
      this.baseColumn.progress,
      this.baseColumn.effectiveDays,
      {
        title: '请求期望完成时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'expectTime'
      },
      this.baseColumn.templateName,
      this.baseColumn.procDefName,
      this.baseColumn.operatorObjType,
      this.baseColumn.operatorObj,
      this.baseColumn.createdBy,
      {
        title: '任务提交时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'reportTime'
      },
      this.baseColumn.action
    ]

    // 已处理-任务处理
    this.hasProcessedTaskColumn = [
      this.baseColumn.id,
      {
        title: '任务名称',
        sortable: 'custom',
        minWidth: 120,
        key: 'taskName',
        render: (h, params) => {
          return (
            <span
              style="cursor:pointer;"
              onClick={() => {
                this.handleDbClick(params.row)
              }}
            >
              {params.row.taskName}
            </span>
          )
        }
      },
      this.baseColumn.name,
      this.baseColumn.status,
      this.baseColumn.curNode,
      this.baseColumn.handler,
      this.baseColumn.progress,
      this.baseColumn.effectiveDays,
      {
        title: '任务期望完成时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'taskExpectTime'
      },
      this.baseColumn.templateName,
      this.baseColumn.procDefName,
      this.baseColumn.operatorObjType,
      this.baseColumn.operatorObj,
      this.baseColumn.createdBy,
      {
        title: '任务提交时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'taskCreatedTime'
      },
      {
        title: '处理时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'taskApprovalTime'
      },
      this.baseColumn.action
    ]

    // 已处理-请求定版
    this.hasProcessedColumn = [
      this.baseColumn.id,
      this.baseColumn.name,
      this.baseColumn.status,
      this.baseColumn.curNode,
      this.baseColumn.handler,
      this.baseColumn.progress,
      this.baseColumn.effectiveDays,
      {
        title: '请求期望完成时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'expectTime'
      },
      {
        title: '退回原因',
        sortable: 'custom',
        minWidth: 150,
        key: 'rollbackDesc'
      },
      this.baseColumn.templateName,
      this.baseColumn.procDefName,
      this.baseColumn.operatorObjType,
      this.baseColumn.operatorObj,
      this.baseColumn.createdBy,
      {
        title: '任务提交时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'reportTime'
      },
      {
        title: '处理时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'approvalTime'
      },
      this.baseColumn.action
    ]

    // 我提交的-所有、被退回
    this.submitAllColumn = [
      this.baseColumn.id,
      this.baseColumn.name,
      this.baseColumn.status,
      this.baseColumn.curNode,
      this.baseColumn.handler,
      this.baseColumn.progress,
      this.baseColumn.effectiveDays,
      {
        title: '期望完成时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'expectTime'
      },
      {
        title: '退回原因',
        sortable: 'custom',
        minWidth: 150,
        key: 'rollbackDesc'
      },
      this.baseColumn.templateName,
      this.baseColumn.procDefName,
      this.baseColumn.operatorObjType,
      this.baseColumn.operatorObj,
      this.baseColumn.createdBy,
      {
        title: '请求提交时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'reportTime'
      },
      this.baseColumn.action
    ]

    // 我提交的-其他、我撤回的
    this.submitColumn = [
      this.baseColumn.id,
      this.baseColumn.name,
      this.baseColumn.status,
      this.baseColumn.curNode,
      this.baseColumn.handler,
      this.baseColumn.progress,
      this.baseColumn.effectiveDays,
      {
        title: '期望完成时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'expectTime'
      },
      this.baseColumn.templateName,
      this.baseColumn.procDefName,
      this.baseColumn.operatorObjType,
      this.baseColumn.operatorObj,
      this.baseColumn.createdBy,
      {
        title: '请求提交时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'reportTime'
      },
      this.baseColumn.action
    ]

    // 我暂存的
    this.draftColumn = [
      this.baseColumn.id,
      this.baseColumn.name,
      this.baseColumn.status,
      this.baseColumn.curNode,
      this.baseColumn.templateName,
      this.baseColumn.procDefName,
      this.baseColumn.operatorObjType,
      this.baseColumn.operatorObj,
      {
        title: '定版处理人',
        sortable: 'custom',
        minWidth: 140,
        key: 'handler',
        render: (h, params) => {
          return (
            <div style="display:flex;flex-direction:column">
              <span>{params.row.handler}</span>
              <span>{params.row.handleRole}</span>
            </div>
          )
        }
      },
      this.baseColumn.createdBy,
      {
        title: '创建时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'createdTime'
      },
      {
        title: '更新时间',
        sortable: 'custom',
        minWidth: 150,
        key: 'updatedTime'
      },
      this.baseColumn.action
    ]
  }
}
