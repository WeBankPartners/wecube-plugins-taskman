<template>
  <div style="width: 400px;">
    <div style="margin-bottom: 8px;">
      <Alert type="warning" show-icon>{{ $t('tw_flow_tips') }}</Alert>
      <!-- <span class="custom-title">{{ $t('workflow_name') }}</span> -->
      <span class="custom-display" @click="jumpToFlowDetail">{{ flowData.procInstName }} {{ flowData.operator }}</span>
    </div>
    <div id="graphcontain">
      <div class="graph-container" id="flow" style="height: 96%"></div>
      <Button class="reset-button" size="small" @click="ResetFlow">ResetZoom</Button>
    </div>

    <Modal
      :title="$t('select_an_operation')"
      v-model="workflowActionModalVisible"
      :footer-hide="true"
      :mask-closable="false"
      :scrollable="true"
    >
      <div class="workflowActionModal-container" style="text-align: center; margin-top: 20px">
        <Button
          type="info"
          v-show="['InProgress', 'Faulted', 'Timeouted', 'Completed', 'Risky'].includes(currentNodeStatus)"
          @click="workFlowActionHandler('showlog')"
          style="margin-left: 10px"
          >{{ $t('show_log') }}</Button
        >
      </div>
    </Modal>

    <Modal v-model="showNodeDetail" :fullscreen="nodeDetailFullscreen" width="1000" :styles="{ top: '50px' }">
      <p slot="header">
        <span>{{ nodeTitle }}</span>
        <Icon v-if="!nodeDetailFullscreen" @click="zoomModal" class="header-icon" type="ios-expand" />
        <Icon v-else @click="nodeDetailFullscreen = false" class="header-icon" type="ios-contract" />
      </p>
      <div v-if="!isTargetNodeDetail" :style="[{ overflow: 'auto' }, fullscreenModalContentStyle]">
        <h5>Data:</h5>
        <pre style="margin: 0 6px 6px" v-html="nodeDetailResponseHeader"></pre>
        <h5>requestObjects:</h5>
        <Table :columns="nodeDetailColumns" :max-height="tableMaxHeight" tooltip="true" :data="nodeDetailIO"> </Table>
      </div>
      <div v-else :style="[{ overflow: 'auto', margin: '0 6px 6px' }, fullscreenModalContentStyle]">
        <!-- v-html="nodeDetail" -->
        <json-viewer :value="nodeDetail" :expand-depth="5"></json-viewer>
      </div>
    </Modal>
  </div>
</template>
<script>
import { getFlowByInstanceId, getNodeContextByNodeId } from '@/api/server'
import * as d3 from 'd3-selection'
// eslint-disable-next-line no-unused-vars
import * as d3Graphviz from 'd3-graphviz'
import { addEvent, removeEvent } from './event.js'
export default {
  props: {
    flowId: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      currentModelNodeRefs: [],
      showNodeDetail: false,
      isTargetNodeDetail: false,
      nodeDetailFullscreen: false,
      fullscreenModalContentStyle: { 'max-height': '400px' },
      tableMaxHeight: 250,
      nodeTitle: null,
      nodeDetail: null,
      graph: {},
      flowGraph: {},
      flowData: {},
      currentNodeTitle: null,

      currentFlowNodeId: '',

      isEnqueryPage: false,
      workflowActionModalVisible: false,

      tartetModels: [],
      retryTartetModels: [],
      catchTartetModels: [],
      nodeDetailColumns: [
        {
          title: 'inputs',
          key: 'inputs',
          render: (h, params) => {
            const jsonData = params.row.inputs
            return (
              <div style="white-space: nowrap; overflow: auto;">
                {jsonData.map((data, index) => (
                  <div key={index}>
                    {'{'}
                    {Object.entries(data).map(([key, value]) => (
                      <div>
                        {key}: {value}
                      </div>
                    ))}
                    {'},'}
                  </div>
                ))}
              </div>
            )
          }
        },
        {
          title: 'outputs',
          key: 'outputs',
          render: (h, params) => {
            const strOutput = JSON.stringify(params.row.outputs)
              .split(',')
              .join(',<br/>')
            return h(
              'div',
              {
                domProps: {
                  innerHTML: `<pre>${strOutput}</pre>`
                }
              },
              []
            )
          }
        }
      ],
      nodeDetailIO: [],
      nodeDetailResponseHeader: null,
      currentFailedNodeID: '',
      timer: null,
      flowDetailTimer: null,
      catchNodeTableList: [],
      retryCatchNodeTableList: [],

      stopSuccess: false,
      currentInstanceStatus: true,
      pluginInfo: ''
    }
  },
  computed: {
    currentNodeStatus () {
      if (!this.flowData.flowNodes) {
        return ''
      }
      const found = this.flowData.flowNodes.find(_ => _.nodeId === this.currentFailedNodeID)
      if (found) {
        return found.status
      } else {
        return ''
      }
    }
  },
  watch: {
    nodeDetailFullscreen: function (tag) {
      tag ? (this.fullscreenModalContentStyle = {}) : (this.fullscreenModalContentStyle['max-height'] = '400px')
    },
    flowId: {
      handler (val) {
        if (val) {
          this.orchestrationSelectHandler()
        }
      },
      immediate: true
    }
  },
  destroyed () {
    clearInterval(this.timer)
  },
  methods: {
    jumpToFlowDetail() {
      if (process.env.PLUGIN === 'plugin') {
        window.sessionStorage.currentPath = '' // 先清空session缓存页面，不然打开新标签页面会回退到缓存的页面
        const path = `${window.location.origin}/#/implementation/workflow-execution/view-execution?id=${this.flowId}&from=noraml`
        window.open(path, '_blank')
      }
    },
    orchestrationSelectHandler () {
      this.stop()
      this.isEnqueryPage = true
      this.$nextTick(async () => {
        // 定时拿状态
        this.start()
        let { statusCode, data } = await getFlowByInstanceId(this.flowId)
        const flowData = data || {}
        if (statusCode === 'OK') {
          this.flowData = {
            ...flowData,
            flowNodes: flowData.taskNodeInstances
          }
          removeEvent('.retry', 'click', this.retryHandler)
          this.initFlowGraph(true)
        }
      })
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

      this.flowGraph.graphviz
        .transition()
        .renderDot(nodesString)
        .on('end', () => {
          if (this.isEnqueryPage) {
            removeEvent('.retry', 'click', this.retryHandler)
            addEvent('.retry', 'click', this.retryHandler)
            d3.selectAll('.retry').attr('cursor', 'pointer')
          } else {
            removeEvent('.retry', 'click', this.retryHandler)
          }
        })
      this.bindFlowEvent()
    },
    start () {
      if (this.timer === null) {
        this.getStatus()
      }
      if (this.timer !== null) {
        this.stop()
      }
      this.timer = setInterval(() => {
        this.getStatus()
      }, 5000)
    },
    stop () {
      clearInterval(this.timer)
    },
    async getStatus () {
      let { statusCode, data } = await getFlowByInstanceId(this.flowId)
      if (statusCode === 'OK') {
        if (
          !this.flowData.flowNodes ||
          (this.flowData.flowNodes && this.comparativeData(this.flowData.flowNodes, data.taskNodeInstances))
        ) {
          this.flowData = {
            ...data,
            flowNodes: data.taskNodeInstances
          }
          removeEvent('.retry', 'click', this.retryHandler)
          this.initFlowGraph(true)
        }
        if (data.status === 'Completed' || data.status === 'InternallyTerminated') {
          this.stopSuccess = true
          this.stop()
        }
      }
    },
    comparativeData (old, newData) {
      let isNew = false
      newData.forEach(_ => {
        const found = old.find(d => d.nodeId === _.nodeId)
        if (found && found.status !== _.status) {
          isNew = true
        }
      })
      return isNew
    },
    retryHandler (e) {
      this.currentFailedNodeID = e.target.parentNode.getAttribute('id')
      this.workflowActionModalVisible = true
    },
    async workFlowActionHandler (type) {
      const found = this.flowData.flowNodes.find(_ => _.nodeId === this.currentFailedNodeID)
      if (!found) {
        return
      }
      if (type === 'showlog') {
        this.flowGraphMouseenterHandler(this.currentFailedNodeID)
      }
    },
    bindFlowEvent () {
      if (this.isEnqueryPage !== true) {
        addEvent('.flow', 'mouseover', e => {
          e.preventDefault()
          e.stopPropagation()
          d3.selectAll('g').attr('cursor', 'pointer')
        })
        removeEvent('.flow', 'click', this.flowNodesClickHandler)
        addEvent('.flow', 'click', this.flowNodesClickHandler)
      } else {
        removeEvent('.flow', 'click', this.flowNodesClickHandler)
      }
    },
    flowGraphMouseenterHandler (id) {
      console.log(id)
      // Task_0f9a25l
      clearTimeout(this.flowDetailTimer)
      this.flowDetailTimer = setTimeout(async () => {
        const found = this.flowData.flowNodes.find(_ => _.nodeId === id)
        this.nodeTitle = (found.orderedNo ? found.orderedNo + '、' : '') + found.nodeName
        const { statusCode, data } = await getNodeContextByNodeId(found.procInstId, found.id)
        if (statusCode === 'OK') {
          this.workflowActionModalVisible = false
          this.nodeDetailResponseHeader = JSON.parse(JSON.stringify(data))
          this.pluginInfo = this.nodeDetailResponseHeader.pluginInfo
          delete this.nodeDetailResponseHeader.requestObjects
          this.nodeDetailResponseHeader = JSON.stringify(this.replaceParams(this.nodeDetailResponseHeader))
            .split(',')
            .join(',<br/>')
          this.nodeDetailResponseHeader = this.nodeDetailResponseHeader.replace(
            'errorMessage',
            "<span style='color:red'>errorMessage</span>"
          )
          // const reg = new RegExp(data.errorMessage, 'g')
          this.nodeDetailResponseHeader = this.nodeDetailResponseHeader.replace(
            data.errorMessage,
            "<span style='color:red'>" + data.errorMessage + '</span>'
          )
          this.nodeDetailIO = data.requestObjects.map(ro => {
            ro['inputs'] = this.replaceParams(ro['inputs'])
            ro['outputs'] = this.replaceParams(ro['outputs'])
            return ro
          })
        }
        this.nodeDetailFullscreen = false
        this.isTargetNodeDetail = false
        this.showNodeDetail = true
        this.tableMaxHeight = 250
      }, 0)
    },
    replaceParams (obj) {
      let placeholder = new Array(16).fill('&nbsp;')
      placeholder.unshift('<br/>')
      for (let key in obj) {
        if (obj[key] !== null && typeof obj[key] === 'string') {
          obj[key] = obj[key].replace('\r\n', placeholder.join(''))
        }
      }
      return obj
    },
    flowNodesClickHandler (e) {
      e.preventDefault()
      e.stopPropagation()
      let g = e.currentTarget
      this.currentFlowNodeId = g.id
      const currentNode = this.flowData.flowNodes.find(_ => {
        return _.nodeId === this.currentFlowNodeId
      })
      this.currentNodeTitle = `${currentNode.orderedNo}、${currentNode.nodeName}`
      this.renderFlowGraph()
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
    },
    zoomModal () {
      this.tableMaxHeight = document.body.scrollHeight - 410
      this.nodeDetailFullscreen = true
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
}

.header-icon {
  margin: 3px 40px 0 0 !important;
}
#graphcontain {
  position: relative;
  border: 1px solid #d3cece;
  border-radius: 3px;
  padding: 5px;
  height: calc(100vh - 150px);
}
.model_target .ivu-modal-content-drag {
  right: 40px;
}
.graph-container {
  overflow: auto;
}
.header-icon {
  float: right;
  margin: 3px 20px 0 0;
}
.reset-button {
  position: absolute;
  right: 20px;
  bottom: 20px;
  font-size: 12px;
}
</style>
