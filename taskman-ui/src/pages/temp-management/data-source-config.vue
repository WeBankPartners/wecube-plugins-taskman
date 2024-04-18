<template>
  <div>
    <Modal v-model="showModel" :title="$t('data_set')" :mask-closable="false" :closable="false">
      <div>
        <Row>
          <Col span="12">{{ $t('display_name') }}</Col>
          <Col span="10">{{ $t('value') }}</Col>
        </Row>
        <Row v-for="(item, itemIndex) in dataSource" :key="itemIndex" style="margin:6px 0">
          <Col span="12">
            <Input v-model.trim="item.label" style="width:90%"></Input>
          </Col>
          <Col span="10">
            <Input v-model.trim="item.value" style="width:90%"></Input>
          </Col>
          <Col span="2">
            <Button
              type="error"
              ghost
              @click="deleteItem(itemIndex)"
              size="small"
              style="vertical-align: sub;cursor: pointer"
              icon="md-trash"
            ></Button>
          </Col>
        </Row>
        <div style="text-align: right;margin-right: 16px;cursor: pointer">
          <Button type="primary" ghost @click="addItem" size="small" icon="md-add"></Button>
        </div>
      </div>
      <template #footer>
        <Button @click="showModel = false">{{ $t('cancel') }}</Button>
        <Button @click="okSelect" type="primary">{{ $t('confirm') }}</Button>
      </template>
    </Modal>
  </div>
</template>

<script>
export default {
  name: '',
  data () {
    return {
      showModel: false,
      dataSource: []
    }
  },
  methods: {
    loadPage (dataSource) {
      this.dataSource = dataSource
      this.showModel = true
    },
    dataValidateFirst () {
      let res = true
      const infoSet = new Set()
      this.dataSource.forEach(item => {
        if (infoSet.has(item.label)) {
          res = false
        } else {
          infoSet.add(item.label)
        }
      })
      if (!res) {
        this.$Message.warning(this.$t('tw_duplicate_data'))
      }
      return res
    },
    okSelect () {
      const isCanBeSave = this.dataValidateFirst()
      if (!isCanBeSave) {
        return
      }
      this.$emit(
        'setDataOptions',
        this.dataSource.filter(ds => ds.label !== '')
      )
      this.showModel = false
    },
    addItem () {
      this.dataSource.push({
        label: '',
        value: ''
      })
    },
    deleteItem (itemIndex) {
      this.dataSource.splice(itemIndex, 1)
    }
  }
}
</script>
