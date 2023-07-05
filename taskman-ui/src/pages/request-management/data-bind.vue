<template>
  <div style="width:50%;margin: 0 auto;min-width: 700px">
    <Tabs>
      <template v-for="node in nodes">
        <TabPane :label="node.nodeName" :name="node.nodeId" :key="node.nodeId">
          <ul>
            <CheckboxGroup v-model="node.bindData">
              <template v-for="(entity, entityIndex) in node.classification">
                <Divider orientation="left" :key="entityIndex">{{ entity }}</Divider>
                <template v-for="item in filterBindData(node)">
                  <li
                    v-if="item.entityName === entity"
                    :key="item.id + entityIndex"
                    style="width: 46%;display: inline-block;margin: 6px"
                  >
                    <Checkbox :label="item.id" :disabled="$parent.formDisable">
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
        </TabPane>
      </template>
    </Tabs>
    <div style="text-align: center;margin-top:12px">
      <Button @click="saveRequest" :disabled="$parent.formDisable" type="primary">{{ $t('temporary_storage') }}</Button>
      <Button @click="rollbackRequest" type="error" :disabled="$parent.formDisable" v-if="$parent.isHandle">{{
        $t('go_back')
      }}</Button>
      <Button @click="startRequest" :disabled="$parent.formDisable" v-if="$parent.isHandle">{{
        $t('final_version')
      }}</Button>
      <Button @click="checkHistory" v-if="requestHistory" type="success" style="margin-left:30px">{{
        $t('pr-vision')
      }}</Button>
    </div>
  </div>
</template>

<script>
import {
  getTemplateNodesForRequest,
  getBindData,
  saveRequest,
  getBindRelate,
  updateRequestStatus,
  startRequest,
  requestParent
} from '@/api/server.js'
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
      requestHistory: false
    }
  },
  mounted () {
    this.finalData.procDefId = this.$parent.procDefId
    this.finalData.procDefKey = this.$parent.procDefKey
    this.getTemplateNodesForRequest()
    this.getBindData()
    this.requestParent()
  },
  methods: {
    async requestParent () {
      const { data } = await requestParent(this.$parent.requestId)
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
      this.$Modal.confirm({
        title: this.$t('confirm') + this.$t('final_version'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          await this.saveRequest()
          const { statusCode } = await startRequest(this.$parent.requestId, this.finalData)
          if (statusCode === 'OK') {
            this.$Notice.success({
              title: this.$t('successful'),
              desc: this.$t('successful')
            })
            this.$router.push({ path: '/taskman/request-mgmt' })
          }
        },
        onCancel: () => {}
      })
    },
    async getBindRelate () {
      const { statusCode, data } = await getBindRelate(this.$parent.requestId)
      if (statusCode === 'OK') {
        if (data.taskNodeBindInfos.length > 0) {
          this.nodes.forEach(node => {
            const find = data.taskNodeBindInfos.find(d => d.nodeId === node.nodeId)
            if (!find) {
              node.bindData = []
              node.bindOptions = []
              node.classification = []
            } else {
              node.bindData = find.boundEntityValues.filter(bound => bound.bindFlag === 'Y').map(b => b.oid)
              node.bindOptions = find.boundEntityValues || []
              let classification = new Set()
              node.bindOptions.forEach(d => {
                classification.add(d.entityName)
              })
              node.classification = Array.from(classification)
            }
          })
        }
      }
    },
    async rollbackRequest () {
      this.$Modal.confirm({
        title: this.$t('confirm') + this.$t('go_back'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          await this.saveRequest()
          const { statusCode } = await updateRequestStatus(this.$parent.requestId, 'Draft')
          if (statusCode === 'OK') {
            this.$router.push({ path: '/taskman/request-mgmt' })
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
    async saveRequest () {
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
      const { statusCode } = await saveRequest(this.$parent.requestId, this.finalData)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
      }
    },
    async getTemplateNodesForRequest () {
      const { statusCode, data } = await getTemplateNodesForRequest(this.$parent.requestTemplate)
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
      const { statusCode, data } = await getBindData(this.$parent.requestId)
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
.ivu-checkbox-group-item {
  width: 100%;
}
.ivu-tooltip {
  width: 90%;
}
.text-ellipsis {
  width: 90%;
  white-space: nowrap;
  overflow: hidden;
  display: inline-block;
  text-overflow: ellipsis;
  vertical-align: bottom;
}
</style>
