<template>
  <div ref="maxheight">
    <div class="msg-form-container">
      <div class="left">
        <div :style="{ height: MODALHEIGHT + 32 + 'px', overflow: 'auto', padding: '0 8px' }">
          <!--自定义表单项-->
          <Divider orientation="left" size="small">{{ $t('custom_form') }}</Divider>
          <CustomDraggable :sortable="$parent.isCheck !== 'Y'" :clone="cloneDog"></CustomDraggable>
          <!--表单项组件库-->
          <Divider orientation="left" size="small">{{ $t('tw_template_library') }}</Divider>
          <ComponentLibraryList ref="libraryList" formType="requestInfo"></ComponentLibraryList>
        </div>
      </div>
      <!--表单预览-->
      <div class="center">
        <div :style="{ height: MODALHEIGHT + 32 + 'px', overflow: 'auto', paddingBottom: '10px' }">
          <div class="title">
            <div class="title-text">
              {{ $t('request_form_details') }}
              <span class="underline"></span>
            </div>
          </div>
          <Divider orientation="left" size="small">{{ $t('tw_information_form') }}</Divider>
          <template v-for="(item, itemIndex) in finalElement">
            <div :key="itemIndex" style="border: 2px dashed #A2EF4D; margin: 8px 0; padding: 8px;min-height: 48px;">
              <draggable
                class="dragArea"
                style="min-height: 40px;"
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
                  <Checkbox v-model="element.checked" style="margin:0;"></Checkbox>
                  <div class="require">
                    <Icon v-if="element.required === 'yes'" size="8" type="ios-medical" />
                  </div>
                  <div
                    class="custom-title"
                    :style="
                      ['calculate', 'textarea'].includes(element.elementType)
                        ? 'vertical-align: top;word-break: break-all;'
                        : ''
                    "
                  >
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
                    :multiple="element.multiple === 'yes'"
                    class="custom-item"
                  >
                    <Option v-for="item in computedOption(element)" :value="item.value" :key="item.label">{{
                      item.label
                    }}</Option>
                  </Select>
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
                    ghost
                    style="width:24px;display:flex;justify-content:center;"
                  >
                    <Icon type="ios-close" size="24"></Icon>
                  </Button>
                </div>
              </draggable>
            </div>
            <Button
              :key="itemIndex + '-'"
              :disabled="getAddComponentDisabled(item.attrs)"
              @click="createComponentLibrary"
              size="small"
              >新建组件</Button
            >
          </template>
        </div>
      </div>
      <!--属性设置-->
      <div class="right">
        <div :style="{ height: MODALHEIGHT + 32 + 'px', overflow: 'auto' }">
          <Collapse v-model="openPanel">
            <Panel name="1">
              {{ $t('general_attributes') }}
              <div slot="content">
                <Form
                  ref="attrForm"
                  :model="editElement"
                  :rules="ruleForm"
                  :label-width="80"
                  :disabled="editElement.controlSwitch === 'yes'"
                >
                  <FormItem :label="$t('display_name')">
                    <Input
                      v-model="editElement.title"
                      @on-change="paramsChanged"
                      :disabled="$parent.isCheck === 'Y'"
                      placeholder=""
                    ></Input>
                  </FormItem>
                  <FormItem :label="$t('tw_code')" prop="name" style="margin-bottom:20px;">
                    <Input
                      v-model="editElement.name"
                      @on-change="paramsChanged"
                      :disabled="$parent.isCheck === 'Y' || editElement.entity !== ''"
                      placeholder=""
                    ></Input>
                  </FormItem>
                  <FormItem :label="$t('tw_form_type')">
                    <Select
                      v-model="editElement.elementType"
                      :disabled="true"
                      @on-change="editElement.defaultValue = ''"
                    >
                      <Option value="input">Input</Option>
                      <Option value="select">Select</Option>
                      <Option value="textarea">Textarea</Option>
                      <Option value="wecmdbEntity">{{ $t('tw_entity_data_items') }}</Option>
                      <Option value="datePicker">DatePicker</Option>
                    </Select>
                  </FormItem>
                  <!--数据集-->
                  <FormItem
                    v-if="editElement.elementType === 'select' && editElement.entity === ''"
                    :label="$t('tw_options')"
                  >
                    <Input :value="getDataOptionsDisplay" disabled style="width:calc(100% - 28px)"></Input>
                    <Button @click.stop="dataOptionsMgmt" type="success" size="small" icon="md-add"></Button>
                  </FormItem>
                  <!--数据源-->
                  <FormItem v-if="editElement.elementType === 'select' && editElement.entity" :label="$t('tw_options')">
                    <Input v-model="editElement.dataOptions" disabled></Input>
                  </FormItem>
                  <!--模型数据项-->
                  <FormItem v-if="editElement.elementType === 'wecmdbEntity'" :label="$t('tw_options')">
                    <Select
                      v-model="editElement.dataOptions"
                      filterable
                      @on-change="paramsChanged"
                      :disabled="$parent.isCheck === 'Y'"
                    >
                      <Option v-for="i in allEntityList" :value="i" :key="i">{{ i }}</Option>
                    </Select>
                  </FormItem>
                  <div style="display:flex;justify-content:space-between;flex-wrap:wrap;">
                    <!--控制审批/任务-->
                    <Form v-if="['select', 'wecmdbEntity'].includes(editElement.elementType)" :label-width="80">
                      <FormItem label="控制审批/任务">
                        <i-switch
                          v-model="editElement.controlSwitch"
                          true-value="yes"
                          false-value="no"
                          :disabled="$parent.isCheck === 'Y'"
                          @on-change="
                            controlSwitchChange($event)
                            paramsChanged()
                          "
                          size="default"
                        />
                      </FormItem>
                    </Form>
                    <FormItem
                      :label="$t('tw_multiple')"
                      v-if="['select', 'wecmdbEntity'].includes(editElement.elementType)"
                    >
                      <i-switch
                        v-model="editElement.multiple"
                        true-value="yes"
                        false-value="no"
                        :disabled="$parent.isCheck === 'Y'"
                        @on-change="paramsChanged"
                        size="default"
                      />
                    </FormItem>
                    <FormItem :label="$t('editable')">
                      <i-switch
                        v-model="editElement.isEdit"
                        true-value="yes"
                        false-value="no"
                        :disabled="$parent.isCheck === 'Y'"
                        @on-change="paramsChanged"
                        size="default"
                      />
                    </FormItem>
                    <FormItem :label="$t('required')">
                      <i-switch
                        v-model="editElement.required"
                        true-value="yes"
                        false-value="no"
                        :disabled="$parent.isCheck === 'Y' || editElement.controlSwitch === 'yes'"
                        @on-change="paramsChanged"
                        size="default"
                      />
                    </FormItem>
                    <FormItem :label="$t('display')">
                      <i-switch
                        v-model="editElement.inDisplayName"
                        true-value="yes"
                        false-value="no"
                        :disabled="$parent.isCheck === 'Y'"
                        @on-change="paramsChanged"
                        size="default"
                      />
                    </FormItem>
                    <FormItem :label="$t('tw_default_empty')">
                      <i-switch
                        v-model="editElement.defaultClear"
                        true-value="yes"
                        false-value="no"
                        :disabled="$parent.isCheck === 'Y'"
                        @on-change="paramsChanged"
                        size="default"
                      />
                    </FormItem>
                  </div>
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
                <Form :label-width="80" label-position="left" :disabled="editElement.controlSwitch === 'yes'">
                  <FormItem label="" :label-width="0">
                    <HiddenCondition
                      ref="hiddenCondition"
                      :disabled="$parent.isCheck === 'Y'"
                      :finalElement="finalElement"
                      :editElement="editElement"
                      v-model="editElement.hiddenCondition"
                    ></HiddenCondition>
                  </FormItem>
                  <FormItem :label="$t('validation_rules')">
                    <Input
                      v-model="editElement.regular"
                      :disabled="$parent.isCheck === 'Y'"
                      :placeholder="$t('only_supports_regular')"
                      @on-change="paramsChanged"
                    ></Input>
                  </FormItem>
                  <FormItem :label="$t('data_item') + $t('constraints')">
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
      </div>
    </div>
    <!--数据集弹框-->
    <DataSourceConfig ref="dataSourceConfigRef" @setDataOptions="setDataOptions"></DataSourceConfig>
    <!--组件库弹框-->
    <ComponentLibraryModal
      ref="library"
      formType="requestInfo"
      v-model="componentVisible"
      :checkedList="componentCheckedList"
      @fetchList="$refs.libraryList.handleSearch()"
    />
    <div class="footer">
      <div class="content">
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
  </div>
</template>

<script>
import { saveRequsetForm, getRequestFormTemplateData, getAllDataModels, cleanFilterData } from '@/api/server.js'
import draggable from 'vuedraggable'
import CustomDraggable from './components/custom-draggable.vue'
import DataSourceConfig from './data-source-config.vue'
import ComponentLibraryModal from './components/component-library-modal.vue'
import ComponentLibraryList from './components/component-library-list.vue'
import HiddenCondition from './components/hidden-condition.vue'
import { uniqueArr, deepClone, findFirstDuplicateIndex } from '@/pages/util'
export default {
  name: 'form-select',
  components: {
    draggable,
    CustomDraggable,
    DataSourceConfig,
    ComponentLibraryModal,
    ComponentLibraryList,
    HiddenCondition
  },
  data () {
    return {
      requestTemplateId: '',
      isParmasChanged: false, // 参数变化标志位
      MODALHEIGHT: 200,
      finalElement: [
        {
          itemGroup: 'requestInfo',
          itemGroupName: this.$t('tw_information_form'),
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
        multiple: 'no',
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
        dataOptions: '[]',
        refEntity: '',
        refPackageName: '',
        controlSwitch: 'no', // 控制审批/任务(下拉类型才有)
        hiddenCondition: [] // 隐藏条件
      },
      allEntityList: [],
      formId: '', // 缓存表单id，供编辑使用
      componentVisible: false, // 组件库弹窗
      componentCheckedList: [], // 当前选中组件库数据
      ruleForm: {
        // 校验编码不能重复
        name: [
          {
            required: true,
            validator: (rule, value, callback) => {
              const arr = this.finalElement[0].attrs.map(i => i.name)
              const index = findFirstDuplicateIndex(arr)
              if (index > -1 && this.finalElement[0].attrs[index].name === this.editElement.name) {
                return callback(new Error('编码不能重复'))
              } else {
                callback()
              }
            },
            trigger: 'change'
          }
        ]
      }
    }
  },
  props: ['isCheck'],
  computed: {
    // 数据集回显
    getDataOptionsDisplay () {
      const options = JSON.parse(this.editElement.dataOptions || '[]')
      const labelArr = options.map(item => item.label)
      return labelArr.join(',')
    },
    // 新增组件库按钮禁用
    getAddComponentDisabled () {
      return function (val) {
        const checkedList = val.filter(item => item.checked) || []
        return checkedList.length === 0
      }
    }
  },
  mounted () {
    const clientHeight = document.documentElement.clientHeight
    this.MODALHEIGHT = clientHeight - this.$refs.maxheight.getBoundingClientRect().top - 90
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
      // 获取模型数据项下拉值
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
      this.openPanel = ['1', '2']
      this.$refs.attrForm.validateField('name') // 编码重复校验
      this.$refs.hiddenCondition.removeConditionsByAttrs(this.editElement.hiddenCondition) // 隐藏条件删除多余属性
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
    controlSwitchChange (val) {
      // 关闭【控制审批任务】开关，清除数据
      if (val === 'no') {
        cleanFilterData(this.requestTemplateId, 'data')
      } else if (val === 'yes') {
        this.editElement.required = 'yes'
      }
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
    log () {
      this.finalElement.forEach(l => {
        // 从组件库拖拽进来的表单组，数据需要额外处理
        const cloneAttrs = deepClone(l.attrs || [])
        var deleteIdx = ''
        l.attrs.forEach((attr, idx) => {
          for (let key of Object.keys(attr)) {
            if (!isNaN(Number(key))) {
              deleteIdx = idx
              cloneAttrs.push(attr[key])
            }
          }
        })
        if (typeof deleteIdx === 'number') {
          cloneAttrs.splice(deleteIdx, 1)
        }
        l.attrs = cloneAttrs
        // 组件库拖拽有重复元素，给出提示，并过滤数据
        const { arr, sameArr } = uniqueArr(deepClone(l.attrs))
        l.attrs = arr
        const titleArr = sameArr.map(i => i.title) || []
        const message = titleArr.join('、')
        if (message) {
          this.$Notice.warning({
            title: this.$t('warning'),
            render: h => {
              return (
                <div style="word-break:break-all;">
                  表单已有表单项<span style="color: red;">{message}</span>,已过滤
                </div>
              )
            }
          })
        } else {
          this.$Notice.success({
            title: this.$t('successful'),
            desc: this.$t('successful')
          })
        }
        // 处理拖拽进来的表单项
        l.attrs.forEach(attr => {
          attr.itemGroup = l.itemGroup
          attr.itemGroupName = l.itemGroupName
          if (attr.isActive) {
            this.editElement = attr
            this.openPanel = ['1', '2']
          }
        })
      })
    },
    // 获取模型数据项下拉值
    async getAllDataModels () {
      const { data, statusCode } = await getAllDataModels()
      if (statusCode === 'OK') {
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
    },
    // #region 普通select数据集配置逻辑
    dataOptionsMgmt () {
      let newDataOptions = JSON.parse(this.editElement.dataOptions || '[]')
      this.$refs.dataSourceConfigRef.loadPage(newDataOptions)
    },
    setDataOptions (options) {
      if (options && options.length > 0) {
        this.editElement.dataOptions = JSON.stringify(options)
      } else {
        this.editElement.dataOptions = ''
      }
    },
    computedOption (element) {
      let res = []
      if (element.elementType === 'select') {
        res = JSON.parse(element.dataOptions || '[]')
      } else if (element.elementType === 'wecmdbEntity') {
      }
      return res
    },
    // 新增组件库
    createComponentLibrary () {
      this.componentVisible = true
      this.componentCheckedList = this.finalElement[0].attrs.filter(i => i.checked === true)
      this.$refs.library.init()
    }
  }
}
</script>
<style scoped lang="scss">
.msg-form-container {
  display: flex;
  .left {
    width: 360px;
    border: 1px solid #dcdee2;
  }
  .center {
    flex: 1;
    border: 1px solid #dcdee2;
    padding: 0 16px;
    width: 57%;
    margin: 0 4px;
  }
  .right {
    width: 360px;
    border: 1px solid #dcdee2;
  }
}
.active-zone {
  color: #338cf0;
}
.ivu-form-item {
  margin-bottom: 16px;
}
.list-group-item- {
  display: flex;
  align-items: center;
  justify-content: space-around;
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
  width: 125px;
  display: flex;
  align-items: center;
  text-align: left;
  word-wrap: break-word;
}
.custom-item {
  width: calc(100% - 190px);
  display: inline-block;
}
.btn-footer-margin {
  margin: 0 6px;
}

.footer {
  position: fixed; /* 使用 fixed 定位，使其固定在页面底部 */
  left: 0;
  bottom: 0;
  width: 100%; /* 撑满整个页面宽度 */
  background-color: white; /* 设置背景色，可根据需求修改 */
  z-index: 10;
}

.content {
  text-align: center; /* 居中内容 */
  padding: 10px; /* 可根据需求调整内容与边框的间距 */
}

.require {
  color: #ed4014;
  width: 6px;
  display: flex;
  align-items: center;
}
</style>
