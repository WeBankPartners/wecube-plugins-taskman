<template>
  <div class="workbench-request-create">
    <Row class="w-header">
      <Col span="18" class="steps">
        <span class="title">请求进度</span>
        <Steps :current="0" style="width:400px">
          <Step v-for="(i, index) in steps" :key="index" :content="i.name">
            <template #icon>
              <Icon size="28" :type="i.icon" />
            </template>
          </Step>
        </Steps>
      </Col>
      <Col span="6" class="btn-group">
        <Button @click="handleDraft">保存草稿</Button>
        <Button type="primary" @click="handlePublish" class="primary">发布</Button>
      </Col>
    </Row>
    <Row class="content">
      <Col span="16" class="split-line">
        <Form :model="form" label-position="right" :label-width="120">
          <HeaderTitle title="发布信息">
            <FormItem label="请求名称" required>
              <Input placeholder="请输入" style="width:400px;" />
            </FormItem>
            <FormItem label="发布描述">
              <Input placeholder="请输入" style="width:400px;" />
            </FormItem>
          </HeaderTitle>
          <HeaderTitle title="发布目标对象">
            <FormItem label="选择操作单元" required>
              <Select v-model="form.unit" clearable filterable style="width:250px;">
                <Option v-for="i in 5" :value="i" :key="i">{{ i }}</Option>
              </Select>
              <Button type="primary" @click="handleChooseExample" class="primary">点击勾选操作实例</Button>
            </FormItem>
            <FormItem v-if="chooseExampleData.length" label="已选择">
              <RadioGroup v-model="form.exampleType">
                <Radio v-for="(i, idx) in chooseExampleData" :label="i.parent" :key="idx" border>{{
                  `${i.title}(${i.selectCount})个`
                }}</Radio>
              </RadioGroup>
              <Tabs value="name1" style="margin-top:20px;">
                <TabPane label="标签一" name="name1"></TabPane>
                <TabPane label="标签二" name="name2"></TabPane>
                <TabPane label="标签三" name="name3"></TabPane>
              </Tabs>
              <Table size="small" :columns="tableColumns" :data="tableData"></Table>
            </FormItem>
          </HeaderTitle>
        </Form>
      </Col>
      <Col span="8">
        <div>1111111111</div>
      </Col>
    </Row>
    <ChooseExampleDrawer
      v-if="chooseExampleVisible"
      :visible.sync="chooseExampleVisible"
      @getData="getChooseExampleData"
    ></ChooseExampleDrawer>
  </div>
</template>

<script>
import HeaderTitle from '../components/header-title.vue'
import ChooseExampleDrawer from '../publish/choose-example.vue'
export default {
  components: {
    HeaderTitle,
    ChooseExampleDrawer
  },
  data () {
    return {
      form: {
        unit: '',
        exampleType: 1
      },
      steps: [
        { name: '提起请求', status: 'process', icon: 'ios-pin' },
        { name: '请求定版', status: 'wait', icon: 'md-radio-button-on' },
        { name: '任务1审批', status: 'wait', icon: 'md-radio-button-on' },
        { name: '任务2审批', status: 'wait', icon: 'md-radio-button-on' },
        { name: '请求完成', status: 'wait', icon: 'md-radio-button-on' }
      ],
      chooseExampleVisible: false,
      chooseExampleData: [], // 勾选的的实例
      exampleTabName: '', // 当前选中实例tab
      exampleTabList: [],
      tableColumns: [
        {
          title: '序号',
          key: 'no'
        },
        {
          title: 'IP',
          key: 'ip'
        },
        {
          title: '部署包',
          key: 'package'
        },
        {
          title: 'DCN',
          key: 'dcn'
        },
        {
          title: this.$t('t_action'),
          key: 'action',
          width: 130,
          fixed: 'right',
          align: 'center',
          render: (h, params) => {
            return (
              <div>
                <Icon size="24" type="md-trash" color="#ed4014" />
                <Icon size="24" type="md-create" />
                <Icon size="24" type="md-eye" />
              </div>
            )
          }
        }
      ],
      tableData: [
        { no: 1, ip: '198.168.150.12', package: 'dsad.dada.dasda.dsa', dcn: '321312313312' },
        { no: 1, ip: '198.168.150.12', package: 'dsad.dada.dasda.dsa', dcn: '321312313312' },
        { no: 1, ip: '198.168.150.12', package: 'dsad.dada.dasda.dsa', dcn: '321312313312' },
        { no: 1, ip: '198.168.150.12', package: 'dsad.dada.dasda.dsa', dcn: '321312313312' }
      ]
    }
  },
  methods: {
    // 操作实例弹窗
    handleChooseExample () {
      this.chooseExampleVisible = true
    },
    // 获取操作实例弹窗数据
    getChooseExampleData (data) {
      this.chooseExampleData = data
    },
    // 保存草稿
    handleDraft () {},
    // 发布
    handlePublish () {}
  }
}
</script>

<style lang="scss" scoped>
.workbench-request-create {
  .w-header {
    padding-bottom: 20px;
    margin-bottom: 20px;
    border-bottom: 1px solid #e8eaec;
    .steps {
      display: flex;
      align-items: center;
      .title {
        font-size: 14px;
        font-weight: bold;
        margin-right: 20px;
      }
    }
    .btn-group {
      text-align: right;
    }
  }
  .content {
    min-height: 500px;
    .split-line {
      border-right: 1px solid #e8eaec;
      padding-right: 20px;
    }
  }
  .primary {
    margin-left: 20px;
  }
}
</style>
<style lang="scss">
.workbench-request-create {
  .ivu-steps-content {
    padding-left: 0px;
    padding-top: 5px;
    font-size: 12px;
    color: #3d3c38 !important;
  }
  .ivu-radio {
    display: none;
  }
  .ivu-radio-wrapper {
    border-radius: 20px;
    font-size: 12px;
    color: #000;
    background: #fff;
  }
  .ivu-radio-wrapper-checked.ivu-radio-border {
    border-color: #2d8cf0;
    background: #2d8cf0;
    color: #fff;
  }
  .ivu-form-item {
    margin-bottom: 25px;
  }
  .ivu-form-item-content {
    line-height: 20px;
  }
}
</style>
