<template>
  <div>
    <Row type="flex">
      <Col span="24">
        <Row>
          <Col span="12">
            <div class="basci-info-left">
              <div class="title">
                <div class="title-text">
                  1.{{ $t('tw_template_information') }}
                  <span class="underline"></span>
                </div>
              </div>
              <div class="basci-info-content">
                <Form :label-width="120">
                  <FormItem :label="$t('name')">
                    <Input
                      v-model="basicInfo.name"
                      maxlength="30"
                      show-word-limit
                      @on-change="paramsChanged"
                      style="width: 96%"
                    />
                    <!-- <Input v-model="basicInfo.name" style="width: 96%" @on-change="paramsChanged"></Input> -->
                    <span style="color: red">*</span>
                    <div
                      v-if="isParmasChanged && (basicInfo.name.length === 0 || basicInfo.name.length > 30)"
                      style="color: red"
                    >
                      {{ $t('name') }}{{ $t('tw_limit_30') }}
                    </div>
                  </FormItem>
                  <FormItem :label="$t('group')">
                    <Select
                      v-model="basicInfo.group"
                      style="width: 96%;"
                      filterable
                      @on-change="paramsChanged"
                      @on-open-change="getGroupOptions"
                    >
                      <Option v-for="item in groupOptions" :value="item.id" :key="item.id">{{ item.name }}</Option>
                    </Select>
                    <span style="color: red">*</span>
                    <div v-if="isParmasChanged && basicInfo.group === ''" style="color: red">
                      {{ $t('group') }}{{ $t('can_not_be_empty') }}
                    </div>
                  </FormItem>
                  <FormItem :label="$t('scene_type')">
                    <Select v-model="basicInfo.type" style="width: 96%;" filterable @on-change="paramsChanged">
                      <Option v-for="item in typeOptions" :value="item.value" :key="item.label">{{
                        item.label
                      }}</Option>
                    </Select>
                    <span style="color: red">*</span>
                  </FormItem>
                  <!-- 属主角色 -->
                  <FormItem :label="$t('mgmtRolesNew')">
                    <Select
                      v-model="basicInfo.mgmtRoles"
                      @on-open-change="getManagementRoles(false)"
                      filterable
                      style="width: 96%;"
                      @on-change="changeMgmtRole"
                    >
                      <Option v-for="item in mgmtRolesOptions" :value="item.id" :key="item.id">{{
                        item.displayName
                      }}</Option>
                    </Select>
                    <span style="color: red">*</span>
                    <div v-if="isParmasChanged && basicInfo.mgmtRoles === ''" style="color: red">
                      {{ $t('mgmtRolesNew') }}{{ $t('can_not_be_empty') }}
                    </div>
                  </FormItem>
                  <!-- 模版属主 -->
                  <FormItem :label="$t('handlerNew')">
                    <Select
                      v-model="basicInfo.handler"
                      filterable
                      style="width: 96%;"
                      @on-change="paramsChanged"
                      @on-open-change="getHandlerRoles(false)"
                    >
                      <Option v-for="item in handlerRolesOptions" :value="item.id" :key="item.id">{{
                        item.displayName
                      }}</Option>
                    </Select>
                  </FormItem>
                  <!-- 使用角色 -->
                  <FormItem :label="$t('useRoles')">
                    <Select
                      v-model="basicInfo.useRoles"
                      @on-open-change="getUserRoles(false)"
                      filterable
                      multiple
                      style="width: 96%;"
                    >
                      <Option v-for="item in useRolesOptions" :value="item.id" :key="item.id">{{
                        item.displayName
                      }}</Option>
                    </Select>
                    <span style="color: red">*</span>
                    <div v-if="isParmasChanged && basicInfo.useRoles.length === 0" style="color: red">
                      {{ $t('useRoles') }}{{ $t('can_not_be_empty') }}
                    </div>
                  </FormItem>
                  <!-- 请求时效 -->
                  <FormItem :label="$t('request_time_limit')">
                    <Select v-model="basicInfo.expireDay" filterable style="width: 96%;" @on-change="paramsChanged">
                      <Option v-for="item in expireDayOptions" :value="item" :key="item"
                        >{{ item }}{{ $t('day') }}</Option
                      >
                    </Select>
                    <span style="color: red">*</span>
                  </FormItem>
                  <!-- 标签 -->
                  <FormItem :label="$t('tags')">
                    <Select
                      v-model="basicInfo.tags"
                      @on-open-change="getTags"
                      filterable
                      allow-create
                      style="width: 96%;"
                      @on-change="paramsChanged"
                      @on-create="handleCreate"
                    >
                      <Option v-for="(item, tagIndex) in tagOptions" :value="item.value" :key="tagIndex">{{
                        item.label
                      }}</Option>
                    </Select>
                  </FormItem>
                  <!-- 描述 -->
                  <FormItem :label="$t('description')">
                    <Input
                      v-model="basicInfo.description"
                      maxlength="200"
                      show-word-limit
                      type="textarea"
                      @on-change="paramsChanged"
                      style="width: 96%"
                    />
                    <div v-if="isParmasChanged && basicInfo.description.length > 200" style="color: red">
                      {{ $t('description') }}{{ $t('tw_limit_200') }}
                    </div>
                  </FormItem>
                </Form>
              </div>
            </div>
          </Col>
          <Col span="12">
            <div class="basci-info-right">
              <div class="title">
                <div class="title-text">
                  2.{{ $t('tw_template_configuration') }}
                  <span class="underline"></span>
                </div>
              </div>
              <div class="basci-info-content">
                <Form :label-width="120">
                  <FormItem :label="$t('tw_start_node')">
                    <i-switch v-model="basicInfo.pendingSwitch" @on-change="pendingSwitchChange" />
                    <Tooltip content="">
                      <div slot="content" style="white-space: normal;">
                        {{ $t('tw_versioning_node_tip') }}
                      </div>
                      <Icon type="ios-alert-outline" />
                    </Tooltip>
                  </FormItem>
                  <template v-if="basicInfo.pendingSwitch">
                    <div>
                      <!-- 处理角色 -->
                      <FormItem :label="$t('handle_role')">
                        <Select
                          v-model="basicInfo.pendingRole"
                          filterable
                          style="width: 96%;"
                          @on-change="pendingRoleChange"
                          @on-open-change="getUserRoles(false)"
                        >
                          <Option v-for="item in useRolesOptions" :value="item.id" :key="item.id">{{
                            item.displayName
                          }}</Option>
                        </Select>
                      </FormItem>
                      <!-- 处理人 -->
                      <FormItem :label="$t('handler')">
                        <Select
                          v-model="basicInfo.pendingHandler"
                          filterable
                          style="width: 96%;"
                          @on-open-change="getPendingHandlerRoles(false)"
                        >
                          <Option v-for="item in pendingHandlerOptions" :value="item.id" :key="item.id">{{
                            item.displayName
                          }}</Option>
                        </Select>
                      </FormItem>
                      <!-- 节点时效 -->
                      <FormItem :label="$t('tw_node_validity_period')">
                        <Select v-model="basicInfo.pendingExpireDay" filterable style="width: 96%;">
                          <Option v-for="item in expireDayOptions" :value="item" :key="item"
                            >{{ item }}{{ $t('day') }}</Option
                          >
                        </Select>
                        <span style="color: red">*</span>
                      </FormItem>
                    </div>
                  </template>
                  <!-- [中间]关联编排 -->
                  <FormItem :label="$t('tw_middle_node')">
                    <i-switch v-model="showFlow" @on-change="showSwitchChange" />
                  </FormItem>
                  <template v-if="showFlow">
                    <div>
                      <!-- 编排 -->
                      <FormItem :label="$t('procDefId')">
                        <Select
                          v-model="basicInfo.procDefId"
                          filterable
                          style="width: 96%;"
                          :disabled="basicInfo.mgmtRoles === ''"
                          @on-open-change="getProcess"
                        >
                          <Option
                            v-for="item in procOptions"
                            :value="item.procDefId"
                            :key="item.procDefId"
                            :label="item.procDefName + ' [' + item.version + ']'"
                          >
                            <span>{{ item.procDefName }} [{{ item.version }}]</span>
                          </Option>
                        </Select>
                        <span style="color: red">*</span>
                        <div v-if="isParmasChanged && basicInfo.procDefId === ''" style="color: red">
                          {{ $t('procDefId') }}{{ $t('can_not_be_empty') }}
                        </div>
                      </FormItem>
                    </div>
                  </template>
                  <!-- [结束]确认节点 -->
                  <FormItem :label="$t('tw_end_node')">
                    <i-switch v-model="basicInfo.confirmSwitch" @on-change="basicInfo.confirmExpireDay = 1" />
                  </FormItem>
                  <template v-if="basicInfo.confirmSwitch">
                    <div>
                      <FormItem :label="$t('tw_node_validity_period')">
                        <Select v-model="basicInfo.confirmExpireDay" filterable style="width: 96%;">
                          <Option v-for="item in expireDayOptions" :value="item" :key="item"
                            >{{ item }}{{ $t('day') }}</Option
                          >
                        </Select>
                        <span style="color: red">*</span>
                      </FormItem>
                    </div>
                  </template>
                </Form>
              </div>
            </div>
          </Col>
        </Row>
      </Col>
    </Row>
    <div class="footer">
      <div class="content">
        <Button
          v-if="isCheck !== 'Y'"
          @click="createTemp(false)"
          type="info"
          :disabled="isSaveBtnDisable"
          class="btn-footer-margin"
          >{{ $t('save') }}</Button
        >
        <Button
          @click="gotoNext"
          type="primary"
          :disabled="isCheck !== 'Y' && isSaveBtnDisable"
          class="btn-footer-margin"
          >{{ $t('next') }}</Button
        >
      </div>
    </div>
  </div>
</template>

<script>
import dayjs from 'dayjs'
import {
  getTempGroupList,
  getManagementRoles,
  getTemplateTags,
  getHandlerRoles,
  getUserRoles,
  getProcess,
  createTemp,
  updateTemp,
  getTemplateList
} from '@/api/server'
export default {
  name: 'BasicInfo',
  data () {
    return {
      isParmasChanged: false, // 参数变化标志位，控制右侧panel显示逻辑
      basicInfo: {
        id: '',
        name: '',
        group: '', // 模版组
        mgmtRoles: '', // 属主角色
        handler: '', // 模版属主
        useRoles: [], // 使用角色
        description: '',
        tags: '', // 标签
        type: 2, // 模版场景类型
        expireDay: 1, // 请求时效
        pendingSwitch: true, // 定版节点
        pendingRole: '', // 定版角色
        pendingHandler: '', // 定版处理人
        pendingExpireDay: 1, // 定版节点时效
        procDefId: '',
        confirmSwitch: false, // 结束确认节点
        confirmExpireDay: 1 // 确认过期时间

        // packageName: '',
        // entityName: '',
        // procDefKey: '',
        // procDefName: ''
      },
      groupOptions: [], // 组列表
      typeOptions: [
        { label: this.$t('tw_request'), value: 2 },
        { label: this.$t('tw_publish'), value: 1 },
        { label: this.$t('tw_question'), value: 3 },
        { label: this.$t('tw_event'), value: 4 },
        { label: this.$t('fork'), value: 5 }
      ],
      tagOptions: [], // 待选标签列表
      tmpTagOptions: [], // 缓存标签列表，供新增时使用
      mgmtRolesOptions: [], // 属主角色列表
      handlerRolesOptions: [], // 模版属主：数据角色下的用户
      useRolesOptions: [], // 使用角色
      expireDayOptions: [1, 2, 3, 4, 5, 6, 7], // 请求时效
      pendingHandlerOptions: [], // 处理人
      showFlow: false, // 是否显示编排选项
      procOptions: [], // 编排列表
      isSaveBtnDisable: true
    }
  },
  watch: {
    showFlow: {
      handler (val) {
        this.isSaveBtnDisable = this.isSaveBtnActive()
      }
    },
    basicInfo: {
      handler (val) {
        this.isSaveBtnDisable = this.isSaveBtnActive()
      },
      immediate: true,
      deep: true
    }
  },
  props: ['isCheck', 'requestTemplateId'],
  mounted () {
    this.basicInfo.name = `${this.$t('tw_template')}_${dayjs().format('YYMMDDHHmmss')}`
    this.loadPage(this.requestTemplateId)
  },
  methods: {
    async loadPage (val) {
      this.requestTemplateId = val
      this.isParmasChanged = false
      this.getProcess()
      await this.getManagementRoles(true)
      this.getTemplateData()
      this.getGroupOptions()
      this.getUserRoles(true)
    },
    // 控制保存按钮
    isSaveBtnActive () {
      let res = false
      if (this.basicInfo.name.length === 0 || this.basicInfo.name.length > 30) {
        return true
      }
      if (this.basicInfo.description.length > 200) {
        return true
      }
      if (this.basicInfo.group === '') {
        return true
      }
      if (this.basicInfo.mgmtRoles === '') {
        return true
      }
      if (this.basicInfo.useRoles.length === 0) {
        return true
      }
      if (this.showFlow && !this.basicInfo.procDefId) {
        return true
      }
      return res
    },
    async createTemp (isGoToNext) {
      if (this.isCheck === 'Y') {
        this.$emit('gotoStep', this.requestTemplateId, 'forward')
        return
      }
      const cacheData = JSON.parse(JSON.stringify(this.basicInfo))
      cacheData.mgmtRoles = [cacheData.mgmtRoles]
      const process = this.procOptions.find(item => item.procDefId === cacheData.procDefId)
      if (process) {
        cacheData.packageName = process.rootEntity.packageName
        cacheData.entityName = process.rootEntity.name
        cacheData.procDefKey = process.procDefKey
        cacheData.procDefName = process.procDefName
        cacheData.procDefVersion = process.version
        cacheData.OperatorObjType = process.rootEntity.displayName
      } else {
        cacheData.packageName = ''
        cacheData.entityName = ''
        cacheData.procDefKey = ''
        cacheData.procDefName = ''
        cacheData.procDefVersion = ''
        cacheData.OperatorObjType = ''
      }
      const method = this.basicInfo.id === '' ? createTemp : updateTemp
      const { statusCode, data } = await method(cacheData)
      if (statusCode === 'OK') {
        if (isGoToNext) {
          this.$emit('gotoStep', data.id, 'forward')
        } else {
          this.$Notice.success({
            title: this.$t('successful'),
            desc: this.$t('successful')
          })
          this.loadPage(data.id)
        }
      }
    },
    async getTemplateData () {
      if (!!this.requestTemplateId === false) {
        this.getHandlerRoles()
        this.getPendingHandlerRoles()
        return
      }
      const params = {
        filters: [
          {
            name: 'id',
            operator: 'eq',
            value: this.requestTemplateId
          }
        ],
        paging: false
      }
      const { statusCode, data } = await getTemplateList(params)
      if (statusCode === 'OK') {
        if (data.contents.length === 1) {
          const templateData = data.contents[0]
          this.basicInfo = templateData
          this.basicInfo.mgmtRoles = templateData.mgmtRoles[0].id
          this.basicInfo.useRoles = templateData.useRoles.map(role => role.id)
          this.showFlow = templateData.procDefId !== ''
          this.getTags()
          this.getHandlerRoles()
          this.getPendingHandlerRoles()
        }
      }
    },
    gotoNext () {
      this.createTemp(true)
    },
    // 获取模版组信息
    async getGroupOptions () {
      const params = {
        filters: [],
        paging: false
      }
      const { statusCode, data } = await getTempGroupList(params)
      if (statusCode === 'OK') {
        this.groupOptions = data.contents
        if (this.basicInfo.id === '' && this.groupOptions.length > 0) {
          this.basicInfo.group = this.groupOptions[0].id
        }
      }
    },
    // 获取标签数据
    async getTags (val) {
      if (this.basicInfo.group === '') {
        return
      }
      const { statusCode, data } = await getTemplateTags(this.basicInfo.group)
      if (statusCode === 'OK') {
        const totalData = this.unique(this.tmpTagOptions.concat(data))
        this.tagOptions = totalData.map(d => {
          return {
            label: d,
            value: d
          }
        })
      }
    },
    unique (arr) {
      return Array.from(new Set(arr))
    },
    // 新增一个标签
    handleCreate (v) {
      this.tmpTagOptions.push(v)
      this.basicInfo.tags = v
    },
    // 属主角色
    async getManagementRoles (type) {
      const { statusCode, data } = await getManagementRoles()
      if (statusCode === 'OK') {
        this.mgmtRolesOptions = data
        if (type && this.basicInfo.id === '' && this.mgmtRolesOptions.length > 0) {
          this.basicInfo.mgmtRoles = this.mgmtRolesOptions[0].id
        }
      }
    },
    // 模版属主
    async getHandlerRoles (type = true) {
      const params = {
        params: {
          roles: this.basicInfo.mgmtRoles
        }
      }
      const { statusCode, data } = await getHandlerRoles(params)
      if (statusCode === 'OK') {
        this.handlerRolesOptions = data.map(d => {
          return {
            displayName: d,
            id: d
          }
        })
        if (type && this.basicInfo.id === '' && this.handlerRolesOptions.length > 0) {
          this.basicInfo.handler = this.handlerRolesOptions[0].id
          this.getPendingHandlerRoles(false)
        }
      }
    },
    // 处理人
    async getPendingHandlerRoles (type = true) {
      const params = {
        params: {
          roles: this.basicInfo.pendingRole
        }
      }
      const { statusCode, data } = await getHandlerRoles(params)
      if (statusCode === 'OK') {
        this.pendingHandlerOptions = data.map(d => {
          return {
            displayName: d,
            id: d
          }
        })
        if (type && this.basicInfo.id === '' && this.pendingHandlerOptions.length > 0) {
          this.basicInfo.pendingHandler = this.pendingHandlerOptions[0].id
        }
      }
    },
    async pendingRoleChange () {
      this.basicInfo.pendingHandler = ''
      await this.getPendingHandlerRoles()
      if (this.basicInfo.id === '' && this.pendingHandlerOptions.length > 0) {
        this.basicInfo.pendingHandler = this.pendingHandlerOptions[0].id
      }
    },
    // 使用角色
    async getUserRoles (type) {
      const { statusCode, data } = await getUserRoles()
      if (statusCode === 'OK') {
        this.useRolesOptions = data
        if (type && this.basicInfo.id === '' && this.useRolesOptions.length > 0) {
          this.basicInfo.useRoles.push(this.useRolesOptions[0].id)
          this.basicInfo.pendingRole = this.useRolesOptions[0].id
          this.getPendingHandlerRoles()
        }
      }
    },
    changeMgmtRole () {
      this.basicInfo.handler = ''
      this.getHandlerRoles()
      this.basicInfo.procDefId = ''
      this.getProcess()
      this.paramsChanged()
    },
    // 编排列表
    async getProcess () {
      // this.procOptions = []
      const { statusCode, data } = await getProcess(this.basicInfo.mgmtRoles)
      if (statusCode === 'OK') {
        this.procOptions = data
      }
    },
    // 改变管理编排
    showSwitchChange (val) {
      this.$Modal.confirm({
        title: this.$t('tw_middle_node'),
        content: this.$t('tw_workflow_change_tip'),
        'z-index': 1000000,
        onOk: async () => {
          this.basicInfo.procDefId = ''
          this.isParmasChanged = true
        },
        onCancel: () => {
          this.showFlow = !this.showFlow
        }
      })
    },
    // 定版节点切换响应
    pendingSwitchChange () {
      this.basicInfo.pendingRole = ''
      this.basicInfo.pendingHandler = ''
    },
    paramsChanged () {
      this.isParmasChanged = true
    }
  }
}
</script>
<style>
.ivu-input[disabled],
fieldset[disabled] .ivu-input {
  color: #757575 !important;
}
.ivu-select-input[disabled] {
  color: #757575;
  -webkit-text-fill-color: #757575;
}
</style>
<style lang="scss" scoped>
.ivu-form-item {
  margin-bottom: 16px;
}
.basci-info-right {
  margin-left: 60px;
  height: calc(100vh - 260px);
}

.basci-info-left {
  @extend .basci-info-right;
  border-right: 1px solid #dcdee2;
}

.title {
  font-size: 16px;
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

.basci-info-content {
  margin: 16px 64px;
}
.btn-footer-margin {
  margin: 0 6px;
}

.footer {
  position: fixed; /* 使用 fixed 定位，使其固定在页面底部 */
  left: 0;
  bottom: 0;
  width: 100%; /* 撑满整个页面宽度 */
  background-color: white; /* 设置背景色，可根据需求修改 */
  z-index: 10;
}

.content {
  text-align: center; /* 居中内容 */
  padding: 10px; /* 可根据需求调整内容与边框的间距 */
}
</style>
