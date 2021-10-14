<template>
  <div>
    <div style="width:40%;margin: 0 auto;">
      <Form :label-width="100">
        <FormItem :label="$t('name')">
          <Input v-model="formData.name" style="width:90%" type="text"> </Input>
          <Icon size="10" style="color:#ed4014" type="ios-medical" />
        </FormItem>
        <FormItem :label="$t('description')">
          <Input v-model="formData.description" style="width:90%" type="text"> </Input>
        </FormItem>
      </Form>
    </div>
    <Divider plain>{{ $t('form_settings') }}</Divider>
    <Row>
      <Col span="6" style="border-right: 1px solid #dcdee2;padding: 0 16px">
        <Divider plain>{{ $t('input_items') }}</Divider>
        <Select v-model="selectedFormItem" @on-change="changeSelectedForm" multiple filterable>
          <OptionGroup v-for="item in formItemOptions" :label="item.description" :key="item.id">
            <Option v-for="attr in item.attributes" :value="attr.id" :key="attr.id">{{ attr.description }}</Option>
          </OptionGroup>
        </Select>
        <Divider plain>{{ $t('custom_form') }}</Divider>
        <draggable
          class="dragArea list-group"
          :list="list1"
          :group="{ name: 'people', pull: 'clone', put: false }"
          :clone="cloneDog"
        >
          <div class="list-group-item" style="width:100%" v-for="element in list1" :key="element.id">
            <Input v-if="element.elementType === 'input'" placeholder="" style="width:84%" />
            <Input v-if="element.elementType === 'textarea'" type="textarea" placeholder="" style="width:84%" />
            <Select v-if="element.elementType === 'select'" placeholder="" style="width:84%"> </Select>
          </div>
        </draggable>
      </Col>
      <Col span="12" style="padding: 16px">
        <div :style="{ 'max-height': MODALHEIGHT + 'px', overflow: 'auto' }">
          <template v-for="(item, itemIndex) in list2">
            <div :key="item.itemGroup" style="border: 1px solid #dcdee2;margin-bottom: 8px;padding: 8px;">
              {{ item.itemGroupName }}
              <draggable class="dragArea list-group" :list="item.attrs" group="people" @change="log">
                <div
                  @click="selectElement(itemIndex, eleIndex)"
                  class="list-group-item"
                  :style="{ width: element.width === 12 ? '50%' : '100%' }"
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
                    @click.stop="removeForm(itemIndex, eleIndex)"
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
                <!-- <FormItem :label="$t('tags')">
                  <Input v-model="editElement.tag" placeholder=""></Input>
                </FormItem> -->
                <FormItem :label="$t('显示')">
                  <Select v-model="editElement.inDisplayName">
                    <Option value="yes">yes</Option>
                    <Option value="no">no</Option>
                  </Select>
                </FormItem>
                <FormItem :label="$t('width')">
                  <Select v-model="editElement.width">
                    <Option :value="12">12</Option>
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
      <Button type="primary" @click="saveForm">{{ $t('save') }}{{ $t('data_item') }}</Button>
      <Button @click="next">{{ $t('next') }}</Button>
    </div>
  </div>
</template>

<script>
import { getSelectedForm, saveRequsetForm, getRequestFormTemplateData } from '@/api/server.js'
import draggable from 'vuedraggable'
let idGlobal = 8
export default {
  name: 'form-select',
  data () {
    return {
      MODALHEIGHT: 200,
      formData: {
        id: '',
        name: '',
        description: '',
        items: [],
        updatedTime: ''
      },
      selectedFormItem: [],
      selectFormItemSet: [], // 使用entity对输入、输出项分类
      formItemOptions: [], // 树形数据
      selectedFormItemOptions: [],

      list1: [
        {
          id: 1,
          isCustom: true,
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
          isCustom: true,
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
          isCustom: true,
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
      list2: [],
      editElement: {
        isCustom: true,
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
      }
    }
  },
  // // props: ['requestTemplateId'],
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 300
    if (!!this.$parent.requestTemplateId !== false) {
      this.getSelectedForm()
      this.getInitData()
    }
  },
  methods: {
    log (log) {
      console.log(log)
      this.list2.forEach(l => {
        l.attrs.forEach(attr => {
          attr.itemGroup = l.itemGroup
          attr.itemGroupName = l.itemGroupName
        })
      })
    },
    async getInitData () {
      const { statusCode, data } = await getRequestFormTemplateData(this.$parent.requestTemplateId)
      if (statusCode === 'OK') {
        this.formData = { ...data }
        if (data.items !== null && data.items.length > 0) {
          this.selectedFormItem = data.items.filter(item => item.attrDefId !== '').map(attr => attr.attrDefId)
          let customItem = data.items.filter(item => item.attrDefId === '')
          if (customItem.length > 0) {
            customItem = customItem.map(custom => {
              custom.isCustom = true
              return custom
            })
            this.list2.unshift({
              // tag: 'Custom',
              itemGroup: 'Custom',
              itemGroupName: 'Custom',
              attrs: customItem
            })
          }
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
      const { statusCode, data } = await saveRequsetForm(this.$parent.requestTemplateId, res)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.formData = { ...data }
        data.items.forEach(item => {
          let findAttrs = this.list2.find(l => l.itemGroup === item.itemGroup)
          let findAttr = findAttrs.attrs.find(attr => attr.name === item.name)
          findAttr.id = item.id
        })
      }
    },
    changeSelectedForm () {
      let remove = []
      const test1 = []
        .concat(...this.list2.map(l => l.attrs))
        .filter(l => l.isCustom === false)
        .map(m => m.attrDefId)
      test1.forEach(t => {
        let tmp = t
        if (!this.selectedFormItem.includes(t.substring(tmp))) {
          remove.push(tmp)
        }
      })
      remove.forEach(r => {
        let findTag = this.selectedFormItemOptions.find(xItem => xItem.id === r)
        let findAttr = this.list2.find(l => l.itemGroup === findTag.entityPackage + '.' + findTag.entityName).attrs
        const findIndex = findAttr.findIndex(l => l.id === r)
        findAttr.splice(findIndex, 1)
      })
      this.selectedFormItem.forEach(item => {
        const seleted = this.selectedFormItemOptions.find(xItem => xItem.id === item)
        let itemGroup = seleted.entityPackage + '.' + seleted.entityName
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
          elementType: seleted.dataType === 'str' ? 'input' : '',
          id: 'c_' + seleted.id,
          isCustom: false,
          inDisplayName: 'no',
          isEdit: 'yes',
          isOutput: 'no',
          isView: 'yes',
          name: seleted.name,
          regular: '',
          sort: 0,
          title: seleted.description,
          width: 24
        }
        const tagExist = this.list2.find(l => l.itemGroup === itemGroup)
        if (tagExist) {
          const find = tagExist.attrs.find(attr => attr.attrDefId === item)
          if (!find) {
            tagExist.attrs.push(attr)
          }
        } else {
          this.list2.push({
            // tag: tag,
            itemGroup: itemGroup,
            itemGroupName: itemGroup,
            attrs: [attr]
          })
        }
      })
    },
    selectElement (itemIndex, eleIndex) {
      this.editElement = this.list2[itemIndex].attrs[eleIndex]
    },
    async getSelectedForm () {
      this.formItemOptions = []
      this.selectedFormItemOptions = []
      const { statusCode, data } = await getSelectedForm(this.$parent.requestTemplateId)
      if (statusCode === 'OK') {
        let entitySet = new Set()
        let formItemOptions = []
        data.forEach(d => {
          const itemGroup = d.entityPackage + '.' + d.entityName
          if (entitySet.has(itemGroup)) {
            let find = formItemOptions.find(f => f.packageName + '.' + f.name === itemGroup)
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
          }
        })
        // this.formItemOptions = formItemOptions
        // this.selectedFormItemOptions = data
      }
    },
    removeForm (itemIndex, eleIndex) {
      this.list2[itemIndex].attrs.splice(eleIndex, 1)
      this.selectedFormItem.splice(eleIndex, 1)
    },
    cloneDog (val) {
      let newItem = JSON.parse(JSON.stringify(val))
      newItem.id = idGlobal++
      // newItem.tag = 'Custom'
      newItem.itemGroup = 'Custom'
      newItem.itemGroupName = 'Custom'
      newItem.title = newItem.title + idGlobal
      const find = this.list2.find(l => l.itemGroup === 'Custom')
      if (find) {
        find.attrs.push(newItem)
      } else {
        this.list2.push({
          // tag: 'Custom',
          itemGroup: 'Custom',
          itemGroupName: 'Custom',
          attrs: [newItem]
        })
      }
    },
    next () {
      this.$emit('formSelectNextStep')
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
  display: inline-block;
  margin: 8px 0;
}
</style>
