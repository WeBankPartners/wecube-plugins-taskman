<template>
  <div class="workbench-data-card">
    <Card
      v-for="(i, index) in data"
      :key="index"
      border
      style="width:250px;margin-right:20px;cursor:pointer;"
      @click="fetchData(i)"
    >
      <div class="content">
        <div class="header">
          <img :src="i.icon" />
          <span>{{ `${i.label}(${i.total})` }}</span>
        </div>
        <div class="data">
          <div class="list">
            <span>{{ i.publishNum }}</span>
            <span>发布</span>
          </div>
          <div class="list">
            <span>{{ i.requestNum }}</span>
            <span>请求</span>
          </div>
        </div>
      </div>
    </Card>
  </div>
</template>

<script>
import { overviewData } from '@/api/server'
export default {
  data() {
    return {
      data: [
        {
          type: 'pending',
          label: '待处理的',
          icon: require('@/images/warning.png'),
          requestNum: 10,
          publishNum: 5
        },
        {
          type: 'hasProcessed',
          label: '已处理的',
          icon: require('@/images/checked.png'),
          requestNum: 10,
          publishNum: 5
        },
        {
          type: 'submit',
          label: '我提交的',
          icon: require('@/images/upload.png'),
          requestNum: 10,
          publishNum: 5
        },
        {
          type: 'draft',
          label: '我保存的',
          icon: require('@/images/save.png'),
          requestNum: 10,
          publishNum: 5
        },
        {
          type: 'collect',
          label: '收藏模板',
          icon: require('@/images/star.png'),
          requestNum: 10,
          publishNum: 5
        }
      ]
    }
  },
  mounted() {
    this.getData()
  },
  methods: {
    async getData() {
      const { statusCode, data } = await overviewData() 
      if (statusCode === 'OK') {
        for (let key in data) {
          this.data.forEach(i => {
            if (i.type === key) {
              const numArr = data[key] && data[key].split(';') || []
              i.publishNum = numArr[0]
              i.requestNum = numArr[1]
              i.total = Number(i.publishNum) + Number(i.requestNum)
            }
          })
        }
      }
    },
    fetchData({ type }) {
      this.$emit('fetchData', type)
    }
  }
}
</script>

<style lang="scss">
.workbench-data-card .ivu-card-body {
  padding: 10px;
}

</style>
<style lang="scss" scoped>
.workbench-data-card {
  display: flex;
  .content {
    .header {
      display: flex;
      align-items: flex-end;
      img {
        width: 24px;
        height: 24px;
        margin-right: 5px;
      }
      span {
        color: rgba(16, 16, 16, 1);
        font-size: 14px;
        font-family: PingFangSC-regular;
        font-weight: 500;
      }
    }
    .data {
      display: flex;
      justify-content: space-around;
      align-items: center;
      margin-top: 10px;
      .list {
        display: flex;
        flex-direction: column;
        align-items: center;
        span {
          color: rgba(16, 16, 16, 1);
          font-size: 14px;
          font-family: PingFangSC-regular;
        }
      }
    }
  }
}
</style>