<template>
  <div class="workbench-entity-table">
    <div class="workbench-entity-table-radio-group">
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
        {{ `${item.itemGroupName}` }}<span class="count">{{ item.value.length }}</span>
      </div>
    </div>
    <div class="form-table">
      <div v-for="(value, index) in tableData" :key="index" class="table-item">
        <div class="number">{{ index + 1 }}</div>
        <div class="form">
          <Form :model="value" ref="form" label-position="left" :label-width="100">
            <Row type="flex" justify="start" :key="index">
              <template v-for="i in formOptions">
                <Col v-if="!value[i.name + 'Hidden']" :key="i.id" :span="i.width || 24">
                  <FormItem
                    :key="i.id"
                    :label="i.title"
                    :prop="i.name"
                    :required="i.required === 'yes'"
                    :rules="
                      i.required === 'yes'
                        ? [
                            {
                              required: true,
                              message: `${i.title}${$t('can_not_be_empty')}`,
                              trigger: ['change', 'blur']
                            }
                          ]
                        : []
                    "
                  >
                    <template v-if="['diffVariable', 'password', 'object'].includes(i.inputType)">
                      <div style="display:flex;align-items:center;margin-top:5px;">
                        <Icon
                          size="18"
                          type="ios-apps-outline"
                          color="#2d8cf0"
                          style="font-weight:bold;cursor:pointer;"
                          @click="handleOpenCmdbDetail(i, value)"
                        />
                        {{ value[i.name] || '1111111111111111' }}
                      </div>
                    </template>
                    <template v-else>
                      <!--输入框-->
                      <Input
                        v-if="i.elementType === 'input'"
                        v-model.trim="value[i.name]"
                        :disabled="formDisabled(i)"
                        style="width: calc(100% - 20px)"
                      ></Input>
                      <Input
                        v-else-if="i.elementType === 'textarea'"
                        v-model.trim="value[i.name]"
                        type="textarea"
                        :disabled="formDisabled(i)"
                        style="width: calc(100% - 20px)"
                      ></Input>
                      <LimitSelect
                        v-if="i.elementType === 'select' || i.elementType === 'wecmdbEntity'"
                        v-model="value[i.name]"
                        :displayName="i.elementType === 'wecmdbEntity' ? 'displayName' : i.entity ? 'key_name' : 'label'"
                        :displayValue="i.elementType === 'wecmdbEntity' ? 'id' : i.entity ? 'guid' : 'value'"
                        :options="value[i.name + 'Options']"
                        :disabled="formDisabled(i)"
                        :multiple="i.multiple === 'Y' || i.multiple === 'yes'"
                        style="width: calc(100% - 20px)"
                        @open-change="handleRefOpenChange(i, value, index)"
                      >
                      </LimitSelect>
                      <!--自定义分析类型-->
                      <Input
                        v-else-if="i.elementType === 'calculate'"
                        :value="value[i.name]"
                        type="textarea"
                        :disabled="true"
                        style="width: calc(100% - 20px)"
                      ></Input>
                      <!--日期时间类型-->
                      <DatePicker
                        v-else-if="i.elementType === 'datePicker'"
                        :value="value[i.name]"
                        @on-change="$event => handleTimeChange($event, value, i.name)"
                        format="yyyy-MM-dd HH:mm:ss"
                        :disabled="formDisabled(i)"
                        type="datetime"
                        style="width: calc(100% - 20px)"
                      >
                      </DatePicker>
                    </template>
                  </FormItem>
                </Col>
              </template>
            </Row>
          </Form>
        </div>
        <div v-if="!formDisable && tableData.length > 1 && isAdd" class="button">
          <Icon type="md-trash" color="#ed4014" size="24" @click="handleDeleteRow(index)" />
        </div>
      </div>
    </div>
    <div v-if="isAdd" class="add-row">
      <!--添加一行-->
      <Button v-if="activeItem.itemGroupRule === 'new'" type="primary" @click="addRow">{{ $t('tw_add_row') }}</Button>
      <!--选择已有数据添加一行-->
      <Select
        ref="addRowSelect"
        v-else-if="['exist', 'exist_empty'].includes(activeItem.itemGroupRule)"
        v-model="addRowSource"
        filterable
        clearable
        :placeholder="$t('tw_addRow_exist')"
        style="width:450px;"
        prefix="md-add-circle"
        @on-open-change="
          flag => {
            if (flag) getCmdbEntityList()
          }
        "
        @on-change="addRow"
      >
        <template #prefix>
          <Icon type="md-add-circle" color="#2d8cf0" :size="24"></Icon>
        </template>
        <template v-for="i in addRowSourceOptions">
          <Option :key="i.id" :value="i.id">{{ i.displayName }}</Option>
        </template>
      </Select>
    </div>
    <!--cmdb数据预览-->
    <BaseDrawer
      :title="cmdbInfo.title"
      :visible.sync="cmdbInfo.visible"
      realWidth="800"
      :scrollable="true"
    >
      <template slot="content">
        <div v-if="cmdbInfo.attr.inputType === 'diffVariable'">
          <div style="text-align: justify;word-break: break-word;overflow-y:auto;max-height:500px">
            <div style="text-align: left;">
              <Alert type="warning">如出现页面值未显示，请点击刷新按钮</Alert>
            </div>
            <div>{{ '111111111111111111111111111111' }}</div>
          </div>
        </div>
        <div v-else-if="cmdbInfo.attr.inputType === 'object'">
          <json-viewer :value="{ a: 1, b: { c: 2 }}" :expand-depth="5"></json-viewer>
        </div>
        <div v-else-if="cmdbInfo.attr.inputType === 'password'">
          <div style="text-align: justify;word-break: break-word;">
            {{ '111111111111111111111111111111' }}
          </div>
        </div>
      </template>
      <template slot="footer">
        <Button type="primary" @click="handleRefreshCmdbData">{{ '刷新' }}</Button>
        <Button type="default" @click="cmdbInfo.visible = false" style="margin-left:10px;">{{ '关闭' }}</Button>
      </template>
    </BaseDrawer>
  </div>
</template>

<script>
import LimitSelect from '@/pages/components/limit-select.vue'
import { getRefOptions, getWeCmdbOptions, saveFormData, getExpressionData } from '@/api/server'
import { debounce, deepClone, fixArrStrToJsonArray } from '@/pages/util'
import { evaluateCondition } from '../evaluate'
import JsonViewer from 'vue-json-viewer'
export default {
  components: {
    LimitSelect,
    JsonViewer
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
    }
  },
  data () {
    return {
      requestData: [],
      activeTab: '',
      activeItem: {}, // 当前选中数据
      refKeys: [], // 引用类型字段集合select类型
      calculateKeys: [], // 自定义计算分析类型集合
      formOptions: [],
      tableData: [],
      addRowSource: '',
      addRowSourceOptions: [],
      worklfowDataIdsObj: {}, // 编排类表单默认下发数据dataId集合
      cmdbInfo: {
        title: '',
        visible: false,
        attr: {},
        value: ''
      }
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
    },
    formDisabled () {
      return function (attr) {
        return attr.isEdit === 'no' || (attr.autofillable === 'yes' && attr.autoFillType === 'forced')
      }
    }
  },
  watch: {
    data: {
      handler (val) {
        if (val && val.length) {
          this.requestData = deepClone(val)
          this.requestData.forEach(item => {
            if (!item.value) {
              item.value = []
            }
            // 无表单数据，不是选择已有数据添加一行，默认添加一行
            if (item.value.length === 0 && this.autoAddRow && !['exist', 'exist_empty'].includes(item.itemGroupRule)) {
              this.handleAddRow(item)
            }
            // 备份编排类表单系统下发数据id集合
            if (['exist', 'exist_empty'].includes(item.itemGroupRule) && item.itemGroupType === 'workflow') {
              const list = item.value || []
              const ids = list.map(item => {
                return item.dataId
              })
              this.$set(this.worklfowDataIdsObj, item.formTemplateId, ids)
            }
            // 选择已有数据添加一行设置为【默认不选】，第一次加载数据表单清空value，审批和任务表单不需要
            if (item.itemGroupRule === 'exist_empty' && item.value && item.value.length && this.isAdd) {
              const firstFlag = item.value.some(i => i.entityData && !i.entityData.hasOwnProperty('_id'))
              if (firstFlag) {
                item.value = []
              }
            }
          })
          this.activeTab = this.requestData[0].entity || this.requestData[0].itemGroup
          this.activeItem = this.requestData[0]
          this.initTableData()
        } else {
          this.requestData = []
          this.activeTab = ''
          this.activeItem = {}
          this.formOptions = []
          this.tableData = []
        }
      },
      deep: true,
      immediate: true
    },
    tableData: {
      handler (val) {
        if (val) {
          val.forEach(item => {
            // 表单隐藏逻辑
            Object.keys(item).forEach(key => {
              const find = this.formOptions.find(i => i.name === key) || {}
              if (find.hiddenCondition && find.required === 'no') {
                const conditions = find.hiddenCondition || []
                item[key + 'Hidden'] = conditions.every(j => {
                  return evaluateCondition(j, item[j.name])
                })
              }
            })
          })
        }
      },
      immediate: true,
      deep: true
    }
  },
  methods: {
    // ref类型下拉框每次展开调用接口
    handleRefOpenChange (titleObj, row, index) {
      this.getRefOptions(titleObj, row, index, false)
    },
    // 保存当前表单组的数据
    async saveCurrentTabData (item) {
      await saveFormData(this.requestId, item)
    },
    // 切换tab刷新表格数据，加上防抖避免切换过快显示异常问题
    handleTabChange: debounce(function (item) {
      // 切换表单组，保存当前表单组数据
      if (this.isAdd) {
        const data = this.requestData.find(r => r.entity === this.activeTab || r.itemGroup === this.activeTab)
        if (!this.requiredCheck(data)) {
          return this.$Message.warning(`【${data.itemGroup}】${this.$t('required_tip')}`)
        } else {
          this.saveCurrentTabData(data)
        }
      }

      this.activeTab = item.entity || item.itemGroup
      this.activeItem = item
      this.addRowSource = ''
      this.addRowSourceOptions = []
      this.initTableData()
    }, 100),
    async initTableData () {
      // 当前选择tab数据
      const data = this.requestData.find(r => r.entity === this.activeTab || r.itemGroup === this.activeTab)
      this.refKeys = []
      this.calculateKeys = []
      data.title.forEach(t => {
        if (t.elementType === 'select' || t.elementType === 'wecmdbEntity') {
          this.refKeys.push(t.name)
        }
        if (t.elementType === 'calculate') {
          this.calculateKeys.push(t.name)
        }
      })
      this.formOptions = data.title

      // table数据初始化
      this.tableData = data.value.map(v => {
        // 缓存RefOptions数据，不需要每次调用
        this.refKeys.forEach(rfk => {
          if (!(v.entityData[rfk + 'Options'] && v.entityData[rfk + 'Options'].length > 0)) {
            v.entityData[rfk + 'Options'] = []
          }
        })

        // 自定义计算分析类型取值
        this.calculateKeys.forEach(key => {
          // 后台有返回值
          if (v.entityData[key] && v.entityData[key].indexOf('[') > -1) {
            let jsonData = JSON.parse(v.entityData[key]) || []
            if (jsonData.length > 0) {
              const displayNameArr = jsonData.map(item => {
                return item.displayName || ''
              })
              v.entityData[key] = displayNameArr.join('；')
            } else {
              v.entityData[key] = '' // 后端可能返回'[]'这种数据
            }
          }
          // 添加一行的数据，并且有cmdb数据id，调用接口获取
          if (!v.entityData[key] && v.addFlag && v.dataId) {
            const titleObj = data.title.find(t => t.name === key)
            this.getExpressionData(titleObj, v)
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
            this.getRefOptions(titleObj, row, index, true)
          }
        })
      })
    },
    async getRefOptions (titleObj, row, index, first) {
      // 模板自定义下拉类型
      if (titleObj.elementType === 'select' && titleObj.entity === '') {
        if (!first) return
        row[titleObj.name + 'Options'] = fixArrStrToJsonArray(titleObj.dataOptions)
        this.$set(this.tableData, index, row)
        return
      }
      // cmdb模型数据项下拉类型
      if (titleObj.elementType === 'wecmdbEntity') {
        // if (!first) return
        const [packageName, ciType] = (titleObj.dataOptions && titleObj.dataOptions.split(':')) || []
        if (!packageName || !ciType) return
        const { status, data } = await getWeCmdbOptions(packageName, ciType, {})
        if (status === 'OK') {
          row[titleObj.name + 'Options'] = data
          this.$set(this.tableData, index, row)
        }
        return
      }
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
        // 数据表单【表单隐藏标识】放到了row里面，需要删除
        if (key.indexOf('Hidden') > -1) {
          delete cache[key]
        }
      })
      this.refKeys.forEach(k => {
        delete cache[k + 'Options']
      })
      delete cache._checked
      delete cache._disabled
      // const filterValue = row[titleObj.name]
      const attrName = titleObj.entity + '__' + titleObj.name
      const attr = titleObj.id
      const params = {
        filters: [
          // {
          //   name: 'guid',
          //   operator: 'in',
          //   value: Array.isArray(filterValue) ? filterValue : [filterValue]
          // }
        ],
        paging: false,
        dialect: {
          associatedData: {
            ...cache
          }
        }
      }
      const { statusCode, data } = await getRefOptions(this.requestId, attr, params, attrName)
      if (statusCode === 'OK') {
        row[titleObj.name + 'Options'] = data
        this.$set(this.tableData, index, row)
      }
    },
    // 获取自定义计算分析类型的值
    async getExpressionData (titleObj, value) {
      const { statusCode, data } = await getExpressionData(titleObj.id, value.dataId)
      if (statusCode === 'OK') {
        const displayNameArr = data.map(item => {
          return item.displayName || ''
        })
        value.entityData[titleObj.name] = displayNameArr.join('；')
      }
    },
    // 删除行数据
    handleDeleteRow (index) {
      this.$Modal.confirm({
        title: this.$t('confirm') + this.$t('delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          this.tableData.splice(index, 1)
          this.requestData.forEach(item => {
            if (item.entity === this.activeTab || item.itemGroup === this.activeTab) {
              item.value.splice(index, 1)
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
      } else if (['exist', 'exist_empty'].includes(this.activeItem.itemGroupRule)) {
        if (this.addRowSource) {
          const source = this.addRowSourceOptions.find(i => i.id === this.addRowSource)
          const data = this.requestData.find(r => r.entity === this.activeTab || r.itemGroup === this.activeTab)
          this.handleAddRow(data, source)
          this.initTableData()
          this.$refs.addRowSelect.clearSingleSelect()
        }
      }
    },
    // 添加一条行数据
    handleAddRow (data, source = null) {
      let entityData = {}
      data.title.forEach(item => {
        // 选择已有数据添加一行，填充默认值
        if (source) {
          if (source.hasOwnProperty(item.name)) {
            entityData[item.name] = source[item.name]
          } else if (!source.hasOwnProperty(item.name) && item.defaultClear === 'no') {
            entityData[item.name] = item.defaultValue
          } else {
            entityData[item.name] = ''
          }
        } else {
          // 模板自带默认值
          if (item.defaultClear === 'no') {
            entityData[item.name] = item.defaultValue
          } else {
            entityData[item.name] = ''
          }
        }
        if (item.elementType === 'select' || item.elementType === 'wecmdbEntity') {
          entityData[item.name + 'Options'] = []
        }
      })
      const idStr = new Date().getTime().toString() + Math.floor(Math.random() * 1000)
      let obj = {
        dataId: source ? source.id || '' : '',
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
    // 获取【选择已有数据添加一行】下拉列表
    async getCmdbEntityList () {
      const { packageName, entity } = this.activeItem
      const { status, data } = await getWeCmdbOptions(packageName, entity, {})
      if (status === 'OK') {
        this.addRowSourceOptions = data || []
        // 过滤下拉框数据(1.编排类表单，下拉框只能选择系统下发的数据2.自选类表单，下拉框可以选全量数据
        // 3.下拉框数据和表单已存在的数据做ID去重)
        let ids = []
        if (this.activeItem.value && this.activeItem.value.length > 0) {
          ids = this.activeItem.value.map(item => {
            return item.dataId
          })
        }
        if (this.activeItem.itemGroupType === 'workflow') {
          this.addRowSourceOptions = this.addRowSourceOptions.filter(item =>
            this.worklfowDataIdsObj[this.activeItem.formTemplateId].includes(item.id)
          )
        }
        this.addRowSourceOptions = this.addRowSourceOptions.filter(item => !ids.includes(item.id))
      }
    },
    // 请求表单数据必填项校验
    requiredCheck (data) {
      let result = true
      let requiredName = []
      data.title.forEach(t => {
        if (t.required === 'yes') {
          requiredName.push(t.name)
        }
      })
      data.value.forEach(v => {
        requiredName.forEach(key => {
          let val = v.entityData[key]
          if (Array.isArray(val)) {
            if (val.length === 0) {
              result = false
            }
          } else {
            if (val === '' || val === undefined) {
              result = false
            }
          }
        })
      })
      return result
    },
    // 表单组必填校验
    validTable (index) {
      if (index !== '') {
        if (this.activeTab === (this.requestData[index].entity || this.requestData[index].itemGroup)) {
          return
        }
        this.activeTab = this.requestData[index].entity || this.requestData[index].itemGroup
        this.activeItem = this.requestData[index]
        this.initTableData()
        this.addRowSource = ''
        this.addRowSourceOptions = []
      }
    },
    handleTimeChange (e, value, name) {
      // 时间选择器默认填充当前时分秒
      // if (e && e.split(' ') && e.split(' ')[1] === '00:00:00') {
      //   value[name] = `${e.split(' ')[0]} ${dayjs().format('HH:mm:ss')}`
      // } else {
      //   value[name] = e
      // }
      value[name] = e
    },
    handleOpenCmdbDetail (attr, value) {
      this.cmdbInfo.title = attr.title
      this.cmdbInfo.visible = true
      this.cmdbInfo.attr = attr
      if (attr.inputType === 'diffVariable') {
        value[attr.name] = this.formatCmdbData(value[attr.name])
      }
    },
    async handleRefreshCmdbData () {
      // const { data } = await queryCiData({
      //   id: this.ciTypeId,
      //   queryObject: {
      //     dialect: { queryMode: 'new' },
      //     filters: [{ name: 'guid', operator: 'eq', value: this.cmdbInfo.attr.name }],
      //     paging: false
      //   }
      // })
      // const res = await this.formatCmdbData(data.contents[0], this.cmdbInfo.attr.name)
      // this.$nextTick(() => {
      //   this.tableDetailInfo.info = res
      //   this.tableDetailInfo.isShow = true
      // })
    },
    formatCmdbData (row, key) {
      const vari = row[key].split('\u0001=\u0001')
      const keys = vari[0].split(',\u0001')
      const values = vari[1].split(',\u0001')
      let res = []
      for (let i = 0; i < keys.length; i++) {
        res.push({
          key: (keys[i] || '').replace('\u0001', ''),
          value: (values[i] || '').replace('\u0001', '')
        })
      }
      res = res.sort((first, second) => {
        const firstKey = first.key.toLocaleUpperCase()
        const secondKey = second.key.toLocaleUpperCase()
        if (firstKey < secondKey) {
          return -1
        } else if (firstKey > secondKey) {
          return 1
        } else {
          return 0
        }
      })
      return res
    }
  }
}
</script>

<style lang="scss">
.workbench-entity-table {
  width: 100%;
  &-radio-group {
    display: flex;
    flex-wrap: wrap;
    .radio {
      padding: 5px 15px;
      border-radius: 32px;
      font-size: 14px;
      cursor: pointer;
      margin-right: 10px;
      margin-bottom: 15px;
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
    margin-bottom: 15px !important;
  }
  .ivu-form-item-label {
    word-wrap: break-word;
    padding: 10px 10px 10px 0;
  }
  .ivu-form-item-required .ivu-form-item-label:before {
    display: inline;
    margin-right: 2px;
  }
  .ivu-form-item-error-tip {
    padding-top: 2px;
    font-size: 12px;
  }
}
</style>
