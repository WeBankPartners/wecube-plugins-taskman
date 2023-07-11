<template>
  <div :style="{ 'max-height': MODALHEIGHT + 'px', overflow: 'auto' }">
    <template v-for="tempG in templateGroup">
      <Card :key="tempG.groupId" style="width:48%;display: inline-block;vertical-align: text-top;margin:1%">
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
          <Card :key="tagIndex" v-if="tempG.isShow" style=";margin:1%">
            <p slot="title">
              {{ tag.tag || $t('unclassified') }}
            </p>
            <div
              @click="choiceTemplate(temp)"
              :class="['diy-tag', temp.status === 'created' ? 'red-style' : '']"
              class="diy-tag"
              v-for="temp in filterData(tag.templates)"
              :key="temp.id"
            >
              <Tooltip content="" max-width="200">
                {{ temp.name }}
                <div slot="content">
                  <p>{{ $t('version') }}:{{ temp.version }}</p>
                  <p>{{ $t('description') }}: {{ temp.description }}</p>
                </div>
              </Tooltip>
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
  props: ['filterWord'],
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 200
    this.getTemplate()
  },
  methods: {
    filterData (templates) {
      const res = templates.filter(t => t.name.toUpperCase().includes(this.filterWord.toUpperCase()))
      return res
    },
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
