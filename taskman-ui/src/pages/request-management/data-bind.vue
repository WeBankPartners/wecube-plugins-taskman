<template>
  <div style="width:50%;margin: 0 auto;">
    <Form :label-width="200">
      <template v-for="node in nodes">
        <FormItem :label="node.nodeName" :key="node.nodeId">
          <Select v-model="node.bindData" multiple filterable :disabled="$parent.formDisable" style="width:100%">
            <Option v-for="item in filterBindData(node)" :value="item.id" :key="item.id">{{ item.displayName }}</Option>
          </Select>
        </FormItem>
      </template>
    </Form>
    <div style="text-align: center;margin-top:48px">
      <Button @click="saveRequest" :disabled="$parent.formDisable" type="primary">{{ $t('save') }}</Button>
      <Button @click="commitRequest" v-if="$parent.isHandle" :disabled="$parent.formDisable">{{ $t('提交') }}</Button>
    </div>
  </div>
</template>

<script>
import { getTemplateNodes, getBindData, saveRequest, getBindRelate, updateRequestStatus } from '@/api/server.js'
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
    async commitRequest () {
      await this.saveRequest()
      const { statusCode } = await updateRequestStatus(this.$parent.requestId, 'Pending')
      if (statusCode === 'OK') {
        this.$router.push({ path: '/request' })
      }
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
