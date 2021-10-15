<template>
  <div style="margin: 24px">
    <div>
      <Row>
        <Col span="4">
          <Input v-model="name" style="width:90%" type="text" :placeholder="$t('name')"> </Input>
        </Col>
        <Col span="4">
          <Button @click="requestList" type="primary">{{ $t('search') }}</Button>
          <Button @click="addTemplate">{{ $t('initiate_request') }}</Button>
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
import { requestList, deleteRequest, terminateRequest } from '@/api/server'
export default {
  name: '',
  data () {
    return {
      MODALHEIGHT: 500,
      name: '',
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
          asc: true,
          field: 'updatedTime'
        }
      },
      tableColumns: [
        {
          title: this.$t('name'),
          key: 'name'
        },
        {
          title: this.$t('emergency'),
          key: 'emergency'
        },
        {
          title: this.$t('template'),
          key: 'requestTemplate'
        },
        {
          title: this.$t('status'),
          key: 'status'
        },
        {
          title: this.$t('tm_updated_time'),
          key: 'updatedTime'
        },
        {
          title: this.$t('action'),
          key: 'action',
          width: 200,
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
                {params.row.status === 'InProgress' && (
                  <Button
                    onClick={() => this.terminateRequest(params.row)}
                    style="margin-left: 8px"
                    type="info"
                    size="small"
                  >
                    {this.$t('terminate')}
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
    this.requestList()
  },
  methods: {
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
        this.requestList()
      }
    },
    deleteRequest (row) {
      this.$Modal.confirm({
        title: this.$t('confirm_delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          let res = await deleteRequest(row.id)
          this.$Modal.remove()
          if (res.statusCode === 'OK') {
            this.success()
            this.requestList()
          }
        },
        onCancel: () => {}
      })
    },
    editTemplate (row) {
      this.$router.push({ path: '/requestManagementIndex', query: { requestTemplateId: row.id } })
    },
    addTemplate () {
      this.$router.push({ path: '/requestManagementIndex', params: { requestTemplateId: '' } })
    },
    changePageSize (pageSize) {
      this.pagination.pageSize = pageSize
      this.requestList()
    },
    changPage (current) {
      this.pagination.currentPage = current
      this.requestList()
    },
    async requestList () {
      this.payload.filters = []
      if (this.name) {
        this.payload.filters.push({
          name: 'name',
          operator: 'contains',
          value: this.name
        })
      }
      this.payload.pageable.pageSize = this.pagination.pageSize
      this.payload.pageable.startIndex = (this.pagination.currentPage - 1) * this.pagination.pageSize
      const { statusCode, data } = await requestList(this.payload)
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
