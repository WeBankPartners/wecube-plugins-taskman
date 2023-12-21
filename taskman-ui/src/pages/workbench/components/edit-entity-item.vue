<template>
  <div>
    <Drawer
      :title="$t('edit')"
      v-model="drawerVisible"
      width="720"
      :mask-closable="true"
      @on-close="handleCancel"
      class="workbench-edit-drawer"
    >
      <div class="content" :style="{ maxHeight: maxHeight + 'px' }">
        <Form :model="value" ref="form" label-position="left" :label-width="120">
          <template v-for="(i, index) in options">
            <FormItem
              :label="i.title"
              :prop="i.name"
              :key="index"
              :required="i.required === 'yes'"
              :rules="i.required === 'yes' ? [{ required: true, message: `${i.title}为空`, trigger: 'change' }] : []"
              style="margin-bottom:20px;"
            >
              <!--输入框-->
              <Input
                v-if="i.elementType === 'input'"
                v-model="value[i.name]"
                :disabled="i.isEdit === 'no' || disabled"
                style="width: calc(100% - 30px)"
              ></Input>
              <Input
                v-else-if="i.elementType === 'textarea'"
                v-model="value[i.name]"
                type="textarea"
                :disabled="i.isEdit === 'no' || disabled"
                style="width: calc(100% - 30px)"
              ></Input>
              <!--下拉选择-->
              <Select
                v-else-if="i.elementType === 'select' && i.entity"
                v-model="value[i.name]"
                :multiple="i.multiple === 'Y'"
                filterable
                clearable
                :disabled="i.isEdit === 'no' || disabled"
                style="width: calc(100% - 30px)"
              >
                <template v-for="item in value[i.name + 'Options']">
                  <Option :key="item.guid" :value="item.guid">{{ item.key_name }}</Option>
                </template>
              </Select>
              <Select
                v-else-if="i.elementType === 'select' && !i.entity"
                v-model="value[i.name]"
                :multiple="i.multiple === 'Y'"
                filterable
                clearable
                :disabled="i.isEdit === 'no' || disabled"
                style="width: calc(100% - 30px)"
              >
                <template v-for="item in value[i.name + 'Options']">
                  <Option :key="item" :value="item">{{ item }}</Option>
                </template>
              </Select>
            </FormItem>
          </template>
        </Form>
      </div>
      <div v-if="!disabled" class="drawer-footer">
        <Button style="margin-right: 8px" @click="handleCancel">{{ $t('cancel') }}</Button>
        <Button type="primary" class="primary" @click="handleSubmit">{{ $t('confirm') }}</Button>
      </div>
    </Drawer>
  </div>
</template>

<script>
import { debounce } from '@/pages/util'
export default {
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    value: {
      type: Object,
      default: () => {}
    },
    options: {
      type: Array,
      default: () => []
    },
    disabled: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    formData () {
      return this.value
    },
    drawerVisible: {
      get () {
        return this.visible
      },
      set (val) {
        this.$emit('update:visible', val)
      }
    }
  },
  data () {
    return {
      maxHeight: 500
    }
  },
  mounted () {
    this.maxHeight = document.body.clientHeight - 150
    window.addEventListener(
      'resize',
      debounce(() => {
        this.maxHeight = document.body.clientHeight - 150
      }, 100)
    )
  },
  methods: {
    handleSubmit () {
      this.$refs.form.validate(valid => {
        if (valid) {
          this.$emit('update:visible', false)
          this.$emit('submit')
        }
      })
    },
    handleCancel () {
      this.$emit('update:visible', false)
    }
  }
}
</script>

<style lang="scss" scoped>
.content {
  min-height: 500px;
  padding: 20px;
  overflow-y: auto;
}
.drawer-footer {
  width: 100%;
  position: absolute;
  bottom: 0;
  left: 0;
  border-top: 1px solid #e8e8e8;
  padding: 10px 16px;
  text-align: center;
  background: #fff;
}
</style>
<style lang="scss">
.workbench-edit-drawer {
  .ivu-tree-title-selected,
  .ivu-tree-title-selected:hover {
    background-color: transparent;
  }
}
</style>
