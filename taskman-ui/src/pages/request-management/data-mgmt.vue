<template>
  <div class=" ">
    <Button @click="addRow" type="primary">{{ $t('add') }}</Button>
    <Table size="small" :columns="tableColumns" :data="tableData"></Table>
    <Modal v-model="formConfig.isShow" :title="$t('data_management')">
      <Form :label-width="100">
        <template v-for="formItem in formConfig.form">
          <FormItem v-if="formItem.elementType === 'input'" :label="formItem.title" :key="formItem.name">
            <Input v-model="formConfig.data.entityData[formItem.name]"></Input>
          </FormItem>
          <FormItem v-if="formItem.elementType === 'select'" :label="formItem.title" :key="formItem.name">
            <Select v-model="formConfig.data.entityData[formItem.name]">
              <Option v-for="item in cityList" :value="item.value" :key="item.value">{{ item.label }}</Option>
            </Select>
          </FormItem>
        </template>
      </Form>
      <div slot="footer">
        <Button @click="cancel">{{ $t('cancel') }}</Button>
        <Button @click="ok" type="primary">{{ $t('save') }}</Button>
      </div>
    </Modal>
  </div>
</template>

<script>
import { saveEntityData } from '@/api/server'
export default {
  name: '',
  data () {
    return {
      model1: '',
      cityList: [],
      requestId: '',
      rootEntityId: '',
      dataArray: [], // 所有数据
      oriData: null,
      formConfig: {
        isShow: false,
        isAdd: false,
        form: [],
        data: {
          dataId: '',
          displayName: '',
          entityData: {},
          entityName: '',
          fullDataId: '',
          id: '',
          packageName: '',
          previousIds: [],
          succeedingIds: []
        }
      },
      tableColumns: [],
      tableData: []
    }
  },
  methods: {
    cancel () {
      this.formConfig.isShow = false
    },
    async ok () {
      this.formConfig.data.entityName = this.oriData.entity
      this.formConfig.data.packageName = this.oriData.packageName
      let find = this.dataArray.find(
        d => d.entity === this.oriData.entity && d.packageName === this.oriData.packageName
      )
      if (this.formConfig.isAdd) {
        this.formConfig.data.entityDataOp = 'create'
        find.value.push(this.formConfig.data)
      } else {
        this.formConfig.data.entityDataOp = 'update'
        let singleDataIndex = find.value.findIndex(v => v.id === this.formConfig.data.id)
        find.value[singleDataIndex] = this.formConfig.data
      }
      this.saveData(this.dataArray)
    },
    initData (rootEntityId, dataArray, data, requestId) {
      this.rootEntityId = rootEntityId
      this.dataArray = dataArray
      this.requestId = requestId
      this.oriData = data
      this.formConfig.form = data.title
      this.tableColumns = data.title.map(t => {
        const col = {
          title: t.title,
          key: t.name
        }
        this.formConfig.data.entityData[t.name] = ''
        return col
      })
      this.tableColumns.push({
        title: this.$t('action'),
        key: 'action',
        fixed: 'right',
        width: 160,
        render: (h, params) => {
          return (
            <div>
              <Button onClick={() => this.editRow(params.row)} size="small" type="primary">
                {this.$t('edit')}
              </Button>
              <Button style="margin-left: 4px" onClick={() => this.deleteRow(params.row)} size="small" type="error">
                {this.$t('delete')}
              </Button>
            </div>
          )
        }
      })
      this.tableData = data.value.map(v => {
        v.entityData._id = v.id
        return v.entityData
      })
    },
    addRow () {
      this.formConfig.data = {
        dataId: '',
        displayName: '',
        entityData: {},
        entityName: '',
        fullDataId: '',
        id: '',
        packageName: '',
        previousIds: [],
        succeedingIds: []
      }
      this.formConfig.isShow = true
      this.formConfig.isAdd = true
    },
    editRow (rowData) {
      const find = this.oriData.value.find(v => v.id === rowData._id)
      this.formConfig.data = { ...find }
      this.formConfig.isAdd = false
      this.formConfig.isShow = true
    },
    deleteRow (rowData) {
      this.$Modal.confirm({
        title: this.$t('confirm_to_delete'),
        content: name,
        onOk: async () => {
          let find = this.dataArray.find(
            d => d.entity === this.oriData.entity && d.packageName === this.oriData.packageName
          )
          let singleDataIndex = find.value.findIndex(v => v.id === rowData._id)
          find.value.splice(singleDataIndex, 1)
          console.log(find.value)
          this.saveData(this.dataArray)
        },
        onCancel: () => {}
      })
    },
    async saveData (data) {
      const params = {
        rootEntityId: this.rootEntityId,
        data: data
      }
      const { statusCode } = await saveEntityData(this.requestId, params)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.cancel()
        this.$emit('getEntityData')
      }
    }
  },
  components: {}
}
</script>

<style scoped lang="scss"></style>
