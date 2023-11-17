<template>
  <div class="header-title">
    <div class="header">
      <Icon
        v-if="expand"
        size="26"
        type="md-arrow-dropdown"
        @click="handleExpand"
        style="cursor:pointer"
      />
      <Icon
        v-else
        size="26"
        type="md-arrow-dropright"
        @click="handleExpand"
        style="cursor:pointer"
      />
      <div class="title">{{ title }}<span class="underline"></span></div>
    </div>
    <div class="content">
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
    }
  },
  data() {
    return {
      expand: true
    }
  },
  methods: {
    handleExpand() {
      this.expand = !this.expand
    }
  }
}
</script>

<style lang="scss" scoped>
.header-title {
  .header {
    display: flex;
    align-items: center;
    .title {
      font-size: 15px;
      font-weight: 500;
      color: #282e38;
      margin: 0 10px;
      .underline {
        display: block;
        margin-top: -10px;
        margin-left: -6px;
        width: 100%;
        padding: 0 6px;
        height: 10px;
        border-radius: 12px;
        background-color: #c6eafe;
        box-sizing: content-box;
      }
    }
  }
  .content {
    padding: 20px 10px;
  }
}
</style>