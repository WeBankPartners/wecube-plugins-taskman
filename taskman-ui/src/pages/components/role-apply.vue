<template>
  <div>
    <Modal
      v-model="showModal"
      :mask-closable="false"
      :fullscreen="isfullscreen"
      :footer-hide="true"
      :width="1000"
      :title="$t('tw_apply_roles')"
    >
      <div slot="header" class="custom-modal-header">
        <span>
          {{ $t('tw_apply_roles') }}
        </span>
        <Icon v-if="isfullscreen" @click="isfullscreen = !isfullscreen" class="fullscreen-icon" type="ios-contract" />
        <Icon v-else @click="isfullscreen = !isfullscreen" class="fullscreen-icon" type="ios-expand" />
      </div>
      <div>
        <div class="title" style="margin-top:0px">
          <div class="title-text">
            {{ $t('tw_apply') }}
            <span class="underline"></span>
          </div>
        </div>
        <div>
          <Form :label-width="100" inline>
            <FormItem :label="$t('manageRole')" required>
              <Select
                v-model="selectedRole"
                @on-open-change="getApplyRoles"
                multiple
                filterable
                :max-tag-count="3"
                style="width:300px;margin-right:24px;"
                :placeholder="$t('tw_apply_roles')"
              >
                <Option v-for="role in roleList" :value="role.id" :key="role.id">{{ role.displayName }}</Option>
              </Select>
            </FormItem>
            <FormItem :label="$t('role_invalidDate')">
              <DatePicker
                type="datetime"
                :value="expireTime"
                @on-change="
                  val => {
                    expireTime = val
                  }
                "
                :placeholder="$t('role_invalidDatePlaceholder')"
                :options="{
                  disabledDate(date) {
                    return date && date.valueOf() < Date.now() - 86400000
                  }
                }"
                style="width:300px;margin-right:24px;"
              ></DatePicker>
              <Button type="primary" :disabled="selectedRole.length === 0" @click="apply">{{ $t('tw_apply') }}</Button>
            </FormItem>
          </Form>
        </div>
        <div class="title" style="margin-top:0px">
          <div class="title-text">
            {{ $t('tw_application_record') }}
            <span class="underline"></span>
          </div>
        </div>
        <Tabs type="card" :value="activeTab" @on-click="tabChange">
          <TabPane :label="$t('tw_pending')" name="pending"></TabPane>
          <TabPane label="生效中" name="inEffect"></TabPane>
          <TabPane label="已过期" name="expire"></TabPane>
          <TabPane label="已拒绝" name="deny"></TabPane>
        </Tabs>
        <div>
          <Table
            :height="tableHeight"
            size="small"
            :columns="this.activeTab === 'pending' ? pendingColumns : processedColumns"
            :data="tableData"
          ></Table>
        </div>
      </div>
      <div slot="footer">
        <Button @click="showModal = false">{{ $t('cancel') }}</Button>
      </div>
    </Modal>
    <!--有效期续期弹框-->
    <Modal v-model="timeModalVisible" :title="$t('tw_expireReset_tips')" :mask-closable="false">
      <Form :label-width="120">
        <FormItem :label="$t('role_invalidDate')">
          <DatePicker
            type="datetime"
            :value="modalExpireTime"
            @on-change="
              val => {
                modalExpireTime = val
              }
            "
            :placeholder="$t('role_invalidDatePlaceholder')"
            :options="{
              disabledDate(date) {
                return date && date.valueOf() < Date.now() - 86400000
              }
            }"
            style="width:300px"
          ></DatePicker>
        </FormItem>
      </Form>
      <template #footer>
        <Button @click="timeModalVisible = false">{{ $t('cancel') }}</Button>
        <Button @click="handleExtendTime" type="primary">{{ $t('confirm') }}</Button>
      </template>
    </Modal>
  </div>
</template>
<script>
import { getApplyRoles, startApply, getApplyList } from '@/api/server.js'
import dayjs from 'dayjs'
export default {
  data () {
    return {
      showModal: false,
      isfullscreen: false,
      selectedRole: [],
      expireTime: '', // 角色过期时间
      roleList: [],
      activeTab: 'pending',
      tableData: [],
      timeModalVisible: false,
      modalExpireTime: '',
      pendingColumns: [
        {
          title: this.$t('tw_account'),
          key: 'createdBy'
        },
        {
          title: this.$t('tw_apply_roles'),
          key: 'roleId',
          render: (h, params) => {
            return <div>{params.row.role.displayName}</div>
          }
        },
        {
          title: this.$t('tw_application_time'),
          key: 'createdTime'
        },
        {
          title: this.$t('role_invalidDate'),
          key: 'expireTime',
          render: (h, params) => {
            return (
              <div style={this.getExpireStyle(params.row)}>
                <span>{this.getExpireTips(params.row)}</span>
                {['preExpried', 'expire'].includes(params.row.status) && (
                  <Icon
                    type="md-time"
                    size="24"
                    color="#808695"
                    style="cursor:pointer;margin-left:5px"
                    onClick={() => {
                      this.openTimeModal(params.row)
                    }}
                  />
                )}
              </div>
            )
          }
        }
      ],
      processedColumns: [
        {
          title: this.$t('tw_apply_roles'),
          key: 'roleId',
          render: (h, params) => {
            return <div>{params.row.role.displayName}</div>
          }
        },
        {
          title: this.$t('tw_application_time'),
          key: 'createdTime'
        },
        {
          title: `${this.$t('tw_approver')}(${this.$t('tw_role_administrator')})`,
          key: 'updatedBy',
          width: 210
        },
        {
          title: this.$t('role_invalidDate'),
          key: 'expireTime',
          render: (h, params) => {
            return (
              <div style={this.getExpireStyle(params.row)}>
                <span>{this.getExpireTips(params.row)}</span>
                {['preExpried', 'expire'].includes(params.row.status) && (
                  <Icon
                    type="md-time"
                    size="24"
                    color="#808695"
                    style="cursor:pointer;margin-left:5px"
                    onClick={() => {
                      this.openTimeModal(params.row)
                    }}
                  />
                )}
              </div>
            )
          }
        }
      ]
    }
  },
  computed: {
    tableHeight () {
      const innerHeight = window.innerHeight
      return this.isfullscreen ? innerHeight - 300 : 400
    },
    getExpireStyle () {
      return function ({ status }) {
        let color = ''
        if (status === 'preExpried') {
          color = '#ff9900'
        } else if (status === 'expire') {
          color = '#ed4014'
        } else {
          color = '#19be6b'
        }
        return { color: color, display: 'flex', alignItems: 'center' }
      }
    },
    getExpireTips () {
      return function ({ status, expireTime }) {
        let text = ''
        if (status === 'preExpried') {
          text = `${expireTime}将到期`
        } else if (status === 'expire') {
          text = `${expireTime}已到期`
        } else if (expireTime) {
          text = `${expireTime}到期`
        } else if (!expireTime) {
          text = `永久有效`
        }
        return text
      }
    }
  },
  methods: {
    openModal () {
      this.showModal = true
      this.selectedRole = []
      this.getTableData()
    },
    tabChange (val) {
      this.activeTab = val
      this.getTableData()
    },
    async getTableData () {
      let statusArr = []
      if (this.activeTab === 'pending') {
        statusArr = ['init']
      } else if (this.activeTab === 'inEffect') {
        statusArr = ['inEffect']
      } else if (this.activeTab === 'expire') {
        statusArr = ['expire']
      } else if (this.activeTab === 'deny') {
        statusArr = ['deny']
      }
      const params = {
        filters: [
          {
            name: 'status',
            operator: 'in',
            value: statusArr
          }
        ],
        paging: true,
        pageable: {
          startIndex: 0,
          pageSize: 10000
        },
        sorting: [
          {
            asc: false,
            field: 'createdTime'
          }
        ]
      }
      const { status, data } = await getApplyList(params)
      if (status === 'OK') {
        this.tableData = data.contents || []
      }
    },
    async apply () {
      if (this.expireTime && !dayjs(this.expireTime).isAfter(dayjs())) {
        return this.$Message.warning(this.$t('role_invalidDateValidate'))
      }
      let data = {
        userName: localStorage.getItem('username'),
        roleIds: this.selectedRole,
        expireTime: this.expireTime
      }
      const { status } = await startApply(data)
      if (status === 'OK') {
        this.selectedRole = []
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('tw_apply_success')
        })
        this.getTableData()
      }
    },
    async getApplyRoles () {
      const params = {
        all: 'N', // Y:所有(包括未激活和已删除的) N:激活的
        roleAdmin: false
      }
      const { status, data } = await getApplyRoles(params)
      if (status === 'OK') {
        this.roleList = data || []
      }
    },
    openTimeModal (row) {
      this.timeModalVisible = true
      this.editRow = row
    },
    async handleExtendTime () {
      if (this.modalExpireTime && !dayjs(this.modalExpireTime).isAfter(dayjs())) {
        return this.$Message.warning(this.$t('role_invalidDateValidate'))
      }
      let data = {
        userName: localStorage.getItem('username'),
        roleIds: [this.editRow.role.id],
        expireTime: this.modalExpireTime
      }
      const { status } = await startApply(data)
      if (status === 'OK') {
        this.timeModalVisible = false
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('tw_apply_success')
        })
        this.getTableData()
      }
    }
  }
}
</script>
<style lang="scss" scoped>
.tagContainers {
  overflow: auto;
  height: calc(100vh - 650px);
}
.item-style {
  padding: 2px 4px;
  border: 1px dashed #e8eaec;
  margin: 6px;
  font-size: 12px;
  border-radius: 4px;
  cursor: pointer;
  width: 80%;
  display: inline-block;
}
.active-item {
  background-color: #2db7f5;
}

.title {
  font-size: 14px;
  font-weight: bold;
  margin: 12px 0;
  display: inline-block;
  .title-text {
    display: inline-block;
    margin-left: 6px;
  }
  .underline {
    display: block;
    margin-top: -10px;
    margin-left: -6px;
    width: 100%;
    padding: 0 6px;
    height: 12px;
    border-radius: 12px;
    background-color: #c6eafe;
    box-sizing: content-box;
  }
}

.custom-modal-header {
  line-height: 20px;
  font-size: 16px;
  color: #17233d;
  font-weight: 500;
  .fullscreen-icon {
    float: right;
    margin-right: 28px;
    font-size: 18px;
    cursor: pointer;
  }
}
</style>
