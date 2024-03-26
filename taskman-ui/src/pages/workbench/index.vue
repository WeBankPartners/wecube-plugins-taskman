<!--工作台-->
<template>
  <div class="workbench">
    <div class="create-time">
      <DatePicker
        :value="queryTime"
        @on-change="
          val => {
            handleQueryDateRange(val)
          }
        "
        type="daterange"
        placement="bottom-end"
        format="yyyy-MM-dd"
        :placeholder="$t('tw_created_time')"
        style="width: 200px"
        :clearable="false"
      />
    </div>
    <DataCard
      ref="dataCard"
      :initTab="initTab"
      :initAction="initAction"
      @initFetch="initData"
      @fetchData="handleOverviewChange"
    ></DataCard>
    <div class="data-tabs">
      <!--tab标签添加name=workbench属性，解决自定义home页面bug-->
      <Tabs
        name="workbench"
        v-if="['myPending', 'pending', 'hasProcessed'].includes(tabName)"
        v-model="type"
        @on-click="handleTypeChange"
      >
        <!--审批-->
        <TabPane :label="approveLabel" name="3" tab="workbench"></TabPane>
        <!--任务处理-->
        <TabPane :label="taskLabel" name="2" tab="workbench"></TabPane>
        <!--请求定版-->
        <TabPane :label="pendingLabel" name="1" tab="workbench"></TabPane>
        <!--请求确认-->
        <TabPane :label="confirmLabel" name="4" tab="workbench"></TabPane>
      </Tabs>
      <Tabs name="workbench" v-if="['submit'].includes(tabName)" v-model="rollback" @on-click="handleRollbackChange">
        <!--所有-->
        <TabPane :label="$t('tw_all_tab')" name="0" tab="workbench"></TabPane>
        <!--被退回-->
        <TabPane :label="$t('tw_return_tab')" name="1" tab="workbench"></TabPane>
        <!--本人撤回-->
        <TabPane :label="$t('tw_recall_tab')" name="3" tab="workbench"></TabPane>
        <!--其他-->
        <TabPane :label="$t('tw_other_tab')" name="2" tab="workbench"></TabPane>
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
import HotLink from './components/hot-link.vue'
import DataCard from './components/data-card.vue'
import BaseSearch from '../components/base-search.vue'
import CollectTable from './collect-table.vue'
import { getPlatformList, recallRequest, pendingHandle, deleteRequest, reRequest } from '@/api/server'
import { deepClone } from '@/pages/util/index'
import dayjs from 'dayjs'
import column from './column.js'
import search from './search.js'
export default {
  components: {
    HotLink,
    DataCard,
    BaseSearch,
    CollectTable
  },
  mixins: [column, search],
  data () {
    return {
      queryTime: [
        dayjs()
          .subtract(3, 'month')
          .format('YYYY-MM-DD'),
        dayjs().format('YYYY-MM-DD')
      ], // 全局时间过滤器
      tabName: 'myPending', // pending(myPending本人处理/pending本组处理),hasProcessed已处理,submit我提交的,draft我的暂存,collect收藏
      actionName: '1', // 1发布,2请求,3问题,4事件,5变更
      initTab: '',
      initAction: '',
      type: '0', // 0所有,1请求定版,2任务处理,3审批,4请求确认
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
        taskCreatedTime: [], // 任务提交时间
        taskHandleUpdatedTime: [], // 任务更新时间
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
      sorting: {}, // 表格默认排序
      taskLabel: () => {
        return (
          <div>
            <span>{this.$t('tw_task_tab')}</span>
            {['myPending', 'pending'].includes(this.tabName) && this.getPendingNumber('Task') > 0 && (
              <span class="badge">{this.getPendingNumber('Task')}</span>
            )}
          </div>
        )
      },
      approveLabel: () => {
        return (
          <div>
            <span>{this.$t('tw_approval')}</span>
            {['myPending', 'pending'].includes(this.tabName) && this.getPendingNumber('Approve') > 0 && (
              <span class="badge">{this.getPendingNumber('Approve')}</span>
            )}
          </div>
        )
      },
      pendingLabel: () => {
        return (
          <div>
            <span>{this.$t('tw_pending_tab')}</span>
            {['myPending', 'pending'].includes(this.tabName) && this.getPendingNumber('Check') > 0 && (
              <span class="badge">{this.getPendingNumber('Check')}</span>
            )}
          </div>
        )
      },
      confirmLabel: () => {
        return (
          <div>
            <span>{this.$t('tw_confirm_tab')}</span>
            {['myPending', 'pending'].includes(this.tabName) && this.getPendingNumber('Confirm') > 0 && (
              <span class="badge">{this.getPendingNumber('Confirm')}</span>
            )}
          </div>
        )
      }
    }
  },
  computed: {
    getPendingNumber () {
      return function (type) {
        return Number(this.$refs.dataCard.pendingNumObj[this.actionName][type]) || 0
      }
    }
  },
  watch: {
    // 切换请求发布，重新获取筛选条件配置
    actionName () {
      this.getFilterOptions()
    }
  },
  mounted () {
    this.initTab = this.$route.query.tabName || 'myPending'
    this.initAction = this.$route.query.actionName || '1'
  },
  methods: {
    handleQueryDateRange (val) {
      this.queryTime = val
      this.handleQuery(true)
    },
    // 初始化加载数据(链接携带参数，跳转到指定标签)
    initData (val, action) {
      this.tabName = val
      this.actionName = action
      const type = this.$route.query.type
      const rollback = this.$route.query.rollback
      if (['myPending', 'pending', 'hasProcessed'].includes(val)) {
        if (['1', '2', '3', '4'].includes(type)) {
          this.type = type
        } else {
          this.type = '3'
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
        this.handleQuery(true)
      }
    },
    // 点击视图卡片触发查询
    handleOverviewChange (val, action) {
      this.tabName = val
      this.actionName = action
      if (['myPending', 'pending', 'hasProcessed'].includes(val)) {
        this.type = '3'
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
      }
    },
    // 切换type
    handleTypeChange () {
      this.getTypeConfig()
      this.handleReset()
      this.handleQuery()
    },
    // 待处理、已处理表格差异化配置
    getTypeConfig () {
      if (this.tabName === 'pending' || this.tabName === 'myPending') {
        this.tableColumn = this.pendingTaskColumn
        if (['1', '4'].includes(this.type)) {
          this.searchOptions = this.pendingSearch
        } else if (['2', '3'].includes(this.type)) {
          this.searchOptions = this.pendingTaskSearch
        }
      } else if (this.tabName === 'hasProcessed') {
        this.tableColumn = this.hasProcessedTaskColumn
        if (['1', '4'].includes(this.type)) {
          this.searchOptions = this.hasProcessedSearch
        } else if (['2', '3'].includes(this.type)) {
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
    // 我提交的表格差异化配置
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
      if (this.tabName === 'pending' || this.tabName === 'myPending') {
        this.sorting = {
          asc: false,
          field: 'taskHandleUpdatedTime'
        }
      } else if (this.tabName === 'hasProcessed') {
        this.sorting = {
          asc: false,
          field: 'taskApprovalTime'
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
      if (this.tabName === 'pending' || this.tabName === 'myPending') {
        dateTransferArr = ['taskExpectTime', 'taskHandleUpdatedTime']
      } else if (this.tabName === 'hasProcessed') {
        dateTransferArr = ['taskExpectTime', 'taskCreatedTime', 'taskApprovalTime']
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
        queryTimeStart: this.queryTime[0] && this.queryTime[0] + ' 00:00:00',
        queryTimeEnd: this.queryTime[1] && this.queryTime[1] + ' 23:59:59',
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
      if (!['myPending', 'pending', 'hasProcessed'].includes(this.tabName)) {
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
    handleQuery (all = false) {
      this.pagination.currentPage = 1
      this.getList()
      this.$refs.dataCard.getData(all)
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
      const path = this.detailRouteMap[this.actionName]
      const url = `/taskman/workbench/${path}`
      this.$router.push({
        path: url,
        query: {
          requestId: row.id,
          requestTemplate: row.templateId,
          isCheck: 'Y',
          isHandle: 'N',
          jumpFrom: this.tabName,
          type: this.tabName === 'submit' ? this.rollback : this.type
        }
      })
    },
    // 表格操作-处理(任务、审批、定版、请求确认)
    async handleEdit (row) {
      const path = this.detailRouteMap[this.actionName]
      const url = `/taskman/workbench/${path}`
      this.$router.push({
        path: url,
        query: {
          requestId: row.id,
          requestTemplate: row.templateId,
          taskHandleId: row.taskHandleId, // 任务处理ID
          taskId: row.taskId,
          isCheck: 'N',
          isHandle: 'Y',
          jumpFrom: this.tabName,
          type: this.tabName === 'submit' ? this.rollback : this.type
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
          const params = {
            taskId: row.taskId,
            taskHandleId: row.taskHandleId,
            latestUpdateTime: (row.taskUpdatedTime && String(new Date(row.taskUpdatedTime).getTime())) || '',
            changeReason: type
          }
          const { statusCode } = await pendingHandle(params)
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
    // 表格操作-重新发起
    async handleRepub (row) {
      const { statusCode, data } = await reRequest(row.id)
      if (statusCode === 'OK') {
        const path = this.createRouteMap[this.actionName]
        const url = `/taskman/workbench/${path}`
        this.$router.push({
          path: url,
          query: {
            requestId: data.id,
            requestTemplate: data.requestTemplate,
            jumpFrom: this.tabName,
            type: this.tabName === 'submit' ? this.rollback : this.type
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
            // 刷新本人处理本组处理数量
            this.$refs.dataCard.getData(false, true)
          }
        },
        onCancel: () => {}
      })
    },
    // 表格操作-草稿去发起
    hanldeLaunch (row) {
      const path = this.createRouteMap[this.actionName]
      const url = `/taskman/workbench/${path}`
      this.$router.push({
        path: url,
        query: {
          requestId: row.id,
          requestTemplate: row.templateId,
          jumpFrom: this.tabName,
          type: this.tabName === 'submit' ? this.rollback : this.type
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
        ['Pending', 'InProgress', 'InApproval', 'Confirm'].includes(row.status) &&
        ['myPending', 'pending'].includes(this.tabName)
      ) {
        this.handleEdit(row)
      } else if (this.tabName === 'draft') {
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
  position: relative;
  .create-time {
    width: 200px;
    position: absolute;
    top: -38px;
    right: 0px;
  }
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
  .data-tabs {
    margin-top: 10px;
  }
}
</style>
<style lang="scss">
.workbench {
  .badge {
    position: absolute;
    display: inline-block;
    font-size: 11px;
    background-color: #f56c6c;
    border-radius: 10px;
    color: #fff;
    height: 18px;
    line-height: 18px;
    padding: 0 6px;
    text-align: center;
    white-space: nowrap;
    margin-left: 5px;
  }
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
  }
  .ivu-btn-small {
    font-size: 12px;
  }
  .ivu-tooltip-inner {
    max-width: 1000px;
  }
  .ivu-badge-count {
    padding: 0 4px;
    font-size: 10px;
  }
  .ivu-tabs-tab {
    padding: 6px 20px !important;
  }
}
</style>
