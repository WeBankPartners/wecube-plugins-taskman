<template>
  <div class="workbench-entity-table">
    <div class="radio-group">
      <div
        v-for="(item, index) in requestData"
        :key="index"
        @click="handleTabChange(item)"
        :class="{
          radio: true,
          custom: item.itemGroupType === 'custom',
          workflow: item.itemGroupType === 'workflow',
          optional: item.itemGroupType === 'optional',
          'fix-old-data': !item.itemGroupType
        }"
        :style="activeStyle(item)"
      >
        {{ `${item.itemGroup}` }}<span class="count">{{ item.value.length }}</span>
      </div>
    </div>
    <div class="form-table">
      <div v-for="(value, index) in tableData" :key="index" class="table-item">
        <div class="number">{{ index + 1 }}</div>
        <div class="form">
          <Form :model="value" ref="form" label-position="left" :label-width="100">
            <Row :key="index">
              <Col v-for="(i, index) in formOptions" :key="index" :span="i.width || 24">
                <FormItem
                  :label="i.title"
                  :prop="i.name"
                  :required="i.required === 'yes'"
                  :rules="
                    i.required === 'yes' ? [{ required: true, message: `${i.title}为空`, trigger: 'change' }] : []
                  "
                >
                  <!--输入框-->
                  <Input
                    v-if="i.elementType === 'input'"
                    v-model="value[i.name]"
                    :disabled="i.isEdit === 'no' || formDisable"
                    style="width: calc(100% - 20px)"
                  ></Input>
                  <Input
                    v-else-if="i.elementType === 'textarea'"
                    v-model="value[i.name]"
                    type="textarea"
                    :disabled="i.isEdit === 'no' || formDisable"
                    style="width: calc(100% - 20px)"
                  ></Input>
                  <LimitSelect
                    v-if="i.elementType === 'select' || i.elementType === 'wecmdbEntity'"
                    v-model="value[i.name]"
                    :displayName="i.elementType === 'wecmdbEntity' ? 'displayName' : 'key_name'"
                    :displayValue="i.elementType === 'wecmdbEntity' ? 'id' : 'guid'"
                    :objectOption="!!i.entity || i.elementType === 'wecmdbEntity'"
                    :options="value[i.name + 'Options']"
                    :disabled="i.isEdit === 'no' || formDisable"
                    :multiple="i.multiple === 'Y'"
                    style="width: calc(100% - 20px)"
                  >
                  </LimitSelect>
                  <Input
                    v-else-if="i.elementType === 'calculate'"
                    :value="i.routineExpression"
                    type="textarea"
                    :disabled="true"
                    style="width: calc(100% - 20px)"
                  ></Input>
                  <DatePicker
                    v-else-if="i.elementType === 'datePicker'"
                    v-model="value[i.name]"
                    format="yyyy-MM-dd"
                    :disabled="i.isEdit === 'no' || formDisable"
                    type="date"
                    style="width: calc(100% - 20px)"
                  >
                  </DatePicker>
                </FormItem>
              </Col>
            </Row>
          </Form>
        </div>
        <div v-if="!formDisable && tableData.length > 1" class="button">
          <Icon type="md-trash" color="#ed4014" size="24" @click="handleDeleteRow" />
        </div>
      </div>
    </div>
    <div class="add-row">
      <Button v-if="isAdd && type === '2' && activeItem.itemGroupRule === 'new'" type="primary" @click="addRow">{{
        $t('tw_add_row')
      }}</Button>
      <!--选择已有数据添加一行-->
      <Select
        v-if="isAdd && type === '2' && activeItem.itemGroupRule === 'exist'"
        v-model="addRowSource"
        filterable
        clearable
        placeholder="选择已有数据添加一行"
        style="width:450px;"
        prefix="md-add-circle"
        @on-open-change="getCmdbEntityList"
        @on-change="addRow"
      >
        <template #prefix>
          <Icon type="md-add-circle" color="#2d8cf0" :size="24"></Icon>
        </template>
        <template v-for="(i, index) in addRowSourceOptions">
          <Option :key="index" :value="i.id">{{ i.displayName }}</Option>
        </template>
      </Select>
    </div>
  </div>
</template>

<script>
import { getRefOptions, getWeCmdbOptions } from '@/api/server'
import { debounce, deepClone } from '@/pages/util'
import EntityItem from './edit-entity-item.vue'
import LimitSelect from '@/pages/components/limit-select.vue'
export default {
  components: {
    EntityItem,
    LimitSelect
  },
  props: {
    data: {
      type: Array,
      default: () => []
    },
    requestId: {
      type: String,
      default: ''
    },
    // 是否创建页面
    isAdd: {
      type: Boolean,
      default: false
    },
    // 无数据时，是否默认添加一行
    autoAddRow: {
      type: Boolean,
      default: false
    },
    formDisable: {
      type: Boolean,
      default: false
    },
    // 类型(1发布2请求)
    type: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      requestData: [],
      activeTab: '',
      activeItem: '',
      refKeys: [], // 引用类型字段集合select类型
      formOptions: [],
      tableData: [],
      addRowSource: '',
      addRowSourceOptions: []
    }
  },
  computed: {
    activeStyle () {
      return function (item) {
        let color = '#fff'
        if (this.activeTab === item.entity || this.activeTab === item.itemGroup) {
          if (item.itemGroupType === 'workflow') {
            color = '#ebdcb4'
          } else if (item.itemGroupType === 'custom') {
            color = 'rgba(184, 134, 248, 0.6)'
          } else if (item.itemGroupType === 'optional') {
            color = 'rgba(129, 179, 55, 0.6)'
          } else if (!item.itemGroupType) {
            color = 'rgb(45, 140, 240)'
          }
        }
        return { background: color }
      }
    }
  },
  watch: {
    data: {
      handler (val) {
        if (val && val.length) {
          this.requestData = deepClone(val)
          this.requestData.forEach(item => {
            // 新增时，没有配置数据，默认添加一行
            if (item.value.length === 0 && this.autoAddRow) {
              this.handleAddRow(item)
            }
          })
          this.activeTab = this.requestData[0].entity || this.requestData[0].itemGroup
          this.activeItem = this.requestData[0]
          this.initTableData()
        }
      },
      deep: true,
      immediate: true
    }
  },
  methods: {
    // 提交时，定位到没有填写必填项的页签
    validTable (index) {
      if (index !== '') {
        if (this.activeTab === (this.requestData[index].entity || this.requestData[index].itemGroup)) {
          return
        }
        this.activeTab = this.requestData[index].entity || this.requestData[index].itemGroup
        this.initTableData()
      }
    },
    // 切换tab刷新表格数据，加上防抖避免切换过快显示异常问题
    handleTabChange: debounce(function (item) {
      this.activeTab = item.entity || item.itemGroup
      this.activeItem = item
      this.initTableData()
      this.addRowSource = ''
      this.addRowSourceOptions = []
    }, 100),
    async initTableData () {
      // 当前选择tab数据
      const data = this.requestData.find(r => r.entity === this.activeTab || r.itemGroup === this.activeTab)
      this.oriData = data
      // select类型集合
      this.refKeys = []
      data.title.forEach(t => {
        if (t.elementType === 'select' || t.elementType === 'wecmdbEntity') {
          this.refKeys.push(t.name)
        }
      })
      this.formOptions = data.title
      // table数据初始化
      this.tableData = data.value.map(v => {
        this.refKeys.forEach(rfk => {
          // 缓存RefOptions数据，不需要每次调用
          if (!(v.entityData[rfk + 'Options'] && v.entityData[rfk + 'Options'].length > 0)) {
            v.entityData[rfk + 'Options'] = []
          }
        })
        if (!v.entityData._id) {
          v.entityData._id = v.id
        }
        return v.entityData
      })
      // 下拉类型数据初始化(待优化，调用接口太多)
      this.tableData.forEach((row, index) => {
        this.refKeys.forEach(rfk => {
          if (!(row[rfk + 'Options'] && row[rfk + 'Options'].length > 0)) {
            const titleObj = data.title.find(f => f.name === rfk)
            this.getRefOptions(titleObj, row, index)
          }
        })
      })
    },
    async getRefOptions (titleObj, row, index) {
      // taskman模板管理配置的普通下拉类型(值用逗号拼接)
      if (titleObj.elementType === 'select' && titleObj.entity === '') {
        row[titleObj.name + 'Options'] = (titleObj.dataOptions && titleObj.dataOptions.split(',')) || []
        this.$set(this.tableData, index, row)
        return
      }
      // taskman模板管理配置的引用下拉类型
      if (titleObj.elementType === 'wecmdbEntity') {
        const [packageName, ciType] = (titleObj.dataOptions && titleObj.dataOptions.split(':')) || []
        const { status, data } = await getWeCmdbOptions(packageName, ciType, {})
        if (status === 'OK') {
          row[titleObj.name + 'Options'] = data
          this.$set(this.tableData, index, row)
        }
        return
      }
      // if (titleObj.refEntity === '') {
      //   row[titleObj.name + 'Options'] = titleObj.selectList
      //   this.$set(this.tableData, index, row)
      //   return
      // }
      let cache = JSON.parse(JSON.stringify(row))
      cache[titleObj.name] = ''
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
        // 删除掉值为空的数据
        if (!cache[key] || (Array.isArray(cache[key]) && cache[key].length === 0)) {
          delete cache[key]
        }
      })
      this.refKeys.forEach(k => {
        delete cache[k + 'Options']
      })
      delete cache._checked
      delete cache._disabled
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
        title: this.$t('confirm') + this.$t('delete'),
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
    // 手动添加一行数据
    addRow () {
      if (this.activeItem.itemGroupRule === 'new') {
        const data = this.requestData.find(r => r.entity === this.activeTab || r.itemGroup === this.activeTab)
        this.handleAddRow(data)
        this.initTableData()
      } else if (this.activeItem.itemGroupRule === 'exist') {
        if (this.addRowSource) {
          const source = this.addRowSourceOptions.find(i => i.id === this.addRowSource)
          this.handleAddRow(this.activeItem, source)
          this.initTableData()
        }
      }
    },
    // 添加一条行数据
    handleAddRow (data, source) {
      let entityData = {}
      data.title.forEach(item => {
        // 选择已有数据添加一行，填充默认值
        if (source) {
          for (let key of Object.keys(source)) {
            if (key === item.name) {
              entityData[item.name] = source[key]
            }
          }
        } else {
          entityData[item.name] = ''
        }

        if (item.elementType === 'select' || item.elementType === 'wecmdbEntity') {
          entityData[item.name + 'Options'] = []
        }
      })
      const idStr = new Date().getTime().toString()
      let obj = {
        dataId: '',
        displayName: '',
        entityData: { ...entityData, _id: '' },
        entityName: data.entity,
        entityDataOp: 'create',
        fullDataId: '',
        id: idStr,
        packageName: data.packageName,
        previousIds: [],
        succeedingIds: [],
        addFlag: true // 前端添加一行标识，提交时需删除
      }
      data.value.push(obj)
    },
    async getCmdbEntityList () {
      const { packageName, entity } = this.activeItem
      const { status, data } = await getWeCmdbOptions(packageName, entity, {})
      if (status === 'OK') {
        this.addRowSourceOptions = data || []
        // 过滤下拉框数据(表单dataId和下拉框数据id相同)
        if (this.activeItem.value) {
          let dataIds = []
          this.activeItem.value.forEach(i => i.dataId && dataIds.push(i.dataId))
          this.addRowSourceOptions = this.addRowSourceOptions.filter(i => !dataIds.includes(i.id))
        }
      }
    }
  }
}
</script>

<style lang="scss">
.workbench-entity-table {
  width: 100%;
  .radio-group {
    display: flex;
    margin-bottom: 15px;
    .radio {
      padding: 5px 15px;
      border-radius: 32px;
      font-size: 14px;
      cursor: pointer;
      margin-right: 10px;
    }
    .custom {
      border: 1px solid #b886f8;
      color: #b886f8;
    }
    .workflow {
      border: 1px solid #cba43f;
      color: #cba43f;
    }
    .optional {
      border: 1px solid #81b337;
      color: #81b337;
    }
    .fix-old-data {
      border: 1px solid #dcdee2;
      color: #000;
    }
  }
  .count {
    font-weight: bold;
    font-size: 14px;
    margin-left: 10px;
  }
  .add-row {
    margin-top: 10px;
  }
  .form-table {
    position: relative;
    .table-item {
      display: flex;
      border-left: 1px dashed #dcdee2;
      border-right: 1px dashed #dcdee2;
      border-bottom: 1px dashed #dcdee2;
      border-radius: 4px;
      &:first-child {
        border-top: 1px dashed #dcdee2;
      }
      .number {
        width: 50px;
        display: flex;
        justify-content: center;
        align-items: center;
        border-right: 1px dashed #dcdee2;
      }
      .form {
        padding: 20px 0 10px 20px;
        flex: 1;
      }
      .button {
        cursor: pointer;
        width: 50px;
        display: flex;
        justify-content: center;
        align-items: center;
        border-left: 1px solid #dcdee2;
      }
    }
  }
  .ivu-form-item {
    margin-bottom: 10px !important;
  }
  .ivu-form-item-label {
    word-wrap: break-word;
    padding: 10px 10px 10px 0;
  }
}
</style>
