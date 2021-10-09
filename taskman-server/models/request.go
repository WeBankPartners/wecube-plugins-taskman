package models

type EntityQueryResult struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    []*EntityDataObj `json:"data"`
}

type EntityDataObj struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
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
	Id            string                 `json:"id"`
	EntityData    map[string]interface{} `json:"entityData"`
	PreviousIds   []string               `json:"previousIds"`
	SucceedingIds []string               `json:"succeedingIds"`
}

type RequestTable struct {
	Id              string `json:"id" xorm:"id"`
	Name            string `json:"name" xorm:"name"`
	Form            string `json:"form" xorm:"form"`
	RequestTemplate string `json:"requestTemplate" xorm:"request_template"`
	ProcInstanceId  string `json:"procInstanceId" xorm:"proc_instance_id"`
	Reporter        string `json:"reporter" xorm:"reporter"`
	ReportTime      string `json:"reportTime" xorm:"report_time"`
	Emergency       string `json:"emergency" xorm:"emergency"`
	ReportRole      string `json:"reportRole" xorm:"report_role"`
	AttachFile      string `json:"attachFile" xorm:"attach_file"`
	Status          string `json:"status" xorm:"status"`
	Cache           string `json:"cache" xorm:"cache"`
	Result          string `json:"result" xorm:"result"`
	CreatedBy       string `json:"createdBy" xorm:"created_by"`
	CreatedTime     string `json:"createdTime" xorm:"created_time"`
	UpdatedBy       string `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime     string `json:"updatedTime" xorm:"updated_time"`
	DelFlag         int    `json:"delFlag" xorm:"del_flag"`
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
	FullEntityDataId string                         `json:"fullEntityDataId"`
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
	Entity string                   `json:"entity"`
	Title  []*FormItemTemplateTable `json:"title"`
	Value  []map[string]string      `json:"value"`
}
