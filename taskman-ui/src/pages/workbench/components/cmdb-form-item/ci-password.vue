<template>
  <div class="cmdb-ci-password">
    <div v-if="realPassword" style="display:flex;align-items:center;">
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
    <div v-else class="no-data-wrap">
      <span class="text">{{ $t('tw_password_empty') }}</span>
      <Icon type="ios-build-outline" v-if="!disabled" @click="resetPassword" class="operation-icon-confirm" />
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
import { getEncryptKey } from '@/api/server'
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
  props: ['formData', 'panalData', 'disabled'],
  mounted () {
    this.realPassword = this.panalData[this.formData.propertyName]
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
      // if (this.useLocalValue || !this.panalData.guid) {
      //   this.realPassword = this.panalData[this.formData.propertyName]
      // } else {
      //   const { statusCode, data } = await queryPassword(
      //     this.formData.ciTypeId,
      //     this.panalData.guid,
      //     this.formData.propertyName,
      //     {}
      //   )
      //   if (statusCode === 'OK') {
      //     this.realPassword = data
      //   }
      // }
      this.realPassword = this.panalData[this.formData.propertyName]
      this.isShowPassword = !this.isShowPassword
    },
    async getEncryptKey () {
      const { statusCode, data } = await getEncryptKey()
      if (statusCode === 'OK') {
        this.encryptKey = data
      }
    }
  }
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
</style>
<style lang="scss">
.cmdb-ci-password .encrypt-password .ivu-input-suffix {
  display: none;
}
</style>
