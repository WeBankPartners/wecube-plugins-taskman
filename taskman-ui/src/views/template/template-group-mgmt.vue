<template>
  <div>
    <PluginTable
      :tableColumns="requestColumns"
      :tableData="requestTableData"
      :tableOuterActions="tableOuterActions"
      :pagination="requestPagination"
      @actionFun="actionFun"
      @handleSubmit="handleSubmit"
      @sortHandler="sortHandler"
      @pageChange="requestPageChange"
      @pageSizeChange="requestPageSizeChange"
    />
    <Modal
      v-model="requestModalVisible"
      :title="currentTitle"
      footer-hide
      width="50"
      @on-cancel="requestModalHide"
    >
      <div style="width:600px;margin:0 auto;">
        <Form ref="requestForm" :model="requestForm" :label-width="110">
          <FormItem :label="$t('name')" prop="name">
            <Input v-model="requestForm.name" :placeholder="$t('service_request_name')"></Input>
          </FormItem>
          <FormItem :label="$t('service_request_role')">
            <Select v-model="requestForm.manageRoleId">
              <Option
                v-for="role in currentUserRoles"
                :key="role.name"
                :value="role.name"
              >{{role.displayName}}</Option>
            </Select>
          </FormItem>
          <FormItem :label="$t('describe')">
            <Input type="textarea" v-model="requestForm.description" :placeholder="$t('describe')"></Input>
          </FormItem>
          <FormItem>
            <Button type="primary" :loading="requestLoading" @click="requestSubmit">{{$t('submit')}}</Button>
            <Button style="margin-left: 8px" @click="requestCancel">{{$t('cancle')}}</Button>
          </FormItem>
        </Form>
      </div>
    </Modal>
  </div>
</template>
<script>
import {
  addRequestTemplateGroup,
  deleteRequestTemplateGroup,
  editRequestTemplateGroup,
  searchRequestTemplateGroup,
  getRoleList,
} from "../../api/server";
import PluginTable from "../../components/table";
export default {
  components: {
    PluginTable
  },
  data() {
    return {
      currentTitle: '',
      allTemplates: [],
      currentUserRoles: [],
      entityData: [],
      requestModalVisible: false,
      requestLoading: false,
      requestForm: {
        name: "",
        emergency: "",
        description: "",
        attachFileId: null,
        templateId:'',
        manageRoleId:'',
        rootDataId: ''
      },
      modelForm: {
        createdBy: "",
        description: "",
        id: "",
        manageRole: "",
        name: "",
        updatedBy: ""
      },
      requestColumns: [
        {
          title: this.$t("name"),
          key: "name",
          inputKey: "name",
          component: "Input",
          inputType: "text",
          sortable: 'custom',
        },
        // {
        //   title: this.$t("reporter"),
        //   key: "reporter",
        //   inputKey: "reporter",
        //   component: "Input",
        //   inputType: "text",
        //   sortable: 'custom',
        // },
        {
          title: this.$t("describe"),
          key: "description",
          inputKey: "description",
          component: "Input",
          inputType: "text"
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
                  <Button
                    type="error"
                    size="small"
                    onClick={() => this.deleteGroup(params.row)}
                  >
                    {this.$t("delete")}
                  </Button>
              </span>
            );
          }
        }
      ],
      requestTableData: [],
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
      requestPayload: {
        filters: {},
        pageable: {
          pageSize: 10,
          startIndex: 0
        },
        paging: true
      },
      requestPagination: {
        currentPage: 1,
        pageSize: 10,
        total: 0
      },
      isAdd: false
    }
  },
  methods: {
    actionFun(type, data) {
      switch (type) {
        case "add":
          this.requestModalVisible = true;
          this.isAdd = true;
          this.currentTitle = '新增'
          this.requestForm = {
            name: "",
            emergency: "",
            description: "",
            attachFileId: null,
            templateId:'',
            manageRoleId:'',
            rootDataId: ''
          }
          break;
      }
    },

    async getRoleList () {
      const {status, message, data} = await getRoleList()
      if (status === 'OK') {
        this.currentUserRoles = data
      }
    },

    requestModalHide() {
      this.requestModalVisible = false;
    },

    requestCancel() {
      this.requestModalVisible = false;
      this.requestForm.name = "";
      this.requestForm.description = "";
      this.requestForm.manageRoleId = "";
    },
    deleteGroup (row) {
      this.$Modal.warning({
          title: 'Warning',
          content: 'Delete ?',
          onOk: async () => {
            const {status} = await deleteRequestTemplateGroup(row.id);
            if (status === "OK") {
              this.$Notice.success({
              title: 'Success',
              desc: 'Success'
            })
            this.getData();
            }
          },
          onCancel: () => {}
      });
    },
    edit (row) {
      this.isAdd = false;
      this.currentTitle = '编辑'
      this.requestForm = {...row, manageRoleId: row.manageRole}
      this.requestModalVisible = true;
    },
    requestSubmit() {
      const payload = {
        ...this.requestForm,
        manageRoleName: this.currentUserRoles.find(_ => _.name == this.requestForm.manageRoleId).displayName
      }
      this.$refs.requestForm.validate(async valid => {
        if (valid) {
          this.requestLoading = true
          const { status } = await addRequestTemplateGroup(payload)
          this.requestLoading = false
          if (status === "OK") {
            this.requestCancel();
            this.getData();
            this.$Notice.success({
              title: 'Success',
              desc: 'Success'
            })
          }
        }
      });
    },

    async getData () {
      const payload = {
        page: this.requestPagination.currentPage,
        pageSize: this.requestPagination.pageSize,
        data: this.requestPayload.filters
      }
      const { data, status} = await searchRequestTemplateGroup(payload)
      if(status === 'OK') {
        this.requestTableData = data.contents
        this.requestPagination.total = data.pageInfo.totalRows;
      }
    },

    handleSubmit(filters) {
      this.requestPayload.filters = filters;
      this.getData();
    },
    requestPageChange(current) {
      this.requestPagination.currentPage = current;
      this.getData();
    },
    requestPageSizeChange(size) {
      this.requestPagination.pageSize = size;
      this.getData();
    },
    sortHandler (data) {
      if (data.order === 'normal') {
        delete this.requestPayload.sorting
      } else {
        this.requestPayload.sorting = {
          asc: data.order === 'asc',
          field: data.key
        }
      }
      // this.getData()
    },
  },
  mounted() {
    this.getData()
    this.getRoleList()
  }
  }
</script>
<style lang="scss" scoped>

</style>