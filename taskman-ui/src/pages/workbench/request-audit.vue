<template>
  <div class="workbench-request-audit">
    <Tabs :value="activeTab" @on-click="handleChangeTab">
      <TabPane label="所有" name="0"></TabPane>
      <TabPane label="发布" name="1"></TabPane>
      <TabPane label="请求" name="2"></TabPane>
    </Tabs>
    <BaseSearch :options="searchOptions" v-model="form" @search="handleQuery"></BaseSearch>
    <Button size="small" @click="handleExport" type="success" :loading="exportFlag" style="margin-bottom:10px;"
      >导出</Button
    >
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
import axios from 'axios'
import BaseSearch from '@/pages/components/base-search.vue'
import { getPublishList } from '@/api/server'
import { deepClone } from '@/pages/util/index'
import { getCookie } from '@/pages/util/cookie'
import column from './column'
import search from './search'
import dayjs from 'dayjs'
export default {
  components: {
    BaseSearch
  },
  mixins: [column, search],
  data () {
    return {
      activeTab: '0',
      form: {
        name: '', // 请求名
        id: '',
        templateId: [], // 模板ID
        status: ['Completed', 'Termination', 'Faulted'], // 状态
        operatorObjType: [], // 操作对象类型
        procDefName: [], // 使用编排
        createdBy: [], // 创建人
        expectTime: [], // 期望时间
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
      tableColumns: [],
      headers: {},
      exportFlag: false
    }
  },
  mounted () {
    const accessToken = getCookie('accessToken')
    this.headers = {
      Authorization: 'Bearer ' + accessToken
    }
    this.tableColumns = deepClone(this.submitAllColumn)
    this.tableColumns = this.tableColumns.filter(item => item.key !== 'rollbackDesc' && item.key !== 'action')
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
    getParams () {
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
        tab: 'commit', // 已提交数据，不包括草稿
        action: Number(this.activeTab),
        ...form,
        startIndex: (this.pagination.currentPage - 1) * this.pagination.pageSize,
        pageSize: this.pagination.pageSize
      }
      if (this.sorting) {
        params.sorting = this.sorting
      }
      return params
    },
    async getList () {
      this.loading = true
      const params = this.getParams()
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
    // 点击名称，id，任务名快捷跳转
    handleDbClick (row) {
      this.hanldeView(row)
    },
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
    handleExport () {
      if (this.exportFlag) return
      this.exportFlag = true
      const params = this.getParams()
      axios({
        method: 'post',
        url: `/taskman/api/v1//request/export`,
        data: params,
        headers: this.headers,
        responseType: 'blob'
      })
        .then(response => {
          this.isExport = false
          if (response.status < 400) {
            const arr = response.headers['content-disposition'].split('=') || []
            const fileName = arr[1]
            let blob = new Blob([response.data])
            if ('msSaveOrOpenBlob' in navigator) {
              window.navigator.msSaveOrOpenBlob(blob, fileName)
            } else {
              if ('download' in document.createElement('a')) {
                // 非IE下载
                let elink = document.createElement('a')
                elink.download = fileName
                elink.style.display = 'none'
                elink.href = URL.createObjectURL(blob)
                document.body.appendChild(elink)
                elink.click()
                URL.revokeObjectURL(elink.href) // 释放URL 对象
                document.body.removeChild(elink)
              } else {
                // IE10+下载
                navigator.msSaveOrOpenBlob(blob, fileName)
              }
            }
            this.$Notice.success({
              title: this.$t('successful'),
              desc: this.$t('successful')
            })
          }
        })
        .catch(() => {
          this.$Message.warning('Error')
        })
        .finally(() => {
          this.exportFlag = false
        })
    }
  }
}
</script>

<style lang="scss" scoped>
.workbench-request-audit {
  width: 100%;
}
</style>
<style lang="scss">
.workbench-request-audit {
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
