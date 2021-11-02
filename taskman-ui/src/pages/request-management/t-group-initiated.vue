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
          <Select v-model="status" multiple clearable filterable style="width:90%" :placeholder="$t('status')">
            <template v-for="option in requestStatus">
              <Option :label="option" :value="option" :key="option"> </Option>
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
import { requestListForDraftInitiated, getTemplateList } from '@/api/server'
export default {
  name: '',
  data () {
    return {
      MODALHEIGHT: 500,
      name: '',
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
          key: 'id'
        },
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
          title: this.$t('handle_role'),
          key: 'handleRoles',
          render: (h, params) => {
            const handleRoles = params.row.handleRoles.length === 1 ? params.row.handleRoles[0] : ''
            return <span>{handleRoles}</span>
          }
        },
        {
          title: this.$t('handler'),
          key: 'handler'
        },
        {
          title: this.$t('status'),
          key: 'status'
        },
        {
          title: this.$t('estimated_finish_time'),
          key: 'expireTime'
        },
        {
          title: this.$t('expected_completion_time'),
          key: 'expectTime'
        },
        {
          title: this.$t('report_time'),
          key: 'reportTime'
        },
        {
          title: this.$t('t_action'),
          key: 'action',
          width: 100,
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
              </div>
            )
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
        'Timeouted',
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
    changePageSize (pageSize) {
      this.pagination.pageSize = pageSize
      this.requestListForDraftInitiated()
    },
    changPage (current) {
      this.pagination.currentPage = current
      this.requestListForDraftInitiated()
    },
    async requestListForDraftInitiated () {
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
      if (this.status.length > 0) {
        this.payload.filters.push({
          name: 'status',
          operator: 'in',
          value: this.status
        })
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

<style scoped lang="scss">
.header-icon {
  float: right;
  margin: 3px 40px 0 0 !important;
}
</style>
