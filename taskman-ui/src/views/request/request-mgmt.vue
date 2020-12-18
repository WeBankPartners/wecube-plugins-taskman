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
            <Select filterable @on-open-change="getTemplates" @on-change="templateChanged" v-model="requestForm.templateId">
              <Option v-for="tem in allTemplates" :key="tem.id" :value="tem.id">{{tem.name}}</Option>
            </Select>
          </FormItem>
          </Col>
          <Col span="12">
          <FormItem :label="$t('target_object')">
            <Select filterable @on-open-change="getEntityDataByTemplateId" v-model="requestForm.rootDataId">
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
          <Row>
            <Form ref="form" :label-width="110">
              <Row v-for="(fields, i) in Object.values(currentFields)" :key="i">
                <TaskFormItem v-for="(item, index) in fields" :index="index" :item="item" :key="index"></TaskFormItem>
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
  getEntityDataByTemplateId
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
      allTemplates: [],
       requestPagination: {
        currentPage: 1,
        pageSize: 10,
        total: 0
      },
      requestForm: {
        name: "",
        emergency: "",
        description: "",
        attachFileId: null,
        templateId:'',
        roleName:'',
        rootDataId: ''
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
  methods: {

    details (row) {

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
        this.requestForm.rootDataId = ""
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
          this.currentFields['entitys'] = data.items.sort(this.compare)
          this.currentFieldsBackUp = JSON.parse(JSON.stringify(this.currentFields))
          this.$nextTick(() => {
            this.currentFields = {}
            this.currentFields = JSON.parse(JSON.stringify(this.currentFieldsBackUp))
          })
        }
      }
    },
    async getEntityDataByTemplateId (v) {
      if (v && this.requestForm.templateId && this.requestForm.templateId.length > 0) {
        const found  = this.allTemplates.find(_ => _.id === this.requestForm.templateId)
        const { data, status } = await getEntityDataByTemplateId(found.procDefKey) //getEntityDataByTemplateId
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
    requestCancel() {
      this.requestModalVisible = false;
      this.requestForm.name = "";
      this.requestForm.emergency = "";
      this.requestForm.description = "";
      this.requestForm.templateId = "";
      this.requestForm.roleId = "";
      this.requestForm.rootDataId = "";
    },
    requestSubmit() {
      this.$refs.requestForm.validate(async valid => {
        if (valid) {
          this.requestLoading = true
          const { status } = await createServiceRequest(this.requestForm);
          this.requestLoading = false
          if (status === "OK") {
            this.requestCancel();
            this.getData();
            this.requestForm.attachFileId = null;
            this.$refs.upload.clearFiles();
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
  }
}
</script>
<style lang="scss" scoped>

</style>


