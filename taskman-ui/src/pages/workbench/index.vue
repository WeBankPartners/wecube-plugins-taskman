<!--工作台-->
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
    <!-- <div class="hot-link">
      <HotLink></HotLink>
    </div> -->
    <div class="data-card">
      <DataCard
        ref="dataCard"
        :initTab="initTab"
        :initAction="initAction"
        @initFetch="initData"
        @fetchData="handleOverviewChange"
      ></DataCard>
    </div>
    <div class="data-tabs">
      <Tabs v-if="['pending', 'hasProcessed'].includes(tabName)" v-model="type" @on-click="handleTypeChange">
        <!--任务处理-->
        <TabPane :label="$t('tw_task_tab')" name="2"></TabPane>
        <!--请求定版-->
        <TabPane :label="$t('tw_pending_tab')" name="1"></TabPane>
      </Tabs>
      <Tabs v-if="['submit'].includes(tabName)" v-model="rollback" @on-click="handleRollbackChange">
        <!--所有-->
        <TabPane :label="$t('tw_all_tab')" name="0"></TabPane>
        <!--被退回-->
        <TabPane :label="$t('tw_return_tab')" name="1"></TabPane>
        <!--本人撤回-->
        <TabPane :label="$t('tw_recall_tab')" name="3"></TabPane>
        <!--其他-->
        <TabPane :label="$t('tw_other_tab')" name="2"></TabPane>
      </Tabs>
      <CollectTable v-if="tabName === 'collect'" ref="collect" :actionName="actionName"></CollectTable>
      <template v-else>
        <!--搜索条件-->
        <BaseSearch ref="search" :options="searchOptions" v-model="form" @search="handleQuery"></BaseSearch>
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
import { getPlatformList, tansferToMe, recallRequest, changeTaskStatus, deleteRequest, reRequest } from '@/api/server'
import { deepClone } from '@/pages/util/index'
import column from './column.js'
import search from './search.js'
export default {
  components: {
    HotLink,
    DataCard,
    BaseSearch,
    WorkBench,
    CollectTable
  },
  mixins: [column, search],
  data () {
    return {
      tabName: 'pending', // pending待处理,hasProcessed已处理,submit我提交的,draft我的暂存,collect收藏
      actionName: '1', // 1发布,2请求(3问题,4事件,5变更)
      initTab: '',
      initAction: '',
      type: '0', // 0所有,1请求定版,2任务处理
      rollback: '0', // 0所有,1已退回,2其他,3被撤回
      form: {
        name: '', // 请求名
        taskName: '', // 任务名
        id: '',
        templateId: [], // 模板ID
        status: [], // 状态
        operatorObjType: [], // 操作对象类型
        procDefName: [], // 使用编排
        createdBy: [], // 创建人
        createdTime: [],
        updatedTime: [],
        expectTime: [], // 期望时间
        reportTime: [], // 请求提交时间
        approvalTime: [], // 请求处理时间
        taskReportTime: [], // 任务提交时间
        taskApprovalTime: [], // 任务审批时间
        taskExpectTime: [] // 任务期望时间
      },
      searchOptions: [], // 表格筛选配置
      tableColumn: [],
      tableData: [],
      loading: false,
      pagination: {
        total: 0,
        currentPage: 1,
        pageSize: 10
      },
      sorting: {} // 表格默认排序
    }
  },
  watch: {
    // 切换请求发布，重新获取筛选条件配置
    actionName () {
      this.getFilterOptions()
    }
  },
  mounted () {
    this.initTab = this.$route.query.tabName || 'pending'
    this.initAction = this.$route.query.actionName || '1'
  },
  methods: {
    // 初始化加载数据(链接携带参数，跳转到指定标签)
    initData (val, action) {
      this.tabName = val
      this.actionName = action
      const type = this.$route.query.type
      const rollback = this.$route.query.rollback
      if (['pending', 'hasProcessed'].includes(val)) {
        if (['1', '2'].includes(type)) {
          this.type = type
        } else {
          this.type = '2'
        }
        this.rollback = ''
        this.getTypeConfig()
      } else if (val === 'submit') {
        if (['1', '2', '3'].includes(rollback)) {
          this.rollback = rollback
        } else {
          this.rollback = '0'
        }
        this.type = ''
        this.getRollbackConfig()
        this.tableColumn = this.submitAllColumn
        this.searchOptions = this.submitSearch
      } else if (val === 'draft') {
        this.tableColumn = this.draftColumn
        this.searchOptions = this.draftSearch
      }
      if (val !== 'collect') {
        this.handleReset()
        this.handleQuery()
      } else {
        this.$nextTick(() => {
          this.$refs.collect.handleQuery()
          this.$refs.dataCard.getData()
        })
      }
    },
    // 点击视图卡片触发查询
    handleOverviewChange (val, action) {
      this.tabName = val
      this.actionName = action
      if (['pending', 'hasProcessed'].includes(val)) {
        this.type = '2'
        this.rollback = ''
        this.getTypeConfig()
      } else if (val === 'submit') {
        this.rollback = '0'
        this.type = ''
        this.getRollbackConfig()
      } else if (val === 'draft') {
        this.tableColumn = this.draftColumn
        this.searchOptions = this.draftSearch
      }
      if (val !== 'collect') {
        this.handleReset()
        this.handleQuery()
      } else {
        this.$nextTick(() => {
          this.$refs.collect.handleQuery()
          this.$refs.dataCard.getData()
        })
      }
    },
    // 切换type
    handleTypeChange () {
      this.getTypeConfig()
      this.handleReset()
      this.handleQuery()
    },
    getTypeConfig () {
      if (this.tabName === 'pending') {
        if (this.type === '1') {
          this.tableColumn = this.pendingColumn
          this.searchOptions = this.pendingSearch
        } else if (this.type === '2') {
          this.tableColumn = this.pendingTaskColumn
          this.searchOptions = this.pendingTaskSearch
        }
      } else if (this.tabName === 'hasProcessed') {
        if (this.type === '1') {
          this.tableColumn = this.hasProcessedColumn
          this.searchOptions = this.hasProcessedSearch
        } else if (this.type === '2') {
          this.tableColumn = this.hasProcessedTaskColumn
          this.searchOptions = this.hasProcessedTaskSearch
        }
      }
    },
    // 切换rollback
    handleRollbackChange () {
      this.getRollbackConfig()
      this.handleReset()
      this.handleQuery()
    },
    getRollbackConfig () {
      if (this.tabName === 'submit') {
        if (this.rollback === '1' || this.rollback === '0') {
          this.tableColumn = this.submitAllColumn
          this.searchOptions = this.submitSearch
        } else if (this.rollback === '2' || this.rollback === '3') {
          this.tableColumn = this.submitColumn
          this.searchOptions = this.submitSearch
        }
      }
    },
    // 重置表单
    handleReset () {
      const resetObj = {}
      Object.keys(this.form).forEach(key => {
        if (Array.isArray(this.form[key])) {
          resetObj[key] = []
        } else {
          resetObj[key] = ''
        }
        // 处理时间类型默认值
        this.searchOptions.forEach(i => {
          if (i.component === 'custom-time' && i.initValue) {
            i.dateType = 1
          } else {
            i.dateType = 4
          }
        })
        // 点击清空按钮需要给默认值的表单选项
        const initOptions = this.searchOptions.filter(i => i.initValue !== undefined)
        initOptions.forEach(i => {
          resetObj[i.key] = i.initValue
        })
      })
      this.form = resetObj
    },
    // 表格排序
    sortTable (col) {
      this.sorting = {
        asc: col.order === 'asc',
        field: col.key
      }
      this.getList(true)
    },
    // 表格默认排序
    initSortTable () {
      if (this.tabName === 'pending') {
        if (this.type === '1') {
          this.sorting = {
            asc: false,
            field: 'reportTime'
          }
        } else if (this.type === '2') {
          this.sorting = {
            asc: false,
            field: 'taskCreatedTime'
          }
        }
      } else if (this.tabName === 'hasProcessed') {
        if (this.type === '1') {
          this.sorting = {
            asc: false,
            field: 'approvalTime'
          }
        } else if (this.type === '2') {
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
    },
    async getList (dynamicSort) {
      this.loading = true
      const form = deepClone(this.form)
      // 过滤掉多余时间
      var dateTransferArr = []
      if (this.tabName === 'pending' && this.type === '2') {
        dateTransferArr = ['taskExpectTime', 'taskReportTime']
      } else if (this.tabName === 'pending' && this.type === '1') {
        dateTransferArr = ['expectTime', 'reportTime']
      } else if (this.tabName === 'hasProcessed' && this.type === '2') {
        dateTransferArr = ['taskExpectTime', 'taskReportTime', 'taskApprovalTime']
      } else if (this.tabName === 'hasProcessed' && this.type === '1') {
        dateTransferArr = ['expectTime', 'reportTime', 'approvalTime']
      } else if (this.tabName === 'submit') {
        dateTransferArr = ['expectTime', 'reportTime']
      } else if (this.tabName === 'draft') {
        dateTransferArr = ['createdTime', 'updatedTime', 'expectTime']
      }
      dateTransferArr.forEach(item => {
        if (form[item] && form[item].length > 0 && form[item][0]) {
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
        type: Number(this.type),
        rollback: Number(this.rollback),
        ...form,
        startIndex: (this.pagination.currentPage - 1) * this.pagination.pageSize,
        pageSize: this.pagination.pageSize
      }
      // 获取默认排序
      if (!dynamicSort) {
        this.initSortTable()
      }
      if (this.sorting) {
        params.sorting = this.sorting
      }
      // 过滤掉多余属性
      if (!['pending', 'hasProcessed'].includes(this.tabName)) {
        delete params.type
      }
      if (!['submit'].includes(this.tabName)) {
        delete params.rollback
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
    // 表格操作-查看
    hanldeView (row) {
      const path = this.actionName === '1' ? 'detailPublish' : 'detailRequest'
      const url = `/taskman/workbench/${path}`
      this.$router.push({
        path: url,
        query: {
          requestId: row.id,
          requestTemplate: row.templateId,
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
        await changeTaskStatus('start', row.taskId, new Date(row.taskUpdatedTime).getTime())
      }
      const path = this.actionName === '1' ? 'detailPublish' : 'detailRequest'
      const url = `/taskman/workbench/${path}`
      this.$router.push({
        path: url,
        query: {
          requestId: row.id,
          requestTemplate: row.templateId,
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
        title: this.$t('confirm'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          // 请求定版的新接口，任务处理的老接口
          let res = null
          if (row.status === 'Pending') {
            res = await tansferToMe(row.id, new Date(row.updatedTime).getTime())
          } else if (row.status === 'InProgress') {
            res = await changeTaskStatus(type, row.taskId, new Date(row.taskUpdatedTime).getTime())
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
            jumpFrom: ''
          }
        })
      }
    },
    // 表格操作撤回
    async handleRecall (row) {
      this.$Modal.confirm({
        title: this.$t('confirm'),
        content: this.$t('tw_recall_tips'),
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
          jumpFrom: this.tabName === 'draft' ? 'my_drafts' : 'my_submit',
          rollback: this.rollback
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
    // 点击名称，id，任务名快捷跳转
    handleDbClick (row) {
      if (
        this.username === row.handler &&
        ['Pending', 'InProgress'].includes(row.status) &&
        this.tabName === 'pending'
      ) {
        this.handleEdit(row)
      } else if (row.status === 'Draft' && this.tabName !== 'hasProcessed') {
        this.hanldeLaunch(row)
      } else if (['Termination', 'Completed', 'Faulted'].includes(row.status) && this.tabName === 'submit') {
        this.handleRepub(row)
      } else if (row.status !== 'Draft') {
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
    margin-top: 10px;
  }
}
</style>
<style lang="scss">
.workbench {
  .ivu-progress-outer {
    display: flex;
    align-items: center;
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
  .ivu-tooltip-inner {
    max-width: 1000px;
  }
}
</style>
