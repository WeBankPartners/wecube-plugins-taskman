<template>
  <div class="taskman-cmdb-diff-variable">
    <div class="inline">
      <div class="text">{{ data || $t('tw_no_data') }}</div>
      <Button
        v-if="data"
        @click="showDetail = true"
        type="primary"
        ghost
        icon="md-eye"
      ></Button>
    </div>
    <!--详情弹框-->
    <Modal :z-index="2000" v-model="showDetail" :title="$t('tw_diff_variable')" @on-ok="showDetail = false" width="1100">
      <div
        v-for="(val, index) in detailInfo"
        :key="index"
        @click="choiceKey(val)"
        :style="remarkedKeys.includes(val.key) ? 'background:#d9d9d9' : ''"
      >
        <div style="width: 300px;display:inline-block;word-break: break-all;margin:4px 0;vertical-align: top;text-align:right;cursor:pointer">
          <span :style="!['', 'NULL'].includes(val.value) ? '' : 'color:red'">
            {{val.key}}
          </span>
        </div>
        <div style="width: 740px;display:inline-block;word-break: break-all;margin:4px 0;">
          ：{{val.value}}
        </div>
      </div>
    </Modal>
  </div>
</template>

<script>
export default {
  props: ['data'],
  data () {
    return {
      remarkedKeys: [],
      detailInfo: null,
      showDetail: false
    }
  },
  mounted () {
    this.detailInfo = this.formatData(this.data)
  },
  methods: {
    choiceKey (chioceObj) {
      const key = chioceObj.key
      if (this.remarkedKeys.includes(key)) {
        // 元素存在于数组中，移除它
        const index = this.remarkedKeys.indexOf(key)
        this.remarkedKeys.splice(index, 1)
      } else {
        // 元素不存在于数组中，添加它
        this.remarkedKeys.push(key)
      }
    },
    formatData (val) {
      if (!val) return []
      const vari = val.split('\u0001=\u0001')
      const keys = vari[0].split(',\u0001')
      const values = vari[1].split(',\u0001')
      let res = []
      for (let i = 0; i < keys.length; i++) {
        res.push({
          key: (keys[i] || '').replace('\u0001', ''),
          value: (values[i] || '').replace('\u0001', '')
        })
      }
      res = res.sort((first, second) => {
        const firstKey = first.key.toLocaleUpperCase()
        const secondKey = second.key.toLocaleUpperCase()
        if (firstKey < secondKey) {
          return -1
        } else if (firstKey > secondKey) {
          return 1
        } else {
          return 0
        }
      })
      return res
    }
  }
}
</script>

<style lang="scss">
.taskman-cmdb-diff-variable {
  width: 100%;
  .inline {
    display: flex;
    align-items: center;
    height: 34px;
    .text {
      font-size: 13px;
      color:#515a6e;
      display: block;
      max-width: 450px;
      text-overflow: ellipsis;
      overflow: hidden;
      white-space: nowrap;
      margin-right: 6px;
    }
  }
}
</style>
