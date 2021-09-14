<template>
  <div style="width:40%;margin: 0 auto;">
    <ValidationObserver ref="observer">
      <Form :label-width="100">
        <template v-for="item in formConfig.itemConfigs">
          <ValidationProvider :rules="item.rules" :name="item.value" v-slot="{ errors }" :key="item.value">
            <FormItem v-if="['text', 'password'].includes(item.type)" :label="$t(item.label)" :error="errors[0]">
              <Input
                v-model="formConfig.values[item.value]"
                style="width:90%"
                :type="item.type"
                :placeholder="item.placeholder"
              >
              </Input>
              <Icon v-if="item.rules" size="10" style="color:#ed4014" type="ios-medical" />
            </FormItem>
            <FormItem v-if="['textarea'].includes(item.type)" :label="$t(item.label)" :error="errors[0]">
              <Input
                v-model="formConfig.values[item.value]"
                style="width:90%"
                :type="item.type"
                :rows="item.rows"
                :placeholder="item.placeholder"
              >
              </Input>
              <Icon v-if="item.rules" size="10" style="color:#ed4014" type="ios-medical" />
            </FormItem>
            <FormItem v-if="['number'].includes(item.type)" :label="$t(item.label)" :error="errors[0]">
              <InputNumber
                :max="item.max || 1000"
                :min="item.min || 1"
                v-model="formConfig.values[item.value]"
              ></InputNumber>
              <Icon v-if="item.rules" size="10" style="color:#ed4014" type="ios-medical" />
            </FormItem>
            <FormItem v-if="['select'].includes(item.type)" :label="$t(item.label)" :error="errors[0]">
              <Select
                v-model="formConfig.values[item.value]"
                clearable
                filterable
                style="width:90%"
                :multiple="item.multiple"
                :placeholder="item.placeholder"
              >
                <template v-for="option in formConfig[item.options]">
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
          </ValidationProvider>
        </template>
        <FormItem>
          <Button @click="resetParams">{{ $t('reset') }}</Button>
          <Button @click="createTemp" type="primary">{{ $t('next') }}</Button>
        </FormItem>
      </Form>
    </ValidationObserver>
  </div>
</template>

<script>
import { ValidationObserver } from 'vee-validate'
import { getTempGroupList, getManagementRoles, getUserRoles, getProcess, createTemp, updateTemp } from '@/api/server'
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
            options: 'mgmtRolesOptions',
            labelKey: 'displayName',
            valueKey: 'id',
            multiple: true,
            type: 'select',
            placeholder: ''
          },
          {
            label: 'useRoles',
            value: 'useRoles',
            rules: 'required',
            options: 'useRolesOptions',
            labelKey: 'displayName',
            valueKey: 'id',
            multiple: true,
            type: 'select',
            placeholder: ''
          },
          { label: 'tags', value: 'tags', type: 'text' },
          { label: 'description', value: 'description', rows: 2, type: 'textarea' }
        ],
        values: {
          id: '',
          name: '',
          group: '',
          mgmtRoles: [],
          useRoles: [],
          description: '',
          tags: '',
          packageName: '',
          entityName: '',
          procDefId: '',
          procDefKey: '',
          procDefName: ''
        },
        groupOptions: [],
        procOptions: [],
        mgmtRolesOptions: [],
        useRolesOptions: []
      }
    }
  },
  mounted () {
    this.getGroupOptions()
    this.getManagementRoles()
    this.getProcess()
    this.getUserRoles()
  },
  methods: {
    async createTemp () {
      if (!this.$refs.observer.flags.valid) {
        return
      }
      const method = this.formConfig.values.id === '' ? createTemp : updateTemp
      const process = this.formConfig.procOptions.find(item => item.procDefId === this.formConfig.values.procDefId)
      this.formConfig.values.packageName = process.rootEntity.packageName
      this.formConfig.values.entityName = process.rootEntity.name
      this.formConfig.values.procDefKey = process.procDefKey
      this.formConfig.values.procDefName = process.procDefName
      const { statusCode, data } = await method(this.formConfig.values)
      if (statusCode === 'OK') {
        this.formConfig.values = { ...data }
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
    },
    getProcOptions () {},
    getMgmtRoles () {},
    getUseRoles () {}
  },
  components: {
    ValidationObserver
  }
}
</script>

<style scoped lang="scss"></style>
