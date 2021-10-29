<template>
  <div style="width:40%;margin: 0 auto;min-width: 700px">
    <Form :label-width="100">
      <template v-for="item in formConfig.itemConfigs">
        <FormItem v-if="['text', 'password'].includes(item.type)" :label="$t(item.label)" :key="item.value">
          <Input
            v-model="formConfig.values[item.value]"
            style="width:90%"
            :type="item.type"
            :disabled="$parent.formDisable || $parent.jumpFrom === 'group_handle'"
            :placeholder="item.placeholder"
          >
          </Input>
          <Icon v-if="item.rules" size="10" style="color:#ed4014" type="ios-medical" />
        </FormItem>
        <FormItem v-if="['select'].includes(item.type)" :label="$t(item.label)" :key="item.value">
          <Select
            v-model="formConfig.values[item.value]"
            filterable
            style="width:90%"
            :disabled="$parent.formDisable || $parent.jumpFrom === 'group_handle'"
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
      </template>
      <FormItem :label="$t('expected_completion_time')">
        <DatePicker
          type="date"
          :value="formConfig.values.expectTime"
          placeholder="Select date"
          style="width: 200px"
        ></DatePicker>
      </FormItem>

      <FormItem>
        <Button @click="createRequest" type="primary">{{ $t('next') }}</Button>
      </FormItem>
    </Form>
  </div>
</template>

<script>
import { createRequest, updateRequest, getRequestInfo } from '@/api/server'
export default {
  name: 'BasicInfo',
  data () {
    return {
      formConfig: {
        isAdd: true,
        itemConfigs: [
          { label: 'name', value: 'name', rules: 'required', type: 'text' },
          {
            label: this.$t('emergency'),
            value: 'emergency',
            rules: 'required',
            options: 'emergencyOptions',
            labelKey: 'label',
            valueKey: 'value',
            multiple: false,
            type: 'select',
            placeholder: ''
          }
          // {
          //   label: this.$t('expected_completion_time'),
          //   value: 'expireDay',
          //   rules: 'required',
          //   options: 'emergencyOptions',
          //   labelKey: 'label',
          //   valueKey: 'value',
          //   multiple: false,
          //   type: 'select',
          //   placeholder: ''
          // }
        ],
        values: {
          id: '',
          name: '',
          emergency: 3,
          requestTemplate: '',
          expectTime: ''
        },
        emergencyOptions: [
          { label: '1', value: 1 },
          { label: '2', value: 2 },
          { label: '3', value: 3 },
          { label: '4', value: 4 },
          { label: '5', value: 5 }
        ]
      }
    }
  },
  mounted () {
    if (this.$parent.requestId) {
      this.getRequestInfo()
    }
  },
  methods: {
    async getRequestInfo () {
      const { statusCode, data } = await getRequestInfo(this.$parent.requestId)
      if (statusCode === 'OK') {
        this.$parent.procDefId = ''
        this.$parent.procDefKey = ''
        this.formConfig.values.name = data.name
        this.formConfig.values.emergency = data.emergency
        this.formConfig.values.id = data.id
      }
    },
    async createRequest () {
      if (this.$parent.formDisable || this.$parent.jumpFrom === 'group_handle') {
        this.$emit('basicForm', this.formConfig.values.id)
        return
      }
      this.formConfig.values.requestTemplate = this.$parent.requestTemplate
      const method = this.formConfig.values.id ? updateRequest : createRequest
      const { statusCode, data } = await method(this.$parent.requestId, this.formConfig.values)
      if (statusCode === 'OK') {
        this.$emit('basicForm', data.id)
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
      }
    }
  },
  components: {}
}
</script>

<style scoped lang="scss"></style>
