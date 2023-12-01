<template>
  <div class="workbench-publish-create">
    <Row class="w-header">
      <Col span="18" class="steps">
        <span class="title">请求进度：</span>
        <Steps :current="0" style="max-width:600px;">
          <Step v-for="(i, index) in progressList" :key="index" :content="i.name">
            <template #icon>
              <Icon size="24" :type="i.icon" :color="i.color" />
            </template>
            <div class="role" slot="content">
              <div class="word-eclipse">{{ i.name }}</div>
              <span>{{ i.handler }}</span>
            </div>
          </Step>
        </Steps>
      </Col>
      <Col span="6" class="btn-group">
        <Button v-if="this.type !== 'detail'" @click="handleDraft" style="margin-right:10px;">保存草稿</Button>
        <Button v-if="this.type !== 'detail'" type="primary" @click="handlePublish">发布</Button>
      </Col>
    </Row>
    <Row class="content">
      <Col span="16" class="split-line">
        <Form :model="form" label-position="right" :label-width="120">
          <HeaderTitle v-if="type === 'detail'" title="请求信息">
            <Row>
              <Col :span="4">请求ID：</Col>
              <Col :span="8">{{ detailInfo.id }}</Col>
              <Col :span="4">请求类型：</Col>
              <Col :span="8">{{ detailInfo.requestType }}</Col>
            </Row>
            <Row style="margin-top:10px;">
              <Col :span="4">创建时间：</Col>
              <Col :span="8">{{ detailInfo.createdTime }}</Col>
              <Col :span="4">期望完成时间：</Col>
              <Col :span="8">{{ detailInfo.expectTime }}</Col>
            </Row>
            <Row style="margin-top:10px;">
              <Col :span="4">请求进度：</Col>
              <Col :span="8">{{ detailInfo.progress }}</Col>
              <Col :span="4">请求状态：</Col>
              <Col :span="8">{{ detailInfo.status }}</Col>
            </Row>
            <Row style="margin-top:10px;">
              <Col :span="4">当前节点：</Col>
              <Col :span="8">{{ detailInfo.curNode }}</Col>
              <Col :span="4">当前处理人：</Col>
              <Col :span="8">{{ detailInfo.handler }}</Col>
            </Row>
            <Row style="margin-top:10px;">
              <Col :span="4">创建人：</Col>
              <Col :span="8">{{ detailInfo.createdBy }}</Col>
              <Col :span="4">创建人角色：</Col>
              <Col :span="8">{{ detailInfo.role }}</Col>
            </Row>
            <Row style="margin-top:10px;">
              <Col :span="4">使用模板：</Col>
              <Col :span="8">{{ detailInfo.templateName }}</Col>
              <Col :span="4">模板组：</Col>
              <Col :span="8">{{ detailInfo.templateGroupName }}</Col>
            </Row>
            <Row style="margin-top:10px;">
              <Col :span="4">请求描述：</Col>
              <Col :span="8">{{ detailInfo.description }}</Col>
            </Row>
          </HeaderTitle>
          <HeaderTitle v-else title="发布信息">
            <FormItem label="请求名称" required>
              <Input v-model="form.name" placeholder="请输入" style="width:400px;" />
            </FormItem>
            <FormItem label="发布描述">
              <Input v-model="form.description" placeholder="请输入" style="width:400px;" />
            </FormItem>
            <FormItem :label="$t('expected_completion_time')">
              <DatePicker
                type="datetime"
                :value="form.expectTime"
                @on-change="
                  val => {
                    form.expectTime = val
                  }
                "
                placeholder="Select date"
                :options="{
                  disabledDate(date) {
                    return date && date.valueOf() < Date.now() - 86400000
                  }
                }"
                style="width:400px;"
              ></DatePicker>
            </FormItem>
          </HeaderTitle>
          <HeaderTitle title="发布目标对象">
            <FormItem label="选择操作单元" required>
              <Select
                v-model="form.rootEntityId"
                :disabled="type === 'detail'"
                clearable
                filterable
                style="width:300px;"
              >
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
        <!-- <StaticFlow :templateId="templateId" ref="staticFlowRef"></StaticFlow> -->
        <DynamicFlow ref="staticFlowRef"></DynamicFlow>
      </Col>
    </Row>
    <EditDrawer
      v-if="editVisible"
      v-model="editData"
      :options="editOptions"
      :visible.sync="editVisible"
      :disabled="viewDisabled"
      @submit="submitEditRow"
    ></EditDrawer>
  </div>
</template>

<script>
import HeaderTitle from '../components/header-title.vue'
import EditDrawer from './edit-item.vue'
import {
  getCreateInfo,
  getProgressInfo,
  getRootEntity,
  getEntityData,
  getRefOptions,
  savePublishData,
  getPublishInfo,
  updateRequestStatus
} from '@/api/server'
import StaticFlow from './flow/static-flow.vue'
import DynamicFlow from './flow/dynamic-flow.vue'
import { deepClone, debounce } from '@/pages/util'
const statusIcon = {
  1: 'md-pin',
  2: 'md-radio-button-on',
  3: 'ios-checkmark-circle-outline',
  4: 'ios-close-circle-outline'
}
const statusColor = {
  1: '#ffa500',
  2: '#8189a5',
  3: '#19be6b',
  4: '#ed4014'
}
export default {
  components: {
    HeaderTitle,
    StaticFlow,
    DynamicFlow,
    EditDrawer
  },
  data () {
    return {
      type: this.$route.query.type, // detail查看
      templateId: '',
      requestId: '',
      activeTab: '',
      refKeys: [], // 引用类型字段集合select类型
      detailInfo: {}, // 详情信息
      form: {
        name: '',
        description: '',
        expectTime: '',
        rootEntityId: '', // 目标对象
        data: []
      },
      rootEntityOptions: [],
      progressList: [],
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
  watch: {
    'form.rootEntityId' (val) {
      if (val) {
        this.getEntityData()
      }
    }
  },
  async mounted () {
    this.templateId = this.$route.query.templateId || ''
    this.$refs.staticFlowRef.orchestrationSelectHandler(this.templateId)
    this.requestId = this.$route.query.requestId || ''
    // 查看调用详情接口
    if (this.requestId) {
      await this.getPublishInfo()
    } else {
      await this.getCreateInfo()
    }
    this.getProgressInfo()
    this.getEntity()
    // this.getEntityData()
  },
  methods: {
    // 切换tab刷新表格数据，加上防抖避免切换过快显示异常问题
    handleTabChange: debounce(function () {
      this.initTableData()
    }, 100),
    // 创建发布,使用模板ID获取详情数据
    async getCreateInfo () {
      const params = {
        requestTemplate: this.templateId,
        role: this.$route.query.role
      }
      const { statusCode, data } = await getCreateInfo(params)
      if (statusCode === 'OK') {
        this.requestId = data.id
        this.form.name = data.name
      }
    },
    // 获取请求进度
    async getProgressInfo () {
      const params = {
        templateId: this.templateId,
        requestId: ''
      }
      const { statusCode, data } = await getProgressInfo(params)
      if (statusCode === 'OK') {
        this.progressList = data
        this.progressList.forEach(item => {
          item.icon = statusIcon[item.status]
          item.color = statusColor[item.status]
          switch (item.node) {
            case 'sendRequest':
              item.name = '提起请求'
              break
            case 'requestPending':
              item.name = '请求定版'
              break
            case 'requestComplete':
              item.name = '请求完成'
              break
            default:
              item.name = item.node
              break
          }
        })
      }
    },
    // 获取发布信息
    async getPublishInfo () {
      const { statusCode, data } = await getPublishInfo(this.requestId)
      if (statusCode === 'OK') {
        this.detailInfo = data.request || {}
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
        this.rootEntityOptions = data.data || []
        this.form.rootEntityId = this.rootEntityOptions[0] && this.rootEntityOptions[0].guid
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
        // 没有数据，默认添加一行
        this.requestData.forEach(item => {
          if (item.value.length === 0) {
            this.handleAddRow(item)
          }
        })
        this.activeTab = this.activeTab || data.data[0].entity || data.data[0].itemGroup
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
          align: 'left',
          minWidth: 200
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
                disabled={t.isEdit === 'no' || this.type === 'detail'}
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
                disabled={t.isEdit === 'no' || this.type === 'detail'}
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
                disabled={t.isEdit === 'no' || this.type === 'detail'}
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
    // 删除行数据
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
    // 添加一条行数据
    handleAddRow (data) {
      let entityData = {}
      data.title.forEach(item => {
        entityData[item.name] = ''
        if (item.elementType === 'select') {
          entityData[item.name + 'Options'] = []
        }
      })
      let obj = {
        dataId: '',
        displayName: '',
        entityData: entityData,
        entityName: data.entity,
        entityDataOp: 'create',
        fullDataId: '',
        id: '',
        packageName: data.packageName,
        previousIds: [],
        succeedingIds: []
      }
      data.value.push(obj)
    },
    handleEditRow (row) {
      this.viewDisabled = false
      this.editVisible = true
      this.editData = deepClone(row)
    },
    submitEditRow () {
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
    async handleDraft () {
      if (this.form.rootEntityId === '') {
        this.$Message.warning(this.$t('root_entity') + this.$t('can_not_be_empty'))
        return
      }
      this.form.data = this.requestData
      const item = this.rootEntityOptions.find(item => item.guid === this.form.rootEntityId)
      this.form.entityName = item.key_name
      if (this.requestDataCheck()) {
        const { statusCode } = await savePublishData(this.requestId, this.form)
        if (statusCode === 'OK') {
          this.$Notice.success({
            title: this.$t('successful'),
            desc: this.$t('successful')
          })
        }
        return statusCode
      } else {
        this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('required_tip')
        })
      }
    },
    // 发布
    async handlePublish () {
      this.$Modal.confirm({
        title: this.$t('confirm') + this.$t('commit'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const draftResult = await this.handleDraft()
          console.log('111111111111', draftResult)
          if (draftResult === 'OK') {
            const { statusCode } = await updateRequestStatus(this.requestId, 'Pending')
            if (statusCode === 'OK') {
              this.$router.push({ path: '/taskman/workbench' })
            }
          }
        },
        onCancel: () => {}
      })
    },
    requestDataCheck () {
      let result = true
      this.requestData.forEach(requestData => {
        let requiredName = []
        requestData.title.forEach(t => {
          if (t.required === 'yes') {
            requiredName.push(t.name)
          }
        })
        requestData.value.forEach(v => {
          requiredName.forEach(key => {
            let val = v.entityData[key]
            if (Array.isArray(val)) {
              if (val.length === 0) {
                result = false
              }
            } else {
              if (val === '') {
                result = false
              }
            }
          })
        })
      })
      return result
    }
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
      .word-eclipse {
        max-width: 100px;
        text-overflow: ellipsis;
        overflow: hidden;
        white-space: nowrap;
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
    padding-left: 0px !important;
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
