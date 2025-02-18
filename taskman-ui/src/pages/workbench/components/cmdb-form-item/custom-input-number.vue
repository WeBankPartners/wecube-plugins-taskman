<template>
  <div class="taskman-custom-input-number">
    <!--敏感字段-->
    <div v-if="column.sensitive && column.sensitive === 'yes'" class="flex-row">
      <Input
        v-if="isShowReal"
        :value="originVal === attrs.value ? getRealValue : attrs.value"
        disabled
        placeholder=""
        style="width: 500px;"
      />
      <Input v-else :value="!attrs.value ? '' : '******'" disabled placeholder="" style="width: 500px;" />
      <Button
        @click="handlePreviewData"
        :disabled="getCmdbQueryPermission === false && originVal === attrs.value"
        :icon="isShowReal ? 'md-eye-off' : 'md-eye'"
      ></Button>
      <Button
        v-if="attrs.disabled === false"
        @click="handleEditData"
        type="primary"
        icon="md-create"
      ></Button>
    </div>
    <!--非敏感字段-->
    <div v-else class="flex-row">
      <InputNumber
        v-bind="attrs"
        :max="99999999"
        :min="-99999999"
        :precision="0"
        placeholder=""
        @input="handleInputChange"
        style="width:500px;"
      />
    </div>
    <!--敏感字段编辑弹框-->
    <Modal v-model="isShowEditModal" :title="$t('edit')">
      <InputNumber
        v-model="editValue"
        :max="99999999"
        :min="-99999999"
        :precision="0"
        placeholder=""
        style="width: 450px;"
      />
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
