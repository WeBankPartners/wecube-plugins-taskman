<template>
  <div style="margin: 24px">
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
            :placeholder="$t('procDefId')"
          >
            <template v-for="option in templateOptions">
              <Option :label="option.name" :value="option.id" :key="option.id"> </Option>
            </template>
          </Select>
        </Col>
        <Col span="4">
          <Select v-model="status" clearable filterable style="width:90%" :placeholder="$t('status')">
            <template v-for="option in requestStatus">
              <Option :label="option" :value="option" :key="option"> </Option>
            </template>
          </Select>
        </Col>
        <Col span="4">
          <Button @click="requestList" type="primary">{{ $t('search') }}</Button>
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
import { requestList, terminateRequest, getTemplateList } from '@/api/server'
export default {
  name: '',
  data () {
    return {
      MODALHEIGHT: 500,
      name: '',
      requestTemplate: '',
      status: '',
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
          title: this.$t('name'),
          key: 'name'
        },
        {
          title: this.$t('emergency'),
          key: 'emergency'
        },
        {
          title: this.$t('template'),
          key: 'requestTemplateName'
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
                  onClick={() => this.checkTemplate(params.row)}
                  style="margin-left: 8px"
                  type="primary"
                  size="small"
                >
                  {this.$t('look_over')}
                </Button>
                {params.row.status === 'InProgress' && (
                  <Button
                    onClick={() => this.terminateRequest(params.row)}
                    style="margin-left: 8px"
                    type="warning"
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
      tableData: [],
      requestStatus: ['Pending', 'InProgress', 'Intervene', 'Termination', 'Done']
    }
  },
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 200
    this.requestList()
  },
  methods: {
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
    async terminateRequest (row) {
      this.$Modal.confirm({
        title: this.$t('confirm_termination'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          let res = await terminateRequest(row.id)
          if (res.statusCode === 'OK') {
            this.success()
            this.requestList()
          }
        },
        onCancel: () => {}
      })
    },
    checkTemplate (row) {
      this.$router.push({
        path: '/requestManagementIndex',
        query: { requestId: row.id, requestTemplate: row.requestTemplate, isAdd: 'N', isCheck: 'Y', isHandle: 'N' }
      })
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
      this.payload.filters = [{ name: 'status', operator: 'in', value: this.requestStatus }]
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
      if (this.status) {
        this.payload.filters.push({
          name: 'status',
          operator: 'eq',
          value: this.status
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