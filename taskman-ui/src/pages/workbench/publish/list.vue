<template>
  <div class="workbench-publish-history">
    <Tabs :value="activeTab" @on-click="handleChangeTab">
      <TabPane label="已提交" name="commit"></TabPane>
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
import { getPublishList, reRequest, recallRequest } from '@/api/server'
import dayjs from 'dayjs'
import column from '../column'
import search from '../search'
import { deepClone } from '@/pages/util/index'
export default {
  components: {
    BaseSearch
  },
  mixins: [column, search],
  data () {
    return {
      activeTab: 'commit',
      tabName: 'submit',
      form: {
        name: '', // 请求名
        id: '',
        templateId: [], // 模板ID
        status: [], // 状态
        operatorObjType: [], // 操作对象类型
        procDefName: [], // 使用编排
        createdBy: [], // 创建人
        expectTime: [
          dayjs()
            .subtract(3, 'month')
            .format('YYYY-MM-DD'),
          dayjs().format('YYYY-MM-DD')
        ], // 期望时间
        reportTime: [
          dayjs()
            .subtract(3, 'month')
            .format('YYYY-MM-DD'),
          dayjs().format('YYYY-MM-DD')
        ] // 请求提交时间
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
      searchOptions: [],
      tableColumns: []
    }
  },
  mounted () {
    this.tableColumns = deepClone(this.submitAllColumn)
    this.searchOptions = this.submitSearch
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
      const dateTransferArr = ['expectTime', 'reportTime']
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
        action: 1,
        ...form,
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
    handleDbClick (row) {
      this.hanldeView(row)
    },
    // 表格操作-查看
    hanldeView (row) {
      const url = `/taskman/workbench/createPublish`
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
        const url = `/taskman/workbench/createPublish`
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
      const url = `/taskman/workbench/createPublish`
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
