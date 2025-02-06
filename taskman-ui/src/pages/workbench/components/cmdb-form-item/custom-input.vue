<template>
  <div class="taskman-custom-input">
    <!--敏感字段-->
    <div v-if="column.sensitive === 'yes'" class="flex-row">
      <template v-if="operation === 'detail'">
        <Input v-if="isShowReal" :value="initValue === attrs.value ? getRealValue : attrs.value" disabled placeholder="" style="width: 500px;" />
        <Input v-else :value="initValue" disabled placeholder="" style="width: 500px;" />
      </template>
      <template v-if="operation === 'edit'">
        <Input v-bind="attrs" placeholder="" style="width: 500px;" @input="handleInputChange"></Input>
      </template>
      <Button
        @click="handlePreviewData"
        :disabled="getCmdbQueryPermission === false"
        :icon="isShowReal ? 'md-eye-off' : 'md-eye'"
        size="small"
      ></Button>
      <Button
        @click="handleEditData"
        :disabled="attrs.disabled"
        type="primary"
        icon="md-create"
        size="small"
      ></Button>
      <Button
        v-if="column.autofillable === 'yes' && column.autoFillType === 'suggest'"
        @click="$emit('input', 'suggest#')"
        :disabled="attrs.disabled"
        icon="md-checkmark"
        size="small"
      ></Button>
    </div>
    <!--非敏感字段-->
    <div v-else class="flex-row">
      <Input v-bind="attrs" placeholder="" style="width: 500px;" @input="handleInputChange"></Input>
      <Button
        v-if="column.autofillable === 'yes' && column.autoFillType === 'suggest'"
        @click="$emit('input', 'suggest#')"
        :disabled="attrs.disabled"
        icon="md-checkmark"
        size="small"
      ></Button>
    </div>
  </div>
</template>
<script>
export default {
  props: {
    // Input props属性
    attrs: {
      type: Object,
      default: () => {}
    },
    // 表单项属性
    column: {
      type: Object,
      default: () => {}
    },
    // 所有敏感字段原文及权限
    allSensitiveData: {
      type: Array,
      default: () => []
    },
    // 当前行数据ID
    dataId: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      isShowReal: false, // 是否显示原文
      operation: 'detail', // 查看or编辑
      initValue: '' // 初始值
    }
  },
  computed: {
    getCmdbQueryPermission () {
      const obj = this.allSensitiveData.find(item => item.attrName === this.column.inputKey && item.guid === this.dataId) || {}
      return obj.queryPermission
    },
    getRealValue () {
      const obj = this.allSensitiveData.find(item => item.attrName === this.column.inputKey && item.guid === this.dataId) || {}
      return obj.value
    }
  },
  mounted () {
    this.initValue = this.attrs.value
  },
  methods: {
    handleInputChange (val) {
      this.$emit('input', val)
    },
    handlePreviewData () {
      this.operation = 'detail'
      this.isShowReal = !this.isShowReal
    },
    handleEditData () {
      this.operation = 'edit'
      this.isShowReal = true
    }
  }
}
</script>

<style scoped lang="scss">
.flex-row {
  display: flex;
  align-items: center;
}
Button {
  margin-left: 5px;
}
</style>
