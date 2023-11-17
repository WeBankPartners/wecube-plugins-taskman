<template>
  <div class="workbench-request-history">
    <Tabs :value="activeTab" @on-click="handleChangeTab">
      <TabPane label="已发布" name="publish"></TabPane>
      <TabPane label="草稿箱" name="draft"></TabPane>
    </Tabs>
    <BaseSearch
      :options="searchOptions"
      v-model="searchForm"
      @search="handleQuery"
    ></BaseSearch>
    <Table
      size="small"
      :columns="tableColumns"
      :data="tableData"
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
export default {
  components: {
    BaseSearch
  },
  data() {
    return {
      activeTab: 'publish',
      searchForm: {
        name: '',
        object: '',
        package: '',
        id: '',
        tag: [],
        status: [],
        createTime: []
      },
      pagination: {
        pageSize: 10,
        currentPage: 1,
        total: 0
      },
      searchOptions: [
        {
          key: 'name',
          placeholder: '请求名称',
          component: 'input'
        },
        {
          key: 'object',
          placeholder: '操作对象',
          component: 'select',
          list: [{label: '已完成', value: 1},{label: '未完成', value: 2},{label: '进行中', value: 3}]
        },
        {
          key: 'package',
          placeholder: '部署包',
          component: 'input'
        },
        {
          key: 'id',
          placeholder: '请求ID',
          component: 'input',
        },
        {
          key: 'tag',
          placeholder: '当前节点',
          component: 'select',
          multiple: true,
          list: [{label: '已完成', value: 1},{label: '未完成', value: 2},{label: '进行中', value: 3}]
        },
        {
          key: 'status',
          placeholder: '请求状态',
          component: 'select',
          multiple: true,
          list: [{label: '已完成', value: 1},{label: '未完成', value: 2},{label: '进行中', value: 3}]
        },
        {
          key: 'createTime',
          label: '创建时间',
          labelWidth: 70,
          component: 'custom-time'
        }
      ],
      tableData: [
        { name: '1111111' },
        { name: '1111111' },
        { name: '1111111' },
        { name: '1111111' },
        { name: '1111111' },
        { name: '1111111' }
      ],
      tableColumns: [
        {
          title: '请求名称',
          key: 'name',
        },
        {
          title: '请求ID',
          key: 'id'
        },
        {
          title: '请求状态',
          key: 'status'
        },
        {
          title: '进度',
          key: 'progress'
        },
        {
          title: '当前节点',
          key: 'tag'
        },
        {
          title: '发布操作对象',
          key: 'object'
        },
        {
          title: '部署实例',
          key: 'example'
        },
        {
          title: '部署包',
          key: 'package'
        },
        {
          title: '创建人',
          key: 'create'
        },
        {
          title: '创建时间',
          key: 'createTime'
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
                  onClick={() => { this.hanldeView(params.row) }}
                  style="margin-right:5px;"
                >
                  查看
                </Button>
                <Button
                  type="primary"
                  size="small"
                  onClick={() => { this.handleRepub(params.row) }}
                >
                  重新发布
                </Button>
              </div>
            )
          }
        }
      ]
    }
  },
  methods: {
    getList() {

    },
    handleQuery() {
      this.pagination.currentPage = 1
      this.getList()
    },
    handlePage(val) {
      this.pagination.currentPage = val
      this.getList()
    },
    handlePageSize(val) {
      this.pagination.currentPage = 1
      this.pagination.pageSize = val
      this.getList()
    },
    handleChangeTab(val) {
      this.activeTab = val
    },
    //查看
    hanldeView(row) {

    },
    //重新发布
    handleRepub(row) {

    }
  }
}
</script>

<style lang="scss" scoped>
.workbench-request-history {
  width: 100%;
}
</style>