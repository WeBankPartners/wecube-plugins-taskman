<template>
  <div style="width:40%;margin: 0 auto;min-width:700px">
    <!-- <ValidationObserver ref="observer"> -->
    <Form :label-width="100">
      <template v-for="(item, itemIndex) in formConfig.itemConfigs">
        <!-- <ValidationProvider :rules="item.rules" :name="item.value" v-slot="{ errors }" :key="item.value"> -->
        <FormItem v-if="['text', 'password'].includes(item.type)" :label="$t(item.label)" :key="item.value">
          <Input
            v-model="formConfig.values[item.value]"
            style="width:90%"
            :type="item.type"
            :disabled="$parent.isCheck === 'Y'"
            :placeholder="item.placeholder"
          >
          </Input>
          <Icon v-if="item.rules" size="10" style="color:#ed4014" type="ios-medical" />
        </FormItem>
        <FormItem v-if="['textarea'].includes(item.type)" :label="$t(item.label)" :key="item.value">
          <Input
            v-model="formConfig.values[item.value]"
            style="width:90%"
            :type="item.type"
            :disabled="$parent.isCheck === 'Y'"
            :rows="item.rows"
            :placeholder="item.placeholder"
          >
          </Input>
          <Icon v-if="item.rules" size="10" style="color:#ed4014" type="ios-medical" />
        </FormItem>
        <FormItem v-if="['number'].includes(item.type)" :label="$t(item.label)" :key="item.value">
          <InputNumber
            :max="item.max || 1000"
            :min="item.min || 1"
            :disabled="$parent.isCheck === 'Y'"
            v-model="formConfig.values[item.value]"
          ></InputNumber>
          <Icon v-if="item.rules" size="10" style="color:#ed4014" type="ios-medical" />
        </FormItem>
        <FormItem v-if="['select'].includes(item.type)" :label="$t(item.label)" :key="item.value">
          <Select
            v-model="formConfig.values[item.value]"
            clearable
            filterable
            :disabled="$parent.isCheck === 'Y'"
            @on-open-change="execut(item.onOpenChange)"
            style="width:90%"
            :multiple="item.multiple"
            :placeholder="item.placeholder"
          >
            <template v-for="option in getOptions(item.options)">
              <Option
                :label="item.labelKey ? option[item.labelKey] : option.label"
                :value="item.valueKey ? option[item.valueKey] : option.value"
                :key="item.valueKey ? option[item.valueKey] : option.value"
              >
              </Option>
            </template>
          </Select>
          <Icon v-if="item.rules" size="10" style="color:#ed4014" type="ios-medical" />
        </FormItem>
        <FormItem v-if="['create_select'].includes(item.type)" :label="$t(item.label)" :key="itemIndex">
          <Select
            v-model="formConfig.values[item.value]"
            @on-open-change="execut(item.onOpenChange)"
            filterable
            allow-create
            :disabled="$parent.isCheck === 'Y'"
            style="width:90%"
            @on-create="handleCreate1"
          >
            <Option v-for="(item, tagIndex) in formConfig[item.options]" :value="item.value" :key="tagIndex">{{
              item.label
            }}</Option>
          </Select>
        </FormItem>
        <!-- </ValidationProvider> -->
      </template>
    </Form>
    <div style="text-align: center">
      <Button @click="resetParams" :disabled="$parent.isCheck === 'Y'">{{ $t('reset') }}</Button>
      <Button @click="createTemp" type="primary">{{ $t('next') }}</Button>
    </div>
    <!-- </ValidationObserver> -->
  </div>
</template>

<script>
import { ValidationObserver } from 'vee-validate'
import {
  getTempGroupList,
  getManagementRoles,
  getTemplateTags,
  getHandlerRoles,
  getUserRoles,
  getProcess,
  createTemp,
  updateTemp,
  getTemplateList
} from '@/api/server'
export default {
  name: 'BasicInfo',
  data () {
    return {
      formConfig: {
        isAdd: true,
        itemConfigs: [
          { label: 'name', value: 'name', rules: 'required', type: 'text' },
          {
            label: 'group',
            value: 'group',
            rules: 'required',
            onOpenChange: 'getGroupOptions',
            options: 'groupOptions',
            labelKey: 'name',
            valueKey: 'id',
            multiple: false,
            type: 'select',
            placeholder: ''
          },
          {
            label: 'procDefId',
            value: 'procDefId',
            rules: 'required',
            onOpenChange: 'getProcess',
            options: 'procOptions',
            labelKey: 'procDefName',
            valueKey: 'procDefId',
            multiple: false,
            type: 'select',
            placeholder: ''
          },
          {
            label: 'mgmtRoles',
            value: 'mgmtRoles',
            rules: 'required',
            onOpenChange: 'getManagementRoles',
            options: 'mgmtRolesOptions',
            labelKey: 'displayName',
            valueKey: 'id',
            multiple: false,
            type: 'select',
            placeholder: ''
          },
          {
            label: 'handler',
            value: 'handler',
            rules: '',
            onOpenChange: 'getHandlerRoles',
            options: 'handlerRolesOptions',
            labelKey: 'displayName',
            valueKey: 'id',
            multiple: false,
            type: 'select',
            placeholder: ''
          },
          {
            label: 'useRoles',
            value: 'useRoles',
            rules: 'required',
            onOpenChange: 'getUserRoles',
            options: 'useRolesOptions',
            labelKey: 'displayName',
            valueKey: 'id',
            multiple: true,
            type: 'select',
            placeholder: ''
          },
          {
            label: 'tags',
            value: 'tags',
            onOpenChange: 'getTags',
            options: 'tagOptions',
            labelKey: 'label',
            valueKey: 'value',
            type: 'create_select'
          },
          { label: 'description', value: 'description', rows: 2, type: 'textarea' }
        ],
        values: {
          id: '',
          name: '',
          group: '',
          mgmtRoles: '',
          handler: '',
          useRoles: [],
          description: '',
          tags: '',
          packageName: '',
          entityName: '',
          procDefId: '',
          procDefKey: '',
          procDefName: ''
        },
        handlerRolesOptions: [],
        groupOptions: [],
        procOptions: [],
        mgmtRolesOptions: [],
        useRolesOptions: [],
        tagOptions: [],
        tmpTagOptions: []
      }
    }
  },
  mounted () {
    this.getInitData()
  },
  methods: {
    handleCreate1 (v) {
      this.formConfig.tmpTagOptions.push(v)
    },
    async getTags (val) {
      if (this.formConfig.values.group === '') {
        // this.$Notice.warning({
        //   title: this.$t('warning'),
        //   desc: '请选选择模板组',
        // })
        return
      }
      const { statusCode, data } = await getTemplateTags(this.formConfig.values.group)
      if (statusCode === 'OK') {
        const totalData = this.unique(this.formConfig.tmpTagOptions.concat(data))
        this.formConfig.tagOptions = totalData.map(d => {
          return {
            label: d,
            value: d
          }
        })
      }
    },
    unique (arr) {
      return Array.from(new Set(arr))
    },
    getOptions (options) {
      return this.formConfig[options]
    },
    execut (method) {
      this[method]()
    },
    async getHandlerRoles () {
      const params = {
        params: {
          roles: this.formConfig.values.mgmtRoles
        }
      }
      const { statusCode, data } = await getHandlerRoles(params)
      if (statusCode === 'OK') {
        this.formConfig.handlerRolesOptions = data.map(d => {
          return {
            displayName: d,
            id: d
          }
        })
      }
    },
    async getInitData () {
      this.getGroupOptions()
      this.getManagementRoles()
      this.getProcess()
      this.getUserRoles()
      if (!!this.$parent.requestTemplateId === false) {
        return
      }
      const params = {
        filters: [
          {
            name: 'id',
            operator: 'eq',
            value: this.$parent.requestTemplateId
          }
        ],
        paging: false
      }
      const { statusCode, data } = await getTemplateList(params)
      if (statusCode === 'OK') {
        this.formConfig.isAdd = false
        this.formConfig.values = { ...data.contents[0] }
        this.formConfig.values.mgmtRoles = data.contents[0].mgmtRoles[0].id
        this.formConfig.values.useRoles = data.contents[0].useRoles.map(role => role.id)
        this.getHandlerRoles()
        this.getTags()
      }
    },
    async createTemp () {
      // if (!this.$refs.observer.flags.valid) {
      //   return
      // }
      if (this.$parent.isCheck === 'Y') {
        this.$emit('basicInfoNextStep', this.formConfig.values)
        return
      }
      let cacheFromValue = JSON.parse(JSON.stringify(this.formConfig.values))
      const method = cacheFromValue.id === '' ? createTemp : updateTemp
      const process = this.formConfig.procOptions.find(item => item.procDefId === this.formConfig.values.procDefId)
      cacheFromValue.packageName = process.rootEntity.packageName
      cacheFromValue.entityName = process.rootEntity.name
      cacheFromValue.procDefKey = process.procDefKey
      cacheFromValue.procDefName = process.procDefName
      cacheFromValue.mgmtRoles = [cacheFromValue.mgmtRoles]
      const { statusCode, data } = await method(cacheFromValue)
      if (statusCode === 'OK') {
        this.formConfig.values = { ...data }
        this.formConfig.values.mgmtRoles = this.formConfig.values.mgmtRoles[0]
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.$emit('basicInfoNextStep', data)
      }
    },
    resetParams () {
      const keys = Object.keys(this.formConfig.values)
      keys.forEach(k => {
        let value = this.formConfig.values[k]
        if (typeof value === 'string') {
          this.formConfig.values[k] = ''
        }
        if (typeof value === 'object' && Array.isArray(value)) {
          this.formConfig.values[k] = []
        }
      })
    },
    async getProcess () {
      const { statusCode, data } = await getProcess()
      if (statusCode === 'OK') {
        this.formConfig.procOptions = data
      }
    },
    async getManagementRoles () {
      const { statusCode, data } = await getManagementRoles()
      if (statusCode === 'OK') {
        this.formConfig.mgmtRolesOptions = data
      }
    },
    async getUserRoles () {
      const { statusCode, data } = await getUserRoles()
      if (statusCode === 'OK') {
        this.formConfig.useRolesOptions = data
      }
    },
    async getGroupOptions () {
      const params = {
        filters: [],
        paging: false
      }
      const { statusCode, data } = await getTempGroupList(params)
      if (statusCode === 'OK') {
        this.formConfig.groupOptions = data.contents
      }
    }
  },
  components: {
    ValidationObserver
  }
}
</script>

<style scoped lang="scss"></style>
