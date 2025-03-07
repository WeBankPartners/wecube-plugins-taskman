<!--
 * @Author: wanghao7717 792974788@qq.com
 * @Date: 2024-10-18 17:55:45
 * @LastEditors: wanghao7717 792974788@qq.com
 * @LastEditTime: 2025-02-19 16:14:36
-->
<template>
  <div class="cmdb-entity-table">
    <template v-if="column.component === 'WeCMDBCIPassword'">
      <WeCMDBCIPassword
        :formData="column"
        :panalData="value"
        :allSensitiveData="allSensitiveData"
        :rowData="rowData"
        :disabled="isGroupEditDisabled(column, value)"
        @input="(v) => {value[column.inputKey] = v}"
      />
    </template>
    <template v-else-if="column.component === 'WeCMDBDiffVariable'">
      <Diffvariable :data="value[column.inputKey]" />
    </template>
    <template v-else-if="column.component === 'Input' && column.inputType === 'multiText'">
      <MultiConfig
        :title="column.title"
        :inputKey="column.inputKey"
        :disabled="isGroupEditDisabled(column, value)"
        :data="JSON.parse(JSON.stringify(value[column.inputKey]))"
        type="text"
        @input="(v) => {value[column.inputKey] = v}"
      ></MultiConfig>
    </template>
    <template v-else-if="column.component === 'Input' && column.inputType === 'multiInt'">
      <MultiConfig
        :title="column.title"
        :inputKey="column.inputKey"
        :disabled="isGroupEditDisabled(column, value)"
        :data="JSON.parse(JSON.stringify(value[column.inputKey]))"
        type="number"
        @input="(v) => {value[column.inputKey] = v}"
      ></MultiConfig>
    </template>
    <template v-else-if="column.component === 'Input' && column.inputType === 'multiObject'">
      <MultiConfig
        :title="column.title"
        :inputKey="column.inputKey"
        :disabled="isGroupEditDisabled(column, value)"
        :data="JSON.parse(JSON.stringify(value[column.inputKey]))"
        type="json"
        @input="(v) => {value[column.inputKey] = v}"
      ></MultiConfig>
    </template>
    <template v-else-if="column.component === 'Input' && column.inputType === 'autofillRule'">
      <Input
        v-bind="getInputProps(column, value)"
        placeholder=""
        @input="(v) => {setValueHandler(v.trim(), column, value)}"
      ></Input>
    </template>
    <template v-else-if="column.component === 'Input' && column.inputType === 'int'">
      <div style="display:flex;">
        <InputNumber
          v-bind="getInputProps(column, value)"
          :max="99999999"
          :min="-99999999"
          :precision="0"
          placeholder=""
          @input="(v) => {setValueHandler(v, column, value)}"
        />
      </div>
    </template>
    <template v-else-if="column.component === 'Input' && column.inputType !== 'object'">
      <CustomInput
        :attrs="getInputProps(column, value)"
        :column="column"
        :allSensitiveData="allSensitiveData"
        :rowData="rowData"
        @input="(v) => {setValueHandler(v.trim(), column, value)}"
      ></CustomInput>
    </template>
    <template v-else-if="column.component === 'Input' && column.inputType === 'object'">
      <JsonConfig
        :title="column.title"
        :inputKey="column.inputKey"
        :disabled="isGroupEditDisabled(column, value)"
        :jsonData="typeof value[column.inputKey] === 'object' ? value[column.inputKey] : JSON.parse(value[column.inputKey] || '{}')"
        @input="(v) => {setValueHandler(v, column, value)}"
      ></JsonConfig>
    </template>
    <template v-else-if="column.component === 'WeCMDBSelect'">
      <WeCMDBSelect
        :value="column.isRefreshable
          ? column.inputType === 'multiSelect'
            ? []
            : ''
          : column.inputType === 'multiSelect'
            ? Array.isArray(value[column.inputKey])
              ? value[column.inputKey]
              : ''
            : formatValue(column, value[column.inputKey])"
        :disabled="isGroupEditDisabled(column, value)"
        :filterParams="column.referenceFilter
          ? {
            attrId: column.ciTypeAttrId,
            params: value
          }
          : null"
        :isMultiple="column.isMultiple"
        :options="column.options"
        :enumId="column.referenceId ? column.referenceId : null"
        @input="(v) => {setValueHandler(v, column, value)}"
        @change="(v) => {
          if (column.onChange) {
            this.$emit('getGroupList', v)
          }
        }"
      ></WeCMDBSelect>
    </template>
    <template v-else-if="column.component === 'WeCMDBRefSelect'">
      <WeCMDBRefSelect
        :value="column.isRefreshable
          ? ''
          : column.inputType === 'multiRef'
            ? Array.isArray(value[column.inputKey])
              ? value[column.inputKey]
              : []
            : value[column.inputKey] || ''"
        :disabled="isGroupEditDisabled(column, value)"
        :filterParams="column.referenceFilter
          ? {
            attrId: column.ciTypeAttrId,
            params: value
          }
          : {
            params: value
          }"
        :ciTypeAttrId="column.ciTypeAttrId"
        :column="column"
        @input="(v) => {setValueHandler(v, column, value)}"
        @change="(v) => {
          if (column.onChange) {
            this.$emit('getGroupList', v)
          }
        }"
      ></WeCMDBRefSelect>
    </template>
    <template v-else-if="column.component === 'DatePicker'">
      <DatePicker
        :value="value[column.inputKey]"
        :disabled="isGroupEditDisabled(column, value)"
        :editable="column.editable === 'yes'"
        type="datetime"
        @input="(v) => {setValueHandler(v, column, value)}"
        @change="(v) => {
          if (column.onChange) {
            this.$emit('getGroupList', v)
          }
        }"
      >
      </DatePicker>
    </template>
  </div>
</template>
<script>
import dayjs from 'dayjs'
import WeCMDBCIPassword from './ci-password.vue'
import MultiConfig from './multi-config.vue'
import JsonConfig from './json-config.vue'
import WeCMDBSelect from './cmdb-select.vue'
import WeCMDBRefSelect from './cmdb-ref-select/index'
import Diffvariable from './diff-variable.vue'
import CustomInput from './custom-input.vue'
export default {
  components: {
    WeCMDBCIPassword,
    MultiConfig,
    JsonConfig,
    WeCMDBSelect,
    WeCMDBRefSelect,
    Diffvariable,
    CustomInput
  },
  props: {
    options: {
      type: Array,
      default: () => []
    },
    column: {
      type: Object,
      default: () => {}
    },
    value: {
      type: Object,
      default: () => {}
    },
    allSensitiveData: {
      type: Array,
      default: () => []
    },
    rowData: {
      type: Object,
      default: () => {}
    },
    disabled: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {}
  },
  computed: {
    isGroupEditDisabled () {
      return function (attr, item) {
        let attrGroupEditDisabled = false
        if (attr.editGroupControl === 'yes') {
          if (attr.editGroupValues.length > 0) {
            let groups = JSON.parse(attr.editGroupValues)
            for (let idx = 0; idx < groups.length; idx++) {
              let group = groups[idx]
              if (attrGroupEditDisabled) {
                break
              }
              const findAttr = this.options.find(el => {
                if (el.propertyName === group.key) {
                  return true
                }
                return false
              })
              if (findAttr && group.value.length > 0) {
                if (!item[findAttr.propertyName]) {
                  // 控制字段未赋值，禁用当前字段
                  attrGroupEditDisabled = true
                } else if (Array.isArray(item[findAttr.propertyName])) {
                  let attrValues = item[findAttr.propertyName]
                  let intersect = attrValues.filter(v => {
                    return group.value.indexOf(v) > -1
                  })
                  // 控制字段是数组且与设置的数据没有交集，禁用当前字段
                  if (intersect.length === 0) {
                    attrGroupEditDisabled = true
                  }
                } else {
                  // 控制字段是单值且不在分组范围，应当禁用当前字段
                  if (group.value.indexOf(item[findAttr.propertyName]) < 0) {
                    attrGroupEditDisabled = true
                  }
                }
              }
            }
          }
        }
        let attrEditDisabled = attr.editable === 'no' || (attr.autofillable === 'yes' && attr.autoFillType === 'forced')
        return attrEditDisabled || attrGroupEditDisabled || this.disabled
      }
    },
    getInputProps () {
      return function (column, value) {
        return {
          ...column,
          editable: column.editable === 'yes',
          type: column.inputType === 'int' ? 'number' : 'text',
          disabled: this.isGroupEditDisabled(column, value),
          value: column.inputType === 'int' ? Number(value[column.inputKey]) : value[column.inputKey]
        }
      }
    },
    formatValue () {
      return function (column, value) {
        // for edit 多选数据会在保存时转为','拼接
        if (column.component === 'WeCMDBSelect' && column.isMultiple) {
          if (value) {
            return Array.isArray(value) ? value : value.split(',')
          } else {
            return null
          }
        } else {
          return value
        }
      }
    }
  },
  methods: {
    // 赋值操作
    setValueHandler (v, col, row) {
      let attrsWillReset = []
      if (['select', 'ref', 'extRef', 'multiSelect', 'multiRef'].indexOf(col.inputType) > -1) {
        this.options.forEach(_ => {
          if (_.uiFormOrder > col.uiFormOrder && _.referenceFilter) {
            if (['multiSelect', 'multiRef'].indexOf(_.inputType) >= 0) {
              attrsWillReset.push({
                propertyName: _.propertyName,
                value: []
              })
            } else if (['select', 'ref', 'extRef'].indexOf(_.inputType) >= 0) {
              attrsWillReset.push({
                propertyName: _.propertyName,
                value: ''
              })
            }
          }
        })
      } else if (['datetime'].indexOf(col.inputType) >= 0) {
        v = dayjs(v).format('YYYY-MM-DD HH:mm:ss')
      } else if (['autofillRule'].includes(col.inputType)) {
        if (v !== '') {
          attrsWillReset.push({
            propertyName: col.inputKey,
            value: v
          })
        }
      }
      row[col.inputKey] = v
      attrsWillReset.forEach(attr => {
        row[attr.propertyName] = attr.value
      })
    }
  }
}
</script>
<style lang="scss" scoped>
.cmdb-entity-table {
  width: 100%;
}
</style>
