<template>
  <div class=" ">
    <Tabs @on-click="changeNode" :value="currentNode">
      <template v-for="node in nodes">
        <TabPane :label="node.nodeName" :name="node.nodeId" :key="node.nodeId">
          <FormComponent :ref="node.nodeId" :currentNode="currentNode" :node="node"></FormComponent>
        </TabPane>
      </template>
    </Tabs>
  </div>
</template>

<script>
import FormComponent from './task-form-component'
import { getTemplateNodes } from '@/api/server.js'
export default {
  name: '',
  data () {
    return {
      requestTemplateId: '614043ac9379fb1e',
      currentNode: '',
      nodes: []
    }
  },
  mounted () {
    this.getTemplateNodes()
  },
  methods: {
    async getTemplateNodes () {
      const { statusCode, data } = await getTemplateNodes(this.requestTemplateId)
      if (statusCode === 'OK') {
        this.nodes = data
        this.currentNode = data[0].nodeId
        // this.initTab(this.currentNode, data[0])
        // this.$refs[this.currentNode].initData(data[0])
      }
    },
    changeNode (nodeId) {
      this.currentNode = nodeId
      const find = this.nodes.find(n => n.nodeId === this.currentNode)
      this.$refs[this.currentNode][0].initData(this.currentNode, find)
    },
    initTab (currentNode, data) {
      this.$refs[this.currentNode][0].initData(data)
    }
  },
  components: {
    FormComponent
  }
}
</script>

<style scoped lang="scss"></style>
