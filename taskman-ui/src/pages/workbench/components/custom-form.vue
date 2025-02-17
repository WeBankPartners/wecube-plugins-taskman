<!--信息表单-->
<template>
  <div class="workbench-custom-form">
    <Form :model="value" ref="form" :label-position="labelPosition" :label-width="labelWidth">
      <Row type="flex" justify="start" :gutter="20">
        <template v-for="(i, index) in formOptions">
          <Col v-if="!i.hidden" :span="i.width || 24" :key="index">
            <FormItem
              :label="i.title"
              :prop="i.name"
              :key="index"
              :required="i.required === 'yes'"
              :rules="
                i.required === 'yes'
                  ? [{ required: true, message: `${i.title}${$t('can_not_be_empty')}`, trigger: 'change' }]
                  : []
              "
              style="margin-bottom:20px;"
            >
              <!--cmdb表单类型-->
              <template v-if="i.cmdbAttr">
                <CMDBFormItem
                  :options="cmdbOptions"
                  :column="getCMDBColumn(i.name)"
                  :value="value"
                  :allSensitiveData="allSensitiveData"
                  :rowData="rowData"
                  :disabled="i.isEdit === 'no' || disabled"
                  style="width: calc(100% - 20px)"
                />
              </template>
              <!--输入框-->
              <template v-else>
                <Input
                  v-if="i.elementType === 'input'"
                  v-model.trim="value[i.name]"
                  :disabled="i.isEdit === 'no' || disabled"
                  style="width:100%;"
                ></Input>
                <Input
                  v-else-if="i.elementType === 'textarea'"
                  v-model.trim="value[i.name]"
                  type="textarea"
                  :disabled="i.isEdit === 'no' || disabled"
                  style="width:100%;"
                ></Input>
                <LimitSelect
                  v-else-if="i.elementType === 'select' || i.elementType === 'wecmdbEntity'"
                  v-model="value[i.name]"
                  :displayName="i.elementType === 'wecmdbEntity' ? 'displayName' : i.entity ? 'key_name' : 'label'"
                  :displayValue="i.elementType === 'wecmdbEntity' ? 'id' : i.entity ? 'guid' : 'value'"
                  :options="entityData[i.name + 'Options']"
                  :disabled="i.isEdit === 'no' || disabled"
                  :multiple="i.multiple === 'Y' || i.multiple === 'yes'"
                  style="width:100%;"
                >
                </LimitSelect>
                <!--自定义分析类型-->
                <Input
                  v-else-if="i.elementType === 'calculate'"
                  :value="value[i.name]"
                  type="textarea"
                  :disabled="true"
                  style="width:100%;"
                ></Input>
                <!--日期时间类型-->
                <DatePicker
                  v-else-if="i.elementType === 'datePicker'"
                  :value="value[i.name]"
                  @on-change="$event => handleTimeChange($event, value, i.name)"
                  format="yyyy-MM-dd HH:mm:ss"
                  :disabled="i.isEdit === 'no' || disabled"
                  type="datetime"
                  style="width:100%;"
                >
                </DatePicker>
              </template>
            </FormItem>
          </Col>
        </template>
      </Row>
    </Form>
  </div>
</template>

<script>
import LimitSelect from '@/pages/components/limit-select.vue'
import { getRefOptions, getWeCmdbOptions, getCmdbFormPermission } from '@/api/server'
import { evaluateCondition } from '../evaluate'
import { components } from './cmdb-form-item/action.js'
import { fixArrStrToJsonArray, deepClone } from '@/pages/util'
import CMDBFormItem from './cmdb-form-item/index.vue'
export default {
  components: {
    LimitSelect,
    CMDBFormItem
  },
  props: {
    options: {
      type: Array,
      default: () => []
    },
    value: {
      type: Object,
      default: () => {}
    },
    rowData: {
      type: Object,
      default: () => {}
    },
    requestId: {
      type: String,
      default: ''
    },
    labelWidth: {
      type: Number,
      default: 100
    },
    labelPosition: {
      type: String,
      default: 'left'
    },
    disabled: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      refKeys: [],
      entityData: {},
      formOptions: [],
      cmdbOptions: [],
      cmdbSensitiveKeysArr: [], // cmdb表单敏感字段name集合
      allSensitiveData: [] // 当前所有敏感字段数据
    }
  },
  computed: {
    getCMDBColumn () {
      return function (name) {
        return this.cmdbOptions.find(i => i.propertyName === name) || {}
      }
    }
  },
  watch: {
    value: {
      handler (val) {
        if (val) {
          Object.keys(val).forEach(key => {
            this.entityData[key] = val[key]
          })
          this.hideFormItems()
        }
      },
      deep: true,
      immediate: true
    },
    options: {
      handler (val) {
        if (val && val.length) {
          this.formOptions = this.options
          this.hideFormItems()
          // select类型集合
          this.refKeys = []
          val.forEach(t => {
            if (t.elementType === 'select' || t.elementType === 'wecmdbEntity') {
              this.refKeys.push(t.name)
            }
          })
          // value数据初始化
          this.refKeys.forEach(rfk => {
            // 缓存RefOptions数据，不需要每次调用
            if (!(this.entityData[rfk + 'Options'] && this.entityData[rfk + 'Options'].length > 0)) {
              this.$set(this.entityData, rfk + 'Options', [])
            }
          })
          // 下拉类型数据初始化(待优化，调用接口太多)
          this.refKeys.forEach(rfk => {
            if (!(this.entityData[rfk + 'Options'] && this.entityData[rfk + 'Options'].length > 0)) {
              const titleObj = val.find(f => f.name === rfk)
              this.getRefOptions(titleObj)
            }
          })
          // cmdb表单属性初始化
          this.cmdbSensitiveKeysArr = []
          const cmdbOptions = []
          const formOptions = deepClone(this.formOptions)
          formOptions.forEach(item => {
            if (item.cmdbAttr) {
              item.cmdbAttr = JSON.parse(item.cmdbAttr)
              item.cmdbAttr.editable = item.isEdit
              const cmdbObj = Object.assign({}, { ...item.cmdbAttr, titleObj: item })
              cmdbOptions.push(cmdbObj)
              if (item.cmdbAttr.sensitive === 'yes') {
                this.cmdbSensitiveKeysArr.push(item.name)
              }
            }
          })
          if (cmdbOptions && cmdbOptions.length > 0) {
            this.getCMDBInitData(cmdbOptions)
          }
          // cmdb表单权限初始化
          if (Array.isArray(this.cmdbSensitiveKeysArr) && this.cmdbSensitiveKeysArr.length > 0) {
            this.getCmdbFormPermission()
          }
        }
      },
      immediate: true,
      deep: true
    }
  },
  methods: {
    // 表单隐藏逻辑
    hideFormItems () {
      Object.keys(this.value).forEach(key => {
        // 表单隐藏逻辑
        const find = this.formOptions.find(i => i.name === key) || {}
        if (find.hiddenCondition && find.required === 'no') {
          const conditions = find.hiddenCondition || []
          find.hidden = conditions.every(j => {
            return evaluateCondition(j, this.value[j.name])
          })
        }
      })
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
    async getRefOptions (titleObj) {
      // taskman模板管理配置的普通下拉类型(值用逗号拼接)
      if (titleObj.elementType === 'select' && titleObj.entity === '') {
        this.$set(this.entityData, titleObj.name + 'Options', fixArrStrToJsonArray(titleObj.dataOptions))
        return
      }
      // taskman模板管理配置的引用下拉类型
      if (titleObj.elementType === 'wecmdbEntity') {
        const [packageName, ciType] = (titleObj.dataOptions && titleObj.dataOptions.split(':')) || []
        const { status, data } = await getWeCmdbOptions(packageName, ciType, {})
        if (status === 'OK') {
          this.$set(this.entityData, titleObj.name + 'Options', data)
        }
        return
      }
      // if (titleObj.refEntity === '') {
      //   row[titleObj.name + 'Options'] = titleObj.selectList
      //   this.$set(this.tableData, index, row)
      //   return
      // }
      let cache = JSON.parse(JSON.stringify(this.entityData))
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
        // 将对象类型转为字符串
        if (typeof cache[key] === 'object') {
          cache[key] = JSON.stringify(cache[key])
        }
        // 将number类型转为字符串
        if (typeof cache[key] === 'number') {
          cache[key] = cache[key].toString()
        }
      })
      this.refKeys.forEach(k => {
        delete cache[k + 'Options']
      })
      delete cache._checked
      delete cache._disabled
      // const filterValue = this.entityData[titleObj.name]
      const attrName = titleObj.entity + '__' + titleObj.name
      const attr = titleObj.id
      const params = {
        // filters: [
        //   {
        //     name: 'guid',
        //     operator: 'in',
        //     value: Array.isArray(filterValue) ? filterValue : [filterValue]
        //   }
        // ],
        paging: false,
        dialect: {
          associatedData: {
            ...cache
          }
        }
      }
      const { statusCode, data } = await getRefOptions(this.requestId, attr, params, attrName)
      if (statusCode === 'OK') {
        this.$set(this.entityData, titleObj.name + 'Options', data)
      }
    },
    // cmdb表单属性初始化
    getCMDBInitData (data) {
      this.cmdbOptions = []
      let columns = []
      for (let index = 0; index < data.length; index++) {
        let renderKey = data[index].propertyName
        if (!['decommissioned', 'notCreated'].includes(data[index].status)) {
          if (['select', 'multiSelect'].includes(data[index].inputType) && data[index].selectList !== '') {
            const { titleObj } = data[index] || { titleObj: {} }
            const attrName = titleObj.entity + '__' + titleObj.name
            const attr = titleObj.id
            // 异步获取cmdb select和multiSelect下拉框的值
            getRefOptions(this.requestId, attr, {}, attrName)
              .then(res => {
                if (res.statusCode === 'OK') {
                  this.cmdbOptions[index].options = res.data.map(item => {
                    return {
                      label: item.key_name,
                      value: item.guid
                    }
                  })
                }
              })
          }
          const columnItem = {
            ...data[index],
            title: data[index].name,
            key: renderKey,
            inputKey: data[index].propertyName,
            inputType: data[index].inputType,
            referenceId: data[index].referenceId,
            disEditor: !data[index].isEditable,
            disAdded: !data[index].isEditable,
            placeholder: data[index].name,
            component: 'Input',
            referenceFilter: data[index].referenceFilter,
            ciType: { id: data[index].referenceId, name: data[index].name },
            type: 'text',
            isMultiple: ['multiSelect', 'multiRef'].includes(data[index].inputType),
            ...components[data[index].inputType],
            options: data[index].options,
            requestId: this.requestId,
            refKeys: this.refKeys
          }
          columns.push(columnItem)
        }
      }
      this.cmdbOptions = columns
    },
    // 获取cmdb表单权限
    async getCmdbFormPermission () {
      const ciType = this.rowData.entityName
      const values = [this.rowData]
      let params = []
      const guidArr = (values && values.map(v => {
        return {
          dataId: v.dataId,
          tmpId: v.id,
          entityData: v.entityData
        }
      })) || []
      this.cmdbSensitiveKeysArr.forEach(name => {
        guidArr.forEach(item => {
          params.push(
            {
              attrName: name,
              attrVal: item.entityData[name] || '',
              ciType,
              guid: item.dataId,
              requestId: this.requestId,
              tmpId: item.tmpId
            }
          )
        })
      })
      const { statusCode, data } = await getCmdbFormPermission(params)
      if (statusCode === 'OK') {
        this.allSensitiveData = data || []
      }
    }
  }
}
</script>
<style lang="scss">
.workbench-custom-form {
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
