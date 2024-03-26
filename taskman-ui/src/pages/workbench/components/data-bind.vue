<template>
  <div class="workbench-data-bind">
    <Tabs v-if="nodes.length" class="tabs">
      <template v-for="node in nodes">
        <TabPane :label="node.nodeName" :name="node.nodeId" :key="node.nodeId">
          <ul v-if="node.classification && node.classification.length > 0">
            <CheckboxGroup v-model="node.bindData">
              <template v-for="(entity, entityIndex) in node.classification">
                <Divider orientation="left" :key="entityIndex">{{ entity }}</Divider>
                <template v-for="item in filterBindData(node)">
                  <li
                    v-if="item.entityName === entity"
                    :key="item.id + entityIndex"
                    style="width: 46%;display: inline-block;margin: 6px"
                  >
                    <Checkbox :label="item.id" :disabled="formDisable">
                      <Tooltip
                        :content="item.displayName"
                        :delay="500"
                        :disabled="item.displayName.length < 30"
                        max-width="300"
                      >
                        <div class="text-ellipsis">
                          <span v-if="item.id.startsWith('tmp__')" style="color: #338cf0">(new)</span>
                          {{ item.displayName }}
                        </div>
                      </Tooltip>
                    </Checkbox>
                  </li>
                </template>
              </template>
            </CheckboxGroup>
          </ul>
          <div v-else>{{ $t('tw_no_data') }}</div>
        </TabPane>
      </template>
    </Tabs>
    <div v-else class="no-data">{{ $t('tw_unbind_workflow') }}</div>
    <div v-if="showBtn" style="text-align: center;margin-top:12px">
      <!--暂存-->
      <Button @click="saveRequest('save')" :disabled="formDisable">{{ $t('tw_save_draft') }}</Button>
      <!--确认定版-->
      <Button @click="startRequest" :disabled="formDisable" v-if="isHandle" type="primary">{{
        $t('tw_confirm')
      }}</Button>
      <!--回退-->
      <Button @click="rollbackRequest" type="error" :disabled="formDisable" v-if="isHandle">{{ $t('go_back') }}</Button>
      <!-- <Button @click="checkHistory" v-if="requestHistory" type="success" style="margin-left:30px">{{
        $t('pr-vision')
      }}</Button> -->
    </div>
  </div>
</template>

<script>
import {
  getTemplateNodesForRequest,
  getBindData,
  saveRequestNew,
  getBindRelate,
  updateRequestStatus,
  startRequestNew,
  requestParent
} from '@/api/server.js'
import { deepClone } from '@/pages/util/index'
export default {
  name: '',
  data () {
    return {
      nodes: [],
      bindData: [],
      taskNodeBindInfos: [],
      finalData: {
        procDefId: '',
        procDefKey: '',
        rootEntityValue: {
          oid: ''
        },
        taskNodeBindInfos: ''
      },
      requestHistory: false,
      backReason: this.$t('tw_recall_reason'), // 回退说明
      unCheckData: [], // 没有选中的数据
      columns: [
        {
          title: this.$t('tw_tag'),
          key: 'nodeName'
        },
        {
          title: this.$t('tw_data_type'),
          key: 'entityName'
        },
        {
          title: this.$t('tw_data'),
          key: 'entityDisplayName'
        }
      ]
    }
  },
  props: {
    requestId: {
      type: String,
      default: ''
    },
    formDisable: {
      type: Boolean,
      default: false
    },
    isHandle: {
      type: Boolean,
      default: false
    },
    requestTemplate: {
      type: String,
      default: ''
    },
    showBtn: {
      type: Boolean,
      default: true
    },
    actionName: {
      type: String,
      default: '1'
    }
  },
  watch: {
    requestId: {
      handler (val) {
        if (val) {
          this.finalData.procDefId = this.$parent.procDefId
          this.finalData.procDefKey = this.$parent.procDefKey
          this.getTemplateNodesForRequest()
          this.getBindData()
          this.requestParent()
        }
      },
      deep: true,
      immediate: true
    }
  },
  mounted () {},
  methods: {
    async requestParent () {
      const { data } = await requestParent(this.requestId)
      if (data) {
        this.requestHistory = data
      } else {
        this.requestHistory = false
      }
    },
    checkHistory () {
      this.$router.push({
        path: '/requestCheck',
        query: {
          requestId: this.requestHistory,
          requestTemplate: null,
          isAdd: 'N',
          isCheck: 'Y',
          isHandle: 'N',
          jumpFrom: 'group_initiated'
        }
      })
    },
    async startRequest () {
      // 获取没有勾选的数据
      this.unCheckData = []
      this.copyNodes = deepClone(this.nodes)
      this.copyNodes.forEach(n => {
        if (n.bindOptions.length > 0) {
          n.bindOptions = n.bindOptions.map(bData => {
            bData.bindFlag = 'N'
            return bData
          })
          n.bindData.forEach(b => {
            let findData = n.bindOptions.find(bData => bData.oid === b)
            findData.bindFlag = 'Y'
          })
        }
      })
      this.copyNodes.forEach(i => {
        i.bindOptions.forEach(j => {
          if (j.bindFlag === 'N') {
            this.unCheckData.push({
              nodeName: i.nodeName,
              entityName: j.entityName,
              entityDisplayName: j.entityDisplayName
            })
          }
        })
      })
      this.$Modal.confirm({
        title: this.$t('tw_confirm') + this.$t('final_version'),
        'z-index': 1000000,
        width: 500,
        loading: true,
        render: () => {
          if (this.unCheckData.length) {
            return (
              <div style="width:450px;">
                <span>{this.$t('tw_delete_bind_tips')}</span>
                <Table style="margin:20px 0" size="small" columns={this.columns} data={this.unCheckData}></Table>
              </div>
            )
          }
        },
        onOk: async () => {
          this.$Modal.remove()
          const flag = await this.saveRequest('submit')
          if (flag) {
            const { statusCode } = await startRequestNew(this.requestId, this.finalData)
            if (statusCode === 'OK') {
              this.$Notice.success({
                title: this.$t('successful'),
                desc: this.$t('successful')
              })
              this.$router.push({
                path: `/taskman/workbench?tabName=hasProcessed&actionName=${this.actionName}&type=1`
              })
            }
          }
        },
        onCancel: () => {}
      })
    },
    async getBindRelate () {
      const { statusCode, data } = await getBindRelate(this.requestId)
      if (statusCode === 'OK') {
        if (data.taskNodeBindInfos.length > 0) {
          this.nodes.forEach(node => {
            const find = data.taskNodeBindInfos.find(d => d.nodeId === node.nodeId)
            if (!find) {
              node.bindData = []
              node.bindOptions = []
              this.$set(node, 'classification', [])
            } else {
              node.bindData = find.boundEntityValues.filter(bound => bound.bindFlag === 'Y').map(b => b.oid)
              node.bindOptions = find.boundEntityValues || []
              let classification = new Set()
              node.bindOptions.forEach(d => {
                classification.add(d.entityName)
              })
              // node.classification = Array.from(classification)
              this.$set(node, 'classification', Array.from(classification))
            }
          })
        }
      }
    },
    // 定版退回操作
    async rollbackRequest () {
      this.$Modal.confirm({
        title: this.$t('tw_confirm') + this.$t('go_back'),
        'z-index': 1000000,
        loading: false,
        render: () => {
          return (
            <div>
              <div style="font-size:12px;margin-bottom:10px;">{this.$t('tw_rollback_tips')}</div>
              <Input
                type="textarea"
                maxlength={255}
                show-word-limit
                v-model={this.backReason}
                placeholder={this.$t('tw_back_bind_placeholder')}
              ></Input>
            </div>
          )
        },
        onOk: async () => {
          if (!this.backReason.trim()) {
            this.$Notice.warning({
              title: this.$t('warning'),
              desc: this.$t('tw_back_bind_tips')
            })
          } else {
            // this.$Modal.remove()
            const flag = await this.saveRequest('submit')
            if (flag) {
              const params = {
                description: this.backReason
              }
              const { statusCode } = await updateRequestStatus(this.requestId, 'Draft', params)
              if (statusCode === 'OK') {
                this.$Notice.success({
                  title: this.$t('successful'),
                  desc: this.$t('successful')
                })
                this.$router.push({
                  path: `/taskman/workbench?tabName=hasProcessed&actionName=${this.actionName}&type=1`
                })
              }
            }
          }
        },
        onCancel: () => {}
      })
    },
    filterBindData (node) {
      if (!node.bindOptions || node.bindOptions.length === 0) return []
      const oid = node.bindOptions.map(n => n.oid)
      const res = this.bindData.filter(bData => oid.includes(bData.id))
      return res
    },
    async saveRequest (type) {
      let tmpData = []
      this.nodes.forEach(n => {
        let params = {
          boundEntityValues: [],
          nodeDefId: n.nodeDefId,
          nodeId: n.nodeId
        }
        if (n.bindOptions.length > 0) {
          n.bindOptions = n.bindOptions.map(bData => {
            bData.bindFlag = 'N'
            return bData
          })
          n.bindData.forEach(b => {
            let findData = n.bindOptions.find(bData => bData.oid === b)
            findData.bindFlag = 'Y'
          })
          params.boundEntityValues = n.bindOptions
          tmpData.push(params)
        }
      })
      this.finalData.taskNodeBindInfos = tmpData
      // type为save(暂存)、submit(撤回、确认定版)
      const { statusCode } = await saveRequestNew(this.requestId, type, this.finalData)
      if (statusCode === 'OK') {
        type === 'save' &&
          this.$Notice.success({
            title: this.$t('successful'),
            desc: this.$t('successful')
          })
        return true
      } else {
        return false
      }
    },
    async getTemplateNodesForRequest () {
      const { statusCode, data } = await getTemplateNodesForRequest(this.requestTemplate)
      if (statusCode === 'OK') {
        this.nodes = data
          .filter(d => d.taskCategory !== '')
          .map(n => {
            n.bindData = []
            return n
          })
        this.getBindRelate()
      }
    },
    async getBindData () {
      const { statusCode, data } = await getBindData(this.requestId)
      if (statusCode === 'OK') {
        this.finalData.rootEntityValue.oid = data.rootEntityId
        data.data.forEach(ele => {
          let displayKeyName = []
          let attrValues = []
          ele.title.forEach(t => {
            attrValues.push({
              attrDefId: t.attrDefId,
              attrName: t.attrDefName,
              dataType: t.attrDefDataType,
              dataValue: ''
            })
            if (t.inDisplayName === 'yes') {
              displayKeyName.push(t.name)
            }
          })
          ele.value.forEach(v => {
            let tmpAttrValue = JSON.parse(JSON.stringify(attrValues))
            tmpAttrValue.forEach(attr => {
              attr.dataValue = v.entityData[attr.attrName]
            })
            v.bindInfo = {
              attrValues: tmpAttrValue,
              bindFlag: 'N',
              entityDataId: v.dataId,
              entityDataOp: v.entityDataOp,
              entityDataState: '',
              entityDefId: '',
              entityName: v.entityName,
              fullEntityDataId: v.fullDataId,
              oid: v.id,
              packageName: v.packageName
            }
            this.bindData.push(v)
          })
        })
      }
    }
  },
  components: {}
}
</script>

<style scoped lang="scss">
.text-ellipsis {
  width: 450px;
  white-space: nowrap;
  overflow: hidden;
  display: inline-block;
  text-overflow: ellipsis;
  vertical-align: bottom;
}
.no-data {
  height: 60px;
  line-height: 30px;
  color: #515a6e;
}
</style>
<style lang="scss">
.workbench-data-bind {
  width: 100%;
  .ivu-checkbox-group-item {
    width: 100%;
  }
  .ivu-tooltip {
    width: 100% !important;
  }
}
</style>
