<template>
  <div class="">
    {{ currentNode }}
    <Tabs @on-click="changeNode" :value="currentNode">
      <template v-for="node in nodes">
        <TabPane :label="node.nodeName" :name="node.nodeId" :key="node.nodeId">
          <FormComponent
            :ref="node.nodeId"
            :currentNode="currentNode"
            :node="node"
            :requestTemplateId="requestTemplateId"
          ></FormComponent>
        </TabPane>
      </template>
    </Tabs>
  </div>
</template>

<script>
import FormComponent from './task-form-component'
import { getTemplateNodes, confirmTemplate } from '@/api/server.js'
export default {
  name: '',
  data () {
    return {
      requestTemplateId: '',
      currentNode: '',
      nodes: []
    }
  },
  mounted () {
    this.requestTemplateId = this.$parent.requestTemplateId
    this.getTemplateNodes()
  },
  methods: {
    async confirmTemplate () {
      const { statusCode } = await confirmTemplate(this.requestTemplateId)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
      }
    },
    async getTemplateNodes () {
      const { statusCode, data } = await getTemplateNodes(this.requestTemplateId)
      if (statusCode === 'OK') {
        this.nodes = data.filter(item => item.taskCategory === 'SUTN')
        // this.nodes = data
        this.$nextTick(() => {
          this.currentNode = this.nodes[0].nodeId
          this.initTab(this.currentNode, this.nodes[0])
        })
        // this.$refs[this.currentNode].initData(data[0])
      }
    },
    changeNode (nodeId) {
      this.currentNode = nodeId
      const find = this.nodes.find(n => n.nodeId === this.currentNode)
      console.log(this.$refs[this.currentNode][0].initData)
      this.$refs[this.currentNode][0].initData(this.currentNode, find, this.requestTemplateId)
    },
    initTab (currentNode, data) {
      this.$refs[this.currentNode][0].initData(currentNode, data, this.requestTemplateId)
    }
  },
  components: {
    FormComponent
  }
}
</script>

<style scoped lang="scss"></style>
