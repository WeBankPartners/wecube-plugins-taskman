<template>
  <div ref="maxheight">
    <Row type="flex">
      <Col span="24">
        <span
          v-for="(approval, approvalIndex) in approvalNodes"
          :key="approval.id"
          style="margin-right:6px;line-height: 40px;"
        >
          <div
            :class="approval.id === activeEditingNode.id ? 'node-active' : 'node-normal'"
            @click="editNode(approval, true)"
          >
            <span>{{ approval.name }}</span>
            <Icon
              v-if="isCheck !== 'Y' && approvalNodes.length > 1"
              @click.stop="removeNode(approval)"
              type="md-close"
              color="#ed4014"
              style="vertical-align: sub;margin-left: 4px;"
              :size="18"
            />
          </div>
          <div v-if="isCheck !== 'Y'" class="ivu-divider-dashed dash-line"></div>
          <Button
            v-if="isCheck !== 'Y'"
            class="custom-add-btn"
            @click.stop="addApprovalNode(approval.sort + 1)"
            type="success"
            size="small"
            shape="circle"
            icon="md-add"
          ></Button>
          <div v-if="approvalIndex !== approvalNodes.length - 1" class="ivu-divider-dashed dash-line"></div>
        </span>
        <ApprovalFormNode
          ref="approvalFormNodeRef"
          :isCheck="isCheck"
          @jumpToNode="jumpToNode"
          @reloadParentPage="loadPage"
          @nodeStatus="nodeStatus"
          @setFormConfigStatus="changeFormConfigStatus"
          @dataFormFilterChange="setIsEditDisabled"
        ></ApprovalFormNode>
        <template v-if="isShowFormConfig">
          <div class="title" style="font-size: 16px;">
            <div class="title-text">
              {{ $t('tw_form_configuration') }}
              <span class="underline"></span>
            </div>
          </div>
          <Row style="margin-bottom: 56px;">
            <Col span="5" style="border: 1px solid #dcdee2;">
              <div :style="{ height: MODALHEIGHT + 'px', overflow: 'auto', padding: '0 8px' }">
                <!--自定义表单项-->
                <Divider orientation="left" size="small">{{ $t('custom_form') }}</Divider>
                <CustomDraggable :sortable="$parent.isCheck !== 'Y'" :clone="cloneDog"></CustomDraggable>
                <!--表单项组件库-->
                <template v-if="finalElement[0].itemGroupName">
                  <Divider orientation="left" size="small">{{ $t('tw_template_library') }}</Divider>
                  <ComponentLibraryList
                    ref="libraryList"
                    :formType="finalElement[0].itemGroupName"
                    :groupType="nextGroupInfo.itemGroupType"
                  ></ComponentLibraryList>
                </template>
              </div>
            </Col>
            <!--表单预览-->
            <Col span="14" style="border: 1px solid #dcdee2; padding: 0 16px; width: 57%; margin: 0 4px;">
              <div :style="{ height: MODALHEIGHT + 'px', overflow: 'auto', paddingBottom: '10px' }">
                <div class="title">
                  <div class="title-text">
                    {{ $t('tw_approval_content') }}
                    <span class="underline"></span>
                  </div>
                </div>
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
                <div style="min-height:300px;">
                  <template v-if="finalElement.length === 1 && finalElement[0].itemGroup !== ''">
                    <template v-for="(item, itemIndex) in finalElement">
                      <div
                        :key="itemIndex"
                        style="border: 2px dotted #A2EF4D; margin: 8px 0; padding: 8px;min-height: 48px;"
                      >
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
                              v-if="['calculate', 'textarea'].includes(element.elementType)"
                              :disabled="element.isEdit === 'no'"
                              v-model="element.defaultValue"
                              type="textarea"
                              class="custom-item"
                            />
                            <Select
                              v-if="element.elementType === 'select'"
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
                            >
                            </Select>
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
                      <div
                        v-if="isCheck !== 'Y'"
                        :key="itemIndex + '-'"
                        style="display:flex;justify-content:space-between;"
                      >
                        <Button
                          :disabled="getAddComponentDisabled(item.attrs)"
                          @click="createComponentLibrary"
                          size="small"
                          >新建组件</Button
                        >
                        <div>
                          <Button type="primary" size="small" ghost @click="saveGroup(9, activeEditingNode)">{{
                            $t('save')
                          }}</Button>
                          <Button size="small" @click="reloadGroup">{{ $t('tw_restore') }}</Button>
                        </div>
                      </div>
                    </template>
                  </template>
                </div>
                <div class="title">
                  <div class="title-text">
                    {{ $t('tw_approval_result') }}
                    <span class="underline"></span>
                  </div>
                </div>
                <div style="margin: 8px 0 0 8px;">
                  <Form :label-width="90">
                    <FormItem :label="$t('t_action')">
                      <Select style="width:94%">
                        <Option value="1">{{ $t('tw_approve') }}</Option>
                        <Option value="2">{{ $t('tw_reject') }}</Option>
                        <Option value="3">{{ $t('tw_send_back') }}</Option>
                        <Option value="4">{{ $t('tw_unrelated') }}</Option>
                      </Select>
                    </FormItem>
                    <FormItem :label="$t('process_comments')">
                      <Input type="textarea" :rows="2" style="width:94%"></Input>
                    </FormItem>
                  </Form>
                </div>
              </div>
              <Modal v-model="showSelectModel" :title="$t('tw_reference_global_forms')" :mask-closable="false">
                <Form :label-width="120">
                  <FormItem :label="$t('tw_global_forms')">
                    <Select style="width: 80%" v-model="itemGroup" filterable>
                      <Option v-for="item in groupOptions" :value="item.id" :key="item.id">{{
                        item.itemGroupName
                      }}</Option>
                    </Select>
                  </FormItem>
                </Form>
                <template #footer>
                  <Button @click="showSelectModel = false">{{ $t('cancel') }}</Button>
                  <Button @click="okSelect" :disabled="itemGroup === ''" type="primary">{{ $t('confirm') }}</Button>
                </template>
              </Modal>
            </Col>
            <!--属性设置-->
            <Col span="5" style="border: 1px solid #dcdee2">
              <div :style="{ height: MODALHEIGHT + 'px', overflow: 'auto' }">
                <Collapse v-model="openPanel">
                  <Panel name="1">
                    {{ $t('general_attributes') }}
                    <div slot="content">
                      <Form :label-width="80" ref="attrForm" :model="editElement" :rules="ruleForm">
                        <FormItem :label="$t('display_name')">
                          <Input
                            v-model="editElement.title"
                            @on-change="paramsChanged"
                            :disabled="$parent.isCheck === 'Y' || Boolean(editElement.copyId)"
                            placeholder=""
                          ></Input>
                        </FormItem>
                        <FormItem :label="$t('tw_code')" prop="name" style="margin-bottom:20px;">
                          <Input
                            v-model="editElement.name"
                            @on-change="paramsChanged"
                            :disabled="
                              $parent.isCheck === 'Y' || editElement.entity !== '' || Boolean(editElement.copyId)
                            "
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
                            @click.stop="dataOptionsMgmt"
                            :disabled="$parent.isCheck === 'Y' || Boolean(editElement.copyId)"
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
                            :disabled="$parent.isCheck === 'Y' || editElement.isEditDisabled"
                            @on-change="paramsChanged"
                            size="default"
                          />
                        </FormItem>
                        <FormItem :label="$t('required')">
                          <i-switch
                            v-model="editElement.required"
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
                          <Select
                            v-model="editElement.width"
                            @on-change="paramsChanged"
                            :disabled="$parent.isCheck === 'Y'"
                          >
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
                      <Form :label-width="80" label-position="left">
                        <FormItem :label="$t('validation_rules')">
                          <Input
                            v-model="editElement.regular"
                            :disabled="$parent.isCheck === 'Y'"
                            :placeholder="$t('only_supports_regular')"
                            @on-change="paramsChanged"
                          ></Input>
                        </FormItem>
                        <FormItem label="" :label-width="0">
                          <HiddenCondition
                            :disabled="$parent.isCheck === 'Y'"
                            :finalElement="finalElement"
                            v-model="editElement.hiddenCondition"
                            :name="editElement.name"
                          ></HiddenCondition>
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
        </template>
      </Col>
      <!-- 自定义表单配置 -->
      <RequestFormDataCustom
        ref="requestFormDataCustomRef"
        @reloadParentPage="reloadGroup"
        :isCheck="isCheck"
        module="other"
        v-show="['custom'].includes(itemGroupType)"
      ></RequestFormDataCustom>
      <!-- 编排表单配置 -->
      <RequestFormDataWorkflow
        ref="requestFormDataWorkflowRef"
        :isCustomItemEditable="false"
        @reloadParentPage="reloadGroup"
        :isCheck="isCheck"
        module="other"
        v-show="['workflow', 'optional'].includes(itemGroupType)"
      ></RequestFormDataWorkflow>
      <!--数据集弹框-->
      <DataSourceConfig ref="dataSourceConfigRef" @setDataOptions="setDataOptions"></DataSourceConfig>
      <!--组件库弹框-->
      <ComponentLibraryModal
        ref="library"
        :checkedList="componentCheckedList"
        :formType="finalElement[0].itemGroupName"
        :groupType="nextGroupInfo.itemGroupType"
        v-model="componentVisible"
        @fetchList="$refs.libraryList.handleSearch()"
      />
    </Row>
    <div class="footer">
      <div class="content" :style="isShowFormConfig ? '' : 'margin-top:48px'">
        <Button
          :disabled="isCheck !== 'Y' && isTopButtonDisable"
          @click="gotoForward"
          ghost
          type="primary"
          class="btn-footer-margin"
          >{{ $t('forward') }}</Button
        >
        <Button
          :disabled="isCheck !== 'Y' && isTopButtonDisable"
          v-if="isCheck !== 'Y'"
          @click="saveApprovalFromNode"
          type="info"
          class="btn-footer-margin"
          >{{ $t('save') }}</Button
        >
        <Button
          :disabled="isCheck !== 'Y' && isTopButtonDisable"
          @click="gotoNext"
          type="primary"
          class="btn-footer-margin"
          >{{ $t('next') }}</Button
        >
      </div>
    </div>
  </div>
</template>

<script>
import draggable from 'vuedraggable'
import ApprovalFormNode from './approval-form-node.vue'
import DataSourceConfig from './data-source-config.vue'
import RequestFormDataCustom from './request-form-data-custom.vue'
import RequestFormDataWorkflow from './request-form-data-workflow.vue'
import CustomDraggable from './components/custom-draggable.vue'
import {
  getApprovalNode,
  addApprovalNode,
  removeApprovalNode,
  getApprovalGlobalForm,
  copyItemGroup,
  removeEmptyDataForm,
  getApprovalNodeGroups,
  deleteRequestGroupForm,
  getAllDataModels,
  saveRequestGroupCustomForm
} from '@/api/server.js'
import ComponentLibraryModal from './components/component-library-modal.vue'
import ComponentLibraryList from './components/component-library-list.vue'
import HiddenCondition from './components/hidden-condition.vue'
import { uniqueArr, deepClone, findFirstDuplicateIndex } from '@/pages/util'
export default {
  components: {
    ApprovalFormNode,
    RequestFormDataCustom,
    RequestFormDataWorkflow,
    draggable,
    CustomDraggable,
    DataSourceConfig,
    ComponentLibraryModal,
    ComponentLibraryList,
    HiddenCondition
  },
  data () {
    return {
      isParmasChanged: false, // 参数变化标志位，控制右侧panel显示逻辑
      MODALHEIGHT: 400,
      isTopButtonDisable: true, // 下一步，上一步等的控制
      approvalNodes: [
        {
          id: '',
          sort: 1,
          requestTemplate: '',
          name: `${this.$t('tw_approval')}1`,
          expireDay: 1,
          description: '',
          roleType: '',
          roleObjs: []
        }
      ],
      activeEditingNode: {}, // 标记整在编辑的节点
      isShowFormConfig: false,
      dataFormInfo: {
        associationWorkflow: false,
        formTemplateId: '',
        groups: []
      },
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
        hiddenCondition: [] // 隐藏条件
      },
      activeTag: {
        // 正在编辑的组
        itemGroupIndex: -1,
        attrIndex: -1
      },
      allEntityList: [],
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
      groups: [], // 已存在组
      showSelectModel: false, // 组选择框
      itemGroupType: '', // 选中的组类型
      itemGroup: '', // 选中的组信息
      nextNodeInfo: {}, // 缓存待切换节点信息
      displayLastGroup: false, // 控制group显示，在新增时显示最后一个，其余显示当前值
      nextGroupInfo: {},
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
  props: ['isCheck', 'requestTemplateId'],
  mounted () {
    const clientHeight = document.documentElement.clientHeight
    this.MODALHEIGHT = clientHeight - this.$refs.maxheight.getBoundingClientRect().top - 90
    this.removeEmptyDataForm()
    this.loadPage()
  },
  methods: {
    // 数据表单过滤项有值，需要禁用审批表单对应表单项"可编辑"属性
    setIsEditDisabled () {
      if (this.$refs.approvalFormNodeRef.filterFormList && this.$refs.approvalFormNodeRef.filterFormList.length > 0) {
        const dataFormObj = this.$refs.approvalFormNodeRef.filterFormList.find(i => i.type === 2)
        const handleTemplates = this.$refs.approvalFormNodeRef.activeApprovalNode.handleTemplates
        dataFormObj.items.forEach(j => {
          const key = `${dataFormObj.itemGroup}-${j.name}`
          const hasValue = handleTemplates.some(val => {
            return Array.isArray(val.filterRule[key]) ? val.filterRule[key].length > 0 : Boolean(val.filterRule[key])
          })
          // 遍历所有表单组，找到需要禁用的表单项
          this.dataFormInfo.groups.forEach(group => {
            group.items.forEach(item => {
              const key1 = `${group.itemGroup}-${item.name}`
              if (key1 === key && hasValue) {
                item.isEdit = 'no'
                this.$set(item, 'isEditDisabled', true)
              }
              if (key1 === key && !hasValue) {
                this.$set(item, 'isEditDisabled', false)
              }
            })
          })
        })
      }
    },
    async removeEmptyDataForm () {
      await removeEmptyDataForm(this.requestTemplateId)
    },
    loadPage (id = '') {
      this.isParmasChanged = false
      this.getApprovalNode(id)
      this.getAllDataModels()
    },
    async getApprovalNode (id = '') {
      const { statusCode, data } = await getApprovalNode(this.requestTemplateId, 'approve')
      if (statusCode === 'OK') {
        if (!id && data.ids.length === 0) {
          this.addApprovalNode(1)
        } else {
          this.approvalNodes = data.ids
          this.activeEditingNode = id === '' ? this.approvalNodes[0] : this.approvalNodes.find(node => node.id === id)
          this.editNode(this.activeEditingNode, false)
        }
      }
    },
    async addApprovalNode (sort) {
      const params = {
        type: 'approve',
        requestTemplate: this.requestTemplateId,
        name: this.$t('tw_approval') + sort,
        expireDay: 1,
        sort: sort
      }
      if (sort === 1) {
        // 新增第一个节点
        const { statusCode, data } = await addApprovalNode(params)
        if (statusCode === 'OK') {
          this.getApprovalNode(data.taskTemplate.id)
        }
      } else {
        const nodeStatus = this.$refs.approvalFormNodeRef.panalStatus()
        if (nodeStatus === 'canSave') {
          this.$refs.approvalFormNodeRef.saveNode(3)
          this.saveGroup(3, this.activeEditingNode)
          const { statusCode, data } = await addApprovalNode(params)
          if (statusCode === 'OK') {
            this.getApprovalNode(data.taskTemplate.id)
          }
        }
      }
    },
    isLoadLastGroup (val) {
      this.displayLastGroup = val
      // this.loadPage()
    },
    async removeNode (node) {
      this.$Modal.confirm({
        title: this.$t('confirm_delete'),
        'z-index': 1000000,
        loading: true,
        okText: this.$t('tw_request_confirm'),
        onOk: async () => {
          this.$Modal.remove()
          const { statusCode } = await removeApprovalNode(this.requestTemplateId, node.id)
          if (statusCode === 'OK') {
            if (this.activeEditingNode.id === node.id) {
              this.getApprovalNode()
            } else {
              this.getApprovalNode(this.activeEditingNode.id)
            }
          }
        },
        onCancel: () => {}
      })
    },
    // 在弹窗关闭、保存、还原状态下回显group内容
    reloadGroup () {
      this.isParmasChanged = false
      this.getApprovalNodeGroups(this.activeEditingNode)
    },
    editNode (node, isNeedSaveFirst = true) {
      if (isNeedSaveFirst && this.isCheck !== 'Y') {
        const nodeStatus = this.$refs.approvalFormNodeRef.panalStatus()
        if (nodeStatus === 'canSave') {
          this.$refs.approvalFormNodeRef.saveNode(3)
          this.beforeEditNode(node)
        }
      } else {
        this.activeEditingNode = node
        let params = {
          requestTemplateId: this.requestTemplateId,
          id: node.id
        }
        this.$refs.approvalFormNodeRef.loadPage(params)
        this.getApprovalNodeGroups(node)
      }
    },
    beforeEditNode (node) {
      this.saveGroup(3, node)
    },
    async getApprovalNodeGroups (node) {
      const { statusCode, data } = await getApprovalNodeGroups(this.requestTemplateId, node.id)
      if (statusCode === 'OK') {
        this.dataFormInfo = data
        this.setIsEditDisabled()
        let groups = this.dataFormInfo.groups
        if (groups.length !== 0) {
          if (this.displayLastGroup) {
            const group = groups[groups.length - 1]
            this.editGroupCustomItems(group, false)
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
        } else {
          this.cancelGroup()
        }
      }
    },
    jumpToNode () {
      this.cancelGroup()
      let params = {
        requestTemplateId: this.requestTemplateId,
        id: this.nextNodeInfo.id
      }
      this.activeEditingNode = this.nextNodeInfo
      this.$refs.approvalFormNodeRef.loadPage(params)
    },
    beforeSelectItemGroup () {
      if (this.finalElement[0].itemGroupId === '') {
        this.selectItemGroup()
      } else {
        this.saveGroup(7)
      }
    },
    // 查询可添加的组
    async selectItemGroup () {
      this.itemGroup = ''
      this.groupOptions = []
      const { statusCode, data } = await getApprovalGlobalForm(this.requestTemplateId)
      if (statusCode === 'OK') {
        // workflow  2.编排数据, 3.optional 自选数据项表单,  custom 1.自定义表单
        this.groupOptions = data.filter(d => {
          const findIndex = this.dataFormInfo.groups.findIndex(group => group.itemGroupName === d.itemGroupName)
          if (findIndex === -1) {
            return d
          }
        })
        this.showSelectModel = true
      }
    },
    // 选择可添加的组
    async okSelect () {
      this.showSelectModel = false
      let params = {
        requestTemplateId: this.requestTemplateId,
        taskTemplateId: this.activeEditingNode.id || '',
        itemGroupId: this.itemGroup
      }
      const { statusCode } = await copyItemGroup(params)
      if (statusCode === 'OK') {
        this.displayLastGroup = true
        this.getApprovalNodeGroups(this.activeEditingNode)
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
        const message = titleArr.join(',')
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
    paramsChanged () {
      this.isParmasChanged = true
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
    },
    updateFinalElement (groupItem) {
      this.finalElement = [
        {
          itemGroupId: groupItem.itemGroupId,
          formTemplateId: this.dataFormInfo.formTemplateId,
          requestTemplateId: this.requestTemplateId,
          itemGroup: groupItem.itemGroup,
          itemGroupName: groupItem.itemGroupName,
          attrs: groupItem.items || []
        }
      ]
      this.openPanel = ''
    },
    // 编辑组弹出信息
    editGroupItem (groupItem) {
      if (this.isCheck === 'Y') {
        this.openDrawer(groupItem)
      } else {
        this.saveGroup(5, groupItem)
      }
    },
    openDrawer (groupItem) {
      this.editGroupCustomItems(groupItem)
      // 隐藏自定义表单项配置
      // this.cancelGroup()
      if (groupItem.itemGroupType === 'custom') {
        this.itemGroupType = groupItem.itemGroupType
        let params = {
          requestTemplateId: this.requestTemplateId,
          taskTemplateId: this.activeEditingNode.id || '',
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
          taskTemplateId: this.activeEditingNode.id || '',
          isAdd: false,
          itemGroupName: groupItem.itemGroupName,
          itemGroupType: groupItem.itemGroupType,
          itemGroupId: groupItem.itemGroupId,
          itemGroup: groupItem.itemGroup
        }
        this.$refs.requestFormDataWorkflowRef.loadPage(params)
      }
    },
    async removeGroupItem (groupItem) {
      this.nextGroupInfo = groupItem
      this.saveGroup(6)
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
            this.getApprovalNodeGroups(this.activeEditingNode)
          }
        },
        onCancel: () => {}
      })
    },
    // // 删除组
    // async removeGroupItem (groupItem) {
    //   this.$Modal.confirm({
    //     title: this.$t('confirm_delete'),
    //     'z-index': 1000000,
    //     loading: true,
    //     onOk: async () => {
    //       this.$Modal.remove()
    //       const { statusCode } = await deleteRequestGroupForm(groupItem.itemGroupId, this.requestTemplateId)
    //       if (statusCode === 'OK') {
    //         this.$Notice.success({
    //           title: this.$t('successful'),
    //           desc: this.$t('successful')
    //         })
    //         this.loadPage()
    //       }
    //     },
    //     onCancel: () => {}
    //   })
    // },
    // showDeleteTip (groupItem) {
    //   this.$Modal.confirm({
    //     title: this.$t('confirm_delete'),
    //     'z-index': 1000000,
    //     loading: true,
    //     onOk: async () => {
    //       this.$Modal.remove()
    //       const { statusCode } = await deleteRequestGroupForm(groupItem.itemGroupId, this.requestTemplateId)
    //       if (statusCode === 'OK') {
    //         this.$Notice.success({
    //           title: this.$t('successful'),
    //           desc: this.$t('successful')
    //         })
    //         this.loadPage()
    //       }
    //     },
    //     onCancel: () => {}
    //   })
    // },
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
    // 选中自定义表单项
    async selectElement (itemIndex, eleIndex) {
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
      this.$refs.attrForm.validateField('name')
    },
    // 删除自定义表单项
    removeForm (itemIndex, eleIndex, element) {
      this.finalElement[itemIndex].attrs.splice(eleIndex, 1)
      this.openPanel = ''
      this.paramsChanged()
    },
    // 保存自定义表单项
    async saveGroup (nextStep, elememt) {
      // nextStep 1新增 2下一步 3切换tab 4 切换到目标group 5切换到目标group打开弹窗 6删除组 7选择组 8上一步
      let finalData = JSON.parse(JSON.stringify(this.finalElement[0]))
      if (finalData.itemGroupId === '') {
        this.loadPage(elememt.id)
        return
      }
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
          this.loadPage()
        } else if (nextStep === 2) {
          this.$emit('gotoStep', this.requestTemplateId, 'forward')
        } else if ([3].includes(nextStep)) {
          if (elememt.id) {
            this.loadPage(elememt.id)
          }
        } else if ([4].includes(nextStep)) {
          // this.activeEditingNode = elememt
          this.updateFinalElement(elememt)
          this.getApprovalNodeGroups(this.activeEditingNode)
        } else if ([9].includes(nextStep)) {
          // this.updateFinalElement(elememt)
          this.getApprovalNodeGroups(this.activeEditingNode)
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
    cancelGroup () {
      this.finalElement = [
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
      this.openPanel = ''
    },
    async gotoNext () {
      if (this.isCheck === 'Y') {
        this.$emit('gotoStep', this.requestTemplateId, 'forward')
        return
      }
      const nodeStatus = this.$refs.approvalFormNodeRef.panalStatus()
      if (nodeStatus === 'canSave') {
        await this.$refs.approvalFormNodeRef.saveNode(3)
        let finalData = JSON.parse(JSON.stringify(this.finalElement[0]))
        if (finalData.itemGroupId === '') {
          this.$emit('gotoStep', this.requestTemplateId, 'forward')
        } else {
          this.saveGroup(2, {})
        }
      }
    },
    gotoForward () {
      if (this.isCheck === 'Y') {
        this.$emit('gotoStep', this.requestTemplateId, 'backward')
        return
      }
      const nodeStatus = this.$refs.approvalFormNodeRef.panalStatus()
      if (nodeStatus === 'canSave') {
        this.$refs.approvalFormNodeRef.saveNode(3)
        let finalData = JSON.parse(JSON.stringify(this.finalElement[0]))
        if (finalData.itemGroupId === '') {
          this.$emit('gotoStep', this.requestTemplateId, 'backward')
        } else {
          this.saveGroup(8, {})
        }
      }
    },
    saveApprovalFromNode () {
      this.$refs.approvalFormNodeRef.saveNode(1)
      this.saveGroup(9, this.activeEditingNode)
    },
    nodeStatus (status) {
      this.isTopButtonDisable = status
    },
    // 控制配置显示状态
    changeFormConfigStatus (status) {
      this.isShowFormConfig = status
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
    // #endregion
  }
}
</script>
<style>
.ivu-input[disabled],
fieldset[disabled] .ivu-input {
  color: #757575 !important;
}
.ivu-select-input[disabled] {
  color: #757575;
  -webkit-text-fill-color: #757575;
}
.ivu-btn-icon-only.ivu-btn-circle > .ivu-icon {
  vertical-align: middle !important;
}
</style>
<style lang="scss" scoped>
.ivu-form-item {
  margin-bottom: 16px;
}
.active-zone {
  color: #338cf0;
}
.basci-info-right {
  height: calc(100vh - 260px);
}

.basci-info-left {
  @extend .basci-info-right;
  border-right: 1px solid #dcdee2;
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
    height: 12px;
    border-radius: 12px;
    background-color: #c6eafe;
    box-sizing: content-box;
  }
}

.basci-info-content {
  margin: 16px 64px;
}
.list-group-item- {
  display: flex;
  align-items: center;
  justify-content: space-around;
  margin: 8px 0;
}
.dash-line {
  display: inline-block;
  width: 24px;
  vertical-align: middle;
  margin: 0 4px;
  border-color: #dcdee2;
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

.node-normal {
  height: 32px;
  line-height: 32px;
  display: inline-block;
  border: 1px solid #e8eaec;
  border-radius: 3px;
  background: #f7f7f7;
  font-size: 12px;
  vertical-align: middle;
  opacity: 1;
  overflow: hidden;
  padding: 0 12px;
  cursor: pointer;
}
.node-active {
  height: 32px;
  line-height: 32px;
  display: inline-block;
  border-radius: 3px;
  background: #2d8cf0;
  font-size: 12px;
  vertical-align: middle;
  opacity: 1;
  overflow: hidden;
  color: #fff;
  padding: 0 12px;
  cursor: pointer;
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
