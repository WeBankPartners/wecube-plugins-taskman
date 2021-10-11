package models

type TaskMetaResult struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Data    TaskMetaResultData `json:"data"`
}

type TaskMetaResultData struct {
	FormMetaId    string                `json:"formMetaId"`
	FormItemMetas []*TaskMetaResultItem `json:"formItemMetas"`
}

type TaskMetaResultItem struct {
	FormItemMetaId string `json:"formItemMetaId"`
	PackageName    string `json:"packageName"`
	EntityName     string `json:"entityName"`
	AttrName       string `json:"attrName"`
}
