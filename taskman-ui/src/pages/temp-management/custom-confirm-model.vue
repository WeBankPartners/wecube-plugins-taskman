<template>
  <Modal v-model="showConfirm" closable :title="config.title">
    <p>{{ config.content }}</p>
    <div slot="footer">
      <Button type="text" @click="showConfirm = false">{{ $t('cancel') }}</Button>
      <Button type="primary" ghost @click="cancelHandler">{{ config.cancelText }}</Button>
      <Button type="primary" @click="okHandler">{{ config.okText }}</Button>
    </div>
  </Modal>
</template>
<script>
export default {
  data () {
    return {
      showConfirm: false,
      config: {
        title: '',
        content: '',
        okText: '',
        cancelText: '',
        okFunc: () => {},
        cancelFunc: () => {}
      }
    }
  },
  methods: {
    open (config) {
      this.config = config
      this.showConfirm = true
    },
    okHandler () {
      this.config.okFunc()
      this.showConfirm = false
    },
    cancelHandler () {
      this.config.cancelFunc()
      this.showConfirm = false
    }
  }
}
</script>
