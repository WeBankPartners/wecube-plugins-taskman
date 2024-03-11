<template>
  <div class="workbench-upload">
    <Upload
      v-if="!onlyShowFile"
      :action="uploadUrl"
      :before-upload="handleUpload"
      :show-upload-list="false"
      with-credentials
      :max-size="10240"
      :on-exceeded-size="handleMaxSize"
      :headers="headers"
      :on-success="uploadSucess"
      :on-error="uploadFailed"
    >
      <Button icon="md-cloud-upload" :disabled="formDisable">{{ $t('upload_attachment') }}</Button>
    </Upload>
    <div :style="{ marginTop: onlyShowFile ? '0px' : '10px', display: 'flex' }">
      <Tag
        v-for="file in attachFiles"
        :key="file.id"
        type="border"
        :closable="!formDisable"
        checkable
        @on-close="removeFile(file)"
        @on-change="downloadFile(file)"
        color="primary"
        style="margin-right:15px;"
        >{{ file.name }}</Tag
      >
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import { deleteAttach } from '@/api/server'
import { getCookie } from '@/pages/util/cookie'
export default {
  props: {
    // request请求，task任务
    type: {
      type: String,
      default: 'request'
    },
    id: {
      type: String,
      default: ''
    },
    taskHandleId: {
      type: String,
      default: ''
    },
    formDisable: {
      type: Boolean,
      default: false
    },
    files: {
      type: Array,
      default: () => []
    },
    onlyShowFile: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      uploadUrl: '',
      headers: {},
      attachFiles: []
    }
  },
  watch: {
    id: {
      handler (val) {
        if (val && this.type === 'request') {
          this.uploadUrl = `/taskman/api/v1/request/attach-file/upload/${val}`
          const accessToken = getCookie('accessToken')
          this.headers = {
            Authorization: 'Bearer ' + accessToken
          }
        }
        if (val && this.taskHandleId && this.type === 'task') {
          this.uploadUrl = `/taskman/api/v1/task/attach-file/${val}/upload/${this.taskHandleId}`
          const accessToken = getCookie('accessToken')
          this.headers = {
            Authorization: 'Bearer ' + accessToken
          }
        }
      },
      immediate: true
    },
    taskHandleId: {
      handler (val) {
        if (val && this.id && this.type === 'task') {
          this.uploadUrl = `/taskman/api/v1/task/attach-file/${this.id}/upload/${val}`
          const accessToken = getCookie('accessToken')
          this.headers = {
            Authorization: 'Bearer ' + accessToken
          }
        }
      },
      immediate: true
    },
    files: {
      handler (val) {
        this.attachFiles = val
      },
      deep: true,
      immediate: true
    }
  },
  mounted () {},
  methods: {
    handleUpload (file) {
      this.$Message.info(this.$t('upload_tip'))
      return true
    },
    removeFile (file) {
      this.$Modal.confirm({
        title: this.$t('confirm_to_delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          this.$Modal.remove()
          const { statusCode, data } = await deleteAttach(file.id)
          if (statusCode === 'OK') {
            this.attachFiles = data
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
    async uploadSucess (item) {
      this.$Notice.success({
        title: 'Successful',
        desc: 'Successful'
      })
      this.attachFiles = item.data
    },
    handleMaxSize () {
      this.$Notice.error({
        title: 'Error',
        desc: this.$t('tw_upload_10_error')
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.workbench-upload {
  width: 100%;
  .file-list {
    display: flex;
  }
}
</style>
