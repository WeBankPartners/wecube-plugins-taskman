<template>
  <div class="taskman-cmdb-multi-config">
    <!--查看-->
    <Tooltip
      v-if="disabled"
      max-width="350"
      style="width: 100%"
      placement="bottom-start"
      :content="type === 'json' ? JSON.stringify(originData) : JSON.stringify(formaMultiData)"
    >
      <div class="inline">
        <span class="text">{{ type === 'json' ? originData : formaMultiData }}</span>
        <Icon v-if="type === 'json'" type="md-eye" @click="showDetail = true" class="operation-icon-confirm" />
      </div>
    </Tooltip>
    <!--编辑-->
    <Button v-else type="primary" @click="showConfig" :disabled="disabled">{{ $t('tw_config') }}</Button>
    <!--非json编辑框-->
    <Modal v-model="showModal" :title="$t('tw_config')" @on-ok="confirmData" @on-cancel="cancel">
      <template v-for="(item, itemIndex) in multiData">
        <div :key="itemIndex" style="margin:4px">
          <Input v-model="item.value" :type="type" v-if="type !== 'json'" style="width:360px"></Input>
          <Button @click="addItem" type="primary" size="small" icon="ios-add" style="margin:0 4px"></Button>
          <Button @click="deleteItem(itemIndex)" v-if="multiData.length !== 1" size="small" type="error" icon="ios-trash"></Button>
        </div>
      </template>
    </Modal>
    <!--json编辑框-->
    <Modal :z-index="2000" v-model="showJsonModal" :title="$t('tw_json_edit')" @on-ok="confirmJsonData" width="800">
      <Button type="primary" @click="addNewJson">新增一组</Button>
      <div style="max-height:500px; overflow:auto">
        <template v-for="(item, itemIndex) in originData">
          <Tree :ref="'jsonTree' + itemIndex" :jsonData="item" :key="itemIndex"></Tree>
        </template>
      </div>
    </Modal>
    <!--详情弹框-->
    <Modal :z-index="2000" v-model="showDetail" :title="title" @on-ok="showDetail = false" width="700">
      <div style="max-height:500px;overflow:auto">
        <JsonViewer :value="type === 'json' ? originData : multiData" :expand-depth="5" class="taskman-json-viewer"></JsonViewer>
      </div>
    </Modal>
  </div>
</template>

<script>
import Tree from './tree'
import JsonViewer from 'vue-json-viewer'
export default {
  name: '',
  data () {
    return {
      showModal: false,
      multiData: [
        {
          value: ''
        }
      ],

      showJsonModal: false,
      originData: [],
      showDetail: false
    }
  },
  props: ['title', 'inputKey', 'data', 'type', 'disabled'],
  computed: {
    formaMultiData () {
      const res = this.multiData.map(item => {
        if (this.type === 'number') {
          return Number(item.value)
        }
        return item.value
      })
      return res
    }
  },
  mounted () {
    let tmp = this.data ? this.data : []
    if (this.type === 'json') {
      this.originData = tmp || []
    } else {
      this.multiData =
        tmp &&
        tmp.map(d => {
          return {
            value: d
          }
        })
      if (this.multiData.length === 0) {
        this.multiData = [
          {
            value: ''
          }
        ]
      }
    }
  },
  methods: {
    confirmJsonData () {
      this.$emit('input', this.originData, this.inputKey)
      this.showJsonModal = false
      this.showModal = false
    },
    addNewJson () {
      this.originData.push({})
    },
    showConfig () {
      if (this.type === 'json') {
        this.showJsonModal = true
      }
      this.showModal = true
    },
    addItem () {
      this.multiData.push({
        value: ''
      })
    },
    deleteItem (index) {
      this.multiData.splice(index, 1)
    },
    confirmData () {
      const res = this.multiData.map(item => {
        if (this.type === 'number') {
          return Number(item.value)
        }
        return item.value
      })
      this.showJsonModal = false
      this.showModal = false
      this.$emit('input', res, this.inputKey)
    },
    cancel () {}
  },
  components: {
    Tree,
    JsonViewer
  }
}
</script>

<style lang="scss">
.taskman-cmdb-multi-config {
  .inline {
    display: flex;
    align-items: center;
    .text {
      font-size: 13px;
      color:#515a6e;
      display: block;
      max-width: 400px;
      text-overflow: ellipsis;
      overflow: hidden;
      white-space: nowrap;
    }
    .operation-icon-confirm {
      font-size: 16px;
      border: 1px solid #57a3f3;
      color: #57a3f3;
      border-radius: 4px;
      width: 32px;
      line-height: 28px;
      cursor: pointer;
      margin-left: 5px;
    }
  }
}
</style>
