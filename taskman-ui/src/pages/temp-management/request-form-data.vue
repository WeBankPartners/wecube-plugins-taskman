<template>
  <div>
    <Row>
      <!--自定义表单项-->
      <Col span="5" style="border: 1px solid #dcdee2; padding: 0 16px">
        <div :style="{ height: MODALHEIGHT + 32 + 'px', overflow: 'auto' }">
          <Divider plain>{{ $t('custom_form') }}</Divider>
          <CustomDraggable :sortable="$parent.isCheck !== 'Y'" :clone="cloneDog"></CustomDraggable>
        </div>
      </Col>
      <!--表单预览-->
      <Col span="14" style="border: 1px solid #dcdee2; padding: 0 16px; width: 57%; margin: 0 4px">
        <div :style="{ height: MODALHEIGHT + 30 + 'px', overflow: 'auto' }">
          <Divider>{{ $t('tw_preview') }}</Divider>
          <div class="title">
            <div class="title-text">
              {{ $t('request_form_details') }}
              <span class="underline"></span>
            </div>
          </div>
          <Divider orientation="left" size="small">{{ $t('tw_data_form') }}</Divider>
          <Form :label-width="100" v-if="dataFormInfo.associationWorkflow">
            <FormItem>
              <span slot="label" style="font-size: 12px;">{{ $t('tw_choose_object') }}</span>
              <Select style="width: 30%">
                <Option v-for="item in []" :value="item.id" :key="item.id">{{ item.displayName }}</Option>
              </Select>
            </FormItem>
          </Form>
          <div>
            <div class="radio-group">
              <div
                v-for="(groupItem, index) in dataFormInfo.groups"
                :key="index"
                :class="{
                  'radio-group-radio': true,
                  'radio-group-custom': groupItem.itemGroupType === 'custom',
                  'radio-group-workflow': groupItem.itemGroupType === 'workflow',
                  'radio-group-optional': groupItem.itemGroupType === 'optional'
                }"
                :style="activeStyle(groupItem)"
              >
                <Icon @click="editGroupItem(groupItem)" type="md-create" color="#2d8cf0" :size="20" />
                <span @click="editGroupCustomItems(groupItem)">
                  {{ `${groupItem.itemGroupName}` }}
                </span>
                <Icon
                  v-if="isCheck !== 'Y'"
                  @click="removeGroupItem(groupItem)"
                  type="md-close"
                  color="#ed4014"
                  :size="20"
                />
              </div>
              <span>
                <Button
                  v-if="isCheck !== 'Y'"
                  style="margin-top: -5px;"
                  @click="beforeSelectItemGroup"
                  type="primary"
                  shape="circle"
                  icon="md-add"
                ></Button>
              </span>
            </div>
          </div>
          <template v-if="finalElement.length === 1 && finalElement[0].itemGroup !== ''">
            <div
              v-for="(item, itemIndex) in finalElement"
              :key="itemIndex"
              style="border: 2px dotted #A2EF4D; margin: 8px 0; padding: 8px;min-height: 48px;"
            >
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
                    <span v-if="element.required === 'yes'" style="color: red;">
                      *
                    </span>
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
                    v-if="['calculate', 'textarea'].includes(element.elementType)"
                    :disabled="element.isEdit === 'no'"
                    v-model="element.defaultValue"
                    type="textarea"
                    :rows="2"
                    class="custom-item"
                  />
                  <Select
                    v-if="element.elementType === 'select'"
                    :disabled="element.isEdit === 'no'"
                    class="custom-item"
                    :multiple="element.multiple === 'yes'"
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
                    icon="ios-close"
                    ghost
                  ></Button>
                </div>
              </draggable>
            </div>
            <div style="text-align: right;">
              <Button v-if="isCheck !== 'Y'" type="primary" size="small" ghost @click="saveGroup(1)">{{
                $t('save')
              }}</Button>
              <Button size="small" @click="restoreGroup">{{ $t('tw_restore') }}</Button>
            </div>
          </template>
        </div>
      </Col>
      <!--属性设置-->
      <Col span="5" style="border: 1px solid #dcdee2">
        <div :style="{ height: MODALHEIGHT + 32 + 'px', overflow: 'auto' }">
          <Collapse v-model="openPanel">
            <Panel name="1">
              {{ $t('general_attributes') }}
              <div slot="content">
                <Form :label-width="80" :disabled="editElement.controlSwitch === 'yes'">
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
                      <Option value="wecmdbEntity">{{ $t('tw_entity_data_items') }}</Option>
                      <Option value="datePicker">DatePicker</Option>
                      <Option value="calculate">Calculate</Option>
                    </Select>
                  </FormItem>
                  <!--数据集-->
                  <FormItem
                    v-if="editElement.elementType === 'select' && editElement.entity === ''"
                    :label="$t('data_set')"
                  >
                    <!-- <Input v-model="editElement.dataOptions" disabled style="width:calc(100% - 38px)"></Input> -->
                    <Input :value="getDataOptionsDisplay" disabled style="width:calc(100% - 38px)"></Input>
                    <Button
                      :disabled="$parent.isCheck === 'Y'"
                      @click.stop="dataOptionsMgmt"
                      type="primary"
                      icon="md-add"
                    ></Button>
                  </FormItem>
                  <!--数据源-->
                  <FormItem
                    v-if="editElement.elementType === 'select' && editElement.entity"
                    :label="$t('data_source')"
                  >
                    <Input v-model="editElement.dataOptions" disabled></Input>
                  </FormItem>
                  <!--模型数据项-->
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
                  <!--控制审批/任务-->
                  <Form :label-width="80">
                    <FormItem v-if="['select', 'wecmdbEntity'].includes(editElement.elementType)" label="控制审批/任务">
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
                  <FormItem :label="$t('defaults')">
                    <Input
                      v-model="editElement.defaultValue"
                      :disabled="$parent.isCheck === 'Y' || editElement.defaultClear === 'yes'"
                      placeholder=""
                      @on-change="paramsChanged"
                    ></Input>
                  </FormItem>
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
                <Form :label-width="80" :disabled="editElement.controlSwitch === 'yes'">
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
                <Form :label-width="80" :disabled="editElement.controlSwitch === 'yes'">
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
    <Modal v-model="showSelectModel" :title="$t('tw_create_form')" :mask-closable="false">
      <div style="margin: 40px 0 60px 0">
        <Form :label-width="120">
          <FormItem :label="$t('tw_form_template_type')">
            <Select
              style="width: 80%"
              v-model="itemGroup"
              ref="selectRef"
              v-if="showSelectModel"
              filterable
              clearable
              @on-open-change="clearQuery"
            >
              <OptionGroup v-for="itemGroup in groupOptions" :label="itemGroup.formTypeName" :key="itemGroup.formType">
                <Option v-for="item in itemGroup.entities" :value="item" :key="item">{{ item }}</Option>
              </OptionGroup>
            </Select>
          </FormItem>
        </Form>
      </div>
      <template #footer>
        <Button @click="showSelectModel = false">{{ $t('cancel') }}</Button>
        <Button @click="okSelect" :disabled="!itemGroup" type="primary">{{ $t('confirm') }}</Button>
      </template>
    </Modal>
    <!-- 自定义表单配置 -->
    <RequestFormDataCustom
      ref="requestFormDataCustomRef"
      @reloadParentPage="isLoadLastGroup"
      :isCheck="isCheck"
      module="data-form"
      v-show="['custom'].includes(itemGroupType)"
    ></RequestFormDataCustom>

    <!-- 编排表单配置 -->
    <RequestFormDataWorkflow
      ref="requestFormDataWorkflowRef"
      :isCustomItemEditable="true"
      @reloadParentPage="isLoadLastGroup"
      :isCheck="isCheck"
      module="data-form"
      v-show="['workflow', 'optional'].includes(itemGroupType)"
    ></RequestFormDataWorkflow>
    <DataSourceConfig ref="dataSourceConfigRef" @setDataOptions="setDataOptions"></DataSourceConfig>
    <div class="footer">
      <div class="content">
        <Button @click="gotoForward" ghost type="primary" class="btn-footer-margin">{{ $t('forward') }}</Button>
        <Button @click="gotoNext" type="primary" class="btn-footer-margin">{{ $t('next') }}</Button>
      </div>
    </div>
  </div>
</template>

<script>
import {
  getAllDataModels,
  getEntityByTemplateId,
  getRequestDataForm,
  deleteRequestGroupForm,
  saveRequestGroupCustomForm,
  cleanFilterData
} from '@/api/server.js'
import draggable from 'vuedraggable'
import RequestFormDataCustom from './request-form-data-custom.vue'
import RequestFormDataWorkflow from './request-form-data-workflow.vue'
import CustomDraggable from './components/custom-draggable.vue'
import DataSourceConfig from './data-source-config.vue'
export default {
  name: 'form-select',
  components: {
    draggable,
    RequestFormDataCustom,
    RequestFormDataWorkflow,
    CustomDraggable,
    DataSourceConfig
  },
  data () {
    return {
      isParmasChanged: false, // 参数变化标志位
      MODALHEIGHT: 200,
      dataFormInfo: {
        associationWorkflow: false,
        formTemplateId: '',
        groups: []
      },
      groups: [], // 已存在组
      showSelectModel: false, // 组选择框
      itemGroupType: '', // 选中的组类型
      itemGroup: '', // 选中的组信息
      groupOptions: [],
      finalElement: [
        // 待编辑组信息
        {
          itemGroupId: '',
          formTemplateId: '',
          requestTemplateId: '',
          itemGroup: '',
          itemGroupName: '',
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
        controlSwitch: 'no' // 控制审批/任务(下拉类型才有)
      },
      allEntityList: [],
      groupStyle: {
        custom: {
          border: '1px solid #b886f8',
          color: '#b886f8',
          'border-radius': '14px'
        },
        workflow: {
          border: '1px solid #cba43f',
          color: '#cba43f',
          'border-radius': '14px'
        },
        optional: {
          border: '1px solid #81b337',
          color: '#81b337',
          'border-radius': '14px'
        }
      },
      displayLastGroup: false, // 控制group显示，在新增时显示最后一个，其余显示当前值\
      nextGroupInfo: {}
    }
  },
  props: ['isCheck'],
  computed: {
    activeStyle () {
      return function (item) {
        let color = '#fff'
        let obj = this.finalElement[0]
        if (item.itemGroupId === obj.itemGroupId) {
          if (item.itemGroupType === 'workflow') {
            color = '#ebdcb4'
          } else if (item.itemGroupType === 'custom') {
            color = 'rgba(184, 134, 248, 0.6)'
          } else if (item.itemGroupType === 'optional') {
            color = 'rgba(129, 179, 55, 0.6)'
          } else if (!item.itemGroupType) {
            color = 'rgb(45, 140, 240)'
          }
        }
        return { background: color }
      }
    },
    getDataOptionsDisplay () {
      const options = JSON.parse(this.editElement.dataOptions || '[]')
      const labelArr = options.map(item => item.label)
      return labelArr.join(',')
    }
  },
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 400
  },
  methods: {
    async loadPage (requestTemplateId) {
      if (requestTemplateId) {
        this.requestTemplateId = requestTemplateId
      }
      this.isParmasChanged = false
      const { statusCode, data } = await getRequestDataForm(this.requestTemplateId)
      if (statusCode === 'OK') {
        this.dataFormInfo = data
        let groups = this.dataFormInfo.groups
        if (this.displayLastGroup) {
          const group = groups[groups.length - 1]
          this.editGroupCustomItems(group, false)
          // this.displayLastGroup = false
        } else {
          let itemGroupId = this.finalElement[0].itemGroupId
          const findGroup = groups.find(form => form.itemGroupId === itemGroupId)
          if (findGroup) {
            this.editGroupCustomItems(findGroup, false)
          } else {
            if (groups.length > 0) {
              this.editGroupCustomItems(groups[0], false)
            }
          }
        }
      }
      this.getAllDataModels()
    },
    // 放弃编辑，重新加载group
    restoreGroup () {
      this.loadPage()
    },
    isLoadLastGroup (val) {
      this.displayLastGroup = val
      this.loadPage()
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
    log (item) {
      this.finalElement.forEach(l => {
        l.attrs.forEach(attr => {
          attr.itemGroup = l.itemGroup
          attr.itemGroupName = l.itemGroupName
          if (attr.id === this.specialId) {
            this.editElement = attr
            if (this.editElement.multiple === 'Y') {
              this.editElement.multiple = 'yes'
            } else if (this.editElement.multiple === 'N') {
              this.editElement.multiple = 'no'
            }
            this.openPanel = '1'
          }
        })
      })
    },
    // 查询可添加的组
    beforeSelectItemGroup () {
      if (this.finalElement[0].itemGroupId === '') {
        this.selectItemGroup()
      } else {
        this.saveGroup(7)
      }
    },
    async selectItemGroup () {
      this.itemGroup = ''
      // this.groupOptions = []
      const { statusCode, data } = await getEntityByTemplateId(this.requestTemplateId)
      if (statusCode === 'OK') {
        // workflow  2.编排数据, 3.optional 自选数据项表单,  custom 1.自定义表单
        this.$nextTick(() => {
          const existGroup = this.dataFormInfo.groups
            .filter(g => g.itemGroupType !== 'custom')
            .map(g => {
              return g.itemGroupName
            })
          this.groupOptions = data.map(d => {
            if (d.formType === 'custom') {
              d.formTypeName = this.$t('tw_custom_form')
              d.entities = ['custom']
            } else if (d.formType === 'workflow') {
              d.formTypeName = this.$t('tw_orchestration')
              let entities = d.entities || []
              d.entities = entities.filter(entity => !existGroup.includes(entity))
            } else if (d.formType === 'optional') {
              d.formTypeName = this.$t('tw_custom')
              let entities = d.entities || []
              d.entities = entities.filter(entity => !existGroup.includes(entity))
            }
            return d
          })
          this.showSelectModel = true
        })
      }
    },
    // 选择可添加的组
    okSelect () {
      this.showSelectModel = false
      if (this.itemGroup === 'custom') {
        this.itemGroupType = 'custom'
        let params = {
          requestTemplateId: this.requestTemplateId,
          // formTemplateId: this.requestTemplateId,
          isAdd: true,
          itemGroupType: this.itemGroupType,
          itemGroup: this.itemGroup + (this.dataFormInfo.groups.length + 1),
          itemGroupId: '',
          itemGroupSort: this.dataFormInfo.groups.length + 1
        }
        this.$refs.requestFormDataCustomRef.loadPage(params)
      } else {
        this.groupOptions.forEach(group => {
          const findEntity = group.entities.findIndex(e => e === this.itemGroup)
          if (findEntity !== -1) {
            this.itemGroupType = group.formType
          }
        })
        this.$nextTick(() => {
          if (['workflow', 'optional'].includes(this.itemGroupType)) {
            let params = {
              requestTemplateId: this.requestTemplateId,
              // formTemplateId: this.requestTemplateId,
              isAdd: true,
              itemGroupType: this.itemGroupType,
              itemGroup: this.itemGroup,
              itemGroupId: '',
              itemGroupSort: this.dataFormInfo.groups.length + 1
            }
            this.$refs.requestFormDataWorkflowRef.loadPage(params)
          }
        })
      }
    },
    // 编辑组自定义属性
    editGroupCustomItems (groupItem, isNeedSaveFirst = true) {
      this.nextGroupInfo = groupItem
      this.displayLastGroup = false
      if (isNeedSaveFirst && this.isCheck !== 'Y') {
        this.saveGroup(4, groupItem)
      } else {
        this.updateFinalElement(groupItem)
      }
      // if (this.isParmasChanged) {
      //   this.$Modal.confirm({
      //     title: `${this.$t('confirm_discarding_changes')}`,
      //     content: `${this.finalElement[0].itemGroupName}:${this.$t('params_edit_confirm')}`,
      //     'z-index': 1000000,
      //     okText: this.$t('save'),
      //     cancelText: this.$t('abandon'),
      //     onOk: async () => {
      //       // this.saveGroup(1)
      //       this.saveGroup(4, groupItem)
      //     },
      //     onCancel: () => {
      //       this.updateFinalElement(groupItem)
      //     }
      //   })
      // } else {
      //   this.updateFinalElement(groupItem)
      // }
    },
    updateFinalElement (groupItem) {
      this.finalElement = [
        {
          itemGroupId: groupItem.itemGroupId,
          formTemplateId: this.requestTemplateId,
          requestTemplateId: this.requestTemplateId,
          itemGroup: groupItem.itemGroup,
          itemGroupName: groupItem.itemGroupName,
          attrs: groupItem.items || []
        }
      ]
      this.openPanel = ''
    },
    // 删除组
    async removeGroupItem (groupItem) {
      this.nextGroupInfo = groupItem
      this.saveGroup(6)
      // this.$Modal.confirm({
      //   title: this.$t('confirm_delete'),
      //   'z-index': 1000000,
      //   loading: true,
      //   onOk: async () => {
      //     this.$Modal.remove()
      //     const { statusCode } = await deleteRequestGroupForm(groupItem.itemGroupId, this.requestTemplateId)
      //     if (statusCode === 'OK') {
      //       this.$Notice.success({
      //         title: this.$t('successful'),
      //         desc: this.$t('successful')
      //       })
      //       this.cancelGroup()
      //       this.loadPage()
      //     }
      //   },
      //   onCancel: () => {}
      // })
    },
    async confirmRemoveGroupItem () {
      this.$Modal.confirm({
        title: this.$t('confirm_delete'),
        'z-index': 1000000,
        loading: true,
        okText: this.$t('tw_request_confirm'),
        onOk: async () => {
          this.$Modal.remove()
          const { statusCode } = await deleteRequestGroupForm(this.nextGroupInfo.itemGroupId, this.requestTemplateId)
          if (statusCode === 'OK') {
            this.$Notice.success({
              title: this.$t('successful'),
              desc: this.$t('successful')
            })
            this.cancelGroup()
            this.loadPage()
          }
        },
        onCancel: () => {}
      })
    },
    // 编辑组信息
    editGroupItem (groupItem) {
      if (this.isCheck === 'Y') {
        this.openDrawer(groupItem)
      } else {
        this.saveGroup(5, groupItem)
      }
    },
    openDrawer (groupItem) {
      this.editGroupCustomItems(groupItem)
      if (groupItem.itemGroupType === 'custom') {
        this.itemGroupType = groupItem.itemGroupType
        let params = {
          requestTemplateId: this.requestTemplateId,
          taskTemplateId: '',
          isAdd: false,
          itemGroupName: groupItem.itemGroupName,
          itemGroupType: groupItem.itemGroupType,
          itemGroupId: groupItem.itemGroupId,
          itemGroup: ''
        }
        this.$refs.requestFormDataCustomRef.loadPage(params)
      }
      if (['workflow', 'optional'].includes(groupItem.itemGroupType)) {
        this.itemGroupType = groupItem.itemGroupType
        let params = {
          requestTemplateId: this.requestTemplateId,
          taskTemplateId: '',
          isAdd: false,
          itemGroupName: groupItem.itemGroupName,
          itemGroupType: groupItem.itemGroupType,
          itemGroupId: groupItem.itemGroupId,
          itemGroup: groupItem.itemGroup
        }
        this.$refs.requestFormDataWorkflowRef.loadPage(params)
      }
    },
    // 获取wecmdb下拉类型entity值
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
    // 选中自定义表单项
    selectElement (itemIndex, eleIndex) {
      this.finalElement[itemIndex].attrs.forEach(item => {
        item.isActive = false
      })
      this.finalElement[itemIndex].attrs[eleIndex].isActive = true
      this.editElement = this.finalElement[itemIndex].attrs[eleIndex]
      if (this.editElement.multiple === 'Y') {
        this.editElement.multiple = 'yes'
      } else if (this.editElement.multiple === 'N') {
        this.editElement.multiple = 'no'
      }
      this.openPanel = '1'
    },
    // 删除自定义表单项
    removeForm (itemIndex, eleIndex, element) {
      this.finalElement[itemIndex].attrs.splice(eleIndex, 1)
      this.openPanel = ''
      this.paramsChanged()
    },
    // 保存自定义表单项
    async saveGroup (nextStep, elememt) {
      // nextStep 1新增 2下一步 3切换tab 4 切换到目标group 5切换到目标group打开弹窗 6删除group 7选择组 8上一步
      let finalData = JSON.parse(JSON.stringify(this.finalElement[0]))
      finalData.items = finalData.attrs.map(attr => {
        if (attr.id.startsWith('c_')) {
          attr.id = ''
        }
        return attr
      })
      delete finalData.attrs
      finalData.items.forEach((item, itemIndex) => {
        item.sort = itemIndex + 1
      })
      const { statusCode } = await saveRequestGroupCustomForm(finalData)
      if (statusCode === 'OK') {
        if (![2, 3, 4, 5, 6, 7, 8].includes(nextStep)) {
          this.$Notice.success({
            title: this.$t('successful'),
            desc: this.$t('successful')
          })
        }
        this.isParmasChanged = false
        if (nextStep === 1) {
          // this.cancelGroup()
          this.loadPage()
        } else if (nextStep === 2) {
          this.$emit('gotoStep', this.requestTemplateId, 'forward')
        } else if (nextStep === 3) {
          this.$emit('changTab', 'msgForm')
        } else if (nextStep === 4) {
          this.updateFinalElement(this.nextGroupInfo)
          this.loadPage()
        } else if (nextStep === 5) {
          this.openDrawer(elememt)
        } else if (nextStep === 6) {
          this.confirmRemoveGroupItem()
        } else if (nextStep === 7) {
          this.selectItemGroup()
        } else if (nextStep === 8) {
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
    cancelGroup () {
      this.finalElement = [
        // 待编辑组信息
        {
          itemGroupId: '',
          formTemplateId: '',
          requestTemplateId: '',
          itemGroup: '',
          itemGroupName: '',
          attrs: []
        }
      ]
      this.isParmasChanged = false
    },
    panalStatus () {
      return this.isParmasChanged
    },
    tabChange () {
      if (this.isCheck !== 'Y') {
        this.saveGroup(3)
      }
    },
    gotoNext () {
      if (this.isCheck === 'Y') {
        this.$emit('gotoStep', this.requestTemplateId, 'forward')
        return
      }
      let finalData = JSON.parse(JSON.stringify(this.finalElement[0]))
      if (finalData.itemGroupId === '') {
        this.$emit('gotoStep', this.requestTemplateId, 'forward')
      } else {
        this.saveGroup(2)
      }
    },
    gotoForward () {
      if (this.isCheck === 'Y') {
        this.$emit('gotoStep', this.requestTemplateId, 'backward')
        return
      }
      let finalData = JSON.parse(JSON.stringify(this.finalElement[0]))
      if (finalData.itemGroupId === '') {
        this.$emit('gotoStep', this.requestTemplateId, 'backward')
      } else {
        this.saveGroup(8)
      }
    },
    // 此处的select在选中再点时，会将选中值当做条件过滤，这里清空query
    clearQuery () {
      this.$refs.selectRef.query = ''
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
    }
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
.ivu-select-dropdown {
  max-height: 300px !important;
}
.radio-group {
  margin-bottom: 15px;
}
/* 偶发样式在发布后不存在 */
.radio-group-radio {
  padding: 5px 15px;
  border-radius: 32px;
  font-size: 12px;
  cursor: pointer;
  margin: 4px;
  display: inline-block;
}
.radio-group-custom {
  border: 1px solid #b886f8;
  color: #b886f8;
}
.radio-group-workflow {
  border: 1px solid #cba43f;
  color: #cba43f;
}
.radio-group-optional {
  border: 1px solid #81b337;
  color: #81b337;
}
</style>
<style scoped lang="scss">
.active-zone {
  color: #338cf0;
}
.ivu-form-item {
  margin-bottom: 16px;
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

.radio-group {
  margin-bottom: 15px;
}
.radio-group-radio {
  padding: 5px 15px;
  border-radius: 32px;
  font-size: 12px;
  cursor: pointer;
  margin: 4px;
  display: inline-block;
}
.radio-group-custom {
  border: 1px solid #b886f8;
  color: #b886f8;
}
.radio-group-workflow {
  border: 1px solid #cba43f;
  color: #cba43f;
}
.radio-group-optional {
  border: 1px solid #81b337;
  color: #81b337;
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
</style>
