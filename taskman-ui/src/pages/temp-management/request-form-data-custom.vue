<template>
  <Drawer
    :title="$t('tw_configuration_form')"
    :closable="false"
    :mask-closable="false"
    :width="550"
    v-model="openFormConfig"
  >
    <div>
      <Form :label-width="120">
        <FormItem :label="$t('tw_form_type')">
          <Button shape="circle" style="border:1px solid #dd6da6; color: #dd6da6;">{{ $t('tw_custom_form') }}</Button>
        </FormItem>
        <FormItem :label="$t('tw_form_name')">
          <Input v-model.trim="group.itemGroupName" style="width: 96%;" @on-change="paramsChanged"></Input>
          <span style="color: red">*</span>
          <span v-if="!group.itemGroupName" style="color: red"
            >{{ $t('tw_form_name') }}{{ $t('can_not_be_empty') }}</span
          >
        </FormItem>
        <FormItem :label="$t('tw_add_a_new_row')">
          <Select v-model="group.itemGroupRule" style="width: 96%" filterable disabled>
            <Option v-for="item in groupRules" :value="item.value" :key="item.value">{{ item.label }}</Option>
          </Select>
          <span style="color: red">*</span>
        </FormItem>
      </Form>
    </div>
    <div class="demo-drawer-footer">
      <Button
        v-if="isCheck !== 'Y'"
        type="primary"
        style="margin-right: 8px"
        :disabled="isSaveBtnActive"
        @click="saveGroupDrawer"
        >{{ $t('save') }}</Button
      >
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
        { label: this.$t('tw_enter_new_data'), value: 'new' }, // 输入新数据
        { label: this.$t('tw_select_data_all'), value: 'exist' }, // 选择已有数据-默认全选
        { label: this.$t('tw_select_data_empty'), value: 'exist_empty' } // 选择已有数据-默认不选
      ],
      group: {
        requestTemplateId: '',
        formTemplateId: '',
        itemGroup: '', // 组id
        itemGroupType: '', // 组类型
        itemGroupName: '', // 组名称
        itemGroupSort: -1, // 组顺序
        itemGroupRule: 'new', // 新增一行
        systemItems: [], // 预制表单字段
        customItems: [] // 自定义分析字段
      }
    }
  },
  props: ['isCheck', 'requestTemplateId', 'module'],
  computed: {
    // 控制保存按钮
    isSaveBtnActive () {
      if (!this.group.itemGroupName) {
        return true
      } else {
        return false
      }
    }
  },
  methods: {
    async loadPage (params) {
      this.isParmasChanged = false
      if (params.isAdd) {
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
      if (this.isCheck === 'Y') {
        this.openFormConfig = false
        return
      }
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
