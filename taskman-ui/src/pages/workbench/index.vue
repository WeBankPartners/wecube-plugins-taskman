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
      <DataCard ref="dataCard" :initTab="initTab" :initAction="initAction" @fetchData="handleOverviewChange"></DataCard>
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
          :columns="tableColumn"
          :data="tableData"
          @on-sort-change="sortTable"
        >
        </Table>
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
import {
  getPlatformList,
  getTemplateList,
  tansferToMe,
  recallRequest,
  changeTaskStatus,
  deleteRequest,
  reRequest,
  getPlatformFilter
} from '@/api/server'
import { deepClone } from '@/pages/util/index'
import column from './column.js'
export default {
  components: {
    HotLink,
    DataCard,
    BaseSearch,
    WorkBench,
    CollectTable
  },
  mixins: [column],
  data () {
    return {
      tabName: 'pending', // pending待处理,hasProcessed已处理,submit我提交的,draft我的暂存,collect收藏
      actionName: '1', // 1发布,2请求(3问题,4事件,5变更)
      initTab: '',
      initAction: '',
      form: {
        type: 0, // 0所有,1请求定版,2任务处理
        rollback: 0, // 0所有,1已退回,2其他,3被撤回
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
      filterOptions: [],
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
      tableColumn: [],
      tableData: [],
      loading: false,
      pagination: {
        total: 0,
        currentPage: 1,
        pageSize: 10
      },
      sorting: {}, // 表格默认排序
      username: window.localStorage.getItem('username')
    }
  },
  watch: {
    'form.type': {
      handler (val) {
        if (this.tabName === 'pending') {
          if (val === 1) {
            this.tableColumn = this.pendingColumn
          } else if (val === 2) {
            this.tableColumn = this.pendingTaskColumn
          }
        } else if (this.tabName === 'hasProcessed') {
          if (val === 1) {
            this.tableColumn = this.hasProcessedColumn
          } else if (val === 2) {
            this.tableColumn = this.hasProcessedTaskColumn
          }
        }
      },
      immediate: true
    },
    'form.rollback': {
      handler (val) {
        if (this.tabName === 'submit') {
          if (val === 1 || val === 0) {
            this.tableColumn = this.submitAllColumn
          } else if (val === 2 || val === 3) {
            this.tableColumn = this.submitColumn
          }
        }
      },
      immediate: true
    }
  },
  mounted () {
    this.initTab = this.$route.query.tabName || 'pending'
    this.initAction = this.$route.query.actionName || '1'
    this.getFilterOptions()
  },
  methods: {
    // 点击视图卡片触发查询
    handleOverviewChange (val, action) {
      this.tabName = val
      this.actionName = action || '1'

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
        // 表格列及排序初始化
        this.tableColumn = this.submitAllColumn
        this.form.type = 0
        this.searchOptions[0].key = 'rollback'
        this.searchOptions[0].initValue = 0
        this.searchOptions[0].list = [
          { label: '所有', value: 0 },
          { label: '被退回', value: 1 },
          { label: '本人撤回', value: 3 },
          { label: '其他', value: 2 }
        ]
      } else if (val === 'draft') {
        this.form.type = 0
        this.form.rollback = 0
        // 表格列及排序初始化
        this.tableColumn = this.draftColumn
        this.searchOptions[0].initValue = 0
        this.searchOptions[0].hidden = true
      }

      if (val !== 'collect') {
        this.handleQuery()
      } else {
        this.$nextTick(() => {
          this.$refs.collect.handleQuery()
          this.$refs.dataCard.getData()
        })
      }
    },
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
      if (form.type === '') {
        form.type = 0
      }
      if (form.rollback === '') {
        form.rollback = 0
      }
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
        tab: this.tabName,
        action: Number(this.actionName),
        ...form,
        startIndex: (this.pagination.currentPage - 1) * this.pagination.pageSize,
        pageSize: this.pagination.pageSize
      }
      // 获取默认排序
      if (this.tabName === 'pending') {
        if (this.form.type === 1) {
          this.sorting = {
            asc: false,
            field: 'reportTime'
          }
        } else if (this.form.type === 2) {
          this.sorting = {
            asc: false,
            field: 'taskCreatedTime'
          }
        }
      } else if (this.tabName === 'hasProcessed') {
        if (this.form.type === 1) {
          this.sorting = {
            asc: false,
            field: 'approvalTime'
          }
        } else if (this.form.type === 2) {
          this.sorting = {
            asc: false,
            field: 'taskApprovalTime'
          }
        }
      } else if (this.tabName === 'submit') {
        this.sorting = {
          asc: false,
          field: 'reportTime'
        }
      } else if (this.tabName === 'draft') {
        this.sorting = {
          asc: false,
          field: 'updatedTime'
        }
      }
      if (this.sorting) {
        params.sorting = this.sorting
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
    // 获取搜索条件的下拉值
    async getFilterOptions () {
      const { statusCode, data } = await getPlatformFilter({ startTime: '' })
      if (statusCode === 'OK') {
        this.filterOptions = data
        this.searchOptions.forEach(item => {
          if (item.key === 'operatorObjType') {
            item.list =
              data.operatorObjTypeList &&
              data.operatorObjTypeList.map(item => {
                return {
                  label: item,
                  value: item
                }
              })
          } else if (item.key === 'templateId') {
            item.list =
              data.templateList &&
              data.templateList.map(item => {
                return {
                  label: item.templateName,
                  value: item.templateId
                }
              })
          } else if (item.key === 'procDefName') {
            item.list =
              data.procDefNameList &&
              data.procDefNameList.map(item => {
                return {
                  label: item,
                  value: item
                }
              })
          } else if (item.key === 'handler') {
            item.list =
              data.handlerList &&
              data.handlerList.map(item => {
                return {
                  label: item,
                  value: item
                }
              })
          } else if (item.key === 'createdBy') {
            item.list =
              data.createdByList &&
              data.createdByList.map(item => {
                return {
                  label: item,
                  value: item
                }
              })
          }
        })
      }
    },
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
          jumpFrom: this.tabName
        }
      })
    },
    // 表格操作-处理(任务处理和请求定版)
    async handleEdit (row) {
      // 处理任务需要更新任务状态
      if (row.status === 'InProgress') {
        await changeTaskStatus('start', row.taskId)
      }
      const path = this.actionName === '1' ? 'createPublish' : 'createRequest'
      const url = `/taskman/workbench/${path}`
      this.$router.push({
        path: url,
        query: {
          requestId: row.id,
          requestTemplate: row.templateId,
          isAdd: 'N',
          isCheck: 'N',
          isHandle: 'Y',
          enforceDisable: 'N',
          jumpFrom: 'group_handle'
        }
      })
    },
    // 表格操作-转给我/认领
    async handleTransfer (row, type) {
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
            res = await changeTaskStatus(type, row.taskId)
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
    // 表格操作撤回
    async handleRecall (row) {
      this.$Modal.confirm({
        title: this.$t('confirm') + '撤回',
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const { statusCode } = await recallRequest(row.id)
          if (statusCode === 'OK') {
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
    },
    handleDbClick (row) {
      if (
        this.username === row.handler &&
        ['Pending', 'InProgress'].includes(row.status) &&
        this.tabName === 'pending'
      ) {
        this.handleEdit(row)
      } else if (row.status === 'Draft') {
        this.hanldeLaunch(row)
      } else {
        this.hanldeView(row)
      }
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
}
</style>
