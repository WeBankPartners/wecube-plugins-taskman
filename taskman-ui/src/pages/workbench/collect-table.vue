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
import ScrollTag from '@/pages/components/scroll-tag.vue'
import { collectTemplateList, uncollectTemplate, getTemplateFilter } from '@/api/server'
import { deepClone } from '@/pages/util/index'
export default {
  components: {
    BaseSearch,
    ScrollTag
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
          placeholder: this.$t('tw_template_id'),
          component: 'input'
        },
        {
          key: 'name',
          placeholder: this.$t('tw_template_name'),
          component: 'input'
        },
        {
          key: 'templateGroupId',
          placeholder: this.$t('tm_template_group'),
          multiple: true,
          component: 'select',
          list: []
        },
        {
          key: 'operatorObjType',
          placeholder: this.$t('tw_operator_type'),
          multiple: true,
          component: 'select',
          list: []
        },
        {
          key: 'procDefName',
          placeholder: this.$t('tw_template_flow'),
          multiple: true,
          component: 'select',
          list: []
        },
        {
          key: 'manageRole',
          placeholder: this.$t('tw_template_owner_role'),
          multiple: true,
          component: 'select',
          list: []
        },
        {
          key: 'owner',
          placeholder: this.$t('tw_template_owner'),
          multiple: true,
          component: 'select',
          list: []
        },
        {
          key: 'useRole',
          placeholder: this.$t('useRoles'),
          multiple: true,
          component: 'select',
          list: []
        },
        {
          key: 'tags',
          placeholder: this.$t('tags'),
          multiple: true,
          component: 'select',
          list: []
        },
        {
          key: 'createdTime',
          label: this.$t('tw_created_time'),
          dateType: 4,
          labelWidth: 85,
          component: 'custom-time'
        },
        {
          key: 'updatedTime',
          label: this.$t('tw_update_time'),
          dateType: 4,
          labelWidth: 85,
          component: 'custom-time'
        }
      ],
      tableColumns: [
        {
          title: this.$t('tw_template_id'),
          minWidth: 150,
          key: 'parentId'
        },
        {
          title: this.$t('tw_template_name'),
          width: 200,
          key: 'name',
          render: (h, params) => {
            return (
              <span>
                {`${params.row.name}【${params.row.version}】`}
                {/* <Tag>{params.row.version}</Tag> */}
              </span>
            )
          }
        },
        {
          title: this.$t('tm_template_group'),
          sortable: 'custom',
          minWidth: 120,
          key: 'templateGroup'
        },
        {
          title: this.$t('tw_template_flow'),
          minWidth: 180,
          key: 'procDefName'
        },
        {
          title: this.$t('tw_operator_type'),
          resizable: true,
          sortable: 'custom',
          minWidth: 140,
          key: 'operatorObjType',
          render: (h, params) => {
            return (
              params.row.operatorObjType && (
                <Tooltip content={params.row.operatorObjType} placement="top">
                  <Tag>{params.row.operatorObjType}</Tag>
                </Tooltip>
              )
            )
          }
        },
        {
          title: this.$t('tw_template_status'),
          minWidth: 120,
          key: 'status',
          render: (h, params) => {
            const list = [
              { label: this.$t('tw_template_status_use'), value: 1, color: '#19be6b' },
              { label: this.$t('tw_template_status_disable'), value: 2, color: '#c5c8ce' },
              { label: this.$t('tw_template_status_role'), value: 3, color: '#ed4014' }
            ]
            const item = list.find(i => i.value === params.row.status)
            return item && <Tag color={item.color}>{item.label}</Tag>
          }
        },
        {
          title: '审批列表',
          minWidth: 160,
          key: 'approves',
          render: (h, params) => {
            return <ScrollTag list={params.row.approves} />
          }
        },
        {
          title: '任务节点',
          minWidth: 160,
          key: 'tasks',
          render: (h, params) => {
            return <ScrollTag list={params.row.tasks} />
          }
        },
        {
          title: '属主角色/人',
          sortable: 'custom',
          key: 'manageRole',
          minWidth: 130,
          render: (h, params) => {
            return (
              <div style="display:flex;flex-direction:column">
                <span>{params.row.manageRole}</span>
                <span>{params.row.owner}</span>
              </div>
            )
          }
        },
        {
          title: this.$t('useRoles'),
          sortable: 'custom',
          minWidth: 130,
          key: 'useRole'
        },
        {
          title: this.$t('tags'),
          minWidth: 130,
          key: 'tags',
          render: (h, params) => {
            return (
              params.row.tags && (
                <Tooltip content={params.row.tags} placement="top">
                  <Tag>{params.row.tags}</Tag>
                </Tooltip>
              )
            )
          }
        },
        {
          title: this.$t('tw_human_task'),
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
          title: this.$t('tw_created_time'),
          sortable: 'custom',
          minWidth: 150,
          key: 'createdTime'
        },
        {
          title: this.$t('tw_update_time'),
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
                <Tooltip content={this.$t('tw_launch')} placement="top">
                  <Button
                    type="success"
                    size="small"
                    onClick={() => {
                      this.hanldeCreate(params.row)
                    }}
                    class={disableClass}
                    style="margin-right: 5px"
                  >
                    <Icon type="ios-send" size="16"></Icon>
                  </Button>
                </Tooltip>
                <Tooltip content={this.$t('tw_uncollec_tooltip')} placement="top">
                  <Button
                    type="warning"
                    size="small"
                    onClick={() => {
                      this.handleUnStar(params.row)
                    }}
                  >
                    <Icon type="ios-star-half" size="16"></Icon>
                  </Button>
                </Tooltip>
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
      createRouteMap: {
        '1': 'createPublish',
        '2': 'createRequest',
        '3': 'createProblem',
        '4': 'createEvent',
        '5': 'createChange'
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
      // 模板禁用提示
      if (row.status === 2) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('tw_template_disable_tips')
        })
        // 模板权限移除提示
      } else if (row.status === 3) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('tw_template_role_tips')
        })
      }
      const path = this.createRouteMap[this.actionName]
      const url = `/taskman/workbench/${path}`
      this.$router.push({
        path: url,
        query: {
          requestTemplate: row.id,
          role: row.manageRole
        }
      })
    },
    // 取消收藏
    handleUnStar (row) {
      this.$Modal.confirm({
        title: this.$t('tw_confirm') + this.$t('tw_uncollec_tooltip'),
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
            this.$parent.$refs.dataCard.getData()
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
