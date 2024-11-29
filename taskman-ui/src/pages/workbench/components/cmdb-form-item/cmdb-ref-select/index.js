/*
 * @Author: wanghao7717 792974788@qq.com
 * @Date: 2024-10-21 19:33:55
 * @LastEditors: wanghao7717 792974788@qq.com
 * @LastEditTime: 2024-11-28 18:27:07
 */
import { getRefOptions, getWeCmdbOptions } from '@/api/server'
import CustomMultipleSelect from './custom-select.vue'
import JsonViewer from 'vue-json-viewer'
import './index.scss'
export default {
  name: 'WeCMDBRefSelect',
  components: { CustomMultipleSelect, JsonViewer },
  props: {
    value: {},
    ciTypeAttrId: '',
    disabled: { default: () => false },
    filterParams: {},
    guidFilters: { default: () => null },
    guidFilterEnabled: { default: () => false },
    column: { type: Object, default: () => ({}) }
  },
  watch: {
    value: {
      handler (val) {
        this.selected = val
      },
      immediate: true
    },
    column: {
      handler (val) {
        if (val && Object.keys(val).length > 0) {
          this.getAllDataWithoutPaging()
        }
      },
      immediate: true,
      deep: true
    }
  },
  data () {
    return {
      allTableDataWithoutPaging: [],
      selected: [],
      selectDisabled: true,
      firstInput: true,
      firstChange: true,
      showDetail: false,
      detailData: []
    }
  },
  methods: {
    async getAllDataWithoutPaging () {
      const cache = this.filterParams ? JSON.parse(JSON.stringify(this.filterParams.params)) : {}
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
        // 数据表单【表单隐藏标识】放到了row里面，需要删除
        if (key.indexOf('Hidden') > -1) {
          delete cache[key]
        }
        // 将对象类型转为字符串
        if (typeof cache[key] === 'object') {
          cache[key] = JSON.stringify(cache[key])
        }
      })
      this.column.refKeys.forEach(k => {
        delete cache[k + 'Options']
      })
      delete cache._checked
      delete cache._disabled
      const { titleObj } = this.column || { titleObj: {} }
      const attrName = titleObj.entity + '__' + titleObj.name
      const attr = titleObj.id
      const params = {
        filters: [],
        paging: false,
        dialect: {
          associatedData: {
            ...cache
          }
        }
      }
      const res = await getRefOptions(this.column.requestId, attr, params, attrName)
      if (res.statusCode === 'OK') {
        let data = res.data
        if (this.guidFilterEnabled && this.guidFilters && Array.isArray(this.guidFilters)) {
          data = data.filter(el => {
            if (this.guidFilters.indexOf(el.guid) >= 0) {
              return true
            }
            return false
          })
        }
        this.selectDisabled = false
        this.allTableDataWithoutPaging = data
      }
    },
    handleInput (v) {
      if (this.column.isMultiple && this.firstInput) {
        this.firstInput = false
        return
      }
      this.selected = v
    },
    selectChangeHandler (val) {
      if (this.column.isMultiple && this.firstChange) {
        this.firstChange = false
        return
      }
      if (Array.isArray(val)) {
        this.$emit('input', val)
      } else {
        this.$emit('input', val || '')
      }
    },
    getFilterRulesOptions (val) {
      this.firstInput = false
      this.firstChange = false
      if (val) {
        this.getAllDataWithoutPaging()
      }
    },
    // ref勾选数据详情查看
    async showRefModal (e) {
      e.preventDefault()
      e.stopPropagation()
      this.detailData = []
      this.showDetail = true
      const { titleObj } = this.column || { titleObj: {} }
      const { refEntity, refPackageName } = titleObj || { refEntity: '', refPackageName: '' }
      if (!refPackageName || !refEntity) return
      const { status, data } = await getWeCmdbOptions(refPackageName, refEntity, {
        filters: [
          {
            name: 'guid',
            operator: 'in',
            value: Array.isArray(this.selected) ? this.selected : [this.selected]
          }
        ]
      })
      if (status === 'OK') {
        const contents = data || []
        if (this.column.isMultiple) {
          contents.forEach(item => {
            if (this.selected.includes(item.guid)) {
              this.detailData.push({
                title: item.key_name,
                value: item
              })
            }
          })
        } else {
          contents.forEach(item => {
            if (item.guid === this.selected) {
              this.detailData.push({
                title: item.key_name,
                value: item
              })
            }
          })
        }
      }
    },
    closeModal () {
      this.showDetail = false
    },
    refreshDiffVariable () {}
  },
  render (h) {
    let renderOptions = this.allTableDataWithoutPaging.map(_ => {
      return (
        <Option value={_.guid || ''} key={_.guid}>
          {_.key_name}
        </Option>
      )
    })
    return (
      <div class="cmdb-ref-select">
        {!this.column.isMultiple && (
          <Select
            onInput={this.handleInput}
            value={this.selected}
            disabled={this.selectDisabled || this.disabled}
            style="width:100%"
            filterable
            clearable
            on-on-change={this.selectChangeHandler}
            on-on-open-change={this.getFilterRulesOptions}
            max-tag-count={2}
          >
            <span slot="prefix" onClick={e => this.showRefModal(e)}>
              @
            </span>
            {renderOptions}
          </Select>
        )}
        {// 引用多选下拉框组件封装
          this.column.isMultiple && (
            <CustomMultipleSelect
              options={this.allTableDataWithoutPaging}
              onShowRefModal={e => this.showRefModal(e)}
              onChange={val => this.$emit('input', val)}
              onOpenChange={this.getFilterRulesOptions}
              disabled={this.selectDisabled || this.disabled}
              v-model={this.selected}
            ></CustomMultipleSelect>
          )}
        <Modal value={this.showDetail} footer-hide={true} title={this.column.title} width={1100}>
          <div style="overflow: auto;max-height:500px;overflow:auto">
            <Collapse>
              {this.detailData.map(column => {
                return (
                  <Panel name={column.title}>
                    {column.title}
                    <p slot="content">
                      <JsonViewer value={column.value} expand-depth={1} class="ref-select-jsonviewer"></JsonViewer>
                    </p>
                  </Panel>
                )
              })}
            </Collapse>
          </div>
          <div style="margin-top:20px;height: 30px">
            <Button style="float: right;margin-right: 20px" onClick={() => this.closeModal()}>
              {this.$t('tw_close')}
            </Button>
            {/* <Button style="float: right;margin-right: 20px" type="primary" onClick={() => this.refreshDiffVariable()}>
              {this.$t('tw_refresh')}
            </Button> */}
          </div>
        </Modal>
      </div>
    )
  }
}
