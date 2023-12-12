package models

// CollectTemplateTable  收藏模板
type CollectTemplateTable struct {
	Id              string `json:"id" xorm:"id"`
	RequestTemplate string `json:"requestTemplate" xorm:"request_template"` // 收藏模板ID
	User            string `json:"user" xorm:"user"`                        // 收藏用户
	Role            string `json:"role" xorm:"role"`                        // 收藏模板时角色
	Type            int    `json:"type" xorm:"type"`                        // 类型:0表示请求 1发布
	CreatedTime     string `json:"createdTime" xorm:"created_time"`
}
