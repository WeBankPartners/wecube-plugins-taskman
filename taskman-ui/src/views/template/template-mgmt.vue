<template>
  <div>
    <!-- <quillEditor></quillEditor> -->
    <!-- 表单设计 -->
    <Row v-show="!isEdit">
      <PluginTable
        :tableColumns="requestColumns"
        :tableData="requestTableData"
        :tableOuterActions="tableOuterActions"
        :pagination="requestPagination"
        @actionFun="actionFun"
        @handleSubmit="handleSubmit"
        @pageChange="requestPageChange"
        @pageSizeChange="requestPageSizeChange"
      />
    </Row>
    <Row v-show="isEdit" style="margin-bottom: 20px">
      <div style="margin-bottom: 10px"><Button @click="goBackQuery" icon="ios-arrow-back">返回</Button></div>
      <Row>

      <Col span= "22"> 
        <Steps style="margin-bottom: 20px" :current="currentStep" status="process">
          <Step title="设置模板基础信息" icon="md-add-circle" class="create_template_1"></Step>
          <Step title="选择表单项" icon="md-apps" class="create_template_2"></Step>
          <Step title="设置请求表单" icon="md-cog" class="create_template_3"></Step>
          <Step title="设置任务表单" icon="md-settings" class="create_template_4"></Step>
        </Steps>
       </Col>
       <Col style="text-align: right" span="2">
        <Button @click="releaseTemplate" type="primary" >模板发布</Button>
       </Col>
      </Row>
      <hr/>
    </Row>
    <Row v-show="isEdit && currentStep === 0">
      <div style="width:600px;margin:0 auto;">
        <Form :model="templateForm" :label-width="100">
          <FormItem label="名称">
            <Input v-model="templateForm.name"></Input>
          </FormItem>
          <FormItem label="模板组">
            <!-- <Input v-model="templateForm.templateGroup"></Input> -->
            <Select v-model="templateForm.requestTempGroup">
              <Option v-for="tem in templateGroupList" :key="tem.id" :value="tem.id" :label="tem.name"></Option>
            </Select>
          </FormItem>
          <FormItem label="流程">
            <Select label-in-value v-model="templateForm.procDefId">
              <Option v-for="process in allProcessDefinitionKeys" :key="process.id" :value="process.procDefId" :label="process.procDefName"></Option>
            </Select>
          </FormItem>
          <FormItem label="管理角色">
            <Select multiple v-model="templateForm.manageRoles">
              <Option v-for="role in allRolesList" :key="role.id" :value="role.name" :label="role.displayName"></Option>
            </Select>
          </FormItem>
          <FormItem label="使用角色">
            <Select multiple v-model="templateForm.useRoles">
              <Option v-for="role in allRolesList" :key="role.id" :value="role.name" :label="role.displayName"></Option>
            </Select>
          </FormItem>
          <FormItem label="标签">
            <Input v-model="templateForm.tags"></Input>
          </FormItem>
          <FormItem label="描述">
            <Input type="textarea" v-model="templateForm.description"></Input>
          </FormItem>
          <FormItem>
            <Button @click="templateFormHandleReset('formDynamic')">重置</Button>
            <Button type="primary" @click="templateFormHandleSubmit" style="margin-left: 8px">下一步</Button>
        </FormItem>
        </Form>
      </div>
    </Row>
    <Row v-show="isEdit && currentStep === 1">
      <div style="width:1000px;margin:10px auto;">
        <Tree :data="attrsTreeData" style="margin:o auto;" @on-check-change="attrsChangedHandler" multiple show-checkbox></Tree>

        <!-- <Table style="margin-bottom:10px;" ref="selection" @on-selection-change="attrsChangedHandler" border :data="attrsData" :columns="attrsColumns"></Table> -->
        <!-- <Button >重置</Button> -->
        <Button type="primary" @click="attrsSetHandler" style="margin-left: 8px;margin:10px auto;">下一步</Button>
      </div>
    </Row>
     <!-- v-show="isEdit && (currentStep === 2 || currentStep === 3)" -->
    <Row v-show="isEdit && (currentStep === 2 || currentStep === 3)">
      <Row>
        <Row v-show="isEdit && currentStep === 3">
          <Col style="margin-bottom:20px;font-size:14px;font-weight:600" span="1">
            任务节点:
          </Col>
          <Col span="20">
            <RadioGroup @on-change="taskNodeChanged" v-model="currentTaskNode" type="button" size="small">
              <Radio v-for="node in procTaskNodes" :key="node.nodeId" :label="node.nodeName"></Radio>
            </RadioGroup>
          </Col>
        </Row>
        <hr style="margin-bottom:20px" />
        <Form v-show="isEdit && currentStep === 2" :model="requestForm" :label-width="100">
          <Col span="6">
            <FormItem label="名称">
              <Input v-model="requestForm.name"></Input>
            </FormItem>
          </Col>
          <Col span="6">
            <FormItem label="描述">
              <Input v-model="requestForm.description"></Input>
            </FormItem>
          </Col>
          <Col span="12">
            <FormItem label="输入项">
              <TreeSelect
                v-model="requestForm.inputAttrDef"
                :maxTagCount="3"
                placeholder="输入项"
                :data="attrsSelections"
                @change="requestFormFieldChanged($event,0)"
                :clearable="true"
                style="width:100%"
              ></TreeSelect>
              <!-- <Select multiple @on-change="requestFormFieldChanged($event,0)" v-model="requestForm.inputAttrDef">
                <Option v-for="(attr,index) in attrsSelections" :key="index" :value="attr.name" :label="attr.displayName"></Option>
              </Select> -->
            </FormItem>
          </Col>
        </Form>
        <Form v-show="isEdit && currentStep === 3" :model="taskForm" :label-width="100">
          <Col span="6">
            <FormItem label="名称">
              <Input v-model="taskForm.name"></Input>
            </FormItem>
          </Col>
          <Col span="6">
            <FormItem label="描述">
              <Input v-model="taskForm.description"></Input>
            </FormItem>
          </Col>
          <Col span="12">
            <FormItem label="输入项">
              <TreeSelect
                v-model="taskForm.inputAttrDef"
                :maxTagCount="3"
                placeholder="输入项"
                :data="taskAttrsSelections"
                :clearable="true"
                style="width:100%"
              ></TreeSelect>
              <!-- <Select multiple v-model="taskForm.inputAttrDef">
                <Option v-for="(attr,index) in taskAttrsSelections" :key="index" :value="attr.name" :label="attr.displayName"></Option>
              </Select> -->
            </FormItem>
          </Col>
          <Col span="6">
            <FormItem label="管理角色">
              <Select multiple v-model="taskForm.manageRoles">
                <Option v-for="role in allRolesList" :key="role.id" :value="role.name" :label="role.displayName"></Option>
              </Select>
            </FormItem>
          </Col>
          <Col span="6">
            <FormItem label="处理角色">
              <Select multiple v-model="taskForm.useRoles">
                <Option v-for="role in allRolesList" :key="role.id" :value="role.name" :label="role.displayName"></Option>
              </Select>
            </FormItem>
          </Col>
          <Col span="12">
            <FormItem label="输出项">
              <TreeSelect
                v-model="taskForm.outputAttrDef"
                :maxTagCount="3"
                placeholder="输出项"
                :data="taskAttrsSelections"
                @change="requestFormFieldChanged($event,1)"
                :clearable="true"
                style="width:100%"
              ></TreeSelect>
              <!-- <Select multiple @on-change="requestFormFieldChanged($event,1)" v-model="taskForm.outputAttrDef">
                <Option v-for="(attr,index) in taskAttrsSelections" :key="index" :value="attr.name" :label="attr.displayName"></Option>
              </Select> -->
            </FormItem>
          </Col>
        </Form>
      </Row>
      <hr />
      <Col span="4" style="border-right: 1px solid rgb(224, 223, 222);height: calc(100vh - 100px);overflow:auto">
        <Divider>拖动设置目标对象顺序</Divider>
        <div ref="entity" style="padding-right:10px;">
          <p style="font-size:16px;background:bisque;margin-bottom:5px;text-align: center" v-for="(entity,index) in entityList" :id="entity" :key="index">{{entity}}</p>
        </div>
        <Divider>自定义表单项</Divider>
        <Row ref="fields">
          <Col span="6" v-for="(comp, index) in componentsList" :id="index" :key="index">
            <div class="components-box">
              {{comp.label}}
            </div>
          </Col>
        </Row>
      </Col>
      <Col span="14" style="height: calc(100vh - 100px);overflow:auto;background: rgb(248, 238, 226);padding:20px">
        <Form :model="formItem" ref="form" :label-width="100">
          <TaskFormItem :isDesign="true" @delete="deleteHandler" @click.native="handleMouseClick(item)" @mouseenter.native="handleMouseEnter(item)" @mouseleave.native="handleMouseLeave(item)" v-for="(item, index) in formFields" :index="index" :item="item" :key="index"></TaskFormItem>
        </Form>
      </Col>
      <Col style="border-left: 1px solid rgb(224, 223, 222);height: calc(100vh - 100px);overflow:auto" span="6">
        <Collapse>
          <Panel name="1">
              通用属性
              <div slot="content">
                <Form :model="currentField" :label-width="100">
                  <FormItem label="字段名">
                    <Input v-model="currentField.name"></Input>
                  </FormItem>
                  <FormItem label="显示名">
                    <Input v-model="currentField.title"></Input>
                  </FormItem>
                  <FormItem label="组件类型">
                    <Select v-model="currentField.elementType">
                      <Option v-for="(comp, index) in componentsList" :key="index" :value="comp.type" :label="comp.label"></Option>
                    </Select>
                  </FormItem>
                  <FormItem label="默认值">
                    <Input v-model="currentField.defaultValue"></Input>
                  </FormItem>
                  <FormItem label="宽度">
                    <Input v-model="currentField.width"></Input>
                  </FormItem>
                </Form>
              </div>
          </Panel>
          <Panel name="2">
              扩展属性
              <div slot="content">
                <Form :model="currentField" :label-width="100">
                  <FormItem label="校验规则">
                    <Input placeholder="仅支持正则" v-model="currentField.regular"></Input>
                  </FormItem>
                </Form>
              </div>
          </Panel>
          <Panel name="3">
              数据项
              <div slot="content">
                <p  v-if="currentField.elementType !== 'PluginSelect'">当前表单项没有数据项</p>
                <Form v-if="currentField.elementType === 'PluginSelect'" :model="currentField" :label-width="100">
                  <FormItem label="数据类型">
                    <Select v-model="currentField.attrDefDataType">
                      <Option value="ref" label="引用类型"></Option>
                      <Option value="str" label="自定义"></Option>
                    </Select>
                  </FormItem>
                  <FormItem v-if="currentField.attrDefDataType === 'ref'" label="目标对象">
                    <Select v-model="currentField.refEntity">
                      <Option v-for="(comp, index) in allEntityList" :key="index" :value="comp.name" :label="comp.displayName"></Option>
                    </Select>
                  </FormItem>
                  <!-- <FormItem v-if="currentField.attrDataType === 'ref'" label="过滤规则">
                    <Input v-model="currentField.entityFilters"></Input>
                    <FilterRule
                      v-model="currentField.entityFilters"
                      :disabled="currentField.entityId.length === 0"
                      :rootEntity="currentField.entityId"
                      :allDataModelsWithAttrs="allEntityType"
                    ></FilterRule>
                  </FormItem> -->
                </Form>
                <div v-if="currentField.elementType === 'PluginSelect' && currentField.attrDefDataType === 'str'">
                  <Row>
                    <Col span="11">Label</Col>
                    <Col span="10">Value</Col>
                  </Row>
                  <Row v-for="(opt, index) in currentField.dataOptions" :key="index" style="margin-bottom:10px">
                    <Col span="10"><Input v-model="opt.label" size="small" ></Input></Col>
                    <Col span="10" offset="1"><Input v-model="opt.value" size="small" ></Input></Col>
                    <Col style="text-align: right" span="2" offset="1">
                      <Button size="small" ghost type="error" @click="deleteDataOptionsItem(index)" icon="ios-trash-outline"></Button>
                    </Col>
                  </Row>
                  <Row>
                    <Button size="small" icon="md-add" type="primary" @click="addDataOptionsItem" style="width:100%"></Button>
                  </Row>
                </div>
              </div>
          </Panel>
        </Collapse>
      </Col>
    </Row>
    <Row v-show="isEdit && (currentStep === 2 || currentStep === 3)">
      <div style="width:1000px;margin:10px auto;text-align: center">
        <!-- <Button >重置</Button> -->
        <Button type="primary" @click="saveCurrentForm" style="margin-left: 8px">保存当前表单</Button>
      </div>
    </Row>
  </div>
</template>
<script>
// import quillEditor from "../../components/quillEditor"
import Sortable from 'sortablejs'
import {
  getRoleList,
  getProcessDefinitionKeysList,
  getAllTemplateGroup,
  saveRequestTemplate,
  searchRequestTemplate,
  deleteRequestTemplate,
  getTaskNodesEntitys,
  saveFormTemplate,
  saveTaskTemplate,
  getAllDataModels,
  getFormTemplateDetail,
  releaseRequestTemplate
} from "../../api/server.js"
import PluginTable from "../../components/table";
import FilterRule from "../../components/filter-rule";
import TreeSelect from '../../components/tree-select.vue'
import {addEvent, removeEvent} from "../../util/event"
export default {
  components: {
    PluginTable,
    FilterRule,
    TreeSelect
  },
  data() {
    return {
      formTemplateId: '',
      isEdit: false,
      isAdd: false,
      currentStep: 0,
      currentTaskNode: '',
      formItem: {
        input: ''
      },
      templateForm: {
        name: '',
        manageRoles: [],
        useRoles: [],
        procDefId: '',
        tags:'',
        description: '',
        requestTempGroup: ''
      },
      requestPayload: {
        filters: {},
        pageable: {
          pageSize: 10,
          startIndex: 0
        },
        paging: true
      },
      tableOuterActions: [
        {
          label: this.$t("add"),
          props: {
            type: "success",
            icon: "md-add",
            disabled: false
          },
          actionType: "add"
        }
      ],
      requestPagination: {
        currentPage: 1,
        pageSize: 10,
        total: 0
      },
      requestTableData: [],
      requestColumns: [
        {
          title: this.$t("name"),
          key: "name",
          inputKey: "name",
          component: "Input",
          inputType: "text"
        },
        {
          title: '流程',
          key: "procDefName",
          inputKey: "procDefName",
          component: "Input",
          inputType: "text",
          isNotFilterable: true
        },
        {
          title: '管理角色',
          key: "manageRolesName",
          inputKey: "manageRoles",
          isMultiple: true,
          span:5,
          component: "PluginSelect",
          options: []
        },
        {
          title: '使用角色',
          key: "useRolesName",
          inputKey: "useRoles",
          isMultiple: true,
          span:5,
          component: "PluginSelect",
          options: []
        },
        {
          title: '标签',
          key: "tags",
          inputKey: "tags",
          component: "Input",
          inputType: "text",
        },
        {
          title: this.$t("describe"),
          key: "description",
          inputKey: "description",
          component: "Input",
          inputType: "text",
          isNotFilterable: true
        },
        {
          title: this.$t("action"),
          key: "action",
          width: 150,
          align: "center",
          isNotFilterable: true,
          render: (h, params) => {
            return (
              <span>
                  <Button
                    type="primary"
                    size="small"
                    style="margin-right:10px"
                    onClick={() => this.edit(params.row)}
                  >
                    {this.$t("edit")}
                  </Button>
              </span>
            );
          }
        }
      ],
      currentField: {},
      componentsList: [
        {
          label: "选择框",
          type: "PluginSelect",
        },
        {
          label: "输入框",
          type: "Input",
        },
        {
          label: "多行文本",
          type: "Textarea",
        },
        {
          label: "富文本",
          type: "TaskManQuillEditor"
        },
      ],
      formFields: [],
      entityList: [],
      list:[],
      currentEntityList:[],
      currentFieldList: [],
      allRolesList: [],
      allProcessDefinitionKeys: [],
      templateGroupList: [],
      currentTemplateId: '',
      attrsData: [],
      attrsTreeData: [],
      procTaskNodes:[],
      attrsColumns: [
        {
          type: 'selection',
          width: 60,
          align: 'center',
          fixed: 'left'
        },
        {
          title: '包名',
          key: "packageName",
        },
        {
          title: '对象',
          key: "entity",
        },
        {
          title: '属性',
          key: "name",
        },
        {
          title: '描述',
          key: "description",
        }
      ],
      attrsSelections: [],
      taskAttrsSelections: [],
      requestForm: {
        name:'',
        description:'',
        inputAttrDef: []
      },
      taskForm: {
        name:'',
        description:'',
        inputAttrDef: [],
        outputAttrDef: [],
        useRoles: [],
        manageRoles:[]
      },
      allEntityType: [],
      allEntityList: [],
      currentFormTemplate: {},
      currentFirstTaskFormTemplate: {},
    }
  },
  methods: {
    async releaseTemplate () {
      const process = this.allProcessDefinitionKeys.find(key => key.procDefId === this.templateForm.procDefId)
      const payload = {
        procDefKey: process.procDefKey,
        id: this.currentTemplateId
      }
      const {status, message, data} = await releaseRequestTemplate(payload)
      if (status === 'OK') {
        this.$Notice.success({
          title: 'Success',
          desc: 'Success'
        })
      }
    },
    async getFormTemplateDetail (type, id, isEdit) {
      const {status, message, data} = await getFormTemplateDetail(type, id)
      if (status === 'OK') {
        this.isAdd = false
        this.formTemplateId = data.id ? data.id : undefined
        this.formFields = data.items ? data.items.map( _ => {
            return {
              ..._,
              dataOptions: _.dataOptions.length > 0 ? JSON.parse(_.dataOptions) : [],
              isCustom: !(_.packageName.length > 0),
              isActive:false,
              isHover: false
            }
          }) : []
          this.currentFieldList = this.formFields
        if (type === 0) {
          this.attrsSelections = data.otherAttrDef && data.otherAttrDef.length > 0 ? JSON.parse(data.otherAttrDef) : []
          this.taskAttrsSelections = this.attrsSelections.concat([{
            title: '自定义',
            expand: true,
            children:this.formFields.filter(_ => _.isCustom).map(f => {return {...f, displayName: f.displayName ? f.displayName : f.title}})
          }])
          this.entityList = data.targetEntitys ? JSON.parse(data.targetEntitys) : []
          this.list = this.entityList
          this.currentEntityList = this.entityList
        }
        if (isEdit) {
          this.attrsTreeData.forEach(i => {
            i.children.forEach(child => {
              this.attrsSelections.forEach(j => {
                j.children.forEach(d => {
                  if (d.id === child.id) {
                    child.checked = true
                  }
                })
              })
            })
          })
        }
        if (type === 0) {
          this.requestForm.name = data.name
          this.requestForm.inputAttrDef = data.inputAttrDef ? JSON.parse(data.inputAttrDef) : []
          this.requestForm.description = data.description
        }
        if (type === 1) {
          this.taskForm.name = data.name
          this.taskForm.inputAttrDef = data.inputAttrDef ? JSON.parse(data.inputAttrDef) : []
          this.taskForm.description = data.description
          this.taskForm.useRoles = data.useRoles
          this.taskForm.outputAttrDef = data.outputAttrDef ? JSON.parse(data.outputAttrDef) : []
          this.taskForm.manageRoles = data.manageRoles
        }
        this.formFieldSortHandler(this.currentEntityList)
      }
    },
    async edit (row) {
      this.templateForm = {
        ...row,
        manageRoles: row.manageRoles ? row.manageRoles.map(role => role.roleName) : [],
        useRoles: row.useRoles ? row.useRoles.map(role => role.roleName) : []
      }
      this.currentTemplateId = row.id
      await this.getTaskNodesEntitys(this.templateForm.procDefId)
      this.getFormTemplateDetail(0,row.id,true)
      this.isEdit = true
      this.currentStep = 0
      removeEvent('.ivu-steps-title', 'mouseover', this.handleStepMouseover)
      removeEvent('.ivu-steps-title', 'click', this.handleStepClick)
      this.addElementEvent()
    },
    deleteDataOptionsItem (index) {
      this.currentField.dataOptions.splice(index, 1)
    },
    addDataOptionsItem () {
      this.currentField.dataOptions.push({label:'',value:''})
    },
    taskNodeChanged (v) {
      this.currentTaskNode = v
      const id = this.procTaskNodes.find(node => node.nodeName === this.currentTaskNode).nodeDefId
      this.getFormTemplateDetail(1,id)
      this.currentField = {}
    },
    goBackQuery () {
      this.isEdit = false
    },
    addElementEvent () {
      addEvent('.ivu-steps-title', 'mouseover', this.handleStepMouseover)
      addEvent('.ivu-steps-title', 'click', this.handleStepClick)
    },
    handleStepMouseover (e) {
      e.preventDefault()
      e.stopPropagation()
      e.currentTarget.style.cursor = 'pointer'
    },
    handleStepClick (e) {
      e.preventDefault()
      e.stopPropagation()
      this.currentStep = e.currentTarget.parentNode.parentNode.classList[0].split('_')[2] * 1 - 1 
      this.currentField = {}
      if (this.currentStep > 1) {
        this.$nextTick(() => {
          this.formFieldSortHandler(this.currentEntityList)
        })
      }
      if (this.currentStep === 3) {
        this.formFields = []
        this.currentFieldList = this.formFields
        this.currentTaskNode = this.procTaskNodes[0].nodeName
        const id = this.procTaskNodes.find(node => node.nodeName === this.currentTaskNode).nodeDefId
        this.getFormTemplateDetail(1,id)
      }
      if (this.currentStep === 2) {
        this.formFields = []
        this.currentFieldList = this.formFields
        this.getFormTemplateDetail(0,this.currentTemplateId)
      }
    },
    async getAllDataModels () {
      const { data, status } = await getAllDataModels()
      if (status === 'OK') {
        this.allEntityType = data.map(_ => {
          // handle result sort by name
          const pluginPackageEntities = _.entities ? _.entities.sort(function (a, b) {
              var s = a.name.toLowerCase()
              var t = b.name.toLowerCase()
              if (s < t) return -1
              if (s > t) return 1
            }) : []
          this.allEntityList = this.allEntityList.concat(pluginPackageEntities)
          return {
            ..._,
            pluginPackageEntities: pluginPackageEntities
          }
        })
      }
    },
    async saveCurrentForm () {
      // saveFormTemplate
      let payload = {}
      if (this.currentStep === 2) {
        payload = {
          ...this.requestForm,
          inputAttrDef: JSON.stringify(this.requestForm.inputAttrDef),
          otherAttrDef: JSON.stringify(this.attrsSelections),
          targetEntitys: JSON.stringify(this.entityList),
          tempId: this.currentTemplateId,
          id: this.formTemplateId,
          formItems: this.formFields.map((i,index) => {
            return {
              ...i,
              sort:index,
              dataOptions: JSON.stringify(i.dataOptions)
            }
          })
        }
      }
      if (this.currentStep === 3) {
        const process = this.allProcessDefinitionKeys.find(key => key.procDefId === this.templateForm.procDefId)
        payload = {
          ...process,
          ...this.taskForm,
          id: this.formTemplateId,
          inputAttrDef: JSON.stringify(this.taskForm.inputAttrDef),
          otherAttrDef: JSON.stringify(this.taskAttrsSelections),
          outputAttrDef: JSON.stringify(this.taskForm.outputAttrDef),
          nodeDefId: this.procTaskNodes.find(node => node.nodeName === this.currentTaskNode).nodeDefId,
          nodeName: this.currentTaskNode,
          manageRoles:this.taskForm.manageRoles.length>0? this.taskForm.manageRoles.map(role => {
              const found = this.allRolesList.find(r => r.name === role)
              return {
                roleName: found.name,
                displayName: found.displayName
              }
            }):[],
            useRoles:this.taskForm.useRoles.length > 0 ? this.taskForm.useRoles.map(role => {
              const found = this.allRolesList.find(r => r.name === role)
              return {
                roleName: found.name,
                displayName: found.displayName
              }
            }):[],
          form: {
            ...this.taskForm,
            targetEntitys: JSON.stringify(this.entityList),
            manageRoles:this.taskForm.manageRoles.length>0? this.taskForm.manageRoles.map(role => {
              const found = this.allRolesList.find(r => r.name === role)
              return {
                roleName: found.name,
                displayName: found.displayName
              }
            }):[],
            useRoles:this.taskForm.useRoles.length>0? this.taskForm.useRoles.map(role => {
              const found = this.allRolesList.find(r => r.name === role)
              return {
                roleName: found.name,
                displayName: found.displayName
              }
            }):[],
            tempId: this.currentTemplateId,
            inputAttrDef: JSON.stringify(this.taskForm.inputAttrDef),
            otherAttrDef: JSON.stringify(this.taskAttrsSelections),
            outputAttrDef: JSON.stringify(this.taskForm.outputAttrDef)
          },
          tempId: this.currentTemplateId,
          formItems: this.formFields.map((i,index) => {
            return {
              ...i,
              sort:index,
              dataOptions: JSON.stringify(i.dataOptions)
            }
          })
        }
      }
      const {status, message, data} = this.currentStep === 2 ? await saveFormTemplate(payload) : await saveTaskTemplate(payload)
      if (status === 'OK') {
        this.$Notice.success({
          title: 'Success',
          desc: 'Success'
        })
        if (this.currentStep === 2) {
          this.taskAttrsSelections = this.attrsSelections.concat([{
            title: '自定义',
            expand: true,
            children:this.formFields.filter(_ => _.isCustom).map(f => {return {...f, displayName: f.displayName ? f.displayName : f.title}})
          }])
          this.currentStep++
          this.formFields = []
          this.currentFieldList = []
          this.currentField = {}
          this.currentTaskNode = this.procTaskNodes[0].nodeName
        }
      }
    },
    requestFormFieldChanged (val,type) {
      //this.attrsSelections val  this.formFields  isCustom taskAttrsSelections
      const isAttrs = this.formFields.filter(field => !field.isCustom)
      isAttrs.forEach(attr => {
        const found = val.filter(e => !e.isEntity).find(v => attr.name === v.name)
        if (!found) {
          const index = this.formFields.indexOf(attr)
          this.formFields.splice(index, 1)
        }
      })
      const selections = type === 0 ? [].concat(...this.attrsSelections.map(_ => _.children)) : [].concat(...this.taskAttrsSelections.map(_ => _.children))
      val.filter(e => !e.isEntity).forEach(item => {
        const field = this.formFields.find(f => f.name === item.name)
        if (!field) {
          const attr = selections.find(a => a.name === item.name)
          this.formFields.push({
            ...attr,
            entity: attr.entity,
            packageName: attr.packageName,
            elementType: attr.dataType === 'ref' || attr.dataType === 'multiRef' ? 'PluginSelect' : 'Input',
            width: 24,
            title: attr.description,
            defaultValue: "",
            isHover: false,
            isActive: false,
            isCustom: attr.isCustom ? true : false,
            dataOptions: [],
            entityFilters:'',
            refEntity:'',
            regular:'',
            attrDefDataType: 'str',
            attrDefId: attr.id,
            formTemplateId: this.currentTemplateId
          })
        }
      })
      const tem = []
      this.formFields.concat(this.currentFieldList.filter(field => field.isCustom)).forEach(_ => {
        const field = tem.find(i => i.name === _.name)
        if (!field) {
          tem.push(_)
        }
      })
      this.formFields = tem
      this.currentFieldList = this.formFields
      this.$nextTick(() => {
          this.formFieldSortHandler(this.currentEntityList)
        })
    },
    attrsSetHandler () {
      this.currentStep++
    },
    attrsChangedHandler (selection, current) {
      let data = []
      selection.filter(i => !i.isEntity).forEach(sel => {
        const found = data.find(d => d.title === sel.entity)
        if (found) {
          found.children.push({
            ...sel,
            checked: false,
            nodeKey: null
          })
        } else {
          data.push({
            title: sel.entity,
            checked: false,
            isEntity: true,
            expand: true,
            children: [{
              ...sel,
              checked: false,
              nodeKey: null
            }]
          })
        }
      })
      this.attrsSelections = data
      this.taskAttrsSelections = data
    },
    actionFun(type, data) {
      switch (type) {
        case "add":
          this.templateForm = {
            name: '',
            manageRoles: [],
            useRoles: [],
            procDefId: '',
            tags:'',
            description: '',
            requestTempGroup: ''
          }
          this.requestForm = {
            name:'',
            description:'',
            inputAttrDef: []
          },
          this.taskForm = {
            name:'',
            description:'',
            inputAttrDef: [],
            outputAttrDef: [],
            useRoles: [],
            manageRoles:[]
          }
          this.attrsData = []
          this.attrsTreeData = []
          this.formTemplateId = ''
          this.currentTemplateId = ''
          this.attrsSelections = []
          this.taskAttrsSelections = []
          this.isEdit = true
          this.isAdd = true
          this.currentStep = 0
          removeEvent('.ivu-steps-title', 'mouseover', this.handleStepMouseover)
          removeEvent('.ivu-steps-title', 'click', this.handleStepClick)
          this.addElementEvent()
          break;
      }
    },
    requestPageChange(current) {
      this.requestPagination.currentPage = current;
      this.getData();
    },
    requestPageSizeChange(size) {
      this.requestPagination.pageSize = size;
      this.getData();
    },
    async getData () {
      const f = this.requestPayload.filters
      const filters = {
        ...f,
        manageRoles:f.manageRoles && f.manageRoles.length > 0 ? f.manageRoles.map(role => {
          const found = this.allRolesList.find(r => r.name === role)
          return {
            roleName: found.name,
            displayName: found.displayName
          }
        }):[],
        useRoles:f.useRoles && f.useRoles.length > 0 ? f.useRoles.map(role => {
          const found = this.allRolesList.find(r => r.name === role)
          return {
            roleName: found.name,
            displayName: found.displayName
          }
        }):[]
      }
      const payload = {
        page: this.requestPagination.currentPage,
        pageSize: this.requestPagination.pageSize,
        data: filters
      }
      const { data, status} = await searchRequestTemplate(payload)
      if(status === 'OK') {
        this.requestTableData = data.contents.map(_ => {
          return {
            ..._,
            manageRolesName: _.manageRoles ? _.manageRoles.map(role => role.displayName).join(', ') : '',
            useRolesName: _.useRoles ? _.useRoles.map(role => role.displayName).join(', ') : ''
          }
        })
        this.requestPagination.total = data.pageInfo.totalRows;
      }
    },

    handleSubmit(filters) {
      this.requestPayload.filters = filters;
      this.getData();
    },
    async getAllTemplateGroup () {
      const {status, message, data} = await getAllTemplateGroup()
      if (status === 'OK') {
        this.templateGroupList = data
      }
    },
    async getRoleList () {
      const {status, message, data} = await getRoleList()
      if (status === 'OK') {
        this.allRolesList = data
        const roles =  data.map(role => {
          return {
            value: role.name,
            label: role.displayName
          }
        })
        this.requestColumns.forEach(column => {
          if (column.key === 'manageRolesName' || column.key === 'useRolesName') {
            column.options = roles
          }
        })
      }
    },
    async getProcessDefinitionKeysList () {
      const {status, message, data} = await getProcessDefinitionKeysList()
      if (status === 'OK') {
        this.allProcessDefinitionKeys = data
      }
    },
    templateFormHandleReset () {

    },
    async templateFormHandleSubmit () {
      const process = this.allProcessDefinitionKeys.find(key => key.procDefId === this.templateForm.procDefId)
      const payload = {
        ...this.templateForm,
        procDefName: process.procDefName,
        procDefKey: process.procDefKey,
        manageRoles:this.templateForm.manageRoles.length>0? this.templateForm.manageRoles.map(role => {
          const found = this.allRolesList.find(r => r.name === role)
          return {
            roleName: found.name,
            displayName: found.displayName
          }
        }):[],
        useRoles:this.templateForm.useRoles.length>0? this.templateForm.useRoles.map(role => {
          const found = this.allRolesList.find(r => r.name === role)
          return {
            roleName: found.name,
            displayName: found.displayName
          }
        }):[]
      }
      const {status, message, data} = await saveRequestTemplate(payload)
      if (status === 'OK') {
        this.currentTemplateId = data.id
        this.$Notice.success({
          title: 'Success',
          desc: 'Success'
        })
        this.currentStep++
        if (this.attrsSelections.length > 0) {
          // let objData = this.$refs.selection.objData
          this.attrsTreeData.forEach(i => {
            i.children.forEach(child => {
              this.attrsSelections.forEach(j => {
                if (j.id === child.id) {
                  child.checked = true
                }
              })
            })
          })
          // Object.keys(objData).forEach(i => {
          //   this.attrsSelections.forEach(j => {
          //     if (j.id === objData[i].id) {
          //       objData[i]._isChecked = true
          //     }
          //   })
          // })
        } else {
          this.getTaskNodesEntitys(this.templateForm.procDefId)
        }
      }
    },
    async getTaskNodesEntitys (id) {
      const nodes = await getTaskNodesEntitys(id)
      // this.procTaskNodes = nodes.data ? nodes.data.filter(node => node.taskCategory === 'SUTN') : []
      this.procTaskNodes = nodes.data
      let entitys = new Set()
      const entityData = nodes.data ? nodes.data : []
      // entityData.filter(f=>f.boundEntity).forEach(node => {
      //   const entity = node.boundEntity
      //   entitys.add(entity.name)
      //   const found = this.attrsData.find(_ => _.name == entity.name)
      //   if (!found) {
      //     this.attrsData.push(entity)
      //   }
      // })
      this.attrsData = [].concat(...this.procTaskNodes.filter(f=>f.boundEntity && f.boundEntity.attributes).map(node => {
        const entity = node.boundEntity
        entitys.add(entity.name)
        return node.boundEntity.attributes.map(attr => {
          return {
            ...attr,
            packageName:entity.packageName,
            entity:entity.name
          }
        })
      }))
      this.attrsTreeData = []
      this.procTaskNodes.filter(f=>f.boundEntity && f.boundEntity.attributes).forEach(node => {
        const entity = node.boundEntity
        entitys.add(entity.name)
        const found = this.attrsTreeData.find(e => e.title === entity.name)
        if (!found) {
          this.attrsTreeData.push({
            title: entity.name,
            checked: false,
            isEntity: true,
            children: node.boundEntity.attributes.map(attr => {
              return {
                ...attr,
                title: attr.name,
                packageName:entity.packageName,
                entity:entity.name,
                checked: false
              }
            })
          })
        }
      })
      this.entityList = Array.from(entitys)
      this.list = Array.from(entitys)
      this.currentEntityList = this.entityList
    },
    handleMouseEnter (item) {
      if (item.isActive) {
        return
      }
      item.isHover = true
    },
    handleMouseLeave (item) {
      item.isHover = false
    },
    handleMouseClick (item) {
      if(item.isActive) {
        item.isActive = false
        return
      }
      this.formFields.forEach(formField => {formField.isActive = false})
      item.isActive = true
      item.isHover = false
      this.currentField = item
    },
    deleteHandler (index) {
      // console.log(index)
    },
    formFieldSortHandler (entityList) {
      let fields = []
      this.formFields = []
      entityList.forEach(entity => {
        this.currentFieldList.forEach(formField => {
          if (formField.entity === entity) {
            fields.push(formField)
          }
        })
      })
      this.$nextTick(() => {
        this.formFields = fields.concat(this.currentFieldList.filter(field => field.isCustom))
        this.currentFieldList = this.formFields
      })
    },
    createSortable(el, items) {
      return new Sortable(el, {
        group: {
          name: 'component',
          pull: 'clone',
          put: false
        },
        sort: false,
        animation: 150,
        setData: (dataTransfer, dragEl) => {
          // const index = parseInt(dragEl.dataset.id)
          // dragEl.__item__ = items[index]
        }
      })
    },
    createEntitySortable(el, items) {
      return new Sortable(el, {
        group: {
          name: 'component',
          pull: 'clone',
          put: false
        },
        animation: 150,
        onUpdate: event => {
          let entitys = []
          for (let i = 0; i < el.children.length; i++) {
            entitys.push(el.children[i].attributes.id.value)
          }
          this.currentEntityList = entitys
          this.formFieldSortHandler(entitys)
        }
      })
    },
    createFormSortable(el) {
      let isAdd = false
      this.sortable = new Sortable(el, {
        group: 'component',
        animation: 150,
        onStart: () => {
          this.dragging = true
          this.showHelper = null
        },
        onEnd: (e) => {
          this.dragging = false
        },
        onSort: (e) => {
          // 添加也会触发onSort， 用个变量去来区分
          if(!isAdd) {
            let fields = []
            for (let i = 0; i < el.children.length; i++) {
              const id = el.children[i].attributes.id.value.split('_')[1]
              fields.push(this.formFields[id])
            }
            this.currentFieldList = fields
          }
          // this.formFields = []
          this.formFieldSortHandler(this.currentEntityList)
          isAdd = false
        },
        onAdd: (e) => {
          const id = e.item.attributes.id.value
          e.item.parentNode.removeChild(e.item)
          function randomString(len) {
        　　len = len || 32;
            let $chars = 'ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyzoLlgqVvUuI';
            let maxPos = $chars.length;
            let pwd = '';
        　　for (let i = 0; i < len; i++) {
        　　　　pwd += $chars.charAt(Math.floor(Math.random() * maxPos));
        　　}
        　　return pwd;
          }
          // this.$nextTick(() => {
            const item = {
              elementType: this.componentsList[id*1].type,
              width: 24,
              title: '自定义',
              name: randomString(16),
              defaultValue: "",
              isHover: false,
              isActive: false,
              dataOptions: [],
              entityFilters:'',
              entityId:'',
              regular:'',
              attrDataType: '',
              attrDefId: '',
              formTemplateId: this.currentTemplateId,
              packageName: '',
              isCustom: true,
              entity: '',
              sort: e.newIndex
            }
            this.formFields.splice(e.newIndex, 0, item)
            this.currentFieldList = this.formFields
            isAdd = true
            // this.formFieldSortHandler(this.currentEntityList)
          // })
        }
      })
    }
  },
  async created () {
    this.getRoleList()
    this.getProcessDefinitionKeysList()
    this.getAllTemplateGroup()
    this.getAllDataModels()
    this.formFieldSortHandler(this.entityList)
    this.getData()
  },
  mounted() {
    this.fieldsSortable = this.createSortable(this.$refs.fields.$el, this.componentsList)
    this.formSortable = this.createFormSortable(this.$refs.form.$el)
    this.entitySortable = this.createEntitySortable(this.$refs.entity, this.list)
  },
  beforeDestroy() {
    this.fieldsSortable && this.fieldsSortable.destroy()
    this.formSortable && this.formSortable.destroy()
    this.entitySortable && this.entitySortable.destroy()
    removeEvent('.ivu-steps-title', 'mouseover', this.handleStepMouseover)
    removeEvent('.ivu-steps-title', 'click', this.handleStepClick)
  }
}
</script>
<style lang="scss">
.components-box {
    background: bisque;
    width: 60px;
    height: 60px;
    text-align: center;
    line-height: 60px;
    cursor: move;
  }
.mask {
  height: 100%;
}
</style>
