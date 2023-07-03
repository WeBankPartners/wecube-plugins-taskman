<template>
  <div>
    <div>
      <Form :label-width="100">
        <Col :span="6">
          <FormItem :label="$t('name')">
            <Input v-model="formData.name" :disabled="$parent.isCheck === 'Y'" style="width: 90%" type="text"> </Input>
            <Icon size="10" style="color: #ed4014" type="ios-medical" />
          </FormItem>
        </Col>
        <Col :span="6">
          <FormItem :label="$t('request_time_limit')">
            <Select v-model="formData.expireDay" :disabled="$parent.isCheck === 'Y'" style="width: 90%" filterable>
              <Option v-for="item in expireDayOptions" :value="item" :key="item">{{ item }}{{ $t('day') }}</Option>
            </Select>
            <Icon size="10" style="color: #ed4014" type="ios-medical" />
          </FormItem>
        </Col>
        <Col :span="6">
          <FormItem :label="$t('description')">
            <Input v-model="formData.description" :disabled="$parent.isCheck === 'Y'" style="width: 90%" type="text">
            </Input>
          </FormItem>
        </Col>
      </Form>
    </div>
    <Divider plain>{{ $t('form_settings') }}</Divider>
    <Row>
      <Col span="6" style="border: 1px solid #dcdee2; padding: 0 16px">
        <div :style="{ height: MODALHEIGHT + 32 + 'px', overflow: 'auto' }">
          <Divider plain>{{ $t('preset') }}{{ $t('form_item') }}</Divider>
          <template v-for="item in formItemOptions">
            <div :key="item.id">
              <label>{{ item.displayName }}:</label>
              <Select
                v-model="item.seletedAttrs"
                @on-change="changeSelectedForm()"
                multiple
                filterable
                :disabled="$parent.isCheck === 'Y'"
                :key="item.id"
              >
                <Button type="success" @click="selectAll(item)" size="small" long>{{ $t('select_all') }}</Button>
                <Option v-for="attr in item.attributes" :value="attr.id" :key="attr.id">{{ attr.description }}</Option>
              </Select>
            </div>
          </template>
          <Divider plain>{{ $t('custom_form') }}</Divider>
          <draggable
            class="dragArea"
            :list="customElement"
            :group="{ name: 'people', pull: 'clone', put: false }"
            :sort="$parent.isCheck !== 'Y'"
            :clone="cloneDog"
          >
            <div class="list-group-item-" style="width: 100%" v-for="element in customElement" :key="element.id">
              <Input v-if="element.elementType === 'input'" :placeholder="$t('t_input')" />
              <Input v-if="element.elementType === 'textarea'" type="textarea" :placeholder="$t('textare')" />
              <Select v-if="element.elementType === 'select'" :placeholder="$t('select')"></Select>
              <div v-if="element.elementType === 'group'" style="width: 100%; height: 80px; border: 1px solid #5ea7f4">
                <span style="margin: 8px; color: #bbbbbb"> Item Group </span>
              </div>
            </div>
          </draggable>
        </div>
      </Col>
      <Col span="12" style="border: 1px solid #dcdee2; padding: 16px; width: 48%; margin: 0 4px">
        <div :style="{ height: MODALHEIGHT + 'px', overflow: 'auto' }">
          <template v-for="(item, itemIndex) in finalElement">
            <div :key="itemIndex" style="border: 1px solid #dcdee2; margin-bottom: 8px; padding: 8px">
              <span v-if="itemIndex !== editItemGroupNameIndex">
                {{ item.itemGroupName }}
                <Button
                  @click.stop="editItemGroupName(itemIndex)"
                  type="primary"
                  size="small"
                  :disabled="$parent.isCheck === 'Y'"
                  ghost
                  icon="md-create"
                ></Button>
              </span>
              <span v-else>
                <p>GroupName &nbsp;<Input v-model="item.itemGroupName" style="width: 300px"></Input></p>
                <p style="display: inline-block">
                  GroupId &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
                  <Input
                    :disabled="!(item.attrs.length === 0 || item.attrs[0].entity === '')"
                    v-model="item.itemGroup"
                    style="width: 300px"
                  ></Input>
                </p>
                <Button
                  @click.stop="confirmItemGroupName(itemIndex)"
                  type="primary"
                  size="small"
                  :disabled="$parent.isCheck === 'Y'"
                  ghost
                  icon="md-checkmark"
                ></Button>
              </span>
              <draggable
                class="dragArea"
                :list="item.attrs"
                :sort="$parent.isCheck !== 'Y'"
                group="people"
                @change="log(item)"
              >
                <div
                  @click="selectElement(itemIndex, eleIndex)"
                  :class="['list-group-item-', element.isActive ? 'active-zone' : '']"
                  :style="{ width: (element.width / 24) * 100 + '%' }"
                  v-for="(element, eleIndex) in item.attrs"
                  :key="element.id"
                >
                  <div>
                    <Icon v-if="element.required === 'yes'" size="8" style="color: #ed4014" type="ios-medical" />
                    {{ element.title }}:
                  </div>
                  <Input
                    v-if="element.elementType === 'input'"
                    :disabled="element.isEdit === 'no'"
                    v-model="element.defaultValue"
                    placeholder=""
                    style="width: calc(100% - 30px)"
                  />
                  <Input
                    v-if="element.elementType === 'textarea'"
                    :disabled="element.isEdit === 'no'"
                    v-model="element.defaultValue"
                    type="textarea"
                    style="width: calc(100% - 30px)"
                  />
                  <Select
                    v-if="element.elementType === 'select'"
                    :disabled="element.isEdit === 'no'"
                    v-model="element.defaultValue"
                    style="width: calc(100% - 30px)"
                  ></Select>
                  <Button
                    @click.stop="removeForm(itemIndex, eleIndex, element)"
                    type="error"
                    size="small"
                    :disabled="$parent.isCheck === 'Y'"
                    ghost
                    icon="ios-close"
                  ></Button>
                </div>
              </draggable>
            </div>
          </template>
        </div>
      </Col>
      <Col span="6" style="border: 1px solid #dcdee2">
        <div :style="{ height: MODALHEIGHT + 32 + 'px', overflow: 'auto' }">
          <Collapse v-model="openPanel">
            <Panel name="1">
              {{ $t('general_attributes') }}
              <div slot="content">
                <Form :label-width="80">
                  <FormItem :label="$t('field_name')">
                    <Input v-model="editElement.name" :disabled="$parent.isCheck === 'Y'" placeholder=""></Input>
                  </FormItem>
                  <FormItem :label="$t('display_name')">
                    <Input v-model="editElement.title" :disabled="$parent.isCheck === 'Y'" placeholder=""></Input>
                  </FormItem>
                  <FormItem :label="$t('data_type')">
                    <Select
                      v-model="editElement.elementType"
                      :disabled="$parent.isCheck === 'Y'"
                      @on-change="editElement.defaultValue = ''"
                    >
                      <Option value="input">Input</Option>
                      <Option value="select">Select</Option>
                      <Option value="textarea">Textarea</Option>
                    </Select>
                  </FormItem>
                  <FormItem
                    v-if="editElement.elementType === 'select'"
                    :label="editElement.entity === '' ? $t('data_set') : $t('data_source')"
                  >
                    <Input
                      v-model="editElement.dataOptions"
                      :disabled="$parent.isCheck === 'Y'"
                      placeholder="eg:a,b"
                    ></Input>
                  </FormItem>
                  <FormItem :label="$t('defaults')">
                    <Input
                      v-model="editElement.defaultValue"
                      :disabled="$parent.isCheck === 'Y'"
                      placeholder=""
                    ></Input>
                  </FormItem>
                  <!-- <FormItem :label="$t('tags')">
                    <Input v-model="editElement.tag" placeholder=""></Input>
                  </FormItem> -->
                  <FormItem :label="$t('display')">
                    <Select v-model="editElement.inDisplayName" :disabled="$parent.isCheck === 'Y'">
                      <Option value="yes">yes</Option>
                      <Option value="no">no</Option>
                    </Select>
                  </FormItem>
                  <FormItem :label="$t('editable')">
                    <Select v-model="editElement.isEdit" :disabled="$parent.isCheck === 'Y'">
                      <Option value="yes">yes</Option>
                      <Option value="no">no</Option>
                    </Select>
                  </FormItem>
                  <FormItem :label="$t('required')">
                    <Select v-model="editElement.required" :disabled="$parent.isCheck === 'Y'">
                      <Option value="yes">yes</Option>
                      <Option value="no">no</Option>
                    </Select>
                  </FormItem>
                  <FormItem :label="$t('width')">
                    <Select v-model="editElement.width" :disabled="$parent.isCheck === 'Y'">
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
                    <Input
                      v-model="editElement.regular"
                      :disabled="$parent.isCheck === 'Y'"
                      :placeholder="$t('only_supports_regular')"
                    ></Input>
                  </FormItem>
                </Form>
              </div>
            </Panel>
            <Panel name="3">
              {{ $t('data_item') }}
              <div slot="content">
                <Form :label-width="80">
                  <FormItem :label="$t('constraints')">
                    <Select v-model="editElement.isRefInside" :disabled="$parent.isCheck === 'Y'">
                      <Option value="yes">yes</Option>
                      <Option value="no">no</Option>
                    </Select>
                  </FormItem>
                </Form>
              </div>
            </Panel>
          </Collapse>
        </div>
      </Col>
    </Row>
    <div style="text-align: center; margin-top: 8px">
      <Button type="primary" @click="saveForm" :disabled="$parent.isCheck === 'Y'"
        >{{ $t('save') }}{{ $t('data_item') }}</Button
      >
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
      openPanel: '',
      formData: {
        id: '',
        name: '',
        expireDay: 1,
        description: '',
        items: [],
        updatedTime: ''
      },
      expireDayOptions: [1, 2, 3, 4, 5, 6, 7],
      selectedFormItem: [],
      selectFormItemSet: [], // 使用entity对输入、输出项分类
      formItemOptions: [], // 树形数据
      selectedFormItemOptions: [],

      customElement: [
        {
          id: 4,
          elementType: 'group',
          itemGroup: '',
          itemGroupName: ''
        },
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
          dataOptions: '',
          regular: '',
          inDisplayName: 'yes',
          isEdit: 'yes',
          multiple: 'N',
          selectList: [],
          isRefInside: 'no',
          required: 'no',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: '',
          refEntity: '',
          refPackageName: ''
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
          dataOptions: '',
          regular: '',
          inDisplayName: 'yes',
          isEdit: 'yes',
          multiple: 'N',
          selectList: [],
          isRefInside: 'no',
          required: 'no',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: '',
          refEntity: '',
          refPackageName: ''
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
          dataOptions: '',
          regular: '',
          inDisplayName: 'yes',
          isEdit: 'yes',
          multiple: 'N',
          selectList: [],
          isRefInside: 'no',
          required: 'no',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: '',
          refEntity: '',
          refPackageName: ''
        }
      ],
      editItemGroupNameIndex: null,
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
        inDisplayName: 'yes',
        isEdit: 'yes',
        multiple: 'N',
        selectList: [],
        isRefInside: 'no',
        required: 'no',
        isOutput: 'no',
        isView: 'yes',
        name: '',
        regular: '',
        sort: 0,
        title: '',
        width: 24,
        dataOptions: '',
        refEntity: '',
        refPackageName: ''
      },
      activeTag: {
        itemGroupIndex: -1,
        attrIndex: -1
      },
      specialId: ''
    }
  },
  // // props: ['requestTemplateId'],
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 400
    if (!!this.$parent.requestTemplateId !== false) {
      this.getSelectedForm()
      this.getInitData()
    }
  },
  methods: {
    selectAll (item) {
      item.attributes.forEach(attr => {
        if (!item.seletedAttrs.includes(attr.id)) {
          item.seletedAttrs.push(attr.id)
        }
      })
      this.changeSelectedForm()
    },
    editItemGroupName (index) {
      this.editItemGroupNameIndex = index
    },
    confirmItemGroupName (index) {
      let editGroup = this.finalElement[index]
      editGroup.attrs.forEach(attr => {
        attr.itemGroupName = editGroup.itemGroupName
        attr.itemGroup = editGroup.itemGroup
      })
      this.editItemGroupNameIndex = null
    },
    onMove (e, originalEvent) {
      // 不允许停靠
      if (e.relatedContext.element.id === 1) return false
      // // 不允许拖拽
      // if (e.draggedContext.element.id === 4) return false
      return true
    },
    log (item) {
      this.finalElement.forEach(l => {
        l.attrs.forEach(attr => {
          attr.itemGroup = l.itemGroup
          attr.itemGroupName = l.itemGroupName
          if (attr.id === this.specialId) {
            this.editElement = attr
            this.openPanel = '1'
          }
        })
      })
    },
    async getInitData () {
      const { statusCode, data } = await getRequestFormTemplateData(this.$parent.requestTemplateId)
      if (statusCode === 'OK') {
        this.formData = { ...data }
        if (data.items !== null && data.items.length > 0) {
          let itemGroupSet = new Set()
          data.items.sort(this.compare('sort'))
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
          this.finalElement.forEach(fEle => {
            let formIO = this.formItemOptions.find(f => f.packageName + ':' + f.name === fEle.itemGroup)
            if (formIO) {
              formIO.seletedAttrs = []
              formIO.seletedAttrs = fEle.attrs.map(attr => attr.attrDefId)
            }
          })
        }
      }
    },
    compare (prop) {
      return function (obj1, obj2) {
        var val1 = obj1[prop]
        var val2 = obj2[prop]
        if (!isNaN(Number(val1)) && !isNaN(Number(val2))) {
          val1 = Number(val1)
          val2 = Number(val2)
        }
        if (val1 < val2) {
          return -1
        } else if (val1 > val2) {
          return 1
        } else {
          return 0
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
      tmp.sort(this.compare('sort'))
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
          let findAttrs = this.finalElement.find(l => l.itemGroup === item.itemGroup)
          let findAttr = findAttrs.attrs.find(attr => attr.name === item.name)
          findAttr.id = item.id
        })
      }
    },
    changeSelectedForm () {
      this.selectedFormItem = []
      this.formItemOptions.forEach(f => {
        this.selectedFormItem = this.selectedFormItem.concat(f.seletedAttrs)
      })
      let remove = []
      const test1 = []
        .concat(...this.finalElement.map(l => l.attrs))
        .filter(l => l.entity !== '')
        .map(m => m.attrDefId)
      test1.forEach(t => {
        let tmp = t
        if (!this.selectedFormItem.includes(t.substring(tmp))) {
          remove.push(tmp)
        }
      })
      remove.forEach(r => {
        let findTag = this.selectedFormItemOptions.find(xItem => xItem.id === r)
        let findAttr = this.finalElement.find(l => l.itemGroup === findTag.entityPackage + ':' + findTag.entityName)
          .attrs
        const findIndex = findAttr.findIndex(l => l.attrDefId === r)
        findAttr.splice(findIndex, 1)
      })
      this.selectedFormItem.forEach(item => {
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
          inDisplayName: 'yes',
          isEdit: 'yes',
          multiple: seleted.multiple,
          selectList: [],
          isRefInside: 'no',
          required: 'no',
          isOutput: 'no',
          isView: 'yes',
          name: seleted.name,
          regular: '',
          sort: 0,
          title: seleted.description,
          width: 24,
          dataOptions: seleted.dataOptions,
          refEntity: seleted.refEntityName,
          refPackageName: seleted.refPackageName
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
      this.openPanel = '1'
    },
    async getSelectedForm () {
      this.formItemOptions = []
      this.selectedFormItemOptions = []
      const { statusCode, data } = await getSelectedForm(this.$parent.requestTemplateId)
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
              attributes: [d],
              seletedAttrs: []
            })
          }
        })
        this.formItemOptions = formItemOptions
        this.selectedFormItemOptions = data
      }
    },
    removeForm (itemIndex, eleIndex, element) {
      this.finalElement[itemIndex].attrs.splice(eleIndex, 1)
      const formItemOptionIndex = this.formItemOptions.findIndex(
        fio => fio.packageName + ':' + fio.name === element.itemGroup
      )
      if (this.formItemOptions[formItemOptionIndex]) {
        const seletedAttrs = this.formItemOptions[formItemOptionIndex].seletedAttrs
        seletedAttrs.splice(eleIndex, 1)
      }
      if (element.id === this.editElement.id) {
        this.openPanel = ''
      }
    },
    cloneDog (val) {
      if (this.$parent.isCheck === 'Y') return
      if (val.elementType === 'group') {
        this.finalElement.push({
          itemGroup: 'itemGroup' + idGlobal++,
          itemGroupName: 'itemGroup' + idGlobal++,
          attrs: []
        })
        return
      }
      let newItem = JSON.parse(JSON.stringify(val))
      newItem.id = 'c_' + idGlobal++
      newItem.title = newItem.title + idGlobal
      newItem.name = newItem.name + idGlobal
      this.specialId = newItem.id
      return newItem
    },
    async next () {
      if (this.$parent.isCheck === 'Y') {
        this.$emit('formSelectNextStep')
      } else {
        await this.saveForm()
        this.$emit('formSelectNextStep')
      }
    }
  },
  components: {
    draggable
  }
}
</script>
<style>
.ivu-input[disabled],
fieldset[disabled] .ivu-input {
  color: #757575 !important;
}
.ivu-select-input[disabled] {
  color: #757575 !important;
  -webkit-text-fill-color: #757575 !important;
}
.ivu-select-disabled .ivu-select-selection {
  color: #757575 !important;
}
</style>
<style scoped lang="scss">
.active-zone {
  color: #338cf0;
}
.ivu-form-item {
  margin-bottom: 8px;
}
.list-group-item- {
  display: inline-block;
  margin: 2px 0;
}
</style>
