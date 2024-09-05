<template>
  <div style="width: 400px;">
    <div style="margin-bottom: 8px;">
      <Alert type="warning" show-icon>{{ $t('tw_flow_tips') }}</Alert>
      <!-- <span class="custom-title">{{ $t('workflow_name') }}</span> -->
      <span class="custom-display" @click="jumpToFlowDetail">
        {{ flowData.procDefName }}
        <span v-if="flowData.procDefVersion">{{ `【${flowData.procDefVersion}】` }}</span>
      </span>
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
    },
    flowId: {
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
    // 打开编排模板详情页
    jumpToFlowDetail() {
      if (process.env.VUE_APP_PLUGIN === 'plugin') {
        window.sessionStorage.currentPath = '' // 先清空session缓存页面，不然打开新标签页面会回退到缓存的页面
        const path = `${window.location.origin}/#/collaboration/workflow-mgmt?flowId=${this.flowId}&editFlow=false&flowListTab=deployed`
        window.open(path, '_blank')
      }
    },
    orchestrationSelectHandler () {
      this.currentFlowNodeId = ''
      this.currentModelNodeRefs = []
      this.getFlowOutlineData()
    },
    async getFlowOutlineData () {
      let { statusCode, data } = await getFlowByTemplateId(this.requestTemplate)
      if (statusCode === 'OK') {
        this.flowData = data || {}
        this.initFlowGraph()
      }
    },
    ResetFlow () {
      if (this.flowGraph.graphviz) {
        this.flowGraph.graphviz.resetZoom()
      }
    },
    renderFlowGraph (excution) {
      // 节点颜色
      const statusColor = {
        Completed: '#5DB400',
        deployed: '#7F8A96',
        InProgress: '#3C83F8',
        Faulted: '#FF6262',
        Risky: '#BF22E0',
        Timeouted: '#F7B500',
        NotStarted: '#7F8A96',
        wait: '#7F8A96'
      }
      const nodes = this.flowData
        && this.flowData.flowNodes
        && this.flowData.flowNodes
          .filter(i => i.status !== 'predeploy')
          .map(_ => {
            const shapeMap = {
              start: 'circle', // 开始
              end: 'doublecircle', // 结束
              abnormal: 'doublecircle', // 异常
              decision: 'diamond', // 判断开始
              decisionMerge: 'diamond', // 判断结束
              fork: 'Mdiamond', // 并行开始
              merge: 'Mdiamond', // 并行结束
              human: 'tab', // 人工
              automatic: 'rect', // 自动
              data: 'cylinder', // 数据
              subProc: 'component', // 子编排
              date: 'Mcircle', // 固定时间
              timeInterval: 'Mcircle' // 时间间隔
            }
            if (['start', 'end', 'abnormal', 'date', 'timeInterval'].includes(_.nodeType)) {
              const className = 'retry'
              const defaultLabel = _.nodeType
              return `${_.nodeId} [label="${
                _.nodeName || defaultLabel
              }", width="0.8", class="flow ${className}", fixedsize=true, style="${
                excution ? 'filled' : 'none'
              }" fillcolor="${excution ? statusColor[_.status] || '#7F8A96' : '#7F8A96'}" shape="${
                shapeMap[_.nodeType]
              }", id="${_.nodeId}"]`
            }
            // const className = _.status === 'Faulted' || _.status === 'Timeouted' ? 'retry' : 'normal'
            let className = 'retry'
            if (['decision'].includes(_.nodeType) && _.status === 'Faulted') {
              className = ''
            }
            const isModelClick = this.currentModelNodeRefs.indexOf(_.orderedNo) > -1
            return `${_.nodeId} [fixedsize=false label="${
              (_.orderedNo ? _.orderedNo + ' ' : '') + _.nodeName
            }" class="flow ${className}" style="${excution || isModelClick ? 'filled' : 'none'}" fillcolor="${
              excution
                ? statusColor[_.status] || '#7F8A96'
                : isModelClick
                  ? '#ff9900'
                  : _.nodeId === this.currentFlowNodeId
                    ? '#5DB400'
                    : '#7F8A96'
            }"  shape="${shapeMap[_.nodeType]}" id="${_.nodeId}" ]`
          })
      const genEdge = () => {
        const lineName = {}
        this.flowData.nodeLinks
          && this.flowData.nodeLinks.forEach(link => {
            lineName[link.source + link.target] = link.name
          })
        const pathAry = []
        this.flowData
          && this.flowData.flowNodes
          && this.flowData.flowNodes.forEach(_ => {
            if (_.succeedingNodeIds.length > 0) {
              let current = []
              current = _.succeedingNodeIds.map(to => {
                const toNodeItem = this.flowData.flowNodes.find(i => i.nodeId === to) || {}
                const edgeColor = statusColor[toNodeItem.status] || '#505a68'
                // 修复判断分支多连线不能区分颜色问题
                if (_.nodeType === 'decision') {
                  return (
                    '"'
                    + _.nodeId
                    + '"'
                    + ' -> '
                    + `${'"' + to + '"'} [label="${lineName[_.nodeId + to]}" color="${edgeColor}" ]`
                  )
                }
                return (
                  '"'
                  + _.nodeId
                  + '"'
                  + ' -> '
                  + `${'"' + to + '"'} [label="${lineName[_.nodeId + to]}" color="${
                    excution ? statusColor[_.status] : 'black'
                  }"]`
                )
              })
              pathAry.push(current)
            }
          })
        return pathAry.flat().toString()
          .replace(/,/g, ';')
      }
      const nodesToString = Array.isArray(nodes) ? nodes.toString().replace(/,/g, ';') + ';' : ''
      const nodesString = 'digraph G {'
        + 'bgcolor="transparent";'
        + 'splines="polyline"'
        + 'Node [fontname=Arial, width=1.8, height=0.45, color="#505a68", fontsize=12]'
        + 'Edge [fontname=Arial, color="#505a68", fontsize=10];'
        + nodesToString
        + genEdge()
        + '}'
      this.flowGraph.graphviz
        .transition()
        .renderDot(nodesString)
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
  height: calc(100vh - 150px);
}
.reset-button {
  position: absolute;
  right: 20px;
  bottom: 20px;
  font-size: 12px;
}
</style>
