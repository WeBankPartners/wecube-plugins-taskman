<template>
  <div>
    <PluginTable
      v-if="isQuery"
      :tableColumns="requestColumns"
      :tableData="requestTableData"
      :tableOuterActions="tableOuterActions"
      :pagination="requestPagination"
      @actionFun="actionFun"
      @handleSubmit="handleSubmit"
      @pageChange="requestPageChange"
      @pageSizeChange="requestPageSizeChange"
    />
    <Button v-show="!isQuery" style="margin-top: 10px;margin-bottom: 10px" @click="queryRequest" type="primary">查看已提交请求</Button>
    <div v-if="!isQuery && requestStep === 1" class="template-container">
      <Card style="margin-bottom: 10px" v-for="group in allTemplatesTree" :key="group.requestTempGroupName">
        <p slot="title">{{group.requestTempGroupName}}</p>
        <div style="padding-left:20px" v-for="tag in group.children">
          <p style="font-size: 17px;font-weight: 500;background-color: #f8f8f9;border-radius: 5px;">{{tag.tag}}</p>
          <div style="padding-left:20px;padding-top: 10px;padding-bottom:10px">
            <Tag v-for="tem in tag.children" @on-change="templateChanged($event,tem.id)" checkable style="cursor: pointer;border-color: #adc6ff;background: #f0f5ff;border: 1px solid #e8eaec;" color="geekblue" :key="tem.id" size="medium">{{tem.name}}</Tag>
          </div>
        </div>
      </Card>
    </div>
    <div v-show="!isQuery && requestStep === 2" style="padding:20px">
      <Form ref="requestForm" :rules="ruleValidate" :model="requestForm" :label-width="110">
        <Row>
          <!-- <Col span="12">
        <FormItem :label="$t('template')">
          <Select filterable @on-open-change="getTemplates" @on-change="templateChanged" v-model="requestForm.requestTempId">
            <Option v-for="tem in allTemplates" :key="tem.id" :value="tem.id">{{tem.name}}</Option>
          </Select>
        </FormItem>
        </Col> -->
        <Col span="12">
        <FormItem :label="$t('target_object')">
          <Select filterable @on-open-change="getEntityDataByTemplateId" @on-change="workflowProcessPrevieEntities" v-model="requestForm.rootEntity">
            <Option v-for="tem in entityData" :key="tem.guid" :value="tem.guid">{{tem.displayName}}</Option>
          </Select>
        </FormItem>
        </Col>
        <Col span="12">
        <FormItem :label="$t('service_request_name')" prop="name">
          <Input v-model="requestForm.name" :placeholder="$t('service_request_name')"></Input>
        </FormItem>
        </Col>
        <Col span="12">
        <!-- <FormItem :label="$t('service_request_role')">
          <Select @on-open-change="getRolesByCurrentUser" v-model="requestForm.roleName">
            <Option
              v-for="role in currentUserRoles"
              :key="role.name"
              :value="role.name"
            >{{role.displayName}}</Option>
          </Select>
        </FormItem> -->
        <FormItem :label="$t('emergency_level')">
          <Select v-model="requestForm.emergency">
            <Option value="normal">{{$t('not_urgent')}}</Option>
            <Option value="urgent">{{$t('emergency')}}</Option>
          </Select>
        </FormItem>
        </Col>
        <Col span="12">
        <!-- <FormItem :label="$t('reqest_attachment')">
          <Upload
            :on-success="uploadSuccess"
            ref="upload"
            action="/service-mgmt/v1/service-requests/attach-file"
          >
            <Button icon="ios-cloud-upload-outline">{{$t('upload_attachment')}}</Button>
          </Upload>
        </FormItem> -->
        </Col>
        <Col span="12">
        <FormItem :label="$t('describe')">
          <Input type="textarea" v-model="requestForm.description" :placeholder="$t('describe')"></Input>
        </FormItem>
        </Col>
        </Row>
        <Divider>申请数据</Divider>
        <!-- <hr style="margin-bottom:10px"/> -->
        <Row>
          <Form ref="form" :label-width="110">
            <div style="border: 1px solid #808695;border-radius: 5px;margin-bottom:5px;padding:10px" v-for="(key, i) in targetEntityList" :key="i">
              <Row v-for="(row, d) in currentFields[key]" :key="d" >
                <Col span="24">
                  <TaskFormItem v-for="(item, index) in row.fields" :index="index" v-model="item.value" :item="item.attribute" :key="index"></TaskFormItem>
                  <hr v-if="d > 0" style="margin-bottom:10px"/>
                </Col>
              </Row>
                <hr style="margin-bottom:10px"/>
              <div style="text-align:right;">
                <Button @click="addRowData(key)">新增数据</Button>
              </div>
            </div>
          </Form>
        </Row>
        <div style="text-align:center;margin-top: 20px">
          <!-- <Button type="primary" :loading="requestLoading" @click="requestSubmit">下一步</Button> -->
          <Button type="primary" @click="goNextStep">下一步</Button>
          <Button style="margin-left: 15px" @click="requestCancel">{{$t('cancle')}}</Button>
        </div>
      </Form>
    </div>
    <div v-show="!isQuery && requestStep === 3" style="padding:20px">
      <Row>
        <Col style="margin-bottom:20px;font-size:18px;font-weight:600" span="2">
          任务节点:
        </Col>
        <Col span="20">
          <RadioGroup @on-change="taskNodeChanged" v-model="currentTaskNode" type="button">
            <Radio v-for="node in procTaskNodes" :key="node.nodeId" :label="node.nodeName"></Radio>
          </RadioGroup>
        </Col>
      </Row>
      <Row>
        <div style="text-align:center;padding-left:100px;padding-right:100px">
          <Table border @on-selection-change="selectionChanged" :data="nodeData" :columns="targetModelColums"></Table>
        </div>
        <div style="text-align:center;margin-top: 10px">
          <Button type="primary" :loading="requestLoading" @click="requestSubmit">提交</Button>
          <Button style="margin-left: 15px" @click="goBack">上一步</Button>
        </div>
      </Row>
    </div>
    <Modal
      v-model="detailModalVisible"
      title="详情"
      footer-hide
      fullscreen
      @on-cancel="detailCancel"
    >
      <div style="padding:20px">
        <Form ref="requestForm" :model="detailForm" :label-width="110">
          <Row>
            <Col span="12">
          <FormItem :label="$t('template')">
            <Select disabled v-model="detailForm.requestTempId">
              <Option v-for="tem in allTemplates" :key="tem.id" :value="tem.id">{{tem.name}}</Option>
            </Select>
          </FormItem>
          </Col>
          <Col span="12">
          <FormItem :label="$t('target_object')">
            <Select disabled v-model="detailForm.rootEntity">
              <Option v-for="tem in detailEntityData" :key="tem.guid" :value="tem.guid">{{tem.displayName}}</Option>
            </Select>
          </FormItem>
          </Col>
          <Col span="12">
          <FormItem :label="$t('service_request_name')" prop="name">
            <Input disabled v-model="detailForm.name" :placeholder="$t('service_request_name')"></Input>
          </FormItem>
          </Col>
          <Col span="12">
          <FormItem :label="$t('emergency_level')">
            <Select disabled v-model="detailForm.emergency">
              <Option value="normal">{{$t('not_urgent')}}</Option>
              <Option value="urgent">{{$t('emergency')}}</Option>
            </Select>
          </FormItem>
          </Col>
          <Col span="12">
          <FormItem :label="$t('describe')">
            <Input type="textarea" v-model="detailForm.description" :placeholder="$t('describe')"></Input>
          </FormItem>
          </Col>
          </Row>
          <hr style="margin-bottom:10px"/>
          <Row>
            <Form ref="form" :label-width="110">
              <Row v-for="(fields, i) in Object.values(detailFields)" :key="i">
                <TaskFormItem v-for="(item, index) in fields" :index="index" v-model="currentForm[item.name]" :item="item" :key="index"></TaskFormItem>
              </Row>
            </Form>
          </Row>
          <div style="position:absolute; bottom:10px;text-align:center;width:95%">
            <Button style="margin-left: 15px" @click="detailCancel">{{$t('cancle')}}</Button>
          </div>
        </Form>
      </div>
    </Modal>
  </div>
</template>
<script>
import PluginTable from "../../components/table";
import {
  searchRequest,
  saveRequestInfo,
  getFormTemplateDetail,
  getTargetOptions,
  requestTemplateAvailable,
  getEntityDataByTemplateId,
  getTaskNodesEntitys,
  workflowProcessPrevieEntities,
  getRequestInfoDetails
} from "../../api/server.js"
export default {
  components: {
    PluginTable
  },
  data () {
    return {
      isQuery: true,
      requestStep: 1,
      requestDataObj: {},
      currentTaskNode: '',
      oldTaskNode: '',
      procTaskNodes: [],
      nodeData: [],
      targetModelColums: [
        {
          type: 'selection',
          width: 60,
          align: 'center'
        },
        {
          title: 'PackageName',
          key: 'packageName'
        },
        {
          title: 'EntityName',
          key: 'entityName'
        },
        {
          title: 'DisplayName',
          key: 'displayName'
        }
      ],
      currentFields: {},
      currentFieldsBackUp: {},
      targetEntityList: [],
      requestModalVisible: false,
      detailModalVisible: false,
      ruleValidate: {
        name: [
          {
            required: true,
            message: "The name cannot be empty",
            trigger: "blur"
          }
        ]
      },
      requestLoading: false,
      entityData: [],
      detailEntityData: [],
      allTemplates: [],
      allTemplatesTree: {},
       requestPagination: {
        currentPage: 1,
        pageSize: 10,
        total: 0
      },
      detailFields: {},
      detailForm: {
        name: "",
        emergency: "",
        description: "",
        attachFileId: null,
        requestTempId:'',
        roleName:'',
        rootEntity: ''
      },
      requestForm: {
        name: "",
        emergency: "",
        description: "",
        attachFileId: null,
        requestTempId:'',
        roleName:'',
        rootEntity: ''
      },
      requestPayload: {
        filters: {},
        pageable: {
          pageSize: 10,
          startIndex: 0
        },
        paging: true
      },
      currentForm: {},
      tableOuterActions: [
        {
          label: '发起请求上报',
          props: {
            type: "success",
            icon: "md-add",
            disabled: false
          },
          actionType: "add"
        }
      ],
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
          title: this.$t("status"),
          key: "status",
          inputKey: "status",
          component: "PluginSelect",
          inputType: "select",
          options: [
            {
              value: "Summitted",
              label: this.$t("summitted")
            },
            {
              value: "Processing",
              label: this.$t("processing")
            },
            {
              value: "Done",
              label: this.$t("done")
            }
          ]
        },
        {
          title: this.$t("emergency_level"),
          key: "emergency",
          inputKey: "emergency",
          component: "PluginSelect",
          options: [
            {
              value: "normal",
              label: this.$t("not_urgent")
            },
            {
              value: "urgent",
              label: this.$t("emergency")
            }
          ],
          inputType: "select"
        },
        {
          title: this.$t("reporter"),
          key: "reporter",
          inputKey: "reporter",
          component: "Input",
          inputType: "text",
          sortable: 'custom',
        },
        {
          title: this.$t("reporting_time"),
          key: "reportTime",
          inputKey: "reportTime",
          component: "DatePicker",
          type: "datetimerange",
          inputType: "date",
          sortable: 'custom',
          isNotFilterable: true
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
                    onClick={() => this.details(params.row)}
                  >
                    详情
                  </Button>
              </span>
            );
          }
        }
      ],
    }
  },
  watch: {
    currentFields: {
      handler(v) {

      },
      deep: true
    }
  },
  methods: {
    queryRequest () {
      this.isQuery = true;
      this.requestStep = 1
    },
    goBack () {
      this.requestStep--
    },
    async getTaskNodesEntitys (id) {
      const nodes = await getTaskNodesEntitys(id)
      this.procTaskNodes = nodes.data.filter(_ => _.taskCategory)
    },
    taskNodeChanged (v) {
      this.requestDataObj[this.oldTaskNode] = JSON.parse(JSON.stringify(this.nodeData))
      this.nodeData = this.requestDataObj[this.currentTaskNode]
      this.oldTaskNode = v
    },
    selectionChanged (selection) {
      this.nodeData.forEach(_ => {
        _._checked = false
        selection.forEach(d => {
          if (d.dataId === _.dataId) {
            _._checked = true
          }
        })
      })
    },
    addRowData (key) {
      const obj = JSON.parse(JSON.stringify(this.currentFields[key][0]))
      this.currentFields[key].push(obj)
    },
    async details (row) {
      const {status, message, data} = await getRequestInfoDetails(row.id)
      if (status === 'OK') {
        this.detailForm = {...data}
        this.detailFields['entitys'] = data.formItemInfos.sort(this.compare).map(i => {
            return {
              ...i,
              isMultiple: true,
              options: i.dataOptions && i.dataOptions.length > 0 ? JSON.parse(i.dataOptions) : []
            }
          })
          const found  = this.allTemplates.find(_ => _.id === this.detailForm.requestTempId)
        const entityData = await getEntityDataByTemplateId(found.procDefId) //getEntityDataByTemplateId
        this.detailEntityData = []
        if (entityData.status === "OK") {
          this.detailEntityData = entityData.data
        }
        const preview = await workflowProcessPrevieEntities(found.procDefId, this.detailForm.rootEntity)
        if (preview.status === "OK") {
          preview.data.entityTreeNodes.forEach(node => {
            const found = this.detailFields.entitys.find(entity => entity.name === node.entityName)
            if (found && !found.options) {
              found.options = [{...node,label:node.displayName,value:node.dataId}]
            } 
            if (found && found.options) {
              found.options.push({...node,label:node.displayName,value:node.dataId})
            }
            if (this.currentForm[node.entityName]) {this.currentForm[node.entityName].push(node.dataId)}
          })
        }
        this.detailModalVisible = true
      }

    },
    async workflowProcessPrevieEntities (v) {
      // workflowProcessPrevieEntities
      if (v) {
        const processKey = this.allTemplates.find(_ => _.id === this.requestForm.requestTempId).procDefId
        // const processKey = 'sjqH9YVJ2DP'
        const {status, message, data} = await workflowProcessPrevieEntities(processKey, v)
        if (status === 'OK') {
          Object.keys(this.currentForm).forEach(field => {
            this.currentForm[field] = []
          })
          let formatData = {}
          data.entityTreeNodes.forEach(node => {
            const found = this.currentFields[node.entityName]
            if (found) {
              const foundData = formatData[node.entityName]
              if (foundData) {
                formatData[node.entityName].push({
                  entity: node.entityName,
                  entityData: {...node},
                  fields: this.currentFields[node.entityName][0].fields.map(field => {
                    return {
                      ...field,
                      value: node.entityData[field.attribute.name]
                    }
                  })
                })
              } else {
                formatData[node.entityName]=[{
                  entity: node.entityName,
                  entityData: {...node},
                  fields: this.currentFields[node.entityName][0].fields.map(field => {
                    return {
                      ...field,
                      value: node.entityData[field.attribute.name]
                    }
                  })
                }]
              }
            }
            if (this.currentForm[node.entityName]) {this.currentForm[node.entityName].push(node.dataId)}
          })
          formatData.customize = this.currentFields.customize
          this.currentFieldsBackUp = JSON.parse(JSON.stringify(formatData))
          this.$nextTick(() => {
            this.currentFields = {}
            this.currentFields = JSON.parse(JSON.stringify(this.currentFieldsBackUp))
          })
        }
      }
    },
    actionFun(type, data) {
      switch (type) {
        case "add":
          // this.requestModalVisible = true;
          this.isQuery = false;
          this.currentFields = {}
          break;
      }
    },
    uploadSuccess(res, file, fileList) {
      this.requestForm.attachFileId = res.data;
    },
    async getTemplates() {
      const { data } = await requestTemplateAvailable();
      this.allTemplates = data;
      // this.allTemplatesTree
      let treeData = []
      data.forEach(_ => {
        const found = treeData.find(t => t.requestTempGroup === _.requestTempGroup)
        if (found) {
          const foundTag = found.children.find(f => f.tag === _.tags)
          if (foundTag) {
            foundTag.children.push(_)
          } else {
            found.children.push({tag: _.tags, children:[{..._}]})
          }
        } else {
          treeData.push({
            requestTempGroup: _.requestTempGroup,
            requestTempGroupName: _.requestTempGroupName,
            children:[{
              tag: _.tags,
              children: [{..._}]
            }]
          })
        }
      })
      this.allTemplatesTree = treeData
    },
    compare (a, b) {
        if (a.sort < b.sort) {
          return -1
        }
        if (a.sort > b.sort) {
          return 1
        }
        return 0
      },
    async templateChanged (checked,v) {
      if (v) {
        this.requestForm.requestTempId = v
        this.entityData = []
        this.requestForm.rootEntity = ""
        this.currentFields = {}
        this.currentFields.customize = [{
          entity: 'customize',
          fields: []
        }]
        const procDefId = this.allTemplates.find(_ => _.id === this.requestForm.requestTempId).procDefId
        this.getTaskNodesEntitys(procDefId)
        const {status, message, data} = await getFormTemplateDetail(0, v)
        if (status === 'OK') {
          //currentFields
          this.targetEntityList = []
          let entity = new Set()
          data.items.forEach(item => {
            item.entity && item.entity.length > 0 && entity.add(item.entity) // 按对象排序
            if (this.currentFields[item.entity]) {
              this.currentFields[item.entity][0].fields.push({
                attribute: {
                  ...item,
                  options: item.dataOptions > 0 ? JSON.parse(item.dataOptions) : []
                },
                // isMultiple: true,
                value: ''
              })
            } else {
              if (item.entity && item.entity.length > 0) {
                this.currentFields[item.entity] = [{
                  entity: item.entity,
                  fields: [{
                    attribute: {
                      ...item,
                      options: item.dataOptions > 0 ? JSON.parse(item.dataOptions) : []
                    },
                    // isMultiple: true,
                    value: ''
                  }]
                }]
                
              } else {
                this.currentFields['customize'][0].fields.push({
                  attribute: {
                    ...item,
                    options: item.dataOptions > 0 ? JSON.parse(item.dataOptions) : []
                  },
                  // isMultiple: true,
                  value: ''
                })
              }
            }
          })
          this.targetEntityList = Array.from(entity)
          this.targetEntityList.push('customize')
          data.items.forEach(_ => {
            this.currentForm[_.name] = []
          })
          this.currentFieldsBackUp = JSON.parse(JSON.stringify(this.currentFields))
          this.$nextTick(() => {
            this.currentFields = {}
            this.currentFields = JSON.parse(JSON.stringify(this.currentFieldsBackUp))
            this.requestStep++
          })
        }
      }
    },
    async getEntityDataByTemplateId (v) {
      if (v && this.requestForm.requestTempId && this.requestForm.requestTempId.length > 0) {
        const found  = this.allTemplates.find(_ => _.id === this.requestForm.requestTempId)
        const { data, status } = await getEntityDataByTemplateId(found.procDefId) //getEntityDataByTemplateId
        this.entityData = []
        if (status === "OK") {
          this.entityData = data
        }
      }
    },
    requestModalHide() {
      this.requestModalVisible = false;
      this.currentFields = {}
    },
    detailCancel () {
      this.detailModalVisible = false;
    },
    generateUUID() {
    let d = new Date().getTime();
      if (window.performance && typeof window.performance.now === "function") {
          d += performance.now(); //use high-precision timer if available
      }
      let uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
          let r = (d + Math.random() * 16) % 16 | 0;
          d = Math.floor(d / 16);
          return (c == 'x' ? r : (r & 0x3 | 0x8)).toString(16);
      });
      return uuid;
    },
    goNextStep () {
      // requestDataObj
      this.requestDataObj = {}
      this.procTaskNodes.forEach(node => {
        this.currentFields[node.boundEntity.name].forEach(row => {
          const data = {
            nodeId: node.nodeId,
            nodeDefId: node.nodeDefId,
            nodeName: node.nodeName,
            ...row.entityData,
            oid: row.entityData? row.entityData.dataId : this.generateUUID(),
            _checked: true
          }
          delete data.entityData
          const attrValues = []
          row.fields.forEach(field => {
            const obj = {
              attrDefId: field.attribute.attrDefId,
              dataType: field.attribute.dataType,
              dataValue: field.value,
              itemTempId: field.attribute.tempId,
              name: field.attribute.name
            }
            attrValues.push(obj)
          })
          data.attrValues = attrValues
          if (this.requestDataObj[node.nodeName]) {
            this.requestDataObj[node.nodeName].push(data)
          } else {
            this.requestDataObj[node.nodeName] = [data]
          }
        })
      })
      this.currentTaskNode = this.procTaskNodes[0].nodeName
      this.oldTaskNode = this.currentTaskNode
      this.nodeData = this.requestDataObj[this.currentTaskNode]
      this.requestStep++
    },
    requestCancel() {
      this.requestStep = 1
      this.requestModalVisible = false;
      this.requestForm.name = "";
      this.requestForm.emergency = "";
      this.requestForm.description = "";
      this.requestForm.requestTempId = "";
      this.requestForm.roleId = "";
      this.requestForm.rootEntity = "";
    },
    requestSubmit() {
      let formItems = [].concat(...Object.values(this.requestDataObj))
      const payload = {
        ...this.requestForm,
        entities: formItems
      }
      this.$refs.requestForm.validate(async valid => {
        if (valid) {
          this.requestLoading = true
          const { status } = await saveRequestInfo(payload);
          this.requestLoading = false
          if (status === "OK") {
            this.requestCancel();
            this.getData();
            this.requestStep = 1
            this.isQuery = true
            this.requestForm.attachFileId = null;
            // this.$refs.upload.clearFiles();
            this.$Notice.success({
              title: 'Success',
              desc: 'Success'
            })
          }
        }
      });
    },
    handleSubmit(filters) {
      this.requestPayload.filters = filters;
      this.getData();
    },
    async getData () {
      const f = this.requestPayload.filters
      const payload = {
        page: this.requestPagination.currentPage,
        pageSize: this.requestPagination.pageSize,
        data: f
      }
      const { data, status} = await searchRequest(payload)
      if(status === 'OK') {
        this.requestTableData = data.contents
        this.requestPagination.total = data.pageInfo.totalRows;
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
  },
  mounted() {
    this.getData();
    this.getTemplates()
  }
}
</script>
<style lang="scss" scoped>

</style>


