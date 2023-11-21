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
        <FormItem label="标签" :label-width="50">
          <Input v-model="form.tagName" @on-change="filterData" style="width:300px" placeholder="请输入模板标签">
            <template #suffix>
              <Icon type="ios-search" />
            </template>
          </Input>
        </FormItem>
        <FormItem label="展示未发布模板" :label-width="120">
          <i-Switch v-model="form.unPublishShow" @on-change="filterData" />
        </FormItem>
      </Form>
    </div>
    <div class="wrapper">
      <div v-if="cardList.length" class="template">
        <Card v-for="(i, index) in cardList" :key="index" style="width:100%;margin-bottom:20px;">
          <div class="w-header" slot="title">
            <Icon size="28" type="ios-people" />
            <div class="title">
              {{ i.manageRole }}
              <span class="underline"></span>
            </div>
            <Icon size="28" type="md-arrow-dropdown" style="cursor:pointer;" @click="handleExpand(i)" />
          </div>
          <div v-show="i.expand">
            <div v-for="j in i.groups" :key="j.groupId" class="content">
              <div class="sub-header">
                <Icon size="24" type="ios-folder" />
                <span class="title">{{ j.groupName }}</span>
              </div>
              <Table
                @on-row-click="handleChooseTemplate"
                size="small"
                :columns="tableColumns"
                :data="j.templates"
                style="margin:10px 0 20px 0"
              >
              </Table>
            </div>
          </div>
        </Card>
      </div>
      <div v-else class="template-no-data">
        暂无数据
      </div>
      <div class="list">
        <Card style="width:300px;min-height:360px;">
          <div class="w-header">
            我的收藏 <span>{{ collectList.length }}</span>
          </div>
          <div v-for="i in collectList" :key="i.id" class="item">
            {{ i.name }}
            <Tooltip content="取消收藏" placement="top-start">
              <Icon
                style="cursor:pointer"
                size="18"
                type="ios-star"
                color="#ebac42"
                @click="handleStar({ ...i, collectFlag: 1 })"
              />
            </Tooltip>
          </div>
        </Card>
      </div>
    </div>
  </div>
</template>

<script>
import { getTemplateTree, collectTemplate, uncollectTemplate, collectTemplateList } from '@/api/server'
export default {
  data () {
    return {
      form: {
        templateName: '',
        tagName: '',
        unPublishShow: true
      },
      // 收藏列表
      collectList: [],
      // 模板数据
      cardList: [],
      originCardList: [],
      tableColumns: [
        {
          title: '模板名称',
          key: 'name',
          render: (h, params) => {
            return (
              <div>
                <span style="margin-right:5px">{params.row.name}</span>
                {params.row.collectFlag === 0 && (
                  <Tooltip content="收藏" placement="top-start">
                    <Icon
                      style="cursor:pointer"
                      size="18"
                      type="ios-star-outline"
                      onClick={e => {
                        e.stopPropagation()
                        this.handleStar(params.row)
                      }}
                    />
                  </Tooltip>
                )}
                {params.row.collectFlag === 1 && (
                  <Tooltip content="取消收藏" placement="top-start">
                    <Icon
                      style="cursor:pointer"
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
              </div>
            )
          }
        },
        {
          title: '状态',
          key: 'status',
          render: (h, params) => {
            return (
              <Tag color={params.row.status === 'created' ? '#85888e' : 'success'}>
                {{ created: '未发布', confirm: '已发布' }[params.row.status]}
              </Tag>
            )
          }
        },
        {
          title: '标签',
          key: 'tags',
          render: (h, params) => {
            return params.row.tags && <Tag>{params.row.tags}</Tag>
          }
        },
        {
          title: '属主',
          key: 'updatedBy'
        }
      ]
    }
  },
  mounted () {
    this.getTemplateData()
    this.getCollectTemplate()
  },
  methods: {
    // 获取模板数据
    async getTemplateData () {
      const { statusCode, data } = await getTemplateTree()
      if (statusCode === 'OK') {
        this.cardList =
          Array.isArray(data) &&
          data.map(i => {
            i.expand = true
            return i
          })
        this.originCardList = this.cardList
      }
    },
    // 展开收缩卡片
    handleExpand (item) {
      item.expand = !item.expand
    },
    // 选中一条模板数据
    handleChooseTemplate (row) {
      this.$router.push({
        path: `/taskman/workbench/createPublish?id=${row.id}`
      })
    },
    // 收藏or取消收藏模板
    async handleStar ({ id, collectFlag }) {
      const method = collectFlag ? uncollectTemplate : collectTemplate
      const { statusCode } = await method(id)
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('successful'),
          desc: this.$t('successful')
        })
        this.getTemplateData()
        this.getCollectTemplate()
      }
    },
    // 收藏模板列表
    async getCollectTemplate () {
      const params = {
        startIndex: 1,
        pageSize: 500
      }
      const { statusCode, data } = await collectTemplateList(params)
      if (statusCode === 'OK') {
        this.collectList = data.contents || []
      }
    },
    // 搜索过滤模板数据
    filterData () {
      const { templateName, tagName, unPublishShow } = this.form
      this.cardList = JSON.parse(JSON.stringify(this.originCardList))
      this.cardList = this.cardList.filter(i => {
        i.groups =
          (Array.isArray(i.groups) &&
            i.groups.filter(j => {
              j.templates =
                (Array.isArray(j.templates) &&
                  j.templates.filter(k => {
                    // 根据模板名、标签名、模版发布状态组合搜索
                    const nameFlag = k.name.toLowerCase().indexOf(templateName.toLowerCase()) > -1
                    const tagFlag = k.tags.toLowerCase().indexOf(tagName.toLowerCase()) > -1
                    const statusFlag = unPublishShow ? true : k.status === 'confirm'
                    return nameFlag && tagFlag && statusFlag
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
  .ivu-form-item {
    margin-bottom: 10px !important;
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
      width: 300px;
      .w-header {
        font-size: 16px;
        font-weight: 700;
        span {
          font-size: 15px;
          font-weight: 400;
        }
      }
      .item {
        padding: 5px 0 3px 5px;
        font-size: 14px;
        color: #515a6e;
        .collect {
          cursor: pointer;
          margin-left: 5px;
        }
      }
    }
    .template {
      padding: 0px 20px 0 0;
      width: calc(100% - 320px);
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
