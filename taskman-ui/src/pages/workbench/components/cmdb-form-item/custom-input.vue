<template>
  <div class="taskman-custom-input">
    <!--敏感字段-->
    <div v-if="column.sensitive === 'yes'" class="flex-row">
      <template v-if="operation === 'detail'">
        <Input v-if="isShowReal" disabled :value="getRealValue" placeholder="" style="width: 50%;" />
        <Input v-else disabled value="******" placeholder="" style="width: 50%;" />
      </template>
      <template v-if="operation === 'edit'">
        <Input v-bind="attrs" placeholder="" style="width: 50%;" @input="handleInputChange"></Input>
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
      <Input v-bind="attrs" placeholder="" style="width: 50%;" @input="handleInputChange"></Input>
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
    attrs: {
      type: Object,
      default: () => {}
    },
    column: {
      type: Object,
      default: () => {}
    },
    allSensitiveData: {
      type: Array,
      default: () => []
    },
    dataId: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      isShowReal: false, // 显示原文
      operation: 'detail' // 查看or编辑
    }
  },
  computed: {
    getCmdbQueryPermission () {
      const obj = this.allSensitiveData.find(item => item.attrName === this.column.name && item.guid === this.dataId) || {}
      return obj.queryPermission
    },
    getRealValue () {
      const obj = this.allSensitiveData.find(item => item.attrName === this.column.name && item.guid === this.dataId) || {}
      return obj.value
    }
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
