<template>
  <Row type="flex">
    <Col span="20" offset="2">
      <div style="margin: 24px 0">
        <span v-for="approval in approvalNodes" :key="approval.id">
          <Tag
            color="primary"
            checkable
            @on-change="editNode(approval)"
            style="cursor: pointer"
            :checked="approval.id === activeEditingNode.id ? true : false"
            :type="approval.id === activeEditingNode.id ? 'none' : 'border'"
          >
            {{ approval.name }}
          </Tag>
          <Button
            @click.stop="removeNopde(approval)"
            type="error"
            size="small"
            ghost
            shape="circle"
            icon="md-trash"
            style="position: relative;right: 16px;bottom: 16px;"
          ></Button>
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
        <ApprovalFormNode ref="approvalFormNodeRef" @reloadParentPage="loadPage"></ApprovalFormNode>
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
              <Form :label-width="120" v-if="dataFormInfo.associationWorkflow">
                <FormItem :label="$t('tw_choose_object')">
                  <Select style="width: 30%">
                    <Option v-for="item in []" :value="item.id" :key="item.id">{{ item.displayName }}</Option>
                  </Select>
                </FormItem>
              </Form>
              <div style="margin-top:16px">
                <div
                  v-for="groupItem in dataFormInfo.groups"
                  :key="groupItem.itemGroupId"
                  style="display: inline-block;"
                >
                  <Button
                    @click.stop="editGroupItem(groupItem)"
                    type="primary"
                    size="small"
                    ghost
                    shape="circle"
                    style="position: relative;left: 16px;bottom: 16px;"
                    icon="md-create"
                  ></Button>
                  <Button
                    shape="circle"
                    :style="groupStyle[groupItem.itemGroupType]"
                    @click="editGroupCustomItems(groupItem)"
                    >{{ groupItem.itemGroupName }}</Button
                  >
                  <Button
                    @click.stop="removeGroupItem(groupItem)"
                    type="error"
                    size="small"
                    ghost
                    shape="circle"
                    icon="md-trash"
                    style="position: relative;right: 16px;bottom: 16px;"
                  ></Button>
                </div>
                <span>
                  <Icon
                    @click="selectItemGroup"
                    style="cursor: pointer;"
                    type="md-add-circle"
                    size="24"
                    color="#2d8cf0"
                  />
                </span>
              </div>
              <template v-if="finalElement.length === 1 && finalElement[0].itemGroup !== ''">
                <div
                  v-for="(item, itemIndex) in finalElement"
                  :key="itemIndex"
                  style="border: 1px solid #dcdee2; margin-bottom: 8px; padding: 8px"
                >
                  <span style="font-weight: 600;">
                    {{ item.itemGroupName }}
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
                      <Select
                        v-if="element.elementType === 'wecmdbEntity'"
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
                <div style="text-align: right;">
                  <Button type="primary" size="small" ghost @click="saveGroup">{{ $t('save') }}</Button>
                  <Button size="small" @click="cancelGroup">{{ $t('cancel') }}</Button>
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
  </Row>
</template>

<script>
import draggable from 'vuedraggable'
import ApprovalFormNode from './approval-form-node.vue'
import { getApprovalNode, addApprovalNode, removeApprovalNode, getApprovalGlobalForm } from '@/api/server.js'
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
        }
      ],
      dataFormInfo: {
        associationWorkflow: false,
        formTemplateId: '',
        groups: []
      },
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
          border: '1px solid #dd6da6',
          color: '#dd6da6'
        },
        workflow: {
          border: '1px solid #ba89f8',
          color: '#ba89f8'
        },
        optional: {
          border: '1px solid #81b337',
          color: '#81b337'
        }
      },
      groups: [], // 已存在组
      showSelectModel: false, // 组选择框
      itemGroupType: '', // 选中的组类型
      itemGroup: '' // 选中的组信息
    }
  },
  props: ['requestTemplateId'],
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 400
    this.loadPage()
  },
  methods: {
    loadPage () {
      this.getApprovalNode()
    },
    async getApprovalNode () {
      const { statusCode, data } = await getApprovalNode(this.requestTemplateId)
      if (statusCode === 'OK') {
        if (data.ids.length === 0) {
          this.addApprovalNode(1)
        } else {
          this.approvalNodes = data.ids
          this.activeEditingNode = this.approvalNodes[0]
          this.$refs.approvalFormNodeRef.loadPage({
            requestTemplateId: this.requestTemplateId,
            id: this.approvalNodes[0].id
          })
        }
      }
    },
    async addApprovalNode (sort) {
      const params = {
        requestTemplate: this.requestTemplateId,
        name: '审批',
        sort: sort
      }
      const { statusCode } = await addApprovalNode(params)
      if (statusCode === 'OK') {
        this.getApprovalNode()
      }
    },
    async removeNopde (node) {
      this.$Modal.confirm({
        title: this.$t('confirm_delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const { statusCode } = await removeApprovalNode(node.id)
          if (statusCode === 'OK') {
            this.getApprovalNode()
          }
        },
        onCancel: () => {}
      })
    },
    editNode (node) {
      console.log(234)
      this.activeEditingNode = node
      this.$refs.approvalFormNodeRef.loadPage({
        requestTemplateId: this.requestTemplateId,
        id: node.id
      })
    },
    // 查询可添加的组
    async selectItemGroup () {
      console.log(45)
      this.itemGroup = ''
      this.groupOptions = []
      const { statusCode, data } = await getApprovalGlobalForm(this.requestTemplateId)
      if (statusCode === 'OK') {
        console.log(77, data)
        // workflow  2.编排数据, 3.optional 自选数据项表单,  custom 1.自定义表单
        // this.$nextTick(() => {
        //   this.groupOptions = data.map(d => {
        //     if (d.formType === 'custom') {
        //       d.entities = ['custom']
        //     } else {
        //       d.entities = d.entities || []
        //     }
        //     return d
        //   })
        //   this.showSelectModel = true
        // })
      }
    },
    cloneDog (val) {
      if (this.$parent.isCheck === 'Y') return
      let newItem = JSON.parse(JSON.stringify(val))
      newItem.id = 'c_' + idGlobal++
      newItem.title = newItem.title + idGlobal
      newItem.name = newItem.name + idGlobal
      this.specialId = newItem.id
      this.paramsChanged()
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
    }
  },
  components: {
    ApprovalFormNode,
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
.basci-info-right {
  height: calc(100vh - 260px);
}

.basci-info-left {
  @extend .basci-info-right;
  border-right: 1px solid #dcdee2;
}

.title {
  font-size: 16px;
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
</style>
