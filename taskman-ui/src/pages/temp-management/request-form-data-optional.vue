<template>
  <div>
    <Row>
      <Col span="6" style="border: 1px solid #dcdee2; padding: 0 16px">
        <div :style="{ height: MODALHEIGHT + 32 + 'px', overflow: 'auto' }">
          <Divider plain>{{ $t('custom_form') }}</Divider>
          <draggable
            class="dragArea"
            :list="customElement"
            :group="{ name: 'people', pull: 'clone', put: false }"
            :sort="$parent.isCheck !== 'Y'"
            :clone="cloneDog"
          >
            <div class="list-group-item-" style="width: 100%" v-for="element in customElement" :key="element.id">
              <Input v-if="element.elementType === 'input'" :placeholder="$t('t_input')" />
              <Input v-if="element.elementType === 'textarea'" type="textarea" :placeholder="$t('textare')" />
              <Select v-if="element.elementType === 'select'" :placeholder="$t('select')"></Select>
              <Select v-if="element.elementType === 'wecmdbEntity'" placeholder="模型数据项"></Select>
              <div v-if="element.elementType === 'group'" style="width: 100%; height: 80px; border: 1px solid #5ea7f4">
                <span style="margin: 8px; color: #bbbbbb"> Item Group </span>
              </div>
            </div>
          </draggable>
        </div>
      </Col>
      <Col span="12" style="border: 1px solid #dcdee2; padding: 0 16px; width: 48%; margin: 0 4px">
        <div :style="{ height: MODALHEIGHT + 30 + 'px', overflow: 'auto' }">
          <Divider>预览</Divider>
          <div class="title">
            <div class="title-text">
              {{ $t('root_entity') }}
              <span class="underline"></span>
            </div>
          </div>
          <Form :label-width="120">
            <FormItem :label="$t('tw_choose_object')">
              <Select style="width: 30%">
                <Option v-for="item in []" :value="item.id" :key="item.id">{{ item.displayName }}</Option>
              </Select>
            </FormItem>
          </Form>
          <div>
            <Icon @click="selectItemGroup" style="cursor: pointer;" type="md-add-circle" size="24" color="#2d8cf0" />
          </div>
          <!-- <template v-for="(item, itemIndex) in finalElement">
            <div :key="itemIndex" style="border: 1px solid #dcdee2; margin-bottom: 8px; padding: 8px">
              <span style="font-weight: 600;">
                {{ item.itemGroupName }}
              </span>
              <draggable
                class="dragArea"
                :list="item.attrs"
                :sort="$parent.isCheck !== 'Y'"
                group="people"
                @change="log(item)"
              >
                <div
                  @click="selectElement(itemIndex, eleIndex)"
                  :class="['list-group-item-', element.isActive ? 'active-zone' : '']"
                  :style="{ width: (element.width / 24) * 100 + '%' }"
                  v-for="(element, eleIndex) in item.attrs"
                  :key="element.id"
                >
                  <div>
                    <Icon v-if="element.required === 'yes'" size="8" style="color: #ed4014" type="ios-medical" />
                    {{ element.title }}:
                  </div>
                  <Input
                    v-if="element.elementType === 'input'"
                    :disabled="element.isEdit === 'no'"
                    v-model="element.defaultValue"
                    placeholder=""
                    style="width: calc(100% - 30px)"
                  />
                  <Input
                    v-if="element.elementType === 'textarea'"
                    :disabled="element.isEdit === 'no'"
                    v-model="element.defaultValue"
                    type="textarea"
                    style="width: calc(100% - 30px)"
                  />
                  <Select
                    v-if="element.elementType === 'select'"
                    :disabled="element.isEdit === 'no'"
                    v-model="element.defaultValue"
                    style="width: calc(100% - 30px)"
                  ></Select>
                  <Select
                    v-if="element.elementType === 'wecmdbEntity'"
                    :disabled="element.isEdit === 'no'"
                    v-model="element.defaultValue"
                    style="width: calc(100% - 30px)"
                  ></Select>
                  <Button
                    @click.stop="removeForm(itemIndex, eleIndex, element)"
                    type="error"
                    size="small"
                    :disabled="$parent.isCheck === 'Y'"
                    ghost
                    icon="ios-close"
                  ></Button>
                </div>
              </draggable>
            </div>
          </template> -->
        </div>
      </Col>
    </Row>
    <Modal v-model="showSelectModel" title="选择组信息" :mask-closable="false">
      <Form :label-width="120">
        <FormItem :label="$t('选择组')">
          <Select style="width: 80%" v-model="selectGroup" filterable>
            <OptionGroup v-for="itemGroup in groupOptions" :label="itemGroup.formType" :key="itemGroup.formType">
              <Option v-for="item in itemGroup.entities" :value="item" :key="item">{{ item }}</Option>
            </OptionGroup>
          </Select>
        </FormItem>
      </Form>
      <template #footer>
        <Button @click="cancelSelect">{{ $t('cancel') }}</Button>
        <Button @click="okSelect" :disabled="selectGroup === ''" type="primary">{{ $t('confirm') }}</Button>
      </template>
    </Modal>
  </div>
</template>

<script>
import { getEntityByTemplateId } from '@/api/server.js'
import draggable from 'vuedraggable'
let idGlobal = 18
export default {
  name: 'form-select',
  data () {
    return {
      isParmasChanged: false, // 参数变化标志位
      MODALHEIGHT: 200,
      customElement: [
        // 预制自定义表单项目
        {
          id: 1,
          name: 'input',
          title: 'Input',
          elementType: 'input',
          defaultValue: '',
          defaultClear: 'no',
          // tag: '',
          itemGroup: '',
          itemGroupName: '',
          packageName: '',
          entity: '',
          width: 24,
          dataOptions: '',
          regular: '',
          inDisplayName: 'yes',
          isEdit: 'yes',
          multiple: 'N',
          selectList: [],
          isRefInside: 'no',
          required: 'no',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: '',
          refEntity: '',
          refPackageName: ''
        },
        {
          id: 3,
          name: 'textarea',
          title: 'Textarea',
          elementType: 'textarea',
          defaultClear: 'no',
          defaultValue: '',
          // tag: '',
          itemGroup: '',
          itemGroupName: '',
          packageName: '',
          entity: '',
          width: 24,
          dataOptions: '',
          regular: '',
          inDisplayName: 'yes',
          isEdit: 'yes',
          multiple: 'N',
          selectList: [],
          isRefInside: 'no',
          required: 'no',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: '',
          refEntity: '',
          refPackageName: ''
        },
        {
          id: 2,
          name: 'select',
          title: 'Select',
          elementType: 'select',
          defaultValue: '',
          defaultClear: 'no',
          // tag: '',
          itemGroup: '',
          itemGroupName: '',
          packageName: '',
          entity: '',
          width: 24,
          dataOptions: '',
          regular: '',
          inDisplayName: 'yes',
          isEdit: 'yes',
          multiple: 'N',
          selectList: [],
          isRefInside: 'no',
          required: 'no',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: '',
          refEntity: '',
          refPackageName: ''
        },
        {
          id: 5,
          name: 'wecmdbEntity',
          title: 'WecmdbEntity',
          elementType: 'wecmdbEntity',
          defaultValue: '',
          defaultClear: 'no',
          // tag: '',
          itemGroup: '',
          itemGroupName: '',
          packageName: '',
          entity: '',
          width: 24,
          dataOptions: '',
          regular: '',
          inDisplayName: 'yes',
          isEdit: 'yes',
          multiple: 'N',
          selectList: [],
          isRefInside: 'no',
          required: 'no',
          isView: 'yes',
          isOutput: 'no',
          sort: 0,
          attrDefId: '',
          attrDefName: '',
          attrDefDataType: '',
          refEntity: '',
          refPackageName: ''
        }
      ],
      showSelectModel: false, // 组选择框
      selectGroup: '', // 选中的组信息
      groupOptions: []
    }
  },
  props: ['requestTemplateId'],
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 400
  },
  methods: {
    async loadPage () {
      this.isParmasChanged = false
      // const { statusCode, data } = await getRequestFormTemplateData(this.requestTemplateId)
      // if (statusCode === 'OK') {
      //   if (data.items !== null && data.items.length > 0) {
      //     data.items.sort(this.compare('sort'))
      //     this.finalElement = [{
      //       itemGroup: data.items[0].itemGroup,
      //       itemGroupName: data.items[0].itemGroupName,
      //       attrs: data.items
      //     }]
      //   }
      // }
    },
    cloneDog (val) {
      if (this.$parent.isCheck === 'Y') return
      let newItem = JSON.parse(JSON.stringify(val))
      newItem.id = 'c_' + idGlobal++
      newItem.title = newItem.title + idGlobal
      newItem.name = newItem.name + idGlobal
      this.specialId = newItem.id
      this.paramsChanged()
      return newItem
    },
    log (item) {
      this.finalElement.forEach(l => {
        l.attrs.forEach(attr => {
          attr.itemGroup = l.itemGroup
          attr.itemGroupName = l.itemGroupName
          if (attr.id === this.specialId) {
            this.editElement = attr
            this.openPanel = '1'
          }
        })
      })
    },
    async selectItemGroup () {
      this.selectGroup = ''
      this.groupOptions = []
      const { statusCode, data } = await getEntityByTemplateId(this.requestTemplateId)
      if (statusCode === 'OK') {
        // workflow  2.编排数据, 3.optional 自选数据项表单,  custom 1.自定义表单
        this.$nextTick(() => {
          this.groupOptions = data.map(d => {
            if (d.formType === 'custom') {
              d.entities = ['custom']
            } else {
              d.entities = d.entities || []
            }
            return d
          })
          this.showSelectModel = true
        })
      }
    },
    okSelect () {
      console.log(this.selectGroup)
      this.showSelectModel = false
    },
    cancelSelect () {
      this.showSelectModel = false
    }
  },
  components: {
    draggable
  }
}
</script>
<style>
.ivu-input[disabled],
fieldset[disabled] .ivu-input {
  color: #757575 !important;
}
.ivu-select-input[disabled] {
  color: #757575 !important;
  -webkit-text-fill-color: #757575 !important;
}
.ivu-select-disabled .ivu-select-selection {
  color: #757575 !important;
}
</style>
<style scoped lang="scss">
.active-zone {
  color: #338cf0;
}
.ivu-form-item {
  margin-bottom: 8px;
}
.list-group-item- {
  display: inline-block;
  margin: 2px 0;
}
.title {
  font-size: 12px;
  font-weight: bold;
  margin: 0 10px;
  display: inline-block;
  .title-text {
    display: inline-block;
    margin-left: 6px;
  }
  .underline {
    display: block;
    margin-top: -10px;
    margin-left: -6px;
    width: 100%;
    padding: 0 6px;
    height: 8px;
    border-radius: 8px;
    background-color: #c6eafe;
    box-sizing: content-box;
  }
}
</style>
