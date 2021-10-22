<template>
  <div style="width:55%;margin: 0 auto;">
    <CheckboxGroup v-model="attrData">
      <Tabs :value="activeTab">
        <template v-for="item in attrOptions">
          <TabPane :label="item.description" :name="item.id" :key="item.id">
            <ul>
              <li v-for="attr in item.attributes" :key="attr.id" style="width: 46%;display: inline-block;margin: 4px">
                <Checkbox :label="attr.id">
                  <span>{{ attr.description }}({{ attr.name }})</span>
                </Checkbox>
              </li>
            </ul>
          </TabPane>
        </template>
      </Tabs>
    </CheckboxGroup>
    <div style="margin-top:16px;">
      <Button @click="saveAttrs" :disabled="attrData.length === 0" type="primary">{{ $t('next') }}</Button>
    </div>
  </div>
</template>

<script>
import { getFormList, saveAttrs, getRequestTemplateAttrs } from '@/api/server.js'
export default {
  name: 'form-select',
  data () {
    return {
      activeTab: '',
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
        if (this.attrOptions.length > 0) {
          this.activeTab = this.attrOptions[0].id
        }
      }
    }
  },
  components: {}
}
</script>

<style scoped lang="scss"></style>
