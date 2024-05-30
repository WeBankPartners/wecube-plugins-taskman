<template>
  <div class="component-library-list">
    <Row>
      <Col :span="12" style="padding-right:5px;">
        <Input :value="getFormTypeDisplay" disabled placeholder="请选择表单类型" />
      </Col>
      <Col :span="12">
        <Select v-model="query.createdBy" clearable filterable placeholder="创建人" @on-change="handleSearch">
          <Option v-for="(i, index) in userList" :value="i.username" :key="index">{{ i.username }}</Option>
        </Select>
      </Col>
      <Col :span="24">
        <div class="group">
          <Input v-model.trim="query.name" clearable placeholder="组件名" class="input" @on-change="handleSearch" />
          <Icon type="ios-list" :size="36" class="icon" @click="handleViewTable" />
        </div>
      </Col>
    </Row>
    <!--拖拽列表-->
    <Card :bordered="false" dis-hover :padding="0" style="min-height:100px;">
      <draggable class="list" :list="data" :group="{ name: 'people', pull: 'clone', put: false }" :clone="cloneGroup">
        <Card v-for="i in data" :key="i.id" class="list-item">
          <div class="card-content">
            <span class="name">{{ i.name }}</span>
            <span class="person">
              <Icon type="md-person" size="16" color="#ff9900" style="margin-right:2px;" />
              {{ i.createdBy }}
            </span>
          </div>
        </Card>
        <div v-if="!loading && !data.length" class="no-data">
          {{ $t('tw_no_data') }}
        </div>
        <Spin fix v-if="loading"></Spin>
      </draggable>
    </Card>
    <Page
      v-show="data.length"
      class="list-page"
      :total="pagination.total"
      @on-change="handlePage"
      :current="pagination.currentPage"
      show-total
      size="small"
    />
    <ComponentLibraryModal
      ref="library"
      v-model="visible"
      :isAdd="false"
      @fetchList="handleSearch"
    ></ComponentLibraryModal>
  </div>
</template>
<script>
import ComponentLibraryModal from './component-library-modal.vue'
import draggable from 'vuedraggable'
import { getTemplateLibraryList, getLibraryFormTypeList, getAllUser } from '@/api/server'
import { debounce, deepClone } from '@/pages/util'
export default {
  components: {
    ComponentLibraryModal,
    draggable
  },
  props: {
    // 表单组类型
    groupType: {
      type: String,
      default: ''
    },
    // 表单类型
    formType: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      visible: false,
      query: {
        formType: '',
        createdBy: window.localStorage.getItem('username'),
        name: ''
      },
      pagination: {
        total: 0,
        currentPage: 1,
        pageSize: 10
      },
      data: [],
      loading: false,
      formTypeList: [],
      userList: []
    }
  },
  computed: {
    getFormTypeDisplay () {
      if (this.query.formType === 'requestInfo') {
        // 信息表单
        return this.$t('tw_information_form')
      } else {
        return this.query.formType
      }
    }
  },
  watch: {
    formType: {
      handler (val) {
        if (val) {
          if (this.groupType === 'custom') {
            this.query.formType = 'custom'
          } else {
            this.query.formType = val
          }
          this.getList()
        }
      },
      immediate: true
    }
  },
  mounted () {
    // this.getFormTypeList()
    this.getCreatedByList()
  },
  methods: {
    cloneGroup (val) {
      const arr = deepClone(val.items || [])
      // val.items[0].isActive = true
      arr.forEach(i => {
        // 拖拽添加的表单项，需要添加字段formItemLibrary
        i.formItemLibrary = i.id
        // 拖拽添加的表单项，前端自定义添加c_，保存的时候会清空c_开头的id
        i.id = 'c_' + i.id
      })
      return arr
    },
    // 获取表单类型下拉列表
    async getFormTypeList () {
      const { statusCode, data } = await getLibraryFormTypeList()
      if (statusCode === 'OK') {
        const arr = data || []
        this.formTypeList = arr.map(i => {
          return {
            label: i === 'requestInfo' ? this.$t('tw_information_form') : i,
            value: i
          }
        })
        // 将信息表单置于数组第一个
        const index = this.formTypeList.findIndex(i => i.value === 'requestInfo')
        const item = this.formTypeList.splice(index, 1)[0]
        this.formTypeList.unshift(item)
      }
    },
    // 获取创建人下拉列表
    async getCreatedByList () {
      const { status, data } = await getAllUser()
      if (status === 'OK') {
        this.userList = data || []
      }
    },
    async getList () {
      this.loading = true
      const params = {
        name: this.query.name,
        formType: this.query.formType,
        createdBy: this.query.createdBy,
        startIndex: (this.pagination.currentPage - 1) * this.pagination.pageSize,
        pageSize: this.pagination.pageSize
      }
      const { statusCode, data } = await getTemplateLibraryList(params)
      if (statusCode === 'OK') {
        this.data = data.contents || []
        this.pagination.total = data.pageInfo.totalRows
      }
      this.loading = false
    },
    handleSearch: debounce(function () {
      this.pagination.currentPage = 1
      this.getList()
    }, 300),
    handlePage (val) {
      this.pagination.currentPage = val
      this.getList()
    },
    handleViewTable () {
      this.visible = true
      this.$refs.library.init()
    }
  }
}
</script>

<style lang="scss">
.component-library-list {
  .ivu-card-body {
    padding: 0px;
  }
}
</style>
<style scoped lang="scss">
.component-library-list {
  width: 100%;
  .group {
    display: flex;
    align-items: center;
    margin-top: 10px;
    .input {
      flex: 1;
    }
    .icon {
      cursor: pointer;
    }
  }
  .list {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    &-item {
      display: flex;
      flex-direction: column;
      width: calc(50% - 5px);
      margin-top: 10px;
      height: 65px;
      cursor: pointer;
      .card-content {
        height: 65px;
        padding: 10px;
        font-size: 12px;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        .name {
          color: #2b85e4;
          max-width: 100%;
          text-overflow: ellipsis;
          overflow: hidden;
          white-space: nowrap;
          // display: -webkit-box;
          // -webkit-line-clamp: 2;
          // -webkit-box-orient: vertical;
        }
        .person {
          display: flex;
          align-items: flex-start;
        }
        .tag {
          width: fit-content;
          padding: 2px 6px;
          background: #5cadff;
          color: #fff;
          border-radius: 4px;
        }
      }
    }
    .no-data {
      width: 100%;
      height: 100px;
      display: flex;
      justify-content: center;
      align-items: center;
    }
    &-page {
      text-align: right;
      margin: 10px 5px;
      width: 100%;
    }
  }
}
</style>
