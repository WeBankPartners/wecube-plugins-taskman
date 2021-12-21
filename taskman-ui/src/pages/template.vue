<template>
  <div>
    <div>
      <Row>
        <Col span="4">
          <Input v-model="name" style="width:90%" type="text" :placeholder="$t('name')"> </Input>
        </Col>
        <Col span="4">
          <Input v-model="tags" style="width:90%" type="text" :placeholder="$t('tags')"> </Input>
        </Col>
        <Col span="4">
          <Select
            v-model="mgmtRoles"
            @on-open-change="getInitRole"
            clearable
            filterable
            multiple
            :placeholder="$t('manageRole')"
            style="width:90%"
          >
            <Option v-for="item in roleOptions" :value="item.id" :key="item.id">{{ item.displayName }}</Option>
          </Select>
        </Col>
        <Col span="2">
          <Select v-model="status" clearable :placeholder="$t('status')" style="width:90%">
            <Option value="confirm" key="confirm">{{ this.$t('status_confirm') }}</Option>
            <Option value="created" key="created">{{ this.$t('status_created') }}</Option>
          </Select>
        </Col>
        <Col span="4">
          <Button @click="getTemplateList" type="primary">{{ $t('search') }}</Button>
          <Button @click="addTemplate" type="success">{{ $t('add') }}</Button>
        </Col>
        <Upload
          :action="uploadUrl"
          :before-upload="handleUpload"
          :show-upload-list="false"
          with-credentials
          style="display:inline-block;float:right;margin-right:16px"
          :headers="headers"
          :on-success="uploadSucess"
          :on-error="uploadFailed"
        >
          <Button>{{ $t('upload') }}</Button>
        </Upload>
      </Row>
    </div>
    <Table
      style="margin: 24px 0;"
      border
      @on-sort-change="sortTable"
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
import axios from 'axios'
import { getCookie } from '@/pages/util/cookie'
import { getTemplateList, deleteTemplate, forkTemplate, getManagementRoles, confirmUploadTemplate } from '@/api/server'
export default {
  name: '',
  data () {
    return {
      MODALHEIGHT: 500,
      name: '',
      status: '',
      mgmtRoles: [],
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
          resizable: true,
          width: 200,
          sortable: 'custom',
          fixed: 'left',
          key: 'name'
        },
        {
          title: this.$t('version'),
          minWidth: 80,
          sortable: 'custom',
          key: 'version'
        },
        {
          title: this.$t('tags'),
          sortable: 'custom',
          minWidth: 130,
          key: 'tags'
        },
        {
          title: this.$t('status'),
          minWidth: 80,
          sortable: 'custom',
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
          resizable: true,
          width: 300,
          sortable: 'custom',
          key: 'description'
        },
        {
          title: this.$t('mgmtRoles'),
          minWidth: 130,
          key: 'mgmtRoles',
          render: (h, params) => {
            const displayName = params.row.mgmtRoles.map(role => role.displayName).join(',')
            return <span>{displayName}</span>
          }
        },
        {
          title: this.$t('useRoles'),
          minWidth: 130,
          key: 'mgmtRoles',
          render: (h, params) => {
            const displayName = params.row.useRoles.map(role => role.displayName).join(',')
            return <span>{displayName}</span>
          }
        },
        {
          title: this.$t('tm_updated_time'),
          sortable: 'custom',
          minWidth: 130,
          key: 'updatedTime'
        },
        {
          title: this.$t('t_action'),
          key: 'action',
          fixed: 'right',
          width: 180,
          align: 'center',
          render: (h, params) => {
            const operationOptions = params.row.operateOptions
            return (
              <div style="text-align: left">
                {operationOptions.includes('edit') && (
                  <Button
                    onClick={() => this.editTemplate(params.row)}
                    style="margin-left: 6px"
                    type="primary"
                    size="small"
                  >
                    {this.$t('edit')}
                  </Button>
                )}
                {operationOptions.includes('delete') && (
                  <Button
                    onClick={() => this.deleteTemplate(params.row)}
                    style="margin-left: 6px"
                    type="error"
                    size="small"
                  >
                    {this.$t('delete')}
                  </Button>
                )}
                {operationOptions.includes('query') && (
                  <Button
                    onClick={() => this.checkTemplate(params.row)}
                    style="margin-left: 6px"
                    type="info"
                    size="small"
                  >
                    {this.$t('detail')}
                  </Button>
                )}
                {operationOptions.includes('fork') && (
                  <Button
                    onClick={() => this.forkTemplate(params.row)}
                    style="margin-left: 6px"
                    type="warning"
                    size="small"
                  >
                    {this.$t('fork')}
                  </Button>
                )}
                {operationOptions.includes('export') && (
                  <Button
                    onClick={() => this.exportTemplate(params.row)}
                    style="margin-left: 6px"
                    type="success"
                    size="small"
                  >
                    {this.$t('download')}
                  </Button>
                )}
              </div>
            )
          }
        }
      ],
      tableData: [],
      roleOptions: [],
      uploadUrl: '/taskman/api/v1/request-template/import',
      headers: {}
    }
  },
  mounted () {
    const accessToken = getCookie('accessToken')
    this.headers = {
      Authorization: 'Bearer ' + accessToken
    }
    this.MODALHEIGHT = document.body.scrollHeight - 200
    this.getTemplateList()
  },
  methods: {
    handleUpload (file) {
      if (!file.name.endsWith('.json')) {
        this.$Notice.warning({
          title: 'Warning',
          desc: 'Must be a json file'
        })
        return false
      }
      return true
    },
    uploadFailed (val, response) {
      console.log(val)
      this.$Notice.error({
        title: 'Error',
        desc: response.statusMessage
      })
    },
    async uploadSucess (val) {
      if (val.statusCode === 'CONFIRM') {
        this.$Modal.confirm({
          title: this.$t('confirm_import'),
          'z-index': 1000000,
          loading: true,
          onOk: async () => {
            this.$Modal.remove()
            const { statusCode } = await confirmUploadTemplate(val.data)
            if (statusCode === 'OK') {
              this.$Notice.success({
                title: 'Successful',
                desc: 'Successful'
              })
              this.getTemplateList()
            }
          },
          onCancel: () => {}
        })
      } else if (val.statusCode === 'OK') {
        this.$Notice.success({
          title: 'Successful',
          desc: 'Successful'
        })
        this.getTemplateList()
      } else {
        this.$Notice.warning({
          title: 'Warning',
          desc: val.statusMessage
        })
      }
    },
    async exportTemplate (row) {
      axios({
        method: 'GET',
        url: `/taskman/api/v1/request-template/export/${row.id}`,
        headers: this.headers,
        responseType: 'blob'
      })
        .then(response => {
          this.isExport = false
          if (response.status < 400) {
            let fileName = `${row.name}.json`
            let blob = new Blob([response.data])
            if ('msSaveOrOpenBlob' in navigator) {
              window.navigator.msSaveOrOpenBlob(blob, fileName)
            } else {
              if ('download' in document.createElement('a')) {
                // 非IE下载
                let elink = document.createElement('a')
                elink.download = fileName
                elink.style.display = 'none'
                elink.href = URL.createObjectURL(blob)
                document.body.appendChild(elink)
                elink.click()
                URL.revokeObjectURL(elink.href) // 释放URL 对象
                document.body.removeChild(elink)
              } else {
                // IE10+下载
                navigator.msSaveOrOpenBlob(blob, fileName)
              }
            }
          }
        })
        .catch(() => {
          this.$Message.warning('Error')
        })
    },
    sortTable (col) {
      const sorting = {
        asc: col.order === 'asc',
        field: col.key
      }
      this.getTemplateList(sorting)
    },
    async getInitRole () {
      const { statusCode, data } = await getManagementRoles()
      if (statusCode === 'OK') {
        this.roleOptions = data
      }
    },
    success () {
      this.$Notice.success({
        title: this.$t('successful'),
        desc: this.$t('successful')
      })
    },
    forkTemplate (row) {
      this.$Modal.confirm({
        title: this.$t('confirm_change'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          let res = await forkTemplate(row.id)
          if (res.statusCode === 'OK') {
            this.success()
            this.getTemplateList()
          }
        },
        onCancel: () => {}
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
      this.$router.push({ path: '/templateManagementIndex', query: { requestTemplateId: row.id, isCheck: 'N' } })
    },
    checkTemplate (row) {
      this.$router.push({ path: '/templateManagementIndex', query: { requestTemplateId: row.id, isCheck: 'Y' } })
    },
    addTemplate () {
      this.$router.push({ path: '/templateManagementIndex', params: { requestTemplateId: '', isCheck: 'N' } })
    },
    changePageSize (pageSize) {
      this.pagination.pageSize = pageSize
      this.getTemplateList()
    },
    changPage (current) {
      this.pagination.currentPage = current
      this.getTemplateList()
    },
    async getTemplateList (sorting) {
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
      if (this.status) {
        this.payload.filters.push({
          name: 'status',
          operator: 'eq',
          value: this.status
        })
      }
      if (this.mgmtRoles.length > 0) {
        this.payload.filters.push({
          name: 'mgmtRoles',
          operator: 'eq',
          value: this.mgmtRoles
        })
      }
      if (sorting) {
        this.payload.sorting = sorting
      }
      this.payload.pageable.pageSize = this.pagination.pageSize
      this.payload.pageable.startIndex = (this.pagination.currentPage - 1) * this.pagination.pageSize
      const { statusCode, data } = await getTemplateList(this.payload)
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
