<template>
  <div class="workbench-custom-form">
    <Form :model="value" ref="form" label-position="right" :label-width="120">
      <template v-for="(i, index) in options">
        <FormItem
          :label="i.title"
          :prop="i.name"
          :key="index"
          :required="i.required === 'yes'"
          :rules="i.required === 'yes' ? [{ required: true, message: `${i.title}为空`, trigger: 'change' }] : []"
          style="margin-bottom:20px;"
        >
          <!--输入框-->
          <Input
            v-if="i.elementType === 'input'"
            v-model="value[i.name]"
            :disabled="i.isEdit === 'no'"
            style="width:60%;"
          ></Input>
          <Input
            v-else-if="i.elementType === 'textarea'"
            v-model="value[i.name]"
            type="textarea"
            :disabled="i.isEdit === 'no'"
            style="width:60%;"
          ></Input>
          <LimitSelect
            v-else-if="i.elementType === 'select' || i.elementType === 'wecmdbEntity'"
            v-model="value[i.name]"
            :displayName="i.elementType === 'wecmdbEntity' ? 'displayName' : 'key_name'"
            :displayValue="i.elementType === 'wecmdbEntity' ? 'id' : 'guid'"
            :objectOption="!!i.entity || i.elementType === 'wecmdbEntity'"
            :options="entityData[i.name + 'Options']"
            :disabled="i.isEdit === 'no'"
            :multiple="i.multiple === 'Y'"
            style="width:60%;"
          >
          </LimitSelect>
        </FormItem>
      </template>
    </Form>
  </div>
</template>

<script>
import LimitSelect from '@/pages/components/limit-select.vue'
import { getRefOptions, getWeCmdbOptions } from '@/api/server'
import { deepClone } from '@/pages/util'
export default {
  components: {
    LimitSelect
  },
  props: {
    requestId: {
      type: String,
      default: ''
    },
    value: {
      type: Object,
      default: () => {}
    },
    options: {
      type: Array,
      default: () => []
    }
  },
  data () {
    return {
      refKeys: [],
      entityData: {}
    }
  },
  watch: {
    value: {
      handler (val) {
        if (val) {
          this.entityData = deepClone(val)
        }
      },
      deep: true
    },
    options: {
      handler (val) {
        if (val && val.length) {
          // select类型集合
          this.refKeys = []
          val.forEach(t => {
            if (t.elementType === 'select' || t.elementType === 'wecmdbEntity') {
              this.refKeys.push(t.name)
            }
          })
          // value数据初始化
          this.refKeys.forEach(rfk => {
            // 缓存RefOptions数据，不需要每次调用
            if (!(this.entityData[rfk + 'Options'] && this.entityData[rfk + 'Options'].length > 0)) {
              this.$set(this.entityData, rfk + 'Options', [])
            }
          })
          // 下拉类型数据初始化(待优化，调用接口太多)
          this.refKeys.forEach(rfk => {
            if (!(this.entityData[rfk + 'Options'] && this.entityData[rfk + 'Options'].length > 0)) {
              const titleObj = val.find(f => f.name === rfk)
              this.getRefOptions(titleObj)
            }
          })
        }
      },
      immediate: true,
      deep: true
    }
  },
  methods: {
    async getRefOptions (titleObj) {
      // taskman模板管理配置的普通下拉类型(值用逗号拼接)
      if (titleObj.elementType === 'select' && titleObj.entity === '') {
        // this.entityData[titleObj.name + 'Options'] = (titleObj.dataOptions && titleObj.dataOptions.split(',')) || []
        this.$set(
          this.entityData,
          titleObj.name + 'Options',
          (titleObj.dataOptions && titleObj.dataOptions.split(',')) || []
        )
        return
      }
      // taskman模板管理配置的引用下拉类型
      if (titleObj.elementType === 'wecmdbEntity') {
        const [packageName, ciType] = (titleObj.dataOptions && titleObj.dataOptions.split(':')) || []
        const { status, data } = await getWeCmdbOptions(packageName, ciType, {})
        if (status === 'OK') {
          this.$set(this.entityData, titleObj.name + 'Options', data)
        }
        return
      }
      // if (titleObj.refEntity === '') {
      //   row[titleObj.name + 'Options'] = titleObj.selectList
      //   this.$set(this.tableData, index, row)
      //   return
      // }
      let cache = JSON.parse(JSON.stringify(this.entityData))
      cache[titleObj.name] = ''
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
        // 删除掉值为空的数据
        if (!cache[key] || (Array.isArray(cache[key]) && cache[key].length === 0)) {
          delete cache[key]
        }
      })
      this.refKeys.forEach(k => {
        delete cache[k + 'Options']
      })
      delete cache._checked
      delete cache._disabled
      const filterValue = this.entityData[titleObj.name]
      const attr = titleObj.entity + '__' + titleObj.name
      const params = {
        filters: [
          {
            name: 'guid',
            operator: 'in',
            value: Array.isArray(filterValue) ? filterValue : [filterValue]
          }
        ],
        paging: false,
        dialect: {
          associatedData: {
            ...cache
          }
        }
      }
      const { statusCode, data } = await getRefOptions(this.requestId, attr, params)
      if (statusCode === 'OK') {
        this.$set(this.entityData, titleObj.name + 'Options', data)
      }
    }
  }
}
</script>
