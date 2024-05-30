<template>
  <Modal
    v-model="visible"
    title="隐藏条件"
    :mask-closable="false"
    :closable="false"
    :width="700"
    class="hidden-condition-modal"
  >
    <div class="content">
      <Row :gutter="5" style="margin-bottom:5px;">
        <Col span="6">表单项</Col>
        <Col span="6">符号</Col>
        <Col span="12">{{ $t('value') }}</Col>
      </Row>
      <Row v-for="(i, index) in hiddenCondition" :key="index" :gutter="5">
        <Form :model="hiddenCondition[index]" :rules="rule" :ref="'form' + index">
          <!--表单项-->
          <Col :span="6">
            <FormItem label="" prop="name">
              <Select v-model="hiddenCondition[index].name" :disabled="disabled" @on-change="handleNameChange">
                <Option v-for="j in nameList" :key="j.id" :value="j.name">{{ j.title }}</Option>
              </Select>
            </FormItem>
          </Col>
          <!--符号-->
          <Col :span="6">
            <FormItem label="" prop="operator">
              <Select
                v-model="hiddenCondition[index].operator"
                :disabled="disabled"
                @on-change="handleOperatorChange($event, index)"
              >
                <Option v-for="(j, index) in getOperatorList(selectAttrs[index])" :key="index" :value="j.value">{{
                  j.label
                }}</Option>
              </Select>
            </FormItem>
          </Col>
          <!--隐藏值-->
          <Col v-if="!['empty', 'notEmpty'].includes(hiddenCondition[index].operator)" :span="10">
            <FormItem label="" prop="value">
              <template v-if="selectAttrs[index].elementType === 'datePicker'">
                <DatePicker
                  v-if="hiddenCondition[index].operator === 'range'"
                  :value="hiddenCondition[index].value"
                  :disabled="disabled"
                  @on-change="
                    val => {
                      hiddenCondition[index].value = val
                    }
                  "
                  type="datetimerange"
                  placement="bottom-end"
                  format="yyyy-MM-dd HH:mm:ss"
                  style="width:100%;"
                />
                <DatePicker
                  v-else
                  :value="hiddenCondition[index].value"
                  @on-change="
                    val => {
                      hiddenCondition[index].value = val
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
                v-else
                v-model="hiddenCondition[index].value"
                :placeholder="
                  ['containsAll', 'containsAny', 'notContains'].includes(hiddenCondition[index].operator)
                    ? 'eg：a,b,c'
                    : ''
                "
                :disabled="disabled"
              ></Input>
              <!-- <Select
                v-if="selectAttrs[index].elementType === 'select' || selectAttrs[index].elementType === 'wecmdbEntity'"
                v-model="editElement.hiddenCondition[index].value"
                :disabled="disabled"
                placeholder="值"
                @on-change="paramsChanged"
                @open-change="getRefOptions($event, index)"
              >
                <Option
                  v-for="(j, idx)  refOptionMap[editElement.hiddenCondition[index].name]"
                  :key="idx"
                  :value="selectAttrs[index].elementType === 'wecmdbEntity' ? 'displayName' : j.entity ? 'key_name' : 'label'"
                >
                  {{ j }}
                </Option>
              </Select> -->
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
          <Button type="primary" ghost @click="handleAddRow" size="small" icon="md-add"></Button>
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
import { deepClone } from '../../util'
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
    name: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      visible: false,
      hiddenCondition: [], // 隐藏条件
      selectAttrs: [], // 隐藏条件对应表单项配置
      nameList: [], // 表单项下拉列表
      operatorList: [
        // 符号下拉列表
        { label: '等于', value: 'eq', condition: ['input', 'singleSelect', 'dateTime'] },
        { label: '不等于', value: 'neq', condition: ['input', 'singleSelect', 'dateTime'] },
        { label: '小于', value: 'lt', condition: ['input', 'singleSelect', 'dateTime'] },
        { label: '大于', value: 'gt', condition: ['input', 'singleSelect', 'dateTime'] },
        { label: '包含', value: 'contains', condition: ['input', 'singleSelect'] },
        { label: '匹配开始', value: 'startsWith', condition: ['input', 'singleSelect'] },
        { label: '匹配结束', value: 'endsWith', condition: ['input', 'singleSelect'] },
        { label: '包含全部', value: 'containsAll', condition: ['singleSelect', 'multipleSelect'] },
        { label: '包含任意', value: 'containsAny', condition: ['singleSelect', 'multipleSelect'] },
        { label: '不包含', value: 'notContains', condition: ['singleSelect', 'multipleSelect'] },
        { label: '在范围内', value: 'range', condition: ['dateTime'] },
        { label: '为空', value: 'empty', condition: ['input', 'singleSelect', 'multipleSelect', 'dateTime'] },
        { label: '不为空', value: 'notEmpty', condition: ['input', 'singleSelect', 'multipleSelect', 'dateTime'] }
      ],
      rule: {
        name: [{ required: true, message: '请选择表单项', trigger: 'change' }],
        operator: [{ required: true, message: '请选择符号', trigger: 'change' }],
        value: [
          {
            validator: (rule, value, callback) => {
              if (!value || (Array.isArray(value) && value.length === 0)) {
                return callback(new Error('请输入过滤值'))
              } else {
                callback()
              }
            },
            required: true,
            message: '请输入过滤值',
            trigger: 'change'
          }
        ]
      }
    }
  },
  computed: {
    // 获取符号可选下拉列表
    getOperatorList () {
      return function (item) {
        if (!item) return this.operatorList
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
        return this.operatorList.filter(i => i.condition.includes(type))
      }
    }
  },
  watch: {
    hiddenCondition: {
      handler (val) {
        if (val && val.length > 0) {
          this.selectAttrs = []
          const deleteIndexArr = []
          this.hiddenCondition.forEach((i, index) => {
            // 新加一行空数据处理
            if (!i.name || !i.operator) {
              this.selectAttrs.push({ elementType: 'input' })
            } else {
              const findIndex = this.finalElement[0].attrs.findIndex(j => j.name === i.name)
              if (findIndex > -1) {
                this.selectAttrs.push(this.finalElement[0].attrs[findIndex])
              } else {
                deleteIndexArr.push(index)
              }
            }
          })
          // 删除不存在预览区的过滤条件
          for (let idx of deleteIndexArr) {
            this.hiddenCondition.splice(idx, 1)
          }
        }
      },
      deep: true
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
      this.nameList = this.finalElement[0].attrs.filter(i => i.name !== this.name)
      this.$nextTick(() => {
        this.$refs.form0[0].resetFields()
      })
    },
    handleNameChange () {
      // 过滤已选表单项
      // this.nameList = this.nameList.filter(i => {
      //   const arr = this.hiddenCondition.map(m => m.name)
      //   return !arr.includes(i.name)
      // })
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
    },
    handleDeleteItem (index) {
      this.hiddenCondition.splice(index, 1)
    },
    async getRefOptions (flag, index) {
      if (!flag) return
      const name = this.value.hiddenCondition[index].name
      const item = this.finalElement[0].attrs.find(i => i.name === name)
      // 模板自定义下拉类型
      if (item.elementType === 'select' && item.entity === '') {
        this.$set(this.refOptionMap, item.name, JSON.parse(item.dataOptions || '[]'))
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
