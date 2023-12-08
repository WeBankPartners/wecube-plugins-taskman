<template>
  <div class="workbench-collect-table">
    <!--搜索条件-->
    <BaseSearch :options="searchOptions" v-model="form" @search="handleQuery"></BaseSearch>
    <!--表格分页-->
    <Table border size="small" :loading="loading" :columns="tableColumns" :data="tableData"></Table>
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
  </div>
</template>

<script>
import BaseSearch from '@/pages/components/base-search.vue'
import { collectTemplateList, uncollectTemplate } from '@/api/server'
export default {
  components: {
    BaseSearch
  },
  props: {
    getTemplateList: {
      type: Function,
      default: () => {}
    },
    actionName: {
      type: String,
      default: '1'
    }
  },
  data () {
    return {
      form: {
        query: '', // ID或名称模糊搜索
        templateId: '', // 模板ID
        status: [] // 状态
      },
      searchOptions: [
        {
          key: 'query',
          placeholder: '名称',
          component: 'input'
        },
        {
          key: 'templateId',
          placeholder: '模板',
          component: 'remote-select',
          remote: this.getTemplateList
        },
        {
          key: 'status',
          placeholder: '状态',
          component: 'select',
          multiple: true,
          list: [
            { label: 'Pending', value: 'Pending' },
            { label: 'InProgress', value: 'InProgress' },
            { label: 'InProgress(Faulted)', value: 'InProgress(Faulted)' },
            { label: 'Termination', value: 'Termination' },
            { label: 'Completed', value: 'Completed' },
            { label: 'InProgress(Timeouted)', value: 'InProgress(Timeouted)' },
            { label: 'Faulted', value: 'Faulted' }
          ]
        }
      ],
      tableColumns: [
        {
          title: '模板ID',
          minWidth: 150,
          key: 'parentId'
        },
        {
          title: '模板名称',
          width: 250,
          key: 'name',
          render: (h, params) => {
            return (
              <span>
                {params.row.name}
                <Tag>{params.row.version}</Tag>
              </span>
            )
          }
        },
        {
          title: '模板组',
          sortable: 'custom',
          minWidth: 120,
          key: 'templateGroup'
        },
        {
          title: '使用编排',
          minWidth: 200,
          key: 'procDefName'
        },
        {
          title: '属主角色',
          sortable: 'custom',
          key: 'manageRole',
          minWidth: 120
        },
        {
          title: '属主',
          sortable: 'custom',
          minWidth: 120,
          key: 'owner'
        },
        {
          title: '使用角色',
          sortable: 'custom',
          minWidth: 120,
          key: 'useRole'
        },
        {
          title: '标签',
          minWidth: 160,
          key: 'tags',
          render: (h, params) => {
            return params.row.tags && <Tag>{params.row.tags}</Tag>
          }
        },
        {
          title: '人工任务',
          minWidth: 160,
          key: 'workNode',
          render: (h, params) => {
            return (
              params.row.workNode &&
              params.row.workNode.map(i => {
                return <Tag color="#2d8cf0">{i}</Tag>
              })
            )
          }
        },
        {
          title: '创建时间',
          sortable: 'custom',
          minWidth: 160,
          key: 'createdTime'
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
                  type="info"
                  size="small"
                  onClick={() => {
                    this.hanldeCreate(params.row)
                  }}
                  style="margin-right:5px;"
                >
                  发起
                </Button>
                <Button
                  type="warning"
                  size="small"
                  onClick={() => {
                    this.handleUnStar(params.row)
                  }}
                >
                  取消收藏
                </Button>
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
      }
    }
  },
  mounted () {
    this.getList()
  },
  methods: {
    async getList () {
      this.loading = true
      const params = {
        ...this.form,
        action: Number(this.actionName),
        startIndex: (this.pagination.currentPage - 1) * this.pagination.pageSize,
        pageSize: this.pagination.pageSize
      }
      const { statusCode, data } = await collectTemplateList(params)
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
    changPage (val) {
      this.pagination.currentPage = val
      this.getList()
    },
    changePageSize (val) {
      this.pagination.currentPage = 1
      this.pagination.pageSize = val
      this.getList()
    },
    // 发起
    hanldeCreate (row) {
      const path = this.actionName === '1' ? 'createPublish' : 'createRequest'
      const url = `/taskman/workbench/${path}`
      this.$router.push({
        path: url,
        query: {
          requestTemplate: row.id,
          role: row.manageRole,
          isAdd: 'Y',
          isCheck: 'N',
          isHandle: 'N',
          jumpFrom: ''
        }
      })
    },
    // 取消收藏
    handleUnStar (row) {
      this.$Modal.confirm({
        title: this.$t('confirm') + '取消收藏',
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const { statusCode } = await uncollectTemplate(row.id)
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

<style lang="scss" scoped></style>
