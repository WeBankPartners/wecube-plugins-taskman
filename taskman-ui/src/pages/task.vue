<template>
  <div>
    <div>
      <Row>
        <Col span="4">
          <Input v-model="name" style="width: 90%" type="text" :placeholder="$t('name')"> </Input>
        </Col>
        <Col span="6">
          <Select v-model="status" multiple clearable filterable style="width: 90%" :placeholder="$t('status')">
            <template v-for="option in ['created', 'marked', 'doing', 'done']">
              <Option :label="option" :value="option" :key="option"> </Option>
            </template>
          </Select>
        </Col>
        <Col span="4">
          <Input v-model="handler" style="width:90%" type="text" :placeholder="$t('handler')"> </Input>
        </Col>
        <Col span="4">
          <Button @click="onSearch" type="primary">{{ $t('search') }}</Button>
        </Col>
      </Row>
    </div>
    <Table
      style="margin: 24px 0"
      border
      size="small"
      @on-sort-change="sortTable"
      :columns="tableColumns"
      :data="tableData"
      :max-height="MODALHEIGHT"
    ></Table>
    <Page
      style="float: right"
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
      status: ['created', 'marked', 'doing'],
      handler: '',
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
          resizable: true,
          width: 200,
          fixed: 'left',
          sortable: 'custom',
          key: 'name'
        },
        {
          title: this.$t('task_source'),
          sortable: 'custom',
          minWidth: 100,
          key: 'reporter'
        },
        {
          title: this.$t('request_name'),
          resizable: true,
          width: 200,
          key: '',
          render: (h, params) => {
            return <span>{params.row.requestObj.name}</span>
          }
        },
        {
          title: this.$t('handle_role'),
          minWidth: 100,
          key: 'handleRoles',
          render: (h, params) => {
            const handleRoles = params.row.handleRoles.length === 1 ? params.row.handleRoles[0] : ''
            return <span>{handleRoles}</span>
          }
        },
        {
          title: this.$t('status'),
          sortable: 'custom',
          key: 'status',
          minWidth: 100
        },
        {
          title: this.$t('handler'),
          sortable: 'custom',
          minWidth: 100,
          key: 'handler'
        },
        {
          title: this.$t('created_time'),
          sortable: 'custom',
          minWidth: 160,
          key: 'createdTime'
        },
        {
          title: this.$t('expire_time'),
          sortable: 'custom',
          minWidth: 160,
          key: 'expireTime'
        },
        {
          title: this.$t('hourglass'),
          minWidth: 160,
          key: 'expireTime',
          render: (h, params) => {
            let pColor = '#2bc453'
            if (params.row.status !== 'done') {
              if (params.row.expireObj.percent > 50.0) {
                pColor = '#2d8cf0'
              }
              if (params.row.expireObj.percent > 75.0) {
                pColor = '#f90'
              }
              if (params.row.expireObj.percent > 100.0) {
                pColor = '#ed4014'
              }
            } else {
              pColor = '#9fa8b2'
            }
            return (
              <Tooltip content={params.row.expireObj.useDay + '/' + params.row.expireObj.totalDay} style="width: 100%">
                <Progress percent={params.row.expireObj.percent} stroke-color={pColor} />
              </Tooltip>
            )
          }
        },
        {
          title: this.$t('t_action'),
          key: 'action',
          width: 160,
          fixed: 'right',
          align: 'center',
          render: (h, params) => {
            const operationOptions = params.row.operationOptions
            return (
              <div style="text-align: left">
                <Button onClick={() => this.checkTask(params.row)} style="margin-left: 8px" type="success" size="small">
                  {this.$t('detail')}
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
    sortTable (col) {
      const sorting = {
        asc: col.order === 'asc',
        field: col.key
      }
      this.taskList(sorting)
    },
    success () {
      this.$Notice.success({
        title: this.$t('successful'),
        desc: this.$t('successful')
      })
    },
    async markTask (row) {
      const { statusCode } = await changeTaskStatus('mark', row.id, new Date(row.updatedTime).getTime())
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.taskList()
      }
    },
    async startTask (row) {
      await changeTaskStatus('start', row.id, new Date(row.updatedTime).getTime())
      this.$router.push({ path: '/taskMgmtIndex', query: { taskId: row.id, enforceDisable: 'N' } })
    },
    async checkTask (row) {
      this.$router.push({ path: '/taskMgmtIndex', query: { taskId: row.id, enforceDisable: 'Y' } })
    },
    onSearch () {
      this.pagination.currentPage = 1
      this.taskList()
    },
    changePageSize (pageSize) {
      this.pagination.currentPage = 1
      this.pagination.pageSize = pageSize
      this.taskList()
    },
    changPage (current) {
      this.pagination.currentPage = current
      this.taskList()
    },
    async taskList (sorting) {
      this.payload.filters = []
      if (this.name) {
        this.payload.filters.push({
          name: 'name',
          operator: 'contains',
          value: this.name
        })
      }
      if (this.handler) {
        this.payload.filters.push({
          name: 'handler',
          operator: 'contains',
          value: this.handler
        })
      }
      if (this.status.length > 0) {
        this.payload.filters.push({
          name: 'status',
          operator: 'in',
          value: this.status
        })
      }
      if (sorting) {
        this.payload.sorting = sorting
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
<style>
.ivu-table-cell {
  padding-left: 8px !important;
  padding-right: 8px !important;
}
</style>
<style scoped lang="scss">
.header-icon {
  float: right;
  margin: 3px 40px 0 0 !important;
}
</style>
