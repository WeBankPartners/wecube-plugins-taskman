<template>
  <div>
    <Drawer
      title="操作实例勾选"
      v-model="visible"
      width="720"
      :mask-closable="true"
      :styles="styles"
      @on-close="handleCancel"
      class="workbench-choose-example"
    >
      <div class="content" :style="{ maxHeight: maxHeight + 'px' }">
        <Tree :data="treeData" show-checkbox @on-check-change="selectTreeData"></Tree>
      </div>
      <div class="drawer-footer">
        <Button style="margin-right: 8px" @click="handleCancel">取消</Button>
        <Button type="primary" class="primary" @click="handleSubmit">确定</Button>
      </div>
    </Drawer>
  </div>
</template>

<script>
import { debounce } from '@/pages/util'
import { getEntityData } from '@/api/server'
export default {
  props: {
    visible: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      requestData: [],
      styles: {
        height: 'calc(100% - 55px)',
        overflow: 'auto',
        paddingBottom: '53px',
        position: 'static'
      },
      numberStyle: {
        display: 'inline-block',
        background: '#2d8cf0',
        color: '#fff',
        padding: '0px 5px',
        borderRadius: '2px',
        margin: '0 2px',
        fontSize: '12px'
      },
      maxHeight: 500,
      treeData: [
        {
          title: '子系统',
          expand: false,
          parent: 1,
          selectCount: 0,
          render: h => {
            return (
              <div>
                已选<span style={this.numberStyle}>{this.treeData[0].selectCount}</span>个子系统
              </div>
            )
          },
          children: [
            { title: '198.168.15.30', value: '111', type: 1 },
            { title: '198.168.15.30', value: '111', type: 1 },
            { title: '198.168.15.30', value: '111', type: 1 }
          ]
        },
        {
          title: '单元',
          expand: false,
          parent: 2,
          selectCount: 0,
          render: h => {
            return (
              <div>
                已选<span style={this.numberStyle}>{this.treeData[1].selectCount}</span>个单元
              </div>
            )
          },
          children: [
            { title: '198.168.15.30', value: '111', type: 2 },
            { title: '198.168.15.30', value: '111', type: 2 },
            { title: '198.168.15.30', value: '111', type: 2 }
          ]
        },
        {
          title: '实例',
          expand: true,
          parent: 3,
          selectCount: 0,
          render: h => {
            return (
              <div>
                已选<span style={this.numberStyle}>{this.treeData[2].selectCount}</span>个实例
              </div>
            )
          },
          children: [
            { title: '198.168.15.30', value: '111', type: 3, checked: true },
            { title: '198.168.15.30', value: '111', type: 3, checked: true },
            { title: '198.168.15.30', value: '111', type: 3, checked: true }
          ]
        }
      ]
    }
  },
  mounted () {
    this.maxHeight = document.body.clientHeight - 150
    window.addEventListener(
      'resize',
      debounce(() => {
        this.maxHeight = document.body.clientHeight - 150
      }, 100)
    )
    this.getEntityData()
  },
  methods: {
    async getEntityData () {
      let params = {
        params: {
          requestId: this.$parent.requestId || '6557454d5324718d',
          rootEntityId: this.$parent.form.rootEntityId
        }
      }
      const { statusCode, data } = await getEntityData(params)
      if (statusCode === 'OK') {
        // this.activeTab = this.activeTab || data.data[0].entity
        this.requestData = data.data
        // this.$nextTick(() => {
        //   const index = this.requestData.findIndex(r => r.entity === this.activeTab || r.itemGroup === this.activeTab)
        //   this.initTable(index)
        // })
      }
    },
    selectTreeData (arr) {
      this.treeData.forEach(i => {
        i.selectCount = 0
        for (let j of arr) {
          if (i.parent === j.type) {
            i.selectCount++
          }
        }
      })
    },
    handleSubmit () {
      this.$emit('update:visible', false)
      this.$emit('getData', this.treeData)
    },
    handleCancel () {
      this.$emit('update:visible', false)
    }
  }
}
</script>

<style lang="scss" scoped>
.content {
  min-height: 500px;
  padding: 20px;
  border: 2px dashed #e8eaec;
  overflow-y: auto;
}
.drawer-footer {
  width: 100%;
  position: absolute;
  bottom: 0;
  left: 0;
  border-top: 1px solid #e8e8e8;
  padding: 10px 16px;
  text-align: center;
  background: #fff;
}
</style>
<style lang="scss">
.workbench-choose-example {
  .ivu-tree-title-selected,
  .ivu-tree-title-selected:hover {
    background-color: transparent;
  }
}
</style>
