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

    <!-- <Modal
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
    </Modal> -->

    <!--节点操作弹窗(查看)-->
    <BaseDrawer
      :title="$t('select_an_operation')"
      :visible.sync="workflowActionModalVisible"
      realWidth="70%"
      :scrollable="true"
      class="json-viewer"
    >
      <template slot="content">
        <BaseHeaderTitle title="节点信息">
          <template v-if="nodeDetailResponseHeader && Object.keys(nodeDetailResponseHeader).length > 0">
            <json-viewer :value="nodeDetailResponseHeader" :expand-depth="5"></json-viewer>
          </template>
          <div v-else class="no-data">{{ $t('noData') }}</div>
        </BaseHeaderTitle>
        <BaseHeaderTitle title="API调用">
          <Table :columns="nodeDetailColumns" tooltip="true" :data="nodeDetailIO"> </Table>
        </BaseHeaderTitle>
      </template>
    </BaseDrawer>
  </div>
</template>
<script>
import { getFlowByInstanceId, getNodeContextByNodeId } from '@/api/server'
import * as d3 from 'd3-selection'
// eslint-disable-next-line no-unused-vars
import * as d3Graphviz from 'd3-graphviz'
import { addEvent, removeEvent } from './event.js'
import JsonViewer from 'vue-json-viewer'
export default {
  components: {
    JsonViewer
  },
  props: {
    flowId: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      currentModelNodeRefs: [],
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
                  <div key={index} style="margin-left:5px;">
                    {'{'}
                    {Object.entries(data).map(([key, value]) => (
                      <div style="margin-left:5px;">
                        <Icon
                          type="md-search"
                          onClick={() => this.handleClick(key, value)}
                          style="cursor:pointer;color:#2d8cf0"
                        />
                        {key}: <span style="color:#42b983;">{value}</span>
                      </div>
                    ))}
                    {'}'}
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
            const strOutput = params.row.outputs
            const noData = strOutput.every(i => i && Object.keys(i).length === 0)
            if (noData) {
              return <span>-</span>
            }
            return (
              <div style="white-space: nowrap; overflow: auto;">
                {strOutput.map((data, index) => (
                  <div key={index} style="margin-left:5px;">
                    {'{'}
                    {Object.entries(data).map(([key, value]) => (
                      <div style="margin-left:5px;">
                        {key}: <span style="color:#42b983;">{value}</span>
                      </div>
                    ))}
                    {'}'}
                  </div>
                ))}
              </div>
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
    // 打开编排执行详情页
    jumpToFlowDetail () {
      if (process.env.VUE_APP_PLUGIN === 'plugin') {
        window.sessionStorage.currentPath = '' // 先清空session缓存页面，不然打开新标签页面会回退到缓存的页面
        const path = `${window.location.origin}/#/implementation/workflow-execution/view-execution?id=${this.flowId}&from=noraml`
        window.open(path, '_blank')
      }
    },
    // 子编排调用API列表支持跳转预览子编排详情
    viewSubProcExecutionDetail (id) {
      if (process.env.VUE_APP_PLUGIN === 'plugin') {
        window.sessionStorage.currentPath = '' // 先清空session缓存页面，不然打开新标签页面会回退到缓存的页面
        const path = `${window.location.origin}/#/implementation/workflow-execution/view-execution?id=${id}&from=sub`
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
      const nodes = this.flowData &&
        this.flowData.flowNodes &&
        this.flowData.flowNodes
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
        this.flowData.nodeLinks &&
          this.flowData.nodeLinks.forEach(link => {
            lineName[link.source + link.target] = link.name
          })
        const pathAry = []
        this.flowData &&
          this.flowData.flowNodes &&
          this.flowData.flowNodes.forEach(_ => {
            if (_.succeedingNodeIds.length > 0) {
              let current = []
              current = _.succeedingNodeIds.map(to => {
                const toNodeItem = this.flowData.flowNodes.find(i => i.nodeId === to) || {}
                const edgeColor = statusColor[toNodeItem.status] || '#505a68'
                // 修复判断分支多连线不能区分颜色问题
                if (_.nodeType === 'decision') {
                  return (
                    '"' +
                    _.nodeId +
                    '"' +
                    ' -> ' +
                    `${'"' + to + '"'} [label="${lineName[_.nodeId + to]}" color="${edgeColor}" ]`
                  )
                }
                return (
                  '"' +
                  _.nodeId +
                  '"' +
                  ' -> ' +
                  `${'"' + to + '"'} [label="${lineName[_.nodeId + to]}" color="${
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
      const nodesString = 'digraph G {' +
        'bgcolor="transparent";' +
        'splines="polyline"' +
        'Node [fontname=Arial, width=1.8, height=0.45, color="#505a68", fontsize=12]' +
        'Edge [fontname=Arial, color="#505a68", fontsize=10];' +
        nodesToString +
        genEdge() +
        '}'
      this.flowGraph.graphviz
        .transition()
        .renderDot(nodesString)
        .on('end', () => {
          if (this.isEnqueryPage) {
            removeEvent('.retry', 'click', this.retryHandler)
            removeEvent('.normal', 'click', this.normalHandler)
            addEvent('.retry', 'click', this.retryHandler)
            addEvent('.normal', 'click', this.normalHandler)
            d3.selectAll('.retry').attr('cursor', 'pointer')
          } else {
            removeEvent('.retry', 'click', this.retryHandler)
            removeEvent('.normal', 'click', this.normalHandler)
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
            flowNodes: data.taskNodeInstances,
            nodeLinks: data.nodeLinks
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
      this.flowGraphMouseenterHandler(this.currentFailedNodeID)
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
      clearTimeout(this.flowDetailTimer)
      this.flowDetailTimer = setTimeout(async () => {
        const found = this.flowData.flowNodes.find(_ => _.nodeId === id)
        this.nodeTitle = (found.orderedNo ? found.orderedNo + '、' : '') + found.nodeName
        const { statusCode, data } = await getNodeContextByNodeId(found.procInstId, found.id)
        if (statusCode === 'OK') {
          this.nodeDetailResponseHeader = JSON.parse(JSON.stringify(data))
          this.pluginInfo = this.nodeDetailResponseHeader.pluginInfo
          delete this.nodeDetailResponseHeader.requestObjects
          this.nodeDetailIO = data.requestObjects.map(ro => {
            ro['inputs'] = this.replaceParams(ro['inputs'])
            ro['outputs'] = this.replaceParams(ro['outputs'])
            return ro
          })
          // 日志input output表格添加子编排查看按钮
          if (this.nodeDetailResponseHeader && this.nodeDetailResponseHeader.nodeType === 'subProc') {
            const hasFlag = this.nodeDetailColumns.some(i => i.key === 'procDefId')
            if (!hasFlag) {
              this.nodeDetailColumns.push({
                title: '子编排',
                key: 'procDefId',
                width: 200,
                render: (h, params) => {
                  let procDefName = ''
                  let procInsId = ''
                  let version = ''
                  if (Array.isArray(params.row.outputs) && params.row.outputs.length > 0) {
                    procDefName = params.row.outputs[0].procDefName || '-'
                    procInsId = params.row.outputs[0].procInsId || ''
                    version = params.row.outputs[0].version || ''
                  }
                  return (
                    <span
                      style="cursor:pointer;color:#5cadff;"
                      onClick={() => {
                        this.viewSubProcExecutionDetail(procInsId)
                      }}
                    >
                      {procDefName}
                      <Tag style="margin-left:2px">{version}</Tag>
                    </span>
                  )
                }
              })
            }
          } else {
            this.nodeDetailColumns = this.nodeDetailColumns.filter(i => i.key !== 'procDefId')
          }
        }
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
