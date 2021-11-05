<template>
  <div class=" ">
    <Form :label-width="100">
      <FormItem :label="$t('root_entity')">
        <Select
          v-model="rootEntityId"
          filterable
          clearable
          :disabled="$parent.formDisable || $parent.jumpFrom === 'group_handle'"
          style="width:300px"
        >
          <Option v-for="item in rootEntityOptions" :value="item.guid" :key="item.guid">{{ item.key_name }}</Option>
        </Select>
        <Button
          @click="getEntityData"
          :disabled="$parent.formDisable || $parent.jumpFrom === 'group_handle'"
          type="primary"
          >{{ $t('search') }}</Button
        >
      </FormItem>
    </Form>
    <Tabs :value="activeTab" @on-click="changeTab">
      <template v-for="entity in requestData">
        <TabPane :label="entity.entity" :name="entity.entity" :key="entity.entity">
          <DataMgmt ref="dataMgmt" @getEntityData="getEntityData" @backData="backData"></DataMgmt>
        </TabPane>
      </template>
    </Tabs>
    <div style="text-align: center;margin-top:24px">
      <Button @click="saveData" v-if="!($parent.formDisable || $parent.jumpFrom === 'group_handle')" type="primary">{{
        $t('save')
      }}</Button>
      <Button @click="commitRequest" v-if="!($parent.formDisable || $parent.jumpFrom === 'group_handle')">{{
        $t('commit')
      }}</Button>
      <Button @click="nextStep" v-if="['', 'group_handle'].includes($parent.jumpFrom) && !$parent.isAdd">{{
        $t('next')
      }}</Button>
    </div>
  </div>
</template>

<script>
import DataMgmt from './data-mgmt'
import { getRootEntity, getEntityData, saveEntityData, getRequestInfo, updateRequestStatus } from '@/api/server'
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
    backData (data) {
      this.requestData = data
    },
    async commitRequest () {
      this.$Modal.confirm({
        title: this.$t('confirm') + this.$t('commit'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          await this.saveData()
          const { statusCode } = await updateRequestStatus(this.$parent.requestId, 'Pending')
          if (statusCode === 'OK') {
            this.$router.push({ path: '/taskman/request-mgmt' })
          }
        },
        onCancel: () => {}
      })
    },
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
      const result = this.paramsCheck()
      if (result) {
        const { statusCode, data } = await saveEntityData(this.requestId, params)
        if (statusCode === 'OK') {
          this.requestData = data.data
          this.$Notice.success({
            title: this.$t('successful'),
            desc: this.$t('successful')
          })
        }
      } else {
        this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('required_tip')
        })
      }
    },
    paramsCheck () {
      let result = true
      this.requestData.forEach(requestData => {
        let requiredName = []
        requestData.title.forEach(t => {
          if (t.required === 'yes') {
            requiredName.push(t.name)
          }
        })
        requestData.value.forEach(v => {
          requiredName.forEach(key => {
            let val = v.entityData[key]
            if (Array.isArray(val)) {
              if (val.length === 0) {
                result = false
              }
            } else {
              if (val === '') {
                result = false
              }
            }
          })
        })
      })
      return result
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
        this.$parent.formDisable,
        this.$parent.jumpFrom
      )
    }
  }
}
</script>

<style scoped lang="scss"></style>
