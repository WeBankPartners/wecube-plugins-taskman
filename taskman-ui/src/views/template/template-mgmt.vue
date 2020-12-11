<template>
  <div>
    <!-- <quillEditor></quillEditor> -->
    <!-- 表单设计 -->
    <Row>
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
    <Row>
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
    <Row>
      <div style="width:1000px;margin:10px auto;text-align: center">
        <Table style="margin-bottom:10px;" @on-selection-change="attrsChangedHandler" border :data="attrsData" :columns="attrsColumns"></Table>
        <Button >重置</Button>
        <Button type="primary" @click="attrsSetHandler" style="margin-left: 8px">下一步</Button>
      </div>
    </Row>
    <Row>
      <Row>
        <Form :model="requestForm" :label-width="100">
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
              <Select multiple @on-change="requestFormFieldChanged" v-model="requestForm.inputAttrDef">
                <Option v-for="(attr,index) in attrsSelections" :key="index" :value="attr.name" :label="attr.description"></Option>
              </Select>
            </FormItem>
          </Col>
        </Form>
      </Row>
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
          <TaskFormItem @delete="deleteHandler" @click.native="handleMouseClick(item)" @mouseenter.native="handleMouseEnter(item)" @mouseleave.native="handleMouseLeave(item)" v-for="(item, index) in formFields" :index="index" :item="item" :key="index"></TaskFormItem>
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
                <p  v-if="currentField.elementType !== 'Select'">当前表单项没有数据项</p>
                <Form v-if="currentField.elementType === 'Select'" :model="currentField" :label-width="100">
                  <FormItem label="数据类型">
                    <Select v-model="currentField.attrDataType">
                      <Option value="ref" label="引用类型"></Option>
                      <Option value="custom" label="自定义"></Option>
                    </Select>
                  </FormItem>
                  <FormItem v-if="currentField.attrDataType === 'ref'" label="目标对象">
                    <Select v-model="currentField.entityId">
                      <Option v-for="(comp, index) in allEntityList" :key="index" :value="comp.name" :label="comp.displayName"></Option>
                    </Select>
                  </FormItem>
                  <FormItem v-if="currentField.attrDataType === 'ref'" label="过滤规则">
                    <Input v-model="currentField.entityFilters"></Input>
                    <FilterRule
                      v-model="currentField.entityFilters"
                      :disabled="currentField.entityId.length === 0"
                      :rootEntity="currentField.entityId"
                      :allDataModelsWithAttrs="allEntityType"
                    ></FilterRule>
                  </FormItem>
                </Form>
                <div>
                  <Row>
                    <Col span="11">Label</Col>
                    <Col span="10">Value</Col>
                  </Row>
                  <Row style="margin-bottom:10px">
                    <Col span="10"><Input size="small" ></Input></Col>
                    <Col span="10" offset="1"><Input size="small" ></Input></Col>
                    <Col style="text-align: right" span="2" offset="1">
                      <Button size="small" ghost type="error" icon="ios-trash-outline"></Button>
                    </Col>
                  </Row>
                  <Row>
                    <Button size="small" icon="md-add" type="primary" style="width:100%"></Button>
                  </Row>
                </div>
              </div>
          </Panel>
        </Collapse>
      </Col>
    </Row>
    <Row>
      <div style="width:1000px;margin:10px auto;text-align: center">
        <Button >重置</Button>
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
  getAllDataModels
} from "../../api/server.js"
import PluginTable from "../../components/table";
import FilterRule from "../../components/filter-rule";
export default {
  components: {
    PluginTable,
    FilterRule
  },
  data() {
    return {
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
      ],
      currentField: {},
      componentsList: [
        {
          label: "选择框",
          type: "Select",
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
          type: "QuillEditor"
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
        }
      ],
      attrsSelections: [],
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
      allEntityList: []
    }
  },
  methods: {
    async getAllDataModels () {
      const { data, status } = await getAllDataModels()
      if (status === 'OK') {
        this.allEntityType = data.map(_ => {
          // handle result sort by name
          const pluginPackageEntities = _.pluginPackageEntities.sort(function (a, b) {
              var s = a.name.toLowerCase()
              var t = b.name.toLowerCase()
              if (s < t) return -1
              if (s > t) return 1
            })
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
      const payload = {
        ...this.requestForm,
        inputAttrDef: JSON.stringify(this.requestForm.inputAttrDef),
        otherAttrDef: JSON.stringify(this.attrsSelections),
        targetEntitys: JSON.stringify(this.entityList),
        tempId: this.currentTemplateId,
        formItems: this.formFields
      }
      const {status, message, data} = await saveFormTemplate(payload)
      if (status === 'OK') {
        this.$Notice.success({
          title: 'Success',
          desc: 'Success'
        })
      }
    },
    requestFormFieldChanged (val) {
      //this.attrsSelections val  this.formFields  isCustom
      const isAttrs = this.formFields.filter(field => !field.isCustom)
      isAttrs.forEach(attr => {
        const found = val.find(v => attr.name === v)
        if (!found) {
          const index = this.formFields.indexOf(attr)
          this.formFields.splice(index, 1)
        }
      })
      val.forEach(item => {
        const field = this.formFields.find(f => f.name === item)
        if (!field) {
          const attr = this.attrsSelections.find(a => a.name === item)
          this.formFields.push({
            ...attr,
            elementType: 'Input',
            width: 24,
            title: attr.description,
            defaultValue: "",
            isHover: false,
            isActive: false,
            isCustom: false,
            dataOptions: '',
            entityFilters:'',
            entityId:'',
            regular:'',
            attrDataType: '',
            attrDefId: attr.id,
            formTemplateId: this.currentTemplateId
          })
        }
      })
      this.currentFieldList = this.formFields
    },
    attrsSetHandler () {

    },
    attrsChangedHandler (selection) {
      this.attrsSelections = selection
    },
    actionFun(type, data) {
      switch (type) {
        case "add":
          // this.requestModalVisible = true;
          // this.isAdd = true;
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
        manageRoles:f.manageRoles.length>0? f.manageRoles.map(role => {
          const found = this.allRolesList.find(r => r.name === role)
          return {
            roleName: found.name,
            displayName: found.displayName
          }
        }):[],
        useRoles:f.useRoles.length>0? f.useRoles.map(role => {
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
            manageRolesName: _.manageRoles.map(role => role.displayName).join(', '),
            useRolesName: _.useRoles.map(role => role.displayName).join(', ')
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
        const nodes = await getTaskNodesEntitys(this.templateForm.procDefId)
        this.procTaskNodes = nodes.data
        let entitys = new Set()
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
        this.entityList = Array.from(entitys)
        this.list = Array.from(entitys)
        this.$Notice.success({
          title: 'Success',
          desc: 'Success'
        })
      }
    },
    handleMouseEnter (item) {
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
          // console.log(dataTransfer, dragEl)
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
              dataOptions: '',
              entityFilters:'',
              entityId:'',
              regular:'',
              attrDataType: '',
              attrDefId: '',
              formTemplateId: this.currentTemplateId,
              isCustom: true,
              entity: ''
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
    const nodes = await getTaskNodesEntitys('rYsEQg2D2Bu')
    this.procTaskNodes = nodes.data
    let entitys = new Set()
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
    this.entityList = Array.from(entitys)
    this.list = Array.from(entitys)
    this.currentEntityList = this.entityList
    this.formFieldSortHandler(this.entityList)
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
