<template>
  <Modal
    v-model="visible"
    :title="$t('tw_hidden_condition')"
    :mask-closable="false"
    :closable="false"
    :width="700"
    class="hidden-condition-modal"
  >
    <div class="content">
      <Row :gutter="5" style="margin-bottom:5px;">
        <Col span="6">{{ $t('form_item') }}</Col>
        <Col span="6">{{ $t('tw_symbol') }}</Col>
        <Col span="12">{{ $t('value') }}</Col>
      </Row>
      <Row v-for="(i, index) in hiddenCondition" :key="index" :gutter="5">
        <Form :model="i" :rules="rule" :ref="'form' + index">
          <!--表单项-->
          <Col :span="6">
            <FormItem label="" prop="name">
              <Select v-model="i.name" :disabled="disabled" @on-change="handleNameChange(index)">
                <Option v-for="j in nameList" :key="j.id" :value="j.name">{{ j.title }}</Option>
              </Select>
            </FormItem>
          </Col>
          <!--符号-->
          <Col :span="6">
            <FormItem label="" prop="operator">
              <Select v-model="i.operator" :disabled="disabled" @on-change="handleOperatorChange($event, index)">
                <Option v-for="(j, index) in selectAttrs[index].operatorList" :key="index" :value="j.value">{{
                  j.label
                }}</Option>
              </Select>
            </FormItem>
          </Col>
          <!--值-->
          <Col v-if="!['empty', 'notEmpty'].includes(i.operator)" :span="10">
            <FormItem label="" prop="value">
              <template v-if="selectAttrs[index].elementType === 'datePicker'">
                <DatePicker
                  v-if="i.operator === 'range'"
                  :value="i.value"
                  :disabled="disabled"
                  @on-change="
                    val => {
                      i.value = val
                    }
                  "
                  type="datetimerange"
                  placement="bottom-end"
                  format="yyyy-MM-dd HH:mm:ss"
                  style="width:100%;"
                />
                <DatePicker
                  v-else
                  :value="i.value"
                  @on-change="
                    val => {
                      i.value = val
                    }
                  "
                  format="yyyy-MM-dd HH:mm:ss"
                  type="datetime"
                  style="width:100%;"
                >
                </DatePicker>
              </template>
              <!--输入框-->
              <Input
                v-else-if="
                  ['textarea', 'input'].includes(selectAttrs[index].elementType) ||
                    ['lt', 'gt', 'contains', 'startsWith', 'endsWith'].includes(i.operator)
                "
                v-model="i.value"
                :placeholder="['containsAll', 'containsAny', 'notContains'].includes(i.operator) ? 'eg：a,b,c' : ''"
                :disabled="disabled"
              ></Input>
              <Select
                v-else-if="['select', 'wecmdbEntity'].includes(selectAttrs[index].elementType)"
                v-model="i.value"
                :disabled="disabled"
                :multiple="['eq', 'neq'].includes(i.operator) ? false : true"
                clearable
                filterable
                placeholder=""
                @on-open-change="getRefOptions($event, index)"
              >
                <Option
                  v-for="(j, idx) in refOptionMap[i.name]"
                  :key="idx"
                  :label="
                    selectAttrs[index].elementType === 'wecmdbEntity'
                      ? j.displayName
                      : selectAttrs[index].entity
                      ? j.key_name
                      : j.label
                  "
                  :value="
                    selectAttrs[index].elementType === 'wecmdbEntity'
                      ? j.id
                      : selectAttrs[index].entity
                      ? j.guid
                      : j.value
                  "
                />
              </Select>
            </FormItem>
          </Col>
          <Col :span="2">
            <Button
              type="error"
              ghost
              @click="handleDeleteItem(index)"
              size="small"
              style="vertical-align:sub;cursor:pointer;"
              icon="md-trash"
            ></Button>
          </Col>
        </Form>
      </Row>
      <Row :gutter="5" style="cursor:pointer;">
        <Col :span="2" :offset="22">
          <Button type="success" ghost @click="handleAddRow" size="small" icon="md-add"></Button>
        </Col>
      </Row>
    </div>
    <template #footer>
      <Button @click="visible = false">{{ $t('cancel') }}</Button>
      <Button @click="handleSubmit" type="primary">{{ $t('confirm') }}</Button>
    </template>
  </Modal>
</template>

<script>
import { getWeCmdbOptions } from '@/api/server'
import { deepClone, fixArrStrToJsonArray } from '../../util'
export default {
  props: {
    finalElement: {
      type: Array,
      default: () => []
    },
    disabled: {
      type: Boolean,
      default: false
    },
    editElement: {
      type: Object,
      default: () => {}
    }
  },
  data () {
    return {
      visible: false,
      hiddenCondition: [], // 隐藏条件
      selectAttrs: [], // 隐藏条件对应表单项配置
      nameList: [], // 表单项下拉列表
      refOptionMap: {},
      operatorList: [
        // 符号下拉列表
        { label: this.$t('tw_symbol_eq'), value: 'eq', condition: ['input', 'singleSelect', 'dateTime'] },
        { label: this.$t('tw_symbol_neq'), value: 'neq', condition: ['input', 'singleSelect', 'dateTime'] },
        { label: this.$t('tw_symbol_lt'), value: 'lt', condition: ['input', 'singleSelect', 'dateTime'] },
        { label: this.$t('tw_symbol_gt'), value: 'gt', condition: ['input', 'singleSelect', 'dateTime'] },
        { label: this.$t('tw_symbol_contains'), value: 'contains', condition: ['input', 'singleSelect'] },
        { label: this.$t('tw_symbol_startsWith'), value: 'startsWith', condition: ['input', 'singleSelect'] },
        { label: this.$t('tw_symbol_endsWith'), value: 'endsWith', condition: ['input', 'singleSelect'] },
        {
          label: this.$t('tw_symbol_containsAll'),
          value: 'containsAll',
          condition: ['singleSelect', 'multipleSelect']
        },
        {
          label: this.$t('tw_symbol_containsAny'),
          value: 'containsAny',
          condition: ['singleSelect', 'multipleSelect']
        },
        {
          label: this.$t('tw_symbol_notContains'),
          value: 'notContains',
          condition: ['singleSelect', 'multipleSelect']
        },
        { label: this.$t('tw_symbol_range'), value: 'range', condition: ['dateTime'] },
        {
          label: this.$t('tw_symbol_empty'),
          value: 'empty',
          condition: ['input', 'singleSelect', 'multipleSelect', 'dateTime']
        },
        {
          label: this.$t('tw_symbol_notEmpty'),
          value: 'notEmpty',
          condition: ['input', 'singleSelect', 'multipleSelect', 'dateTime']
        }
      ],
      rule: {
        name: [{ required: true, message: this.$t('tw_formItem_placeholder'), trigger: 'change' }],
        operator: [{ required: true, message: this.$t('tw_symbol_placeholder'), trigger: 'change' }],
        value: [
          {
            validator: (rule, value, callback) => {
              if (!value || (Array.isArray(value) && value.length === 0)) {
                return callback(new Error(this.$t('tw_filterValue_placeholder')))
              } else {
                callback()
              }
            },
            required: true,
            message: this.$t('tw_filterValue_placeholder'),
            trigger: 'change'
          }
        ]
      }
    }
  },
  methods: {
    initData (arr) {
      this.hiddenCondition = deepClone(arr)
      if (this.hiddenCondition.length === 0) {
        this.hiddenCondition.push({
          name: '',
          operator: '',
          value: ''
        })
      }
      this.visible = true
      // 表单项下拉列表
      this.nameList = this.finalElement[0].attrs.filter(i => {
        if (i.cmdbAttr) {
          const cmdbAttr = JSON.parse(i.cmdbAttr)
          return i.name !== this.editElement.name && cmdbAttr.sensitive !== 'yes'
        } else {
          return i.name !== this.editElement.name
        }
      })
      // 刷新符号下拉列表
      this.refreshSelectAttrs(true)
      // 获取值为select类型下拉列表
      this.selectAttrs.forEach((item, index) => {
        if (['select', 'wecmdbEntity'].includes(item.elementType)) {
          this.getRefOptions(true, index)
        }
      })
    },
    refreshSelectAttrs (need) {
      this.selectAttrs = []
      this.hiddenCondition.forEach(i => {
        // 新加一行空数据处理
        if (!i.name) {
          this.selectAttrs.push({ elementType: 'input' })
        } else {
          const findIndex = this.finalElement[0].attrs.findIndex(j => j.name === i.name)
          this.selectAttrs.push(this.finalElement[0].attrs[findIndex])
        }
      })
      // 刷新符号下拉列表
      if (need) {
        this.selectAttrs.forEach(item => {
          let type = ''
          if (['input', 'textarea'].includes(item.elementType)) {
            type = 'input'
          } else if (item.elementType === 'datePicker') {
            type = 'dateTime'
          } else if (['select', 'wecmdbEntity'].includes(item.elementType)) {
            if (['Y', 'yes'].includes(item.multiple)) {
              type = 'multipleSelect'
            } else {
              type = 'singleSelect'
            }
          }
          item.operatorList = this.operatorList.filter(i => i.condition.includes(type))
        })
      }
    },
    handleNameChange (index) {
      // 清空符号和值
      this.hiddenCondition[index].operator = ''
      this.hiddenCondition[index].value = ''
      this.refreshSelectAttrs(true)
    },
    handleOperatorChange (val, index) {
      if (val === 'range') {
        this.hiddenCondition[index].value = []
      } else {
        this.hiddenCondition[index].value = ''
      }
    },
    handleAddRow () {
      // 处理接口返回null类型
      if (!this.hiddenCondition) {
        this.hiddenCondition = []
      }
      this.hiddenCondition.push({
        name: '',
        operator: '',
        value: ''
      })
      this.refreshSelectAttrs()
    },
    handleDeleteItem (index) {
      this.hiddenCondition.splice(index, 1)
      this.refreshSelectAttrs()
    },
    async getRefOptions (flag, index) {
      if (!flag) return
      const name = this.hiddenCondition[index].name
      const item = this.finalElement[0].attrs.find(i => i.name === name)
      // 模板自定义下拉类型
      if (item.elementType === 'select' && item.entity === '') {
        this.$set(this.refOptionMap, item.name, fixArrStrToJsonArray(item.dataOptions))
        return
      }
      // cmdb下发
      if (item.elementType === 'select' && item.entity) {
        if (!item.refPackageName || !item.refEntity) return
        const { status, data } = await getWeCmdbOptions(item.refPackageName, item.refEntity, {})
        if (status === 'OK') {
          this.$set(this.refOptionMap, item.name, data || [])
        }
        return
      }
      // 模型数据项
      if (item.elementType === 'wecmdbEntity') {
        const [packageName, ciType] = (item.dataOptions && item.dataOptions.split(':')) || []
        if (!packageName || !ciType) return
        const { status, data } = await getWeCmdbOptions(packageName, ciType, {})
        if (status === 'OK') {
          this.$set(this.refOptionMap, item.name, data || [])
        }
      }
    },
    handleSubmit () {
      let validArr = []
      this.hiddenCondition.forEach((_, index) => {
        const key = `form${index}`
        this.$refs[key][0].validate(valid => {
          validArr.push(valid)
        })
      })
      const validFlag = validArr.every(i => i === true)
      if (validFlag) {
        this.visible = false
        this.$emit('updateData', this.hiddenCondition)
      }
    }
  }
}
</script>

<style lang="scss">
.hidden-condition-modal {
  .content {
    min-height: 100px;
  }
  .ivu-form-item {
    margin-bottom: 5px;
  }
}
</style>
