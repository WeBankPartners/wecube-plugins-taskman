<template>
  <div class=" ">
    <Form :label-width="100">
      <FormItem :label="$t('root_entity')">
        <Select v-model="rootEntityId" style="width:300px">
          <Option v-for="item in rootEntityOptions" :value="item.id" :key="item.id">{{ item.displayName }}</Option>
        </Select>
        <Button @click="getEntityData" type="primary">{{ $t('search') }}</Button>
      </FormItem>
    </Form>
    <Tabs :value="activeTab" @on-click="changeTab">
      <template v-for="entity in requestData">
        <TabPane :label="entity.entity" :name="entity.entity" :key="entity.entity">
          <DataMgmt ref="dataMgmt" @getEntityData="getEntityData"></DataMgmt>
        </TabPane>
      </template>
    </Tabs>
  </div>
</template>

<script>
import DataMgmt from './data-mgmt'
import { getRootEntity, getEntityData } from '@/api/server'
export default {
  name: '',
  data () {
    return {
      requestId: '',
      rootEntityId: '',
      rootEntityOptions: [],
      activeTab: '',
      requestData: []
    }
  },
  components: {
    DataMgmt
  },
  mounted () {
    this.requestId = this.$parent.requestId
    this.getEntity()
    this.getEntityData()
  },
  methods: {
    async getEntity () {
      let params = {
        params: {
          requestId: this.requestId
        }
      }
      const { statusCode, data } = await getRootEntity(params)
      if (statusCode === 'OK') {
        this.rootEntityOptions = data.data
      }
    },
    async getEntityData () {
      let params = {
        params: {
          requestId: this.$parent.requestId,
          rootEntityId: this.rootEntityId
        }
      }
      const { statusCode, data } = await getEntityData(params)
      if (statusCode === 'OK') {
        console.log(data)
        this.activeTab = this.activeTab || data.data[0].entity
        this.requestData = data.data
        this.$nextTick(() => {
          const index = this.requestData.findIndex(r => r.entity === this.activeTab)
          this.initTable(index)
        })
      }
    },
    changeTab (entity) {
      this.activeTab = entity
      const index = this.requestData.findIndex(r => r.entity === this.activeTab)
      this.initTable(index)
    },
    initTable (index) {
      const find = this.requestData.find(r => r.entity === this.activeTab)
      this.$refs.dataMgmt[index].initData(this.rootEntityId, this.requestData, find, this.requestId)
    }
  }
}
</script>

<style scoped lang="scss"></style>
