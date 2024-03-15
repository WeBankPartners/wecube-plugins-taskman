<template>
  <div>
    <Row>
      <Col span="5" style="border: 1px solid #dcdee2; padding: 0 16px">
        <div :style="{ height: MODALHEIGHT + 32 + 'px', overflow: 'auto' }">
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
              <Select v-if="element.elementType === 'wecmdbEntity'" placeholder="模型数据项"></Select>
              <DatePicker
                v-if="element.elementType === 'datePicker'"
                :type="element.type"
                :placeholder="$t('tw_date_picker')"
                style="width:100%"
              ></DatePicker>
              <div v-if="element.elementType === 'group'" style="width: 100%; height: 80px; border: 1px solid #5ea7f4">
                <span style="margin: 8px; color: #bbbbbb"> Item Group </span>
              </div>
            </div>
          </draggable>
        </div>
      </Col>
      <Col span="14" style="border: 1px solid #dcdee2; padding: 0 16px; width: 57%; margin: 0 4px">
        <div :style="{ height: MODALHEIGHT + 30 + 'px', overflow: 'auto' }">
          <Divider>{{ $t('tw_preview') }}</Divider>
          <div class="title">
            <div class="title-text">
              {{ $t('请求信息') }}
              <span class="underline"></span>
            </div>
          </div>
          <template v-for="(item, itemIndex) in finalElement">
            <div :key="itemIndex" style="border: 2px dashed #A2EF4D; margin: 8px 0; padding: 8px;min-height: 48px;">
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
                  <div
                    class="custom-title"
                    :style="
                      ['calculate', 'textarea'].includes(element.elementType)
                        ? 'vertical-align: top;word-break: break-all;'
                        : ''
                    "
                  >
                    <Icon v-if="element.required === 'yes'" size="8" style="color: #ed4014" type="ios-medical" />
                    {{ element.title }}
                  </div>
                  <Input
                    v-if="element.elementType === 'input'"
                    :disabled="element.isEdit === 'no'"
                    v-model="element.defaultValue"
                    placeholder=""
                    class="custom-item"
                  />
                  <Input
                    v-if="element.elementType === 'textarea'"
                    :disabled="element.isEdit === 'no'"
                    v-model="element.defaultValue"
                    type="textarea"
                    :rows="2"
                    class="custom-item"
                  />
                  <Select
                    v-if="element.elementType === 'select'"
                    :disabled="element.isEdit === 'no'"
                    v-model="element.defaultValue"
                    class="custom-item"
                  ></Select>
                  <Select
                    v-if="element.elementType === 'wecmdbEntity'"
                    :disabled="element.isEdit === 'no'"
                    v-model="element.defaultValue"
                    class="custom-item"
                  ></Select>
                  <DatePicker
                    v-if="element.elementType === 'datePicker'"
                    class="custom-item"
                    :type="element.type"
                  ></DatePicker>
                  <Button
                    @click.stop="removeForm(itemIndex, eleIndex, element)"
                    type="error"
                    size="small"
                    :disabled="$parent.isCheck === 'Y'"
                    icon="ios-trash"
                  ></Button>
                </div>
              </draggable>
            </div>
          </template>
        </div>
      </Col>
      <Col span="5" style="border: 1px solid #dcdee2">
        <div :style="{ height: MODALHEIGHT + 32 + 'px', overflow: 'auto' }">
          <Collapse v-model="openPanel">
            <Panel name="1">
              {{ $t('general_attributes') }}
              <div slot="content">
                <Form :label-width="80">
                  <FormItem :label="$t('display_name')">
                    <Input
                      v-model="editElement.title"
                      @on-change="paramsChanged"
                      :disabled="$parent.isCheck === 'Y'"
                      placeholder=""
                    ></Input>
                  </FormItem>
                  <FormItem :label="$t('tw_code')">
                    <Input
                      v-model="editElement.name"
                      @on-change="paramsChanged"
                      :disabled="$parent.isCheck === 'Y' || editElement.entity !== ''"
                      placeholder=""
                    ></Input>
                  </FormItem>
                  <FormItem :label="$t('data_type')">
                    <Select
                      v-model="editElement.elementType"
                      :disabled="true"
                      @on-change="editElement.defaultValue = ''"
                    >
                      <Option value="input">Input</Option>
                      <Option value="select">Select</Option>
                      <Option value="textarea">Textarea</Option>
                      <Option value="wecmdbEntity">模型数据项</Option>
                      <Option value="datePicker">DatePicker</Option>
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
                      @on-change="paramsChanged"
                    ></Input>
                  </FormItem>
                  <!--添加wecmdbEntity类型，根据选择配置生成url(用于获取下拉配置)-->
                  <FormItem v-if="editElement.elementType === 'wecmdbEntity'" :label="$t('data_source')">
                    <Select
                      v-model="editElement.dataOptions"
                      filterable
                      @on-change="paramsChanged"
                      :disabled="$parent.isCheck === 'Y'"
                    >
                      <Option v-for="i in allEntityList" :value="i" :key="i">{{ i }}</Option>
                    </Select>
                  </FormItem>
                  <!-- <FormItem :label="$t('tags')">
                    <Input v-model="editElement.tag" placeholder=""></Input>
                  </FormItem> -->
                  <FormItem :label="$t('display')">
                    <RadioGroup v-model="editElement.inDisplayName" @on-change="paramsChanged">
                      <Radio label="yes" :disabled="$parent.isCheck === 'Y'">{{ $t('tw_yes') }}</Radio>
                      <Radio label="no" :disabled="$parent.isCheck === 'Y'">{{ $t('tw_no') }}</Radio>
                    </RadioGroup>
                  </FormItem>
                  <FormItem :label="$t('editable')">
                    <RadioGroup v-model="editElement.isEdit" @on-change="paramsChanged">
                      <Radio label="yes" :disabled="$parent.isCheck === 'Y'">{{ $t('tw_yes') }}</Radio>
                      <Radio label="no" :disabled="$parent.isCheck === 'Y'">{{ $t('tw_no') }}</Radio>
                    </RadioGroup>
                  </FormItem>
                  <FormItem :label="$t('required')">
                    <RadioGroup v-model="editElement.required" @on-change="paramsChanged">
                      <Radio label="yes" :disabled="$parent.isCheck === 'Y'">{{ $t('tw_yes') }}</Radio>
                      <Radio label="no" :disabled="$parent.isCheck === 'Y'">{{ $t('tw_no') }}</Radio>
                    </RadioGroup>
                  </FormItem>
                  <FormItem :label="$t('tw_default_empty')">
                    <RadioGroup v-model="editElement.defaultClear" @on-change="paramsChanged">
                      <Radio label="yes" :disabled="$parent.isCheck === 'Y'">{{ $t('tw_yes') }}</Radio>
                      <Radio label="no" :disabled="$parent.isCheck === 'Y'">{{ $t('tw_no') }}</Radio>
                    </RadioGroup>
                  </FormItem>
                  <FormItem :label="$t('defaults')">
                    <Input
                      v-model="editElement.defaultValue"
                      :disabled="$parent.isCheck === 'Y' || editElement.defaultClear === 'yes'"
                      placeholder=""
                      @on-change="paramsChanged"
                    ></Input>
                  </FormItem>
                  <FormItem :label="$t('width')">
                    <Select v-model="editElement.width" @on-change="paramsChanged" :disabled="$parent.isCheck === 'Y'">
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
                      @on-change="paramsChanged"
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
                    <Select
                      v-model="editElement.isRefInside"
                      @on-change="paramsChanged"
                      :disabled="$parent.isCheck === 'Y'"
                    >
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
    <div style="text-align: center;margin-top: 16px;">
      <Button @click="gotoForward" ghost type="primary" class="btn-footer-margin">{{ $t('forward') }}</Button>
      <Button
        v-if="isCheck !== 'Y'"
        @click="saveMsgForm(1)"
        type="info"
        :disabled="isSaveBtnActive()"
        class="btn-footer-margin"
        >{{ $t('save') }}</Button
      >
      <Button @click="gotoNext" type="primary" class="btn-footer-margin">{{ $t('next') }}</Button>
    </div>
  </div>
</template>

<script>
import { saveRequsetForm, getRequestFormTemplateData, getAllDataModels } from '@/api/server.js'
import draggable from 'vuedraggable'
export default {
  name: 'form-select',
  data () {
    return {
      requestTemplateId: '',
      isParmasChanged: false, // 参数变化标志位
      MODALHEIGHT: 200,
      customElement: [
        // 预制自定义表单项目
        {
          id: 1,
          name: 'input',
          title: 'Input',
          elementType: 'input',
          defaultValue: '',
          defaultClear: 'no',
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
          defaultClear: 'no',
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
          defaultClear: 'no',
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
          id: 5,
          name: 'wecmdbEntity',
          title: 'WecmdbEntity',
          elementType: 'wecmdbEntity',
          defaultValue: '',
          defaultClear: 'no',
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
          id: 6,
          name: 'datePicker',
          title: 'datePicker',
          elementType: 'datePicker',
          type: 'datetime',
          defaultValue: '',
          defaultClear: 'no',
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
      finalElement: [
        {
          itemGroup: 'requestInfo',
          itemGroupName: this.$t('tw_request_title'),
          attrs: []
        }
      ],
      openPanel: '', // 正在编辑的表单项
      editElement: {
        // 正在编辑的表单项属性
        attrDefDataType: '',
        attrDefId: '',
        attrDefName: '',
        defaultValue: '',
        defaultClear: 'no',
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
      allEntityList: [],
      formId: '' // 缓存表单id，供编辑使用
    }
  },
  props: ['isCheck'],
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 400
  },
  methods: {
    async loadPage (requestTemplateId) {
      this.requestTemplateId = requestTemplateId
      this.isParmasChanged = false
      const { statusCode, data } = await getRequestFormTemplateData(this.requestTemplateId)
      if (statusCode === 'OK') {
        this.formId = data.id
        if (data.items !== null && data.items.length > 0) {
          data.items.sort(this.compare('sort'))
          this.finalElement = [
            {
              itemGroup: data.items[0].itemGroup,
              itemGroupName: data.items[0].itemGroupName,
              attrs: data.items
            }
          ]
        }
      }
      this.getAllDataModels()
    },
    // 保存表单信息
    async saveMsgForm (nextStep) {
      // nextStep 1新增 2下一步 3切换tab 4上一步
      let tmp = [].concat(...JSON.parse(JSON.stringify(this.finalElement)).map(l => l.attrs))
      tmp.forEach((l, index) => {
        l.sort = index + 1
        if (!isNaN(l.id) || l.id.startsWith('c_')) {
          l.id = ''
        }
      })
      tmp.sort(this.compare('sort'))
      let res = {
        id: this.formId,
        items: tmp
      }
      const { statusCode } = await saveRequsetForm(this.requestTemplateId, res)
      if (statusCode === 'OK') {
        if (![2, 3, 4].includes(nextStep)) {
          this.$Notice.success({
            title: this.$t('successful'),
            desc: this.$t('successful')
          })
        }
        this.isParmasChanged = false
        if (nextStep === 1) {
          this.loadPage(this.requestTemplateId)
        } else if (nextStep === 2) {
          this.$emit('gotoStep', this.requestTemplateId, 'forward')
        } else if (nextStep === 3) {
          this.$emit('changTab', 'dataForm')
        } else if (nextStep === 4) {
          this.$emit('gotoStep', this.requestTemplateId, 'backward')
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
    selectElement (itemIndex, eleIndex) {
      this.finalElement[itemIndex].attrs.forEach(item => {
        item.isActive = false
      })
      this.finalElement[itemIndex].attrs[eleIndex].isActive = true
      this.editElement = this.finalElement[itemIndex].attrs[eleIndex]
      this.openPanel = '1'
    },
    removeForm (itemIndex, eleIndex, element) {
      this.finalElement[itemIndex].attrs.splice(eleIndex, 1)
      this.openPanel = ''
      this.paramsChanged()
    },
    // 控制保存按钮
    isSaveBtnActive () {
      let res = false
      return res
    },
    paramsChanged () {
      this.isParmasChanged = true
    },
    panalStatus () {
      return this.isParmasChanged
    },
    tabChange () {
      if (this.isCheck !== 'Y') {
        this.saveMsgForm(3)
      }
    },
    gotoNext () {
      if (this.isCheck === 'Y') {
        this.$emit('gotoStep', this.requestTemplateId, 'forward')
      } else {
        this.saveMsgForm(2)
      }
    },
    gotoForward () {
      if (this.isCheck === 'Y') {
        this.$emit('gotoStep', this.requestTemplateId, 'backward')
      } else {
        this.saveMsgForm(4)
      }
    },
    cloneDog (val) {
      if (this.$parent.isCheck === 'Y') return
      let newItem = JSON.parse(JSON.stringify(val))
      const itemNo = this.generateRandomString()
      newItem.id = 'c_' + itemNo
      newItem.title = newItem.title + itemNo
      newItem.name = newItem.name + itemNo
      newItem.isActive = true
      this.specialId = newItem.id
      this.paramsChanged()
      this.finalElement[0].attrs.forEach(item => {
        item.isActive = false
      })
      return newItem
    },
    generateRandomString () {
      let result = ''
      for (let i = 0; i < 4; i++) {
        result += Math.floor(Math.random() * 10)
      }
      return result
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
    // 获取wecmdb下拉类型entity值
    async getAllDataModels () {
      const { data, status } = await getAllDataModels()
      if (status === 'OK') {
        this.allEntityList = []
        const sortData = data.map(_ => {
          return {
            ..._,
            entities: _.entities.sort(function (a, b) {
              var s = a.name.toLowerCase()
              var t = b.name.toLowerCase()
              if (s < t) return -1
              if (s > t) return 1
            })
          }
        })
        sortData.forEach(i => {
          i.entities.forEach(j => {
            this.allEntityList.push(`${j.packageName}:${j.name}`)
          })
        })
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
  margin: 8px 0;
}
.title {
  font-size: 14px;
  font-weight: bold;
  margin: 12px 0;
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
    height: 8px;
    border-radius: 8px;
    background-color: #c6eafe;
    box-sizing: content-box;
  }
}
.custom-title {
  width: 80px;
  display: inline-block;
  text-align: right;
  word-wrap: break-word;
}
.custom-item {
  width: calc(100% - 130px);
  display: inline-block;
}
.btn-footer-margin {
  margin: 0 6px;
}
</style>
