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
              <Select v-if="element.elementType === 'wecmdbEntity'" :placeholder="$t('tw_entity_data_items')"></Select>
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
              {{ $t('root_entity') }}
              <span class="underline"></span>
            </div>
          </div>
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
                    <Option v-for="item in element.dataOptions.split(',')" :value="item" :key="item">{{ item }}</Option>
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
                      <Option value="wecmdbEntity">{{ $t('tw_entity_data_items') }}</Option>
                      <Option value="datePicker">DatePicker</Option>
                      <Option value="calculate">Calculate</Option>
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
                  <FormItem
                    :label="$t('tw_multiple')"
                    v-if="['select', 'wecmdbEntity'].includes(editElement.elementType)"
                  >
                    <RadioGroup v-model="editElement.multiple" @on-change="paramsChanged">
                      <Radio label="yes" :disabled="$parent.isCheck === 'Y'">{{ $t('tw_yes') }}</Radio>
                      <Radio label="no" :disabled="$parent.isCheck === 'Y'">{{ $t('tw_no') }}</Radio>
                    </RadioGroup>
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
    <Modal v-model="showSelectModel" title="创建表单" :mask-closable="false">
      <div style="margin: 40px 0 60px 0">
        <Form :label-width="120">
          <FormItem :label="$t('表单模版类型')">
            <Select style="width: 80%" v-model="itemGroup" v-if="showSelectModel" filterable>
              <OptionGroup v-for="itemGroup in groupOptions" :label="itemGroup.formTypeName" :key="itemGroup.formType">
                <Option v-for="item in itemGroup.entities" :value="item" :key="item">{{ item }}</Option>
              </OptionGroup>
            </Select>
          </FormItem>
        </Form>
      </div>
      <template #footer>
        <Button @click="showSelectModel = false">{{ $t('cancel') }}</Button>
        <Button @click="okSelect" :disabled="itemGroup === ''" type="primary">{{ $t('confirm') }}</Button>
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
    <div style="text-align: center;margin-top: 16px;">
      <Button @click="gotoForward" ghost type="primary" class="btn-footer-margin">{{ $t('forward') }}</Button>
      <Button @click="gotoNext" type="primary" class="btn-footer-margin">{{ $t('next') }}</Button>
    </div>
  </div>
</template>

<script>
import {
  getAllDataModels,
  getEntityByTemplateId,
  getRequestDataForm,
  deleteRequestGroupForm,
  saveRequestGroupCustomForm
} from '@/api/server.js'
import draggable from 'vuedraggable'
import RequestFormDataCustom from './request-form-data-custom.vue'
import RequestFormDataWorkflow from './request-form-data-workflow.vue'
export default {
  name: 'form-select',
  data () {
    return {
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
          multiple: 'no',
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
          multiple: 'no',
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
          multiple: 'no',
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
          multiple: 'no',
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
          multiple: 'no',
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
        dataOptions: '',
        refEntity: '',
        refPackageName: ''
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
      this.groupOptions = []
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
              d.formTypeName = this.$t('编排entity表单')
              let entities = d.entities || []
              d.entities = entities.filter(entity => !existGroup.includes(entity))
            } else if (d.formType === 'optional') {
              d.formTypeName = this.$t('自选entity表单')
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
    // 编辑组信息
    editGroupItem (groupItem) {
      this.saveGroup(5, groupItem)
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
    },
    // 选中自定义表单项
    selectElement (itemIndex, eleIndex) {
      this.finalElement[itemIndex].attrs.forEach(item => {
        item.isActive = false
      })
      this.finalElement[itemIndex].attrs[eleIndex].isActive = true
      this.editElement = this.finalElement[itemIndex].attrs[eleIndex]
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
    }
  },
  components: {
    draggable,
    RequestFormDataCustom,
    RequestFormDataWorkflow
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
</style>
