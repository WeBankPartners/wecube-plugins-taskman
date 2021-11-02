<template>
  <div style="margin: 24px">
    <div>
      <Row>
        <Col span="4">
          <Input v-model="name" style="width:90%" type="text" :placeholder="$t('name')"> </Input>
        </Col>
        <Col span="4">
          <Input v-model="tags" style="width:90%" type="text" :placeholder="$t('tags')"> </Input>
        </Col>
        <Col span="4">
          <Button @click="getTemplateList" type="primary">{{ $t('search') }}</Button>
          <Button @click="addTemplate">{{ $t('add') }}</Button>
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
import { getTemplateList, deleteTemplate } from '@/api/server'
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
          title: this.$t('version'),
          key: 'version'
        },
        {
          title: this.$t('tags'),
          key: 'tags'
        },
        {
          title: this.$t('status'),
          key: 'status',
          render: (h, params) => {
            const statusArray = {
              confirm: this.$t('status_confirm'),
              created: this.$t('status_created')
            }
            return <span>{statusArray[params.row.status]}</span>
          }
        },
        {
          title: this.$t('description'),
          key: 'description'
        },
        {
          title: this.$t('mgmtRoles'),
          key: 'mgmtRoles',
          render: (h, params) => {
            const displayName = params.row.mgmtRoles.map(role => role.displayName).join(',')
            return <span>{displayName}</span>
          }
        },
        {
          title: this.$t('useRoles'),
          key: 'mgmtRoles',
          render: (h, params) => {
            const displayName = params.row.useRoles.map(role => role.displayName).join(',')
            return <span>{displayName}</span>
          }
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
            return h('div', [
              h(
                'Button',
                {
                  props: {
                    type: 'primary',
                    size: 'small'
                  },
                  style: {
                    'margin-left': '8px'
                  },
                  on: {
                    click: () => {
                      this.editTemplate(params.row)
                    }
                  }
                },
                this.$t('edit')
              ),
              h(
                'Button',
                {
                  props: {
                    type: 'error',
                    size: 'small'
                  },
                  style: {
                    'margin-left': '8px'
                  },
                  on: {
                    click: () => {
                      this.deleteTemplate(params.row)
                    }
                  }
                },
                this.$t('delete')
              )
            ])
          }
        }
      ],
      tableData: []
    }
  },
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 200
    this.getTemplateList()
  },
  methods: {
    success () {
      this.$Notice.success({
        title: this.$t('successful'),
        desc: this.$t('successful')
      })
    },
    deleteTemplate (row) {
      this.$Modal.confirm({
        title: this.$t('confirm_delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const params = {
            params: { id: row.id }
          }
          let res = await deleteTemplate(params)
          if (res.statusCode === 'OK') {
            this.success()
            this.getTemplateList()
          }
        },
        onCancel: () => {}
      })
    },
    editTemplate (row) {
      this.$router.push({ path: '/templateManagementIndex', query: { requestTemplateId: row.id } })
    },
    addTemplate () {
      this.$router.push({ path: '/templateManagementIndex', params: { requestTemplateId: '' } })
    },
    changePageSize (pageSize) {
      this.pagination.pageSize = pageSize
      this.getTemplateList()
    },
    changPage (current) {
      this.pagination.currentPage = current
      this.getTemplateList()
    },
    async getTemplateList () {
      this.payload.filters = []
      if (this.name) {
        this.payload.filters.push({
          name: 'name',
          operator: 'contains',
          value: this.name
        })
      }
      if (this.tags) {
        this.payload.filters.push({
          name: 'tags',
          operator: 'contains',
          value: this.tags
        })
      }
      this.payload.pageable.pageSize = this.pagination.pageSize
      this.payload.pageable.startIndex = (this.pagination.currentPage - 1) * this.pagination.pageSize
      this.$Spin.show()
      const { statusCode, data } = await getTemplateList(this.payload)
      if (statusCode === 'OK') {
        this.$Spin.hide()
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
