<template>
  <div class="new-publish-template">
    <div class="search">
      <Form :model="form" inline label-position="left">
        <FormItem label="模板名" :label-width="60">
          <Input v-model="form.templateName" @on-change="filterData" style="width:300px" placeholder="请输入模板名">
            <template #suffix>
              <Icon type="ios-search" />
            </template>
          </Input>
        </FormItem>
        <FormItem label="操作对象类型" :label-width="100">
          <Select
            v-model="form.operatorObjType"
            @on-change="filterData"
            clearable
            filterable
            placeholder="请选择操作对象类型"
            style="width:300px;"
          >
            <Option v-for="(item, index) in operateOptions" :value="item" :key="index">{{ item }}</Option>
          </Select>
        </FormItem>
      </Form>
    </div>
    <div class="wrapper">
      <div class="template">
        <Tabs v-model="activeName" @on-click="filterData" style="margin-bottom:10px;">
          <TabPane label="已发布模板" name="confirm"></TabPane>
          <TabPane label="我的草稿" name="created"></TabPane>
        </Tabs>
        <template v-if="cardList.length">
          <Card v-for="(i, index) in cardList" :key="index" style="width:100%;margin-bottom:20px;">
            <div class="w-header" slot="title">
              <Icon size="28" type="ios-people" />
              <div class="title">
                {{ i.manageRole }}
                <span class="underline"></span>
              </div>
              <Icon
                v-if="i.expand"
                size="28"
                type="md-arrow-dropdown"
                style="cursor:pointer;"
                @click="handleExpand(i)"
              />
              <Icon v-else size="28" type="md-arrow-dropright" style="cursor:pointer;" @click="handleExpand(i)" />
            </div>
            <div v-show="i.expand">
              <div v-for="j in i.groups" :key="j.groupId" class="content">
                <div class="sub-header">
                  <Icon size="24" type="ios-folder" />
                  <span class="title">{{ j.groupName }}</span>
                </div>
                <Table
                  @on-row-click="
                    row => {
                      handleChooseTemplate(row, row.role)
                    }
                  "
                  size="small"
                  :columns="tableColumns"
                  :data="j.templates"
                  style="margin:10px 0 20px 0"
                >
                </Table>
              </div>
            </div>
          </Card>
        </template>
        <div v-else class="template-no-data">
          暂无数据
        </div>
      </div>
      <div class="list">
        <Card style="width:520px;min-height:600px;">
          <div class="w-header">
            我的收藏 <span>{{ collectList.length }}</span>
          </div>
          <Table
            style="margin:20px 0;"
            :show-header="true"
            :disabled-hover="true"
            @on-row-click="
              row => {
                handleChooseTemplate(row, row.manageRole)
              }
            "
            size="small"
            :columns="collectColumns"
            :data="collectList"
          >
          </Table>
        </Card>
      </div>
    </div>
  </div>
</template>

<script>
import { getTemplateTree, collectTemplate, uncollectTemplate, collectTemplateList } from '@/api/server'
import { debounce } from '@/pages/util'
export default {
  data () {
    return {
      activeName: 'confirm', // confirm已发布，created我的草稿(未发布)
      type: '', // publish发布，request请求
      form: {
        templateName: '',
        operatorObjType: '' // 操作对象类型
      },
      operateOptions: [],
      // 收藏列表
      collectList: [],
      // 模板数据
      cardList: [],
      originCardList: [],
      tableColumns: [
        {
          title: '模板名称',
          key: 'name',
          minWidth: 250,
          render: (h, params) => {
            return (
              <div>
                {params.row.collectFlag === 0 && this.activeName === 'confirm' && (
                  <Tooltip content="收藏" placement="top-start">
                    <Icon
                      style="cursor:pointer;margin-right:5px;"
                      size="18"
                      type="ios-star-outline"
                      onClick={e => {
                        e.stopPropagation()
                        this.handleStar(params.row)
                      }}
                    />
                  </Tooltip>
                )}
                {params.row.collectFlag === 1 && this.activeName === 'confirm' && (
                  <Tooltip content="取消收藏" placement="top-start">
                    <Icon
                      style="cursor:pointer;margin-right:5px;"
                      size="18"
                      type="ios-star"
                      color="#ebac42"
                      onClick={e => {
                        e.stopPropagation()
                        this.handleStar(params.row)
                      }}
                    />
                  </Tooltip>
                )}
                <span style="margin-right:2px">
                  {params.row.name}
                  <Tag style="margin-left:5px">{params.row.version}</Tag>
                </span>
              </div>
            )
          }
        },
        {
          title: '操作对象类型',
          key: 'operatorObjType',
          render: (h, params) => {
            return params.row.operatorObjType && <Tag>{params.row.operatorObjType}</Tag>
          }
        },
        {
          title: '属主/角色',
          key: 'updatedBy',
          render: (h, params) => {
            return (
              <div style="display:flex;flex-direction:column">
                <span>{params.row.handler}</span>
                <span>{params.row.role}</span>
              </div>
            )
          }
        }
      ],
      collectColumns: [
        {
          title: '模板名称',
          key: 'name',
          width: 230,
          render: (h, params) => {
            return (
              <div>
                <Tooltip content="取消收藏" placement="top-start">
                  <Icon
                    style="cursor:pointer;margin-right:5px;"
                    size="18"
                    type="ios-star"
                    color="#ebac42"
                    onClick={e => {
                      e.stopPropagation()
                      this.handleStar({ ...params.row, collectFlag: 1 })
                    }}
                  />
                </Tooltip>
                <span class={{ 'active-link': params.row.status === 1, 'disabled-link': params.row.status !== 1 }}>
                  {params.row.name}
                  <Tag style="margin-left:5px">{params.row.version}</Tag>
                </span>
              </div>
            )
          }
        },
        {
          title: '模板状态',
          width: 120,
          key: 'status',
          render: (h, params) => {
            const list = [
              { label: '可使用', value: 1, color: '#19be6b' },
              { label: '已禁用', value: 2, color: '#c5c8ce' },
              { label: '权限被移除', value: 3, color: '#ed4014' }
            ]
            const item = list.find(i => i.value === params.row.status)
            return item && <Tag color={item.color}>{item.label}</Tag>
          }
        },
        {
          title: '使用角色',
          key: 'useRole'
        }
      ]
    }
  },
  watch: {
    // 解决同路由跳转不触发页面更新问题
    $route (to, from) {
      this.type = this.$route.query.type || ''
      this.getTemplateData()
      this.getCollectTemplate()
    }
  },
  mounted () {
    this.type = this.$route.query.type || ''
    this.getTemplateData()
    this.getCollectTemplate()
  },
  methods: {
    // 获取模板数据
    async getTemplateData () {
      this.operateOptions = []
      const { statusCode, data } = await getTemplateTree()
      const typeMap = {
        0: 'request',
        1: 'publish'
      }
      if (statusCode === 'OK') {
        this.cardList =
          Array.isArray(data) &&
          data.map(i => {
            i.expand = true
            i.groups.forEach(j => {
              j.templates.forEach(m => {
                if (m.operatorObjType && typeMap[m.type] === this.type) {
                  this.operateOptions.push(m.operatorObjType)
                }
              })
            })
            return i
          })
        this.originCardList = this.cardList
        // 数据去重
        this.operateOptions = Array.from(new Set(this.operateOptions))
        this.filterData()
      }
    },
    // 展开收缩卡片
    handleExpand (item) {
      item.expand = !item.expand
    },
    // 选中一条模板数据
    handleChooseTemplate (row, role) {
      if (row.status === 2) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: '该模板已禁用'
        })
      } else if (row.status === 3) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: '该模板使用权限被移除'
        })
      }
      const path = row.type === 0 ? 'createRequest' : 'createPublish'
      const url = `/taskman/workbench/${path}`
      this.$router.push({
        path: url,
        query: {
          requestTemplate: row.id,
          role: role,
          isAdd: 'Y',
          isCheck: 'N',
          isHandle: 'N',
          jumpFrom: ''
        }
      })
    },
    // 收藏or取消收藏模板
    handleStar: debounce(async function ({ id, collectFlag, role }) {
      const method = collectFlag ? uncollectTemplate : collectTemplate
      const params = collectFlag ? id : { templateId: id, role }
      const { statusCode } = await method(params)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.getTemplateData()
        this.getCollectTemplate()
      }
    }, 300),
    // 收藏模板列表
    async getCollectTemplate () {
      const params = {
        action: this.type === 'publish' ? 1 : 2, // 1发布2请求
        startIndex: 0,
        pageSize: 500
      }
      const { statusCode, data } = await collectTemplateList(params)
      if (statusCode === 'OK') {
        this.collectList = data.contents || []
      }
    },
    // 搜索过滤模板数据
    filterData () {
      const { templateName, operatorObjType } = this.form
      this.cardList = JSON.parse(JSON.stringify(this.originCardList))
      this.cardList = this.cardList.filter(i => {
        i.groups =
          (Array.isArray(i.groups) &&
            i.groups.filter(j => {
              j.templates =
                (Array.isArray(j.templates) &&
                  j.templates.filter(k => {
                    const typeMap = {
                      0: 'request',
                      1: 'publish'
                    }
                    // 根据模板名、标签名、模版发布状态组合搜索
                    const nameFilter = k.name.toLowerCase().indexOf(templateName.toLowerCase()) > -1
                    const operatorFilter = operatorObjType ? k.operatorObjType === operatorObjType : true
                    const typeFilter = this.type === typeMap[k.type]
                    return nameFilter && operatorFilter && typeFilter && this.activeName === k.status
                  })) ||
                []
              // 没有模板的组不显示
              return j.templates.length
            })) ||
          []
        // 没有组的角色不显示
        return i.groups.length
      })
    }
  }
}
</script>

<style lang="scss">
.new-publish-template {
  .ivu-card-head {
    border-bottom: 1px solid #e8eaec;
    padding: 5px 10px;
    line-height: 1;
  }
  .content .ivu-table-row {
    cursor: pointer;
  }
  .ivu-form-item {
    margin-bottom: 10px !important;
    display: inline-block !important;
  }
  .active-link:hover {
    margin-right: 2px;
    color: #2d8cf0;
    cursor: pointer;
  }
  .disabled-link:hover {
    margin-right: 2px;
    cursor: not-allowed;
  }
}
</style>
<style lang="scss" scoped>
.new-publish-template {
  .search {
    width: 100%;
    margin-bottom: 6px;
  }
  .wrapper {
    display: flex;
    .list {
      width: 520px;
      .w-header {
        font-size: 16px;
        font-weight: 700;
        span {
          font-size: 15px;
          font-weight: 400;
        }
      }
      .item {
        display: flex;
        padding: 5px 0 3px 5px;
        font-size: 14px;
        color: #515a6e;
        .collect {
          cursor: pointer;
          margin-left: 5px;
        }
        .template {
          width: 280px;
          text-overflow: ellipsis;
          overflow: hidden;
          // white-space: nowrap;
          &:hover {
            color: #5cadff;
            cursor: pointer;
          }
        }
      }
    }
    .template {
      padding: 0px 20px 0 0;
      width: calc(100% - 520px);
      .w-header {
        display: flex;
        align-items: center;
        .title {
          font-size: 16px;
          font-weight: bold;
          margin: 0 10px;
          .underline {
            display: block;
            margin-top: -10px;
            margin-left: -6px;
            width: 100%;
            padding: 0 6px;
            height: 12px;
            border-radius: 12px;
            background-color: #c6eafe;
            box-sizing: content-box;
          }
        }
      }
      .sub-header {
        display: flex;
        align-items: center;
        .title {
          font-size: 14px;
          font-weight: bold;
          margin-left: 5px;
        }
      }
    }
    .template-no-data {
      padding: 0px 20px 0 0;
      width: calc(100% - 320px);
      text-align: center;
      line-height: 300px;
      font-size: 16px;
    }
  }
}
</style>
