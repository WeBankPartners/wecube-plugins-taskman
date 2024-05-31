<template>
  <div class="temp-hidden-condition">
    <div class="title">
      <span>
        隐藏条件
        <span v-if="editElement.required === 'yes'" class="tips">属性必填，无法设置过滤条件</span>
      </span>
      <Button
        v-if="editElement.required === 'no'"
        @click="handleOpenModal"
        type="success"
        size="small"
        icon="md-add"
        class="btn"
      ></Button>
    </div>
    <div v-if="editElement.required === 'no'" class="box">
      <Row v-for="(i, index) in value" :key="index" class="box-item">
        <Col :span="6" class="name">{{ getFormItemTitle(i.name) }}</Col>
        <Col :span="6" class="operator">{{ operatorMap[i.operator] || '-' }}</Col>
        <Col :span="12" class="value">
          {{ getValueDisplay(i.value) || '-' }}
        </Col>
      </Row>
      <div v-if="!value || value.length === 0" class="no-data">
        暂无隐藏条件
      </div>
    </div>
    <HiddenConditionModal
      ref="modal"
      :finalElement="finalElement"
      :data="value"
      :disabled="disabled"
      :editElement="editElement"
      @updateData="handleUpdate"
    />
  </div>
</template>

<script>
import { deepClone } from '../../util'
import HiddenConditionModal from './hidden-condition-modal.vue'
export default {
  components: {
    HiddenConditionModal
  },
  props: {
    value: {
      type: Array,
      default: () => []
    },
    finalElement: {
      type: Array,
      default: () => []
    },
    disabled: {
      type: Boolean,
      default: false
    },
    editElement: {
      type: Object,
      default: () => {}
    }
  },
  data () {
    return {
      operatorMap: {
        eq: '等于',
        neq: '不等于',
        lt: '小于',
        gt: '大于',
        contains: '包含',
        startsWith: '匹配开始',
        endsWith: '匹配结束',
        containsAll: '包含全部',
        containsAny: '包含任意',
        notContains: '不包含',
        range: '在范围内',
        empty: '为空',
        notEmpty: '不为空'
      }
    }
  },
  computed: {
    getFormItemTitle () {
      return function (val) {
        let title = ''
        this.finalElement[0].attrs.forEach(i => {
          if (i.name === val) {
            title = i.title
          }
        })
        return title
      }
    },
    getValueDisplay () {
      return function (val) {
        if (Array.isArray(val)) {
          return val.join('-')
        } else {
          return val
        }
      }
    },
    attrs () {
      return this.finalElement[0].attrs
    }
  },
  watch: {
    attrs: {
      handler (val) {
        if (val) {
          this.removeConditionsByAttrs(this.value)
        }
      }
    }
  },
  methods: {
    handleOpenModal () {
      this.$refs.modal.initData(this.value || [])
    },
    handleUpdate (val) {
      this.$emit('input', val)
    },
    removeConditionsByAttrs (arr) {
      // 如果预览区内对应表单项删除，则清空该条过滤条件
      const deleteNameArr = []
      let hiddenCondition = deepClone(arr || [])
      hiddenCondition.forEach(i => {
        const exist = this.finalElement[0].attrs.some(j => j.name === i.name)
        if (!exist) {
          deleteNameArr.push(i.name)
        }
      })
      hiddenCondition = hiddenCondition.filter(item => !deleteNameArr.includes(item.name))
      this.$emit('input', hiddenCondition)
    }
  }
}
</script>

<style lang="scss">
.temp-hidden-condition {
  .title {
    display: flex;
    align-items: center;
    margin-bottom: 5px;
    .btn {
      margin-left: 25px;
    }
    .tips {
      font-size: 12px;
      margin-left: 20px;
    }
  }
  .box {
    display: flex;
    align-items: center;
    flex-direction: column;
    min-height: 50px;
    border: 1px dashed #d7dadc;
    padding: 0 10px 10px 10px;
    &-item {
      display: flex;
      font-size: 12px;
      line-height: 20px;
      margin-top: 10px;
      width: 100%;
      .name {
        text-align: center;
      }
      .operator {
        text-align: center;
        color: #2db7f5;
      }
      .value {
        text-align: center;
        overflow: hidden;
        max-width: 100%;
        text-overflow: ellipsis;
        overflow: hidden;
        word-break: break-word;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
      }
    }
    .no-data {
      display: flex;
      justify-content: center;
      align-items: center;
      height: 50px;
    }
  }
  .ivu-form-item {
    margin-bottom: 5px;
  }
}
</style>
