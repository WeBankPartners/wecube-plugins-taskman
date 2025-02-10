<template>
  <div class="cmdb-ci-password">
    <div v-if="formData.sensitive === 'yes' && isAdd" class="flex-center">
      <Tooltip
        max-width="200"
        class="ci-password-cell-show-span"
        placement="bottom-start"
        :content="getDisplayValue"
      >
        <div class="password-wrapper">{{ getDisplayValue }}</div>
      </Tooltip>
      <Button
        @click="showPassword"
        :disabled="getCmdbQueryPermission === false && panalData[formData.propertyName]"
        :icon="isShowPassword ? 'md-eye-off' : 'md-eye'"
      ></Button>
      <Button
        v-if="disabled === false"
        @click="resetPassword"
        type="primary"
        icon="md-create"
      ></Button>
    </div>
    <div v-else class="flex-center">
      <Tooltip
        max-width="200"
        class="ci-password-cell-show-span"
        placement="bottom-start"
        :content="panalData[formData.propertyName] || $t('tw_no_data')"
      >
        <div class="password-wrapper">{{ panalData[formData.propertyName] ? (isShowPassword ? panalData[formData.propertyName] : '******') : $t('tw_no_data') }}</div>
      </Tooltip>
      <Button
        @click="showPassword"
        :icon="isShowPassword ? 'md-eye-off' : 'md-eye'"
      ></Button>
      <Button
        v-if="disabled === false"
        @click="resetPassword"
        type="primary"
        icon="md-create"
      ></Button>
    </div>
    <!--密码编辑弹框-->
    <Modal v-model="isShowEditModal" :title="useLocalValue ? $t('tw_enter_password') : $t('tw_password_edit')">
      <Form ref="form" :model="editFormData" :rules="rules" label-position="right" :label-width="120">
        <FormItem :label="useLocalValue ? $t('tw_password') : $t('tw_new_password')" prop="newPassword">
          <Input
            class="encrypt-password"
            password
            :placeholder="$t('tw_new_password_input_placeholder')"
            ref="newPasswordInput"
            type="password"
            v-model="editFormData.newPassword"
          />
        </FormItem>
        <FormItem :label="$t('tw_confirm_password')" prop="comparedPassword">
          <Input
            class="encrypt-password"
            password
            :placeholder="$t('tw_please_input_new_password_again')"
            ref="comparedPasswordInput"
            type="password"
            v-model="editFormData.comparedPassword"
          />
        </FormItem>
      </Form>
      <div slot="footer">
        <Button @click="closeEditModal">{{ $t('tw_close') }}</Button>
        <Button @click="confirm" :loading="modalLoading" type="primary">{{
          useLocalValue ? $t('confirm') : $t('save')
        }}</Button>
      </div>
    </Modal>
  </div>
</template>

<script>
// import { getEncryptKey } from '@/api/server'
// import CryptoJS from 'crypto-js'
export default {
  name: '',
  data () {
    return {
      encryptKey: '',
      realPassword: '',
      originVal: '',
      useLocalValue: false,
      isShowPassword: false,

      isShowEditModal: false,
      editFormData: {
        newPassword: '',
        comparedPassword: ''
      },
      modalLoading: false,
      rules: {
        newPassword: [
          {
            required: true,
            message: this.$t('tw_new_password_input_placeholder'),
            validator: () => !!this.editFormData.newPassword
          }
        ],
        comparedPassword: [
          {
            required: true,
            message: this.$t('tw_new_password_input_placeholder'),
            validator: () => !!this.editFormData.comparedPassword
          },
          {
            required: true,
            message: this.$t('tw_please_input_right_new_password'),
            validator: () => this.editFormData.newPassword === this.editFormData.comparedPassword
          }
        ]
      }
    }
  },
  props: ['formData', 'panalData', 'allSensitiveData', 'rowData', 'disabled', 'isAdd'],
  computed: {
    getCmdbQueryPermission () {
      const obj = this.allSensitiveData.find(item => {
        if (this.rowData.dataId) {
          return item.attrName === this.formData.propertyName && item.guid === this.rowData.dataId
        } else {
          return item.attrName === this.formData.propertyName && item.tmpId === this.rowData.id
        }
      }) || {}
      return obj.queryPermission
    },
    getRealValue () {
      const obj = this.allSensitiveData.find(item => {
        if (this.rowData.dataId) {
          return item.attrName === this.formData.propertyName && item.guid === this.rowData.dataId
        } else {
          return item.attrName === this.formData.propertyName && item.tmpId === this.rowData.id
        }
      }) || {}
      return obj.value
    },
    getDisplayValue () {
      if (!this.panalData[this.formData.propertyName]) return this.$t('tw_no_data')
      if (this.isShowPassword) {
        if (this.originVal === this.panalData[this.formData.propertyName]) {
          return this.realPassword
        } else {
          return this.panalData[this.formData.propertyName]
        }
      } else {
        return '******'
      }
    }
  },
  mounted () {
    this.originVal = this.panalData[this.formData.propertyName]
  },
  methods: {
    resetPassword () {
      this.isShowEditModal = true
    },
    confirm () {
      this.$refs.form.validate(vail => {
        if (vail) {
          this.handleInput()
          this.useLocalValue = true
        }
      })
    },
    async handleInput () {
      // if (this.editFormData.newPassword) {
      //   await this.getEncryptKey()
      //   const key = CryptoJS.enc.Utf8.parse(this.encryptKey)
      //   const config = {
      //     iv: CryptoJS.enc.Utf8.parse(Math.trunc(new Date() / 100000) * 100000000),
      //     mode: CryptoJS.mode.CBC
      //     // padding: CryptoJS.pad.PKcs7
      //   }
      //   this.editFormData.newPassword = CryptoJS.AES.encrypt(this.editFormData.newPassword, key, config).toString()
      // }
      // this.panalData[this.formData.propertyName] = this.editFormData.newPassword
      this.$emit('input', this.editFormData.newPassword)
      this.realPassword = this.editFormData.newPassword
      this.editFormData = {
        newPassword: '',
        comparedPassword: ''
      }
      this.isShowEditModal = false
    },
    closeEditModal () {
      this.isShowEditModal = false
      this.editFormData = {
        newPassword: '',
        comparedPassword: ''
      }
    },
    async showPassword () {
      // this.realPassword = this.panalData[this.formData.propertyName]
      this.realPassword = this.getRealValue
      this.isShowPassword = !this.isShowPassword
    }
    // async getEncryptKey () {
    //   const { statusCode, data } = await getEncryptKey()
    //   if (statusCode === 'OK') {
    //     this.encryptKey = data
    //   }
    // }
  }
}
</script>

<style scoped lang="scss">
.flex-center {
  display: flex;
  align-items: center;
}
.operation-icon-confirm {
  font-size: 16px;
  border: 1px solid #57a3f3;
  color: #57a3f3;
  border-radius: 4px;
  width: 32px;
  line-height: 24px;
  cursor: pointer;
  margin-left: 5px;
}
.password-wrapper {
  text-overflow: ellipsis;
  overflow: hidden;
  width: fit-content;
  white-space: nowrap;
  margin-right: 5px;
}
.no-data-wrap {
  height: 34px;
  display: flex;
  align-items: center;
  .text {
    font-size: 13px;
    color: #515a6e;
  }
}
Button {
  margin-left: 5px;
}
</style>
<style lang="scss">
.cmdb-ci-password .encrypt-password .ivu-input-suffix {
  display: none;
}
</style>
