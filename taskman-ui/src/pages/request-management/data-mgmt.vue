<template>
  <div class=" ">
    <Table size="small" :columns="tableColumns" :data="tableData"></Table>
  </div>
</template>

<script>
export default {
  name: '',
  data () {
    return {
      oriData: null,
      formData: [],
      tableColumns: [],
      tableData: []
    }
  },
  methods: {
    initData (data) {
      this.oriData = data
      this.formData = data.title
      this.tableColumns = data.title.map(t => {
        const col = {
          title: t.title,
          key: t.name
        }
        return col
      })
      this.tableData = data.value.map(v => {
        v.entityData._id = v.id
        return v.entityData
      })
    }
  },
  components: {}
}
</script>

<style scoped lang="scss"></style>
