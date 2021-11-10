<template>
  <div>
    <div>
      <Row>
        <Col span="4">
          <Input v-model="name" style="width:90%" type="text" :placeholder="$t('name')"> </Input>
        </Col>
        <Col span="4">
          <Select
            v-model="requestTemplate"
            @on-open-change="getTemplateList"
            clearable
            filterable
            style="width:90%"
            :placeholder="$t('template')"
          >
            <template v-for="option in templateOptions">
              <Option :label="option.name" :value="option.id" :key="option.id"> </Option>
            </template>
          </Select>
        </Col>
        <Col span="4">
          <Button @click="requestListForDraftInitiated" type="primary">{{ $t('search') }}</Button>
        </Col>
      </Row>
    </div>
    <Table
      style="margin: 24px 0;"
      border
      size="small"
      @on-sort-change="sortTable"
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
import { requestListForDraftInitiated, deleteRequest, terminateRequest, getTemplateList } from '@/api/server'
export default {
  name: '',
  data () {
    return {
      MODALHEIGHT: 500,
      name: '',
      requestTemplate: '',
      tags: '',
      pagination: {
        pageSize: 10,
        currentPage: 1,
        total: 0
      },
      payload: {
        filters: [{ name: 'status', operator: 'in', value: ['draft', 'created'] }],
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
          minWidth: 160,
          fixed: 'left',
          key: 'id'
        },
        {
          title: this.$t('name'),
          resizable: true,
          width: 200,
          sortable: 'custom',
          key: 'name'
        },
        {
          title: this.$t('template'),
          resizable: true,
          width: 200,
          sortable: 'custom',
          key: 'requestTemplateName'
        },
        {
          title: this.$t('priority'),
          minWidth: 80,
          sortable: 'custom',
          key: 'emergency',
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
          title: this.$t('status'),
          minWidth: 80,
          sortable: 'custom',
          key: 'status'
        },
        {
          title: this.$t('handler'),
          sortable: 'custom',
          minWidth: 140,
          key: 'handler'
        },
        {
          title: this.$t('report_time'),
          sortable: 'custom',
          minWidth: 130,
          key: 'report_time'
        },
        {
          title: this.$t('expected_completion_time'),
          sortable: 'custom',
          minWidth: 130,
          key: 'expectTime'
        },
        {
          title: this.$t('estimated_finish_time'),
          sortable: 'custom',
          minWidth: 130,
          key: 'expireTime'
        },
        {
          title: this.$t('t_action'),
          key: 'action',
          width: 200,
          fixed: 'right',
          align: 'center',
          render: (h, params) => {
            return (
              <div style="text-align: left">
                <Button
                  onClick={() => this.editTemplate(params.row)}
                  style="margin-left: 8px"
                  type="primary"
                  size="small"
                >
                  {this.$t('edit')}
                </Button>
                <Button
                  onClick={() => this.deleteRequest(params.row)}
                  style="margin-left: 8px"
                  type="error"
                  size="small"
                >
                  {this.$t('delete')}
                </Button>
              </div>
            )
          }
        }
      ],
      tableData: [],
      templateOptions: []
    }
  },
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 200
    this.requestListForDraftInitiated()
  },
  methods: {
    sortTable (col) {
      const sorting = {
        asc: col.order === 'asc',
        field: col.key
      }
      this.requestListForDraftInitiated(sorting)
    },
    success () {
      this.$Notice.success({
        title: this.$t('successful'),
        desc: this.$t('successful')
      })
    },
    async terminateRequest (row) {
      let res = await terminateRequest(row.id)
      if (res.statusCode === 'OK') {
        this.success()
        this.requestListForDraftInitiated()
      }
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
    deleteRequest (row) {
      this.$Modal.confirm({
        title: this.$t('confirm_delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          let res = await deleteRequest(row.id)
          if (res.statusCode === 'OK') {
            this.success()
            this.requestListForDraftInitiated()
          }
        },
        onCancel: () => {}
      })
    },
    editTemplate (row) {
      this.$router.push({
        path: '/requestManagementIndex',
        query: {
          requestId: row.id,
          requestTemplate: row.requestTemplate,
          isAdd: 'N',
          isCheck: 'N',
          isHandle: 'N',
          jumpFrom: 'my_drafts'
        }
      })
    },
    changePageSize (pageSize) {
      this.pagination.pageSize = pageSize
      this.requestListForDraftInitiated()
    },
    changPage (current) {
      this.pagination.currentPage = current
      this.requestListForDraftInitiated()
    },
    async requestListForDraftInitiated (sorting) {
      this.payload.filters = [{ name: 'status', operator: 'in', value: ['draft', 'created'] }]
      if (this.name) {
        this.payload.filters.push({
          name: 'name',
          operator: 'contains',
          value: this.name
        })
      }
      if (this.requestTemplate) {
        this.payload.filters.push({
          name: 'requestTemplate',
          operator: 'eq',
          value: this.requestTemplate
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
  padding-left: 8px;
  padding-right: 8px;
}
</style>
<style scoped lang="scss">
.header-icon {
  float: right;
  margin: 3px 40px 0 0 !important;
}
</style>
