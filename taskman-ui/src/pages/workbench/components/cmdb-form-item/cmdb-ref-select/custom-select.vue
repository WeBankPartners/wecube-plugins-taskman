<!--
 * @Author: wanghao7717 792974788@qq.com
 * @Date: 2024-10-14 15:05:46
 * @LastEditors: wanghao7717 792974788@qq.com
 * @LastEditTime: 2025-01-06 19:21:23
-->
<template>
  <div>
    <Poptip
      v-model="visible"
      transfer
      popper-class="taskman-cmdb-custom-select-popper"
      placement="bottom-start"
      class="taskman-cmdb-custom-select"
      @on-popper-show="handleOpenChange"
    >
      <!--模拟输入框显示效果-->
      <div ref="input" :class="{ 'taskman-cmdb-custom-select-input': true, 'taskman-cmdb-custom-select-disabled': disabled }">
        <div class="ref" @click.stop="handleOpenRefModal($event)">@</div>
        <div class="tags">
          <template v-if="getSelectedOptions && getSelectedOptions.length > 0">
            <Tag
              v-for="(i, index) in getSelectedOptions"
              v-show="index === 0"
              @click.native="visible = !visible"
              closable
              @on-close="handleRemoveItem($event, i)"
              :key="i.guid"
              style="max-width:calc(100% - 50px)"
              >{{ i.key_name }}</Tag
            >
            <Tag v-if="selected && selected.length > 1" @click.native="visible = !visible">
              + {{ selected.length - 1 }}</Tag
            >
          </template>
          <template v-else>
            <span style="color:#dcdee2;">{{ $t('tw_please_select') }}</span>
          </template>
        </div>
        <div class="icon">
          <Icon v-if="visible" type="ios-arrow-up" />
          <Icon v-else type="ios-arrow-down" />
        </div>
      </div>
      <!--模拟下拉框显示效果-->
      <div
        slot="content"
        class="taskman-cmdb-custom-select-content"
        :style="{ minWidth: width + 'px', width: 'fit-content', maxWidth: '500px' }"
      >
        <!-- 选中结果 -->
        <div class="dropdown-selected">
          <span class="dropdown-selected-title">{{ $t('tw_select_result') }}:</span>
          <Tag v-for="i in getSelectedOptions" :key="i.guid" closable @on-close="handleRemoveItem($event, i)">{{
            i.key_name
          }}</Tag>
        </div>
        <!--关键字搜索-->
        <Input v-model="keyword" @input="handleFilterOptions" :placeholder="$t('tw_keyword_search')" style="width:100%;"></Input>
        <div v-if="getFilterOptions.length > 0" class="dropdown-wrap">
          <div
            v-for="i in getFilterOptions"
            :key="i.guid"
            :class="{ 'dropdown-wrap-item': true, 'dropdown-wrap-item-active': i.checked }"
            @click="handleSelectItem(i)"
          >
            <span>{{ i.key_name }}</span>
            <Icon v-if="i.checked" type="ios-checkmark" size="24" color="#5384FF" />
          </div>
        </div>
        <div v-else class="no-data">{{ $t('tw_no_data') }}</div>
      </div>
    </Poptip>
  </div>
</template>
<script>
export default {
  name: 'CustomSelect',
  props: {
    options: {
      type: Array,
      default: () => []
    },
    value: {
      type: Array | String,
      default: () => []
    },
    disabled: {
      type: Boolean,
      default: false
    },
    title: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      keyword: '',
      optionsData: [],
      selected: [],
      width: 300,
      isMultiple: true,
      visible: false
    }
  },
  computed: {
    getFilterOptions () {
      return this.optionsData.filter(i => i.isShow && !this.selected.includes(i.guid))
    },
    getSelectedOptions () {
      return this.optionsData.filter(i => this.selected.includes(i.guid))
    }
  },
  watch: {
    options: {
      handler (val) {
        if (val && Array.isArray(val)) {
          this.optionsData = JSON.parse(JSON.stringify(val))
          this.optionsData.forEach(i => {
            this.$set(i, 'checked', false)
            this.$set(i, 'isShow', true)
          })
          this.initData()
        }
      },
      immediate: true,
      deep: true
    },
    value: {
      handler (val) {
        if (val && Array.isArray(val)) {
          this.selected = JSON.parse(JSON.stringify(val))
        }
        this.initData()
      },
      immediate: true,
      deep: true
    }
  },
  mounted () {
    this.width = this.$refs.input.clientWidth - 32
  },
  methods: {
    initData () {
      if (this.optionsData.length > 0) {
        this.optionsData.forEach(i => {
          // 多选
          if (this.isMultiple) {
            if (this.selected.includes(i.guid)) {
              i.checked = true
            } else {
              i.checked = false
            }
          } else {
            if (this.selected === i.guid) {
              i.checked = true
            }
          }
        })
      }
    },
    handleFilterOptions () {
      this.optionsData.forEach(i => {
        if (i.key_name.toLowerCase().indexOf(this.keyword.toLowerCase()) > -1) {
          i.isShow = true
        } else {
          i.isShow = false
        }
      })
    },
    handleSelectItem (item) {
      if (this.disabled) return
      if (!this.isMultiple) {
        this.optionsData.forEach(i => {
          if (i.guid !== item.guid) {
            i.checked = false
          }
        })
      }
      this.optionsData.forEach(i => {
        if (i.guid === item.guid) {
          i.checked = !i.checked
          // 多选
          if (this.isMultiple) {
            if (i.checked) {
              this.selected.push(item.guid)
            } else {
              const index = this.selected.findIndex(i => i === item.guid)
              this.selected.splice(index, 1)
            }
          } else {
            if (i.checked) {
              this.selected = item.guid
            } else {
              this.selected = ''
            }
          }
          this.$emit('input', this.selected)
          this.$emit('change', this.selected)
        }
      })
    },
    handleRemoveItem (e, item) {
      e.stopPropagation()
      if (this.disabled) return
      const index = this.selected.findIndex(i => i === item.guid)
      this.selected.splice(index, 1)
      this.$emit('input', this.selected)
      this.$emit('change', this.selected)
    },
    // 点击@符号打开弹框
    handleOpenRefModal (e) {
      this.visible = false
      this.$emit('showRefModal', e)
    },
    // 下拉展开回调
    handleOpenChange () {
      this.keyword = ''
      this.$emit('openChange', true)
    }
  }
}
</script>
<style lang="scss">
.taskman-cmdb-custom-select {
  width: 100%;
  &-input {
    width: 100%;
    height: 32px;
    padding: 0 3px;
    overflow: hidden;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border: 1px solid #dcdee2;
    border-radius: 4px;
    .ref {
      width: 24px;
      cursor: pointer !important;
      color: #000;
    }
    .tags {
      width: 100%;
      display: flex;
      flex-wrap: wrap;
      align-items: center;
      height: 32px;
    }
    .icon {
      width: 20px;
    }
  }
  &-disabled {
    color: #ccc;
    background: #f3f3f3;
    cursor: not-allowed;
  }
  &-content {
    .dropdown {
      &-selected {
        &-title {
          font-weight: bold;
          margin-right: 5px;
        }
        margin-bottom: 5px;
        width: 100;
        display: flex;
        flex-wrap: wrap;
        align-items: center;
      }
      &-wrap {
        max-height: 400px;
        overflow-y: auto;
        display: flex;
        flex-direction: column;
        padding-top: 10px;
        &-item {
          padding: 3px 0px;
          cursor: pointer;
          width: 100%;
          display: flex;
          align-items: center;
          justify-content: space-between;
          &:hover {
            background: #f3f3f3;
          }
          span {
            display: block;
            max-width: 460px;
            overflow: hidden;
            text-overflow: ellipsis;
          }
        }
        &-item-active {
          color: #5384FF;
        }
      }
    }
    .no-data {
      font-size: 12px;
      text-align: center;
      margin-top: 10px;
    }
  }
  .ivu-poptip-rel {
    width: 100%;
  }
  .ivu-select-dropdown {
    display: none;
  }
  .ivu-tag-text {
    display: inline-block;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    color: #515a6e;
  }
  .ivu-tag {
    width: fit-content;
    display: flex;
    align-items: center;
    line-height: 10px !important;
    padding: 4px 6px !important;
  }
}
.taskman-cmdb-custom-select-popper .ivu-poptip-body {
  padding: 8px 16px !important;
}
</style>
