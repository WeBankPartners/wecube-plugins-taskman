<template>
  <Drawer title="配置表单" :closable="false" :mask-closable="false" :width="550" v-model="openFormConfig">
    <div class="wf-entity-config">
      <Form :label-width="120">
        <FormItem :label="$t('表单类型')">
          <Button shape="circle" :style="groupStyle[group.itemGroupType]">{{
            group.itemGroupType === 'workflow' ? '编排数据项' : '自选数据项'
          }}</Button>
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
        <Divider style="margin-top:40px">预制表单字段</Divider>
        <div>
          <Row>
            <Col span="12" v-for="system in group.systemItems" :key="system.id">
              <Checkbox v-model="system.active" @on-change="paramsChanged">{{
                system.description || system.name
              }}</Checkbox>
            </Col>
          </Row>
        </div>
        <Divider style="margin-top:40px">自定义分析字段</Divider>
        <Row>
          <Col span="1">&nbsp;&nbsp;&nbsp;&nbsp;</Col>
          <Col span="6">编码</Col>
          <Col span="16">数据查找规则</Col>
        </Row>
        <Row v-for="(item, itemIndex) in group.customItems" :key="itemIndex" style="margin:6px 0">
          <Col span="1">
            <div style="margin-top: 6px;">
              <Checkbox v-model="item.active" @on-change="paramsChanged"></Checkbox>
            </div>
          </Col>
          <Col span="5">
            <Input v-model="item.name" :disabled="!isCustomItemEditable"></Input>
          </Col>
          <Col span="16">
            <ItemFilterRulesGroup
              :isBatch="false"
              ref="filterRulesGroupRef"
              :currentIndex="itemIndex"
              @filterRuleChanged="singleFilterRuleChanged"
              :disabled="!isCustomItemEditable"
              :routineExpression="item.routineExpression || group.itemGroup"
              :allEntityType="allEntityType"
              :currentSelectedEntity="group.itemGroup"
            >
            </ItemFilterRulesGroup>
          </Col>
          <Col span="2">
            <Button
              type="error"
              :disabled="!isCustomItemEditable"
              @click="deleteCustomItem(itemIndex)"
              size="small"
              icon="md-trash"
            ></Button>
          </Col>
        </Row>
        <div style="text-align: right;margin-right: 18px;">
          <Button
            type="primary"
            ghost
            :disabled="!isCustomItemEditable"
            @click="addCustomItem"
            size="small"
            icon="md-add"
          ></Button>
        </div>
      </Form>
    </div>
    <div class="demo-drawer-footer">
      <Button type="primary" style="margin-right: 16px" :disabled="isSaveBtnActive()" @click="saveGroupDrawer">{{
        $t('save')
      }}</Button>
      <Button @click="cancelGroupDrawer">{{ $t('cancel') }}</Button>
    </div>
  </Drawer>
</template>

<script>
import ItemFilterRulesGroup from './item-filter-rules-group.vue'
import { saveRequestGroupForm, getRequestGroupForm, getAllDataModels } from '@/api/server.js'
export default {
  name: 'workflow',
  data () {
    return {
      allEntityType: [], // 所有模型

      isParmasChanged: false,
      openFormConfig: false, // 配置表单控制
      formTemplateId: '', // 数据表单id
      groupRules: [
        // 新增一行选项
        { label: '输入新数据', value: 'new' },
        { label: '选择已有数据', value: 'exist' }
      ],
      groupStyle: {
        custom: {
          border: '1px solid #b886f8',
          color: '#b886f8'
        },
        workflow: {
          border: '1px solid #cba43f',
          color: '#cba43f'
        },
        optional: {
          border: '1px solid #81b337',
          color: '#81b337'
        }
      },
      defaultItem: {
        // 自定义分析字段信息
        active: true,
        id: '',
        name: '',
        description: '',
        itemGroup: '',
        itemGroupType: '',
        itemGroupName: '',
        itemGroupSort: 1,
        itemGroupRule: 'new',
        formTemplate: '',
        defaultValue: '',
        routineExpression: '',
        sort: 3,
        packageName: '',
        entity: '',
        attrDefId: '',
        attrDefName: '',
        attrDefDataType: '',
        elementType: 'calculate',
        title: '',
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
  props: ['module', 'isCustomItemEditable'],
  methods: {
    // 自定义字段获取所有类型
    async getAllDataModels () {
      let { data, status } = await getAllDataModels()
      if (status === 'OK') {
        this.allEntityType = data
      }
    },
    // 定位规则回传
    singleFilterRuleChanged (val, index) {
      this.group.customItems[index].routineExpression = val
    },
    async loadPage (params) {
      await this.getRequestGroupForm(params)
      this.getAllDataModels()
      if (params.isAdd) {
        this.group.itemGroupSort = params.itemGroupSort
        this.group.itemGroupRule = 'new'
      }
      this.openFormConfig = true
    },
    async getRequestGroupForm (params) {
      console.log(params)
      const { statusCode, data } = await getRequestGroupForm({
        taskTemplateId: params.taskTemplateId,
        requestTemplateId: params.requestTemplateId,
        entity: params.itemGroup,
        formType: params.itemGroupType,
        itemGroupId: params.itemGroupId,
        module: this.module
      })
      if (statusCode === 'OK') {
        this.group = data
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
      let finalData = JSON.parse(JSON.stringify(this.group))
      finalData.systemItems = finalData.systemItems.filter(system => system.active === true)
      finalData.customItems = finalData.customItems.filter(custom => custom.active === true)
      const { statusCode } = await saveRequestGroupForm(finalData)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.openFormConfig = false
        this.$emit('reloadParentPage', finalData.itemGroupId === '')
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
    },
    deleteCustomItem (itemIndex) {
      this.group.customItems.splice(itemIndex, 1)
    },
    addCustomItem () {
      let name = `calculate${this.group.customItems.length}`
      name = this.nameCheck(name)
      let tmpItem = JSON.parse(JSON.stringify(this.defaultItem))
      tmpItem.routineExpression = this.group.itemGroup
      tmpItem.name = name
      tmpItem.title = name
      this.group.customItems.push(tmpItem)
    },
    nameCheck (name) {
      const findIndex = this.group.customItems.findIndex(item => item.name === name)
      if (findIndex > -1) {
        return this.nameCheck(name + '1')
      } else {
        return name
      }
    }
  },
  components: {
    ItemFilterRulesGroup
  }
}
</script>

<style scoped lang="scss">
.wf-entity-config {
  height: calc(100vh - 150px);
  overflow: auto;
}

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
