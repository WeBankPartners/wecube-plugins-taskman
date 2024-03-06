<template>
  <div>
    <Row type="flex">
      <Col span="24" style="padding: 0 20px">
        <div style="margin-bottom: 12px">
          <span v-for="approval in approvalNodes" :key="approval.id" style="margin-right:6px;">
            <div
              :class="approval.id === activeEditingNode.id ? 'node-active' : 'node-normal'"
              @click="editNode(approval)"
            >
              <span>{{ approval.name }}</span>
              <Icon
                v-if="approvalNodes.length > 1"
                @click.stop="removeNopde(approval)"
                type="md-close"
                color="#ed4014"
                style="vertical-align: sub;"
                :size="18"
              />
            </div>
            <Button
              @click.stop="addApprovalNode(approval.sort + 1)"
              type="success"
              size="small"
              ghost
              shape="circle"
              icon="md-add"
            ></Button>
          </span>
        </div>
        <div>
          <ApprovalFormNode
            ref="approvalFormNodeRef"
            @jumpToNode="jumpToNode"
            @reloadParentPage="loadPage"
          ></ApprovalFormNode>
        </div>
        <Divider />
        <div>
          <Row>
            <Col span="6" style="border: 1px solid #dcdee2; padding: 0 16px">
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
                      type="date"
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
            <Col span="12" style="border: 1px solid #dcdee2; padding: 0 16px; width: 48%; margin: 0 4px">
              <div :style="{ height: MODALHEIGHT + 30 + 'px', overflow: 'auto' }">
                <Divider>预览</Divider>
                <div class="title">
                  <div class="title-text">
                    {{ $t('审批内容') }}
                    <span class="underline"></span>
                  </div>
                </div>
                <div style="margin-top:16px">
                  <div class="radio-group">
                    <div
                      v-for="(groupItem, index) in dataFormInfo.groups"
                      :key="index"
                      :class="{
                        radio: true,
                        custom: groupItem.itemGroupType === 'custom',
                        workflow: groupItem.itemGroupType === 'workflow',
                        optional: groupItem.itemGroupType === 'optional'
                      }"
                      :style="activeStyle(groupItem)"
                    >
                      <Icon @click="editGroupItem(groupItem)" type="md-create" color="#2d8cf0" :size="16" />
                      <span @click="editGroupCustomItems(groupItem)">
                        {{ `${groupItem.itemGroup}` }}
                      </span>
                      <Icon @click="removeGroupItem(groupItem)" type="md-close" color="#ed4014" :size="18" />
                    </div>
                    <span>
                      <Button @click="selectItemGroup" type="primary" ghost icon="md-add"></Button>
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
                              ['calculate', 'textarea'].includes(element.elementType) ? 'vertical-align: top;' : ''
                            "
                          >
                            <Icon
                              v-if="element.required === 'yes'"
                              size="8"
                              style="color: #ed4014"
                              type="ios-medical"
                            />
                            {{ element.title }}:
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
                            type="date"
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
                    <div style="text-align: right;">
                      <Button type="primary" size="small" ghost @click="saveGroup">{{ $t('save') }}</Button>
                      <Button size="small" @click="cancelGroup">{{ $t('tw_abandon') }}</Button>
                    </div>
                  </template>
                </div>
                <div class="title">
                  <div class="title-text">
                    {{ $t('审批结果') }}
                    <span class="underline"></span>
                  </div>
                </div>
                <div style="margin-top: 24px;">
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
              <Modal v-model="showSelectModel" title="选择组信息" :mask-closable="false">
                <Form :label-width="120">
                  <FormItem :label="$t('选择组')">
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
            <Col span="6" style="border: 1px solid #dcdee2">
              <div :style="{ height: MODALHEIGHT + 32 + 'px', overflow: 'auto' }">
                <Collapse v-model="openPanel">
                  <Panel name="1">
                    {{ $t('general_attributes') }}
                    <div slot="content">
                      <Form :label-width="80">
                        <FormItem :label="$t('field_name')">
                          <Input
                            v-model="editElement.name"
                            @on-change="paramsChanged"
                            :disabled="$parent.isCheck === 'Y'"
                            placeholder=""
                          ></Input>
                        </FormItem>
                        <FormItem :label="$t('display_name')">
                          <Input
                            v-model="editElement.title"
                            @on-change="paramsChanged"
                            :disabled="$parent.isCheck === 'Y'"
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
      </Col>
      <!-- 自定义表单配置 -->
      <RequestFormDataCustom
        ref="requestFormDataCustomRef"
        @reloadParentPage="loadPage"
        module="other"
        v-show="['custom'].includes(itemGroupType)"
      ></RequestFormDataCustom>

      <!-- 编排表单配置 -->
      <RequestFormDataWorkflow
        ref="requestFormDataWorkflowRef"
        @reloadParentPage="loadPage"
        module="other"
        v-show="['workflow', 'optional'].includes(itemGroupType)"
      ></RequestFormDataWorkflow>
    </Row>
    <div style="text-align: center;margin-top: 16px;">
      <Button @click="gotoNext" type="primary">{{ $t('next') }}</Button>
    </div>
  </div>
</template>

<script>
import draggable from 'vuedraggable'
import ApprovalFormNode from './approval-form-node.vue'
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
  saveRequestGroupCustomForm
} from '@/api/server.js'
let idGlobal = 118
export default {
  name: 'BasicInfo',
  data () {
    return {
      isParmasChanged: false, // 参数变化标志位，控制右侧panel显示逻辑
      MODALHEIGHT: 200,
      approvalNodes: [
        {
          id: '',
          sort: 1,
          requestTemplate: '',
          name: '审批1',
          expireDay: 1,
          description: '',
          roleType: '',
          roleObjs: []
        }
      ],
      activeEditingNode: {}, // 标记整在编辑的节点
      expireDayOptions: [1, 2, 3, 4, 5, 6, 7], // 时效
      roleTypeOptions: [
        // custom.单人自定义 any.协同 all.并行 admin.提交人角色管理员 auto.自动通过
        { label: '单人:自定义', value: 'custom' },
        { label: '协同:任意审批人处理', value: 'any' },
        { label: '并行:全部审批人处理', value: 'all' },
        { label: '提交人角色管理员', value: 'admin' },
        { label: '自动通过', value: 'auto' }
      ],
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
          abc: 'datePicker',
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
      nextNodeInfo: {} // 缓存待切换节点信息
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
  props: ['requestTemplateId'],
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 500
    this.removeEmptyDataForm()
    this.loadPage()
  },
  methods: {
    async removeEmptyDataForm () {
      await removeEmptyDataForm(this.requestTemplateId)
    },
    loadPage (sort) {
      this.cancelGroup()
      this.getApprovalNode(sort)
      this.getAllDataModels()
    },
    async getApprovalNode (sort) {
      const { statusCode, data } = await getApprovalNode(this.requestTemplateId, 'approve')
      if (statusCode === 'OK') {
        if (data.ids.length === 0) {
          this.addApprovalNode(1)
        } else {
          this.approvalNodes = data.ids
          const tmpSort = sort === undefined ? 1 : sort
          this.activeEditingNode = this.approvalNodes.find(node => node.sort === tmpSort)
          this.editNode(this.activeEditingNode)
        }
      }
    },
    async addApprovalNode (sort) {
      this.nextNodeInfo = JSON.parse(JSON.stringify(this.activeEditingNode))
      const res = this.preApprovalNodeChange()
      if (res) {
        return
      }
      const params = {
        type: 'approve',
        requestTemplate: this.requestTemplateId,
        name: this.$t('tw_approval') + sort,
        expireDay: 1,
        sort: sort
      }
      const { statusCode } = await addApprovalNode(params)
      if (statusCode === 'OK') {
        this.getApprovalNode(sort - 1)
      }
    },
    async removeNopde (node) {
      this.nextNodeInfo = JSON.parse(JSON.stringify(this.activeEditingNode))
      const res = this.preApprovalNodeChange()
      if (!res) {
        this.$Modal.confirm({
          title: this.$t('confirm_delete'),
          'z-index': 1000000,
          loading: true,
          onOk: async () => {
            this.$Modal.remove()
            const { statusCode } = await removeApprovalNode(this.requestTemplateId, node.id)
            if (statusCode === 'OK') {
              this.getApprovalNode()
            }
          },
          onCancel: () => {}
        })
      }
    },
    preApprovalNodeChange () {
      if (this.isParmasChanged) {
        this.$Modal.confirm({
          title: `${this.$t('confirm_discarding_changes')}`,
          content: `${this.finalElement[0].itemGroupName}:${this.$t('params_edit_confirm')}`,
          'z-index': 1000000,
          okText: this.$t('save'),
          cancelText: this.$t('abandon'),
          onOk: async () => {
            this.saveGroup()
          },
          onCancel: () => {
            this.cancelGroup()
          }
        })
        return true
      }
      if (this.$refs.approvalFormNodeRef.panalStatus()) {
        this.$refs.approvalFormNodeRef.isNeedConfirm()
        return true
      }
      return false
    },
    editNode (node) {
      this.cancelGroup()
      let params = {
        requestTemplateId: this.requestTemplateId,
        id: node.id
      }
      this.nextNodeInfo = node
      const res = this.preApprovalNodeChange()
      if (!res) {
        this.activeEditingNode = node
        this.$refs.approvalFormNodeRef.loadPage(params)
      }
      this.getApprovalNodeGroups(node)
    },
    async getApprovalNodeGroups (node) {
      const { statusCode, data } = await getApprovalNodeGroups(this.requestTemplateId, node.id)
      if (statusCode === 'OK') {
        this.dataFormInfo = data
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
        taskTemplateId: this.activeEditingNode.id,
        itemGroupId: this.itemGroup
      }
      const { statusCode } = await copyItemGroup(params)
      if (statusCode === 'OK') {
        this.getApprovalNodeGroups(this.activeEditingNode)
      }
    },
    cloneDog (val) {
      if (this.$parent.isCheck === 'Y') return
      let newItem = JSON.parse(JSON.stringify(val))
      newItem.id = 'c_' + idGlobal++
      newItem.title = newItem.title + idGlobal
      newItem.name = newItem.name + idGlobal
      newItem.isActive = true
      this.specialId = newItem.id
      this.paramsChanged()
      this.finalElement[0].attrs.forEach(item => {
        item.isActive = false
      })
      return newItem
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
    editGroupCustomItems (groupItem) {
      if (this.isParmasChanged) {
        this.$Modal.confirm({
          title: `${this.$t('confirm_discarding_changes')}`,
          content: `${this.finalElement[0].itemGroupName}:${this.$t('params_edit_confirm')}`,
          'z-index': 1000000,
          okText: this.$t('save'),
          cancelText: this.$t('abandon'),
          onOk: async () => {
            this.saveGroup()
          },
          onCancel: () => {
            this.getApprovalNodeGroups(this.activeEditingNode)
            this.cancelGroup()
            this.updateFinalElement(groupItem)
          }
        })
      } else {
        this.cancelGroup()
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
    },
    // 编辑组弹出信息
    editGroupItem (groupItem) {
      if (this.isParmasChanged) {
        this.$Modal.confirm({
          title: `${this.$t('confirm_discarding_changes')}`,
          content: `${this.finalElement[0].itemGroupName}:${this.$t('params_edit_confirm')}`,
          'z-index': 1000000,
          okText: this.$t('save'),
          cancelText: this.$t('abandon'),
          onOk: async () => {
            this.saveGroup()
          },
          onCancel: () => {
            this.getApprovalNodeGroups(this.activeEditingNode)
            this.openDrawer(groupItem)
          }
        })
      } else {
        this.openDrawer(groupItem)
      }
    },
    openDrawer (groupItem) {
      // 隐藏自定义表单项配置
      this.cancelGroup()
      if (groupItem.itemGroupType === 'custom') {
        this.itemGroupType = groupItem.itemGroupType
        let params = {
          requestTemplateId: this.requestTemplateId,
          formTemplateId: this.dataFormInfo.formTemplateId,
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
          formTemplateId: this.dataFormInfo.formTemplateId,
          isAdd: false,
          itemGroupName: groupItem.itemGroupName,
          itemGroupType: groupItem.itemGroupType,
          itemGroupId: groupItem.itemGroupId,
          itemGroup: groupItem.itemGroup
        }
        this.$refs.requestFormDataWorkflowRef.loadPage(params)
      }
    },
    // 删除组
    async removeGroupItem (groupItem) {
      if (this.isParmasChanged) {
        this.$Modal.confirm({
          title: `${this.$t('confirm_discarding_changes')}`,
          content: `${this.finalElement[0].itemGroupName}:${this.$t('params_edit_confirm')}`,
          'z-index': 1000000,
          okText: this.$t('save'),
          cancelText: this.$t('abandon'),
          onOk: async () => {
            this.saveGroup()
          },
          onCancel: () => {
            this.showDeleteTip(groupItem)
          }
        })
      } else {
        this.showDeleteTip(groupItem)
      }
    },
    showDeleteTip (groupItem) {
      this.$Modal.confirm({
        title: this.$t('confirm_delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const { statusCode } = await deleteRequestGroupForm(groupItem.itemGroupId, this.activeEditingNode.id)
          if (statusCode === 'OK') {
            this.$Notice.success({
              title: this.$t('successful'),
              desc: this.$t('successful')
            })
            this.loadPage()
          }
        },
        onCancel: () => {}
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
    async saveGroup () {
      let finalData = JSON.parse(JSON.stringify(this.finalElement[0]))
      finalData.items = finalData.attrs.map(attr => {
        if (attr.id.startsWith('c_')) {
          attr.id = ''
        }
        return attr
      })
      delete finalData.attrs
      const { statusCode } = await saveRequestGroupCustomForm(finalData)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.loadPage()
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
    gotoNext () {
      if (this.isParmasChanged || this.$refs.approvalFormNodeRef.panalStatus()) {
        this.$Modal.confirm({
          title: `${this.$t('confirm_discarding_changes')}`,
          content: `${this.$t('params_edit_confirm')}`,
          'z-index': 1000000,
          okText: this.$t('save'),
          cancelText: this.$t('abandon'),
          onOk: async () => {
            this.saveGroup()
            this.$refs.approvalFormNodeRef.saveNode()
          },
          onCancel: () => {
            this.$emit('gotoNextStep', this.requestTemplateId)
          }
        })
      } else {
        this.$emit('gotoNextStep', this.requestTemplateId)
      }
    }
  },
  components: {
    ApprovalFormNode,
    RequestFormDataCustom,
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
  margin: 0 10px;
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
.custom-title {
  width: 90px;
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
  .radio {
    padding: 5px 15px;
    border-radius: 32px;
    font-size: 12px;
    cursor: pointer;
    margin: 4px;
    display: inline-block;
  }
  .custom {
    border: 1px solid #b886f8;
    color: #b886f8;
  }
  .workflow {
    border: 1px solid #cba43f;
    color: #cba43f;
  }
  .optional {
    border: 1px solid #81b337;
    color: #81b337;
  }
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
  border: 1px solid #e8eaec;
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
</style>
