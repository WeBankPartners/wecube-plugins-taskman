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
        <span class="text">{{ type === 'json' ? originData : (formaMultiData && formaMultiData.length === 0 ? $t('tw_no_data') :  formaMultiData) }}</span>
        <Button
          v-if="type === 'json'"
          @click="showDetail = true"
          type="primary"
          ghost
          icon="md-eye"
        ></Button>
      </div>
    </Tooltip>
    <!--编辑-->
    <Button v-else type="primary" @click="showConfig" :disabled="disabled">{{ $t('tw_config') }}</Button>
    <!--非json编辑框-->
    <Modal v-model="showModal" :title="$t('tw_config')">
      <div class="multiconfig-modal-content">
        <template v-for="(item, itemIndex) in multiData">
          <div :key="itemIndex" style="margin:4px">
            <InputNumber
              v-if="type === 'number'"
              :max="99999999"
              :min="-99999999"
              style="width:400px"
              :precision="0"
              v-model="item.value"
            />
            <Input v-else v-model="item.value" :maxlength="255" show-word-limit style="width:400px"></Input>
            <Button @click="addItem" type="primary" size="small" icon="ios-add" style="margin:0 4px"></Button>
            <Button @click="deleteItem(itemIndex)" v-if="multiData.length !== 1" size="small" type="error" icon="ios-trash"></Button>
          </div>
        </template>
      </div>
      <template #footer>
        <Button @click="showModal = false">{{ $t('cancel') }}</Button>
        <Button @click="confirmData" type="primary">{{ $t('confirm') }}</Button>
      </template>
    </Modal>
    <!--json编辑框-->
    <Modal :z-index="2000" v-model="showJsonModal" :title="$t('tw_json_edit')" width="800">
      <Button type="primary" @click="addNewJson">{{ $t('tw_add_group') }}</Button>
      <div style="max-height:500px; overflow:auto">
        <template v-for="(item, itemIndex) in originData">
          <Tree :ref="'jsonTree' + itemIndex" :jsonData="item" :key="itemIndex"></Tree>
        </template>
      </div>
      <template #footer>
        <Button @click="showJsonModal = false">{{ $t('cancel') }}</Button>
        <Button @click="confirmJsonData" type="primary">{{ $t('confirm') }}</Button>
      </template>
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
  components: {
    Tree,
    JsonViewer
  },
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
      const arr = this.multiData.filter(i => i.value !== '') || []
      const res = arr.map(item => {
        if (this.type === 'number') {
          return Number(item.value)
        } else {
          return item.value
        }
      })
      return res
    }
  },
  mounted () {
    this.initData()
  },
  methods: {
    initData () {
      const data = JSON.parse(JSON.stringify(this.data))
      let tmp = data || []
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
    confirmJsonData () {
      this.showJsonModal = false
      this.showModal = false
      this.$emit('input', this.originData)
    },
    addNewJson () {
      this.originData.push({})
    },
    showConfig () {
      if (this.type === 'json') {
        this.showJsonModal = true
      } else {
        this.showModal = true
      }
      this.initData()
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
      const emptyFlag = this.multiData.some(item => item.value === '')
      if (emptyFlag) {
        this.$Message.warning('数据不能为空')
      } else {
        const res = this.multiData.map(item => {
          if (this.type === 'number') {
            return Number(item.value)
          } else {
            return item.value
          }
        })
        this.showJsonModal = false
        this.showModal = false
        this.$emit('input', res)
      }
    }
  }
}
</script>

<style lang="scss">
.taskman-cmdb-multi-config {
  .inline {
    display: flex;
    align-items: center;
    height: 34px;
    .text {
      font-size: 13px;
      color:#515a6e;
      display: block;
      max-width: 450px;
      text-overflow: ellipsis;
      overflow: hidden;
      white-space: nowrap;
      margin-right: 6px;
    }
  }
}
.multiconfig-modal-content {
  max-height: 360px;
  overflow-y: auto;
}
</style>
