<template>
  <div>
    <Modal
      v-model="showModel"
      :title="$t('data_set')"
      :mask-closable="false"
      :closable="false"
      class="data-source-config"
    >
      <div>
        <Row>
          <Col span="12">{{ $t('display_name') }}</Col>
          <Col span="10">{{ $t('value') }}</Col>
        </Row>
        <Row v-for="(item, itemIndex) in dataSource" :key="itemIndex" style="margin:6px 0">
          <Form :model="item" :rules="rule" :ref="'form' + itemIndex">
            <Col span="12">
              <FormItem label="" prop="label">
                <Input v-model.trim="item.label" style="width:90%"></Input>
              </FormItem>
            </Col>
            <Col span="10">
              <FormItem label="" prop="value">
                <Input v-model.trim="item.value" style="width:90%"></Input>
              </FormItem>
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
          </Form>
        </Row>
        <div style="text-align: right;margin-right: 16px;cursor: pointer">
          <Button type="success" ghost @click="addItem" size="small" icon="md-add"></Button>
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
      dataSource: [],
      rule: {
        label: [{ required: true, message: this.$t('display_name'), trigger: 'change' }],
        value: [{ required: true, message: this.$t('value'), trigger: 'change' }]
      }
    }
  },
  methods: {
    loadPage (dataSource) {
      if (dataSource && dataSource.length === 0) {
        this.dataSource = [{ label: '', value: '' }]
      } else {
        this.dataSource = dataSource
      }
      this.showModel = true
    },
    validateSameLabel () {
      let res = true
      const infoSet = new Set()
      const nameArr = this.dataSource.filter(i => i.label)
      nameArr.forEach(item => {
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
    validateSameValue () {
      let res = true
      const infoSet = new Set()
      const valueArr = this.dataSource.filter(i => i.value)
      valueArr.forEach(item => {
        if (infoSet.has(item.value)) {
          res = false
        } else {
          infoSet.add(item.value)
        }
      })
      if (!res) {
        this.$Message.warning(this.$t('tw_duplicate_value'))
      }
      return res
    },
    okSelect () {
      if (!this.validateSameLabel()) {
        return
      }
      if (!this.validateSameValue()) {
        return
      }
      const validArr = []
      this.dataSource.forEach((_, index) => {
        const key = `form${index}`
        this.$refs[key][0].validate(valid => {
          validArr.push(valid)
        })
      })
      const validFlag = validArr.every(i => i === true)
      if (validFlag) {
        this.$emit(
          'setDataOptions',
          this.dataSource.filter(ds => ds.label !== '')
        )
        this.showModel = false
      }
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
<style lang="scss">
.data-source-config {
  .ivu-form-item {
    margin-bottom: 8px;
  }
  .ivu-form-item-error-tip {
    padding-top: 0px;
  }
}
</style>
