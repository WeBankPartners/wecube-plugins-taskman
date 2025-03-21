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
          minWidth: 200,
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
              { label: this.$t('tw_inApproval'), value: 'InApproval', color: '#5384FF' },
              { label: this.$t('status_inProgress'), value: 'InProgress', color: '#5384FF' },
              { label: this.$t('tw_request_confirm'), value: 'Confirm', color: '#b886f8' },
              { label: this.$t('status_inProgress_faulted'), value: 'InProgress(Faulted)', color: '#FF4D4F' },
              { label: this.$t('status_termination'), value: 'Termination', color: '#F29360' },
              { label: this.$t('status_complete'), value: 'Completed', color: '#00CB91' },
              { label: this.$t('status_inProgress_timeouted'), value: 'InProgress(Timeouted)', color: '#FF4D4F' },
              { label: this.$t('status_faulted'), value: 'Faulted', color: '#F29360' },
              { label: this.$t('status_draft'), value: 'Draft', color: '#808695' },
              { label: this.$t('tw_stop'), value: 'Stop', color: '#FF4D4F' }
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
              confirm: this.$t('tw_request_confirm'),
              requestPending: this.$t('tw_request_pending'),
              requestComplete: this.$t('tw_request_complete'),
              Completed: this.$t('tw_request_complete')
            }
            if (map[params.row.curNode] || params.row.curNode) {
              return (
                <Tooltip content={map[params.row.curNode] || params.row.curNode} placement="top">
                  <Tag>{map[params.row.curNode] || params.row.curNode || '-'}</Tag>
                </Tooltip>
              )
            } else {
              return <span>-</span>
            }
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
          minWidth: 170,
          key: 'handler',
          render: (h, params) => {
            let handlerArr = []
            let roleArr = []
            if (this.tabName === 'draft' || (this.tabName === 'submit' && ['1', '3'].includes(this.rollback))) {
              handlerArr = params.row.checkHandler.split(',') || []
              roleArr = params.row.checkHandleRoleDisplay.split(',') || []
            } else {
              handlerArr = params.row.handler.split(',') || []
              roleArr = params.row.handleRoleDisplay.split(',') || []
            }
            return (
              <div style="display:flex;flex-direction:column;">
                {handlerArr.map((item, index) => {
                  return <span>{`${roleArr[index] || '-'} / ${item || '-'}`}</span>
                })}
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
            return (
              <span>
                {['myPending', 'pending', 'hasProcessed'].includes(this.tabName)
                  ? this.$t('tw_taskTime_progress')
                  : this.$t('tw_requestTime_progress')}
              </span>
            )
          },
          minWidth: 140,
          key: 'effectiveDays',
          render: (h, params) => {
            let stayTime = '' // 停留时长
            let totalTime = '' // 预期停留时长
            if (['myPending', 'pending', 'hasProcessed'].includes(this.tabName)) {
              stayTime = params.row.taskStayTime
              totalTime = params.row.taskStayTimeTotal
            } else {
              stayTime = params.row.requestStayTime
              totalTime = params.row.requestStayTimeTotal
            }
            const percent = (stayTime / totalTime) * 100
            const color = percent > 50 ? (percent > 80 ? '#FF4D4F' : '#F29360') : '#00CB91'
            return (
              <div>
                <Progress stroke-color={color} percent={percent > 100 ? 100 : percent}>
                  <span>{`${stayTime}${this.$t('tw_days')}/${totalTime}${this.$t('tw_days')}`}</span>
                </Progress>
                {percent > 100 && (
                  <span style="color:#FF4D4F;display:flex;align-items:center;">
                    <Icon type="md-warning" color="#FF4D4F" />
                    {`${this.$t('tw_exceed')}${(stayTime - totalTime).toFixed(1)}${this.$t('tw_days')}`}
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
            if (params.row.templateName) {
              return (
                <span>
                  {`${params.row.templateName}`}
                  {params.row.version && (
                    <span style="border:1px solid #e8eaec;border-radius:3px;background:#f7f7f7;padding:1px 4px;">
                      {params.row.version}
                    </span>
                  )}
                </span>
              )
            } else {
              return <span>-</span>
            }
          }
        },
        procDefName: {
          title: this.$t('tw_template_flow'),
          sortable: 'custom',
          minWidth: 180,
          key: 'procDefName',
          render: (h, params) => {
            if (params.row.procDefName) {
              return (
                <span>
                  {`${params.row.procDefName}`}
                  {params.row.procDefVersion && (
                    <span style="border:1px solid #e8eaec;border-radius:3px;background:#f7f7f7;padding:1px 4px;">
                      {params.row.procDefVersion}
                    </span>
                  )}
                </span>
              )
            } else {
              return <span>-</span>
            }
          }
        },
        operatorObjType: {
          title: this.$t('tw_operator_type'),
          resizable: true,
          sortable: 'custom',
          minWidth: 150,
          key: 'operatorObjType',
          render: (h, params) => {
            if (params.row.operatorObjType) {
              return (
                params.row.operatorObjType && (
                  <Tooltip content={params.row.operatorObjType} placement="top">
                    <Tag>{params.row.operatorObjType}</Tag>
                  </Tooltip>
                )
              )
            } else {
              return <span>-</span>
            }
          }
        },
        operatorObj: {
          title: this.$t('tw_operator'),
          resizable: true,
          sortable: 'custom',
          minWidth: 150,
          key: 'operatorObj',
          render: (h, params) => {
            return <span>{params.row.operatorObj || '-'}</span>
          }
        },
        createdBy: {
          title: this.$t('tw_reporter'),
          sortable: 'custom',
          minWidth: 170,
          key: 'createdBy',
          render: (h, params) => {
            return <span>{`${params.row.roleDisplay} / ${params.row.createdBy}`}</span>
          }
        },
        // 关联单
        requestRefId: {
          title: this.$t('tw_ref'),
          width: 230,
          key: 'requestRefId',
          render: (h, params) => {
            const { requestRefName, requestRefId, requestRefType } = params.row
            if (requestRefId) {
              return (
                <span
                  style="cursor:pointer;color:#5cadff;"
                  onClick={() => {
                    this.handleViewRefDetail(params.row)
                  }}
                >
                  {requestRefName ? `【${this.typeMap[requestRefType]}】${requestRefName}【${requestRefId}】` : '-'}
                </span>
              )
            } else {
              return <span>-</span>
            }
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
                  ['myPending', 'pending', 'hasProcessed', 'submit'].includes(this.tabName) && (
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
                  ['Pending', 'InProgress', 'InApproval', 'Confirm'].includes(params.row.status) &&
                  ['myPending', 'pending'].includes(this.tabName) && (
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
                  ['Pending', 'InProgress', 'InApproval'].includes(params.row.status) &&
                  ['myPending', 'pending'].includes(this.tabName) &&
                  // 模板指定/提交人指定，该提交人角色的管理员可以展示认领按钮
                  ((['template', 'custom'].includes(params.row.handlerType) &&
                    params.row.roleAdministrator === this.username) ||
                    !['template', 'custom'].includes(params.row.handlerType)) && (
                    <Tooltip content={this.$t('tw_action_claim')} placement="top">
                      <Button
                        type="primary"
                        size="small"
                        onClick={() => {
                          this.handleTransfer(params.row, 'claim')
                        }}
                      >
                        <Icon type="md-person" size="16"></Icon>
                      </Button>
                    </Tooltip>
                  )}
                {// 转给我
                  params.row.handler &&
                  this.username !== params.row.handler &&
                  ['Pending', 'InProgress', 'InApproval'].includes(params.row.status) &&
                  ['myPending', 'pending'].includes(this.tabName) &&
                  // 模板指定/提交人指定，该提交人角色的管理员可以展示转给我按钮
                  ((['template', 'custom'].includes(params.row.handlerType) &&
                    params.row.roleAdministrator === this.username) ||
                    !['template', 'custom'].includes(params.row.handlerType)) && (
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
                  ['Termination', 'Faulted'].includes(params.row.status) && this.tabName === 'submit' && (
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
                // 我提交的定版或审批中revokeBtn为true可退回
                  ['Pending', 'InApproval'].includes(params.row.status) &&
                  params.row.revokeBtn &&
                  this.tabName === 'submit' && (
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
                        style="margin-left:5px;"
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
      hasProcessedTaskColumn: [],
      submitAllColumn: [],
      submitColumn: [],
      draftColumn: [],
      username: window.localStorage.getItem('username'),
      typeMap: {
        1: this.$t('tw_publish'),
        2: this.$t('tw_request'),
        3: this.$t('tw_question'),
        4: this.$t('tw_event'),
        5: this.$t('fork')
      },
      createRouteMap: {
        '1': 'createPublish',
        '2': 'createRequest',
        '3': 'createProblem',
        '4': 'createEvent',
        '5': 'createChange'
      },
      detailRouteMap: {
        '1': 'detailPublish',
        '2': 'detailRequest',
        '3': 'detailProblem',
        '4': 'detailEvent',
        '5': 'detailChange'
      }
    }
  },
  mounted () {
    // 本人/本组处理
    this.pendingTaskColumn = [
      this.baseColumn.id,
      this.baseColumn.name,
      this.baseColumn.status,
      {
        title: this.$t('task_name'),
        sortable: 'custom',
        minWidth: 120,
        key: 'taskName',
        render: (h, params) => {
          const taskNameMap = {
            check: this.$t('tw_pending_tab'),
            confirm: this.$t('tw_request_confirm')
          }
          return (
            <span
              style="cursor:pointer;"
              onClick={() => {
                this.handleDbClick(params.row)
              }}
            >
              {taskNameMap[params.row.taskName] || params.row.taskName || '-'}
            </span>
          )
        }
      },
      {
        title: this.$t('tw_taskUpdated'),
        sortable: 'custom',
        minWidth: 150,
        key: 'taskHandleUpdatedTime'
      },
      {
        title: this.$t('tw_taskEnd'),
        sortable: 'custom',
        minWidth: 150,
        key: 'taskExpectTime'
      },
      this.baseColumn.effectiveDays,
      this.baseColumn.progress,
      this.baseColumn.curNode,
      this.baseColumn.handler,
      this.baseColumn.createdBy,
      this.baseColumn.templateName,
      this.baseColumn.procDefName,
      this.baseColumn.operatorObjType,
      this.baseColumn.operatorObj,
      this.baseColumn.requestRefId,
      this.baseColumn.action
    ]

    // 已处理
    this.hasProcessedTaskColumn = [
      this.baseColumn.id,
      this.baseColumn.name,
      this.baseColumn.status,
      this.baseColumn.curNode,
      {
        title: this.$t('task_name'),
        sortable: 'custom',
        minWidth: 120,
        key: 'taskName',
        render: (h, params) => {
          const taskNameMap = {
            check: this.$t('tw_pending_tab'),
            confirm: this.$t('tw_request_confirm')
          }
          return (
            <span
              style="cursor:pointer;"
              onClick={() => {
                this.handleDbClick(params.row)
              }}
            >
              {taskNameMap[params.row.taskName] || params.row.taskName || '-'}
            </span>
          )
        }
      },
      {
        title: this.$t('tw_taskCreated'),
        sortable: 'custom',
        minWidth: 150,
        key: 'taskCreatedTime'
      },
      {
        title: this.$t('tw_taskEnd'),
        sortable: 'custom',
        minWidth: 150,
        key: 'taskExpectTime'
      },
      {
        title: this.$t('handle_time'),
        sortable: 'custom',
        minWidth: 150,
        key: 'taskApprovalTime'
      },
      this.baseColumn.effectiveDays,
      this.baseColumn.progress,
      this.baseColumn.handler,
      this.baseColumn.createdBy,
      {
        title: this.$t('tw_rollback_reason'),
        isShow: this.type === 1,
        sortable: 'custom',
        minWidth: 150,
        key: 'rollbackDesc',
        render: (h, params) => {
          return <BaseEllipsis content={params.row.rollbackDesc}></BaseEllipsis>
        }
      },
      this.baseColumn.templateName,
      this.baseColumn.procDefName,
      this.baseColumn.operatorObjType,
      this.baseColumn.operatorObj,
      this.baseColumn.requestRefId,
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
          return <BaseEllipsis content={params.row.rollbackDesc}></BaseEllipsis>
        }
      },
      this.baseColumn.templateName,
      this.baseColumn.procDefName,
      this.baseColumn.operatorObjType,
      this.baseColumn.operatorObj,
      this.baseColumn.requestRefId,
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
      this.baseColumn.requestRefId,
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
      this.baseColumn.requestRefId,
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
