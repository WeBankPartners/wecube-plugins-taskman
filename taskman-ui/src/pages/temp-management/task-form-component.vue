<template>
  <div class=" ">
    <Row>
      <Form :label-width="100">
        <Col :span="5">
          <FormItem :label="$t('name')">
            <Input v-model="formData.name" style="width:90%" type="text"> </Input>
            <Icon size="10" style="color:#ed4014" type="ios-medical" />
          </FormItem>
        </Col>
        <Col :span="4">
          <FormItem :label="$t('processing_role')">
            <Select v-model="formData.useRoles" filterable>
              <Option v-for="item in useRolesOptions" :value="item.id" :key="item.id">{{ item.displayName }}</Option>
            </Select>
          </FormItem>
        </Col>
        <Col :span="3">
          <FormItem :label="$t('task_time_limit')">
            <Select v-model="formData.expireDay" filterable>
              <Option v-for="item in expireDayOptions" :value="item" :key="item">{{ item }}{{ $t('day') }}</Option>
            </Select>
          </FormItem>
        </Col>
        <Col :span="5">
          <FormItem :label="$t('Entity')">
            <Select v-model="selectedEntities" multiple filterable>
              <Option v-for="item in entityOptions" :value="item.value" :key="item.value">{{ item.label }}</Option>
            </Select>
          </FormItem>
        </Col>
        <Col :span="5">
          <FormItem :label="$t('description')">
            <Input v-model="formData.description" style="width:90%" type="text"> </Input>
          </FormItem>
        </Col>
      </Form>
    </Row>
    <Divider plain>{{ $t('form_settings') }}</Divider>
    <Row>
      <Col span="6" style="border-right: 1px solid #dcdee2;padding: 0 16px">
        <Divider plain>{{ $t('preset') }}{{ $t('input_items') }}</Divider>
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
        <Divider plain>{{ $t('preset') }}{{ $t('output_items') }}</Divider>
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
        <Divider plain>{{ $t('custom_form') }}</Divider>
        <draggable
          class="dragArea"
          :list="customElement"
          :group="{ name: 'people', pull: 'clone', put: false }"
          :clone="cloneDog"
        >
          <div class="list-group-item" style="width:100%" v-for="element in customElement" :key="element.id">
            <Input v-if="element.elementType === 'input'" :placeholder="$t('input')" />
            <Input v-if="element.elementType === 'textarea'" type="textarea" :placeholder="$t('textare')" />
            <Select v-if="element.elementType === 'select'" :placeholder="$t('select')"></Select>
          </div>
        </draggable>
      </Col>
      <Col span="12" style="padding: 16px">
        <div :style="{ 'max-height': MODALHEIGHT + 'px', overflow: 'auto' }">
          <template v-for="(item, itemIndex) in finalElement">
            <div :key="item.itemGroup" style="border: 1px solid #dcdee2;margin-bottom: 8px;padding: 8px;">
              {{ item.itemGroupName }}
              <draggable class="dragArea" :list="item.attrs" group="people" @change="log">
                <div
                  @click="selectElement(itemIndex, eleIndex)"
                  :class="['list-group-item', element.isActive ? 'active-zone' : '']"
                  :style="{ width: (element.width / 24) * 100 + '%' }"
                  v-for="(element, eleIndex) in item.attrs"
                  :key="element.id"
                >
                  <div>{{ element.title }}:</div>
                  <Input
                    v-if="element.elementType === 'input'"
                    disabled
                    v-model="element.defaultValue"
                    placeholder=""
                    style="width: calc(100% - 30px);"
                  />
                  <Input
                    v-if="element.elementType === 'textarea'"
                    disabled
                    v-model="element.defaultValue"
                    type="textarea"
                    style="width: calc(100% - 30px);"
                  />
                  <Select
                    v-if="element.elementType === 'select'"
                    disabled
                    v-model="element.defaultValue"
                    style="width: calc(100% - 30px);"
                  ></Select>
                  <Button
                    @click.stop="removeForm(element, itemIndex, eleIndex)"
                    type="error"
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
            {{ $t('general_attributes') }}
            <div slot="content">
              <Form :label-width="80">
                <FormItem :label="$t('field_name')">
                  <Input v-model="editElement.name" placeholder=""></Input>
                </FormItem>
                <FormItem :label="$t('display_name')">
                  <Input v-model="editElement.title" placeholder=""></Input>
                </FormItem>
                <FormItem :label="$t('data_type')">
                  <Select v-model="editElement.elementType" @on-change="editElement.defaultValue = ''">
                    <Option value="input">Input</Option>
                    <Option value="select">Select</Option>
                    <Option value="textarea">Textarea</Option>
                  </Select>
                </FormItem>
                <FormItem :label="$t('defaults')">
                  <Input v-model="editElement.defaultValue" placeholder=""></Input>
                </FormItem>
                <!-- <FormItem label="标签">
                  <Input v-model="editElement.tag" placeholder=""></Input>
                </FormItem> -->
                <FormItem :label="$t('display')">
                  <Select v-model="editElement.inDisplayName">
                    <Option value="yes">yes</Option>
                    <Option value="no">no</Option>
                  </Select>
                </FormItem>
                <FormItem :label="$t('editable')">
                  <Select v-model="editElement.isEdit">
                    <Option value="yes">yes</Option>
                    <Option value="no">no</Option>
                  </Select>
                </FormItem>
                <FormItem :label="$t('width')">
                  <Select v-model="editElement.width">
                    <Option :value="6">6</Option>
                    <Option :value="12">12</Option>
                    <Option :value="18">18</Option>
                    <Option :value="24">24</Option>
                  </Select>
                </FormItem>
              </Form>
            </div>
          </Panel>
          <Panel name="2">
            {{ $t('extended_attributes') }}
            <div slot="content">
              <Form :label-width="80">
                <FormItem :label="$t('validation_rules')">
                  <Input v-model="editElement.regular" :placeholder="$t('only_supports_regular')"></Input>
                </FormItem>
              </Form>
            </div>
          </Panel>
          <Panel name="3">
            {{ $t('data_item') }}
            <div slot="content">
              {{ $t('no_data_item') }}
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
import { getSelectedForm, getUserRoles, saveTaskForm, getTaskFormDataByNodeId } from '@/api/server.js'
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
        useRoles: '',
        items: [],
        expireDay: 1,
        updatedTime: ''
      },
      selectedEntities: [],
      entityOptions: [],
      expireDayOptions: [1, 2, 3, 4, 5, 6, 7],
      selectedInputFormItem: [],
      selectedOutputFormItem: [],
      cacheFormItemOptions: [], // 原始数据缓存
      formItemOptions: [], // 树形数据
      selectedFormItemOptions: [],
      useRolesOptions: [],

      customElement: [
        {
          id: 1,
          name: 'input',
          title: 'Input',
          elementType: 'input',
          defaultValue: '',
          // tag: '',
          itemGroup: '',
          itemGroupName: '',
          packageName: '',
          entity: '',
          width: 24,
          regular: '',
          inDisplayName: 'no',
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
          name: 'select',
          title: 'Select',
          elementType: 'select',
          defaultValue: '',
          // tag: '',
          itemGroup: '',
          itemGroupName: '',
          packageName: '',
          entity: '',
          width: 24,
          regular: '',
          inDisplayName: 'no',
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
          name: 'textarea',
          title: 'Textarea',
          elementType: 'textarea',
          defaultValue: '',
          // tag: '',
          itemGroup: '',
          itemGroupName: '',
          packageName: '',
          entity: '',
          width: 24,
          regular: '',
          inDisplayName: 'no',
          isEdit: 'yes',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: ''
        }
      ],
      finalElement: [],
      editElement: {
        attrDefDataType: '',
        attrDefId: '',
        attrDefName: '',
        defaultValue: '',
        // tag: '',
        itemGroup: '',
        itemGroupName: '',
        packageName: '',
        entity: '',
        elementType: 'input',
        id: 0,
        inDisplayName: 'no',
        isEdit: 'yes',
        isOutput: 'no',
        isView: 'yes',
        name: '',
        regular: '',
        sort: 0,
        title: '',
        width: 24
      },
      activeTag: {
        itemGroupIndex: -1,
        attrIndex: -1
      }
    }
  },
  props: ['currentNode', 'node', 'requestTemplateId'],
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 480
    this.nodeId = this.currentNode
    this.nodeData = this.node
    this.initPage()
  },
  watch: {
    selectedEntities: function (val) {
      this.formItemOptions = this.cacheFormItemOptions.filter(c => {
        if (val.includes(c.name)) {
          return c
        }
      })
      let canSelect = []
      this.formItemOptions.forEach(f => {
        f.attributes.forEach(attr => {
          canSelect.push(attr.id)
        })
      })
      this.selectedInputFormItem.forEach((i, index) => {
        if (!canSelect.includes(i)) {
          this.selectedInputFormItem = this.selectedInputFormItem.splice(index, 1)
        }
      })
      this.selectedOutputFormItem.forEach((i, index) => {
        if (!canSelect.includes(i)) {
          this.selectedOutputFormItem = this.selectedOutputFormItem.splice(index, 1)
        }
      })
    }
  },
  methods: {
    log (log) {
      this.finalElement.forEach(l => {
        l.attrs.forEach(attr => {
          attr.itemGroup = l.itemGroup
          attr.itemGroupName = l.itemGroupName
        })
      })
    },
    initData (currentNode, node, requestTemplateId) {
      this.nodeId = currentNode
      this.nodeData = node
      this.initPage()
    },
    initPage () {
      if (this.nodeData.nodeId === this.nodeId) {
        this.formData.nodeDefId = this.nodeData.nodeDefId
        this.formData.nodeDefName = this.nodeData.nodeName
        this.getSelectedForm()
        this.getUserRoles()
        this.getTaskFormDataByNodeId()
      }
    },
    async getTaskFormDataByNodeId () {
      if (!!this.requestTemplateId === false) {
        return
      }
      this.finalElement = []
      const { statusCode, data } = await getTaskFormDataByNodeId(this.requestTemplateId, this.formData.nodeDefId)
      if (statusCode === 'OK') {
        this.formData = { ...data }
        this.formData.nodeDefId = this.nodeData.nodeDefId
        this.formData.nodeDefName = this.nodeData.nodeName
        if (data.items && data.items.length > 0) {
          this.selectedFormItem = data.items.filter(item => item.attrDefId !== '')
          this.selectedFormItem.forEach(sForm => {
            if (!this.selectedEntities.includes(sForm.entity)) {
              this.selectedEntities.push(sForm.entity)
            }
          })
          this.selectedInputFormItem = this.selectedFormItem
            .filter(item => item.isOutput === 'no')
            .map(attr => attr.attrDefId)
          this.selectedOutputFormItem = this.selectedFormItem
            .filter(item => item.isOutput === 'yes')
            .map(attr => attr.attrDefId)
        }
        if (data.items !== null && data.items.length > 0) {
          let itemGroupSet = new Set()
          data.items.forEach(item => {
            if (itemGroupSet.has(item.itemGroup)) {
              let exitEle = this.finalElement.find(ele => ele.itemGroup === item.itemGroup)
              exitEle.attrs.push(item)
            } else {
              itemGroupSet.add(item.itemGroup)
              this.finalElement.unshift({
                itemGroup: item.itemGroup,
                itemGroupName: item.itemGroupName,
                attrs: [item]
              })
            }
          })
          // this.selectedFormItem = data.items.filter(item => item.attrDefId !== '').map(attr => attr.attrDefId)
        }
      }
    },
    async saveForm () {
      if (this.formData.name === '') {
        this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('name') + ' ' + this.$t('can_not_be_empty')
        })
        return
      }
      let tmp = [].concat(...JSON.parse(JSON.stringify(this.finalElement)).map(l => l.attrs))
      tmp.forEach((l, index) => {
        l.sort = index
        if (!isNaN(l.id) || l.id.startsWith('c_')) {
          l.id = ''
        }
      })
      let cloneFormData = JSON.parse(JSON.stringify(this.formData))
      cloneFormData.useRoles = [cloneFormData.useRoles]
      let res = {
        ...cloneFormData,
        items: tmp
      }
      const { statusCode, data } = await saveTaskForm(this.requestTemplateId, res)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.formData = { ...data }
        this.formData.useRoles = data.useRoles[0]
        data.items.forEach(item => {
          let findAttrs = this.finalElement.find(l => l.itemGroup === item.itemGroup)
          let findAttr = findAttrs.attrs.find(attr => attr.name === item.name)
          findAttr.id = item.id
        })
      }
    },
    changeInputSelectedForm () {
      let remove = []
      const test1 = []
        .concat(...this.finalElement.map(l => l.attrs))
        .filter(l => l.entity !== '')
        .map(m => m.attrDefId)
      const allSelectedFormItem = this.selectedInputFormItem.concat(this.selectedOutputFormItem)
      test1.forEach(t => {
        let tmp = t
        if (!allSelectedFormItem.includes(tmp)) {
          remove.push(tmp)
        }
      })
      remove.forEach(r => {
        let findTag = this.selectedFormItemOptions.find(xItem => xItem.id === r)
        let findAttr = this.finalElement.find(l => l.itemGroup === findTag.entityPackage + ':' + findTag.entityName)
          .attrs
        const findIndex = findAttr.findIndex(l => l.id === r)
        findAttr.splice(findIndex, 1)
      })

      this.selectedInputFormItem.forEach(item => {
        const seleted = this.selectedFormItemOptions.find(xItem => xItem.id === item)
        let itemGroup = seleted.entityPackage + ':' + seleted.entityName
        const elementType = {
          str: 'input',
          ref: 'select'
        }
        const attr = {
          attrDefDataType: seleted.dataType,
          attrDefId: seleted.id,
          attrDefName: seleted.name,
          defaultValue: '',
          // tag: tag,
          itemGroup: itemGroup,
          itemGroupName: itemGroup,
          packageName: seleted.entityPackage,
          entity: seleted.entityName,
          elementType: elementType[seleted.dataType],
          id: 'c_' + seleted.id,
          isEdit: 'no',
          isOutput: 'no',
          isView: 'yes',
          name: seleted.name,
          regular: '',
          sort: 0,
          title: seleted.description,
          width: 24
        }
        const tagExist = this.finalElement.find(l => l.itemGroup === itemGroup)
        if (tagExist) {
          const find = tagExist.attrs.find(attr => attr.attrDefId === item)
          if (!find) {
            tagExist.attrs.push(attr)
          }
        } else {
          this.finalElement.push({
            // tag: tag,
            itemGroup: itemGroup,
            itemGroupName: itemGroup,
            attrs: [attr]
          })
        }
      })
    },
    changeOutputSelectedForm () {
      let remove = []
      const test1 = []
        .concat(...this.finalElement.map(l => l.attrs))
        .filter(l => l.entity !== '')
        .map(m => m.attrDefId)
      const allSelectedFormItem = this.selectedInputFormItem.concat(this.selectedOutputFormItem)
      test1.forEach(t => {
        let tmp = t
        if (!allSelectedFormItem.includes(tmp)) {
          remove.push(tmp)
        }
      })
      remove.forEach(r => {
        let findTag = this.selectedFormItemOptions.find(xItem => xItem.id === r)
        let findAttr = this.finalElement.find(l => l.itemGroup === findTag.entityPackage + ':' + findTag.entityName)
          .attrs
        const findIndex = findAttr.findIndex(l => l.id === r)
        findAttr.splice(findIndex, 1)
      })
      this.selectedOutputFormItem.forEach(item => {
        const seleted = this.selectedFormItemOptions.find(xItem => xItem.id === item)
        let itemGroup = seleted.entityPackage + ':' + seleted.entityName
        const elementType = {
          str: 'input',
          ref: 'select'
        }
        const attr = {
          attrDefDataType: seleted.dataType,
          attrDefId: seleted.id,
          attrDefName: seleted.name,
          defaultValue: '',
          // tag: tag,
          itemGroup: itemGroup,
          itemGroupName: itemGroup,
          packageName: seleted.entityPackage,
          entity: seleted.entityName,
          elementType: elementType[seleted.dataType],
          id: 'c_' + seleted.id,
          inDisplayName: 'no',
          isEdit: 'yes',
          isOutput: 'yes',
          isView: 'yes',
          name: seleted.name,
          regular: '',
          sort: 0,
          title: seleted.description,
          width: 24
        }
        const tagExist = this.finalElement.find(l => l.itemGroup === itemGroup)
        if (tagExist) {
          const find = tagExist.attrs.find(attr => attr.attrDefId === item)
          if (!find) {
            tagExist.attrs.push(attr)
          }
        } else {
          this.finalElement.push({
            // tag: tag,
            itemGroup: itemGroup,
            itemGroupName: itemGroup,
            attrs: [attr]
          })
        }
      })
    },
    cloneDog (val) {
      let newItem = JSON.parse(JSON.stringify(val))
      newItem.id = idGlobal++
      // newItem.tag = 'Custom'
      newItem.itemGroup = 'Custom'
      newItem.itemGroupName = 'Custom'
      newItem.title = newItem.title + idGlobal
      const find = this.finalElement.find(l => l.itemGroup === 'Custom')
      if (find) {
        find.attrs.push(newItem)
      } else {
        this.finalElement.push({
          tag: 'Custom',
          itemGroup: 'Custom',
          itemGroupName: 'Custom',
          attrs: [newItem]
        })
      }
    },
    selectElement (itemIndex, eleIndex) {
      if (this.activeTag.itemGroupIndex !== -1 && this.activeTag.attrIndex !== -1) {
        this.finalElement[this.activeTag.itemGroupIndex].attrs[this.activeTag.attrIndex].isActive = false
      }
      this.activeTag = {
        itemGroupIndex: itemIndex,
        attrIndex: eleIndex
      }
      this.editElement = this.finalElement[itemIndex].attrs[eleIndex]
      this.editElement.isActive = true
    },
    removeForm (element, itemIndex, eleIndex) {
      this.finalElement[itemIndex].attrs.splice(eleIndex, 1)
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
      this.entityOptions = []
      const { statusCode, data } = await getSelectedForm(this.requestTemplateId)
      if (statusCode === 'OK') {
        let entitySet = new Set()
        let formItemOptions = []
        data.forEach(d => {
          const itemGroup = d.entityPackage + ':' + d.entityName
          if (entitySet.has(itemGroup)) {
            let find = formItemOptions.find(f => f.packageName + ':' + f.name === itemGroup)
            find.attributes.push(d)
          } else {
            entitySet.add(itemGroup)
            formItemOptions.push({
              description: d.entityDisplayName,
              displayName: d.entityDisplayName,
              id: d.entityId,
              name: d.entityName,
              packageName: d.entityPackage,
              attributes: [d]
            })
            this.entityOptions.push({
              label: d.entityDisplayName,
              value: d.entityName
            })
          }
        })
        this.cacheFormItemOptions = JSON.parse(JSON.stringify(formItemOptions))
        this.formItemOptions = formItemOptions
        this.selectedFormItemOptions = data
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
.active-zone {
  color: #338cf0;
}
.ivu-form-item {
  margin-bottom: 8px;
}
.list-group-item {
  display: inline-block;
  margin: 8px 0;
}
</style>
