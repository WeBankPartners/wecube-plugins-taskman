<template>
  <div class="workbench-request-history">
    <Tabs :value="activeTab" @on-click="handleChangeTab">
      <TabPane label="已发布" name="commit"></TabPane>
      <TabPane label="草稿箱" name="draft"></TabPane>
    </Tabs>
    <BaseSearch :options="searchOptions" v-model="form" @search="handleQuery"></BaseSearch>
    <Table
      size="small"
      :columns="tableColumns"
      :data="tableData"
      :loading="loading"
      @on-sort-change="sortTable"
    ></Table>
    <Page
      style="float:right;margin-top:10px;"
      :total="pagination.total"
      @on-change="handlePage"
      show-sizer
      :current="pagination.currentPage"
      :page-size="pagination.pageSize"
      @on-page-size-change="handlePageSize"
      show-total
    />
  </div>
</template>

<script>
import BaseSearch from '@/pages/components/base-search.vue'
import { getPublishList, reRequest } from '@/api/server'
import { deepClone } from '@/pages/util/index'
import dayjs from 'dayjs'
export default {
  components: {
    BaseSearch
  },
  data () {
    return {
      activeTab: 'commit',
      form: {
        name: '', // ID或名称模糊搜索
        id: '',
        templateId: [], // 模板ID
        status: [], // 状态
        operatorObjType: [], // 操作对象类型
        procDefName: [], // 使用编排
        createdBy: [], // 创建人
        handler: [], // 当前处理人
        createdTime: [],
        updatedTime: [],
        expectTime: []
      },
      tableData: [],
      loading: false,
      sorting: {
        asc: false,
        field: 'reportTime'
      }, // 表格默认排序
      pagination: {
        total: 0,
        currentPage: 1,
        pageSize: 10
      },
      searchOptions: [
        {
          key: 'name',
          placeholder: '名称',
          component: 'input'
        },
        {
          key: 'id',
          placeholder: 'id',
          component: 'input'
        },
        {
          key: 'operatorObjType',
          placeholder: '操作对象类型',
          multiple: true,
          component: 'select',
          list: []
        },
        {
          key: 'templateId',
          placeholder: '模板',
          multiple: true,
          component: 'select',
          list: []
        },
        {
          key: 'procDefName',
          placeholder: '使用编排',
          multiple: true,
          component: 'select',
          list: []
        },
        {
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
        {
          key: 'createdBy',
          placeholder: '创建人',
          component: 'select',
          multiple: true,
          list: []
        },
        {
          key: 'handler',
          placeholder: '处理人',
          component: 'select',
          multiple: true,
          list: []
        },
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
          label: '期望时间',
          dateType: 4,
          labelWidth: 85,
          component: 'custom-time'
        }
      ],
      tableColumns: [
        {
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
        {
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
        {
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
        {
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
        {
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
        {
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
        {
          title: '请求停留时长',
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
        {
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
        {
          title: '使用编排',
          sortable: 'custom',
          minWidth: 150,
          key: 'procDefName'
        },
        {
          title: '操作对象类型',
          resizable: true,
          sortable: 'custom',
          minWidth: 150,
          key: 'operatorObjType',
          render: (h, params) => {
            return params.row.operatorObjType && <Tag>{params.row.operatorObjType}</Tag>
          }
        },
        {
          title: '操作对象',
          resizable: true,
          sortable: 'custom',
          minWidth: 150,
          key: 'operatorObj'
        },
        {
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
        {
          title: '请求提交时间',
          sortable: 'custom',
          minWidth: 150,
          key: 'reportTime'
        },
        {
          title: this.$t('t_action'),
          key: 'action',
          width: 160,
          fixed: 'right',
          align: 'center',
          render: (h, params) => {
            return (
              <div>
                <Button
                  size="small"
                  onClick={() => {
                    this.hanldeView(params.row)
                  }}
                  style="margin-right:5px;"
                >
                  查看
                </Button>
                {['Termination', 'Completed', 'Faulted'].includes(params.row.status) && (
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
                {
                  // params.row.status === 'Pending' && this.tabName === 'submit' && (
                  // <Button
                  //   type="error"
                  //   size="small"
                  //   onClick={() => {
                  //     this.handleRecall(params.row)
                  //   }}
                  // >
                  //   撤回
                  // </Button>
                  // )
                }
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
              </div>
            )
          }
        }
      ]
    }
  },
  mounted () {
    this.getList()
  },
  methods: {
    // 表格排序
    sortTable (col) {
      this.sorting = {
        asc: col.order === 'asc',
        field: col.key
      }
      this.getList()
    },
    async getList () {
      this.loading = true
      const form = deepClone(this.form)
      const dateTransferArr = ['createdTime', 'updatedTime', 'expectTime']
      dateTransferArr.forEach(item => {
        if (form[item] && form[item].length > 0) {
          form[item + 'Start'] = form[item][0] + ' 00:00:00'
          form[item + 'End'] = form[item][1] + ' 23:59:59'
          delete form[item]
        } else {
          form[item + 'Start'] = ''
          form[item + 'End'] = ''
          delete form[item]
        }
      })
      const params = {
        tab: this.activeTab,
        action: 2,
        ...this.form,
        startIndex: (this.pagination.currentPage - 1) * this.pagination.pageSize,
        pageSize: this.pagination.pageSize
      }
      if (this.sorting) {
        params.sorting = this.sorting
      }
      const { statusCode, data } = await getPublishList(params)
      if (statusCode === 'OK') {
        this.tableData = data.contents || []
        this.pagination.total = data.pageInfo.totalRows
      }
      this.loading = false
    },
    handleQuery () {
      this.pagination.currentPage = 1
      this.getList()
    },
    handlePage (val) {
      this.pagination.currentPage = val
      this.getList()
    },
    handlePageSize (val) {
      this.pagination.currentPage = 1
      this.pagination.pageSize = val
      this.getList()
    },
    handleChangeTab (val) {
      this.activeTab = val
      this.handleQuery()
    },
    handleDbClick (row) {},
    // 表格操作-查看
    hanldeView (row) {
      const url = `/taskman/workbench/createRequest`
      this.$router.push({
        path: url,
        query: {
          requestId: row.id,
          requestTemplate: row.templateId,
          isAdd: 'N',
          isCheck: 'Y',
          isHandle: 'N',
          enforceDisable: 'Y',
          jumpFrom: ''
        }
      })
    },
    // 表格操作-重新发起
    async handleRepub (row) {
      const { statusCode, data } = await reRequest(row.id)
      if (statusCode === 'OK') {
        const url = `/taskman/workbench/createRequest`
        this.$router.push({
          path: url,
          query: {
            requestId: data.id,
            requestTemplate: data.requestTemplate,
            isAdd: 'Y',
            isCheck: 'N',
            isHandle: 'N',
            jumpFrom: ''
          }
        })
      }
    },
    // 表格操作-草稿去发起
    hanldeLaunch (row) {
      const url = `/taskman/workbench/createRequest`
      this.$router.push({
        path: url,
        query: {
          requestId: row.id,
          requestTemplate: row.templateId,
          isAdd: 'Y',
          isCheck: 'N',
          isHandle: 'N',
          jumpFrom: ''
        }
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.workbench-request-history {
  width: 100%;
}
</style>
<style lang="scss">
.workbench-request-history {
  .ivu-progress-outer {
    width: 90px !important;
    padding-right: 30px !important;
    margin-right: -33px !important;
  }
  .ivu-progress-inner {
    width: 60px !important;
  }
  .ivu-progress-text {
    color: #515a6e !important;
    min-width: 80px !important;
  }
  .ivu-progress {
    display: flex;
    // flex-direction: column;
  }
  .ivu-btn-small {
    font-size: 12px;
  }
}
</style>
