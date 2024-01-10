<!--请求审计-->
<template>
  <div class="workbench-request-audit">
    <Tabs :value="actionName" @on-click="handleChangeTab">
      <!--所有-->
      <TabPane :label="$t('tw_all_tab')" name="0"></TabPane>
      <!--发布-->
      <TabPane :label="$t('tw_publish')" name="1"></TabPane>
      <!--请求-->
      <TabPane :label="$t('tw_request')" name="2"></TabPane>
    </Tabs>
    <BaseSearch :options="searchOptions" v-model="form" @search="handleQuery"></BaseSearch>
    <Button size="small" @click="handleExport" type="success" :loading="exportFlag" class="export">{{
      $t('download')
    }}</Button>
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
export default {
  components: {
    BaseSearch
  },
  mixins: [column, search],
  data () {
    return {
      actionName: '0', // 0所有、1发布、2请求
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
  watch: {
    // 切换请求发布，重新获取筛选条件配置
    actionName () {
      this.getFilterOptions()
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
    this.searchOptions.forEach(item => {
      // 请求状态设置默认值
      if (item.key === 'status') {
        item.initValue = ['Completed', 'Termination', 'Faulted']
      }
    })
    this.handleReset()
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
        action: Number(this.actionName),
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
      this.actionName = val
      this.handleReset()
      this.handleQuery()
    },
    // 点击名称，id，任务名快捷跳转
    handleDbClick (row) {
      this.hanldeView(row)
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
  .export {
    width: 50px;
    height: 30px;
    position: absolute;
    right: 30px;
    top: 95px;
    font-size: 14px;
  }
}
</style>
<style lang="scss">
.workbench-request-audit {
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
}
</style>
