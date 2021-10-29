package models

type EntityQueryResult struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    []*EntityDataObj `json:"data"`
}

type EntityDataObj struct {
	Id          string `json:"guid"`
	DisplayName string `json:"key_name"`
	IsNew       bool   `json:"isNew"`
}

type EntityTreeResult struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    EntityTreeData `json:"data"`
}

type EntityTreeData struct {
	EntityTreeNodes  []*EntityTreeObj `json:"entityTreeNodes"`
	ProcessSessionId string           `json:"processSessionId"`
}

type EntityTreeObj struct {
	PackageName   string                 `json:"packageName"`
	EntityName    string                 `json:"entityName"`
	DataId        string                 `json:"dataId"`
	DisplayName   string                 `json:"displayName"`
	FullDataId    interface{}            `json:"fullDataId"`
	Id            string                 `json:"id"`
	EntityData    map[string]interface{} `json:"entityData"`
	PreviousIds   []string               `json:"previousIds"`
	SucceedingIds []string               `json:"succeedingIds"`
	EntityDataOp  string                 `json:"entityDataOp"`
}

type RequestTable struct {
	Id                  string `json:"id" xorm:"id"`
	Name                string `json:"name" xorm:"name"`
	Form                string `json:"form" xorm:"form"`
	RequestTemplate     string `json:"requestTemplate" xorm:"request_template"`
	RequestTemplateName string `json:"requestTemplateName" xorm:"-"`
	ProcInstanceId      string `json:"procInstanceId" xorm:"proc_instance_id"`
	ProcInstanceKey     string `json:"procInstanceKey" xorm:"proc_instance_key"`
	Reporter            string `json:"reporter" xorm:"reporter"`
	Handler             string `json:"handler" xorm:"handler"`
	ReportTime          string `json:"reportTime" xorm:"report_time"`
	Emergency           int    `json:"emergency" xorm:"emergency"`
	ReportRole          string `json:"reportRole" xorm:"report_role"`
	AttachFile          string `json:"attachFile" xorm:"attach_file"`
	Status              string `json:"status" xorm:"status"`
	Cache               string `json:"cache" xorm:"cache"`
	BindCache           string `json:"bindCache" xorm:"bind_cache"`
	Result              string `json:"result" xorm:"result"`
	ExpireTime          string `json:"expireTime" xorm:"expire_time"`
	ExpectTime          string `json:"expectTime" xorm:"expect_time"`
	CreatedBy           string `json:"createdBy" xorm:"created_by"`
	CreatedTime         string `json:"createdTime" xorm:"created_time"`
	UpdatedBy           string `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime         string `json:"updatedTime" xorm:"updated_time"`
	DelFlag             int    `json:"delFlag" xorm:"del_flag"`
}

type AttachFileTable struct {
	Id           string `json:"id" xorm:"id"`
	Name         string `json:"name" xorm:"name"`
	S3Url        string `json:"s3Url" xorm:"s3_url"`
	S3BucketName string `json:"s3BucketName" xorm:"s3_bucket_name"`
	S3KeyName    string `json:"s3KeyName" xorm:"s3_key_name"`
	DelFlag      int    `json:"delFlag" xorm:"del_flag"`
}

type RequestCacheData struct {
	ProcDefId         string                         `json:"procDefId"`
	ProcDefKey        string                         `json:"procDefKey"`
	RootEntityValue   RequestCacheEntityValue        `json:"rootEntityValue"`
	TaskNodeBindInfos []*RequestCacheTaskNodeBindObj `json:"taskNodeBindInfos"`
}

type RequestCacheTaskNodeBindObj struct {
	BoundEntityValues []*RequestCacheEntityValue `json:"boundEntityValues"`
	NodeDefId         string                     `json:"nodeDefId"`
	NodeId            string                     `json:"nodeId"`
}

type RequestCacheEntityValue struct {
	AttrValues       []*RequestCacheEntityAttrValue `json:"attrValues"`
	BindFlag         string                         `json:"bindFlag"`
	EntityDataId     string                         `json:"entityDataId"`
	EntityDataOp     string                         `json:"entityDataOp"`
	EntityDataState  string                         `json:"entityDataState"`
	EntityDefId      string                         `json:"entityDefId"`
	EntityName       string                         `json:"entityName"`
	FullEntityDataId interface{}                    `json:"fullEntityDataId"`
	Oid              string                         `json:"oid"`
	PackageName      string                         `json:"packageName"`
	PreviousOids     []string                       `json:"previousOids"`
	Processed        bool                           `json:"processed"`
	SucceedingOids   []string                       `json:"succeedingOids"`
}

type RequestCacheEntityAttrValue struct {
	AttrDefId string      `json:"attrDefId"`
	AttrName  string      `json:"attrName"`
	DataType  string      `json:"dataType"`
	DataValue interface{} `json:"dataValue"`
}

type RequestPreDataTableObj struct {
	PackageName   string                   `json:"packageName"`
	Entity        string                   `json:"entity"`
	ItemGroup     string                   `json:"itemGroup"`
	ItemGroupName string                   `json:"itemGroupName"`
	RefEntity     []string                 `json:"-"`
	SortLevel     int                      `json:"-"`
	Title         []*FormItemTemplateTable `json:"title"`
	Value         []*EntityTreeObj         `json:"value"`
}

type StartInstanceResult struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Data    StartInstanceResultData `json:"data"`
}

type StartInstanceResultData struct {
	Id          int    `json:"id"`
	ProcInstKey string `json:"procInstKey"`
	ProcDefId   string `json:"procDefId"`
	ProcDefKey  string `json:"procDefKey"`
	Status      string `json:"status"`
}

type RequestPreDataSort []*RequestPreDataTableObj

func (s RequestPreDataSort) Len() int {
	return len(s)
}

func (s RequestPreDataSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s RequestPreDataSort) Less(i, j int) bool {
	if s[i].SortLevel == s[j].SortLevel {
		return s[i].Entity < s[j].Entity
	}
	return s[i].SortLevel > s[j].SortLevel
}

type RequestPreDataDto struct {
	RootEntityId string                    `json:"rootEntityId"`
	Data         []*RequestPreDataTableObj `json:"data"`
}

type TerminateInstanceParam struct {
	ProcInstId  string `json:"procInstId"`
	ProcInstKey string `json:"procInstKey"`
}

type EntityNodeBindQueryObj struct {
	NodeDefId string `xorm:"node_def_id"`
	ItemGroup string `xorm:"item_group"`
}

type InstanceStatusQuery struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    InstanceStatusQueryObj `json:"data"`
}

type InstanceStatusQueryObj struct {
	Id                int                        `json:"id"`
	ProcDefId         string                     `json:"procDefId"`
	ProcInstKey       string                     `json:"procInstKey"`
	ProcInstName      string                     `json:"procInstName"`
	Status            string                     `json:"status"`
	TaskNodeInstances []*InstanceStatusQueryNode `json:"taskNodeInstances"`
}

type InstanceStatusQueryNode struct {
	Id        int    `json:"id"`
	NodeId    string `json:"nodeId"`
	NodeDefId string `json:"nodeDefId"`
	NodeName  string `json:"nodeName"`
	NodeType  string `json:"nodeType"`
	Status    string `json:"status"`
}
