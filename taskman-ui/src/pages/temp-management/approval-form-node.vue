<template>
  <div>
    <div class="title">
      <div class="title-text">
        节点配置
        <span class="underline"></span>
      </div>
    </div>
    <div>
      <Form ref="formInline" inline :label-width="120">
        <FormItem label="审批名称">
          <Input type="text" v-model="activeApprovalNode.name" style="width: 94%;"> </Input>
          <span style="color: red">*</span>
          <span v-if="activeApprovalNode.name === ''" style="color: red"
            >{{ $t('审批名称') }}{{ $t('can_not_be_empty') }}</span
          >
        </FormItem>
        <FormItem label="审批时效">
          <Select v-model="activeApprovalNode.expireDay" filterable style="width: 94%;">
            <Option v-for="item in expireDayOptions" :value="item" :key="item">{{ item }}{{ $t('day') }}</Option>
          </Select>
          <span style="color: red">*</span>
        </FormItem>
        <FormItem label="审批描述">
          <Input
            v-model="activeApprovalNode.description"
            style="width: 200%"
            type="textarea"
            :rows="1"
            @on-change="paramsChanged"
          >
          </Input>
        </FormItem>
      </Form>
      <Form ref="formInline" inline :label-width="120">
        <FormItem label="分配">
          <Select v-model="activeApprovalNode.roleType" @on-change="changeRoleType" filterable style="width: 94%;">
            <Option v-for="item in roleTypeOptions" :value="item.value" :key="item.value">{{ item.label }}</Option>
          </Select>
          <span style="color: red">*</span>
        </FormItem>
        <FormItem
          v-if="['custom', 'any', 'all'].includes(activeApprovalNode.roleType)"
          label="审批人"
          style="width:70%"
        >
          <Row>
            <Col span="1">序号</Col>
            <Col span="5">角色设置方式</Col>
            <Col span="5">人员设置方式</Col>
            <Col span="5">角色</Col>
            <Col span="5">人员</Col>
            <Col span="2">操作</Col>
          </Row>
          <Row v-for="(roleObj, roleObjIndex) in activeApprovalNode.roleObjs" :key="roleObjIndex" style="margin: 4px 0">
            <Col span="1">{{ roleObjIndex + 1 }}</Col>
            <Col span="5">
              <Select v-model="roleObj.roleType" filterable>
                <Option v-for="item in approvalRoleTypeOptions" :value="item.value" :key="item.value">{{
                  item.label
                }}</Option>
              </Select>
            </Col>
            <Col span="5">
              <Select v-model="roleObj.handlerType" filterable>
                <Option
                  v-for="item in handlerTypeOptions.filter(h => h.used.includes(roleObj.roleType))"
                  :value="item.value"
                  :key="item.value"
                  >{{ item.label }}</Option
                >
              </Select>
            </Col>
            <Col span="5">
              <Select
                v-model="roleObj.role"
                filterable
                @on-change="getUserByRole(roleObj.role, roleObjIndex)"
                :disabled="!(roleObj.roleType === 'template')"
              >
                <Option v-for="item in useRolesOptions" :value="item.id" :key="item.id">{{ item.displayName }}</Option>
              </Select>
            </Col>
            <Col span="5">
              <Select
                v-model="roleObj.handler"
                filterable
                :disabled="
                  !(roleObj.roleType === 'template' && ['template_suggest', 'template'].includes(roleObj.handlerType))
                "
              >
                <Option v-for="item in roleObj.handlerOptions" :value="item.id" :key="item.id">{{
                  item.displayName
                }}</Option>
              </Select>
            </Col>
            <Col span="2">
              <Button
                @click.stop="removeRoleObjItem(roleObjIndex)"
                type="error"
                size="small"
                ghost
                icon="md-trash"
              ></Button>
            </Col>
          </Row>
          <Button
            v-if="['any', 'all'].includes(activeApprovalNode.roleType)"
            @click.stop="addRoleObjItem"
            type="primary"
            size="small"
            ghost
            icon="md-add"
          ></Button>
        </FormItem>
      </Form>
      <div style="text-align: center;">
        <Button type="primary" @click="saveNode">{{ $t('save') }}</Button>
      </div>
    </div>
  </div>
</template>

<script>
import Vue from 'vue'
import { getUserRoles, getHandlerRoles, updateApprovalNode, getApprovalNodeById } from '@/api/server'
export default {
  name: 'BasicInfo',
  data () {
    return {
      isParmasChanged: false, // 参数变化标志位，控制右侧panel显示逻辑
      requestTemplateId: '',
      activeApprovalNode: {
        id: '',
        sort: 1,
        requestTemplate: '',
        name: '审批1',
        expireDay: 1,
        description: '',
        roleType: 'auto',
        roleObjs: [
          {
            roleType: 'template', // 角色设置方式：template.模板指定 custom.提交人指定
            handlerType: '', // 人员设置方式：template.模板指定 template_suggest.模板建议 custom.提交人指定 custom_suggest.提交人建议 system.组内系统分配 claim.组内主动认领。[template,template_suggest]只当role_type=template才有
            role: '',
            handler: '',
            handlerOptions: [] // 缓存角色下的用户，添加数据时添加，保存时清除
          }
        ]
      },
      expireDayOptions: [1, 2, 3, 4, 5, 6, 7], // 时效
      roleTypeOptions: [
        // custom.单人自定义 any.协同 all.并行 admin.提交人角色管理员 auto.自动通过
        { label: '单人:自定义', value: 'custom' },
        { label: '协同:任意审批人处理', value: 'any' },
        { label: '并行:全部审批人处理', value: 'all' },
        { label: '提交人角色管理员', value: 'admin' },
        { label: '自动通过', value: 'auto' }
      ],
      approvalSingle: {
        roleType: 'template', // 角色设置方式：template.模板指定 custom.提交人指定
        handlerType: '', // 人员设置方式：template.模板指定 template_suggest.模板建议 custom.提交人指定 custom_suggest.提交人建议 system.组内系统分配 claim.组内主动认领。[template,template_suggest]只当role_type=template才有
        role: '',
        handler: '',
        handlerOptions: [] // 缓存角色下的用户，添加数据时添加，保存时清除
      },
      approvalRoleTypeOptions: [
        { label: '模板指定', value: 'template' },
        { label: '提交人指定', value: 'custom' }
      ],
      handlerTypeOptions: [
        { label: '模板指定', value: 'template', used: ['template'] },
        { label: '模板建议', value: 'template_suggest', used: ['template'] },
        { label: '提交人指定', value: 'custom', used: ['template', 'custom'] },
        { label: '提交人建议', value: 'custom_suggest', used: ['template', 'custom'] },
        { label: '组内系统分配', value: 'system', used: ['template', 'custom'] },
        { label: '组内主动认领', value: 'claim', used: ['template', 'custom'] }
      ],
      useRolesOptions: [] // 使用角色
    }
  },
  methods: {
    loadPage (params) {
      console.log(2)
      this.requestTemplateId = params.requestTemplateId
      this.getNodeById(params)
      this.getUserRoles()
    },
    async getNodeById (params) {
      const { statusCode, data } = await getApprovalNodeById(this.requestTemplateId, params.id)
      if (statusCode === 'OK') {
        this.activeApprovalNode = data
        this.mgmtData()
      }
    },
    async saveNode () {
      this.activeApprovalNode.requestTemplate = this.requestTemplateId
      console.log(33, this.activeApprovalNode)
      const { statusCode } = await updateApprovalNode(this.activeApprovalNode)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.$emit('reloadParentPage')
      }
    },
    mgmtData () {
      this.activeApprovalNode.roleObjs.forEach((roleObj, roleObjIndex) => {
        if (roleObj.role !== '') {
          this.getUserByRole(roleObj.role, roleObjIndex)
        }
      })
    },
    // 新增一组审批人
    addRoleObjItem () {
      this.activeApprovalNode.roleObjs.push(JSON.parse(JSON.stringify(this.approvalSingle)))
    },
    removeRoleObjItem (index) {
      this.activeApprovalNode.roleObjs.splice(index, 1)
    },
    // 使用角色
    async getUserRoles () {
      const { statusCode, data } = await getUserRoles()
      if (statusCode === 'OK') {
        this.useRolesOptions = data
      }
    },
    changeRoleType () {
      this.activeApprovalNode.roleObjs = [
        {
          roleType: 'template', // 角色设置方式：template.模板指定 custom.提交人指定
          handlerType: '', // 人员设置方式：template.模板指定 template_suggest.模板建议 custom.提交人指定 custom_suggest.提交人建议 system.组内系统分配 claim.组内主动认领。[template,template_suggest]只当role_type=template才有
          role: '',
          handler: '',
          handlerOptions: [] // 缓存角色下的用户，添加数据时添加，保存时清除
        }
      ]
    },
    async getUserByRole (role, roleObjIndex) {
      const params = {
        params: {
          roles: role
        }
      }
      const { statusCode, data } = await getHandlerRoles(params)
      if (statusCode === 'OK') {
        Vue.set(
          this.activeApprovalNode.roleObjs[roleObjIndex],
          'handlerOptions',
          data.map(d => {
            return {
              displayName: d,
              id: d
            }
          })
        )
      }
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
  margin: 0 10px;
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
</style>
