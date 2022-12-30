<template>
  <div class=" ">
    <Form :label-width="100">
      <FormItem :label="$t('root_entity')">
        <Select
          v-model="rootEntityId"
          filterable
          clearable
          :disabled="$parent.formDisable || $parent.jumpFrom === 'group_handle'"
          style="width:300px"
        >
          <Option v-for="item in rootEntityOptions" :value="item.guid" :key="item.guid">{{ item.key_name }}</Option>
        </Select>
        <Button
          @click="getEntityData"
          :disabled="$parent.formDisable || $parent.jumpFrom === 'group_handle'"
          type="primary"
          >{{ $t('search') }}</Button
        >
        <template v-if="!($parent.formDisable || $parent.jumpFrom === 'group_handle')">
          <Upload
            :action="uploadUrl"
            :show-upload-list="false"
            with-credentials
            style="display:inline-block;margin-left:32px"
            :headers="headers"
            :on-success="uploadSucess"
            :on-error="uploadFailed"
          >
            <Button type="success">{{ $t('upload_attachment') }}</Button>
          </Upload>
        </template>
        <div v-for="file in attachFiles" style="display:inline-block" :key="file.id">
          <Tag
            type="border"
            :closable="!($parent.formDisable || $parent.jumpFrom === 'group_handle')"
            checkable
            @on-close="removeFile(file)"
            @on-change="downloadFile(file)"
            color="primary"
            >{{ file.name }}</Tag
          >
        </div>
      </FormItem>
    </Form>
    <Tabs :value="activeTab" @on-click="changeTab">
      <template v-for="entity in requestData">
        <TabPane :label="entity.itemGroupName" :name="entity.entity || entity.itemGroup" :key="entity.itemGroup">
          <DataMgmt ref="dataMgmt" @getEntityData="getEntityData" @backData="backData"></DataMgmt>
        </TabPane>
      </template>
    </Tabs>
    <div style="text-align: center;margin-top:24px">
      <Button @click="saveData" v-if="!($parent.formDisable || $parent.jumpFrom === 'group_handle')" type="primary">{{
        $t('save')
      }}</Button>
      <Button @click="commitRequest" v-if="!($parent.formDisable || $parent.jumpFrom === 'group_handle')">{{
        $t('commit')
      }}</Button>
      <Button @click="nextStep" v-if="['', 'group_handle'].includes($parent.jumpFrom) && !$parent.isAdd">{{
        $t('next')
      }}</Button>
    </div>
  </div>
</template>

<script>
import DataMgmt from './data-mgmt'
import axios from 'axios'
import { getCookie } from '@/pages/util/cookie'
import {
  getRootEntity,
  getEntityData,
  deleteAttach,
  saveEntityData,
  getRequestInfo,
  updateRequestStatus
} from '@/api/server'
export default {
  name: '',
  data () {
    return {
      requestId: '',
      rootEntityId: '',
      rootEntityOptions: [],
      activeTab: '',
      attachFiles: [],
      requestData: [],
      uploadUrl: '',
      headers: {}
    }
  },
  components: {
    DataMgmt
  },
  mounted () {
    this.requestId = this.$parent.requestId
    this.uploadUrl = `/taskman/api/v1/request/attach-file/upload/${this.requestId}`
    const accessToken = getCookie('accessToken')
    this.headers = {
      Authorization: 'Bearer ' + accessToken
    }
    this.getEntity()
    this.getEntityData()
    if (this.$parent.requestId) {
      this.getRequestInfo()
    }
  },
  methods: {
    removeFile (file) {
      this.$Modal.confirm({
        title: this.$t('confirm_to_delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const { statusCode } = await deleteAttach(file.id)
          if (statusCode === 'OK') {
            this.getRequestInfo()
          }
        },
        onCancel: () => {}
      })
    },
    async downloadFile (file) {
      axios({
        method: 'GET',
        url: `/taskman/api/v1/request/attach-file/download/${file.id}`,
        headers: this.headers,
        responseType: 'blob'
      })
        .then(response => {
          if (response.status < 400) {
            let fileName = `${file.name}`
            let blob = new Blob([response.data])
            if ('msSaveOrOpenBlob' in navigator) {
              window.navigator.msSaveOrOpenBlob(blob, fileName)
            } else {
              if ('download' in document.createElement('a')) {
                // 非IE下载
                let elink = document.createElement('a')
                elink.download = fileName
                elink.style.display = 'none'
                elink.href = URL.createObjectURL(blob)
                document.body.appendChild(elink)
                elink.click()
                URL.revokeObjectURL(elink.href) // 释放URL 对象
                document.body.removeChild(elink)
              } else {
                // IE10+下载
                navigator.msSaveOrOpenBlob(blob, fileName)
              }
            }
          }
        })
        .catch(error => {
          console.log(error)
          this.$Message.warning('Error')
        })
    },
    uploadFailed (val, response) {
      console.log(val)
      this.$Notice.error({
        title: 'Error',
        desc: response.statusMessage
      })
    },
    async uploadSucess () {
      this.$Notice.success({
        title: 'Successful',
        desc: 'Successful'
      })
      this.getRequestInfo()
    },
    backData (data) {
      this.requestData = data
    },
    async commitRequest () {
      this.$Modal.confirm({
        title: this.$t('confirm') + this.$t('commit'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          this.confirmCommitRequest()
          // await this.saveData()
          // const { statusCode } = await updateRequestStatus(this.$parent.requestId, 'Pending')
          // if (statusCode === 'OK') {
          //   this.$router.push({ path: '/taskman/request-mgmt' })
          // }
        },
        onCancel: () => {}
      })
    },
    async confirmCommitRequest () {
      const params = {
        rootEntityId: this.rootEntityId,
        data: this.requestData
      }
      const result = this.paramsCheck()
      if (result) {
        const { statusCode } = await saveEntityData(this.requestId, params)
        if (statusCode === 'OK') {
          const { statusCode } = await updateRequestStatus(this.$parent.requestId, 'Pending')
          if (statusCode === 'OK') {
            this.$router.push({ path: '/taskman/request-mgmt' })
          }
        }
      } else {
        this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('required_tip')
        })
      }
    },
    async getRequestInfo () {
      const { statusCode, data } = await getRequestInfo(this.requestId)
      if (statusCode === 'OK') {
        this.attachFiles = data.attachFiles
        this.rootEntityId = data.cache
        this.getEntityData()
      }
    },
    nextStep () {
      if (!this.$parent.formDisable) {
        this.saveData()
      }
      this.$emit('nextStep')
    },
    async saveData () {
      const params = {
        rootEntityId: this.rootEntityId,
        data: this.requestData
      }
      const result = this.paramsCheck()
      if (result) {
        const { statusCode, data } = await saveEntityData(this.requestId, params)
        if (statusCode === 'OK') {
          this.requestData = data.data
          this.$Notice.success({
            title: this.$t('successful'),
            desc: this.$t('successful')
          })
        }
      } else {
        this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('required_tip')
        })
      }
    },
    paramsCheck () {
      let result = true
      this.requestData.forEach(requestData => {
        let requiredName = []
        requestData.title.forEach(t => {
          if (t.required === 'yes') {
            requiredName.push(t.name)
          }
        })
        requestData.value.forEach(v => {
          requiredName.forEach(key => {
            let val = v.entityData[key]
            if (Array.isArray(val)) {
              if (val.length === 0) {
                result = false
              }
            } else {
              if (val === '') {
                result = false
              }
            }
          })
        })
      })
      return result
    },
    async getEntity () {
      let params = {
        params: {
          requestId: this.requestId
        }
      }
      const { statusCode, data } = await getRootEntity(params)
      if (statusCode === 'OK') {
        this.rootEntityOptions = data.data
      }
    },
    async getEntityData () {
      let params = {
        params: {
          requestId: this.$parent.requestId,
          rootEntityId: this.rootEntityId
        }
      }
      const { statusCode, data } = await getEntityData(params)
      if (statusCode === 'OK') {
        this.activeTab = this.activeTab || data.data[0].entity
        this.requestData = data.data
        this.$nextTick(() => {
          const index = this.requestData.findIndex(r => r.entity === this.activeTab || r.itemGroup === this.activeTab)
          this.initTable(index)
        })
      }
    },
    changeTab (entity) {
      this.activeTab = entity
      const index = this.requestData.findIndex(r => r.entity === this.activeTab || r.itemGroup === this.activeTab)
      this.initTable(index)
    },
    initTable (index) {
      const find = this.requestData.find(r => r.entity === this.activeTab || r.itemGroup === this.activeTab)
      this.$refs.dataMgmt[index].initData(
        this.rootEntityId,
        this.requestData,
        find,
        this.requestId,
        this.$parent.formDisable,
        this.$parent.jumpFrom
      )
    }
  }
}
</script>

<style scoped lang="scss"></style>
