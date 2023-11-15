package models

type QueryRequestFilterObj struct {
	Name     string      `json:"name"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type QueryRequestSorting struct {
	Asc   bool   `json:"asc"`
	Field string `json:"field"`
}

type QueryRequestDialect struct {
	AssociatedData map[string]string `json:"associatedData"`
	QueryMode      string            `json:"queryMode"`
}

type QueryRequestParam struct {
	Filters       []*QueryRequestFilterObj `json:"filters"`
	Dialect       *QueryRequestDialect     `json:"dialect"`
	Paging        bool                     `json:"paging"`
	Pageable      *PageInfo                `json:"pageable"`
	Sorting       *QueryRequestSorting     `json:"sorting"`
	ResultColumns []string                 `json:"resultColumns"`
}

type TransFiltersParam struct {
	IsStruct   bool
	StructObj  interface{}
	Prefix     string
	KeyMap     map[string]string
	PrimaryKey string
}

type PlatformRequestParam struct {
	Tab         string   `json:"tab"`         // 标签,取值有:
	Action      int      `json:"action"`      // 行为, 1表示发布行为,2请求,3问题,4事件, 5变更
	Type        int      `json:"type"`        // 0代表所有,1表示请求定版,2表示任务处理
	Query       string   `json:"query"`       // ID和名称模糊匹配
	TemplateId  string   `json:"templateId"`  // 模版id
	Status      []string `json:"status"`      // 请求状态 Pending InProgress(Faulted)
	OperatorObj string   `json:"operatorObj"` // 操作对象
	CreatedBy   string   `json:"createdBy"`   // 创建人
	StartIndex  int      `json:"startIndex"`
	PageSize    int      `json:"pageSize"`
}
