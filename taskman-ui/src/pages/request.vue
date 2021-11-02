<template>
  <div>
    <div>
      <Row>
        <Col span="4">
          <Button @click="addTemplate" type="success">{{ $t('initiate_request') }}</Button>
        </Col>
      </Row>
    </div>
    <Tabs :value="activeTab">
      <TabPane :label="$t('my_drafts')" name="my_drafts">
        <TDraft></TDraft>
      </TabPane>
      <TabPane :label="$t('group_initiated')" name="group_initiated">
        <TGroupInitiated></TGroupInitiated>
      </TabPane>
      <TabPane :label="$t('group_handle')" name="group_handle">
        <TGroupHandle></TGroupHandle>
      </TabPane>
    </Tabs>
  </div>
</template>
<script>
import TDraft from './request-management/t-draft'
import TGroupHandle from './request-management/t-group-handle'
import TGroupInitiated from './request-management/t-group-initiated'
export default {
  name: '',
  data () {
    return {
      activeTab: 'my_drafts'
    }
  },
  mounted () {
    this.activeTab = this.$route.query.activeTab || 'my_drafts'
  },
  methods: {
    addTemplate () {
      this.$router.push({
        path: '/requestManagementIndex',
        query: { requestId: '', requestTemplate: '', isAdd: 'Y', isCheck: 'N', isHandle: 'N', jumpFrom: '' }
      })
    }
  },
  components: {
    TDraft,
    TGroupHandle,
    TGroupInitiated
  }
}
</script>

<style scoped lang="scss">
.header-icon {
  float: right;
  margin: 3px 40px 0 0 !important;
}
</style>
