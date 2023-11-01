<template>
  <div>
    <Row>
      <Col span="4">
        <Select
          v-model="value"
          v-loadmore="hanldeLoadMore"
          @on-query-change="handleQuery"
          clearable
          filterable
          style="width: 90%"
        >
          <template v-for="i in optionsData">
            <Option :label="i[sourceLabel]" :value="i[sourceValue]" :key="i[sourceValue]"> </Option>
          </template>
        </Select>
      </Col>
    </Row>
  </div>
</template>

<script>
export default {
  name: '',
  directives: {
    loadmore: {
      bind (el, binding) {
        const wrap = el.querySelector('.ivu-select-dropdown')
        wrap.addEventListener('scroll', () => {
          const flag = wrap.scrollHeight - wrap.scrollTop <= wrap.clientHeight
          if (flag) {
            binding.value()
          }
        })
      }
    }
  },
  props: {
    sourceData: {
      type: Array,
      default: () => []
    },
    sourceLabel: {
      type: String,
      default: 'label'
    },
    sourceValue: {
      type: String,
      default: 'value'
    }
  },
  data () {
    return {
      value: '',
      optionsData: [], // 下拉数据
      sourceDataFilter: [], // 过滤后数据
      currentPage: 1,
      pageSize: 20,
      query: ''
    }
  },
  mounted () {
    this.getList()
  },
  methods: {
    getList () {
      if (this.query) {
        this.sourceDataFilter = this.sourceData.filter(item => item[this.sourceLabel].includes(this.query))
      } else {
        this.sourceDataFilter = this.sourceData
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
