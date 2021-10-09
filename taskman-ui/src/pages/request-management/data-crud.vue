<template>
  <div class=" ">
    <Form :label-width="100">
      <FormItem label="目标对象">
        <Select v-model="entityDataId" style="width:200px">
          <Option v-for="item in rootEntityOptions" :value="item.id" :key="item.id">{{ item.displayName }}</Option>
        </Select>
        <Button @click="getEntityData" type="primary">{{ $t('next') }}</Button>
      </FormItem>
    </Form>
    <Tabs :value="activeTab" @on-click="changeTab">
      <template v-for="entity in requestData">
        <TabPane :label="entity.entity" :name="entity.entity" :key="entity.entity">
          <DataMgmt ref="dataMgmt"></DataMgmt>
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
      entityDataId: '',
      rootEntityOptions: [],
      activeTab: '',
      requestData: []
    }
  },
  components: {
    DataMgmt
  },
  mounted () {
    this.getEntity()
  },
  methods: {
    async getEntity () {
      let params = {
        params: {
          requestId: this.$parent.requestId
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
          entityDataId: this.entityDataId
        }
      }
      const { statusCode, data } = await getEntityData(params)
      if (statusCode === 'OK') {
        this.activeTab = data[0].entity
        this.requestData = data
        this.$nextTick(() => {
          this.initTable(0)
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
      this.$refs.dataMgmt[index].initData(find)
    }
  }
}
</script>

<style scoped lang="scss"></style>
