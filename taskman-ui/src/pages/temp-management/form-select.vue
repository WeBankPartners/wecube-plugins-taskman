<template>
  <div class=" ">
    <Form :label-width="100">
      <FormItem>
        <Select v-model="attrData" :label="$t(item.label)">
          <Option v-for="item in attrOptions" :value="item.id" :key="item.id">{{ item.description }}</Option>
        </Select>
      </FormItem>
    </Form>
  </div>
</template>

<script>
import { getFromList } from '@/api/server.js'
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
