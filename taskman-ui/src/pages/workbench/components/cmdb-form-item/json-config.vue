<!--
 * @Author: wanghao7717 792974788@qq.com
 * @Date: 2024-11-19 10:23:32
 * @LastEditors: wanghao7717 792974788@qq.com
 * @LastEditTime: 2025-01-06 16:36:18
-->
<template>
  <div class="taskman-cmdb-json-config">
    <!--查看-->
    <Tooltip v-if="disabled" max-width="350" style="width: 100%" placement="bottom-start" :content="jsonDataString">
      <div class="inline">
        <span class="text">{{ jsonDataString || $t('tw_no_data') }}</span>
        <Icon type="md-eye" @click="showDetail = true" class="operation-icon-confirm" />
      </div>
    </Tooltip>
    <!--编辑-->
    <Button v-else type="primary" :disabled="disabled" @click="showTreeConfig">{{ $t('tw_config') }}</Button>
    <!--编辑弹框-->
    <Modal :z-index="2000" v-model="showEdit" :title="$t('tw_json_edit')" width="800">
      <Button type="primary" v-if="isArray" @click="addNewJson">{{ $t('tw_add_group') }}</Button>
      <div style="max-height:500px;overflow-y:auto;">
        <template v-for="(item, itemIndex) in originData">
          <Tree :ref="'jsonTree' + itemIndex" :jsonData="item" :key="itemIndex"></Tree>
        </template>
      </div>
      <template #footer>
        <Button @click="showEdit = false">{{ $t('cancel') }}</Button>
        <Button @click="confirmJsonData" type="primary">{{ $t('confirm') }}</Button>
      </template>
    </Modal>
    <!--查看弹框-->
    <Modal :z-index="2000" v-model="showDetail" :title="title" @on-ok="showDetail = false" width="700">
      <div style="max-height:500px;overflow:auto">
        <JsonViewer :value="JSON.parse(jsonDataString || '{}')" :expand-depth="5" class="taskman-json-viewer"></JsonViewer>
      </div>
    </Modal>
  </div>
</template>

<script>
import JsonViewer from 'vue-json-viewer'
import Tree from './tree'
export default {
  components: {
    Tree,
    JsonViewer
  },
  data () {
    return {
      showEdit: false,
      jsonDataString: '',
      isArray: false,
      originData: [],
      finalData: null,
      last: null,
      showDetail: false
    }
  },
  props: ['inputKey', 'jsonData', 'disabled', 'title'],
  mounted () {
    this.initData()
  },
  methods: {
    initData () {
      const jsonData = JSON.parse(JSON.stringify(this.jsonData))
      this.isArray = Array.isArray(jsonData)
      this.originData = []
      if (this.isArray) {
        this.originData = jsonData
      } else {
        this.originData.push(jsonData || {})
      }
      const jsonDataString = JSON.stringify(jsonData)
      this.jsonDataString = jsonDataString === '""' ? '' : jsonDataString
    },
    showTreeConfig () {
      this.showEdit = true
      this.initData()
    },
    confirmJsonData () {
      if (this.isArray) {
        this.finalData = []
        const len = this.originData.length
        for (let i = 0; i < len; i++) {
          this.finalData.push(this.$refs['jsonTree' + i][0].jsonJ)
        }
        this.last = this.finalData
      } else {
        this.finalData = this.$refs['jsonTree' + 0][0].jsonJ
        this.last = this.finalData
      }
      this.$emit('input', this.last)
      this.showEdit = false
    },
    addNewJson () {
      this.originData.push({})
    }
  }
}
</script>

<style lang="scss">
.taskman-cmdb-json-config {
  width: 100%;
  .inline {
    display: flex;
    align-items: center;
    height: 34px;
    .text {
      font-size: 13px;
      color:#515a6e;
      display: block;
      max-width: 380px;
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
.taskman-json-viewer {
  .jv-code {
    overflow: hidden;
    padding: 0px 10px;
  }
}
</style>
