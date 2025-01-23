<template>
  <div class="taskman-custom-input">
    <!--敏感字段-->
    <div v-if="column.sensitive === 'yes'" class="flex-row">
      <Input v-bind="attrs" placeholder="" style="width: 50%;" @input="handleInputChange"></Input>
      <Button
        @click="handlePreviewData"
        style="margin-left: 5px;"
        :icon="visible ? 'md-eye-off' : 'md-eye'"
      ></Button>
      <Button
        @click="handleEditData"
        style="margin-left: 5px;"
        type="primary"
        icon="md-create"
      ></Button>
      <Button
        v-if="column.autofillable === 'yes' && column.autoFillType === 'suggest'"
        @click="$emit('input', 'suggest#')"
        icon="md-checkmark"
        style="margin-left: 5px;"
      ></Button>
    </div>
    <!--非敏感字段-->
    <div v-else class="flex-row">
      <Input v-bind="attrs" placeholder="" style="width: 50%;" @input="handleInputChange"></Input>
      <Button
        v-if="column.autofillable === 'yes' && column.autoFillType === 'suggest'"
        @click="$emit('input', 'suggest#')"
        icon="md-checkmark"
        style="margin-left: 5px;"
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
    }
  },
  data () {
    return {
      visible: false
    }
  },
  methods: {
    handleInputChange (val) {
      this.$emit('input', val)
    },
    handlePreviewData () {
      this.visible = !this.visible
    },
    handleEditData () {}
  }
}
</script>

<style scoped lang="scss">
.flex-row {
  display: flex;
}
</style>
