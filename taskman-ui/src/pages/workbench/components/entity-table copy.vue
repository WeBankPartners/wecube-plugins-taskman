<template>
  <div class="workbench-entity-table">
    <RadioGroup v-model="activeTab" @on-change="handleTabChange" style="margin-bottom:20px;">
      <Radio v-for="(item, index) in requestData" :label="item.entity || item.itemGroup" :key="index" border>
        <span
          >{{ `${item.itemGroup}` }}<span class="count">{{ item.value.length }}</span></span
        >
      </Radio>
    </RadioGroup>
    <Table size="small" :columns="tableColumns" :data="tableData" @on-selection-change="handleChooseData"></Table>
    <Button v-if="isAdd && type === '2'" size="small" style="margin-top: 10px;" @click="addRow">{{
      $t('tw_add_row')
    }}</Button>
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
import { getRefOptions, getWeCmdbOptions } from '@/api/server'
import { debounce, deepClone } from '@/pages/util'
import EditDrawer from './edit-entity-item.vue'
import limitSelect from '@/pages/components/limit-select.vue'
export default {
  components: {
    EditDrawer,
    limitSelect
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
      refKeys: [], // 引用类型字段集合select类型
      tableColumns: [],
      tableData: [],
      editVisible: false,
      editOptions: [],
      editData: {},
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
          width: 100,
          fixed: 'right',
          align: 'center',
          render: (h, params) => {
            return (
              <div style="display:flex;justify-content:space-around;">
                {!this.formDisable && (
                  <Tooltip content={this.$t('edit')} placement="top-start">
                    <Icon
                      size="20"
                      type="md-create"
                      style="cursor:pointer;"
                      onClick={() => {
                        this.handleEditRow(params.row)
                      }}
                    />
                  </Tooltip>
                )}
                <Tooltip content={this.$t('detail')} placement="top-start">
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
    data: {
      handler (val) {
        if (val && val.length) {
          this.requestData = deepClone(val)
          this.requestData.forEach(item => {
            // 新增时，没有配置数据，默认添加一行
            if (item.value.length === 0 && this.autoAddRow) {
              this.handleAddRow(item)
            }
            item.value.forEach(j => {
              if (this.formDisable) {
                j.entityData._disabled = true
              }
              j.entityData._checked = true
            })
          })
          this.activeTab = this.activeTab || this.requestData[0].entity || this.requestData[0].itemGroup
          this.initTableData()
        }
      },
      deep: true,
      immediate: true
    }
  },
  mounted () {},
  methods: {
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
    handleTabChange: debounce(function () {
      this.initTableData()
    }, 100),
    // 编辑操作，刷新requestData
    refreshRequestData () {
      this.requestData.forEach(item => {
        if (item.entity === this.activeTab || item.itemGroup === this.activeTab) {
          for (let m of item.value) {
            for (let n of this.tableData) {
              if (m.id === n._id) {
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
      // select类型集合
      this.refKeys = []
      data.title.forEach(t => {
        if (t.elementType === 'select' || t.elementType === 'wecmdbEntity') {
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
          minWidth: 240
        }
        if (t.required === 'yes') {
          column.renderHeader = (h, { column }) => {
            return (
              <span>
                {`${column.title}`}
                <span class="required">{`（${this.$t('required')}）`}</span>
              </span>
            )
          }
        }
        if (t.elementType === 'select' || t.elementType === 'wecmdbEntity') {
          column.render = (h, params) => {
            return (
              <div style="display:flex;align-items:center">
                {t.name === 'deploy_package' && t.required === 'yes' && !params.row[t.name] && (
                  <Icon size="24" color="#ed4014" type="md-apps" />
                )}
                {t.name === 'deploy_package' && t.required === 'yes' && params.row[t.name] && (
                  <Icon size="24" color="#19be6b" type="md-apps" />
                )}
                {
                  <limitSelect
                    value={params.row[t.name]}
                    on-on-change={v => {
                      this.tableData[params.row._index][t.name] = v
                      this.refreshRequestData()
                    }}
                    displayName={t.elementType === 'wecmdbEntity' ? 'displayName' : 'key_name'}
                    displayValue={t.elementType === 'wecmdbEntity' ? 'id' : 'guid'}
                    objectOption={!!t.entity || t.elementType === 'wecmdbEntity'}
                    options={params.row[t.name + 'Options']}
                    disabled={t.isEdit === 'no' || this.formDisable}
                    multiple={t.multiple === 'Y'}
                  />
                }
                {
                  // <Select
                  //   value={params.row[t.name]}
                  //   on-on-change={v => {
                  //     this.tableData[params.row._index][t.name] = v
                  //     this.refreshRequestData()
                  //   }}
                  //   filterable
                  //   clearable
                  //   multiple={t.multiple === 'Y'}
                  //   disabled={t.isEdit === 'no' || this.formDisable}
                  // >
                  //   {Array.isArray(params.row[t.name + 'Options']) &&
                  //     params.row[t.name + 'Options'].map(i => (
                  //       <Option value={t.entity ? i.guid : i} key={t.entity ? i.guid : i}>
                  //         {t.entity ? i.key_name : i}
                  //       </Option>
                  //     ))}
                  // </Select>
                }
              </div>
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
                disabled={t.isEdit === 'no' || this.formDisable}
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
                disabled={t.isEdit === 'no' || this.formDisable}
              />
            )
          }
        }
        this.tableColumns.push(column)
      })

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
      const data = this.requestData.find(r => r.entity === this.activeTab || r.itemGroup === this.activeTab)
      this.handleAddRow(data)
      this.initTableData()
      // const addRow = data.value[data.value.length - 1].entityData
      // this.tableData.push(addRow)
      // this.refKeys.forEach(rfk => {
      //   const titleObj = data.title.find(f => f.name === rfk)
      //   this.getRefOptions(titleObj, addRow, data.value.length - 1)
      // })
    },
    // 添加一条行数据
    handleAddRow (data) {
      let entityData = {}
      data.title.forEach(item => {
        entityData[item.name] = ''
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
    handleEditRow (row) {
      this.viewDisabled = false
      this.editVisible = true
      // 当前选择tab数据
      const data = this.requestData.find(r => r.entity === this.activeTab || r.itemGroup === this.activeTab)
      // 编辑表单的options配置
      this.editOptions = data.title
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
      // 当前选择tab数据
      const data = this.requestData.find(r => r.entity === this.activeTab || r.itemGroup === this.activeTab)
      // 编辑表单的options配置
      this.editOptions = data.title
      this.editData = deepClone(row)
    },
    // 给勾选的表格数据设置_checked属性
    handleChooseData (selection) {
      const selectIds = selection.map(i => i._id)
      this.tableData.forEach(row => {
        if (selectIds.includes(row._id)) {
          row._checked = true
        } else {
          row._checked = false
        }
      })
    }
  }
}
</script>

<style lang="scss">
.workbench-entity-table {
  width: 100%;
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
  .count {
    font-weight: bold;
    font-size: 14px;
    margin-left: 10px;
  }
}
</style>
