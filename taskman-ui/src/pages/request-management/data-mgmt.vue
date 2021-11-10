<template>
  <div class="table-c">
    <Button
      @click="addRow"
      type="success"
      style="margin-bottom: 4px"
      v-if="!(formDisable || jumpFrom === 'group_handle')"
      :disabled="formDisable || jumpFrom === 'group_handle'"
      >{{ $t('add') }}</Button
    >
    <div :style="{ 'max-height': MODALHEIGHT + 'px', overflow: 'auto', width: '100%' }">
      <table width="100%" border="0" cellspacing="0" cellpadding="0">
        <tr>
          <td width="5%" class="padding-style" style="text-align: center">{{ $t('index') }}</td>
          <td width="85%" class="padding-style" style="text-align: center">{{ $t('form') }}</td>
          <td
            width="10%"
            class="padding-style"
            style="text-align: center"
            v-if="!(formDisable || jumpFrom === 'group_handle')"
          >
            {{ $t('t_action') }}
          </td>
        </tr>
        <template v-for="(data, dataIndex) in tableData">
          <tr :key="data.id">
            <td class="padding-style" style="text-align: center">
              <span style="font-size: 18px">
                {{ dataIndex + 1 }}
              </span>
              <span v-if="data._id.startsWith('tmp__')" style="color: #338cf0">new</span>
            </td>
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
                  :disabled="element.isEdit === 'no' || formDisable || jumpFrom === 'group_handle'"
                  style="width: calc(100% - 30px);"
                />
                <Input
                  v-if="element.elementType === 'textarea'"
                  v-model="data[element.name]"
                  type="textarea"
                  :disabled="element.isEdit === 'no' || formDisable || jumpFrom === 'group_handle'"
                  style="width: calc(100% - 30px);"
                />
                <Select
                  v-if="element.elementType === 'select'"
                  v-model="data[element.name]"
                  filterable
                  clearable
                  :multiple="element.multiple === 'Y'"
                  @on-open-change="getRefOptions(element, data, dataIndex)"
                  :disabled="element.isEdit === 'no' || formDisable || jumpFrom === 'group_handle'"
                  style="width: calc(100% - 30px);"
                >
                  <Option v-for="item in data[element.name + 'Options']" :value="item.guid" :key="item.guid">{{
                    item.key_name
                  }}</Option>
                </Select>
              </div>
            </td>
            <td class="padding-style" style="text-align: center" v-if="!(formDisable || jumpFrom === 'group_handle')">
              <Button
                style="margin-left: 4px"
                @click="deleteRow(dataIndex)"
                v-if="!(formDisable || jumpFrom === 'group_handle')"
                :disabled="formDisable || jumpFrom === 'group_handle'"
                size="small"
                type="error"
              >
                {{ $t('delete') }}
              </Button>
            </td>
          </tr>
        </template>
      </table>
    </div>
  </div>
</template>

<script>
import { getRefOptions } from '@/api/server'
export default {
  name: '',
  data () {
    return {
      MODALHEIGHT: 200,
      formDisable: false,
      jumpFrom: '',
      requestId: '',
      rootEntityId: 'host_resource_6152f8039c58e6b0',
      dataArray: [], // 所有数据
      oriData: null,
      form: {},
      tableData: [],
      refKeys: [] // 引用类型字段集合
    }
  },
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 500
  },
  methods: {
    addRow () {
      let entityData = {}
      this.tableColumns.forEach(item => {
        entityData[item.name] = ''
        this.refKeys.forEach(async rfk => {
          entityData[rfk + 'Options'] = []
        })
      })
      let data = {
        dataId: '',
        displayName: '',
        entityData: entityData,
        entityName: this.oriData.entity,
        entityDataOp: 'create',
        fullDataId: '',
        id: '',
        packageName: this.oriData.packageName,
        previousIds: [],
        succeedingIds: []
      }
      let find = this.dataArray.find(d => d.itemGroup === this.oriData.itemGroup)
      find.value.push(data)
      this.initData(this.rootEntityId, this.dataArray, find, this.requestId)
    },
    deleteRow (index) {
      let find = this.dataArray.find(d => d.itemGroup === this.oriData.itemGroup)
      find.value.splice(index, 1)
      this.initData(this.rootEntityId, this.dataArray, find, this.requestId)
    },
    async getRefOptions (formItem, formData, index) {
      if (formItem.refEntity === '') {
        formData[formItem.name + 'Options'] = formItem.selectList
        this.$set(this.tableData, index, formData)
        return
      }
      let cache = JSON.parse(JSON.stringify(formData))
      const keys = Object.keys(cache)
      keys.forEach(key => {
        if (Array.isArray(cache[key])) {
          cache[key] = cache[key].map(c => {
            return {
              guid: c
            }
          })
          cache[key] = JSON.stringify(cache[key])
        }
      })
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
    async initData (rootEntityId, dataArray, data, requestId, formDisable, jumpFrom) {
      this.formDisable = formDisable
      this.jumpFrom = jumpFrom
      this.rootEntityId = rootEntityId
      this.dataArray = dataArray
      this.requestId = requestId
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
      this.$emit('backData', this.dataArray)
    }
  },
  components: {}
}
</script>
<style>
.ivu-table-cell {
  padding-left: 8px;
  padding-right: 8px;
}
</style>
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
