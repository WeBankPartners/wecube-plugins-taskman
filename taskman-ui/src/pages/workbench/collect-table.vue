<template>
  <div class="workbench-collect-table">
    <!--搜索条件-->
    <BaseSearch :options="searchOptions" v-model="form" @search="handleQuery"></BaseSearch>
    <!--表格分页-->
    <Table
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
  </div>
</template>

<script>
import BaseSearch from '@/pages/components/base-search.vue'
import { collectTemplateList, uncollectTemplate, getTemplateFilter } from '@/api/server'
import { deepClone } from '@/pages/util/index'
export default {
  components: {
    BaseSearch
  },
  props: {
    actionName: {
      type: String,
      default: '1'
    }
  },
  data () {
    return {
      form: {
        name: '',
        id: '',
        templateGroupId: [], // 模板组ID
        operatorObjType: [], // 操作对象类型
        procDefName: [], // 使用编排
        manageRole: [], // 属主角色
        owner: [], // 属主
        useRole: [], // 使用角色
        tags: [], // 标签
        createdTime: [],
        updatedTime: []
      },
      searchOptions: [
        {
          key: 'id',
          placeholder: 'id',
          component: 'input'
        },
        {
          key: 'name',
          placeholder: '名称',
          component: 'input'
        },
        {
          key: 'templateGroupId',
          placeholder: '模板组',
          multiple: true,
          component: 'select',
          list: []
        },
        {
          key: 'operatorObjType',
          placeholder: '操作对象类型',
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
          key: 'manageRole',
          placeholder: '属主角色',
          multiple: true,
          component: 'select',
          list: []
        },
        {
          key: 'owner',
          placeholder: '属主',
          multiple: true,
          component: 'select',
          list: []
        },
        {
          key: 'useRole',
          placeholder: '使用角色',
          multiple: true,
          component: 'select',
          list: []
        },
        {
          key: 'tags',
          placeholder: '标签',
          multiple: true,
          component: 'select',
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
          width: 200,
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
          minWidth: 180,
          key: 'procDefName'
        },
        {
          title: '操作对象类型',
          resizable: true,
          sortable: 'custom',
          minWidth: 140,
          key: 'operatorObjType',
          render: (h, params) => {
            return params.row.operatorObjType && <Tag>{params.row.operatorObjType}</Tag>
          }
        },
        {
          title: '模板状态',
          minWidth: 120,
          key: 'status',
          render: (h, params) => {
            const list = [
              { label: '可使用', value: 1, color: '#19be6b' },
              { label: '已禁用', value: 2, color: '#c5c8ce' },
              { label: '权限被移除', value: 3, color: '#ed4014' }
            ]
            const item = list.find(i => i.value === params.row.status)
            return item && <Tag color={item.color}>{item.label}</Tag>
          }
        },
        {
          title: '属主角色',
          sortable: 'custom',
          key: 'manageRole',
          minWidth: 130
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
          minWidth: 130,
          key: 'useRole'
        },
        {
          title: '标签',
          minWidth: 130,
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
          minWidth: 150,
          key: 'createdTime'
        },
        {
          title: '更新时间',
          sortable: 'custom',
          minWidth: 150,
          key: 'updatedTime'
        },
        {
          title: this.$t('t_action'),
          key: 'action',
          width: 160,
          fixed: 'right',
          align: 'center',
          render: (h, params) => {
            const disableClass = `${[2, 3].includes(params.row.status) ? 'gray-button' : ''}`
            return (
              <div>
                <Button
                  type="info"
                  size="small"
                  onClick={() => {
                    this.hanldeCreate(params.row)
                  }}
                  class={disableClass}
                  style="margin-right: 5px"
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
    this.getFilterOptions()
  },
  methods: {
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
      const form = deepClone(this.form)
      const dateTransferArr = ['createdTime', 'updatedTime']
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
        ...form,
        action: Number(this.actionName),
        startIndex: (this.pagination.currentPage - 1) * this.pagination.pageSize,
        pageSize: this.pagination.pageSize
      }
      if (sort) {
        params.sorting = sort
      }
      const { statusCode, data } = await collectTemplateList(params)
      if (statusCode === 'OK') {
        this.tableData = data.contents || []
        this.pagination.total = data.pageInfo.totalRows
      }
      this.loading = false
    },
    // 获取搜索条件的下拉值
    async getFilterOptions () {
      const { statusCode, data } = await getTemplateFilter({ startTime: '' })
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
          } else if (item.key === 'templateGroupId') {
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
          } else if (item.key === 'manageRole') {
            item.list =
              data.manageRoleList &&
              data.manageRoleList.map(item => {
                return {
                  label: item,
                  value: item
                }
              })
          } else if (item.key === 'owner') {
            item.list =
              data.ownerList &&
              data.ownerList.map(item => {
                return {
                  label: item,
                  value: item
                }
              })
          } else if (item.key === 'owner') {
            item.list =
              data.ownerList &&
              data.ownerList.map(item => {
                return {
                  label: item,
                  value: item
                }
              })
          } else if (item.key === 'useRole') {
            item.list =
              data.useRoleList &&
              data.useRoleList.map(item => {
                return {
                  label: item,
                  value: item
                }
              })
          } else if (item.key === 'tags') {
            item.list =
              data.tagList &&
              data.tagList.map(item => {
                return {
                  label: item,
                  value: item
                }
              })
          }
        })
      }
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
      if (row.status === 2) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: '该模板已禁用'
        })
      } else if (row.status === 3) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: '该模板使用权限已移除'
        })
      }
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
<style lang="scss">
.workbench-collect-table {
  .gray-button {
    background-color: #c5c8ce;
    border-color: #c5c8ce;
  }
}
</style>
