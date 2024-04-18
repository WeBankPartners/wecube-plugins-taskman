<template>
  <div class="preview-drag">
    <draggable class="dragArea" :list="list" :sort="sortable" group="people" @change="log(item)">
      <div
        @click="selectElement(itemIndex, eleIndex)"
        :class="['list-group-item-', element.isActive ? 'active-zone' : '']"
        :style="{ width: (element.width / 24) * 100 + '%' }"
        v-for="(element, eleIndex) in item.attrs"
        :key="element.id"
      >
        <div
          class="custom-title"
          :style="
            ['calculate', 'textarea'].includes(element.elementType) ? 'vertical-align: top;word-break: break-all;' : ''
          "
        >
          <span v-if="element.required === 'yes'" style="color: red;">
            *
          </span>
          {{ element.title }}
        </div>
        <Input
          v-if="element.elementType === 'input'"
          :disabled="element.isEdit === 'no'"
          v-model="element.defaultValue"
          placeholder=""
          class="custom-item"
        />
        <Input
          v-if="['calculate', 'textarea'].includes(element.elementType)"
          :disabled="element.isEdit === 'no'"
          v-model="element.defaultValue"
          type="textarea"
          :rows="2"
          class="custom-item"
        />
        <Select
          v-if="element.elementType === 'select'"
          :disabled="element.isEdit === 'no'"
          class="custom-item"
          :multiple="element.multiple === 'yes'"
        >
          <Option v-for="item in element.dataOptions.split(',')" :value="item" :key="item">{{ item }}</Option>
        </Select>
        <Select
          v-if="element.elementType === 'wecmdbEntity'"
          :disabled="element.isEdit === 'no'"
          v-model="element.defaultValue"
          class="custom-item"
        ></Select>
        <DatePicker v-if="element.elementType === 'datePicker'" class="custom-item" :type="element.type"></DatePicker>
        <Button
          @click.stop="removeForm(itemIndex, eleIndex, element)"
          type="error"
          size="small"
          :disabled="$parent.isCheck === 'Y'"
          icon="ios-close"
          ghost
        ></Button>
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
    list: {
      type: Array,
      default: () => []
    },
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
    return {}
  }
}
</script>
<style lang="scss" scoped>
.preview-drag {
  .list-group-item- {
    display: inline-block;
    margin: 8px 0;
  }
}
</style>
