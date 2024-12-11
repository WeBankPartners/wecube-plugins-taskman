<!--
 * @Author: wanghao7717 792974788@qq.com
 * @Date: 2024-12-10 17:20:25
 * @LastEditors: wanghao7717 792974788@qq.com
 * @LastEditTime: 2024-12-11 18:09:12
-->
<template>
  <div class="cmdb-select">
    <Select
      :value="value"
      :multiple="isMultiple"
      :disabled="disabled"
      filterable
      clearable
      @on-change="changeValue"
    >
      <Option v-for="item in opts" :value="item.value || ''" :key="item.value">{{ item.label }}</Option>
    </Select>
  </div>
</template>
<script>
export default {
  name: 'WeCMDBSelect',
  props: {
    value: {},
    isMultiple: { default: () => false },
    options: { default: () => [] },
    filterParams: {},
    disabled: { default: () => false },
    enumId: { default: () => null }
  },
  data () {
    return {
      filterOpts: [],
      enumOpts: []
    }
  },
  watch: {},
  computed: {
    opts () {
      if (this.filterParams) {
        return this.filterOpts
      } else if (this.enumId) {
        return this.enumOpts
      } else {
        return this.options
      }
    }
  },
  mounted () {},
  methods: {
    changeValue (val) {
      this.$emit('input', val || null)
      this.$emit('change', val || null)
    }
  }
}
</script>
