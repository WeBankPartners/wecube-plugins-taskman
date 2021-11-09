<template>
  <div>
    <div>
      <Row>
        <Col span="4">
          <Input v-model="name" style="width:90%" type="text" :placeholder="$t('name')"> </Input>
        </Col>
        <Col span="4">
          <Select
            v-model="manageRole"
            @on-open-change="getInitRole"
            clearable
            filterable
            :placeholder="$t('manageRole')"
            style="width:90%"
          >
            <Option v-for="item in roleOptions" :value="item.id" :key="item.id">{{ item.displayName }}</Option>
          </Select>
        </Col>
        <Col span="4">
          <Button @click="getTempGroupList" type="primary">{{ $t('search') }}</Button>
          <Button @click="addTempGroup">{{ $t('add') }}</Button>
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
    <VVModel :modalConfig="modalConfig" @saveModel="saveModel"></VVModel>
  </div>
</template>

<script>
import VVModel from '@/pages/components/vv-model.vue'
import { getTempGroupList, getManagementRoles, createTempGroup, updateTempGroup, deleteTempGroup } from '@/api/server'
export default {
  name: '',
  data () {
    return {
      MODALHEIGHT: 500,
      name: '',
      manageRole: '',
      roleOptions: [],
      modalConfig: {
        isShow: false,
        isAdd: true,
        title: '',
        itemConfigs: [
          { label: 'name', value: 'name', rules: 'required', type: 'text' },
          {
            label: 'manageRole',
            value: 'manageRole',
            rules: 'required',
            options: 'roleOptions',
            labelKey: 'displayName',
            valueKey: 'id',
            multiple: false,
            type: 'select',
            placeholder: ''
          },
          { label: 'description', value: 'description', rows: 2, type: 'textarea' }
        ],
        values: {
          name: '',
          manageRole: '',
          description: ''
        },
        roleOptions: []
      },
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
          field: 'catId'
        }
      },
      tableColumns: [
        {
          title: this.$t('name'),
          resizable: true,
          width: 400,
          key: 'name'
        },
        {
          title: this.$t('manageRole'),
          key: 'manageRole',
          render: (h, params) => {
            const displayName = params.row.manageRoleObj.displayName
            return <span>{displayName}</span>
          }
        },
        {
          title: this.$t('description'),
          resizable: true,
          width: 300,
          key: 'description'
        },
        {
          title: this.$t('tm_updated_time'),
          key: 'updatedTime'
        },
        {
          title: this.$t('t_action'),
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
                      this.editTemp(params.row)
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
                      this.deleteTemp(params.row)
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
    this.getTempGroupList()
  },
  methods: {
    async getInitRole () {
      const { statusCode, data } = await getManagementRoles()
      if (statusCode === 'OK') {
        this.roleOptions = data
      }
    },
    async getRole () {
      const { statusCode, data } = await getManagementRoles()
      if (statusCode === 'OK') {
        this.modalConfig.roleOptions = data
        this.modalConfig.isShow = true
      }
    },
    success () {
      this.$Notice.success({
        title: this.$t('successful'),
        desc: this.$t('successful')
      })
    },
    deleteTemp (row) {
      this.$Modal.confirm({
        title: this.$t('confirm_delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const params = {
            params: { id: row.id }
          }
          let res = await deleteTempGroup(params)
          if (res.statusCode === 'OK') {
            this.success()
            this.getTempGroupList()
          }
        },
        onCancel: () => {}
      })
    },
    async saveModel (modelData) {
      const method = this.modalConfig.isAdd ? createTempGroup : updateTempGroup
      const { statusCode } = await method(modelData)
      if (statusCode === 'OK') {
        this.success()
        this.getTempGroupList()
        this.modalConfig.isShow = false
      }
    },
    async addTempGroup () {
      this.modalConfig.isAdd = true
      this.modalConfig.title = this.$t('add') + this.$t('tm_template_group')
      this.getRole()
    },
    async editTemp (row) {
      this.modalConfig.values = { ...row }
      this.modalConfig.title = this.$t('edit') + ':' + row.name
      this.modalConfig.isAdd = false
      this.getRole()
    },
    changePageSize (pageSize) {
      this.pagination.pageSize = pageSize
      this.getTempGroupList()
    },
    changPage (current) {
      this.pagination.currentPage = current
      this.getTempGroupList()
    },
    async getTempGroupList () {
      this.payload.filters = []
      if (this.name) {
        this.payload.filters.push({
          name: 'name',
          operator: 'contains',
          value: this.name
        })
      }
      if (this.manageRole) {
        this.payload.filters.push({
          name: 'manageRole',
          operator: 'eq',
          value: this.manageRole
        })
      }
      this.payload.pageable.pageSize = this.pagination.pageSize
      this.payload.pageable.startIndex = (this.pagination.currentPage - 1) * this.pagination.pageSize
      const { statusCode, data } = await getTempGroupList(this.payload)
      if (statusCode === 'OK') {
        this.tableData = data.contents
        this.pagination.total = data.pageInfo.totalRows
      }
    }
  },
  components: {
    VVModel
  }
}
</script>

<style scoped lang="scss">
.header-icon {
  float: right;
  margin: 3px 40px 0 0 !important;
}
</style>
