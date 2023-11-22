<template>
  <div style="min-width: 500px;">
    <div style="margin-bottom: 8px;">
      <span class="custom-title">{{ $t('workflow_name') }}</span>
      <span class="custom-display">test_custom_demo 2023-10-12 20:34:94 (V5)</span>
    </div>
    <div id="graphcontain">
      <div class="graph-container" id="flow" style="height: 96%"></div>
      <Button class="reset-button" size="small" @click="ResetFlow">ResetZoom</Button>
    </div>
  </div>
</template>
<script>
import * as d3 from 'd3-selection'
// eslint-disable-next-line no-unused-vars
import * as d3Graphviz from 'd3-graphviz'
export default {
  data () {
    return {
      flow: 'test_custom_demo 2023-10-12 20:34:94 (V5)',
      flowGraph: {},
      selectedFlow: 'tUrPmy0f3sVj'
    }
  },
  mounted () {
    this.orchestrationSelectHandler()
  },
  methods: {
    orchestrationSelectHandler () {
      this.currentFlowNodeId = ''
      this.currentModelNodeRefs = []
      this.getFlowOutlineData(this.selectedFlow)
    },
    async getFlowOutlineData (id) {
      // if (!id) return
      // let { status, data } = await getFlowOutlineByID(id)
      // if (status === 'OK') {
      //   this.flowData = data
      //   this.initFlowGraph()
      // }
      this.flowData = {
        procDefId: 'tUrPmy0f3sVj',
        procDefKey: 'wecube1698921083968',
        procDefName: '3333333--22',
        procDefVersion: '2',
        status: 'deployed',
        procDefData: null,
        rootEntity: 'wecmdb:subsystem',
        createdTime: null,
        excludeMode: 'N',
        tags: null,
        permissionToRole: null,
        flowNodes: [
          {
            nodeId: 'SubProcess_06n94s2',
            nodeName: '4',
            nodeType: 'subProcess',
            nodeDefId: 'tUrPmyDf3t3P',
            status: 'deployed',
            orderedNo: '8',
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: 'wecmdb:subsystem',
            taskCategory: null,
            serviceId: 'monitor/dashboard(subsystem)/add',
            dynamicBind: 'N',
            description: '',
            previousNodeIds: ['ExclusiveGateway_1y3olpg'],
            succeedingNodeIds: ['SubProcess_0dr3g37']
          },
          {
            nodeId: 'SubProcess_1fh08x5',
            nodeName: '7',
            nodeType: 'subProcess',
            nodeDefId: 'tUrPmyJf3t4n',
            status: 'deployed',
            orderedNo: '6',
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: 'wecmdb:subsystem',
            taskCategory: null,
            serviceId: 'monitor/dashboard(subsystem)/add',
            dynamicBind: 'N',
            description: '',
            previousNodeIds: ['ExclusiveGateway_1y3olpg'],
            succeedingNodeIds: ['SubProcess_0qzm5xz']
          },
          {
            nodeId: 'SubProcess_0qzm5xz',
            nodeName: '8',
            nodeType: 'subProcess',
            nodeDefId: 'tUrPmyNf3t5T',
            status: 'deployed',
            orderedNo: '7',
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: 'wecmdb:subsystem',
            taskCategory: null,
            serviceId: 'monitor/dashboard(subsystem)/add',
            dynamicBind: 'N',
            description: '',
            previousNodeIds: ['SubProcess_1fh08x5'],
            succeedingNodeIds: ['ExclusiveGateway_1qwkoz6']
          },
          {
            nodeId: 'ExclusiveGateway_1022yhr',
            nodeName: '分1',
            nodeType: 'parallelGateway',
            nodeDefId: 'tUrPmyRf3t6c',
            status: 'deployed',
            orderedNo: null,
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: null,
            taskCategory: null,
            serviceId: null,
            dynamicBind: null,
            description: null,
            previousNodeIds: ['SubProcess_04ydxqc'],
            succeedingNodeIds: ['SubProcess_168fu10', 'SubProcess_0fg0vhx']
          },
          {
            nodeId: 'SubProcess_17826ce',
            nodeName: '9',
            nodeType: 'subProcess',
            nodeDefId: 'tUrPmyVf3t7U',
            status: 'deployed',
            orderedNo: '11',
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: 'wecmdb:subsystem',
            taskCategory: null,
            serviceId: 'monitor/dashboard(subsystem)/add',
            dynamicBind: 'N',
            description: '',
            previousNodeIds: ['ExclusiveGateway_1qwkoz6'],
            succeedingNodeIds: ['ExclusiveGateway_1f5jzvi']
          },
          {
            nodeId: 'SubProcess_0dr3g37',
            nodeName: '5',
            nodeType: 'subProcess',
            nodeDefId: 'tUrPmyZf3t8n',
            status: 'deployed',
            orderedNo: '9',
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: 'wecmdb:subsystem',
            taskCategory: null,
            serviceId: 'monitor/dashboard(subsystem)/add',
            dynamicBind: 'N',
            description: '',
            previousNodeIds: ['SubProcess_06n94s2'],
            succeedingNodeIds: ['SubProcess_07usvqz']
          },
          {
            nodeId: 'StartEvent_005xaf4',
            nodeName: '开始',
            nodeType: 'startEvent',
            nodeDefId: 'tUrPmyhf3sYX',
            status: 'deployed',
            orderedNo: null,
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: null,
            taskCategory: null,
            serviceId: null,
            dynamicBind: null,
            description: null,
            previousNodeIds: [],
            succeedingNodeIds: ['SubProcess_04ydxqc']
          },
          {
            nodeId: 'SubProcess_1ivoxai',
            nodeName: '3',
            nodeType: 'subProcess',
            nodeDefId: 'tUrPmykf3sZp',
            status: 'deployed',
            orderedNo: '4',
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: 'wecmdb:subsystem',
            taskCategory: null,
            serviceId: 'monitor/dashboard(subsystem)/add',
            dynamicBind: 'N',
            description: '',
            previousNodeIds: ['SubProcess_03hwg62'],
            succeedingNodeIds: ['ExclusiveGateway_1f5jzvi']
          },
          {
            nodeId: 'SubProcess_0fg0vhx',
            nodeName: '4',
            nodeType: 'subProcess',
            nodeDefId: 'tUrPmypf3t0c',
            status: 'deployed',
            orderedNo: '5',
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: 'wecmdb:subsystem',
            taskCategory: null,
            serviceId: 'monitor/dashboard(subsystem)/add',
            dynamicBind: 'N',
            description: '',
            previousNodeIds: ['ExclusiveGateway_1022yhr'],
            succeedingNodeIds: ['ExclusiveGateway_1y3olpg']
          },
          {
            nodeId: 'SubProcess_168fu10',
            nodeName: '1',
            nodeType: 'subProcess',
            nodeDefId: 'tUrPmysf3t1u',
            status: 'deployed',
            orderedNo: '2',
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: 'wecmdb:subsystem',
            taskCategory: null,
            serviceId: 'monitor/dashboard(subsystem)/add',
            dynamicBind: 'N',
            description: '',
            previousNodeIds: ['ExclusiveGateway_1022yhr'],
            succeedingNodeIds: ['SubProcess_03hwg62']
          },
          {
            nodeId: 'SubProcess_03hwg62',
            nodeName: '2',
            nodeType: 'subProcess',
            nodeDefId: 'tUrPmyxf3t2p',
            status: 'deployed',
            orderedNo: '3',
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: 'wecmdb:subsystem',
            taskCategory: null,
            serviceId: 'monitor/dashboard(subsystem)/add',
            dynamicBind: 'N',
            description: '',
            previousNodeIds: ['SubProcess_168fu10'],
            succeedingNodeIds: ['SubProcess_1ivoxai']
          },
          {
            nodeId: 'SubProcess_07usvqz',
            nodeName: '6',
            nodeType: 'subProcess',
            nodeDefId: 'tUrPmz3f3t9d',
            status: 'deployed',
            orderedNo: '10',
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: 'wecmdb:subsystem',
            taskCategory: null,
            serviceId: 'monitor/dashboard(subsystem)/add',
            dynamicBind: 'N',
            description: '',
            previousNodeIds: ['SubProcess_0dr3g37'],
            succeedingNodeIds: ['ExclusiveGateway_1qwkoz6']
          },
          {
            nodeId: 'ExclusiveGateway_1y3olpg',
            nodeName: '分2',
            nodeType: 'parallelGateway',
            nodeDefId: 'tUrPmz7f3ta2',
            status: 'deployed',
            orderedNo: null,
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: null,
            taskCategory: null,
            serviceId: null,
            dynamicBind: null,
            description: null,
            previousNodeIds: ['SubProcess_0fg0vhx'],
            succeedingNodeIds: ['SubProcess_1fh08x5', 'SubProcess_06n94s2']
          },
          {
            nodeId: 'ExclusiveGateway_1f5jzvi',
            nodeName: '合1',
            nodeType: 'parallelGateway',
            nodeDefId: 'tUrPmzaf3tby',
            status: 'deployed',
            orderedNo: null,
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: null,
            taskCategory: null,
            serviceId: null,
            dynamicBind: null,
            description: null,
            previousNodeIds: ['SubProcess_17826ce', 'SubProcess_1ivoxai'],
            succeedingNodeIds: ['EndEvent_0ch083n']
          },
          {
            nodeId: 'SubProcess_04ydxqc',
            nodeName: '0',
            nodeType: 'subProcess',
            nodeDefId: 'tUrPmzef3tcq',
            status: 'deployed',
            orderedNo: '1',
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: 'wecmdb:subsystem',
            taskCategory: null,
            serviceId: 'monitor/dashboard(subsystem)/add',
            dynamicBind: 'N',
            description: '',
            previousNodeIds: ['StartEvent_005xaf4'],
            succeedingNodeIds: ['ExclusiveGateway_1022yhr']
          },
          {
            nodeId: 'ExclusiveGateway_1qwkoz6',
            nodeName: '合2',
            nodeType: 'parallelGateway',
            nodeDefId: 'tUrPmzhf3tdO',
            status: 'deployed',
            orderedNo: null,
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: null,
            taskCategory: null,
            serviceId: null,
            dynamicBind: null,
            description: null,
            previousNodeIds: ['SubProcess_07usvqz', 'SubProcess_0qzm5xz'],
            succeedingNodeIds: ['SubProcess_17826ce']
          },
          {
            nodeId: 'EndEvent_0ch083n',
            nodeName: '结束',
            nodeType: 'endEvent',
            nodeDefId: 'tUrPmzlf3teA',
            status: 'deployed',
            orderedNo: null,
            procDefId: 'tUrPmy0f3sVj',
            procDefKey: 'wecube1698921083968',
            routineExpression: null,
            taskCategory: null,
            serviceId: null,
            dynamicBind: null,
            description: null,
            previousNodeIds: ['ExclusiveGateway_1f5jzvi'],
            succeedingNodeIds: []
          }
        ]
      }
      this.initFlowGraph()
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
  padding: 10px 12px 10px 0;
  box-sizing: border-box;
}

.custom-display {
  display: inline-block;
  width: 300px;
  height: 32px;
  line-height: 1.5;
  padding: 4px 7px;
  font-size: 14px;
  border: 1px solid #069cec;
  border-radius: 4px;
  color: #069cec;
  // cursor: pointer;
}
#graphcontain {
  border: 1px solid #d3cece;
  border-radius: 3px;
  padding: 5px;
  height: calc(100vh - 220px);
}
.reset-button {
  position: absolute;
  right: 20px;
  bottom: 5px;
  font-size: 12px;
}
</style>
