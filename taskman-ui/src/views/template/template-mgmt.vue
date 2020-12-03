<template>
  <div>
    <!-- <quillEditor></quillEditor> -->
    <!-- 表单设计 -->
    <Row>
      <Col span="4" style="border-right: 1px solid rgb(224, 223, 222);height: calc(100vh - 100px);overflow:auto">
        <Divider>拖动设置目标对象顺序</Divider>
        <div ref="entity" style="padding-right:10px;">
          <p style="font-size:16px;background:bisque;margin-bottom:5px;text-align: center" v-for="(entity,index) in entityList" :id="entity" :key="index">{{entity}}</p>
        </div>
        <Divider>自定义表单项</Divider>
        <Row ref="fields">
          <Col span="6" v-for="(comp, index) in componentsList" :id="index" :key="index">
            <div class="components-box">
              {{comp.label}}
            </div>
          </Col>
        </Row>
      </Col>
      <Col span="16" style="height: calc(100vh - 100px);overflow:auto;background: rgb(248, 238, 226);padding:20px">
        <Form :model="formItem" ref="form" :label-width="100">
          <TaskFormItem @delete="deleteHandler" @click.native="handleMouseClick(item)" @mouseenter.native="handleMouseEnter(item)" @mouseleave.native="handleMouseLeave(item)" v-for="(item, index) in formFields" :index="index" :item="item" :key="index"></TaskFormItem>
        </Form>
      </Col>
      <Col style="border-left: 1px solid rgb(224, 223, 222);height: calc(100vh - 100px);overflow:auto" span="4">
        33333
      </Col>
    </Row>
  </div>
</template>
<script>
// import quillEditor from "../../components/quillEditor"
import Sortable from 'sortablejs'

export default {
  components: {

  },
  data() {
    return {
      formItem: {
        input: ''
      },
      componentsList: [
        {
          label: "选择框",
          type: "Select",
        },
        {
          label: "输入框",
          type: "Input",
        },
        {
          label: "多行文本",
          type: "Textarea",
        },
        {
          label: "富文本",
          type: "QuillEditor"
        },
      ],
      formFields: [
        {
          component: "Input",
          colSpan: 6,
          name: "test",
          label: "测试1",
          defaultValue: "aaaa",
          entity: "unit",
          isHover: false,
          isActive: false
        },
        {
          component: "Select",
          colSpan: 6,
          name: "test",
          label: "测试2",
          defaultValue: "aaaa",
          entity: "unit",
          isHover: false,
          isActive: false
        },
        {
          component: "Select",
          colSpan: 24,
          name: "test",
          label: "测试3",
          defaultValue: "aaaa",
          entity: "network",
          isHover: false,
          isActive: false
        },
        {
          component: "Select",
          colSpan: 24,
          name: "test",
          label: "测试4",
          defaultValue: "aaaa",
          entity: "datacenter",
          isHover: false,
          isActive: false
        },
        {
          component: "Input",
          colSpan: 24,
          name: "test",
          label: "测试5",
          defaultValue: "aaaa",
          entity: "datacenter",
          isHover: false,
          isActive: false
        },
        {
          component: "Select",
          colSpan: 24,
          name: "test",
          label: "测试6",
          defaultValue: "aaaa",
          entity: "datacenter",
          isHover: false,
          isActive: false
        },
        {
          component: "Select",
          colSpan: 24,
          name: "test",
          label: "测试7",
          defaultValue: "aaaa",
          entity: "datacenter",
          isHover: false,
          isActive: false
        },
        {
          component: "Select",
          colSpan: 24,
          name: "test",
          label: "测试8",
          defaultValue: "aaaa",
          entity: "unit",
          isHover: false,
          isActive: false
        },
      ],
      entityList: ['unit', 'datacenter', 'network'],
      list:['unit', 'datacenter', 'network'],
    }
  },
  methods: {
    handleMouseEnter (item) {
      item.isHover = true
    },
    handleMouseLeave (item) {
      item.isHover = false
    },
    handleMouseClick (item) {
      if(item.isActive) {
        item.isActive = false
        return
      }
      this.formFields.forEach(formField => {formField.isActive = false})
      item.isActive = true
    },
    deleteHandler (index) {
      console.log(index)
    },
    createSortable(el, items) {
      return new Sortable(el, {
        group: {
          name: 'component',
          pull: 'clone',
          put: false
        },
        sort: false,
        animation: 150,
        setData: (dataTransfer, dragEl) => {
          console.log(dataTransfer, dragEl)
          // const index = parseInt(dragEl.dataset.id)
          // dragEl.__item__ = items[index]
        }
      })
    },
    createEntitySortable(el, items) {
      return new Sortable(el, {
        group: {
          name: 'component',
          pull: 'clone',
          put: false
        },
        animation: 150,
        onUpdate: event => {
          const newIndex = event.newIndex
          const oldIndex = event.oldIndex
          const $li = el.children[newIndex]
          const $oldLi = el.children[oldIndex]
          // el.removeChild($li)
          // if (newIndex > oldIndex) {
          //   el.insertBefore($li, $oldLi)
          // } else {
          //   el.insertBefore($li, $oldLi.nextSibling)
          // }
          items[event.oldIndex] = items.splice(event.newIndex, 1, items[event.oldIndex])[0]
          console.log(this.entityList,event,items)
        }
      })
    },
    createFormSortable(el) {
      this.sortable = new Sortable(el, {
        group: 'component',
        animation: 150,
        onStart: () => {
          this.dragging = true
          this.showHelper = null
        },
        onEnd: (e) => {
          this.dragging = false
        },
        onSort: (e) => {
          // 添加也会触发onSort， 用个变量去来区分
          console.log(this.formFields,e)
          // if (!isAdd) {
          //   // this.$store.commit('sortFields', e)
          // }
          // isAdd = false
        },
        onAdd: (e) => {
          console.log(1,e)
          const id = e.item.attributes.id.value
          e.item.parentNode.removeChild(e.item)
          
          // if (item) {
          //   isAdd = true
          //   // this.handleDrop(item, e.newIndex)
          // }
          this.$nextTick(() => {
            const item = {
              component: this.componentsList[id*1].type,
              colSpan: 24,
              name: "test",
              label: "测试1aaaa",
              defaultValue: "aaaa",
              entity: null,
              isHover: false,
              isActive: false
            }
            this.formFields.splice(e.newIndex, 0, item)
          })
        }
      })
    }
  },
  mounted() {
    this.fieldsSortable = this.createSortable(this.$refs.fields.$el, this.componentsList)
    this.formSortable = this.createFormSortable(this.$refs.form.$el)
    this.entitySortable = this.createEntitySortable(this.$refs.entity, this.list)
  },
  beforeDestroy() {
    this.fieldsSortable && this.fieldsSortable.destroy()
    this.formSortable && this.formSortable.destroy()
    this.entitySortable && this.entitySortable.destroy()
  }
}
</script>
<style lang="scss">
.components-box {
    background: bisque;
    width: 60px;
    height: 60px;
    text-align: center;
    line-height: 60px;
    cursor: move;
  }
.mask {
  height: 100%;
}
</style>
