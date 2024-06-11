<template>
  <div class="custom-drag">
    <draggable
      class="dragArea"
      :list="customElement"
      :group="{ name: 'people', pull: 'clone', put: false }"
      :sort="sortable"
      :clone="clone"
    >
      <div class="list-group-item-" style="width: 100%" v-for="element in customElement" :key="element.id">
        <Input v-if="element.elementType === 'input'" :placeholder="$t('t_input')" />
        <Input v-if="element.elementType === 'textarea'" type="textarea" :placeholder="$t('textare')" />
        <Select v-if="element.elementType === 'select'" :placeholder="$t('select')"></Select>
        <Select v-if="element.elementType === 'wecmdbEntity'" :placeholder="$t('tw_entity_data_items')"></Select>
        <DatePicker
          v-if="element.elementType === 'datePicker'"
          :type="element.type"
          :placeholder="$t('tw_date_picker')"
          style="width:100%"
        ></DatePicker>
        <div v-if="element.elementType === 'group'" style="width: 100%; height: 80px; border: 1px solid #5ea7f4">
          <span style="margin: 8px; color: #bbbbbb"> Item Group </span>
        </div>
      </div>
    </draggable>
  </div>
</template>

<script>
import draggable from 'vuedraggable'
export default {
  components: {
    draggable
  },
  props: {
    sortable: {
      type: Boolean,
      default: false
    },
    clone: {
      type: Function,
      default: () => {}
    }
  },
  data () {
    return {
      customElement: [
        // 预制自定义表单项目
        {
          id: 1,
          name: 'input',
          title: 'Input',
          elementType: 'input',
          defaultValue: '',
          defaultClear: 'no',
          // tag: '',
          itemGroup: '',
          itemGroupName: '',
          packageName: '',
          entity: '',
          width: 24,
          dataOptions: '',
          regular: '',
          inDisplayName: 'yes',
          isEdit: 'yes',
          multiple: 'no',
          selectList: [],
          isRefInside: 'no',
          required: 'no',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: '',
          refEntity: '',
          refPackageName: '',
          hiddenCondition: [] // 隐藏条件
        },
        {
          id: 3,
          name: 'textarea',
          title: 'Textarea',
          elementType: 'textarea',
          defaultClear: 'no',
          defaultValue: '',
          // tag: '',
          itemGroup: '',
          itemGroupName: '',
          packageName: '',
          entity: '',
          width: 24,
          dataOptions: '',
          regular: '',
          inDisplayName: 'yes',
          isEdit: 'yes',
          multiple: 'no',
          selectList: [],
          isRefInside: 'no',
          required: 'no',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: '',
          refEntity: '',
          refPackageName: '',
          hiddenCondition: [] // 隐藏条件
        },
        {
          id: 2,
          name: 'select',
          title: 'Select',
          elementType: 'select',
          defaultValue: '',
          defaultClear: 'no',
          // tag: '',
          itemGroup: '',
          itemGroupName: '',
          packageName: '',
          entity: '',
          width: 24,
          dataOptions: '',
          regular: '',
          inDisplayName: 'yes',
          isEdit: 'yes',
          multiple: 'no',
          selectList: [],
          isRefInside: 'no',
          required: 'no',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: '',
          refEntity: '',
          refPackageName: '',
          controlSwitch: 'no',
          hiddenCondition: [] // 隐藏条件
        },
        {
          id: 5,
          name: 'wecmdbEntity',
          title: 'WecmdbEntity',
          elementType: 'wecmdbEntity',
          defaultValue: '',
          defaultClear: 'no',
          // tag: '',
          itemGroup: '',
          itemGroupName: '',
          packageName: '',
          entity: '',
          width: 24,
          dataOptions: '',
          regular: '',
          inDisplayName: 'yes',
          isEdit: 'yes',
          multiple: 'no',
          selectList: [],
          isRefInside: 'no',
          required: 'no',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: '',
          refEntity: '',
          refPackageName: '',
          controlSwitch: 'no',
          hiddenCondition: [] // 隐藏条件
        },
        {
          id: 6,
          name: 'datePicker',
          title: 'datePicker',
          elementType: 'datePicker',
          type: 'datetime',
          defaultValue: '',
          defaultClear: 'no',
          // tag: '',
          itemGroup: '',
          itemGroupName: '',
          packageName: '',
          entity: '',
          width: 24,
          dataOptions: '',
          regular: '',
          inDisplayName: 'yes',
          isEdit: 'yes',
          multiple: 'no',
          selectList: [],
          isRefInside: 'no',
          required: 'no',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: '',
          refEntity: '',
          refPackageName: '',
          hiddenCondition: [] // 隐藏条件
        }
      ]
    }
  }
}
</script>
<style lang="scss" scoped>
.custom-drag {
  .list-group-item- {
    display: inline-block;
    margin: 8px 0;
  }
}
</style>
