<template>
  <div class="header-title">
    <div class="w-header">
      <Icon v-if="expand" size="26" type="md-arrow-dropdown" @click="handleExpand" style="cursor:pointer" />
      <Icon v-else size="26" type="md-arrow-dropright" @click="handleExpand" style="cursor:pointer" />
      <div class="title">{{ title }}<span class="underline"></span></div>
      <span class="sub-title">
        <slot name="sub-title"></slot>
      </span>
    </div>
    <div v-show="expand" class="content">
      <slot></slot>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      default: ''
    },
    subTitle: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      expand: true
    }
  },
  methods: {
    handleExpand () {
      this.expand = !this.expand
    }
  }
}
</script>

<style lang="scss" scoped>
.header-title {
  .w-header {
    display: flex;
    align-items: center;
    .title {
      font-size: 16px;
      font-weight: bold;
      color: #515a6e;
      margin: 0 6px;
      .underline {
        display: block;
        margin-top: -12px;
        margin-left: -8px;
        width: 100%;
        padding: 0 8px;
        height: 12px;
        border-radius: 12px;
        background-color: #c6eafe;
        box-sizing: content-box;
      }
    }
    .sub-title {
      font-size: 15px;
    }
  }
  .content {
    padding: 20px 10px;
  }
}
</style>
