<template>
  <div>
    <Row>
      <Col span="4">
        <Select
          ref="select"
          :value="data"
          @on-change="handleSelect"
          v-loadmore="hanldeLoadMore"
          @on-query-change="handleQuery"
          clearable
          filterable
          :disabled="disabled"
          :multiple="multiple"
          :style="{ width: width + 'px' }"
        >
          <template v-for="i in optionsData">
            <Option
              :label="objectOption ? i[displayName] : i"
              :value="objectOption ? i[displayValue] : i"
              :key="objectOption ? i[displayValue] : i"
            >
            </Option>
          </template>
        </Select>
      </Col>
    </Row>
  </div>
</template>

<script>
import { debounce, deepClone } from '@/pages/util'
export default {
  name: '',
  directives: {
    loadmore: {
      bind (el, binding) {
        const wrap = el.querySelector('.ivu-select-dropdown')
        wrap.addEventListener(
          'scroll',
          debounce(() => {
            const flag = wrap.scrollHeight - wrap.scrollTop <= wrap.clientHeight
            if (flag) {
              binding.value()
            }
          }, 300)
        )
      }
    }
  },
  props: {
    // 下拉显示名
    displayName: {
      type: String,
      default: 'label'
    },
    // 下拉绑定值
    displayValue: {
      type: String,
      default: 'value'
    },
    // 下拉选项是否是对象类型
    objectOption: {
      type: Boolean,
      default: false
    },
    value: {
      type: String | Array
    },
    // 下拉选项
    options: {
      type: Array,
      default: () => []
    },
    disabled: {
      type: Boolean,
      default: false
    },
    multiple: {
      type: Boolean,
      default: false
    },
    width: {
      type: Number,
      default: 250
    }
  },
  data () {
    return {
      data: '',
      optionsData: [], // 下拉数据
      sourceDataFilter: [], // 过滤后数据
      currentPage: 1,
      pageSize: 20,
      query: ''
    }
  },
  watch: {
    value: {
      handler (val) {
        this.data = val
      },
      immediate: true,
      deep: true
    },
    options: {
      handler (val) {
        if (val) {
          let initOptions = deepClone(val)
          if (Array.isArray(this.value)) {
            if (this.objectOption) {
              this.value.forEach(i => {
                const index = initOptions.findIndex(j => j[this.displayValue] === i)
                const item = initOptions.splice(index, 1)
                initOptions.unshift(item)
              })
            } else {
              this.value.forEach(i => {
                const index = initOptions.findIndex(j => j === i)
                const item = initOptions.splice(index, 1)
                initOptions.unshift(item)
              })
            }
          } else {
            if (this.objectOption) {
              const index = initOptions.findIndex(j => j[this.displayValue] === this.value)
              const item = initOptions.splice(index, 1)
              initOptions.unshift(item)
            } else {
              const index = initOptions.findIndex(j => j === this.value)
              const item = initOptions.splice(index, 1)
              initOptions.unshift(item)
            }
          }
          this.optionsData = initOptions.slice(0, 1 * this.pageSize)
        }
      },
      immediate: true,
      deep: true
    }
  },
  methods: {
    handleSelect (val) {
      this.data = val
      this.$emit('on-change', val)
    },
    getList () {
      if (this.query) {
        this.sourceDataFilter = this.options.filter(item => item[this.displayName].includes(this.query))
      } else {
        this.sourceDataFilter = this.options
      }
      this.optionsData = this.sourceDataFilter.slice(0, this.currentPage * this.pageSize)
    },
    hanldeLoadMore () {
      if (this.sourceDataFilter.length === this.optionsData.length) return
      this.currentPage++
      this.getList()
    },
    handleQuery (val) {
      this.query = val
      this.currentPage = 1
      this.getList()
    }
  }
}
</script>
