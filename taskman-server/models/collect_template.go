package models

// CollectTemplateTable  收藏模板
type CollectTemplateTable struct {
	Id              string `json:"id" xorm:"id"`
	RequestTemplate string `json:"requestTemplate" xorm:"request_template"` // 收藏模板ID
	User            string `json:"user" xorm:"user"`                        // 收藏用户
	Type            int    `json:"type" xorm:"type"`                        // 类型:0表示请求 1发布
	CreatedTime     string `json:"createdTime" xorm:"created_time"`
}

// QueryCollectTemplateObj 模板查询条件
type QueryCollectTemplateObj struct {
	StartIndex int `json:"startIndex"`
	PageSize   int `json:"pageSize"`
}
