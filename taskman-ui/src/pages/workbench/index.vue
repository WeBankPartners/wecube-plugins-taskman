<template>
  <div class="workbench">
    <!-- <div class="header">
      <div class="left">
        <Icon size="32" type="md-desktop" />
        <span class="ml">个人工作台</span>
      </div>
      <div class="right">
        <Icon v-if="expand" @click="expand = !expand" size="30" type="ios-arrow-down" />
        <Icon v-else @click="expand = !expand" size="30" type="ios-arrow-up" />
        <Icon @click="expand = !expand"  class='ml' size="32" type="ios-refresh" />
      </div>
    </div> -->
    <div class="hot-link">
      <HotLink></HotLink>
    </div>
    <div class="data-card">
      <DataCard ref="dataCard" :parent-action="actionName" @fetchData="handleOverviewChange"></DataCard>
    </div>
    <div class="data-tabs">
      <!-- <Tabs v-model="actionName">
        <TabPane label="发布" name="1"></TabPane>
        <TabPane label="请求" name="2"></TabPane>
      </Tabs> -->
      <CollectTable
        ref="collect"
        v-if="tabName === 'collect'"
        :getTemplateList="getTemplateList"
        :actionName="actionName"
      ></CollectTable>
      <template v-else>
        <!--搜索条件-->
        <BaseSearch :options="searchOptions" v-model="form" @search="handleQuery"></BaseSearch>
        <!--表格分页-->
        <Table
          :border="false"
          size="small"
          :loading="loading"
          :columns="tableColumns"
          :data="tableData"
          @on-sort-change="sortTable"
        ></Table>
        <Page
          style="float:right;margin-top:10px;"
          :total="pagination.total"
          @on-change="changPage"
          show-sizer
          :current="pagination.currentPage"
          :page-size="pagination.pageSize"
          @on-page-size-change="changePageSize"
          show-total
        />
      </template>
    </div>
  </div>
</template>

<script>
import WorkBench from '@/pages/components/workbench-menu.vue'
import HotLink from './components/hot-link.vue'
import DataCard from './components/data-card.vue'
import BaseSearch from '../components/base-search.vue'
import CollectTable from './collect-table.vue'
import { getPlatformList, getTemplateList, tansferToMe, changeTaskStatus, deleteRequest, reRequest } from '@/api/server'
import dayjs from 'dayjs'
export default {
  components: {
    HotLink,
    DataCard,
    BaseSearch,
    WorkBench,
    CollectTable
  },
  data () {
    return {
      tabName: 'pending', // pending待处理,hasProcessed已处理,submit我提交的,draft我的暂存,collect收藏
      actionName: '1', // 1发布,2请求(3问题,4事件,5变更)
      form: {
        type: 0, // 0所有,1请求定版,2任务处理
        rollback: 0, // 0所有,1已退回,2其他
        name: '', // ID或名称模糊搜索
        id: '',
        templateId: [], // 模板ID
        status: [], // 状态
        operatorObj: '', // 操作对象
        createdBy: [] // 创建人
      },
      searchOptions: [
        {
          key: 'type',
          initValue: 0,
          hidden: false,
          component: 'radio-group',
          list: [
            { label: '所有', value: 0 },
            { label: '请求定版', value: 1 },
            { label: '任务处理', value: 2 }
          ]
        },
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
          key: 'operatorObj',
          placeholder: '操作对象',
          component: 'input'
        },
        {
          key: 'templateId',
          placeholder: '模板',
          multiple: true,
          component: 'remote-select',
          remote: this.getTemplateList
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
        }
      ],
      tableColumns: [
        {
          title: '请求ID',
          width: 150,
          key: 'id'
        },
        {
          title: this.$t('name'),
          sortable: 'custom',
          minWidth: 250,
          key: 'name'
        },
        {
          title: '使用模板',
          sortable: 'custom',
          minWidth: 180,
          key: 'templateName'
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
          title: '使用编排',
          minWidth: 150,
          key: 'procDefName'
        },
        {
          title: '请求状态',
          sortable: 'custom',
          key: 'status',
          minWidth: 120,
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
                      ? `${item.label}（被退回）`
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
              sendRequest: '提起请求',
              requestPending: '请求定版',
              requestComplete: '请求完成'
            }
            return <span>{map[params.row.curNode] || params.row.curNode}</span>
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
          renderHeader: () => {
            return <span>{this.form.type === 2 ? '任务停留时长' : '请求停留时长'}</span>
          },
          minWidth: 200,
          key: 'expectTime',
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
          title: '创建人',
          sortable: 'custom',
          minWidth: 160,
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
          title: '当前处理人',
          sortable: 'custom',
          minWidth: 160,
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
          title: '创建时间',
          sortable: 'custom',
          minWidth: 150,
          key: 'createdTime'
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
          key: 'rollbackDesc',
          minWidth: 200
        },
        {
          title: this.$t('t_action'),
          key: 'action',
          minWidth: 160,
          fixed: 'right',
          align: 'center',
          render: (h, params) => {
            return (
              <div>
                {['pending', 'hasProcessed', 'submit'].includes(this.tabName) && (
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
                {this.username === params.row.handler && ['Pending', 'InProgress'].includes(params.row.status) && (
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
                {!params.row.handler && ['Pending', 'InProgress'].includes(params.row.status) && (
                  <Button
                    type="info"
                    size="small"
                    onClick={() => {
                      this.handleTransfer(params.row)
                    }}
                  >
                    认领
                  </Button>
                )}
                {params.row.handler &&
                  this.username !== params.row.handler &&
                  ['Pending', 'InProgress'].includes(params.row.status) && (
                  <Button
                    type="success"
                    size="small"
                    onClick={() => {
                      this.handleTransfer(params.row)
                    }}
                  >
                      转给我
                  </Button>
                )}
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
                {this.tabName === 'draft' && (
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
      ],
      tableData: [],
      loading: false,
      pagination: {
        total: 0,
        currentPage: 1,
        pageSize: 10
      },
      username: window.localStorage.getItem('username')
    }
  },
  watch: {
    // 切换tab视图
    tabName: {
      handler (val) {
        if (val) {
          this.$nextTick(() => {
            this.searchOptions[0].hidden = false
            // 待处理、进行中
            if (['pending', 'hasProcessed'].includes(val)) {
              this.form.type = 2
              this.form.rollback = 0
              this.searchOptions[0].key = 'type'
              this.searchOptions[0].initValue = 2
              this.searchOptions[0].list = [
                { label: '任务处理', value: 2 },
                { label: '请求定版', value: 1 }
              ]
              // 我提交的
            } else if (val === 'submit') {
              this.form.rollback = 0
              this.form.type = 0
              this.searchOptions[0].key = 'rollback'
              this.searchOptions[0].initValue = 0
              this.searchOptions[0].list = [
                { label: '所有', value: 0 },
                { label: '被退回', value: 1 },
                { label: '其他', value: 2 }
              ]
            } else if (val === 'draft') {
              this.form.type = 0
              this.form.rollback = 0
              this.searchOptions[0].hidden = true
            }
            if (val !== 'collect') {
              this.handleQuery()
            } else {
              this.$nextTick(() => {
                this.$refs.collect.handleQuery()
              })
            }
          })
        }
      },
      immediate: true
    },
    // 切换行为action
    actionName () {
      if (this.tabName !== 'collect') {
        this.handleQuery()
      } else {
        this.$nextTick(() => {
          this.$refs.collect.handleQuery()
        })
      }
    }
  },
  mounted () {
    // this.tabName = this.$route.query.tabName
    // this.actionName = this.$route.query.actionName
    // this.getList()
  },
  methods: {
    // 点击视图卡片触发查询
    handleOverviewChange (val, action) {
      this.tabName = val
      this.actionName = action || '1'
    },
    // 表格排序
    sortTable (col) {
      const sorting = {
        asc: col.order === 'asc',
        field: col.key
      }
      this.getList(sorting)
    },
    async getList (sort = { asc: false, field: 'updatedTime' }) {
      this.loading = true
      if (this.form.type === '') {
        this.form.type = 0
      }
      if (this.form.rollback === '') {
        this.form.rollback = 0
      }
      const params = {
        tab: this.tabName,
        action: Number(this.actionName),
        ...this.form,
        startIndex: (this.pagination.currentPage - 1) * this.pagination.pageSize,
        pageSize: this.pagination.pageSize
      }
      if (sort) {
        params.sorting = sort
      }
      const { statusCode, data } = await getPlatformList(params)
      if (statusCode === 'OK') {
        this.tableData = data.contents || []
        this.pagination.total = data.pageInfo.totalRows
      }
      this.loading = false
    },
    handleQuery () {
      this.pagination.currentPage = 1
      this.getList()
      this.$refs.dataCard.getData()
    },
    changPage (val) {
      this.pagination.currentPage = val
      this.getList()
    },
    changePageSize (val) {
      this.pagination.currentPage = 1
      this.pagination.pageSize = val
      this.getList()
    },
    // 获取下拉搜索模板列表
    async getTemplateList () {
      const params = {
        filters: [],
        paging: false
      }
      const { statusCode, data } = await getTemplateList(params)
      if (statusCode === 'OK') {
        const options = data.contents || []
        return options.map(option => {
          return {
            label: `${option.name}(${option.version || '-'})`,
            value: option.id
          }
        })
      }
    },
    // 获取下拉列表的值
    // async getFilterOptions() {
    //   await getPlatformFilter({ startTime: '' })
    // },
    // 表格操作-查看
    hanldeView (row) {
      const path = this.actionName === '1' ? 'createPublish' : 'createRequest'
      const url = `/taskman/workbench/${path}`
      this.$router.push({
        path: url,
        query: {
          requestId: row.id,
          requestTemplate: row.templateId,
          isAdd: 'N',
          isCheck: 'Y',
          isHandle: 'N',
          enforceDisable: 'Y',
          jumpFrom: 'group_handle'
        }
      })
    },
    // 表格操作-处理(任务处理和请求定版)
    handleEdit (row) {
      const path = this.actionName === '1' ? 'createPublish' : 'createRequest'
      const url = `/taskman/workbench/${path}`
      this.$router.push({
        path: url,
        query: {
          requestId: row.id,
          requestTemplate: row.templateId,
          isAdd: 'N',
          isCheck: 'Y',
          isHandle: 'Y',
          enforceDisable: 'N',
          jumpFrom: 'group_handle'
        }
      })
    },
    // 表格操作-转给我
    async handleTransfer (row) {
      this.$Modal.confirm({
        title: this.$t('confirm') + '转给我',
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          // 请求定版的新接口，任务处理的老接口
          let res = null
          if (row.status === 'Pending') {
            res = await tansferToMe(row.id)
          } else if (row.status === 'InProgress') {
            res = await changeTaskStatus('give', row.id)
          }
          if (res.statusCode === 'OK') {
            this.$Notice.success({
              title: this.$t('successful'),
              desc: this.$t('successful')
            })
            this.getList()
          }
        },
        onCancel: () => {}
      })
    },
    // 表格操作-重新发起
    async handleRepub (row) {
      const { statusCode, data } = await reRequest(row.id)
      if (statusCode === 'OK') {
        const path = this.actionName === '1' ? 'createPublish' : 'createRequest'
        const url = `/taskman/workbench/${path}`
        this.$router.push({
          path: url,
          query: {
            requestId: data.id,
            requestTemplate: data.requestTemplate,
            isAdd: 'Y',
            isCheck: 'N',
            isHandle: 'N',
            jumpFrom: 'my_drafts'
          }
        })
      }
    },
    // 表格操作-草稿去发起
    hanldeLaunch (row) {
      const path = this.actionName === '1' ? 'createPublish' : 'createRequest'
      const url = `/taskman/workbench/${path}`
      this.$router.push({
        path: url,
        query: {
          requestId: row.id,
          requestTemplate: row.templateId,
          isAdd: 'Y',
          isCheck: 'N',
          isHandle: 'N',
          jumpFrom: 'my_drafts'
        }
      })
    },
    // 删除草稿
    handleDeleteDraft (row) {
      this.$Modal.confirm({
        title: this.$t('confirm_delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          let res = await deleteRequest(row.id)
          if (res.statusCode === 'OK') {
            this.$Notice.success({
              title: this.$t('successful'),
              desc: this.$t('successful')
            })
            this.getList()
          }
        },
        onCancel: () => {}
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.workbench {
  .header {
    display: flex;
    justify-content: space-between;
    align-content: center;
    .left,
    .right {
      display: flex;
      align-items: center;
      span {
        font-size: 16px;
        font-weight: 600;
      }
      .ml {
        margin-left: 10px;
      }
    }
  }
  .hot-link {
    margin-top: 12px;
  }
  .data-card {
    margin-top: 24px;
  }
  .data-tabs {
    margin-top: 24px;
  }
}
</style>
<style lang="scss">
.workbench {
  .ivu-progress-outer {
    width: 90px !important;
    padding-right: 30px !important;
    margin-right: -33px !important;
  }
  .ivu-progress-inner {
    width: 60px;
  }
  .ivu-progress-text {
    color: #515a6e !important;
    min-width: 80px !important;
  }
  .ivu-progress {
    display: flex;
  }
}
</style>
