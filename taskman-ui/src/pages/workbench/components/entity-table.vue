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
import { getRefOptions } from '@/api/server'
import { debounce, deepClone } from '@/pages/util'
import EditDrawer from './edit-entity-item.vue'
export default {
  components: {
    EditDrawer
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
    isAdd: {
      type: Boolean,
      default: false
    },
    formDisable: {
      type: Boolean,
      default: false
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
          width: 100,
          fixed: 'right',
          align: 'center',
          render: (h, params) => {
            return (
              <div style="display:flex;justify-content:space-around;">
                {!this.formDisable && (
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
                )}
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
    data: {
      handler (val) {
        if (val && val.length) {
          this.requestData = deepClone(val)
          this.requestData.forEach(item => {
            if (item.value.length === 0 && this.isAdd) {
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
              <div style="display:flex;align-items:center">
                {t.name === 'deploy_package' && t.required === 'yes' && !params.row[t.name] && (
                  <Icon size="24" color="#ed4014" type="md-apps" />
                )}
                {t.name === 'deploy_package' && t.required === 'yes' && params.row[t.name] && (
                  <Icon size="24" color="#19be6b" type="md-apps" />
                )}
                <Select
                  value={params.row[t.name]}
                  on-on-change={v => {
                    this.tableData[params.row._index][t.name] = v
                    this.refreshRequestData()
                  }}
                  multiple={t.multiple === 'Y'}
                  disabled={t.isEdit === 'no' || this.formDisable}
                >
                  {Array.isArray(params.row[t.name + 'Options']) &&
                    params.row[t.name + 'Options'].map(i => (
                      <Option value={t.entity ? i.guid : i} key={t.entity ? i.guid : i}>
                        {t.entity ? i.key_name : i}
                      </Option>
                    ))}
                </Select>
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
          v.entityData[rfk + 'Options'] = []
        })
        if (!v.entityData._id) {
          v.entityData._id = v.id
        }
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
