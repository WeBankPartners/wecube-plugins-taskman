<template>
  <div class="taskman-base-search">
    <div class="form">
      <Form :inline="true" :model="value" label-position="right">
        <template v-for="(i, index) in options">
          <FormItem :label="i.label || ''" :prop="i.key" :label-width="i.labelWidth || 0" :key="index">
            <!--输入框-->
            <Input
              v-if="i.component === 'input'"
              v-model="value[i.key]"
              :placeholder="i.placeholder"
              clearable
              :style="{ width: i.width || 200 + 'px' }"
            ></Input>
            <!--下拉选择-->
            <Select
              v-else-if="i.component === 'select'"
              v-model="value[i.key]"
              :placeholder="i.placeholder"
              clearable
              :multiple="i.multiple || false"
              :filterable="i.filterable || false"
              :max-tag-count="1"
              :style="{ width: i.width || 200 + 'px' }"
            >
              <template v-for="(j, idx) in i.list">
                <Option :key="idx" :value="j.value">{{ j.label }}</Option>
              </template>
            </Select>
            <!--获取接口的下拉选择-->
            <Select
              v-else-if="i.component === 'remote-select'"
              v-model="value[i.key]"
              @on-open-change="getRemoteData(i)"
              :placeholder="i.placeholder"
              clearable
              :multiple="i.multiple || false"
              :filterable="i.filterable || false"
              :max-tag-count="1"
              :style="{ width: i.width || 200 + 'px' }"
            >
              <template v-for="(j, idx) in i.list">
                <Option :key="idx" :value="j.value">{{ j.label }}</Option>
              </template>
            </Select>
            <!--标签组-->
            <RadioGroup
              v-else-if="i.component === 'radio-group'"
              v-model="value[i.key]"
              @on-change="$emit('search')"
              style="margin-right:32px;"
            >
              <Radio v-for="(j, idx) in i.list" :label="j.value" :key="idx" border>{{ j.label }}</Radio>
            </RadioGroup>
            <!--自定义时间选择器-->
            <div v-else-if="i.component === 'custom-time'" class="custom-time">
              <RadioGroup
                v-if="dateType !== 4"
                v-model="dateType"
                @on-change="handleDateTypeChange(i.key)"
                style="margin-top:-2px;"
              >
                <Radio v-for="(j, idx) in dateTypeList" :label="j.value" :key="idx" border>{{ j.label }}</Radio>
              </RadioGroup>
              <div v-else>
                <DatePicker
                  :value="value[i.key]"
                  @on-change="
                    val => {
                      handleDateRange(val, i.key)
                    }
                  "
                  type="daterange"
                  placement="bottom-end"
                  format="yyyy-MM-dd"
                  :placeholder="$t('created_time')"
                  style="width: 200px"
                />
                <Icon
                  size="18"
                  style="cursor:pointer;"
                  type="md-close-circle"
                  @click="
                    dateType = ''
                    formData[i.key] = []
                  "
                />
              </div>
            </div>
          </FormItem>
        </template>
      </Form>
    </div>
    <div class="button-group">
      <Button @click="handleSearch" size="small" type="primary">查询</Button>
      <Button @click="handleReset" size="small" style="margin-top:10px;">重置</Button>
    </div>
  </div>
</template>

<script>
import dayjs from 'dayjs'
export default {
  props: {
    options: {
      type: Array,
      default: () => []
    },
    value: {
      type: Object,
      default: () => {}
    }
  },
  computed: {
    formData () {
      return this.value
    }
  },
  data () {
    return {
      dateType: '', // 1-近一个月2-近半年3-近一年4-自定义
      dateTypeList: [
        { label: '近一个月', value: 1 },
        { label: '近半年', value: 2 },
        { label: '近一年', value: 3 },
        { label: '自定义', value: 4 }
      ]
    }
  },
  watch: {},
  methods: {
    handleSearch () {
      this.$emit('search')
    },
    // 重置表单
    handleReset () {
      const resetObj = {}
      Object.keys(this.formData).forEach(key => {
        if (Array.isArray(this.formData[key])) {
          resetObj[key] = []
        } else {
          resetObj[key] = ''
        }
        // 点击清空按钮需要给默认值的表单选项
        const initOptions = this.options.filter(i => i.initValue !== undefined)
        initOptions.forEach(i => {
          resetObj[i.key] = i.initValue
        })
      })
      this.$emit('input', resetObj)
    },
    // 自定义时间控件转化时间格式值
    handleDateTypeChange (key) {
      this.formData[key] = []
      const cur = dayjs().format('YYYY-MM-DD')
      if (this.dateType === 1) {
        const pre = dayjs()
          .subtract(1, 'month')
          .format('YYYY-MM-DD')
        this.formData[key] = [pre, cur]
      } else if (this.dateType === 2) {
        const pre = dayjs()
          .subtract(6, 'month')
          .format('YYYY-MM-DD')
        this.formData[key] = [pre, cur]
      } else if (this.dateType === 3) {
        const pre = dayjs()
          .subtract(1, 'year')
          .format('YYYY-MM-DD')
        this.formData[key] = [pre, cur]
      }
      // 同步更新父组件form数据
      this.$emit('input', this.formData)
    },
    handleDateRange (dateArr, key) {
      this.formData[key] = [...dateArr]
      this.$emit('input', this.formData)
    },
    // 获取远程下拉框数据
    async getRemoteData (i) {
      const res = await i.remote()
      this.$set(i, 'list', res)
    }
  }
}
</script>

<style lang="scss">
.taskman-base-search {
  display: flex;
  padding-bottom: 10px;
  .ivu-form-item {
    margin-bottom: 10px !important;
  }
  .ivu-radio {
    display: none;
  }
  .ivu-radio-wrapper {
    border-radius: 5px;
    height: 30px !important;
    line-height: 30px !important;
    padding: 0 10px;
    font-size: 12px;
    color: #000;
    background: #f5f8fa;
    border: none;
  }
  .ivu-radio-wrapper-checked.ivu-radio-border {
    border-color: #2d8cf0;
    background: #2d8cf0;
    color: #fff;
  }
  .form {
    flex: 1;
  }
  .button-group {
    width: 65px;
    display: flex;
    flex-direction: column;
    border-left: 1px solid #0000000f;
    padding-left: 20px;
    box-sizing: content-box;
    button {
      height: 28px;
      line-height: 28px;
      font-size: 13px;
    }
  }
}
</style>