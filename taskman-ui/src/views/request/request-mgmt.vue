<template>
  <div>
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
    <Modal
      v-model="requestModalVisible"
      :title="$t('request_to_report')"
      footer-hide
      fullscreen
      @on-cancel="requestModalHide"
    >
      <div style="padding:20px">
        <Form ref="requestForm" :rules="ruleValidate" :model="requestForm" :label-width="110">
          <Row>
            <Col span="12">
          <FormItem :label="$t('template')">
            <Select filterable @on-open-change="getTemplates" @on-change="templateChanged" v-model="requestForm.requestTempId">
              <Option v-for="tem in allTemplates" :key="tem.id" :value="tem.id">{{tem.name}}</Option>
            </Select>
          </FormItem>
          </Col>
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
          <hr style="margin-bottom:10px"/>
          <Row>
            <Form ref="form" :label-width="110">
              <Row v-for="(fields, i) in Object.values(currentFields)" :key="i">
                <TaskFormItem v-for="(item, index) in fields" :index="index" v-model="currentForm[item.name]" :item="item" :key="index"></TaskFormItem>
              </Row>
            </Form>
          </Row>
          <div style="position:absolute; bottom:10px;text-align:center;width:95%">
            <Button type="primary" :loading="requestLoading" @click="requestSubmit">{{$t('submit')}}</Button>
            <Button style="margin-left: 15px" @click="requestCancel">{{$t('cancle')}}</Button>
          </div>
        </Form>
      </div>
    </Modal>
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
  workflowProcessPrevieEntities,
  getRequestInfoDetails
} from "../../api/server.js"
export default {
  components: {
    PluginTable
  },
  data () {
    return {
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
          label: this.$t("add"),
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
          this.currentFields.entitys.filter(en => en.entity && en.entity.length > 0).forEach(entity => entity.options = [])
          Object.keys(this.currentForm).forEach(field => {
            this.currentForm[field] = []
          })
          data.entityTreeNodes.forEach(node => {
            const found = this.currentFields.entitys.find(entity => entity.name === node.entityName)
            if (found && !found.options) {
              found.options = [{...node,label:node.displayName,value:node.dataId}]
            } 
            if (found && found.options) {
              found.options.push({...node,label:node.displayName,value:node.dataId})
            }
            if (this.currentForm[node.entityName]) {this.currentForm[node.entityName].push(node.dataId)}
          })
          this.currentFieldsBackUp = JSON.parse(JSON.stringify(this.currentFields))
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
          this.requestModalVisible = true;
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
    async templateChanged (v) {
      if (v) {
        this.entityData = []
        this.requestForm.rootEntity = ""
        this.currentFields = {}
        const {status, message, data} = await getFormTemplateDetail(0, v)
        if (status === 'OK') {
          //currentFields
          this.targetEntityList = JSON.parse(data.targetEntitys)
          // data.items.forEach(item => {
          //   if (this.currentFields[item.entity]) {
          //     this.currentFields[item.entity].push({
          //       ...item,
          //       options: item.dataOptions.length > 0 ? JSON.parse(item.dataOptions) : []
          //     })
          //   } else {
          //     this.currentFields[item.entity] = [{
          //       ...item,
          //       options: item.dataOptions.length > 0 ? JSON.parse(item.dataOptions) : []
          //     }]
          //   }
          // })
          data.items.forEach(_ => {
            this.currentForm[_.name] = []
          })
          this.currentFields['entitys'] = data.items.sort(this.compare).map(i => {
            return {
              ...i,
              isMultiple: true,
              options: i.dataOptions.length > 0 ? JSON.parse(i.dataOptions) : []
            }
          })
          this.currentFieldsBackUp = JSON.parse(JSON.stringify(this.currentFields))
          this.$nextTick(() => {
            this.currentFields = {}
            this.currentFields = JSON.parse(JSON.stringify(this.currentFieldsBackUp))
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
    requestCancel() {
      this.requestModalVisible = false;
      this.requestForm.name = "";
      this.requestForm.emergency = "";
      this.requestForm.description = "";
      this.requestForm.requestTempId = "";
      this.requestForm.roleId = "";
      this.requestForm.rootEntity = "";
    },
    requestSubmit() {
      let formItems = []
      Object.keys(this.currentForm).forEach(_ => {
        const found = this.currentFields.entitys.find(field => field.name === _)
        formItems.push({
          itemTempId: found.id,
          name: _,
          value: JSON.stringify(this.currentForm[_])
        })
      })
      const payload = {
        ...this.requestForm,
        formItems: formItems
      }
      this.$refs.requestForm.validate(async valid => {
        if (valid) {
          this.requestLoading = true
          const { status } = await saveRequestInfo(payload);
          this.requestLoading = false
          if (status === "OK") {
            this.requestCancel();
            this.getData();
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


