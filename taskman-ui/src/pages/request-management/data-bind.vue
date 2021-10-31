<template>
  <div style="width:50%;margin: 0 auto;min-width: 700px">
    <Tabs>
      <template v-for="node in nodes">
        <TabPane :label="node.nodeName" :name="node.nodeId" :key="node.nodeId">
          <ul>
            <CheckboxGroup v-model="node.bindData">
              <li
                v-for="item in filterBindData(node)"
                :key="item.id"
                style="width: 46%;display: inline-block;margin: 6px"
              >
                <Checkbox :label="item.id" :disabled="$parent.formDisable">
                  <Tooltip
                    :content="item.displayName"
                    :delay="500"
                    max-width="300"
                    :disabled="item.displayName.length < 30"
                  >
                    <div class="text-ellipsis">{{ item.displayName }}</div>
                  </Tooltip>
                </Checkbox>
              </li>
            </CheckboxGroup>
          </ul>
        </TabPane>
      </template>
    </Tabs>
    <div style="text-align: center;margin-top:48px">
      <Button @click="saveRequest" :disabled="$parent.formDisable" type="primary">{{ $t('temporary_storage') }}</Button>
      <Button @click="rollbackRequest" type="error" :disabled="$parent.formDisable" v-if="$parent.isHandle">{{
        $t('go_back')
      }}</Button>
      <Button @click="startRequest" :disabled="$parent.formDisable" v-if="$parent.isHandle">{{
        $t('final_version')
      }}</Button>
    </div>
  </div>
</template>

<script>
import {
  getTemplateNodes,
  getBindData,
  saveRequest,
  getBindRelate,
  updateRequestStatus,
  startRequest
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
      }
    }
  },
  mounted () {
    this.finalData.procDefId = this.$parent.procDefId
    this.finalData.procDefKey = this.$parent.procDefKey
    this.getTemplateNodes()
    this.getBindData()
  },
  methods: {
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
            node.bindData = find.boundEntityValues.map(b => b.oid)
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
      return this.bindData
    },
    async saveRequest () {
      let tmpData = []
      this.nodes.forEach(n => {
        let params = {
          boundEntityValues: [],
          nodeDefId: n.nodeDefId,
          nodeId: n.nodeId
        }
        n.bindData.map(b => {
          let findData = this.bindData.find(bData => bData.id === b)
          findData.bindInfo.bindFlag = 'Y'
          params.boundEntityValues.push(findData.bindInfo)
        })
        tmpData.push(params)
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
    async getTemplateNodes () {
      const { statusCode, data } = await getTemplateNodes(this.$parent.requestTemplate)
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
