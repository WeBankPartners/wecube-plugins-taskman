<template>
  <div style="width:50%;margin: 0 auto;">
    <Form :label-width="200">
      <template v-for="node in nodes">
        <FormItem :label="node.nodeName" :key="node.nodeId">
          <Select v-model="node.bindData" multiple filterable style="width:100%">
            <Option v-for="item in filterBindData(node)" :value="item.id" :key="item.id">{{ item.displayName }}</Option>
          </Select>
        </FormItem>
      </template>
    </Form>
    <Button @click="request" style="float:right" type="primary">{{ $t('initiate_request') }}</Button>
  </div>
</template>

<script>
import { getTemplateNodes, getBindData, startRequest } from '@/api/server.js'
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
    filterBindData (node) {
      return this.bindData
    },
    request () {
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
      this.startRequest()
    },
    async startRequest () {
      const { statusCode } = await startRequest(this.$parent.requestId, this.finalData)
      if (statusCode === 'OK') {
        this.$router.push('/')
      }
    },
    async getTemplateNodes () {
      const { statusCode, data } = await getTemplateNodes(this.$parent.requestTemplate)
      // const { statusCode, data } = await getTemplateNodes('6166464841f44a4c')
      if (statusCode === 'OK') {
        this.nodes = data
          .filter(d => d.taskCategory !== '')
          .map(n => {
            n.bindData = []
            return n
          })
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
            if (displayKeyName.length > 0) {
              displayKeyName.forEach(k => {
                v.displayName += v.entityData[k]
              })
              attrValues.forEach(attr => {
                attr.dataValue = v.entityData[attr.attrName]
              })
              v.bindInfo = {
                attrValues: attrValues,
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
            } else {
              v.bindInfo = {
                attrValues: v.displayName,
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

<style scoped lang="scss"></style>
