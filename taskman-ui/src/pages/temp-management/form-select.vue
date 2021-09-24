<template>
  <div style="width:40%;margin: 0 auto;">
    <Form :label-width="100">
      <FormItem :label="$t('attributes')">
        <Select v-model="attrData" multiple filterable>
          <OptionGroup v-for="item in attrOptions" :label="item.description" :key="item.id">
            <Option v-for="attr in item.attributes" :value="attr.id" :key="attr.id">{{ attr.description }}</Option>
          </OptionGroup>
        </Select>
      </FormItem>
      <FormItem>
        <Button @click="saveAttrs" :disabled="attrData.length === 0" type="primary">{{ $t('next') }}</Button>
      </FormItem>
    </Form>
  </div>
</template>

<script>
import { getFormList, saveAttrs, getRequestTemplateAttrs } from '@/api/server.js'
export default {
  name: 'form-select',
  data () {
    return {
      attrData: [],
      attrOptions: []
    }
  },
  props: ['requestTemplateId'],
  mounted () {
    if (!!this.$parent.requestTemplateId !== false) {
      this.getFormList()
      this.getInitData()
    }
  },
  methods: {
    async getInitData () {
      const { statusCode, data } = await getRequestTemplateAttrs(this.$parent.requestTemplateId)
      if (statusCode === 'OK') {
        this.attrData = data.map(d => d.id)
      }
    },
    async saveAttrs () {
      const attrs = [].concat(...this.attrOptions.map(attr => attr.attributes))
      const params = attrs.filter(item => this.attrData.includes(item.id))
      const { statusCode } = await saveAttrs(this.$parent.requestTemplateId, params)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.$emit('formSelectNextStep')
      }
    },
    async getFormList () {
      const { statusCode, data } = await getFormList(this.$parent.requestTemplateId)
      if (statusCode === 'OK') {
        this.attrOptions = data
      }
    }
  },
  components: {}
}
</script>

<style scoped lang="scss"></style>
