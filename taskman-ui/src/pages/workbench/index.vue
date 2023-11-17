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
    <div class="hot-link">
      <HotLink></HotLink>
    </div>
    <div class="data-card">
      <DataCard @fetchData="handleOverviewChange"></DataCard>
    </div>
    <div class="data-tabs">
      <Tabs v-model="actionName" @on-click="handleQuery">
        <TabPane label="发布" name="1"></TabPane>
        <TabPane label="请求" name="2"></TabPane>
      </Tabs>
      <!--搜索条件-->
      <BaseSearch :options="searchOptions" v-model="form" @search="handleQuery"></BaseSearch>
      <!--表格分页-->
      <Table
        border
        size="small"
        :columns="tableColumns"
        :data="tableData"
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
  </div>
</template>

<script>
import WorkBench from '@/pages/components/workbench-menu.vue'
import HotLink from './components/hot-link.vue'
import DataCard from './components/data-card.vue'
import BaseSearch from '../components/base-search.vue'
import { getPlatformList } from '@/api/server'
export default {
  components: {
    HotLink,
    DataCard,
    BaseSearch,
    WorkBench
  },
  data() {
    return {
      tabName: 'pending', //pending待处理,hasProcessed已处理,submit我提交的,draft我的暂存,collect收藏
      actionName: '1', //1发布行为,2请求,3问题,4事件,5变更
      form: {
        type: 0, //0所有,1请求定版,2任务处理
        query: '', //ID或名称模糊搜索
        templateId: '', //模板ID
        status: [], //状态
        operatorObj: '', //操作对象
        createdBy: '' //创建人
      },
      searchOptions: [
        {
          key: 'type',
          initValue: 0,
          component: 'radio-group',
          list: [
            { label: '所有', value: 0 },
            { label: '请求定版', value: 1 },
            { label: '任务处理', value: 2 }
          ]
        },
        {
          key: 'query',
          placeholder: '名称',
          component: 'input'
        },
        {
          key: 'templateId',
          placeholder: '模板',
          component: 'select',
          list: [
            { label: '已完成', value: 1 },
            { label: '未完成', value: 2 },
            { label: '进行中', value: 3 }
          ]
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
        },
        {
          key: 'operatorObj',
          placeholder: '操作对象',
          component: 'input'
        },
        {
          key: 'createdBy',
          placeholder: '创建人',
          component: 'input'
        }
      ],
      tableColumns: [
        {
          title: '请求ID',
          width: 150,
          key: 'id'
        },
        {
          title: this.$t('name'),
          width: 200,
          key: 'name'
        },
        {
          title: '使用模板',
          sortable: 'custom',
          minWidth: 120,
          key: 'templateName'
        },
        {
          title: '操作对象',
          resizable: true,
          width: 200,
          key: 'operatorObj'
        },
        {
          title: '使用编排',
          minWidth: 200,
          key: 'procDefName'
        },
        {
          title: '请求状态',
          sortable: 'custom',
          key: 'status',
          minWidth: 120
        },
        {
          title: '当前节点',
          sortable: 'custom',
          minWidth: 120,
          key: 'curNode'
        },
        {
          title: '进展',
          sortable: 'custom',
          minWidth: 120,
          key: 'progress'
        },
        {
          title: '停留时长',
          sortable: 'custom',
          minWidth: 160,
          key: 'expectTime',
          render: (h, params) => {
            // 需根据createdTime、expectTime计算出来
          }
        },
        {
          title: '创建人',
          sortable: 'custom',
          minWidth: 160,
          key: 'created_by'
        },
        {
          title: '当前处理人',
          sortable: 'custom',
          minWidth: 160,
          key: 'handler'
        },
        {
          title: '创建时间',
          sortable: 'custom',
          minWidth: 160,
          key: 'createdTime'
        },
        {
          title: '期望完成时间',
          sortable: 'custom',
          minWidth: 160,
          key: 'expectTime'
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
      ],
      tableData: [],
      pagination: {
        total: 0,
        currentPage: 1,
        pageSize: 10
      }
    }
  },
  mounted() {
    this.getList()
  },
  methods: {
    handleOverviewChange(val) {
      this.tabName = val
    },
    async getList() {
      const params = {
        tab: this.tabName,
        action: Number(this.actionName),
        ...this.form,
        startIndex: this.pagination.currentPage,
        pageSize: this.pagination.pageSize
      }
      const { statusCode, data } = await getPlatformList(params)
      if (statusCode === 'OK') {
        this.tableData = data.contents || []
      }
    },
    handleQuery() {
      this.pagination.currentPage = 1
      this.getList()
    },
    changPage(val) {
      this.pagination.currentPage = val
      this.getList()
    },
    changePageSize(val) {
      this.pagination.currentPage = 1
      this.pagination.pageSize = val
      this.getList()
    },
    //查看
    hanldeView() {

    },
    //重新发布
    handleRepub() {

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
    .left, .right {
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
    margin-top: 24px;
  }
  .data-card {
    margin-top: 24px;
  }
  .data-tabs {
    margin-top: 24px;
  }
}
</style>