<template>
  <div>
    <div>
      <Row>
        <Col span="4">
          <Input v-model="name" style="width: 90%" type="text" :placeholder="$t('name')"> </Input>
        </Col>
        <Col span="3">
          <Select
            v-model="requestTemplate"
            @on-open-change="getTemplateList"
            clearable
            filterable
            style="width: 90%"
            :placeholder="$t('template')"
          >
            <template v-for="option in templateOptions">
              <Option :label="`${option.name}(${option.version || '-'})`" :value="option.id" :key="option.id"> </Option>
            </template>
          </Select>
        </Col>
        <Col span="3">
          <Select v-model="status" multiple clearable filterable style="width: 90%" :placeholder="$t('status')">
            <template v-for="option in requestStatus">
              <Option :label="option" :value="option" :key="option"> </Option>
            </template>
          </Select>
        </Col>
        <Col span="3">
          <Input v-model="handler" clearable style="width: 90%" type="text" :placeholder="$t('handler')"> </Input>
        </Col>
        <Col span="3">
          <Input v-model="reporter" clearable style="width: 90%" type="text" :placeholder="$t('report_man')"> </Input>
        </Col>
        <Col span="4">
          <DatePicker
            :value="createdTime"
            @on-change="createTimeSelect"
            format="yyyy-MM-dd"
            type="daterange"
            split-panels
            :placeholder="$t('created_time')"
            style="width: 220px"
          ></DatePicker>
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
import { requestListForDraftInitiated, getTemplateList, reRequest } from '@/api/server'
export default {
  name: '',
  data () {
    return {
      MODALHEIGHT: 500,
      name: '',
      handler: '',
      reporter: '', // 提单人
      createdTime: [], // 提单日期
      requestTemplate: '',
      status: [],
      tags: '',
      pagination: {
        pageSize: 10,
        currentPage: 1,
        total: 0
      },
      templateOptions: [],
      payload: {
        filters: [],
        pageable: {
          pageSize: 10,
          startIndex: 0
        },
        paging: true,
        sorting: {
          asc: false,
          field: 'updatedTime'
        }
      },
      tableColumns: [
        {
          title: 'ID',
          minWidth: 130,
          fixed: 'left',
          key: 'id'
        },
        {
          title: this.$t('name'),
          width: 300,
          resizable: true,
          sortable: 'custom',
          key: 'name'
        },
        {
          title: this.$t('priority'),
          key: 'emergency',
          sortable: 'custom',
          width: 80,
          render: (h, params) => {
            const emergencyObj = {
              1: this.$t('high'),
              2: this.$t('medium'),
              3: this.$t('low')
            }
            return <span>{emergencyObj[params.row.emergency]}</span>
          }
        },
        {
          title: this.$t('template'),
          width: 160,
          resizable: true,
          sortable: 'custom',
          key: 'requestTemplateName'
        },
        {
          title: this.$t('createdBy'),
          sortable: 'createdBy',
          minWidth: 140,
          key: 'createdBy'
        },
        {
          title: this.$t('createdTime'),
          sortable: 'createdTime',
          minWidth: 140,
          key: 'createdTime'
        },
        {
          title: this.$t('handle_role'),
          minWidth: 140,
          key: 'handleRoles',
          render: (h, params) => {
            const handleRoles = params.row.handleRoles.length === 1 ? params.row.handleRoles[0] : ''
            return <span>{handleRoles}</span>
          }
        },
        {
          title: this.$t('handler'),
          minWidth: 140,
          sortable: 'custom',
          key: 'handler'
        },
        {
          title: this.$t('report_man'),
          sortable: 'custom',
          minWidth: 140,
          key: 'reporter'
        },
        {
          title: this.$t('status'),
          sortable: 'custom',
          minWidth: 140,
          key: 'status'
        },
        {
          title: this.$t('estimated_finish_time'),
          sortable: 'custom',
          minWidth: 130,
          key: 'expireTime'
        },
        {
          title: this.$t('expected_completion_time'),
          sortable: 'custom',
          minWidth: 130,
          key: 'expectTime'
        },
        {
          title: this.$t('actual_completion_time'),
          sortable: 'custom',
          minWidth: 130,
          key: 'completedTime'
        },
        {
          title: this.$t('report_time'),
          sortable: 'custom',
          minWidth: 130,
          key: 'reportTime'
        },
        {
          title: this.$t('t_action'),
          key: 'action',
          width: 160,
          fixed: 'right',
          align: 'center',
          render: (h, params) => {
            return (
              <div style="text-align: left">
                <Button
                  onClick={() => this.checkTemplate(params.row)}
                  style="margin-left: 8px"
                  type="primary"
                  size="small"
                >
                  {this.$t('detail')}
                </Button>
                {params.row.status.indexOf('Faulted') !== -1 && (
                  <Button
                    onClick={() => this.reRequest(params.row)}
                    style="margin-left: 8px"
                    type="success"
                    size="small"
                  >
                    {this.$t('re-request')}
                  </Button>
                )}
              </div>
            )
            // }
          }
        }
      ],
      tableData: [],
      requestStatus: [
        'Pending',
        'InProgress',
        'InProgress(Faulted)',
        'Termination',
        'Completed',
        'InProgress(Timeouted)',
        'Faulted'
      ],
      timer: null
    }
  },
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 200
    this.requestListForDraftInitiated()
    this.timer = window.setInterval(() => {
      this.requestListForDraftInitiated()
    }, 60000)
  },
  destroyed () {
    window.clearInterval(this.timer)
  },
  methods: {
    createTimeSelect (val) {
      this.createdTime = [...val]
    },
    sortTable (col) {
      const sorting = {
        asc: col.order === 'asc',
        field: col.key
      }
      this.requestListForDraftInitiated(sorting)
    },
    async getTemplateList () {
      const params = {
        filters: [],
        paging: false
      }
      const { statusCode, data } = await getTemplateList(params)
      if (statusCode === 'OK') {
        this.templateOptions = data.contents
      }
    },
    success () {
      this.$Notice.success({
        title: this.$t('successful'),
        desc: this.$t('successful')
      })
    },
    checkTemplate (row) {
      this.$router.push({
        path: row.status === 'Pending' ? '/requestManagementIndex' : '/requestCheck',
        query: {
          requestId: row.id,
          requestTemplate: row.requestTemplate,
          isAdd: 'N',
          isCheck: 'Y',
          isHandle: 'N',
          jumpFrom: 'group_initiated'
        }
      })
    },
    async reRequest (row) {
      await reRequest(row.id)
      const activeTab = 'my_drafts'
      this.$emit('requestTabChange', activeTab)
    },
    onSearch () {
      this.pagination.currentPage = 1
      this.requestListForDraftInitiated()
    },
    changePageSize (pageSize) {
      this.pagination.currentPage = 1
      this.pagination.pageSize = pageSize
      this.requestListForDraftInitiated()
    },
    changPage (current) {
      this.pagination.currentPage = current
      this.requestListForDraftInitiated()
    },
    async requestListForDraftInitiated (sorting) {
      this.payload.filters = [{ name: 'status', operator: 'in', value: this.requestStatus }]
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
      if (this.reporter) {
        this.payload.filters.push({
          name: 'reporter',
          operator: 'contains',
          value: this.reporter
        })
      }
      if (this.createdTime && this.createdTime.length) {
        this.createdTime[0] &&
          this.payload.filters.push({
            name: 'createdTime',
            operator: 'gt',
            value: this.createdTime[0]
          })
        this.createdTime[1] &&
          this.payload.filters.push({
            name: 'createdTime',
            operator: 'lt',
            value: this.createdTime[1]
          })
      }
      if (this.requestTemplate) {
        this.payload.filters.push({
          name: 'requestTemplate',
          operator: 'eq',
          value: this.requestTemplate
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
      const { statusCode, data } = await requestListForDraftInitiated(this.payload)
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
