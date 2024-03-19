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
    </div>
  </div>
</template>
<script>
import { login } from '@/api/server'
export default {
  data () {
    return {
      username: '',
      password: '',
      loading: false
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
      console.log(33, payload)
      const { status, data } = await login(payload)
      if (status === 'OK') {
        console.log(data)
        localStorage.setItem('username', this.username)
        const accessTokenObj = data.find(d => d.tokenType === 'accessToken')
        const refreshTokenObj = data.find(d => d.tokenType === 'refreshToken')
        localStorage.setItem('taskman-accessToken', accessTokenObj.token)
        localStorage.setItem('taskman-refreshToken', refreshTokenObj.token)
        localStorage.setItem('taskman-expiration', refreshTokenObj.expiration)
        this.$router.push('/taskman/workbench')
      }
      this.loading = false
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
