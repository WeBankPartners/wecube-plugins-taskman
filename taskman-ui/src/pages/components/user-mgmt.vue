<template>
  <div>
    <Modal v-model="showModal" :mask-closable="false" :closable="false" :width="800" :title="$t('tw_user_mgmt')">
      <div>
        <Tabs type="card">
          <TabPane :label="$t('tw_pending')">标签一的内容</TabPane>
          <TabPane :label="$t('tw_hasProcessed')">标签二的内容</TabPane>
        </Tabs>
      </div>
      <div slot="footer">
        <Button @click="showModal = fasle">{{ $t('cancel') }}</Button>
      </div>
    </Modal>
  </div>
</template>
<script>
export default {
  data () {
    return {
      showModal: false
    }
  },
  methods: {
    openModal () {
      this.showModal = true
    }
  }
}
</script>
