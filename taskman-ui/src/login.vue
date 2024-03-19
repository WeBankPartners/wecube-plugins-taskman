<template>
  <div>
    <div class="body"></div>
    <div class="header-login">
      <div></div>
    </div>
    <br />
    <div class="login-form">
      <Input type="text" placeholder="username" v-model="username" name="user" @on-enter="login" />

      <Input
        type="password"
        password
        placeholder="password"
        v-model="password"
        name="password"
        @on-enter="login"
        style="margin-top: 20px"
      />
      <Button type="primary" long @click="login" :loading="loading" style="margin-top: 20px">
        Login
      </Button>
      <!-- <Button type="success" long>SUBMIT</Button> -->

      <Modal v-model="showRoleApply" :mask-closable="false" :closable="false" :title="$t('tw_apply_roles')">
        <Form ref="formValidate" :model="formValidate" :rules="ruleValidate" :label-width="80">
          <FormItem label="UM账号" prop="userName">
            <Input v-model="formValidate.userName" disabled></Input>
          </FormItem>
          <FormItem :label="$t('tw_email')" prop="emailAddr">
            <Input v-model="formValidate.emailAddr" :placeholder="$t('tw_email')"></Input>
          </FormItem>
          <FormItem :label="$t('manageRole')" prop="roleIds">
            <Select
              v-model="formValidate.roleIds"
              @on-open-change="getApplyRoles"
              multiple
              filterable
              :placeholder="$t('tw_apply_roles')"
            >
              <Option v-for="role in roleList" :value="role.id" :key="role.id">{{ role.displayName }}</Option>
            </Select>
          </FormItem>
        </Form>
        <div slot="footer">
          <Button @click="handleReset('formValidate')">{{ $t('cancel') }}</Button>
          <Button @click="handleSubmit('formValidate')" type="primary">{{ $t('tw_apply') }}</Button>
        </div>
      </Modal>
    </div>
  </div>
</template>
<script>
import { login, getApplyRoles, startApply } from '@/api/server'
export default {
  data () {
    return {
      username: '',
      password: '',
      loading: false,
      showRoleApply: false,
      formValidate: {
        userName: '',
        emailAddr: '',
        roleIds: []
      },
      ruleValidate: {
        emailAddr: [
          { required: true, message: `${this.$t('tw_email')} ${this.$t('can_not_be_empty')}`, trigger: 'blur' },
          { type: 'email', message: this.$t('tw_email_incorrect_format'), trigger: 'blur' }
        ],
        roleIds: [
          {
            required: true,
            type: 'array',
            min: 1,
            message: `${this.$t('manageRole')} ${this.$t('can_not_be_empty')}`,
            trigger: 'change'
          }
        ]
      },
      roleList: []
    }
  },
  methods: {
    async login () {
      if (!this.username || !this.password) return
      this.loading = true
      const payload = {
        username: this.username,
        password: this.password
      }
      const { status, data } = await login(payload)
      if (status === 'OK') {
        localStorage.setItem('username', this.username)
        const accessTokenObj = data.find(d => d.tokenType === 'accessToken')
        const refreshTokenObj = data.find(d => d.tokenType === 'refreshToken')
        localStorage.setItem('taskman-accessToken', accessTokenObj.token)
        localStorage.setItem('taskman-refreshToken', refreshTokenObj.token)
        localStorage.setItem('taskman-expiration', refreshTokenObj.expiration)
        const needRegister = data.needRegister || false
        if (needRegister) {
          this.showRoleApply = true
          this.formValidate.userName = this.username
        } else {
          this.$router.push('/taskman/workbench')
        }
      }
      this.loading = false
    },
    async getApplyRoles () {
      const params = {
        all: 'N', // Y:所有(包括未激活和已删除的) N:激活的
        roleAdmin: false
      }
      const { status, data } = await getApplyRoles(params)
      if (status === 'OK') {
        this.roleList = data || []
      }
    },
    handleSubmit (name) {
      this.$refs[name].validate(async valid => {
        if (valid) {
          const { status } = await startApply(this.formValidate)
          if (status === 'OK') {
            this.$Notice.success({
              title: this.$t('successful'),
              desc: this.$t('tw_apply_success')
            })
            this.showRoleApply = false
          }
        }
      })
    },
    handleReset (name) {
      this.$refs[name].resetFields()
      this.showRoleApply = false
    },
    clearSession () {
      let localStorage = window.localStorage
      localStorage.removeItem('username')
      window.needReLoad = true
    }
  },
  created () {
    this.clearSession()
  }
}
</script>
<style scoped>
.body {
  position: absolute;
  width: 100%;
  height: 100%;
  background-image: url('./images/bg.jpg');
  background-size: cover;
  -webkit-filter: blur(3px);
  z-index: 0;
}

.header-login {
  position: absolute;
  top: calc(50% - 35px);
  left: calc(50% - 355px);
  z-index: 2;
}

.header-login div {
  /* width: 600px;
  height: 50px;
  background-image: url('../assets/wecube-logo.png');
  background-size: contain;
  background-repeat: no-repeat; */
}

.login-form {
  position: absolute;
  top: calc(50% - 75px);
  left: calc(50% - 50px);
  height: 150px;
  width: 230px;
  padding: 10px;
  z-index: 2;
  text-align: center;
}
</style>
