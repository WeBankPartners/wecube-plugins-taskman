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
          <Option v-for="item in selectedFormItemOptions" :value="item.id" :key="item.id">{{
            item.description
          }}</Option>
        </Select>
        <Divider plain>输出项</Divider>
        <Select v-model="selectedOutputFormItem" @on-change="changeOutputSelectedForm" multiple filterable>
          <Option v-for="item in selectedFormItemOptions" :value="item.id" :key="item.id">{{
            item.description
          }}</Option>
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
        <draggable class="dragArea list-group" :list="list2" group="people">
          <div
            @click="selectElement(eleIndex)"
            class="list-group-item"
            v-for="(element, eleIndex) in list2"
            :key="element.id"
          >
            <div style="width:10%;display:inline-block;text-align:right;padding-right:10px">{{ element.title }}:</div>
            <Input
              v-if="element.elementType === 'input'"
              disabled
              v-model="element.defaultValue"
              placeholder=""
              style="width:86%"
            />
            <Input
              v-if="element.elementType === 'textarea'"
              disabled
              v-model="element.defaultValue"
              type="textarea"
              style="width:86%"
            />
            <Select
              v-if="element.elementType === 'select'"
              disabled
              v-model="element.defaultValue"
              style="width:86%"
            ></Select>
            <Button @click="removeForm(eleIndex)" type="primary" size="small" ghost icon="ios-close"></Button>
          </div>
        </draggable>
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
    <Button type="primary" @click="saveForm">保存当前表单</Button>
  </div>
</template>

<script>
import { getSelectedForm, getManagementRoles, getUserRoles, saveTaskForm } from '@/api/server.js'
import draggable from 'vuedraggable'
let idGlobal = 8
export default {
  name: '',
  data () {
    return {
      nodeId: '',
      nodeData: null,
      requestTemplateId: '614043ac9379fb1e',
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
  props: ['currentNode', 'node'],
  mounted () {
    this.nodeId = this.currentNode
    this.nodeData = this.node
    this.initPage()
  },
  methods: {
    initData (currentNode, node) {
      this.nodeId = currentNode
      this.nodeData = node
      this.initPage()
    },
    initPage () {
      if (this.nodeData.nodeId === this.nodeId) {
        this.formData.nodeDefId = this.nodeData.nodeDefId
        this.formData.nodeDefName = this.nodeData.nodeName
        console.log('请求数据')
        this.getSelectedForm()
        this.getManagementRoles()
        this.getUserRoles()
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
      let tmp = JSON.parse(JSON.stringify(this.list2))
      tmp.forEach((l, index) => {
        l.sort = index
        if (!isNaN(l.id)) {
          l.id = ''
        }
      })
      let res = {
        ...this.formData,
        items: tmp
      }
      const { statusCode, data } = await saveTaskForm(this.requestTemplateId, res)
      if (statusCode === 'OK') {
        this.formData = { ...data }
      }
    },
    changeInputSelectedForm () {
      this.selectedInputFormItem.forEach((item, itemIndex) => {
        if (this.selectedOutputFormItem.includes(item)) {
          this.$Notice.warning({
            title: '警告',
            desc: '输入项与输出项不能选择相同数据'
          })
          this.selectedInputFormItem.splice(itemIndex, 1)
        }
      })
      let remove = []
      const test1 = this.list2.filter(l => l.isCustom === false).map(m => m.id)
      test1.forEach(t => {
        if (!this.selectedInputFormItem.includes(t) && !this.selectedOutputFormItem.includes(t)) {
          remove.push(t)
        }
      })
      remove.forEach(r => {
        const findIndex = this.list2.findIndex(l => l.id === r)
        this.list2.splice(findIndex, 1)
      })
      this.selectedInputFormItem.forEach(item => {
        let find = this.list2.find(d => d.id === item)
        if (find) {
        } else {
          const seleted = this.selectedFormItemOptions.find(xItem => xItem.id === item)
          this.list2.push({
            attrDefDataType: seleted.dataType,
            attrDefId: seleted.id,
            attrDefName: seleted.name,
            defaultValue: '',
            elementType: seleted.dataType === 'str' ? 'input' : '',
            id: seleted.id,
            isCustom: false,
            isEdit: 'yes',
            isOutput: 'no',
            isView: 'yes',
            name: seleted.name,
            regular: '',
            sort: 0,
            title: seleted.description,
            width: 70
          })
        }
      })
    },
    changeOutputSelectedForm () {
      this.selectedOutputFormItem.forEach((item, itemIndex) => {
        if (this.selectedInputFormItem.includes(item)) {
          this.selectedOutputFormItem.splice(itemIndex, 1)
        }
      })
      let remove = []
      const test1 = this.list2.filter(l => l.isCustom === false).map(m => m.id)
      test1.forEach(t => {
        if (!this.selectedInputFormItem.includes(t) && !this.selectedOutputFormItem.includes(t)) {
          remove.push(t)
        }
      })
      remove.forEach(r => {
        const findIndex = this.list2.findIndex(l => l.id === r)
        this.list2.splice(findIndex, 1)
      })
      this.selectedOutputFormItem.forEach(item => {
        let find = this.list2.find(d => d.id === item)
        if (find) {
        } else {
          const seleted = this.selectedFormItemOptions.find(xItem => xItem.id === item)
          this.list2.push({
            attrDefDataType: seleted.dataType,
            attrDefId: seleted.id,
            attrDefName: seleted.name,
            defaultValue: '',
            elementType: seleted.dataType === 'str' ? 'input' : '',
            id: seleted.id,
            isCustom: false,
            isEdit: 'yes',
            isOutput: 'yes',
            isView: 'yes',
            name: seleted.name,
            regular: '',
            sort: 0,
            title: seleted.description,
            width: 70
          })
        }
      })
    },
    cloneDog (val) {
      let newItem = JSON.parse(JSON.stringify(val))
      newItem.id = idGlobal++
      newItem.title = newItem.title + idGlobal
      this.list2.push(newItem)
      this.selectElement(this.list2.length - 1)
    },
    selectElement (eleIndex) {
      this.editElement = this.list2[eleIndex]
    },
    async getSelectedForm () {
      const { statusCode, data } = await getSelectedForm(this.requestTemplateId)
      if (statusCode === 'OK') {
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
.list-group-item {
  width: 90%;
  margin: 16px 0;
}
</style>
