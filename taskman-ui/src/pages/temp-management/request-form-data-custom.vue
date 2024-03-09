<template>
  <Drawer title="配置表单" :closable="false" :mask-closable="false" :width="550" v-model="openFormConfig">
    <div>
      <Form :label-width="120">
        <FormItem :label="$t('表单类型')">
          <Button shape="circle" style="border:1px solid #dd6da6; color: #dd6da6;">自定义表单</Button>
        </FormItem>
        <FormItem :label="$t('表单名')">
          <Input v-model="group.itemGroupName" style="width: 96%;" @on-change="paramsChanged"></Input>
          <span style="color: red">*</span>
          <span v-if="group.itemGroupName === ''" style="color: red"
            >{{ $t('表单名') }}{{ $t('can_not_be_empty') }}</span
          >
        </FormItem>
        <FormItem :label="$t('新增一行')">
          <Select v-model="group.itemGroupRule" style="width: 96%" filterable disabled>
            <Option v-for="item in groupRules" :value="item.value" :key="item.value">{{ item.label }}</Option>
          </Select>
          <span style="color: red">*</span>
        </FormItem>
      </Form>
    </div>
    <div class="demo-drawer-footer">
      <Button type="primary" style="margin-right: 8px" :disabled="isSaveBtnActive()" @click="saveGroupDrawer">{{
        $t('save')
      }}</Button>
      <Button @click="cancelGroupDrawer">{{ $t('cancel') }}</Button>
    </div>
  </Drawer>
</template>

<script>
import { saveRequestGroupForm, getRequestGroupForm } from '@/api/server.js'
export default {
  name: 'custom',
  data () {
    return {
      isParmasChanged: false,
      openFormConfig: false, // 配置表单控制
      groupRules: [
        // 新增一行选项
        { label: '输入新数据', value: 'new' },
        { label: '选择已有数据', value: 'exist' }
      ],
      group: {
        requestTemplateId: '',
        formTemplateId: '',
        itemGroup: '', // 组id
        itemGroupType: '', // 组类型
        itemGroupName: '', // 组名称
        itemGroupSort: -1, // 组顺序
        itemGroupRule: 'exist', // 新增一行
        systemItems: [], // 预制表单字段
        customItems: [] // 自定义分析字段
      }
    }
  },
  props: ['requestTemplateId', 'module'],
  methods: {
    async loadPage (params) {
      if (params.isAdd) {
        this.group = {}
        this.group.requestTemplateId = params.requestTemplateId
        this.group.formTemplateId = params.formTemplateId
        this.group.itemGroupName = params.itemGroup
        this.group.itemGroupType = params.itemGroupType
        this.group.itemGroup = params.itemGroup
        this.group.itemGroupSort = params.itemGroupSort
        this.group.itemGroupRule = 'new'
        this.group.systemItems = []
        this.group.customItems = []
        this.openFormConfig = true
      } else {
        const { statusCode, data } = await getRequestGroupForm({
          formTemplateId: params.formTemplateId,
          requestTemplateId: params.requestTemplateId,
          entity: params.itemGroup,
          formType: params.itemGroupType,
          itemGroupId: params.itemGroupId,
          module: this.module
        })
        if (statusCode === 'OK') {
          this.group = data
          this.group.requestTemplateId = params.requestTemplateId
          this.openFormConfig = true
        }
      }
    },
    paramsChanged () {
      this.isParmasChanged = true
    },
    // 控制保存按钮
    isSaveBtnActive () {
      let res = false
      if (this.group.itemGroupName === '') {
        return true
      }
      return res
    },
    async saveGroupDrawer () {
      const { statusCode } = await saveRequestGroupForm(this.group)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.openFormConfig = false
        this.$emit('reloadParentPage', !this.group.itemGroupId)
      }
    },
    cancelGroupDrawer () {
      if (this.isParmasChanged) {
        this.$Modal.confirm({
          title: `${this.$t('confirm_discarding_changes')}`,
          content: `${this.group.itemGroupName}:${this.$t('params_edit_confirm')}`,
          'z-index': 1000000,
          okText: this.$t('save'),
          cancelText: this.$t('abandon'),
          onOk: async () => {
            this.saveGroupDrawer()
          },
          onCancel: () => {
            this.openFormConfig = false
          }
        })
      } else {
        this.openFormConfig = false
      }
    }
  }
}
</script>

<style scoped lang="scss">
.demo-drawer-footer {
  width: 100%;
  position: absolute;
  bottom: 0;
  left: 0;
  border-top: 1px solid #e8e8e8;
  padding: 10px 16px;
  text-align: right;
  background: #fff;
  text-align: left;
}
</style>
