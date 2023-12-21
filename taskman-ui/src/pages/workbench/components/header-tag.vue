<template>
  <div class="header-tag">
    <Row v-if="showHeader" class="title" :gutter="10">
      <Col :span="4">{{ $t('handle_time') }}</Col>
      <Col :span="3">{{ $t('tw_assume') }}</Col>
      <Col :span="3">{{ $t('handler_role') }}</Col>
      <Col :span="3">{{ $t('handler') }}</Col>
      <Col :span="3">{{ $t('t_action') }}</Col>
      <Col :span="4">{{ $t('tw_note') }}</Col>
      <Col :span="4">{{ $t('tw_attach') }}</Col>
    </Row>
    <Row class="content" :gutter="10">
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
        (days > 0 ? days + this.$t('tw_day') : '') +
        (hours > 0 ? hours + this.$t('tw_hour') : '') +
        (mins > 0 ? mins + this.$t('tw_minute') : '') +
        secs +
        this.$t('tw_second')
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
    margin-top: -36px;
    height: 36px;
    line-height: 36px;
    padding-left: 20px;
    font-size: 12px;
    font-weight: bold;
  }
  .content {
    text-align: left;
    padding-left: 20px;
    background: #f0faff;
    font-size: 12px;
    overflow: hidden;
    word-break: break-all;
  }
}
</style>
