<template>
  <div>
    <Row>
      <Col span="24">
        <Select
          ref="select"
          :value="value"
          @on-change="handleSelect"
          @on-query-change="handleQuery"
          @on-open-change="handleOpenChange"
          clearable
          filterable
          :disabled="disabled"
          :multiple="multiple"
        >
          <template v-for="i in optionsData">
            <Option :label="i[displayName]" :value="i[displayValue]" :key="i[displayValue]"> </Option>
          </template>
          <!--加载更多-->
          <Option v-if="optionsData.length >= 20" style="padding:0px;" label="" value="">
            <div style="width:100%;height:30px;" @click.stop>
              <Icon
                type="ios-more"
                color="#2d8cf0"
                size="24"
                style="margin-right:12px;float:right"
                @click="handleLoadMore()"
              />
            </div>
          </Option>
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
          }, 50)
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
      default: 220
    }
  },
  data () {
    return {
      optionsData: [], // 当前下拉数据
      sourceData: [], // 备份原始下拉数据
      sourceDataFilter: [], // 模糊搜索过滤数据
      currentPage: 1,
      pageSize: 1000,
      query: '', // 模糊搜索条件
      limitQueryFlag: true
    }
  },
  watch: {
    options: {
      handler (val) {
        if (val) {
          // 下拉默认选项值置顶
          this.sourceData = deepClone(val)
          // if (Array.isArray(this.value)) {
          //   if (this.objectOption) {
          //     this.value.forEach(i => {
          //       const index = this.sourceData.findIndex(j => j[this.displayValue] === i)
          //       if (index > 0) {
          //         const item = this.sourceData.splice(index, 1)
          //         this.sourceData.unshift(...item)
          //       }
          //     })
          //   } else {
          //     this.value.forEach(i => {
          //       const index = this.sourceData.findIndex(j => j === i)
          //       if (index > 0) {
          //         const item = this.sourceData.splice(index, 1)
          //         this.sourceData.unshift(...item)
          //       }
          //     })
          //   }
          // } else {
          //   if (this.objectOption) {
          //     const index = this.sourceData.findIndex(j => j[this.displayValue] === this.value)
          //     if (index > 0) {
          //       const item = this.sourceData.splice(index, 1)
          //       this.sourceData.unshift(...item)
          //     }
          //   } else {
          //     const index = this.sourceData.findIndex(j => j === this.value)
          //     if (index > 0) {
          //       const item = this.sourceData.splice(index, 1)
          //       this.sourceData.unshift(...item)
          //     }
          //   }
          // }
          this.optionsData = this.sourceData.slice(0, 1 * this.pageSize)
        }
      },
      immediate: true,
      deep: true
    }
  },
  methods: {
    handleSelect (val) {
      // this.limitQueryFlag = true
      this.$emit('input', val)
      this.$emit('on-change', val)
    },
    getList () {
      // if (this.query) {
      //   if (this.objectOption) {
      //     this.sourceDataFilter = this.sourceData.filter(item => item[this.displayName].includes(this.query))
      //   } else {
      //     this.sourceDataFilter = this.sourceData.filter(item => item.includes(this.query))
      //   }
      // } else {
      //   this.sourceDataFilter = this.sourceData
      // }
      this.sourceDataFilter = this.sourceData
      this.optionsData = this.sourceDataFilter.slice(0, this.currentPage * this.pageSize)
    },
    handleLoadMore () {
      if (this.sourceDataFilter.length === this.optionsData.length) return
      this.currentPage++
      this.getList()
    },
    handleQuery (val) {
      // if (this.limitQueryFlag) {
      //   this.limitQueryFlag = false
      // } else {
      //   this.query = val
      //   this.currentPage = 1
      //   this.getList()
      // }
    },
    handleOpenChange (flag) {
      if (flag) {
        // this.limitQueryFlag = true
        this.$emit('open-change')
      }
    }
  }
}
</script>
