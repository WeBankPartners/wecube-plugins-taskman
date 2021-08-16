<template>
  <div>
    <Select
      :value="value"
      :multiple="isMultiple"
      filterable
      clearable
      @on-change="changeValue"
    >
      <Option v-for="item in list" :value="item.guid" :key="item.guid">{{
        item.displayName
      }}</Option>
    </Select>
  </div>
</template>
<script>
const DEFAULT_TAG_NUMBER = 2;
import { queryReferenceEnumCodes, getTargetOptions, getAllDataModels } from "@/api/server";
export default {
  name: "PluginSelect",
  props: {
    value: {},
    isMultiple: { default: () => false },
    options: { default: () => [] },
    maxTags: { default: () => DEFAULT_TAG_NUMBER },
    packageName: {required: false},
    entityName: {required: true}
  },
  data() {
    return {
      list: [],
      allEntityList: []
    };
  },
  watch: {},
  computed: {
  },
  mounted() {
    this.getAllDataModels()
  },
  methods: {
    async getAllDataModels () {
    const { data, status } = await getAllDataModels()
    if (status === 'OK') {
      this.allEntityType = data.map(_ => {
        // handle result sort by name
        const pluginPackageEntities = _.entities ? _.entities.sort(function (a, b) {
            var s = a.name.toLowerCase()
            var t = b.name.toLowerCase()
            if (s < t) return -1
            if (s > t) return 1
          }) : []
        this.allEntityList = this.allEntityList.concat(pluginPackageEntities)
        return {
          ..._,
          pluginPackageEntities: pluginPackageEntities
        }
      })
      this.getTargetOptions()
    }
  },
    async getTargetOptions () {
      const packageName = this.allEntityList.find(_ => _.name === this.entityName).packageName
      const { status, data } = await getTargetOptions(packageName, this.entityName)
      if (status === 'OK') {
          this.list = data
        }
    },
    formatOptions() {},
    changeValue(val) {
      this.$emit("input", val);
      this.$emit("change", val);
    },
  }
};
</script>
