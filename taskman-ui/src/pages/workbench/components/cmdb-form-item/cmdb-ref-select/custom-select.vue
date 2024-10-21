<!--
 * @Author: wanghao7717 792974788@qq.com
 * @Date: 2024-10-14 15:05:46
 * @LastEditors: wanghao7717 792974788@qq.com
 * @LastEditTime: 2024-10-16 19:50:41
-->
<template>
  <div>
    <Poptip v-model="visible" placement="bottom" class="custom-select" @on-popper-show="handleOpenChange">
      <div ref="input" :class="{ 'custom-select-input': true, 'custom-select-disabled': disabled }">
        <div class="ref" @click.stop="handleOpenRefModal($event)">@</div>
        <div class="tags">
          <Tag
            v-for="i in optionsData"
            v-show="selected.includes(i.guid)"
            closable
            @on-close="handleRemoveItem(i)"
            @click.native="visible = !visible"
            :key="i.guid"
            >{{ i.key_name }}</Tag
          >
          <!-- <Tag v-if="selected && selected.length > 1">+{{ selected.length - 1 }}</Tag> -->
        </div>
        <div class="icon">
          <Icon v-if="visible" type="ios-arrow-up" />
          <Icon v-else type="ios-arrow-down" />
        </div>
      </div>
      <div slot="content" class="custom-select-content" :style="{ width: width + 'px' }">
        <!-- <div class="dropdown-select">
          <span>选中结果:</span>
          <Tag
            v-for="i in optionsData"
            v-show="selected.includes(i.guid)"
            :key="i.guid"
          >{{i.key_name}}</Tag>
        </div> -->
        <Input v-model="keyword" @input="handleFilterOptions" placeholder="关键字搜索" style="width:100%;"></Input>
        <div class="dropdown-wrap">
          <div
            v-for="i in optionsData.filter(i => i.isShow)"
            :key="i.guid"
            :class="{ 'dropdown-wrap-item': true, 'dropdown-wrap-item-active': i.checked }"
            @click="handleSelectItem(i)"
          >
            {{ i.key_name }}<Icon v-if="i.checked" type="ios-checkmark" size="24" color="#2d8cf0e6" />
          </div>
        </div>
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
  watch: {
    options: {
      handler (val) {
        if (val && val.length > 0) {
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
        this.selected = JSON.parse(JSON.stringify(val))
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
    handleRemoveItem (item) {
      if (this.disabled) return
      const index = this.selected.findIndex(i => i === item.guid)
      this.selected.splice(index, 1)
      this.$emit('input', this.selected)
      this.$emit('change', this.selected)
    },
    // 打开引用数据弹框回调
    handleOpenRefModal (e) {
      if (this.disabled) return
      this.$emit('showRefModal', e)
    },
    // 下拉展开回调
    handleOpenChange () {
      this.$emit('openChange', true)
    }
  }
}
</script>
<style lang="scss">
.custom-select {
  width: 100%;
  &-input {
    width: 100%;
    min-height: 32px;
    max-height: 60px;
    padding: 0 3px;
    overflow: hidden;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;
    .ref {
      width: 24px;
    }
    .tags {
      width: 100%;
      display: flex;
      flex-wrap: wrap;
      align-items: center;
      overflow-y: auto;
      min-height: 32px;
      max-height: 60px;
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
      &-select {
        margin-bottom: 5px;
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
        }
        &-item-active {
          color: #2d8cf0e6;
        }
      }
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
    max-width: 120px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    color: #515a6e;
  }
  .ivu-tag {
    width: fit-content;
    display: flex;
    align-items: center;
  }
}
</style>
