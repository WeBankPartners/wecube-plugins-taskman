<template>
  <div class=" ">
    <template v-for="tempG in templateGroup">
      <Card :key="tempG.groupId">
        <p slot="title">
          {{ tempG.groupName }}
        </p>
        <div @click="choiceTemplate(temp)" class="diy-tag" v-for="temp in tempG.templates" :key="temp.id">
          {{ temp.name }}
        </div>
      </Card>
    </template>
  </div>
</template>

<script>
import { getTemplateByUser } from '@/api/server.js'
export default {
  name: '',
  data () {
    return {
      templateGroup: []
    }
  },
  mounted () {
    this.getTemplate()
  },
  methods: {
    async getTemplate () {
      const { statusCode, data } = await getTemplateByUser()
      if (statusCode === 'OK') {
        this.templateGroup = data
      }
    },
    choiceTemplate (temp) {
      console.log(temp)
      this.$emit('choiceTemp', temp)
    }
  },
  components: {}
}
</script>

<style scoped lang="scss">
.diy-tag {
  display: inline-block;
  cursor: pointer;
  margin: 8px;
  border: 1px solid #338cf0;
  color: #338cf0;
  padding: 2px 8px 0;
  border-radius: 4px;
  font-size: 13px;
}
</style>
