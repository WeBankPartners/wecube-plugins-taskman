export default {
  data () {
    return {
      baseColumn: {
        id: {
          title: this.$t('request_id'),
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
          title: this.$t('request_name'),
          sortable: 'custom',
          minWidth: 250,
          key: 'name',
          render: (h, params) => {
            return (
              <div style="display:flex;flex-direction:column;align-items:flex-start;">
                <span
                  style="cursor:pointer;"
                  onClick={() => {
                    this.handleDbClick(params.row)
                  }}
                >
                  {params.row.name}
                </span>
                {/* {this.username === params.row.handler &&
                  ['Pending', 'InProgress'].includes(params.row.status) &&
                  this.tabName === 'pending' && <Tag color="#ed4014">{this.$t('tw_only_me')}</Tag>} */}
              </div>
            )
          }
        },
        status: {
          title: this.$t('tw_request_status'),
          sortable: 'custom',
          key: 'status',
          minWidth: 160,
          render: (h, params) => {
            const list = [
              { label: this.$t('status_pending'), value: 'Pending', color: '#b886f8' },
              { label: this.$t('status_inProgress'), value: 'InProgress', color: '#1990ff' },
              { label: this.$t('status_inProgress_faulted'), value: 'InProgress(Faulted)', color: '#ed4014' },
              { label: this.$t('status_termination'), value: 'Termination', color: '#e29836' },
              { label: this.$t('status_complete'), value: 'Completed', color: '#7ac756' },
              { label: this.$t('status_inProgress_timeouted'), value: 'InProgress(Timeouted)', color: '#ed4014' },
              { label: this.$t('status_faulted'), value: 'Faulted', color: '#e29836' },
              { label: this.$t('status_draft'), value: 'Draft', color: '#808695' }
            ]
            const item = list.find(i => i.value === params.row.status)
            // 被退回的草稿添加标签
            const tagName =
              params.row.rollbackDesc && params.row.status === 'Draft'
                ? `${item.label}(${this.$t('tw_returned_tips')})`
                : item.label
            return (
              item && (
                <Tooltip content={tagName} placement="top">
                  <Tag color={item.color}>{tagName}</Tag>
                </Tooltip>
              )
            )
          }
        },
        curNode: {
          title: this.$t('tw_cur_tag'),
          minWidth: 165,
          key: 'curNode',
          render: (h, params) => {
            const map = {
              waitCommit: this.$t('tw_wait_commit'),
              sendRequest: this.$t('tw_commit_request'),
              requestPending: this.$t('tw_request_pending'),
              requestComplete: this.$t('tw_request_complete'),
              Completed: this.$t('tw_request_complete')
            }
            return (
              <Tooltip content={map[params.row.curNode] || params.row.curNode} placement="top">
                <Tag>{map[params.row.curNode] || params.row.curNode}</Tag>
              </Tooltip>
            )
          }
        },
        handler: {
          renderHeader: () => {
            return (
              // 我暂存的，我提交的(被退回、本人撤回)显示为定版处理人，其余为当前处理人
              <span>
                {this.tabName === 'draft' || (this.tabName === 'submit' && ['1', '3'].includes(this.rollback))
                  ? this.$t('tw_pending_handler')
                  : this.$t('tw_cur_handler')}
              </span>
            )
          },
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
          title: this.$t('tw_progress'),
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
        // 任务停留时长
        effectiveDays: {
          renderHeader: () => {
            return <span>{this.type === '2' ? this.$t('tw_task_stay_time') : this.$t('tw_request_stay_time')}</span>
          },
          minWidth: 140,
          key: 'effectiveDays',
          render: (h, params) => {
            let stayTime = '' // 停留时长
            let totalTime = '' // 预期停留时长
            if (this.type === '2') {
              stayTime = params.row.taskStayTime
              totalTime = params.row.taskStayTimeTotal
            } else {
              stayTime = params.row.requestStayTime
              totalTime = params.row.requestStayTimeTotal
            }
            const percent = (stayTime / totalTime) * 100
            const color = percent > 50 ? (percent > 80 ? '#ed4014' : '#ffbf6b') : '#19be6b'
            return (
              <div>
                <Progress stroke-color={color} percent={percent > 100 ? 100 : percent}>
                  <span>{`${stayTime}${this.$t('tw_days')}/${totalTime}${this.$t('tw_days')}`}</span>
                </Progress>
                {percent > 100 && (
                  <span style="color:#ed4014;display:flex;align-items:center;">
                    <Icon type="md-warning" color="#ed4014" />
                    {`${this.$t('tw_exceed')}${stayTime - totalTime}${this.$t('tw_days')}`}
                  </span>
                )}
              </div>
            )
          }
        },
        templateName: {
          title: this.$t('tw_use_template'),
          sortable: 'custom',
          minWidth: 200,
          key: 'templateName',
          render: (h, params) => {
            return (
              <span>
                {`${params.row.templateName}【${params.row.version}】`}
                {/* <Tag>{params.row.version}</Tag> */}
              </span>
            )
          }
        },
        procDefName: {
          title: this.$t('tw_template_flow'),
          sortable: 'custom',
          minWidth: 150,
          key: 'procDefName'
        },
        operatorObjType: {
          title: this.$t('tw_operator_type'),
          resizable: true,
          sortable: 'custom',
          minWidth: 150,
          key: 'operatorObjType',
          render: (h, params) => {
            return (
              params.row.operatorObjType && (
                <Tooltip content={params.row.operatorObjType} placement="top">
                  <Tag>{params.row.operatorObjType}</Tag>
                </Tooltip>
              )
            )
          }
        },
        operatorObj: {
          title: this.$t('tw_operator'),
          resizable: true,
          sortable: 'custom',
          minWidth: 150,
          key: 'operatorObj'
        },
        createdBy: {
          title: this.$t('createdBy'),
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
          minWidth: 120,
          fixed: 'right',
          align: 'center',
          render: (h, params) => {
            return (
              <div style="display:flex;align-items:center;justify-content:center;">
                {// 查看
                  ['pending', 'hasProcessed', 'submit'].includes(this.tabName) && params.row.status !== 'Draft' && (
                    <Tooltip content={this.$t('tw_action_view')} placement="top">
                      <Button
                        size="small"
                        type="info"
                        onClick={() => {
                          this.hanldeView(params.row)
                        }}
                        style="margin-right:5px;"
                      >
                        <Icon type="md-eye" size="16"></Icon>
                      </Button>
                    </Tooltip>
                  )}
                {// 处理
                  this.username === params.row.handler &&
                  ['Pending', 'InProgress'].includes(params.row.status) &&
                  this.tabName === 'pending' && (
                    <Tooltip content={this.$t('tw_action_handle')} placement="top">
                      <Button
                        type="warning"
                        size="small"
                        onClick={() => {
                          this.handleEdit(params.row)
                        }}
                      >
                        <Icon type="ios-hammer" size="16"></Icon>
                      </Button>
                    </Tooltip>
                  )}
                {// 认领
                  !params.row.handler &&
                  ['Pending', 'InProgress'].includes(params.row.status) &&
                  this.tabName === 'pending' && (
                    <Tooltip content={this.$t('tw_action_claim')} placement="top">
                      <Button
                        type="info"
                        size="small"
                        onClick={() => {
                          this.handleTransfer(params.row, 'mark')
                        }}
                      >
                        <Icon type="ios-hand" size="16"></Icon>
                      </Button>
                    </Tooltip>
                  )}
                {// 转给我
                  params.row.handler &&
                  this.username !== params.row.handler &&
                  ['Pending', 'InProgress'].includes(params.row.status) &&
                  this.tabName === 'pending' && (
                    <Tooltip content={this.$t('tw_action_give')} placement="top">
                      <Button
                        type="success"
                        size="small"
                        onClick={() => {
                          this.handleTransfer(params.row, 'give')
                        }}
                      >
                        <Icon type="ios-hand" size="16"></Icon>
                      </Button>
                    </Tooltip>
                  )}
                {// 重新发起
                  ['Termination', 'Completed', 'Faulted'].includes(params.row.status) && this.tabName === 'submit' && (
                    <Tooltip content={this.$t('tw_action_relaunch')} placement="top">
                      <Button
                        type="success"
                        size="small"
                        onClick={() => {
                          this.handleRepub(params.row)
                        }}
                      >
                        <Icon type="ios-refresh" size="16"></Icon>
                      </Button>
                    </Tooltip>
                  )}
                {// 撤回
                // 我提交的定版状态可退回
                  params.row.status === 'Pending' && this.tabName === 'submit' && (
                    <Tooltip content={this.$t('tw_recall')} placement="top">
                      <Button
                        type="error"
                        size="small"
                        // 非本人创建禁用
                        disabled={params.row.createdBy !== this.username}
                        onClick={() => {
                          this.handleRecall(params.row)
                        }}
                      >
                        <Icon type="ios-redo" size="16"></Icon>
                      </Button>
                    </Tooltip>
                  )}
                {// 去发起
                // 草稿类（不包括已处理的）
                  params.row.status === 'Draft' && this.tabName !== 'hasProcessed' && (
                    <Tooltip
                      content={this.tabName === 'submit' ? this.$t('tw_action_relaunch') : this.$t('tw_action_launch')}
                      placement="top"
                    >
                      <Button
                        type="success"
                        size="small"
                        // 非本人创建禁用
                        disabled={params.row.createdBy !== this.username}
                        onClick={() => {
                          this.hanldeLaunch(params.row)
                        }}
                        style="margin-right:5px;"
                      >
                        <Icon type="ios-send" size="16"></Icon>
                      </Button>
                    </Tooltip>
                  )}
                {// 删除
                  this.tabName === 'draft' && (
                    <Tooltip content={this.$t('delete')} placement="top">
                      <Button
                        type="error"
                        size="small"
                        // 非本人创建禁用
                        disabled={params.row.createdBy !== this.username}
                        onClick={() => {
                          this.handleDeleteDraft(params.row)
                        }}
                      >
                        <Icon type="md-trash" size="16"></Icon>
                      </Button>
                    </Tooltip>
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
      draftColumn: [],
      username: window.localStorage.getItem('username')
    }
  },
  mounted () {
    // 待处理-任务处理
    this.pendingTaskColumn = [
      this.baseColumn.id,
      {
        title: this.$t('task_name'),
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
        title: this.$t('tw_task_commit_time'),
        sortable: 'custom',
        minWidth: 150,
        key: 'taskCreatedTime'
      },
      {
        title: this.$t('tw_task_expect_time'),
        sortable: 'custom',
        minWidth: 150,
        key: 'taskExpectTime'
      },
      this.baseColumn.createdBy,
      this.baseColumn.templateName,
      this.baseColumn.procDefName,
      this.baseColumn.operatorObjType,
      this.baseColumn.operatorObj,
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
        title: this.$t('tw_task_commit_time'),
        sortable: 'custom',
        minWidth: 150,
        key: 'reportTime'
      },
      {
        title: this.$t('tw_request_expect_time'),
        sortable: 'custom',
        minWidth: 150,
        key: 'expectTime'
      },
      this.baseColumn.createdBy,
      this.baseColumn.templateName,
      this.baseColumn.procDefName,
      this.baseColumn.operatorObjType,
      this.baseColumn.operatorObj,
      this.baseColumn.action
    ]

    // 已处理-任务处理
    this.hasProcessedTaskColumn = [
      this.baseColumn.id,
      {
        title: this.$t('task_name'),
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
        title: this.$t('tw_task_commit_time'),
        sortable: 'custom',
        minWidth: 150,
        key: 'taskCreatedTime'
      },
      {
        title: this.$t('tw_task_expect_time'),
        sortable: 'custom',
        minWidth: 150,
        key: 'taskExpectTime'
      },
      this.baseColumn.createdBy,
      this.baseColumn.templateName,
      this.baseColumn.procDefName,
      this.baseColumn.operatorObjType,
      this.baseColumn.operatorObj,
      {
        title: this.$t('handle_time'),
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
        title: this.$t('tw_task_commit_time'),
        sortable: 'custom',
        minWidth: 150,
        key: 'reportTime'
      },
      {
        title: this.$t('tw_request_expect_time'),
        sortable: 'custom',
        minWidth: 150,
        key: 'expectTime'
      },
      this.baseColumn.createdBy,
      {
        title: this.$t('tw_rollback_reason'),
        sortable: 'custom',
        minWidth: 150,
        key: 'rollbackDesc',
        render: (h, params) => {
          return (
            <Tooltip max-width="300" content={params.row.rollbackDesc}>
              <span style="overflow:hidden;text-overflow:ellipsis;display:-webkit-box;-webkit-line-clamp:3;-webkit-box-orient:vertical;">
                {params.row.rollbackDesc}
              </span>
            </Tooltip>
          )
        }
      },
      this.baseColumn.templateName,
      this.baseColumn.procDefName,
      this.baseColumn.operatorObjType,
      this.baseColumn.operatorObj,
      {
        title: this.$t('handle_time'),
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
        title: this.$t('tw_expect_time'),
        sortable: 'custom',
        minWidth: 150,
        key: 'expectTime'
      },
      {
        title: this.$t('tw_rollback_reason'),
        sortable: 'custom',
        minWidth: 150,
        key: 'rollbackDesc',
        render: (h, params) => {
          return (
            <Tooltip max-width="300" content={params.row.rollbackDesc}>
              <span style="overflow:hidden;text-overflow:ellipsis;display:-webkit-box;-webkit-line-clamp:3;-webkit-box-orient:vertical;">
                {params.row.rollbackDesc}
              </span>
            </Tooltip>
          )
        }
      },
      this.baseColumn.templateName,
      this.baseColumn.procDefName,
      this.baseColumn.operatorObjType,
      this.baseColumn.operatorObj,
      this.baseColumn.createdBy,
      {
        title: this.$t('tw_request_commit_time'),
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
        title: this.$t('tw_expect_time'),
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
        title: this.$t('tw_request_commit_time'),
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
      this.baseColumn.handler,
      this.baseColumn.createdBy,
      {
        title: this.$t('tw_created_time'),
        sortable: 'custom',
        minWidth: 150,
        key: 'createdTime'
      },
      {
        title: this.$t('tw_update_time'),
        sortable: 'custom',
        minWidth: 150,
        key: 'updatedTime'
      },
      {
        title: this.$t('tw_expect_time'),
        sortable: 'custom',
        minWidth: 150,
        key: 'expectTime'
      },
      this.baseColumn.action
    ]
  }
}
