<template>
  <div style="margin: 24px">
    <div>
      <Row>
        <Col span="4">
          <Button @click="addTemp">{{ $t('add') }}</Button>
        </Col>
      </Row>
    </div>
  </div>
</template>

<script>
import {} from '@/api/server'
export default {
  name: '',
  data () {
    return {}
  },
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 200
  },
  methods: {
    addTemp () {
      this.$router.push({ path: '/templateManagementIndex' })
    }
  },
  components: {}
}
</script>

<style scoped lang="scss">
.header-icon {
  float: right;
  margin: 3px 40px 0 0 !important;
}
</style>
