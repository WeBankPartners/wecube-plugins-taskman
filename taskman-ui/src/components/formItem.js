import "./formItem.scss"
export default {
  name: "TaskFormItem",
  data () {
    return {
      itemValue: ''
    }
  },
  props: {
    item: {},
    index: {},
    isDesign: {},
    value: {}
  },
  computed: {
    classList () {
      let classes = this.isDesign ? 'form-item-taskman form-item-taskman-isDesign ' : 'form-item-taskman '
      if (this.item.isHover) {
        classes = classes+'form-item-taskman-hover'
      }
      if (this.item.isActive) {
        classes = classes+'form-item-taskman-active'
      }
      return classes
    }
  },
  watch: {
    value(val) {
      this.itemValue = val
    }
  },
  methods: {
   deleteHandler () {
     this.$emit('delete', this.index)
   },
   inputHandler (v) {
    this.itemValue = v
    this.$emit('input',v)
   }
  },
  mounted () {
    this.itemValue = this.item.defaultValue
  },
  render () {
    const { elementType, width, title, options, value, defaultValue, placeholder,isHover, isActive, entity, refEntity,isCustom, isMultiple } = this.item
    return (
      <Col id={'formItem_'+this.index} span={width}>
        <FormItem label-width={200} class={this.classList} label={title}>
          <elementType options={options} entityName = {refEntity} isMultiple={isMultiple} value={this.itemValue} onInput={v => this.inputHandler(v)} placeholder={placeholder}></elementType>
          {(this.isDesign && isActive) && <Button disabled={!isCustom} onClick={() => this.deleteHandler()} class="deleteButton" size="small" ghost type="error" icon="ios-trash-outline"></Button>}
        </FormItem>
      </Col>
    )
  }
}