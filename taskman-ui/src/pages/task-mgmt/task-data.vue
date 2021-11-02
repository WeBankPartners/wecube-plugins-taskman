<template>
  <div class="table-c">
    <table width="100%" border="0" cellspacing="0" cellpadding="0">
      <tr>
        <td width="5%" class="padding-style" style="text-align: center">{{ $t('index') }}</td>
        <td width="95%" class="padding-style" style="text-align: center">{{ $t('form') }}</td>
      </tr>
      <template v-for="(data, dataIndex) in tableData">
        <tr :key="data.id">
          <td class="padding-style" style="text-align: center">{{ dataIndex + 1 }}</td>
          <td class="padding-style">
            <div
              class="list-group-item-"
              :style="{ width: (element.width / 24) * 100 + '%' }"
              v-for="element in form"
              :key="element.id"
            >
              <div>
                <Icon v-if="element.required === 'yes'" size="8" style="color:#ed4014" type="ios-medical" />
                {{ element.title }}:
              </div>
              <Input
                v-if="element.elementType === 'input'"
                v-model="data[element.name]"
                placeholder=""
                :disabled="element.isEdit === 'no' || isDisabled || enforceDisable"
              />
              <Input
                v-if="element.elementType === 'textarea'"
                v-model="data[element.name]"
                type="textarea"
                :disabled="element.isEdit === 'no' || isDisabled || enforceDisable"
              />
              <Select
                v-if="element.elementType === 'select'"
                v-model="data[element.name]"
                :disabled="element.isEdit === 'no' || isDisabled || enforceDisable"
                @on-open-change="getRefOptions(element, data, dataIndex)"
              >
                <Option v-for="item in data[element.name + 'Options']" :value="item.guid" :key="item.guid">{{
                  item.key_name
                }}</Option>
              </Select>
            </div>
          </td>
        </tr>
      </template>
    </table>
  </div>
</template>

<script>
import { getRefOptions } from '@/api/server'
export default {
  name: '',
  data () {
    return {
      oriData: null,
      form: {},
      tableData: [],
      refKeys: [] // 引用类型字段集合
    }
  },
  props: ['data', 'isDisabled', 'enforceDisable', 'requestId'],
  mounted () {
    this.initData(this.data)
  },
  methods: {
    async getRefOptions (formItem, formData, index) {
      if (formItem.refEntity === '') {
        formData[formItem.name + 'Options'] = formItem.selectList
        this.$set(this.tableData, index, formData)
        return
      }
      let cache = JSON.parse(JSON.stringify(formData))
      cache[formItem.name] = ''
      this.refKeys.forEach(k => {
        delete cache[k + 'Options']
      })
      const attr = formItem.entity + '__' + formItem.name
      const params = {
        filters: [],
        paging: false,
        dialect: {
          associatedData: {
            ...cache
          }
        }
      }
      const { statusCode, data } = await getRefOptions(this.requestId, attr, params)
      if (statusCode === 'OK') {
        formData[formItem.name + 'Options'] = data
        this.$set(this.tableData, index, formData)
      }
    },
    async initData (data) {
      this.oriData = data
      this.form = data.title
      this.refKeys = []
      this.tableColumns = data.title.map(t => {
        if (t.elementType === 'select') {
          this.refKeys.push(t.name)
        }
        return t
      })
      this.tableData = await data.value.map(v => {
        this.refKeys.forEach(async rfk => {
          v.entityData[rfk + 'Options'] = []
        })
        v.entityData._id = v.id
        return v.entityData
      })
      this.tableData.forEach((data, index) => {
        this.refKeys.forEach(async rfk => {
          const formItem = this.form.find(f => f.name === rfk)
          this.getRefOptions(formItem, data, index)
        })
      })
    }
  },
  components: {}
}
</script>

<style scoped lang="scss">
.list-group-item- {
  display: inline-block;
  margin: 2px 0;
}
.table-c table {
  border-right: 1px solid #dcdee2;
  border-bottom: 1px solid #dcdee2;
}
.table-c table td {
  border-left: 1px solid #dcdee2;
  border-top: 1px solid #dcdee2;
}
.padding-style {
  padding: 2px;
}
</style>
