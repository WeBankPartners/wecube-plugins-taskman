<template>
  <div class="taskman-custom-input">
    <!--敏感字段-->
    <div v-if="column.sensitive && column.sensitive === 'yes'" class="flex-row">
      <Input
        v-if="isShowReal"
        :value="originVal === attrs.value ? getRealValue : attrs.value"
        disabled
        placeholder=""
      />
      <Input v-else :value="!attrs.value ? '' : '******'" disabled placeholder="" />
      <Button
        @click="handlePreviewData"
        :disabled="getCmdbQueryPermission === false && originVal === attrs.value"
        type="primary"
        ghost
        :icon="isShowReal ? 'md-eye-off' : 'md-eye'"
      ></Button>
      <Button
        v-if="attrs.disabled === false"
        @click="handleEditData"
        type="primary"
        ghost
        icon="md-build"
      ></Button>
      <Button
        v-if="column.autofillable === 'yes' && column.autoFillType === 'suggest'"
        @click="$emit('input', 'suggest#')"
        :disabled="attrs.disabled"
        type="primary"
        ghost
        icon="md-checkmark"
      ></Button>
    </div>
    <!--非敏感字段-->
    <div v-else class="flex-row">
      <Input v-bind="attrs" placeholder="" @input="handleInputChange"></Input>
      <Button
        v-if="column.autofillable === 'yes' && column.autoFillType === 'suggest'"
        @click="$emit('input', 'suggest#')"
        :disabled="attrs.disabled"
        type="primary"
        ghost
        icon="md-checkmark"
      ></Button>
    </div>
    <!--敏感字段编辑弹框-->
    <Modal v-model="isShowEditModal" :title="$t('edit')">
      <Input
        v-model="editValue"
        placeholder=""
        style="width: 450px;"
      ></Input>
      <div slot="footer">
        <Button @click="isShowEditModal = false">{{ $t('cancel') }}</Button>
        <Button @click="handleSave" type="primary">{{ $t('save') }}</Button>
      </div>
    </Modal>
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
    // 当前行数据
    rowData: {
      type: Object,
      default: () => {}
    }
  },
  data () {
    return {
      isShowReal: false, // 是否显示原文
      originVal: '', // 初始值
      editValue: '', // 编辑值
      isShowEditModal: false
    }
  },
  computed: {
    getCmdbQueryPermission () {
      const obj = this.allSensitiveData.find(item => {
        if (this.rowData.dataId) {
          return item.attrName === this.column.inputKey && item.guid === this.rowData.dataId
        } else {
          return item.attrName === this.column.inputKey && item.tmpId === this.rowData.id
        }
      }) || {}
      return obj.queryPermission
    },
    getRealValue () {
      const obj = this.allSensitiveData.find(item => {
        if (this.rowData.dataId) {
          return item.attrName === this.column.inputKey && item.guid === this.rowData.dataId
        } else {
          return item.attrName === this.column.inputKey && item.tmpId === this.rowData.id
        }
      }) || {}
      return obj.value
    }
  },
  mounted () {
    this.originVal = this.attrs.value
  },
  methods: {
    handleInputChange (val) {
      this.$emit('input', val)
    },
    handlePreviewData () {
      this.isShowReal = !this.isShowReal
    },
    handleEditData () {
      this.editValue = ''
      this.isShowEditModal = true
    },
    handleSave () {
      this.$emit('input', this.editValue)
      this.isShowEditModal = false
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
