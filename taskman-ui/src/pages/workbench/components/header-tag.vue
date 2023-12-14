<template>
  <div class="header-tag">
    <Row v-if="showHeader" class="title">
      <Col :span="4">处理时间</Col>
      <Col :span="3">耗时</Col>
      <Col :span="3">处理角色</Col>
      <Col :span="3">处理人</Col>
      <Col :span="3">操作</Col>
      <Col :span="8">备注</Col>
    </Row>
    <Row class="content">
      <Col :span="4">{{ data.handleTime }}</Col>
      <Col :span="3">{{ getDiffTime }}</Col>
      <Col :span="3">{{ data.handleRoleName || '-' }}</Col>
      <Col :span="3">{{ data.handler || '-' }}</Col>
      <Col :span="3">{{ data.choseOption || operation }}</Col>
      <Col :span="8">{{ data.comment }}</Col>
    </Row>
  </div>
</template>

<script>
import dayjs from 'dayjs'
export default {
  props: {
    data: {
      type: Object,
      default: () => {}
    },
    showHeader: {
      type: Boolean,
      default: false
    },
    operation: {
      type: String,
      default: ''
    }
  },
  data () {
    return {}
  },
  computed: {
    getDiffTime () {
      const newDate = dayjs(this.data.handleTime)
      const oldDate = dayjs(this.data.createTime)
      let subtime = (newDate - oldDate) / 1000
      let days = parseInt(subtime / 86400)
      let hours = parseInt(subtime / 3600) - 24 * days
      let mins = parseInt((subtime % 3600) / 60)
      let secs = parseInt(subtime % 60)
      return (
        (days > 0 ? days + '天' : '') +
        (hours > 0 ? hours + '小时' : '') +
        (mins > 0 ? mins + '分钟' : '') +
        secs +
        '秒'
      )
    }
  },
  methods: {}
}
</script>

<style lang="scss" scoped>
.header-tag {
  width: calc(100% - 10px);
  margin-left: 10px;
  .title {
    background: #f7f7f7;
    text-align: left;
    margin-top: -30px;
    height: 30px;
    line-height: 30px;
    padding-left: 20px;
    font-size: 12px;
  }
  .content {
    text-align: left;
    padding-left: 20px;
    background: #f0faff;
    font-size: 12px;
  }
}
</style>
