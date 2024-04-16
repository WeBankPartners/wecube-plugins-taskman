<template>
  <div>
    <Row type="flex">
      <Col span="24" style="padding: 0 20px">
        <div>
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
        </div>
        <div>
          <ApprovalFormNode
            ref="approvalFormNodeRef"
            :isCheck="isCheck"
            @jumpToNode="jumpToNode"
            @reloadParentPage="loadPage"
            @nodeStatus="nodeStatus"
            @setFormConfigStatus="changeFormConfigStatus"
          ></ApprovalFormNode>
        </div>
        <template v-if="isShowFormConfig">
          <div class="title" style="font-size: 16px;">
            <div class="title-text">
              {{ $t('tw_form_configuration') }}
              <span class="underline"></span>
            </div>
          </div>
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
                    <div
                      class="list-group-item-"
                      style="width: 100%"
                      v-for="element in customElement"
                      :key="element.id"
                    >
                      <Input v-if="element.elementType === 'input'" :placeholder="$t('t_input')" />
                      <Input v-if="element.elementType === 'textarea'" type="textarea" :placeholder="$t('textare')" />
                      <Select v-if="element.elementType === 'select'" :placeholder="$t('select')"></Select>
                      <Select
                        v-if="element.elementType === 'wecmdbEntity'"
                        :placeholder="$t('tw_entity_data_items')"
                      ></Select>
                      <DatePicker
                        v-if="element.elementType === 'datePicker'"
                        :type="element.type"
                        :placeholder="$t('tw_date_picker')"
                        style="width:100%"
                      ></DatePicker>
                      <div
                        v-if="element.elementType === 'group'"
                        style="width: 100%; height: 80px; border: 1px solid #5ea7f4"
                      >
                        <span style="margin: 8px; color: #bbbbbb"> Item Group </span>
                      </div>
                    </div>
                  </draggable>
                </div>
              </Col>
              <Col span="14" style="border: 1px solid #dcdee2; padding: 0 16px; width: 57%; margin: 0 4px">
                <div :style="{ height: MODALHEIGHT + 'px', overflow: 'auto' }">
                  <Divider>{{ $t('tw_preview') }}</Divider>
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
                  <div style="min-height: 200px;">
                    <template v-if="finalElement.length === 1 && finalElement[0].itemGroup !== ''">
                      <div
                        v-for="(item, itemIndex) in finalElement"
                        :key="itemIndex"
                        style="border: 2px dotted #A2EF4D; margin: 8px 0; padding: 8px;min-height: 48px;"
                      >
                        <div :key="itemIndex"></div>
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
                              <Icon
                                v-if="element.required === 'yes'"
                                size="8"
                                style="color: #ed4014"
                                type="ios-medical"
                              />
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
                              icon="ios-close"
                              ghost
                            ></Button>
                          </div>
                        </draggable>
                      </div>
                      <div style="text-align: right;">
                        <Button
                          v-if="isCheck !== 'Y'"
                          type="primary"
                          size="small"
                          ghost
                          @click="saveGroup(9, activeEditingNode)"
                          >{{ $t('save') }}</Button
                        >
                        <Button size="small" @click="reloadGroup">{{ $t('tw_restore') }}</Button>
                      </div>
                    </template>
                  </div>
                  <div class="title">
                    <div class="title-text">
                      {{ $t('tw_approval_result') }}
                      <span class="underline"></span>
                    </div>
                  </div>
                  <div style="margin: 12px 0 0 8px;">
                    <Form :label-width="90">
                      <FormItem :label="$t('t_action')">
                        <Select style="width:94%">
                          <Option value="1">{{ $t('tw_approve') }}</Option>
                          <Option value="2">{{ $t('tw_reject') }}</Option>
                          <Option value="3">{{ $t('tw_send_back') }}</Option>
                        </Select>
                      </FormItem>
                      <FormItem :label="$t('tw_comments')">
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
                            <Input v-model="editElement.dataOptions" disabled style="width:70%"></Input>
                            <Button
                              class="custom-add-btn"
                              :disabled="$parent.isCheck === 'Y'"
                              @click.stop="dataOptionsMgmt"
                              type="primary"
                              ghost
                              size="small"
                              icon="ios-create-outline"
                            ></Button>
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
                        <Form :label-width="80">
                          <FormItem :label="$t('validation_rules')">
                            <Input
                              v-model="editElement.regular"
                              :disabled="$parent.isCheck === 'Y'"
                              :placeholder="$t('only_supports_regular')"
                              @on-change="paramsChanged"
                            ></Input>
                          </FormItem>
                          <FormItem :label="$t('tw_data_filtering')">
                            <a-select
                              mode="multiple"
                              v-model="editElement.filterRule"
                              @dropdownVisibleChange="getFilterRuleOption(editElement)"
                              :dropdownMatchSelectWidth="false"
                            >
                              <a-select-option v-for="item in filterRuleOption" :value="item.value" :key="item.value">
                                <Tooltip placement="left">
                                  <p slot="content" style="white-space: normal;">
                                    {{ item.label }}
                                  </p>
                                  {{ item.label }}
                                </Tooltip>
                              </a-select-option>
                            </a-select>
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
          </div>
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
      <DataSourceConfig ref="dataSourceConfigRef" @setDataOptions="setDataOptions"></DataSourceConfig>
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
  saveRequestGroupCustomForm,
  getTargetOptions
} from '@/api/server.js'
export default {
  name: 'BasicInfo',
  data () {
    return {
      isParmasChanged: false, // 参数变化标志位，控制右侧panel显示逻辑
      MODALHEIGHT: 200,
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
          dataOptions: '[]',
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
          dataOptions: '[]',
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
          dataOptions: '[]',
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
          dataOptions: '[]',
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
          dataOptions: '[]',
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
        refPackageName: ''
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
      filterRuleOption: [] // 缓存数据过滤选项
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
    }
  },
  props: ['isCheck', 'requestTemplateId'],
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 500
    this.removeEmptyDataForm()
    this.loadPage()
  },
  methods: {
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
    async selectElement (itemIndex, eleIndex) {
      this.finalElement[itemIndex].attrs.forEach(item => {
        item.isActive = false
      })
      this.finalElement[itemIndex].attrs[eleIndex].isActive = true
      this.editElement = this.finalElement[itemIndex].attrs[eleIndex]
      if (!Array.isArray(this.editElement.filterRule)) {
        this.$set(this.editElement, 'filterRule', JSON.parse(this.editElement.filterRule || '[]'))
      }
      await this.getFilterRuleOption(this.editElement)
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
        if (Array.isArray(item.filterRule)) {
          item.filterRule = JSON.stringify(item.filterRule)
        }
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
      this.editElement.dataOptions = JSON.stringify(options)
      const valueArray = options.map(d => d.value)
      this.editElement.filterRule = this.editElement.filterRule.filter(fr => valueArray.includes(fr))
    },
    computedOption (element) {
      let res = []
      if (element.elementType === 'select') {
        res = JSON.parse(element.dataOptions || '[]')
      } else if (element.elementType === 'wecmdbEntity') {
      }
      return res
    },
    async getFilterRuleOption (element) {
      if (element.elementType === 'select') {
        this.filterRuleOption = JSON.parse(element.dataOptions || '[]')
      } else if (element.elementType === 'wecmdbEntity') {
        if (element.dataOptions !== '' && element.dataOptions.split(':').length === 2) {
          const { status, data } = await getTargetOptions(
            element.dataOptions.split(':')[0],
            element.dataOptions.split(':')[1]
          )
          if (status === 'OK') {
            this.filterRuleOption = data.map(d => {
              return {
                label: d.displayName,
                value: d.id
              }
            })
          }
        }
      }
    }
    // #endregion
  },
  components: {
    ApprovalFormNode,
    RequestFormDataCustom,
    DataSourceConfig,
    RequestFormDataWorkflow,
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
  display: inline-block;
  margin: 8px 0;
}
.custom-add-btn {
  // margin: 0 24px;
}
.dash-line {
  display: inline-block;
  width: 24px;
  vertical-align: middle;
  margin: 0 4px;
  border-color: #dcdee2;
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
</style>
