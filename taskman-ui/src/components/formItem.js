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
    index: {}
  },
  computed: {
    classList () {
      let classes = 'form-item-taskman '
      if (this.item.isHover) {
        classes = 'form-item-taskman form-item-taskman-hover'
      }
      if (this.item.isActive) {
        classes = 'form-item-taskman form-item-taskman-active'
      }
      return classes
    }
  },
  methods: {
   deleteHandler () {
     this.$emit('delete', this.index)
   }
  },
  mounted () {
    this.itemValue = this.item.defaultValue
  },
  render () {
    const { component, colSpan, label, options, value, defaultValue, placeholder,isHover, isActive } = this.item
    return (
      <Col id={'formItem_'+this.index} span={colSpan}>
        <FormItem class={this.classList} label={label}>
          <component value={this.itemValue} placeholder={placeholder}></component>
          {(isHover || isActive) && <Button onClick={() => this.deleteHandler()} class="deleteButton" size="small" ghost type="error" icon="ios-trash-outline"></Button>}
        </FormItem>
      </Col>
    )
  }
}