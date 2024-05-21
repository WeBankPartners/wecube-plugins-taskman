<template>
  <div class="temp-hidden-condition">
    <div class="title">
      <span>隐藏条件</span>
      <Button @click.stop="handleOpenModal" type="primary" size="small" icon="md-add" class="btn"></Button>
    </div>
    <div class="box">
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
      :name="name"
      @updateData="handleUpdate"
    />
  </div>
</template>

<script>
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
    name: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      hiddenCondition: [],
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
    }
  },
  methods: {
    handleOpenModal () {
      this.$refs.modal.initData(this.value || [])
    },
    handleUpdate (val) {
      this.hiddenCondition = val
      this.$emit('input', val)
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
