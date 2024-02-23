<template>
  <Drawer title="配置表单" :closable="false" :width="550" v-model="openFormConfig">
    <div>
      <Form :label-width="120">
        <FormItem :label="$t('表单类型')">
          <Button shape="circle" style="border:1px solid #dd6da6; color: #dd6da6;">编排数据项</Button>
        </FormItem>
        <FormItem :label="$t('表单名')">
          <Input v-model="group.itemGroupName" style="width: 96%;" disabled></Input>
        </FormItem>
        <FormItem :label="$t('数据类型')">
          <Input v-model="group.itemGroupName" style="width: 96%;" disabled></Input>
        </FormItem>
        <FormItem :label="$t('新增一行')">
          <Select v-model="group.itemGroupRule" style="width: 96%" filterable @on-change="paramsChanged">
            <Option v-for="item in groupRules" :value="item.value" :key="item.value">{{ item.label }}</Option>
          </Select>
          <span style="color: red">*</span>
        </FormItem>
        <Divider>预制表单字段</Divider>
        <Divider>自定义分析字段</Divider>
        <Row>
          <Col span="10">字段名称</Col>
          <Col span="10">数据查找规则</Col>
        </Row>
        <Row>
          <span v-for="(item, itemIndex) in group.customItems" :key="itemIndex">
            <Col span="10">
              <Input v-model="item.name"></Input>
            </Col>
            <Col span="10">
              3333
            </Col>
            <Col span="4">
              <Button type="error" @click="deleteCustomItem(itemIndex)" icon="md-trash"></Button>
            </Col>
          </span>
          <Button style="margin-top:4px" type="primary" size="small" ghost @click="addCustomItem">新增</Button>
        </Row>
      </Form>
    </div>
    <div class="demo-drawer-footer">
      <Button type="primary" style="margin-right: 16px" :disabled="isSaveBtnActiv()" @click="saveGroupDrawer">{{
        $t('save')
      }}</Button>
      <Button @click="cancelGroupDrawer">{{ $t('cancel') }}</Button>
    </div>
  </Drawer>
</template>

<script>
import { saveRequestGroupForm, getRequestGroupForm } from '@/api/server.js'
export default {
  name: 'workflow',
  data () {
    return {
      isParmasChanged: false,
      openFormConfig: false, // 配置表单控制
      formTemplateId: '', // 数据表单id
      groupRules: [
        // 新增一行选项
        { label: '输入新数据', value: 'new' },
        { label: '选择已有数据', value: 'exist' }
      ],
      defaultItem: {
        // 自定义分析字段信息
        active: false,
        id: '',
        name: 'default',
        description: '',
        itemGroup: '',
        itemGroupType: '',
        itemGroupName: '',
        ItemGroupSort: 1,
        itemGroupRule: 'new',
        formTemplate: '',
        defaultValue: '',
        sort: 3,
        packageName: '',
        entity: '',
        attrDefId: '',
        attrDefName: '',
        attrDefDataType: '',
        elementType: 'textarea',
        title: 'Textarea16',
        width: 24,
        refPackageName: '',
        refEntity: '',
        dataOptions: '',
        required: 'no',
        regular: '',
        isEdit: 'yes',
        isView: 'yes',
        isOutput: 'no',
        inDisplayName: 'yes',
        isRefInside: 'no',
        multiple: 'N',
        defaultClear: 'no',
        selectList: null
      },
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
  props: ['requestTemplateId'],
  methods: {
    loadPage (params, formTemplateId) {
      if (params.isAdd) {
        this.group.requestTemplateId = params.requestTemplateId
        this.group.formTemplateId = params.formTemplateId
        this.group.itemGroupName = params.itemGroup
        this.group.itemGroupType = params.itemGroupType
        this.group.itemGroup = params.itemGroup
        this.group.itemGroupName = params.itemGroup
        this.group.itemGroupSort = params.itemGroupSort
        this.openFormConfig = true
      }
      this.getRequestGroupForm()
    },
    async getRequestGroupForm () {
      const { statusCode, data } = await getRequestGroupForm(
        this.requestTemplateId,
        '65d7106743737270',
        this.group.itemGroupName
      )
      if (statusCode === 'OK') {
        console.log(22, data)
      }
    },
    paramsChanged () {
      this.isParmasChanged = true
    },
    // 控制保存按钮
    isSaveBtnActiv () {
      let res = false
      if (this.group.itemGroupName === '') {
        return true
      }
      return res
    },
    async saveGroupDrawer () {
      let params = {
        formTemplateId: '',
        groups: [this.group]
      }
      const { statusCode, data } = await saveRequestGroupForm(this.requestTemplateId, params)
      if (statusCode === 'OK') {
        console.log(data)
      }
    },
    cancelGroupDrawer () {},
    deleteCustomItem (itemIndex) {},
    addCustomItem () {
      let tmpItem = JSON.parse(JSON.stringify(this.defaultItem))
      this.group.customItems.push(tmpItem)
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
