<template>
  <div class="cmdb-ci-password">
    <div style="display:flex;align-items:center;">
      <Tooltip
        max-width="200"
        class="ci-password-cell-show-span"
        placement="bottom-start"
        :content="isShowPassword ? realPassword : '******'"
      >
        <div class="password-wrapper">{{ isShowPassword ? realPassword : '******' }}</div>
      </Tooltip>
      <div style="float: right; margin-right: 12px;">
        <Icon :type="isShowPassword ? 'md-eye-off' : 'md-eye'" @click="showPassword" class="operation-icon-confirm" />
        <Icon type="ios-build-outline" v-if="!disabled" @click="resetPassword" class="operation-icon-confirm" />
      </div>
    </div>
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
        <Button @click="confirm" :loading="modalLoading" type="primary">{{
          useLocalValue ? $t('confirm') : $t('save')
        }}</Button>
        <Button @click="closeEditModal">{{ $t('tw_close') }}</Button>
      </div>
    </Modal>
  </div>
</template>

<script>
import { queryPassword, getEncryptKey } from '@/api/server'
import CryptoJS from 'crypto-js'
export default {
  name: '',
  data () {
    return {
      encryptKey: '',
      realPassword: '',
      useLocalValue: false,
      isShowPassword: false,

      isShowEditModal: false,
      editFormData: {
        newPassword: '',
        comparedPassword: ''
      },
      modalLoading: false,
      rules: {
        comparedPassword: [
          {
            message: this.$t('tw_please_input_right_new_password'),
            validator: () => this.editFormData.newPassword === this.editFormData.comparedPassword
          }
        ]
      }
    }
  },
  props: ['formData', 'panalData', 'disabled'],
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
      if (this.editFormData.newPassword) {
        await this.getEncryptKey()
        const key = CryptoJS.enc.Utf8.parse(this.encryptKey)
        const config = {
          iv: CryptoJS.enc.Utf8.parse(Math.trunc(new Date() / 100000) * 100000000),
          mode: CryptoJS.mode.CBC
          // padding: CryptoJS.pad.PKcs7
        }
        this.editFormData.newPassword = CryptoJS.AES.encrypt(this.editFormData.newPassword, key, config).toString()
      }
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
      if (this.useLocalValue || !this.panalData.guid) {
        this.realPassword = this.panalData[this.formData.propertyName]
      } else {
        const { statusCode, data } = await queryPassword(
          this.formData.ciTypeId,
          this.panalData.guid,
          this.formData.propertyName,
          {}
        )
        if (statusCode === 'OK') {
          this.realPassword = data
        }
      }
      this.isShowPassword = !this.isShowPassword
    },
    async getEncryptKey () {
      const { statusCode, data } = await getEncryptKey()
      if (statusCode === 'OK') {
        this.encryptKey = data
      }
    }
  },
  mounted () {}
}
</script>

<style scoped lang="scss">
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
  margin-right: 20px;
}
</style>
<style lang="scss">
.cmdb-ci-password .encrypt-password .ivu-input-suffix {
  display: none;
}
</style>
