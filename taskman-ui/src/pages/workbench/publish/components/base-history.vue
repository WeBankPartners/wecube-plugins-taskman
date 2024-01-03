<template>
  <div class="workbench-publish-history">
    <Tabs :value="activeTab" @on-click="handleChangeTab">
      <!--已提交-->
      <TabPane :label="$t('tw_commit_tab')" name="commit"></TabPane>
      <!--草稿箱-->
      <TabPane :label="$t('tw_draft_tab')" name="draft"></TabPane>
      <!--被退回-->
      <TabPane :label="$t('tw_return_tab')" name="rollback"></TabPane>
      <!--本人撤回-->
      <TabPane :label="$t('tw_recall_tab')" name="revoke"></TabPane>
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
import { getPublishList, reRequest, recallRequest } from '@/api/server'
import column from '../../column'
import search from '../../search'
import { deepClone } from '@/pages/util/index'
export default {
  components: {
    BaseSearch
  },
  mixins: [column, search],
  props: {
    // 1发布,2请求(3问题,4事件,5变更)
    actionName: {
      type: String,
      default: '1'
    }
  },
  data () {
    return {
      activeTab: 'commit', // commit已提交,draft我暂存草稿,rollback我退回,revoke本人撤回
      tabName: 'submit',
      rollback: '0', // 0所有,1已退回,2其他,3被撤回
      form: {
        name: '', // 请求名
        id: '',
        templateId: [], // 模板ID
        status: [], // 状态
        operatorObjType: [], // 操作对象类型
        procDefName: [], // 使用编排
        createdBy: [], // 创建人
        expectTime: [], // 期望时间
        reportTime: [] // 请求提交时间
      },
      tableData: [],
      loading: false,
      sorting: {},
      pagination: {
        total: 0,
        currentPage: 1,
        pageSize: 10
      },
      searchOptions: [],
      tableColumns: []
    }
  },
  mounted () {
    this.tableColumns = deepClone(this.submitAllColumn)
    this.searchOptions = this.submitSearch
    this.handleReset()
    this.getList()
  },
  methods: {
    // 切换类型
    handleChangeTab (val) {
      if (val === 'commit') {
        this.rollback = '0'
        this.tabName = 'submit'
        this.tableColumns = deepClone(this.submitAllColumn)
        this.searchOptions = this.submitSearch
      } else if (val === 'rollback') {
        this.rollback = '1'
        this.tabName = 'submit'
        this.tableColumns = deepClone(this.submitAllColumn)
        this.searchOptions = this.submitSearch
      } else if (val === 'draft') {
        this.tabName = 'draft'
        this.tableColumns = deepClone(this.draftColumn)
        this.searchOptions = this.draftSearch
      } else if (val === 'revoke') {
        this.rollback = '3'
        this.tabName = 'submit'
        this.tableColumns = deepClone(this.submitColumn)
        this.searchOptions = this.submitSearch
      }
      this.activeTab = val
      this.handleReset()
      this.handleQuery()
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
    async getList (dynamicSort) {
      this.loading = true
      const form = deepClone(this.form)
      let dateTransferArr
      if (this.activeTab === 'draft') {
        dateTransferArr = ['createdTime', 'updatedTime', 'expectTime']
      } else {
        dateTransferArr = ['expectTime', 'reportTime']
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
        tab: this.activeTab,
        action: Number(this.actionName),
        ...form,
        startIndex: (this.pagination.currentPage - 1) * this.pagination.pageSize,
        pageSize: this.pagination.pageSize
      }
      // 获取默认排序
      if (!dynamicSort) {
        if (this.activeTab === 'draft') {
          this.sorting = {
            asc: false,
            field: 'updatedTime'
          }
        } else {
          this.sorting = {
            asc: false,
            field: 'reportTime'
          }
        }
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
    // 点击名称，id快捷跳转
    handleDbClick (row) {
      if (row.status === 'Draft') {
        this.hanldeLaunch(row)
      } else if (['Termination', 'Completed', 'Faulted'].includes(row.status) && this.tabName === 'submit') {
        this.handleRepub(row)
      } else {
        this.hanldeView(row)
      }
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
          jumpFrom: ''
        }
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
    // 表格操作-草稿去发起
    hanldeLaunch (row) {
      const path = this.actionName === '1' ? 'createPublish' : 'createRequest'
      const url = `/taskman/workbench/${path}`
      this.$router.push({
        path: url,
        query: {
          requestId: row.id,
          requestTemplate: row.templateId,
          jumpFrom: 'my_drafts'
        }
      })
    },
    // 表格操作撤回
    async handleRecall (row) {
      this.$Modal.confirm({
        title: this.$t('confirm') + this.$t('tw_recall'),
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
    }
  }
}
</script>

<style lang="scss" scoped>
.workbench-publish-history {
  width: 100%;
}
</style>
<style lang="scss">
.workbench-publish-history {
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
