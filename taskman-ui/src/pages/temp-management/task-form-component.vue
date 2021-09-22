<template>
  <div class=" ">
    <Row>
      <Form :label-width="100">
        <Col :span="4">
          <FormItem :label="$t('name')">
            <Input v-model="formData.name" style="width:90%" type="text"> </Input>
            <Icon size="10" style="color:#ed4014" type="ios-medical" />
          </FormItem>
        </Col>
        <Col :span="4">
          <FormItem :label="$t('description')">
            <Input v-model="formData.description" style="width:90%" type="text"> </Input>
          </FormItem>
        </Col>
        <Col :span="4">
          <FormItem :label="$t('mgmtRoles')">
            <Select v-model="formData.mgmtRoles" multiple filterable>
              <Option v-for="item in mgmtRolesOptions" :value="item.id" :key="item.id">{{ item.displayName }}</Option>
            </Select>
          </FormItem>
        </Col>
        <Col :span="4">
          <FormItem :label="$t('useRoles')">
            <Select v-model="formData.useRoles" multiple filterable>
              <Option v-for="item in useRolesOptions" :value="item.id" :key="item.id">{{ item.displayName }}</Option>
            </Select>
          </FormItem>
        </Col>
      </Form>
    </Row>
    <Divider plain>表单设置</Divider>
    <Row>
      <Col span="6" style="border-right: 1px solid #dcdee2;padding: 0 16px">
        <Divider plain>输入项</Divider>
        <Select v-model="selectedInputFormItem" @on-change="changeInputSelectedForm" multiple filterable>
          <OptionGroup v-for="item in formItemOptions" :label="item.description" :key="item.id">
            <Option
              v-for="attr in item.attributes"
              :value="attr.id"
              :disabled="selectedOutputFormItem.includes(attr.id)"
              :key="attr.id"
              >{{ attr.description }}</Option
            >
          </OptionGroup>
        </Select>
        <Divider plain>输出项</Divider>
        <Select v-model="selectedOutputFormItem" @on-change="changeOutputSelectedForm" multiple filterable>
          <OptionGroup v-for="item in formItemOptions" :label="item.description" :key="item.id">
            <Option
              v-for="attr in item.attributes"
              :value="attr.id"
              :disabled="selectedInputFormItem.includes(attr.id)"
              :key="attr.id"
              >{{ attr.description }}</Option
            >
          </OptionGroup>
        </Select>
        <Divider plain>自定义表单</Divider>
        <draggable
          class="dragArea list-group"
          :list="list1"
          :group="{ name: 'people', pull: 'clone', put: false }"
          :clone="cloneDog"
        >
          <div class="list-group-item" v-for="element in list1" :key="element.id">
            <Input v-if="element.elementType === 'input'" placeholder="输入框" style="width:84%" />
            <Input
              v-if="element.elementType === 'textarea'"
              type="textarea"
              placeholder="多行文本框"
              style="width:84%"
            />
            <Select v-if="element.elementType === 'select'" placeholder="选择框" style="width:84%"> </Select>
          </div>
        </draggable>
      </Col>
      <Col span="12" style="padding: 16px">
        <div :style="{ 'max-height': MODALHEIGHT + 'px', overflow: 'auto' }">
          <template v-for="(item, itemIndex) in list2">
            <div :key="item.tag" style="border: 1px solid #dcdee2;margin-bottom: 8px;padding: 8px;">
              {{ item.tag }}
              <draggable class="dragArea list-group" :list="item.attrs" group="people">
                <!-- {{item.attrs}} -->
                <div
                  @click="selectElement(itemIndex, eleIndex)"
                  class="list-group-item"
                  v-for="(element, eleIndex) in item.attrs"
                  :key="element.id"
                >
                  <div style="width:20%;display:inline-block;text-align:right;padding-right:10px">
                    {{ element.title }}:
                  </div>
                  <Input
                    v-if="element.elementType === 'input'"
                    disabled
                    v-model="element.defaultValue"
                    placeholder=""
                    style="width:66%"
                  />
                  <Input
                    v-if="element.elementType === 'textarea'"
                    disabled
                    v-model="element.defaultValue"
                    type="textarea"
                    style="width:66%"
                  />
                  <Select
                    v-if="element.elementType === 'select'"
                    disabled
                    v-model="element.defaultValue"
                    style="width:66%"
                  ></Select>
                  <Button
                    @click.stop="removeForm(element, itemIndex, eleIndex)"
                    type="primary"
                    size="small"
                    ghost
                    icon="ios-close"
                  ></Button>
                </div>
              </draggable>
            </div>
          </template>
        </div>
      </Col>
      <Col span="6" style="border-left: 1px solid #dcdee2;">
        <Collapse>
          <Panel name="1">
            通用属性
            <div slot="content">
              <Form :label-width="80">
                <FormItem label="字段名">
                  <Input v-model="editElement.name" placeholder="Enter something..."></Input>
                </FormItem>
                <FormItem label="显示名">
                  <Input v-model="editElement.title" placeholder="Enter something..."></Input>
                </FormItem>
                <FormItem label="数据类型">
                  <Select v-model="editElement.elementType" @on-change="editElement.defaultValue = ''">
                    <Option value="input">Input</Option>
                    <Option value="select">Select</Option>
                    <Option value="textarea">Textarea</Option>
                  </Select>
                </FormItem>
                <FormItem label="默认值">
                  <Input v-model="editElement.defaultValue" placeholder="Enter something..."></Input>
                </FormItem>
                <!-- <FormItem label="标签">
                  <Input v-model="editElement.tag" placeholder="Enter something..."></Input>
                </FormItem> -->
                <FormItem label="宽度">
                  <Input v-model="editElement.width" placeholder="Enter something..."></Input>
                </FormItem>
              </Form>
            </div>
          </Panel>
          <Panel name="2">
            扩展属性
            <div slot="content">
              <Form :label-width="80">
                <FormItem label="校验规则">
                  <Input v-model="editElement.regular" placeholder="仅支持正则"></Input>
                </FormItem>
              </Form>
            </div>
          </Panel>
          <Panel name="3">
            数据项
            <div slot="content">
              当前表单项没有数据项
              <Form :label-width="80">
                <FormItem label="校验规则">
                  <Input v-model="editElement.regular" placeholder="仅支持正则"></Input>
                </FormItem>
              </Form>
            </div>
          </Panel>
        </Collapse>
      </Col>
    </Row>
    <div style="text-align:center">
      <Button type="primary" @click="saveForm">保存当前表单</Button>
    </div>
  </div>
</template>

<script>
import {
  getSelectedForm,
  getManagementRoles,
  getUserRoles,
  saveTaskForm,
  getTaskFormDataByNodeId
} from '@/api/server.js'
import draggable from 'vuedraggable'
let idGlobal = 8
export default {
  name: '',
  data () {
    return {
      MODALHEIGHT: 500,
      nodeId: '',
      nodeData: null,
      formData: {
        id: '',
        nodeDefId: '',
        nodeDefName: '',
        name: '',
        description: '',
        mgmtRoles: [],
        useRoles: [],
        items: [],
        updatedTime: ''
      },
      selectedInputFormItem: [],
      selectedOutputFormItem: [],
      formItemOptions: [], // 树形数据
      selectedFormItemOptions: [],
      mgmtRolesOptions: [],
      useRolesOptions: [],

      list1: [
        {
          id: 1,
          isCustom: true,
          name: 'input',
          title: 'Input',
          elementType: 'input',
          defaultValue: '',
          tag: '',
          width: 70,
          regular: '',
          isEdit: 'yes',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: ''
        },
        {
          id: 2,
          isCustom: true,
          name: 'select',
          title: 'Select',
          elementType: 'select',
          defaultValue: '',
          tag: '',
          width: 70,
          regular: '',
          isEdit: 'yes',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: ''
        },
        {
          id: 3,
          isCustom: true,
          name: 'textarea',
          title: 'Textarea',
          elementType: 'textarea',
          defaultValue: '',
          tag: '',
          width: 70,
          regular: '',
          isEdit: 'yes',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: ''
        }
      ],
      list2: [],
      editElement: {
        isCustom: true,
        attrDefDataType: '',
        attrDefId: '',
        attrDefName: '',
        defaultValue: '',
        tag: '',
        elementType: 'input',
        id: 0,
        isEdit: 'yes',
        isOutput: 'no',
        isView: 'yes',
        name: '',
        regular: '',
        sort: 0,
        title: '',
        width: 70
      }
    }
  },
  props: ['currentNode', 'node', 'requestTemplateId'],
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 300
    this.nodeId = this.currentNode
    this.nodeData = this.node
    console.log(100)
    this.initPage()
  },
  methods: {
    initData (currentNode, node, requestTemplateId) {
      console.log(0)
      this.nodeId = currentNode
      this.nodeData = node
      this.initPage()
    },
    initPage () {
      console.log(1, this.nodeId)
      if (this.nodeData.nodeId === this.nodeId) {
        this.formData.nodeDefId = this.nodeData.nodeDefId
        this.formData.nodeDefName = this.nodeData.nodeName
        this.getSelectedForm()
        this.getManagementRoles()
        this.getUserRoles()
        this.getTaskFormDataByNodeId()
      }
    },
    async getTaskFormDataByNodeId () {
      if (!!this.requestTemplateId === false) {
        return
      }
      this.list2 = []
      const { statusCode, data } = await getTaskFormDataByNodeId(this.requestTemplateId, this.formData.nodeDefId)
      if (statusCode === 'OK') {
        console.log(data)
        this.formData = { ...data }
        if (data.items && data.items.length > 0) {
          this.selectedFormItem = data.items.filter(item => item.attrDefId !== '')
          this.selectedInputFormItem = this.selectedFormItem
            .filter(item => item.isOutput === 'no')
            .map(attr => attr.attrDefId)
          this.selectedOutputFormItem = this.selectedFormItem
            .filter(item => item.isOutput === 'yes')
            .map(attr => attr.attrDefId)
          console.log(this.selectedFormItem)
          let customItem = data.items.filter(item => item.attrDefId === '')
          if (customItem.length > 0) {
            customItem = customItem.map(custom => {
              custom.isCustom = true
              return custom
            })
            this.list2.unshift({
              tag: 'Custom',
              attrs: customItem
            })
          }
        }
      }
    },
    async saveForm () {
      if (this.formData.name === '') {
        this.$Notice.warning({
          title: '警告',
          desc: '名称不能为空'
        })
        return
      }
      let tmp = [].concat(...JSON.parse(JSON.stringify(this.list2)).map(l => l.attrs))
      tmp.forEach((l, index) => {
        l.sort = index
        if (!isNaN(l.id) || l.id.startsWith('c_')) {
          l.id = ''
        }
      })
      let res = {
        ...this.formData,
        items: tmp
      }
      const { statusCode, data } = await saveTaskForm(this.requestTemplateId, res)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.formData = { ...data }
        data.items.forEach(item => {
          let findAttrs = this.list2.find(l => l.tag === item.tag)
          let findAttr = findAttrs.attrs.find(attr => attr.name === item.name)
          findAttr.id = item.id
        })
      }
    },
    changeInputSelectedForm () {
      let remove = []
      const test1 = []
        .concat(...this.list2.map(l => l.attrs))
        .filter(l => l.isCustom === false)
        .map(m => m.id)
      const allSelectedFormItem = this.selectedInputFormItem.concat(this.selectedOutputFormItem)
      test1.forEach(t => {
        let tmp = t.substring(2)
        if (!allSelectedFormItem.includes(tmp)) {
          remove.push(tmp)
        }
      })
      remove.forEach(r => {
        let findTag = this.selectedFormItemOptions.find(xItem => xItem.id === r)
        let findAttr = this.list2.find(l => l.tag === findTag.entityPackage + '.' + findTag.entityName).attrs
        const findIndex = findAttr.findIndex(l => l.id === r)
        findAttr.splice(findIndex, 1)
      })

      this.selectedInputFormItem.forEach(item => {
        const seleted = this.selectedFormItemOptions.find(xItem => xItem.id === item)
        let tag = seleted.entityPackage + '.' + seleted.entityName
        const attr = {
          attrDefDataType: seleted.dataType,
          attrDefId: seleted.id,
          attrDefName: seleted.name,
          defaultValue: '',
          tag: tag,
          elementType: seleted.dataType === 'str' ? 'input' : '',
          id: 'c_' + seleted.id,
          isCustom: false,
          isEdit: 'yes',
          isOutput: 'no',
          isView: 'yes',
          name: seleted.name,
          regular: '',
          sort: 0,
          title: seleted.description,
          width: 70,
          entityName: seleted.entityName,
          entityPackage: seleted.entityPackage
        }
        console.log(attr)
        const tagExist = this.list2.find(l => l.tag === tag)
        if (tagExist) {
          const find = tagExist.attrs.find(attr => attr.id.substring(2) === item)
          if (!find) {
            tagExist.attrs.push(attr)
          }
        } else {
          this.list2.push({
            tag: tag,
            attrs: [attr]
          })
        }
      })
      console.log(this.list2)
    },
    changeOutputSelectedForm () {
      let remove = []
      const test1 = []
        .concat(...this.list2.map(l => l.attrs))
        .filter(l => l.isCustom === false)
        .map(m => m.id)
      const allSelectedFormItem = this.selectedInputFormItem.concat(this.selectedOutputFormItem)
      test1.forEach(t => {
        let tmp = t.substring(2)
        if (!allSelectedFormItem.includes(tmp)) {
          remove.push(tmp)
        }
      })
      remove.forEach(r => {
        let findTag = this.selectedFormItemOptions.find(xItem => xItem.id === r)
        let findAttr = this.list2.find(l => l.tag === findTag.entityPackage + '.' + findTag.entityName).attrs
        const findIndex = findAttr.findIndex(l => l.id === r)
        findAttr.splice(findIndex, 1)
      })
      this.selectedOutputFormItem.forEach(item => {
        const seleted = this.selectedFormItemOptions.find(xItem => xItem.id === item)
        let tag = seleted.entityPackage + '.' + seleted.entityName
        const attr = {
          attrDefDataType: seleted.dataType,
          attrDefId: seleted.id,
          attrDefName: seleted.name,
          defaultValue: '',
          tag: tag,
          elementType: seleted.dataType === 'str' ? 'input' : '',
          id: 'c_' + seleted.id,
          isCustom: false,
          isEdit: 'yes',
          isOutput: 'yes',
          isView: 'yes',
          name: seleted.name,
          regular: '',
          sort: 0,
          title: seleted.description,
          width: 70,
          entityName: seleted.entityName,
          entityPackage: seleted.entityPackage
        }
        console.log(attr)
        const tagExist = this.list2.find(l => l.tag === tag)
        if (tagExist) {
          const find = tagExist.attrs.find(attr => attr.id.substring(2) === item)
          if (!find) {
            tagExist.attrs.push(attr)
          }
        } else {
          this.list2.push({
            tag: tag,
            attrs: [attr]
          })
        }
      })
      console.log(this.list2)
    },
    cloneDog (val) {
      let newItem = JSON.parse(JSON.stringify(val))
      newItem.id = idGlobal++
      newItem.tag = 'Custom'
      newItem.title = newItem.title + idGlobal
      const find = this.list2.find(l => l.tag === 'Custom')
      if (find) {
        find.attrs.push(newItem)
      } else {
        this.list2.push({
          tag: 'Custom',
          attrs: [newItem]
        })
      }
    },
    selectElement (itemIndex, eleIndex) {
      this.editElement = this.list2[itemIndex].attrs[eleIndex]
    },
    removeForm (element, itemIndex, eleIndex) {
      this.list2[itemIndex].attrs.splice(eleIndex, 1)
      const outputIndex = this.selectedOutputFormItem.findIndex(i => i === element.id.substring(2))
      if (outputIndex > -1) {
        this.selectedOutputFormItem.splice(outputIndex, 1)
        return
      }
      const inputIndex = this.selectedInputFormItem.findIndex(i => i === element.id.substring(2))
      if (inputIndex > -1) {
        this.selectedInputFormItem.splice(inputIndex, 1)
      }
    },
    async getSelectedForm () {
      const { statusCode, data } = await getSelectedForm(this.requestTemplateId)
      if (statusCode === 'OK') {
        let entitySet = new Set()
        let formItemOptions = []
        data.forEach(d => {
          const tag = d.entityPackage + '.' + d.entityName
          if (entitySet.has(tag)) {
            let find = formItemOptions.find(f => f.packageName + '.' + f.name === tag)
            find.attributes.push(d)
          } else {
            entitySet.add(tag)
            formItemOptions.push({
              description: d.entityDisplayName,
              displayName: d.entityDisplayName,
              id: d.entityId,
              name: d.entityName,
              packageName: d.entityPackage,
              attributes: [d]
            })
          }
        })
        this.formItemOptions = formItemOptions
        this.selectedFormItemOptions = data
      }
    },
    async getManagementRoles () {
      const { statusCode, data } = await getManagementRoles()
      if (statusCode === 'OK') {
        this.mgmtRolesOptions = data
      }
    },
    async getUserRoles () {
      const { statusCode, data } = await getUserRoles()
      if (statusCode === 'OK') {
        this.useRolesOptions = data
      }
    }
  },
  components: {
    draggable
  }
}
</script>

<style scoped lang="scss">
.ivu-form-item {
  margin-bottom: 8px;
}
.list-group-item {
  width: 90%;
  margin: 8px 0;
}
</style>
