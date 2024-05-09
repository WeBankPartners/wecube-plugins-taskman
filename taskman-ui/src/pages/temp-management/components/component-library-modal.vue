<template>
  <Modal
    class="component-library-dialog"
    v-model="visible"
    title=""
    :mask-closable="false"
    :closable="false"
    :footer-hide="true"
    width="1200"
  >
    <div class="content">
      <div v-if="isAdd" class="left">
        <span class="title">新建组件</span>
        <Form :label-width="100">
          <FormItem label="组件名：">
            <Input v-model="form.name" />
          </FormItem>
          <FormItem label="表单类型：">
            <span class="form-text">表单类型</span>
          </FormItem>
          <FormItem label="表单项：">
            <span class="form-text">表单项1、表单项2、表单项3</span>
          </FormItem>
        </Form>
        <Button type="primary" @click="handleSave" style="float: right;" :disabled="!form.name">保存</Button>
      </div>
      <div class="right">
        <span class="title">组件列表</span>
        <div class="query">
          <Select v-model="query.type" placeholder="表单类型" style="width: 150px;">
            <Option v-for="(i, index) in typeList" :value="i" :key="index">{{ i }}</Option>
          </Select>
          <Select v-model="query.createdBy" placeholder="创建人" style="width: 150px;">
            <Option v-for="(i, index) in typeList" :value="i" :key="index">{{ i }}</Option>
          </Select>
          <Input v-model="query.name" placeholder="组件名" style="width: 300px;" />
        </div>
        <Table
          style="width:100%;margin-top:20px;"
          :border="false"
          size="small"
          :columns="tableColumns"
          :data="tableData"
        />
      </div>
      <Icon type="md-close" class="close" :size="24" @click="handleClose" />
    </div>
  </Modal>
</template>
<script>
export default {
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    isAdd: {
      type: Boolean,
      default: true
    }
  },
  data () {
    return {
      form: {
        name: ''
      },
      query: {
        type: '',
        createdBy: '',
        name: ''
      },
      typeList: [1, 2, 3],
      tableData: [],
      tableColumns: [
        {
          title: '组件名',
          key: 'name',
          align: 'left',
          minWidth: 200
        },
        {
          title: '表单类型',
          key: 'type',
          align: 'left',
          minWidth: 150
        },
        {
          title: '表单项',
          key: 'type',
          align: 'left',
          minWidth: 150
        },
        {
          title: '创建人',
          key: 'type',
          align: 'left',
          minWidth: 150
        },
        {
          title: this.$t('t_action'),
          key: 'action',
          align: 'center',
          fixed: 'right',
          width: 80,
          render: (h, params) => {
            return (
              <Button
                on-click={() => {
                  this.handleDelete(params.row)
                }}
                type="error"
                size="small"
                ghost
                icon="md-trash"
              ></Button>
            )
          }
        }
      ]
    }
  },
  methods: {
    handleSave () {},
    handleDelete (row) {},
    handleClose () {
      this.$emit('update:visible', false)
    }
  }
}
</script>

<style lang="scss" scoped>
.component-library-dialog {
  .content {
    display: flex;
    position: relative;
    .left {
      width: 360px;
      padding-right: 20px;
      border-right: 1px solid #e8eaec;
      .form-text {
        font-size: 14px;
        color: #515a6e;
      }
    }
    .right {
      flex: 1;
      padding-left: 20px;
    }
    .close {
      position: absolute;
      right: 0px;
      top: 0px;
      cursor: pointer;
    }
  }
  .title {
    display: block;
    font-size: 16px;
    color: #17233d;
    font-weight: 500;
    margin-bottom: 10px;
  }
}
</style>
