<template>
  <div class="workbench-request-create">
    <Row class="back-header">
      <Icon size="24" type="md-arrow-back" style="cursor:pointer" @click="$router.back()" />
      <span>{{ detailInfo.templateName || form.name }}</span>
    </Row>
    <Row class="w-header">
      <Col span="18" class="steps">
        <span class="title">请求进度：</span>
        <Steps :current="0" style="max-width:600px;">
          <Step v-for="(i, index) in progressList" :key="index" :content="i.name">
            <template #icon>
              <Icon size="24" :type="i.icon" :color="i.color" />
            </template>
            <div class="role" slot="content">
              <div class="word-eclipse">{{ i.name }}</div>
              <span>{{ i.handler }}</span>
            </div>
          </Step>
        </Steps>
      </Col>
      <Col span="6" class="btn-group">
        <Button @click="handleDraft(false)" style="margin-right:10px;">保存草稿</Button>
        <Button type="primary" @click="handlePublish">发布</Button>
      </Col>
    </Row>
    <Row class="content">
      <Col span="16" class="split-line">
        <Form :model="form" label-position="right" :label-width="120">
          <template v-if="isAdd">
            <HeaderTitle title="发布信息">
              <FormItem label="请求名称" required>
                <Input v-model="form.name" placeholder="请输入" style="width:400px;" />
              </FormItem>
              <FormItem label="发布描述">
                <Input v-model="form.description" type="textarea" placeholder="请输入" style="width:400px;" />
              </FormItem>
              <FormItem :label="$t('expected_completion_time')">
                <DatePicker
                  type="datetime"
                  :value="form.expectTime"
                  @on-change="
                    val => {
                      form.expectTime = val
                    }
                  "
                  placeholder="Select date"
                  :options="{
                    disabledDate(date) {
                      return date && date.valueOf() < Date.now() - 86400000
                    }
                  }"
                  style="width:400px;"
                ></DatePicker>
              </FormItem>
            </HeaderTitle>
            <HeaderTitle title="发布目标对象">
              <DataCrud
                ref="dataCrud"
                :requestId="requestId"
                :formDisable="formDisable"
                :jumpFrom="jumpFrom"
              ></DataCrud>
            </HeaderTitle>
          </template>
          <template v-else>
            <HeaderTitle title="请求信息">
              <Row>
                <Col :span="4">请求ID：</Col>
                <Col :span="8">{{ detailInfo.id }}</Col>
                <Col :span="4">请求类型：</Col>
                <Col :span="8">{{ { 0: '请求', 1: '发布' }[detailInfo.requestType] }}</Col>
              </Row>
              <Row style="margin-top:10px;">
                <Col :span="4">创建时间：</Col>
                <Col :span="8">{{ detailInfo.createdTime }}</Col>
                <Col :span="4">期望完成时间：</Col>
                <Col :span="8">{{ detailInfo.expectTime }}</Col>
              </Row>
              <Row style="margin-top:10px;">
                <Col :span="4">请求进度：</Col>
                <Col :span="8">
                  <Progress :percent="detailInfo.progress" style="width:150px;" />
                </Col>
                <Col :span="4">请求状态：</Col>
                <Col :span="8">{{ getStatusName(detailInfo.status) }}</Col>
              </Row>
              <Row style="margin-top:10px;">
                <Col :span="4">当前节点：</Col>
                <Col :span="8">{{
                  {
                    sendRequest: '提起请求',
                    requestPending: '请求定版',
                    requestComplete: '请求完成'
                  }[detailInfo.curNode] || detailInfo.curNode
                }}</Col>
                <Col :span="4">当前处理人：</Col>
                <Col :span="8">{{ detailInfo.handler }}</Col>
              </Row>
              <Row style="margin-top:10px;">
                <Col :span="4">创建人：</Col>
                <Col :span="8">{{ detailInfo.createdBy }}</Col>
                <Col :span="4">创建人角色：</Col>
                <Col :span="8">{{ detailInfo.role }}</Col>
              </Row>
              <Row style="margin-top:10px;">
                <Col :span="4">使用模板：</Col>
                <Col :span="8">{{ detailInfo.templateName }}</Col>
                <Col :span="4">模板组：</Col>
                <Col :span="8">{{ detailInfo.templateGroupName }}</Col>
              </Row>
              <Row style="margin-top:10px;">
                <Col :span="4">请求描述：</Col>
                <Col :span="8">{{ detailInfo.description }}</Col>
              </Row>
            </HeaderTitle>
            <HeaderTitle title="处理历史">
              <!-- <DataCrud ref="dataCrud" :requestId="requestId" :formDisable="formDisable" :jumpFrom="jumpFrom"></DataCrud> -->
              <DataBind
                :isHandle="isHandle"
                :requestTemplate="requestTemplate"
                :requestId="requestId"
                :formDisable="formDisable"
              ></DataBind>
            </HeaderTitle>
          </template>
        </Form>
      </Col>
      <!--编排流程-->
      <Col span="8">
        <div style="padding: 0 20px">
          <StaticFlow v-if="isAdd" :templateId="requestTemplate" ref="staticFlowRef"></StaticFlow>
          <DynamicFlow v-else ref="staticFlowRef"></DynamicFlow>
        </div>
      </Col>
    </Row>
  </div>
</template>

<script>
import HeaderTitle from '../components/header-title.vue'
import { getCreateInfo, getProgressInfo, savePublishData, getPublishInfo, updateRequestStatus } from '@/api/server'
import StaticFlow from '../publish/flow/static-flow.vue'
import DynamicFlow from '../publish/flow/dynamic-flow.vue'
import DataCrud from './components/data-crud.vue'
import DataBind from './components/data-bind.vue'
import dayjs from 'dayjs'
const statusIcon = {
  1: 'md-pin',
  2: 'md-radio-button-on',
  3: 'ios-checkmark-circle-outline',
  4: 'ios-close-circle-outline'
}
const statusColor = {
  1: '#ffa500',
  2: '#8189a5',
  3: '#19be6b',
  4: '#ed4014'
}
export default {
  components: {
    HeaderTitle,
    StaticFlow,
    DynamicFlow,
    DataCrud,
    DataBind
  },
  data () {
    return {
      detailInfo: {}, // 详情信息
      form: {
        name: '',
        description: '',
        expectTime: '',
        rootEntityId: '', // 目标对象
        entityName: '',
        data: []
      },
      progressList: [],
      isAdd: this.$route.query.isAdd === 'Y',
      isHandle: this.$route.query.isHandle === 'Y', // 处理标志
      formDisable: this.$route.query.isCheck === 'Y', // 查看标志
      jumpFrom: this.$route.query.jumpFrom, // 入口tab标记
      requestTemplate: this.$route.query.requestTemplate,
      requestId: this.$route.query.requestId,
      procDefId: '',
      procDefKey: ''
    }
  },
  computed: {
    getStatusName () {
      return function (val) {
        const list = [
          { label: this.$t('status_pending'), value: 'Pending', color: '#b886f8' },
          { label: this.$t('status_inProgress'), value: 'InProgress', color: '#1990ff' },
          { label: this.$t('status_inProgress_faulted'), value: 'InProgress(Faulted)', color: '#f26161' },
          { label: this.$t('status_termination'), value: 'Termination', color: '#e29836' },
          { label: this.$t('status_complete'), value: 'Completed', color: '#7ac756' },
          { label: this.$t('status_inProgress_timeouted'), value: 'InProgress(Timeouted)', color: '#f26161' },
          { label: this.$t('status_faulted'), value: 'Faulted', color: '#e29836' },
          { label: this.$t('status_draft'), value: 'Draft', color: '#808695' }
        ]
        const item = list.find(i => i.value === val) || {}
        return item.label
      }
    }
  },
  async mounted () {
    this.$refs.staticFlowRef.orchestrationSelectHandler(this.requestTemplate)
    // 查看调用详情接口
    if (this.requestId) {
      await this.getPublishInfo()
    } else {
      await this.getCreateInfo()
    }
    this.getProgressInfo()
  },
  methods: {
    // 创建发布,使用模板ID获取详情数据
    async getCreateInfo () {
      const params = {
        requestTemplate: this.requestTemplate,
        role: this.$route.query.role
      }
      const { statusCode, data } = await getCreateInfo(params)
      if (statusCode === 'OK') {
        this.requestId = data.id
        this.form.name = data.name
        this.form.expectTime = dayjs()
          .add(data.expireDay || 0, 'day')
          .format('YYYY-MM-DD HH:mm:ss')
      }
    },
    // 获取请求进度
    async getProgressInfo () {
      const params = {
        templateId: this.requestTemplate,
        requestId: this.requestId
      }
      const { statusCode, data } = await getProgressInfo(params)
      if (statusCode === 'OK') {
        this.progressList = data
        this.progressList.forEach(item => {
          item.icon = statusIcon[item.status]
          item.color = statusColor[item.status]
          switch (item.node) {
            case 'sendRequest':
              item.name = '提起请求'
              break
            case 'requestPending':
              item.name = '请求定版'
              break
            case 'requestComplete':
              item.name = '请求完成'
              break
            default:
              item.name = item.node
              break
          }
        })
      }
    },
    // 获取发布信息
    async getPublishInfo () {
      const { statusCode, data } = await getPublishInfo(this.requestId)
      if (statusCode === 'OK') {
        this.detailInfo = data.request || {}
        const { name, description, expectTime } = this.detailInfo
        this.form = Object.assign({}, this.form, {
          name,
          description,
          expectTime
        })
      }
    },
    // 保存草稿
    async handleDraft (noJump) {
      this.form = Object.assign({}, this.form, {
        rootEntityId: this.$refs.dataCrud.rootEntityId,
        entityName: this.$refs.dataCrud.entityName,
        data: this.$refs.dataCrud.requestData
      })
      if (this.form.rootEntityId === '') {
        this.$Message.warning(this.$t('root_entity') + this.$t('can_not_be_empty'))
        return
      }
      if (this.requestDataCheck()) {
        const { statusCode } = await savePublishData(this.requestId, this.form)
        if (statusCode === 'OK') {
          if (noJump) {
            return statusCode
          } else {
            this.$router.push({ path: '/taskman/workbench?tabName=pending' })
          }
        }
      } else {
        this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('required_tip')
        })
      }
    },
    // 发布
    async handlePublish () {
      this.$Modal.confirm({
        title: this.$t('confirm') + this.$t('commit'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const draftResult = await this.handleDraft(true)
          if (draftResult === 'OK') {
            const { statusCode } = await updateRequestStatus(this.requestId, 'Pending')
            if (statusCode === 'OK') {
              this.$router.push({ path: '/taskman/workbench?tabName=pending' })
            }
          }
        },
        onCancel: () => {}
      })
    },
    requestDataCheck () {
      let result = true
      this.form.data.forEach(requestData => {
        let requiredName = []
        requestData.title.forEach(t => {
          if (t.required === 'yes') {
            requiredName.push(t.name)
          }
        })
        requestData.value.forEach(v => {
          requiredName.forEach(key => {
            let val = v.entityData[key]
            if (Array.isArray(val)) {
              if (val.length === 0) {
                result = false
              }
            } else {
              if (val === '') {
                result = false
              }
            }
          })
        })
      })
      return result
    }
  }
}
</script>

<style lang="scss" scoped>
.workbench-request-create {
  .back-header {
    display: flex;
    align-items: center;
    margin-bottom: 10px;
    span {
      font-size: 14px;
      margin-left: 15px;
    }
  }
  .w-header {
    padding-bottom: 20px;
    margin-bottom: 20px;
    border-bottom: 2px dashed #d7dadc;
    .steps {
      display: flex;
      align-items: center;
      .title {
        font-size: 14px;
        font-weight: 500;
        margin-right: 20px;
      }
      .role {
        display: flex;
        flex-direction: column;
      }
      .word-eclipse {
        max-width: 100px;
        text-overflow: ellipsis;
        overflow: hidden;
        white-space: nowrap;
      }
    }
    .btn-group {
      text-align: right;
    }
  }
  .content {
    min-height: 500px;
    .split-line {
      border-right: 2px dashed #d7dadc;
      padding-right: 20px;
    }
    .count {
      font-weight: bold;
      font-size: 14px;
      margin-left: 10px;
    }
    .program {
      padding: 20px;
    }
  }
}
</style>
<style lang="scss">
.workbench-request-create {
  .ivu-steps-content {
    padding-left: 0px !important;
    padding-top: 5px;
    font-size: 12px;
    color: #3d3c38 !important;
  }
  .ivu-steps .ivu-steps-tail > i {
    height: 3px;
    background: #8189a5;
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
  .required {
    color: red;
  }
}
</style>
