<template>
  <div :style="{ 'max-height': MODALHEIGHT + 'px', overflow: 'auto' }">
    <template v-for="tempG in templateGroup">
      <Card :key="tempG.groupId">
        <p slot="title" @click="tempG.isShow = !tempG.isShow">
          {{ tempG.groupName }}
          <template v-if="tempG.isShow">
            <Icon size="20" type="md-arrow-dropdown" />
          </template>
          <template v-else>
            <Icon size="20" type="md-arrow-dropup" />
          </template>
        </p>
        <template v-for="(tag, tagIndex) in tempG.tags">
          <Card :key="tagIndex" v-if="tempG.isShow">
            <p slot="title">
              {{ tag.tag || $t('unclassified') }}
            </p>
            <div
              @click="choiceTemplate(temp)"
              :class="['diy-tag', temp.status === 'created' ? 'red-style' : '']"
              class="diy-tag"
              v-for="temp in tag.templates"
              :key="temp.id"
            >
              {{ temp.name }}
            </div>
          </Card>
        </template>
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
      MODALHEIGHT: 200,
      templateGroup: []
    }
  },
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 400
    this.getTemplate()
  },
  methods: {
    async getTemplate () {
      const { statusCode, data } = await getTemplateByUser()
      if (statusCode === 'OK') {
        this.templateGroup = data.map(d => {
          d.isShow = true
          return d
        })
      }
    },
    choiceTemplate (temp) {
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
  width: auto;
}
.red-style {
  border: 1px solid red;
  color: red;
}
</style>
