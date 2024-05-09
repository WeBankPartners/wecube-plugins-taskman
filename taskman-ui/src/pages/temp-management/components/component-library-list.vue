<template>
  <div class="component-library-list">
    <Row>
      <Col :span="12">
        <Select v-model="query.type" placeholder="表单类型">
          <Option v-for="(i, index) in typeList" :value="i" :key="index">{{ i }}</Option>
        </Select>
      </Col>
      <Col :span="12">
        <Select v-model="query.createdBy" placeholder="创建人">
          <Option v-for="(i, index) in typeList" :value="i" :key="index">{{ i }}</Option>
        </Select>
      </Col>
      <Col :span="24">
        <div class="group">
          <Input v-model="query.name" placeholder="组件名" class="input" />
          <Icon type="ios-list" :size="36" class="icon" @click="handleViewTable" />
        </div>
      </Col>
    </Row>
    <div class="list">
      <div v-for="i in 10" :key="i" class="list-item">{{ i }}</div>
    </div>
    <ComponentLibraryModal :visible.sync="visible" :isAdd="false"></ComponentLibraryModal>
  </div>
</template>
<script>
import ComponentLibraryModal from './component-library-modal.vue'
export default {
  components: {
    ComponentLibraryModal
  },
  data () {
    return {
      visible: false,
      query: {
        type: '',
        createdBy: '',
        name: ''
      },
      typeList: [1, 2, 3]
    }
  },
  methods: {
    handleViewTable () {
      this.visible = true
    }
  }
}
</script>

<style scoped lang="scss">
.component-library-list {
  width: 100%;
  .group {
    display: flex;
    align-items: center;
    margin-top: 10px;
    .input {
      flex: 1;
    }
    .icon {
      cursor: pointer;
    }
  }
  .list {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    &-item {
      width: calc(50% - 5px);
      margin-top: 10px;
      height: 60px;
      border: 1px solid black;
    }
  }
}
</style>
