<template>
  <div>
    <div class="title">
      <div class="title-text">
        {{ $t('tw_node_configuration') }}
        <span class="underline"></span>
      </div>
    </div>
    <div>
      <Form ref="formInline" inline :label-width="100">
        <FormItem :label="$t('name')">
          <Input
            type="text"
            maxlength="30"
            show-word-limit
            v-model="activeApprovalNode.name"
            @on-change="paramsChanged"
            style="width: 160px;"
          >
          </Input>
          <span style="color: red">*</span>
          <div v-if="activeApprovalNode.name === ''" style="color: red">
            {{ $t('name') }}{{ $t('can_not_be_empty') }}
          </div>
        </FormItem>
        <FormItem :label="$t('tw_validity_period')">
          <Select v-model="activeApprovalNode.expireDay" @on-change="paramsChanged" style="width: 160px;">
            <Option v-for="item in expireDayOptions" :value="item" :key="item">{{ item }}{{ $t('day') }}</Option>
          </Select>
          <span style="color: red">*</span>
        </FormItem>
        <FormItem :label="$t('description')">
          <Input
            v-model="activeApprovalNode.description"
            style="width: 200%"
            type="textarea"
            :rows="2"
            @on-change="paramsChanged"
          >
          </Input>
        </FormItem>
      </Form>
      <Form ref="formInline" inline :label-width="100">
        <FormItem :label="$t('tw_allocation')">
          <Select v-model="activeApprovalNode.handleMode" @on-change="changeRoleType" style="width: 160px;">
            <Option v-for="item in roleTypeOptions" :value="item.value" :key="item.value">{{ item.label }}</Option>
          </Select>
          <span style="color: red">*</span>
        </FormItem>
        <FormItem
          v-if="['custom', 'any', 'all'].includes(activeApprovalNode.handleMode)"
          :label="$t('handler')"
          style="width:70%"
        >
          <Row>
            <Col class="cutom-table-border" span="2">{{ $t('index') }}</Col>
            <Col class="cutom-table-border margin-left--1" span="5">{{ $t('tw_role_based_config') }}</Col>
            <Col class="cutom-table-border margin-left--1" span="5">{{ $t('tw_user_based_config') }}</Col>
            <Col class="cutom-table-border margin-left--1" span="5">{{ $t('manageRole') }}</Col>
            <Col class="cutom-table-border margin-left--1" span="4">{{ $t('tw_users') }}</Col>
            <Col class="cutom-table-border margin-left--1" span="2">{{ $t('t_action') }}</Col>
          </Row>
          <Row v-for="(roleObj, roleObjIndex) in activeApprovalNode.handleTemplates" :key="roleObjIndex" style="">
            <Col class="cutom-table-border margin-top--1" span="2">{{ roleObjIndex + 1 }}</Col>
            <Col class="cutom-table-border margin-top--1 margin-left--1" span="5">
              <Select v-model="roleObj.assign" @on-change="paramsChanged">
                <Option v-for="item in approvalRoleTypeOptions" :value="item.value" :key="item.value">{{
                  item.label
                }}</Option>
              </Select>
            </Col>
            <Col class="cutom-table-border margin-top--1 margin-left--1" span="5">
              <Select v-model="roleObj.handlerType" @on-change="paramsChanged">
                <Option
                  v-for="item in handlerTypeOptions.filter(h => h.used.includes(roleObj.assign))"
                  :value="item.value"
                  :key="item.value"
                  >{{ item.label }}</Option
                >
              </Select>
            </Col>
            <Col class="cutom-table-border margin-top--1 margin-left--1" span="5">
              <Select
                v-model="roleObj.role"
                filterable
                @on-change="changeUser(roleObj.role, roleObjIndex, true)"
                @on-open-change="getUserRoles"
                :disabled="isRoleDisable(roleObj, roleObjIndex)"
              >
                <Option v-for="item in useRolesOptions" :value="item.id" :key="item.id">{{ item.displayName }}</Option>
              </Select>
            </Col>
            <Col class="cutom-table-border margin-top--1 margin-left--1" span="4">
              <Select
                v-model="roleObj.handler"
                filterable
                @on-change="paramsChanged"
                @on-open-change="changeUser(roleObj.role, roleObjIndex, false)"
                :disabled="isHandlerDisable(roleObj, roleObjIndex)"
              >
                <Option v-for="item in roleObj.handlerOptions" :value="item.id" :key="item.id">{{
                  item.displayName
                }}</Option>
              </Select>
            </Col>
            <Col class="cutom-table-border margin-top--1 margin-left--1" span="2">
              <Button
                :disabled="activeApprovalNode.handleTemplates.length < 2"
                @click.stop="removeRoleObjItem(roleObjIndex)"
                type="error"
                size="small"
                ghost
                icon="md-trash"
              ></Button>
            </Col>
          </Row>
          <Button
            v-if="isCheck !== 'Y' && ['any', 'all'].includes(activeApprovalNode.handleMode)"
            @click.stop="addRoleObjItem"
            type="primary"
            size="small"
            ghost
            icon="md-add"
            :disabled="isHandlerAddDisable"
          ></Button>
          <span style="color: red" v-if="isHandlerAddDisable">{{ $t('tw_duplicate_data_tip') }}</span>
        </FormItem>
      </Form>
      <div style="text-align: center;">
        <Button
          v-if="isCheck !== 'Y'"
          type="primary"
          :disabled="isSaveNodeDisable || isHandlerAddDisable"
          @click="saveNode(1)"
          >{{ $t('save') }}</Button
        >
      </div>
    </div>
  </div>
</template>

<script>
// import Vue from 'vue'
import {
  getUserRoles,
  getHandlerRoles,
  updateApprovalNode,
  getApprovalNodeById,
  deleteGroupsByNodeid
} from '@/api/server'
export default {
  name: '',
  data () {
    return {
      isParmasChanged: false, // 参数变化标志位，控制右侧panel显示逻辑
      requestTemplateId: '',
      activeApprovalNode: {
        id: '',
        sort: 1,
        requestTemplate: '',
        name: `${this.$t('tw_approval')}1`,
        expireDay: 1,
        description: '',
        handleMode: 'custom',
        handleTemplates: [
          {
            assign: 'template', // 角色设置方式：template.模板指定 custom.提交人指定
            handlerType: 'template_suggest', // 人员设置方式：template.模板指定 template_suggest.模板建议 custom.提交人指定 custom_suggest.提交人建议 system.组内系统分配 claim.组内主动认领。[template,template_suggest]只当role_type=template才有
            role: '',
            handler: '',
            handlerOptions: [] // 缓存角色下的用户，添加数据时添加，保存时清除
          }
        ]
      },
      expireDayOptions: [1, 2, 3, 4, 5, 6, 7], // 时效
      roleTypeOptions: [
        // custom.单人自定义 any.协同 all.并行 admin.提交人角色管理员 auto.自动通过
        { label: this.$t('tw_single'), value: 'custom' },
        { label: this.$t('tw_collaborative'), value: 'any' },
        { label: this.$t('tw_parallel'), value: 'all' },
        { label: this.$t('tw_roleAdmin'), value: 'admin' },
        { label: this.$t('tw_autoWith'), value: 'auto' }
      ],
      approvalSingle: {
        assign: 'template', // 角色设置方式：template.模板指定 custom.提交人指定
        handlerType: 'template_suggest', // 人员设置方式：template.模板指定 template_suggest.模板建议 custom.提交人指定 custom_suggest.提交人建议 system.组内系统分配 claim.组内主动认领。[template,template_suggest]只当role_type=template才有
        role: '',
        handler: '',
        handlerOptions: [] // 缓存角色下的用户，添加数据时添加，保存时清除
      },
      approvalRoleTypeOptions: [
        { label: this.$t('tw_template_assign'), value: 'template' },
        { label: this.$t('tw_reporter_assign'), value: 'custom' }
      ],
      handlerTypeOptions: [
        { label: this.$t('tw_template_assign'), value: 'template', used: ['template'] },
        { label: this.$t('tw_template_suggest'), value: 'template_suggest', used: ['template'] },
        { label: this.$t('tw_reporter_assign'), value: 'custom', used: ['template', 'custom'] },
        { label: this.$t('tw_reporter_suggest'), value: 'custom_suggest', used: ['template', 'custom'] },
        { label: this.$t('tw_group_assign'), value: 'system', used: ['template', 'custom'] },
        { label: this.$t('tw_group_claim'), value: 'claim', used: ['template', 'custom'] }
      ],
      useRolesOptions: [], // 使用角色
      isHandlerAddDisable: true,
      isSaveNodeDisable: true
    }
  },
  props: ['isCheck'],
  watch: {
    activeApprovalNode: {
      handler (val) {
        this.isSaveNodeDisable = this.isSaveBtnActive()
        this.isHandlerAddDisable = this.isAddRoleObjDisable()
        this.$emit('nodeStatus', this.isSaveNodeDisable || this.isHandlerAddDisable)
      },
      immediate: true,
      deep: true
    }
  },
  methods: {
    loadPage (params) {
      this.isParmasChanged = false
      this.requestTemplateId = params.requestTemplateId
      this.getNodeById(params)
      this.getUserRoles()
    },
    async getNodeById (params) {
      const { statusCode, data } = await getApprovalNodeById(this.requestTemplateId, params.id, 'approve')
      if (statusCode === 'OK') {
        this.activeApprovalNode = data
        this.$emit('setFormConfigStatus', !['auto'].includes(this.activeApprovalNode.handleMode))
        this.$set(this.activeApprovalNode, 'handleTemplates', data.handleTemplates)
        // this.activeApprovalNode.handleTemplates = data.thandleTemplates
        this.mgmtData()
      }
    },
    // 控制保存按钮
    isSaveBtnActive () {
      if (this.activeApprovalNode.name === '') {
        return true
      }
      // 前三种分配类型需要设置角色
      if (
        ['custom', 'any', 'all'].includes(this.activeApprovalNode.handleMode) &&
        this.activeApprovalNode.handleTemplates
      ) {
        if (this.activeApprovalNode.handleTemplates.length === 0) {
          return true
        } else {
          let res = false
          for (let i = 0; i < this.activeApprovalNode.handleTemplates.length; i++) {
            const item = this.activeApprovalNode.handleTemplates[i]
            // 人员设置方式 没选
            if (!item.handlerType) {
              res = true
              break
            }
            // 模版建议和模版指定需要选择角色和人员
            if (
              item.assign === 'template' &&
              ['template', 'template_suggest'].includes(item.handlerType) &&
              (!item.role || !item.handler)
            ) {
              res = true
              break
            }
            // 提交人指定/提交人建议/组内系统分配/组内主动认领 需要选择角色
            if (
              item.assign === 'template' &&
              ['custom', 'custom_suggest', 'system', 'claim'].includes(item.handlerType) &&
              !item.role
            ) {
              res = true
              break
            }
          }
          return res
        }
      }
    },
    async saveNode (type, nextNodeId) {
      // type 1自我更新 2转到目标节点 3父级页面调用保存
      this.activeApprovalNode.requestTemplate = this.requestTemplateId
      let tmpData = JSON.parse(JSON.stringify(this.activeApprovalNode))
      if (['admin', 'auto'].includes(tmpData.handleMode)) {
        delete tmpData.handleTemplates
      }
      const { statusCode } = await updateApprovalNode(tmpData)
      if (statusCode === 'OK') {
        if (![2, 3].includes(type)) {
          this.$Notice.success({
            title: this.$t('successful'),
            desc: this.$t('successful')
          })
        }
        this.isParmasChanged = false
        if (type === 1) {
          this.$emit('reloadParentPage', this.activeApprovalNode.id)
        } else if (type === 2) {
          this.$emit('reloadParentPage', nextNodeId)
        }
        if (['auto'].includes(this.activeApprovalNode.handleMode)) {
          this.removeNodeGroups()
        }
      }
    },
    // 在无需表单配置的场景下，删除节点下的组
    async removeNodeGroups () {
      deleteGroupsByNodeid(this.requestTemplateId, this.activeApprovalNode.id)
    },

    // 为父页面提供状态查询
    panalStatus () {
      const nodeStatus = this.isSaveNodeDisable || this.isHandlerAddDisable
      if (nodeStatus) {
        this.$Message.warning(this.$t('tw_node_data_incomplete'))
      }
      return nodeStatus ? 'unableToSave' : 'canSave'
    },
    mgmtData () {
      this.activeApprovalNode.handleTemplates &&
        this.activeApprovalNode.handleTemplates.forEach((roleObj, roleObjIndex) => {
          if (roleObj.role !== '') {
            this.getUserByRole(roleObj.role, roleObjIndex)
          }
        })
    },
    isAddRoleObjDisable () {
      let res = false
      if (this.activeApprovalNode.handleTemplates.length > 1) {
        res = this.hasDuplicateObjects()
      }
      return res
    },
    hasDuplicateObjects () {
      let res = false
      const tmpData = JSON.parse(JSON.stringify(this.activeApprovalNode.handleTemplates)).map((data, dIndex) => {
        return {
          assign: data.assign,
          handler: data.handler,
          handlerType: data.handlerType,
          role: data.role,
          isRoleDisable: this.isRoleDisable(data, dIndex),
          isHandlerDisable: this.isHandlerDisable(data, dIndex)
        }
      })
      const roleAndHanlerCanSelect = tmpData.filter(t => !t.isRoleDisable && !t.isHandlerDisable)
      const roleCanSelect = tmpData.filter(t => !t.isRoleDisable && t.isHandlerDisable)
      if (roleAndHanlerCanSelect.length > 1) {
        let hasDuplicate = this.isDataHasDuplicate(roleAndHanlerCanSelect)
        return hasDuplicate
      }
      if (roleCanSelect.length > 1) {
        let hasDuplicate = this.isDataHasDuplicate(roleCanSelect)
        return hasDuplicate
      }
      return res
    },
    isDataHasDuplicate (arr) {
      const tmpData = JSON.parse(JSON.stringify(arr)).map(data => {
        return {
          handler: data.handler === '' ? this.generateRandomString(10) : data.handler,
          role: data.role === '' ? this.generateRandomString(10) : data.role
        }
      })
      let objectStrings = new Set() // 使用 Set 存储对象的字符串表示
      for (let obj of tmpData) {
        let objString = JSON.stringify(obj) // 将对象转换为字符串
        if (objectStrings.has(objString)) {
          return true // 存在相同的对象
        }
        objectStrings.add(objString)
      }
      return false // 不存在相同的对象
    },
    generateRandomString (length) {
      let result = ''
      const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
      const charactersLength = characters.length
      for (let i = 0; i < length; i++) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength))
      }
      return result
    },
    // 新增一组审批人
    addRoleObjItem () {
      this.activeApprovalNode.handleTemplates.push(JSON.parse(JSON.stringify(this.approvalSingle)))
      this.getUserRoles()
    },
    removeRoleObjItem (index) {
      this.activeApprovalNode.handleTemplates.splice(index, 1)
    },
    // 使用角色
    async getUserRoles () {
      const { statusCode, data } = await getUserRoles()
      if (statusCode === 'OK') {
        this.useRolesOptions = data
      }
    },
    changeRoleType () {
      this.activeApprovalNode.handleTemplates = [
        {
          assign: 'template', // 角色设置方式：template.模板指定 custom.提交人指定
          handlerType: 'template_suggest', // 人员设置方式：template.模板指定 template_suggest.模板建议 custom.提交人指定 custom_suggest.提交人建议 system.组内系统分配 claim.组内主动认领。[template,template_suggest]只当role_type=template才有
          role: '',
          handler: '',
          handlerOptions: [] // 缓存角色下的用户，添加数据时添加，保存时清除
        }
      ]
      this.$emit('setFormConfigStatus', !['auto'].includes(this.activeApprovalNode.handleMode))
      this.paramsChanged()
    },
    changeUser (role, roleObjIndex, isClearHandler) {
      if (isClearHandler) {
        this.activeApprovalNode.handleTemplates[roleObjIndex].handler = ''
      }
      this.isParmasChanged = true
      this.getUserByRole(role, roleObjIndex)
    },
    async getUserByRole (role, roleObjIndex) {
      // 猥琐，下方赋值会使该变量丢失
      const handler = this.activeApprovalNode.handleTemplates[roleObjIndex].handler
      const params = {
        params: {
          roles: role
        }
      }
      const { statusCode, data } = await getHandlerRoles(params)
      if (statusCode === 'OK') {
        this.$set(
          this.activeApprovalNode.handleTemplates[roleObjIndex],
          'handlerOptions',
          data.map(d => {
            return {
              displayName: d,
              id: d
            }
          })
        )
        this.activeApprovalNode.handleTemplates[roleObjIndex].handler = handler
      }
    },
    paramsChanged () {
      this.isParmasChanged = true
    },
    isRoleDisable (roleObj, roleObjIndex) {
      const res = !(roleObj.assign === 'template')
      if (res) {
        this.activeApprovalNode.handleTemplates[roleObjIndex].role = ''
      }
      return res
    },
    isHandlerDisable (roleObj, roleObjIndex) {
      const res = !(roleObj.assign === 'template' && ['template_suggest', 'template'].includes(roleObj.handlerType))
      if (res) {
        this.activeApprovalNode.handleTemplates[roleObjIndex].handler = ''
      }
      return res
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
.basci-info-right {
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

.cutom-table-border {
  border: 1px solid #dcdee2;
  padding: 4px;
  text-align: center;
}
.margin-left--1 {
  margin-left: -1px;
}
.margin-top--1 {
  margin-top: -1px;
}
</style>
