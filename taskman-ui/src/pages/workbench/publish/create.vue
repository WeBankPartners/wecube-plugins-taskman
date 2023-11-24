<template>
  <div class="workbench-publish-create">
    <Row class="w-header">
      <Col span="18" class="steps">
        <span class="title">请求进度：</span>
        <Steps :current="0" style="max-width:500px;">
          <Step v-for="(i, index) in steps" :key="index" :content="i.name">
            <template #icon>
              <Icon size="26" :type="i.icon" :color="i.color" />
            </template>
            <div class="role" slot="content">
              <span>{{ i.name }}</span>
              <span>{{ '管理员' }}</span>
            </div>
          </Step>
        </Steps>
      </Col>
      <Col span="6" class="btn-group">
        <Button @click="handleDraft" style="margin-right:10px;">保存草稿</Button>
        <Button type="primary" @click="handlePublish">发布</Button>
      </Col>
    </Row>
    <Row class="content">
      <Col span="16" class="split-line">
        <Form :model="form" label-position="right" :label-width="120">
          <HeaderTitle title="发布信息">
            <FormItem label="请求名称" required>
              <Input v-model="form.name" placeholder="请输入" style="width:400px;" />
            </FormItem>
            <FormItem label="发布描述">
              <Input v-model="form.description" placeholder="请输入" style="width:400px;" />
            </FormItem>
          </HeaderTitle>
          <HeaderTitle title="发布目标对象">
            <FormItem label="选择操作单元" required>
              <Select v-model="form.rootEntityId" clearable filterable style="width:300px;" @on-change="getEntityData">
                <Option v-for="item in rootEntityOptions" :value="item.guid" :key="item.guid">{{
                  item.key_name
                }}</Option>
              </Select>
            </FormItem>
            <FormItem v-if="requestData.length" label="已选择">
              <RadioGroup v-model="activeTab" @on-change="handleTabChange" style="margin-bottom:20px;">
                <Radio v-for="(item, index) in requestData" :label="item.entity || item.itemGroup" :key="index" border>
                  <span
                    >{{ `${item.itemGroup}` }}<span class="count">{{ item.value.length }}</span></span
                  >
                </Radio>
              </RadioGroup>
              <Table size="small" :columns="tableColumns" :data="tableData"></Table>
            </FormItem>
          </HeaderTitle>
        </Form>
      </Col>
      <!--编排流程-->
      <Col span="8">
        <StaticFlow></StaticFlow>
      </Col>
    </Row>
    <EditDrawer
      v-if="editVisible"
      v-model="editData"
      :options="editOptions"
      :visible.sync="editVisible"
      :disabled="viewDisabled"
      @submit="submitEditDrawer"
    ></EditDrawer>
  </div>
</template>

<script>
import HeaderTitle from '../components/header-title.vue'
import EditDrawer from './edit-item.vue'
import { createRequest, getRootEntity, getEntityData, getRefOptions } from '@/api/server'
import StaticFlow from './flow/static-flow.vue'
import { deepClone, debounce } from '@/pages/util'
export default {
  components: {
    HeaderTitle,
    StaticFlow,
    EditDrawer
  },
  data () {
    return {
      templateId: '',
      requestId: '6557454d5324718d',
      activeTab: '',
      refKeys: [], // 引用类型字段集合select类型
      form: {
        name: '',
        description: '',
        rootEntityId: '', // 目标对象
        exampleType: 1
      },
      rootEntityOptions: [],
      steps: [
        { name: '提起请求', status: 'process', icon: 'md-pin', color: '#ffa500' },
        { name: '请求定版', status: 'wait', icon: 'md-radio-button-on', color: '#8189a5' },
        { name: '任务1审批', status: 'wait', icon: 'md-radio-button-on', color: '#8189a5' },
        { name: '任务2审批', status: 'wait', icon: 'md-radio-button-on', color: '#8189a5' },
        { name: '请求完成', status: 'wait', icon: 'md-radio-button-on', color: '#8189a5' }
      ],
      tableColumns: [],
      tableData: [], // 用于当前表格数据的展示
      requestData: [], // 用于最后提交的所有表格数据
      editVisible: false,
      editOptions: [],
      editData: {},
      editIndex: 0,
      viewDisabled: false,
      initTableColumns: [
        {
          type: 'selection',
          width: 55,
          align: 'center',
          fixed: 'left'
        },
        {
          title: this.$t('t_action'),
          key: 'action',
          width: 120,
          fixed: 'right',
          align: 'center',
          render: (h, params) => {
            return (
              <div style="display:flex;justify-content:space-around;">
                <Tooltip content="删除" placement="top-start">
                  <Icon
                    size="20"
                    type="md-trash"
                    color="#ed4014"
                    style="cursor:pointer;"
                    onClick={() => {
                      this.handleDeleteRow(params.row)
                    }}
                  />
                </Tooltip>
                <Tooltip content="编辑" placement="top-start">
                  <Icon
                    size="20"
                    type="md-create"
                    style="cursor:pointer;"
                    onClick={() => {
                      this.handleEditRow(params.row)
                    }}
                  />
                </Tooltip>
                <Tooltip content="查看" placement="top-start">
                  <Icon
                    size="20"
                    type="md-eye"
                    style="cursor:pointer;"
                    onClick={() => {
                      this.handleViewRow(params.row)
                    }}
                  />
                </Tooltip>
              </div>
            )
          }
        }
      ]
    }
  },
  async mounted () {
    this.templateId = this.$route.query.templateId || ''
    if (this.templateId) {
      await this.getCreateInfo()
    }
    this.getEntity()
    this.getEntityData()
  },
  methods: {
    // 切换tab刷新表格数据，加上防抖避免切换过快显示异常问题
    handleTabChange: debounce(function () {
      this.initTableData()
    }, 100),
    // 创建发布,使用模板ID获取详情数据
    async getCreateInfo () {
      const params = {
        templateId: this.templateId
      }
      const { statusCode, data } = await createRequest(params)
      if (statusCode === 'OK') {
        this.requestId = data.id
        this.form.name = data.name
      }
    },
    // 操作目标对象下拉值
    async getEntity () {
      let params = {
        params: {
          requestId: this.requestId
        }
      }
      const { statusCode, data } = await getRootEntity(params)
      if (statusCode === 'OK') {
        this.rootEntityOptions = data.data
      }
    },
    // 获取目标对象对应数据
    async getEntityData () {
      let params = {
        params: {
          requestId: this.requestId,
          rootEntityId: this.form.rootEntityId
        }
      }
      const { statusCode, data } = await getEntityData(params)
      if (statusCode === 'OK') {
        this.requestData = data.data
        this.activeTab = this.activeTab || data.data[0].entity
        this.initTableData()
      }
    },
    // 编辑操作，刷新requestData
    refreshRequestData () {
      this.requestData.forEach(item => {
        if (item.entity === this.activeTab || item.itemGroup === this.activeTab) {
          for (let m of item.value) {
            for (let n of this.tableData) {
              if ((m.id = n._id)) {
                m.entityData = n
              }
            }
          }
        }
      })
    },
    async initTableData () {
      // 当前选择tab数据
      const data = this.requestData.find(r => r.entity === this.activeTab || r.itemGroup === this.activeTab)
      this.oriData = data
      // 编辑表单的options配置
      this.editOptions = data.title

      // select类型集合
      this.refKeys = []
      data.title.forEach(t => {
        if (t.elementType === 'select') {
          this.refKeys.push(t.name)
        }
      })

      // tableColumns数据初始化
      this.tableColumns = deepClone(this.initTableColumns)
      data.title.forEach(t => {
        let column = {
          title: t.title,
          key: t.name,
          align: 'left'
        }
        if (t.required === 'yes') {
          column.renderHeader = (h, { column }) => {
            return (
              <span>
                {`${column.title}`}
                <span class="required">（必填）</span>
              </span>
            )
          }
        }
        if (t.elementType === 'select') {
          column.render = (h, params) => {
            return (
              <Select
                value={params.row[t.name]}
                on-on-change={v => {
                  this.tableData[params.row._index][t.name] = v
                  this.refreshRequestData()
                }}
                multiple={t.multiple === 'Y'}
                disabled={t.isEdit === 'no'}
              >
                {Array.isArray(params.row[t.name + 'Options']) &&
                  params.row[t.name + 'Options'].map(i => (
                    <Option value={t.entity ? i.guid : i} key={t.entity ? i.guid : i}>
                      {t.entity ? i.key_name : i}
                    </Option>
                  ))}
              </Select>
            )
          }
        } else if (t.elementType === 'input') {
          column.render = (h, params) => {
            return (
              <Input
                value={params.row[t.name]}
                onInput={v => {
                  params.row[t.name] = v
                  // 暂时这么写,为啥给params赋值不会更新tableData？
                  this.tableData[params.row._index][t.name] = v
                }}
                onBlur={() => {
                  this.refreshRequestData()
                }}
                disabled={t.isEdit === 'no'}
              />
            )
          }
        } else if (t.elementType === 'textarea') {
          column.render = (h, params) => {
            return (
              <Input
                value={params.row[t.name]}
                onInput={v => {
                  params.row[t.name] = v
                  this.tableData[params.row._index][t.name] = v
                }}
                onBlur={() => {
                  this.refreshRequestData()
                }}
                type="textarea"
                disabled={t.isEdit === 'no'}
              />
            )
          }
        }
        this.tableColumns.push(column)
      })

      // table数据初始化
      this.tableData = data.value.map(v => {
        this.refKeys.forEach(rfk => {
          v.entityData[rfk + 'Options'] = []
        })
        v.entityData._id = v.id
        return v.entityData
      })

      // 下拉类型数据初始化
      this.tableData.forEach((row, index) => {
        this.refKeys.forEach(rfk => {
          const titleObj = data.title.find(f => f.name === rfk)
          this.getRefOptions(titleObj, row, index)
        })
      })
    },
    async getRefOptions (titleObj, row, index) {
      if (titleObj.elementType === 'select' && titleObj.entity === '') {
        row[titleObj.name + 'Options'] = titleObj.dataOptions.split(',')
        this.$set(this.tableData, index, row)
        return
      }
      // if (titleObj.refEntity === '') {
      //   row[titleObj.name + 'Options'] = titleObj.selectList
      //   this.$set(this.tableData, index, row)
      //   return
      // }
      let cache = JSON.parse(JSON.stringify(row))
      const keys = Object.keys(cache)
      keys.forEach(key => {
        if (Array.isArray(cache[key])) {
          cache[key] = cache[key].map(c => {
            return {
              guid: c
            }
          })
          cache[key] = JSON.stringify(cache[key])
        }
      })
      cache[titleObj.name] = ''
      this.refKeys.forEach(k => {
        delete cache[k + 'Options']
      })
      const filterValue = row[titleObj.name]
      const attr = titleObj.entity + '__' + titleObj.name
      const params = {
        filters: [
          {
            name: 'guid',
            operator: 'in',
            value: Array.isArray(filterValue) ? filterValue : [filterValue]
          }
        ],
        paging: false,
        dialect: {
          associatedData: {
            ...cache
          }
        }
      }
      const { statusCode, data } = await getRefOptions(this.requestId, attr, params)
      if (statusCode === 'OK') {
        row[titleObj.name + 'Options'] = data
        this.$set(this.tableData, index, row)
      }
    },
    handleDeleteRow (row) {
      this.$Modal.confirm({
        title: this.$t('confirm') + '删除',
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          this.tableData.splice(row._index, 1)
          this.requestData.forEach(item => {
            if (item.entity === this.activeTab || item.itemGroup === this.activeTab) {
              item.value.splice(row._index, 1)
            }
          })
        },
        onCancel: () => {}
      })
    },
    handleEditRow (row) {
      this.viewDisabled = false
      this.editVisible = true
      this.editData = deepClone(row)
    },
    submitEditDrawer () {
      this.tableData = this.tableData.map(item => {
        if (item._id === this.editData._id) {
          for (let key in item) {
            item[key] = this.editData[key]
          }
        }
        return item
      })
      this.refreshRequestData()
    },
    handleViewRow (row) {
      this.viewDisabled = true
      this.editVisible = true
      this.editData = deepClone(row)
    },
    // 保存草稿
    handleDraft () {},
    // 发布
    handlePublish () {}
  }
}
</script>

<style lang="scss" scoped>
.workbench-publish-create {
  .w-header {
    padding-bottom: 20px;
    margin-bottom: 20px;
    border-bottom: 2px dashed #d7dadc;
    .steps {
      display: flex;
      align-items: center;
      .title {
        font-size: 14px;
        font-weight: 500;
        margin-right: 20px;
      }
      .role {
        display: flex;
        flex-direction: column;
      }
    }
    .btn-group {
      text-align: right;
    }
  }
  .content {
    min-height: 500px;
    .split-line {
      border-right: 2px dashed #d7dadc;
      padding-right: 20px;
    }
    .count {
      font-weight: bold;
      font-size: 14px;
      margin-left: 10px;
    }
    .program {
      padding: 20px;
    }
  }
}
</style>
<style lang="scss">
.workbench-publish-create {
  .ivu-steps-content {
    padding-left: 0px;
    padding-top: 5px;
    font-size: 12px;
    color: #3d3c38 !important;
  }
  .ivu-steps .ivu-steps-tail > i {
    height: 3px;
    background: #8189a5;
  }
  .ivu-radio {
    display: none;
  }
  .ivu-radio-wrapper {
    border-radius: 20px;
    font-size: 12px;
    color: #000;
    background: #fff;
  }
  .ivu-radio-wrapper-checked.ivu-radio-border {
    border-color: #2d8cf0;
    background: #2d8cf0;
    color: #fff;
  }
  .ivu-form-item {
    margin-bottom: 25px;
  }
  .ivu-form-item-content {
    line-height: 20px;
  }
  .required {
    color: red;
  }
}
</style>
