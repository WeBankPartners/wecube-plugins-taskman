<template>
  <div class="">
    <div style="text-align:end">
      <Button @click="confirmTemplate" type="primary">模板发布</Button>
    </div>
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
import { getTemplateNodes, confirmTemplate } from '@/api/server.js'
export default {
  name: '',
  data () {
    return {
      requestTemplateId: '614479d70dd9be04',
      currentNode: '',
      nodes: []
    }
  },
  mounted () {
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
    },
    publishTemplate () {
      console.log('发布模板！')
    }
  },
  components: {
    FormComponent
  }
}
</script>

<style scoped lang="scss"></style>
