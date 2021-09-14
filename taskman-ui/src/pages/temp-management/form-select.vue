<template>
  <div style="width:40%;margin: 0 auto;">
    <Form :label-width="100">
      <FormItem :label="$t('attributes')">
        <Select v-model="attrData" multiple filterable>
          <Option v-for="item in attrOptions" :value="item.id" :key="item.id">{{ item.description }}</Option>
        </Select>
      </FormItem>
      <FormItem>
        <Button @click="saveAttrs" :disabled="attrData.length === 0" type="primary">{{ $t('next') }}</Button>
      </FormItem>
    </Form>
  </div>
</template>

<script>
import { getFromList, saveAttrs } from '@/api/server.js'
export default {
  name: 'form-select',
  data () {
    return {
      requestTemplateId: '614043ac9379fb1e',
      attrData: [],
      attrOptions: []
    }
  },
  // props: ['requestTemplateId'],
  mounted () {
    this.getFromList()
  },
  methods: {
    async saveAttrs () {
      const params = this.attrOptions.filter(item => this.attrData.includes(item.id))
      const { statusCode } = await saveAttrs(this.requestTemplateId, params)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.$emit('formSelectNextStep')
      }
    },
    async getFromList () {
      const { statusCode, data } = await getFromList(this.requestTemplateId)
      if (statusCode === 'OK') {
        this.attrOptions = data
      }
    }
  },
  components: {}
}
</script>

<style scoped lang="scss"></style>
