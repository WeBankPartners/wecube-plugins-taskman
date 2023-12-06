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
	Tab      string `json:"tab"`      // 标签,取值有:pending 待处理    hasProcessed 已处理  submit 我提交的 draft 我的暂存  collect 收藏
	Action   int    `json:"action"`   // 行为, 1表示发布行为,2请求,3问题,4事件, 5变更
	Type     int    `json:"type"`     // 0代表所有,1表示请求定版,2表示任务处理
	Rollback int    `json:"rollback"` // 0代表所有,1表示被退回,2表示其他(submit 我提交的tab下才有这个筛选生效)
	CommonRequestParam
}

type FilterRequestParam struct {
	StartTime string `json:"startTime"` //开始时间
}

// QueryCollectTemplateObj 模板查询条件
type QueryCollectTemplateParam struct {
	Id               string               `json:"id"`               // ID
	Name             string               `json:"name"`             // Name
	TemplateGroupId  []string             `json:"templateGroupId"`  // 模板组id
	OperatorObjType  []string             `json:"operatorObjType"`  // 操作对象类型
	ProcDefName      []string             `json:"procDefName"`      // 使用编排
	ManageRole       []string             `json:"manageRole"`       // 属主角色
	Owner            []string             `json:"owner"`            // 属主
	UseRole          []string             `json:"useRole"`          // 使用角色
	Tags             []string             `json:"tags" `            // 标签
	CreatedStartTime string               `json:"createdStartTime"` // 创建开始时间
	CreatedEndTime   string               `json:"createdEndTime"`   // 创建结束时间
	StartIndex       int                  `json:"startIndex"`
	PageSize         int                  `json:"pageSize"`
	Sorting          *QueryRequestSorting `json:"sorting"` // 排序字段
}

type RequestHistoryParam struct {
	Tab        string `json:"tab"`        // 标签,取值有: draft 暂存  commit 已经提交
	Permission string `json:"permission"` // 权限,取值有: group 本组,  all 表示所有
	Action     int    `json:"action"`     // 行为: 0表示所有,1表示发布行为,2请求
	CommonRequestParam
}

type CommonRequestParam struct {
	Id               string               `json:"id"`                // ID
	Name             string               `json:"name"`              // Name
	TemplateId       []string             `json:"templateId"`        // 模版id
	Status           []string             `json:"status"`            // 请求状态 Pending InProgress(Faulted)
	OperatorObj      string               `json:"operatorObj"`       // 操作对象
	CreatedBy        []string             `json:"createdBy"`         // 创建人
	OperatorObjType  []string             `json:"operatorObjType"`   // 操作对象类型
	ProcDefName      []string             `json:"procDefName"`       // 使用编排
	Handler          []string             `json:"handler"`           // 当前处理人
	CreatedStartTime string               `json:"createdStartTime" ` // 创建开始时间
	CreatedEndTime   string               `json:"createdEndTime" `   // 创建结束时间
	ExpectStartTime  string               `json:"expectStartTime" `  // 期望开始时间
	ExpectEndTime    string               `json:"expectEndTime" `    // 期望结束时间
	StartIndex       int                  `json:"startIndex"`
	PageSize         int                  `json:"pageSize"`
	Sorting          *QueryRequestSorting `json:"sorting"` // 排序字段
}

type PlatDataParam struct {
	Param      CommonRequestParam
	QueryParam []interface{}
	User       string
	Where      string
	Sql        string
	UserToken  string
}
