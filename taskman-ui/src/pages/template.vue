<template>
  <div>
    <div>
      <Row>
        <Tabs v-model="status" @on-click="onSearch">
          <!--已发布-->
          <TabPane :label="$t('status_confirm')" name="confirm"></TabPane>
          <!--草稿-->
          <TabPane :label="$t('status_draft')" name="created"></TabPane>
          <!--待管理员确认-->
          <TabPane :label="pendingLabel" name="pending"></TabPane>
          <!--已禁用-->
          <TabPane :label="$t('tw_template_status_disable')" name="disable"></TabPane>
          <!--已作废-->
          <TabPane :label="$t('tw_template_status_cancel')" name="cancel"></TabPane>
        </Tabs>
      </Row>
      <Row>
        <Col span="4">
          <Input
            v-model="name"
            style="width: 90%"
            type="text"
            :placeholder="$t('name')"
            clearable
            @on-change="handleInputChange"
          >
          </Input>
        </Col>
        <Col span="4">
          <Input
            v-model="tags"
            style="width: 90%"
            type="text"
            :placeholder="$t('tags')"
            clearable
            @on-change="handleInputChange"
          >
          </Input>
        </Col>
        <Col span="4">
          <Select
            v-model="mgmtRoles"
            @on-open-change="getInitRole"
            clearable
            filterable
            multiple
            :max-tag-count="3"
            :placeholder="$t('tw_template_owner_role')"
            style="width: 90%"
            @on-change="onSearch"
          >
            <Option v-for="item in roleOptions" :value="item.id" :key="item.id">{{ item.displayName }}</Option>
          </Select>
        </Col>
        <Col span="4">
          <Select
            v-model="type"
            clearable
            filterable
            multiple
            :placeholder="$t('tw_useScene')"
            style="width: 90%"
            @on-change="onSearch"
          >
            <Option v-for="item in typeList" :value="item.value" :key="item.value">{{ item.label }}</Option>
          </Select>
        </Col>
        <Col span="4">
          <Button @click="onSearch" type="primary">{{ $t('search') }}</Button>
          <Button @click="handleReset" type="default">{{ $t('reset') }}</Button>
        </Col>
        <div style="display:flex;float:right;">
          <Button @click="addTemplate" type="success">{{ $t('add') }}</Button>
          <Upload
            :action="uploadUrl"
            :before-upload="handleUpload"
            :show-upload-list="false"
            with-credentials
            style="margin-left:10px;"
            :headers="headers"
            :on-success="uploadSucess"
            :on-error="uploadFailed"
          >
            <Button>{{ $t('upload') }}</Button>
          </Upload>
        </div>
      </Row>
    </div>
    <Table
      style="margin: 24px 0"
      @on-sort-change="sortTable"
      size="small"
      :loading="loading"
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
    <Modal v-model="modalShow" width="300">
      <p>{{ $t('confirm_disable') }}</p>
      <template #footer>
        <Button type="error" size="large" long @click="disableInit">{{ $t('confirm') }}</Button>
      </template>
    </Modal>
  </div>
</template>

<script>
import axios from 'axios'
import { getCookie } from '@/pages/util/cookie'
import {
  getTemplateList,
  deleteTemplate,
  forkTemplate,
  copyTemplate,
  getManagementRoles,
  confirmUploadTemplate,
  enableTemplate,
  disableTemplate,
  templateGiveMe,
  updateTemplateStatus,
  templateConfirmCount
} from '@/api/server'
import { debounce } from '@/pages/util'
import ScrollTag from '@/pages/components/scroll-tag.vue'
export default {
  name: '',
  components: {
    ScrollTag
  },
  data () {
    return {
      MODALHEIGHT: 500,
      name: '',
      status: 'confirm',
      mgmtRoles: [],
      type: [],
      tags: '',
      modalShow: false,
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
      username: window.localStorage.getItem('username'),
      tableColumns: [
        {
          title: this.$t('name'),
          resizable: true,
          width: 200,
          sortable: 'custom',
          key: 'name',
          render: (h, params) => {
            return (
              <span
                style="cursor:pointer;"
                onClick={() => {
                  this.handleDbClick(params.row)
                }}
              >
                {params.row.name}
              </span>
            )
          }
        },
        {
          title: this.$t('version'),
          minWidth: 60,
          sortable: 'custom',
          key: 'version',
          render: (h, params) => {
            if (params.row.version) {
              return <Tag>{params.row.version}</Tag>
            } else {
              return <span>-</span>
            }
          }
        },
        {
          title: this.$t('tw_useScene'),
          minWidth: 80,
          key: 'type',
          render: (h, params) => {
            const find = this.typeList.find(i => i.value === String(params.row.type)) || {}
            return (find.label && <Tag>{find.label}</Tag>) || <span>-</span>
          }
        },
        {
          title: this.$t('procDefId'),
          minWidth: 100,
          sortable: 'custom',
          key: 'procDefName',
          render: (h, params) => {
            if (params.row.procDefName) {
              return (
                <span>
                  {`${params.row.procDefName}`}
                  {params.row.procDefVersion && (
                    <span style="border:1px solid #e8eaec;border-radius:3px;background:#f7f7f7;padding:1px 4px;">
                      {params.row.procDefVersion}
                    </span>
                  )}
                </span>
              )
            } else {
              return <span>-</span>
            }
          }
        },
        {
          title: this.$t('tags'),
          sortable: 'custom',
          minWidth: 130,
          key: 'tags',
          render: (h, params) => {
            if (params.row.tags) {
              return <Tag>{params.row.tags}</Tag>
            } else {
              return <span>-</span>
            }
          }
        },
        {
          title: this.$t('description'),
          resizable: true,
          minWidth: 120,
          sortable: 'custom',
          key: 'description',
          render: (h, params) => {
            return (
              <Tooltip max-width="300" content={params.row.description}>
                <span style="overflow:hidden;text-overflow:ellipsis;display:-webkit-box;-webkit-line-clamp:3;-webkit-box-orient:vertical;">
                  {params.row.description || '-'}
                </span>
              </Tooltip>
            )
          }
        },
        {
          title: this.$t('tw_template_owner_role'),
          minWidth: 120,
          key: 'mgmtRoles',
          render: (h, params) => {
            return params.row.mgmtRoles.map(item => {
              return <Tag>{item.displayName}</Tag>
            })
          }
        },
        {
          title: this.$t('useRoles'),
          minWidth: 120,
          key: 'mgmtRoles',
          render: (h, params) => {
            const list = params.row.useRoles.map(item => {
              return item.displayName
            })
            return <ScrollTag list={list} />
          }
        },
        {
          title: this.$t('updatedBy'),
          sortable: 'updatedBy',
          minWidth: 100,
          key: 'updatedBy'
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
          width: 325,
          align: 'center',
          render: (h, params) => {
            return (
              <div style="display:flex;align-items:center;justify-content:center;">
                {/* 转给我 */ this.status === 'created' && this.username !== params.row.updatedBy && (
                  <Tooltip content={this.$t('tw_action_give')} placement="top">
                    <Button
                      type="success"
                      size="small"
                      style="margin-right:5px;"
                      onClick={() => this.giveTemplate(params.row)}
                    >
                      <Icon type="ios-hand" size="16"></Icon>
                    </Button>
                  </Tooltip>
                )}
                {/* 编辑 */ this.status === 'created' && this.username === params.row.updatedBy && (
                  <Tooltip content={this.$t('edit')} placement="top">
                    <Button
                      size="small"
                      type="primary"
                      style="margin-right:5px;"
                      onClick={() => this.editTemplate(params.row)}
                    >
                      <Icon type="md-create" size="16"></Icon>
                    </Button>
                  </Tooltip>
                )}
                {/* 删除 */ this.status === 'created' && this.username === params.row.updatedBy && (
                  <Tooltip content={this.$t('delete')} placement="top">
                    <Button
                      size="small"
                      type="error"
                      style="margin-right:5px;"
                      onClick={() => this.deleteTemplate(params.row)}
                    >
                      <Icon type="md-trash" size="16"></Icon>
                    </Button>
                  </Tooltip>
                )}
                {/* 查看 */ ['pending', 'confirm', 'disable', 'cancel'].includes(this.status) && (
                  <Tooltip content={this.$t('detail')} placement="top">
                    <Button
                      size="small"
                      type="info"
                      style="margin-right:5px;"
                      onClick={() => this.checkTemplate(params.row)}
                    >
                      <Icon type="md-eye" size="16"></Icon>
                    </Button>
                  </Tooltip>
                )}
                {/* 确认发布 */ this.status === 'pending' && this.username === params.row.administrator && (
                  <Tooltip content={this.$t('tw_confirmPublish')} placement="top">
                    <Button
                      size="small"
                      type="success"
                      style="margin-right:5px;"
                      onClick={() => this.confirmTemplate(params.row)}
                    >
                      <Icon type="ios-paper-plane" size="16"></Icon>
                    </Button>
                  </Tooltip>
                )}
                {/* 退回草稿 */ this.status === 'pending' && this.username === params.row.administrator && (
                  <Tooltip content={this.$t('tw_redrawDraft')} placement="top">
                    <Button
                      size="small"
                      type="error"
                      style="margin-right:5px;"
                      onClick={() => this.draftTemplate(params.row)}
                    >
                      <Icon type="ios-redo" size="16"></Icon>
                    </Button>
                  </Tooltip>
                )}
                {/* 复制 */ this.status === 'confirm' && (
                  <Tooltip content={this.$t('copy')} placement="top">
                    <Button
                      size="small"
                      type="primary"
                      style="margin-right:5px;"
                      onClick={() => this.copyTemplate(params.row)}
                    >
                      <Icon type="md-copy" size="16"></Icon>
                    </Button>
                  </Tooltip>
                )}
                {/* 变更 */ (this.status === 'confirm' || (this.status === 'cancel' && params.row.cancelEdit)) && (
                  <Tooltip content={this.$t('fork')} placement="top">
                    <Button
                      size="small"
                      type="warning"
                      style="margin-right:5px;"
                      onClick={() => this.forkTemplate(params.row)}
                    >
                      <Icon type="md-git-branch" size="16"></Icon>
                    </Button>
                  </Tooltip>
                )}
                {/* 导出 */ this.status === 'confirm' && (
                  <Tooltip content={this.$t('download')} placement="top">
                    <Button
                      size="small"
                      type="success"
                      style="margin-right:5px;"
                      onClick={() => this.exportTemplate(params.row)}
                    >
                      <Icon type="md-cloud-download" size="16"></Icon>
                    </Button>
                  </Tooltip>
                )}
                {/* 禁用 */ this.status === 'confirm' && (
                  <Tooltip content={this.$t('disable')} placement="top">
                    <Button
                      size="small"
                      type="error"
                      style="margin-right:5px;"
                      onClick={() => this.disableTemplate(params.row)}
                    >
                      <Icon type="md-lock" size="16"></Icon>
                    </Button>
                  </Tooltip>
                )}
                {/* 作废 */ this.status === 'confirm' && (
                  <Tooltip content={this.$t('void')} placement="top">
                    <Button
                      size="small"
                      type="warning"
                      style="margin-right:5px;"
                      onClick={() => this.cancelTemplate(params.row)}
                    >
                      <Icon type="md-cut" size="16"></Icon>
                    </Button>
                  </Tooltip>
                )}
                {/* 启用 */ this.status === 'disable' && (
                  <Tooltip content={this.$t('enable')} placement="top">
                    <Button
                      size="small"
                      type="success"
                      style="margin-right:5px;"
                      onClick={() => this.enableTemplate(params.row)}
                    >
                      <Icon type="md-unlock" size="16"></Icon>
                    </Button>
                  </Tooltip>
                )}
              </div>
            )
          }
        }
      ],
      tableData: [],
      roleOptions: [],
      uploadUrl: '/taskman/api/v1/request-template/import',
      headers: {},
      backReason: '',
      confirmCount: 0, // 管理员待发布模板数量
      pendingLabel: () => {
        return (
          <div>
            <span>{this.$t('tw_waitAdminComfirm')}</span>
            {(this.confirmCount && <span class="badge">{this.confirmCount}</span>) || ''}
          </div>
        )
      },
      loading: false,
      typeList: [
        {
          value: '1',
          label: this.$t('tw_publish') // 发布
        },
        {
          value: '2',
          label: this.$t('tw_request') // 请求
        },
        {
          value: '3',
          label: this.$t('tw_question') // 问题
        },
        {
          value: '4',
          label: this.$t('tw_event') // 事件
        },
        {
          value: '5',
          label: this.$t('fork') // 变更
        }
      ]
    }
  },
  computed: {
    actionWidth () {
      if (this.status === 'confirm') {
        return 235
      } else if (this.status === 'pending') {
        return 140
      } else {
        return 120
      }
    }
  },
  watch: {
    // 给草稿页签动态插入退回原因
    status: {
      handler (val) {
        this.tableColumns.forEach(item => {
          if (item.key === 'action') {
            item.width = this.actionWidth
          }
        })
        if (val === 'created') {
          const index = this.tableColumns.findIndex(i => i.key === 'mgmtRoles')
          this.tableColumns.splice(index, 0, {
            title: this.$t('tw_rollback_reason'),
            sortable: 'updatedBy',
            minWidth: 150,
            key: 'rollbackDesc',
            render: (h, params) => {
              return (
                <Tooltip max-width="300" content={params.row.rollbackDesc}>
                  <span style="overflow:hidden;text-overflow:ellipsis;display:-webkit-box;-webkit-line-clamp:3;-webkit-box-orient:vertical;">
                    {params.row.rollbackDesc}
                  </span>
                </Tooltip>
              )
            }
          })
        } else {
          const index = this.tableColumns.findIndex(i => i.key === 'rollbackDesc')
          index > 0 && this.tableColumns.splice(index, 1)
        }
      },
      immediate: true
    }
  },
  mounted () {
    this.status = this.$route.query.status || 'confirm'
    const accessToken = getCookie('accessToken')
    this.headers = {
      Authorization: 'Bearer ' + accessToken
    }
    const lang = localStorage.getItem('lang') || 'zh-CN'
    if (lang === 'zh-CN') {
      this.headers['Accept-Language'] = 'zh-CN,zh;q=0.9,en;q=0.8'
    } else {
      this.headers['Accept-Language'] = 'en-US,en;q=0.9,zh;q=0.8'
    }
    this.MODALHEIGHT = document.body.scrollHeight - 200
    this.getTemplateList()
  },
  methods: {
    handleReset () {
      this.name = ''
      this.mgmtRoles = []
      this.type = []
      this.tags = ''
      this.onSearch()
    },
    handleInputChange: debounce(function () {
      this.onSearch()
    }, 300),
    // 转给我
    async giveTemplate (row) {
      const params = {
        request_template_id: row.id,
        latestUpdateTime: new Date(row.updatedTime).getTime().toString()
      }
      const { statusCode } = await templateGiveMe(params)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: 'Successful',
          desc: 'Successful'
        })
        this.getTemplateList()
        this.editTemplate(row)
      }
    },
    // 确认发布
    async confirmTemplate (row) {
      this.$Modal.confirm({
        title: this.$t('tw_confirmPublish'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const params = {
            requestTemplateId: row.id,
            status: 'pending',
            targetStatus: 'confirm'
          }
          const { statusCode } = await updateTemplateStatus(params)
          if (statusCode === 'OK') {
            this.$Notice.success({
              title: 'Successful',
              desc: 'Successful'
            })
            this.status = 'confirm'
            this.onSearch()
          }
        },
        onCancel: () => {}
      })
    },
    // 复制模板
    async copyTemplate (row) {
      const { statusCode } = await copyTemplate(row.id)
      if (statusCode === 'OK') {
        this.status = 'created'
        this.onSearch()
      }
    },
    // 退回草稿
    async draftTemplate (row) {
      this.$Modal.confirm({
        title: this.$t('tw_confirmRedraeDraft'),
        'z-index': 1000000,
        loading: false,
        render: () => {
          return (
            <Input
              type="textarea"
              maxlength={255}
              show-word-limit
              v-model={this.backReason}
              placeholder={this.$t('tw_back_bind_placeholder')}
            ></Input>
          )
        },
        onOk: async () => {
          if (!this.backReason.trim()) {
            this.$Notice.warning({
              title: this.$t('warning'),
              desc: this.$t('tw_back_bind_tips')
            })
          } else {
            const params = {
              requestTemplateId: row.id,
              status: 'pending',
              targetStatus: 'created',
              reason: this.backReason
            }
            const { statusCode } = await updateTemplateStatus(params)
            if (statusCode === 'OK') {
              this.$Notice.success({
                title: this.$t('successful'),
                desc: this.$t('successful')
              })
              this.status = 'created'
              this.onSearch()
            }
          }
        },
        onCancel: () => {}
      })
    },
    // 作废
    async cancelTemplate (row) {
      this.$Modal.confirm({
        title: this.$t('tw_confirmCancel'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const params = {
            requestTemplateId: row.id,
            status: 'confirm',
            targetStatus: 'cancel'
          }
          const { statusCode } = await updateTemplateStatus(params)
          if (statusCode === 'OK') {
            this.$Notice.success({
              title: 'Successful',
              desc: 'Successful'
            })
            this.status = 'cancel'
            this.onSearch()
          }
        },
        onCancel: () => {}
      })
    },
    handleUpload (file) {
      if (!file.name.endsWith('.json')) {
        this.$Notice.warning({
          title: 'Warning',
          desc: 'Must be a json file'
        })
        return false
      }
      this.$Message.info(this.$t('upload_tip'))
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
        // 确认覆盖
        this.$Modal.confirm({
          title: this.$t('confirm_import'),
          content: this.$t('tw_template_cover_tips'),
          'z-index': 1000000,
          loading: true,
          onOk: async () => {
            this.$Modal.remove()
            const { statusCode } = await confirmUploadTemplate(val.data.token)
            if (statusCode === 'OK') {
              this.$Notice.success({
                title: 'Successful',
                desc: 'Successful'
              })
              this.status = 'created'
              this.onSearch()
            }
          },
          onCancel: () => {}
        })
      } else if (val.statusCode === 'OK') {
        // 导入成功
        this.$Notice.success({
          title: 'Successful',
          desc: 'Successful'
        })
        this.status = 'created'
        this.onSearch()
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
    // 变更
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
            this.status = 'created'
            this.onSearch()
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
      this.$router.push({
        path: '/templateManagementIndex',
        query: { requestTemplateId: row.id, isCheck: 'N', parentStatus: this.status }
      })
    },
    checkTemplate (row) {
      this.$router.push({
        path: '/templateManagementIndex',
        query: { requestTemplateId: row.id, isCheck: 'Y', parentStatus: this.status }
      })
    },
    addTemplate () {
      this.$router.push({
        path: '/templateManagementIndex',
        params: { requestTemplateId: '', isCheck: 'N', parentStatus: 'created' }
      })
    },
    onSearch () {
      this.pagination.currentPage = 1
      this.getTemplateList()
    },
    changePageSize (pageSize) {
      this.pagination.currentPage = 1
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
          operator: 'in',
          value: [this.status]
        })
      }
      if (this.mgmtRoles.length > 0) {
        this.payload.filters.push({
          name: 'mgmtRoles',
          operator: 'eq',
          value: this.mgmtRoles
        })
      }
      if (this.type.length > 0) {
        this.payload.filters.push({
          name: 'type',
          operator: 'eq',
          value: this.type
        })
      }
      if (sorting) {
        this.payload.sorting = sorting
      }
      this.payload.pageable.pageSize = this.pagination.pageSize
      this.payload.pageable.startIndex = (this.pagination.currentPage - 1) * this.pagination.pageSize
      this.loading = true
      const { statusCode, data } = await getTemplateList(this.payload)
      this.loading = false
      if (statusCode === 'OK') {
        this.tableData = data.contents
        this.pagination.total = data.pageInfo.totalRows
      }
      this.getConfirmCount()
    },
    disableTemplate (row) {
      this.modalShow = row.id
    },
    async disableInit () {
      const res = await disableTemplate(this.modalShow)
      this.modalShow = false
      if (res.statusCode === 'OK') {
        this.success()
        this.status = 'disable'
        this.onSearch()
      }
    },
    async enableTemplate (row) {
      const res = await enableTemplate(row.id)
      if (res.statusCode === 'OK') {
        this.success()
        this.status = 'confirm'
        this.onSearch()
      }
    },
    async getConfirmCount () {
      const { statusCode, data } = await templateConfirmCount()
      if (statusCode === 'OK') {
        this.confirmCount = data
      }
    },
    // 点击名称快捷跳转
    handleDbClick (row) {
      if (this.status === 'created') {
        this.editTemplate(row)
      } else {
        this.checkTemplate(row)
      }
    }
  }
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
.badge {
  position: absolute;
  display: inline-block;
  font-size: 11px;
  background-color: #f56c6c;
  border-radius: 10px;
  color: #fff;
  height: 18px;
  line-height: 18px;
  padding: 0 6px;
  text-align: center;
  white-space: nowrap;
  margin-left: 5px;
}
</style>
