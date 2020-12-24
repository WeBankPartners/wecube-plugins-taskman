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
          <Row>
            <Form ref="form" :label-width="110">
              <Row>
                <TaskFormItem v-for="(item, index) in currentFields" :index="index" v-model="currentForm[item.name]" :item="item" :key="index"></TaskFormItem>
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
  searchTask,
  getTaskInfoDetails,
  taskInfoReceive,
  taskInfoProcessing,
  getTaskInfoInstance
} from "../../api/server.js"
export default {
  components: {
    PluginTable
  },
  data () {
    return {
      currentFields: [],
      currentTaskId: '',
      currentFieldsBackUp: {},
      requestModalVisible: false,
      requestLoading: false,
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
      tableOuterActions: [],
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
                    onClick={() => this.receiveTask(params.row)}
                  >
                    领取
                  </Button>
                  <Button
                    type="primary"
                    size="small"
                    style="margin-right:10px"
                    onClick={() => this.processTask(params.row)}
                  >
                    处理
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
    async receiveTask (row) {
      const {status} = await taskInfoReceive({id: row.id});
      if (status === "OK") {
        this.getData()
      }
    },
    async processTask (row) {
      const { status, data } = await getTaskInfoDetails({id: row.id})
      if (status === "OK") {
        this.requestModalVisible = true;
        this.currentTaskId = row.id;
        this.currentFields = data.formItemInfo   //currentForm
        this.currentFields.forEach(field => {
          this.currentForm[field.name] = ''
        })
      }
    },
    actionFun(type, data) {
      switch (type) {
        case "add":
          
          break;
      }
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
    requestModalHide() {
      this.requestModalVisible = false;
      this.currentFields = []
      this.currentForm = []
    },
    requestCancel() {
      this.requestModalVisible = false;
      this.currentForm = []
    },
    async requestSubmit() {
      const payload = {
        formItemInfoList:[],
        recordId: ''
      }
      Object.keys(this.currentForm).forEach(key => {
        const found = this.currentFields.find(field => field.name === key)
        payload.formItemInfoList.push({
          itemTempId: found.itemTempId,
          name:key,
          value:this.currentForm[key]
        })
      })
      const {status} = await taskInfoProcessing(payload)
      if (status === 'OK') {
        this.requestModalHide()
        this.getData()
        this.$Notice.success({
              title: 'Success',
              desc: 'Success'
            })
      }
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
      const { data, status} = await searchTask(payload)
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


