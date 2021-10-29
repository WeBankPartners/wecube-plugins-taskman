<template>
  <div style="margin: 24px">
    <div>
      <Row>
        <Col span="4">
          <Input v-model="name" style="width:90%" type="text" :placeholder="$t('name')"> </Input>
        </Col>
        <Col span="4">
          <Select v-model="status" clearable filterable style="width:90%" :placeholder="$t('status')">
            <template v-for="option in ['created', 'marked', 'doing', 'done']">
              <Option :label="option" :value="option" :key="option"> </Option>
            </template>
          </Select>
        </Col>
        <Col span="4">
          <Button @click="taskList" type="primary">{{ $t('search') }}</Button>
        </Col>
      </Row>
    </div>
    <Table
      style="margin: 24px 0;"
      border
      size="small"
      :columns="tableColumns"
      :data="tableData"
      :max-height="MODALHEIGHT"
    ></Table>
    <Page
      style="float:right"
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
import { taskList, changeTaskStatus } from '@/api/server'
export default {
  name: '',
  data () {
    return {
      MODALHEIGHT: 500,
      name: '',
      status: '',
      tags: '',
      pagination: {
        pageSize: 10,
        currentPage: 1,
        total: 0
      },
      payload: {
        filters: [],
        pageable: {
          pageSize: 10,
          startIndex: 0
        },
        paging: true,
        sorting: {
          asc: false,
          field: 'createdTime'
        }
      },
      tableColumns: [
        {
          title: this.$t('name'),
          key: 'name'
        },
        {
          title: this.$t('emergency'),
          key: 'emergency',
          width: 80
        },
        {
          title: this.$t('request_name'),
          key: '',
          render: (h, params) => {
            return <span>{params.row.requestObj.name}</span>
          }
        },
        {
          title: this.$t('template'),
          key: '',
          render: (h, params) => {
            return <span>{params.row.requestObj.requestTemplate}</span>
          }
        },
        {
          title: this.$t('node_name'),
          key: 'nodeName'
        },
        {
          title: this.$t('description'),
          key: 'description'
        },
        {
          title: this.$t('status'),
          key: 'status',
          width: 100
        },
        {
          title: this.$t('tm_updated_time'),
          key: 'updatedTime'
        },
        {
          title: this.$t('action'),
          key: 'action',
          width: 160,
          align: 'center',
          render: (h, params) => {
            const operationOptions = params.row.operationOptions
            return (
              <div style="text-align: left">
                <Button onClick={() => this.checkTask(params.row)} style="margin-left: 8px" type="success" size="small">
                  {this.$t('look_over')}
                </Button>
                {operationOptions.includes('mark') && (
                  <Button
                    onClick={() => this.markTask(params.row)}
                    style="margin-left: 8px"
                    type="primary"
                    size="small"
                  >
                    {this.$t('claim')}
                  </Button>
                )}
                {operationOptions.includes('start') && (
                  <Button onClick={() => this.startTask(params.row)} style="margin-left: 8px" type="info" size="small">
                    {this.$t('handle')}
                  </Button>
                )}
              </div>
            )
          }
        }
      ],
      tableData: []
    }
  },
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 200
    this.taskList()
  },
  methods: {
    success () {
      this.$Notice.success({
        title: this.$t('successful'),
        desc: this.$t('successful')
      })
    },
    async markTask (row) {
      const { statusCode } = await changeTaskStatus('mark', row.id)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.taskList()
      }
    },
    async startTask (row) {
      await changeTaskStatus('start', row.id)
      this.$router.push({ path: '/taskMgmtIndex', query: { taskId: row.id, enforceDisable: 'N' } })
    },
    async checkTask (row) {
      this.$router.push({ path: '/taskMgmtIndex', query: { taskId: row.id, enforceDisable: 'Y' } })
    },
    changePageSize (pageSize) {
      this.pagination.pageSize = pageSize
      this.taskList()
    },
    changPage (current) {
      this.pagination.currentPage = current
      this.taskList()
    },
    async taskList () {
      this.payload.filters = []
      if (this.name) {
        this.payload.filters.push({
          name: 'name',
          operator: 'contains',
          value: this.name
        })
      }
      if (this.status) {
        this.payload.filters.push({
          name: 'status',
          operator: 'eq',
          value: this.status
        })
      }
      this.payload.pageable.pageSize = this.pagination.pageSize
      this.payload.pageable.startIndex = (this.pagination.currentPage - 1) * this.pagination.pageSize
      const { statusCode, data } = await taskList(this.payload)
      if (statusCode === 'OK') {
        this.tableData = data.contents
        this.pagination.total = data.pageInfo.totalRows
      }
    }
  },
  components: {}
}
</script>

<style scoped lang="scss">
.header-icon {
  float: right;
  margin: 3px 40px 0 0 !important;
}
</style>
