<template>
  <div style="width: 400px;">
    <div style="margin-bottom: 8px;">
      <Alert type="warning" show-icon
        >如对以下编排的设计、执行情况有任何疑问。请复制编排名。发送给wecube管理员询问</Alert
      >
      <!-- <span class="custom-title">{{ $t('workflow_name') }}</span> -->
      <span class="custom-display">{{ flowData.procDefName }}</span>
    </div>
    <div id="graphcontain">
      <div class="graph-container" id="flow" style="height: 96%"></div>
      <Button class="reset-button" size="small" @click="ResetFlow">ResetZoom</Button>
    </div>
  </div>
</template>
<script>
import { getFlowByTemplateId } from '@/api/server'
import * as d3 from 'd3-selection'
// eslint-disable-next-line no-unused-vars
import * as d3Graphviz from 'd3-graphviz'
export default {
  props: {
    requestTemplate: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      flowGraph: {},
      flowData: {}
    }
  },
  watch: {
    requestTemplate: {
      handler (val) {
        if (val) {
          this.orchestrationSelectHandler()
        }
      },
      immediate: true
    }
  },
  methods: {
    orchestrationSelectHandler () {
      this.currentFlowNodeId = ''
      this.currentModelNodeRefs = []
      this.getFlowOutlineData()
    },
    async getFlowOutlineData () {
      let { statusCode, data } = await getFlowByTemplateId(this.requestTemplate)
      if (statusCode === 'OK') {
        this.flowData = data
        this.initFlowGraph()
      }
    },
    ResetFlow () {
      if (this.flowGraph.graphviz) {
        this.flowGraph.graphviz.resetZoom()
      }
    },
    renderFlowGraph (excution) {
      const statusColor = {
        Completed: '#5DB400',
        deployed: '#7F8A96',
        InProgress: '#3C83F8',
        Faulted: '#FF6262',
        Risky: '#BF22E0',
        Timeouted: '#F7B500',
        NotStarted: '#7F8A96'
      }
      let nodes =
        this.flowData &&
        this.flowData.flowNodes &&
        this.flowData.flowNodes
          .filter(i => i.status !== 'predeploy')
          .map((_, index) => {
            if (_.nodeType === 'startEvent' || _.nodeType === 'endEvent') {
              const defaultLabel = _.nodeType === 'startEvent' ? 'start' : 'end'
              return `${_.nodeId} [label="${_.nodeName || defaultLabel}", fontsize="10", class="flow",style="${
                excution ? 'filled' : 'none'
              }" color="${excution ? statusColor[_.status] : '#7F8A96'}" shape="circle", id="${_.nodeId}"]`
            } else {
              // const className = _.status === 'Faulted' || _.status === 'Timeouted' ? 'retry' : 'normal'
              const className = 'retry'
              const isModelClick = this.currentModelNodeRefs.indexOf(_.orderedNo) > -1
              return `${_.nodeId} [fixedsize=false label="${(_.orderedNo ? _.orderedNo + ' ' : '') +
                _.nodeName}" class="flow ${className}" style="${excution || isModelClick ? 'filled' : 'none'}" color="${
                excution
                  ? statusColor[_.status]
                  : isModelClick
                    ? '#ff9900'
                    : _.nodeId === this.currentFlowNodeId
                      ? '#5DB400'
                      : '#7F8A96'
              }"  shape="box" id="${_.nodeId}" ]`
            }
          })
      let genEdge = () => {
        let pathAry = []
        this.flowData &&
          this.flowData.flowNodes &&
          this.flowData.flowNodes.forEach(_ => {
            if (_.succeedingNodeIds.length > 0) {
              let current = []
              current = _.succeedingNodeIds.map(to => {
                return (
                  '"' +
                  _.nodeId +
                  '"' +
                  ' -> ' +
                  `${'"' + to + '"'} [color="${excution ? statusColor[_.status] : 'black'}"]`
                )
              })
              pathAry.push(current)
            }
          })
        return pathAry
          .flat()
          .toString()
          .replace(/,/g, ';')
      }
      let nodesToString = Array.isArray(nodes) ? nodes.toString().replace(/,/g, ';') + ';' : ''
      let nodesString =
        'digraph G {' +
        'bgcolor="transparent";' +
        'Node [fontname=Arial, height=".3", fontsize=12];' +
        'Edge [fontname=Arial, color="#7f8fa6", fontsize=10];' +
        nodesToString +
        genEdge() +
        '}'

      this.flowGraph.graphviz.transition().renderDot(nodesString)
    },
    initFlowGraph (excution = false) {
      const graphEl = document.getElementById('flow')
      let graph
      graph = d3.select(`#flow`)
      graph.on('dblclick.zoom', null)
      this.flowGraph.graphviz = graph
        .graphviz()
        .fit(true)
        .zoom(true)
        .height(graphEl.offsetHeight - 10)
        .width(graphEl.offsetWidth - 10)
      this.renderFlowGraph(excution)
    }
  }
}
</script>
<style lang="scss" scoped>
.custom-title {
  text-align: right;
  vertical-align: middle;
  float: left;
  font-size: 14px;
  color: #515a6e;
  line-height: 1;
  padding: 0 12px 10px 0;
  box-sizing: border-box;
}

.custom-display {
  display: inline-block;
  width: 400px;
  height: 32px;
  line-height: 1.5;
  padding: 4px 7px;
  font-size: 14px;
  border: 1px solid #069cec;
  border-radius: 4px;
  color: #069cec;
  cursor: pointer;
}
#graphcontain {
  position: relative;
  border: 1px solid #d3cece;
  border-radius: 3px;
  padding: 5px;
  height: calc(100vh - 220px);
}
.reset-button {
  position: absolute;
  right: 20px;
  bottom: 20px;
  font-size: 12px;
}
</style>
