<template>
  <div class=" ">
    <Form :label-width="100">
      <FormItem :label="$t('root_entity')">
        <Select v-model="rootEntityId" :disabled="$parent.formDisable" style="width:300px">
          <Option v-for="item in rootEntityOptions" :value="item.guid" :key="item.guid">{{ item.key_name }}</Option>
        </Select>
        <Button @click="getEntityData" :disabled="$parent.formDisable" type="primary">{{ $t('search') }}</Button>
      </FormItem>
    </Form>
    <Tabs :value="activeTab" @on-click="changeTab">
      <template v-for="entity in requestData">
        <TabPane :label="entity.entity" :name="entity.entity" :key="entity.entity">
          <DataMgmt ref="dataMgmt" @getEntityData="getEntityData"></DataMgmt>
        </TabPane>
      </template>
    </Tabs>
    <div style="text-align: center;margin-top:48px">
      <Button @click="saveData" :disabled="$parent.formDisable" type="primary">{{ $t('save') }}</Button>
      <Button @click="nextStep">{{ $t('next') }}</Button>
    </div>
  </div>
</template>

<script>
import DataMgmt from './data-mgmt'
import { getRootEntity, getEntityData, saveEntityData, getRequestInfo } from '@/api/server'
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
    if (this.$parent.requestId) {
      this.getRequestInfo()
    }
  },
  methods: {
    async getRequestInfo () {
      const { statusCode, data } = await getRequestInfo(this.requestId)
      if (statusCode === 'OK') {
        this.rootEntityId = data.cache
        this.getEntityData()
      }
    },
    nextStep () {
      if (!this.$parent.formDisable) {
        this.saveData()
      }
      this.$emit('nextStep')
    },
    async saveData () {
      const params = {
        rootEntityId: this.rootEntityId,
        data: this.requestData
      }
      const { statusCode } = await saveEntityData(this.requestId, params)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
      }
    },
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
      this.$refs.dataMgmt[index].initData(
        this.rootEntityId,
        this.requestData,
        find,
        this.requestId,
        this.$parent.formDisable
      )
    }
  }
}
</script>

<style scoped lang="scss"></style>
