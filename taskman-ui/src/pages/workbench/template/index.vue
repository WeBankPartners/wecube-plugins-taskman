<template>
  <div class="new-publish-template">
    <div class="search">
      <Form :model="form" inline label-position="left">
        <!--模板名-->
        <FormItem :label="$t('tw_template_name')" :label-width="lang === 'zh-CN' ? 60 : 115">
          <Input
            v-model="form.templateName"
            @on-change="filterData"
            style="width:300px"
            clearable
            :placeholder="$t('tw_template_placeholder')"
          >
            <template #suffix>
              <Icon type="ios-search" />
            </template>
          </Input>
        </FormItem>
        <!--操作对象类型-->
        <FormItem :label="$t('tw_operator_type')" :label-width="lang === 'zh-CN' ? 100 : 155">
          <Select
            v-model="form.operatorObjType"
            @on-change="filterData"
            clearable
            filterable
            :placeholder="$t('tw_operator_placeholder')"
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
          <!--已发布-->
          <TabPane :label="$t('tw_template_publish_tab')" name="confirm"></TabPane>
          <!--我的草稿-->
          <!-- <TabPane v-if="draftCardList.length" :label="$t('tw_template_draft_tab')" name="created"></TabPane> -->
        </Tabs>
        <Card :bordered="false" dis-hover :padding="0" style="height:400px;">
          <template v-if="cardList.length">
            <Card v-for="(i, index) in cardList" :key="index" style="width:100%;margin-bottom:20px;">
              <div class="w-header" slot="title">
                <Icon size="28" type="ios-people" />
                <div class="title">
                  {{ i.manageRoleDisplay }}
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
                        handleChooseTemplate(row, row.manageRole)
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
          <div v-if="!spinShow && !cardList.length" class="template-no-data">
            {{ $t('tw_no_data') }}
          </div>
          <Spin fix v-if="spinShow"></Spin>
        </Card>
      </div>
      <!--收藏列表-->
      <div class="list">
        <Card style="width:520px;min-height:600px;">
          <div class="w-header">
            {{ $t('tw_my_collect') }} <span>{{ collectList.length || 0 }}</span>
          </div>
          <Table
            style="margin:20px 0;"
            :show-header="true"
            :disabled-hover="true"
            @on-row-click="
              row => {
                handleChooseTemplate(row, row.useRole)
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
import { debounce, deepClone } from '@/pages/util'
export default {
  data () {
    return {
      activeName: 'confirm', // confirm已发布，created我的草稿(未发布)
      type: '', // 1发布,2请求,3问题,4事件,5变更
      form: {
        templateName: '',
        operatorObjType: '' // 操作对象类型
      },
      operateOptions: [],
      collectList: [], // 收藏列表
      cardList: [], // 模板数据
      spinShow: false, // 加载动画
      originCardList: [],
      draftCardList: [], // 草稿页签数据
      tableColumns: [
        {
          title: this.$t('tw_template_name'),
          key: 'name',
          minWidth: 250,
          render: (h, params) => {
            return (
              <div>
                {params.row.collectFlag === 0 && this.activeName === 'confirm' && (
                  <Tooltip content={this.$t('tw_collec_tooltip')} placement="top-start">
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
                  <Tooltip content={this.$t('tw_uncollec_tooltip')} placement="top-start">
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
          title: this.$t('tw_operator_type'),
          key: 'operatorObjType',
          render: (h, params) => {
            if (params.row.operatorObjType) {
              return <Tag>{params.row.operatorObjType}</Tag>
            } else {
              return <span>-</span>
            }
          }
        },
        {
          title: this.$t('tw_owner_role'),
          key: 'updatedBy',
          render: (h, params) => {
            return (
              <div style="display:flex;flex-direction:column">
                <span>{params.row.roleDisplay}</span>
                <span>{params.row.handler}</span>
              </div>
            )
          }
        }
      ],
      collectColumns: [
        {
          title: this.$t('tw_template_name'),
          key: 'name',
          width: 230,
          render: (h, params) => {
            return (
              <div>
                <Tooltip content={this.$t('tw_uncollec_tooltip')} placement="top-start">
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
          title: this.$t('tw_template_status'),
          width: 120,
          key: 'status',
          render: (h, params) => {
            const list = [
              { label: this.$t('tw_template_status_use'), value: 1, color: '#19be6b' },
              { label: this.$t('tw_template_status_disable'), value: 2, color: '#c5c8ce' },
              { label: this.$t('tw_template_status_role'), value: 3, color: '#ed4014' }
            ]
            const item = list.find(i => i.value === params.row.status)
            return item && <Tag color={item.color}>{item.label}</Tag>
          }
        },
        {
          title: this.$t('useRoles'),
          key: 'useRoleDisplay'
        }
      ],
      lang: window.localStorage.getItem('lang'),
      createRouteMap: {
        '1': 'createPublish',
        '2': 'createRequest',
        '3': 'createProblem',
        '4': 'createEvent',
        '5': 'createChange'
      }
    }
  },
  watch: {
    // 解决同路由跳转不触发页面更新问题
    $route (to, from) {
      this.collectList = []
      this.cardList = []
      this.type = this.$route.query.type || ''
      this.activeName = 'confirm'
      this.form.templateName = ''
      this.form.operatorObjType = ''
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
      this.spinShow = true
      const { statusCode, data } = await getTemplateTree()
      this.spinShow = false
      if (statusCode === 'OK') {
        let templateList = data || []
        templateList = templateList.map(i => {
          i.expand = true
          i.groups.forEach(j => {
            j.templates.forEach(template => {
              template.manageRole = i.manageRole
              template.manageRoleDisplay = i.manageRoleDisplay
              if (template.operatorObjType && template.type === Number(this.type)) {
                this.operateOptions.push(template.operatorObjType)
              }
            })
          })
          return i
        })
        // 数据去重
        this.operateOptions = Array.from(new Set(this.operateOptions))
        this.originCardList = deepClone(templateList)
        this.filterData()
        // 获取草稿态数据，没有则隐藏草稿标签
        let cardList = deepClone(this.originCardList)
        this.draftCardList = cardList.filter(i => {
          i.groups =
            (Array.isArray(i.groups) &&
              i.groups.filter(j => {
                j.templates =
                  (Array.isArray(j.templates) &&
                    j.templates.filter(template => {
                      return template.status === 'created' && template.type === Number(this.type)
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
    },
    // 展开收缩卡片
    handleExpand (item) {
      item.expand = !item.expand
    },
    // 发起操作
    handleChooseTemplate (row, role) {
      // 模板禁用提示
      if (row.status === 2) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('tw_template_disable_tips')
        })
        // 模板权限移除提示
      } else if (row.status === 3) {
        return this.$Notice.warning({
          title: this.$t('warning'),
          desc: this.$t('tw_template_role_tips')
        })
      }
      // const path = this.type === '2' ? 'createRequest' : 'createPublish'
      const path = this.createRouteMap[this.type]
      const url = `/taskman/workbench/${path}`
      this.$router.push({
        path: url,
        query: {
          requestTemplate: row.id,
          role: role, // 模板创建人角色
          jumpFrom: ''
        }
      })
    },
    // 收藏or取消收藏模板
    handleStar: debounce(async function ({ id, collectFlag, manageRole }) {
      const method = collectFlag ? uncollectTemplate : collectTemplate
      const params = collectFlag ? id : { templateId: id, role: manageRole }
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
        action: Number(this.type),
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
      this.cardList = deepClone(this.originCardList)
      this.cardList = this.cardList.filter(i => {
        i.groups =
          (Array.isArray(i.groups) &&
            i.groups.filter(j => {
              j.templates =
                (Array.isArray(j.templates) &&
                  j.templates.filter(k => {
                    // 根据模板名、标签名、模版发布状态组合搜索
                    const nameFilter = k.name.toLowerCase().indexOf(templateName.toLowerCase()) > -1
                    const operatorFilter = operatorObjType ? k.operatorObjType === operatorObjType : true
                    const typeFilter = Number(this.type) === k.type
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
