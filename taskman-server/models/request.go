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
	Id              string `json:"id" json:"id"`
	Name            string `json:"name" json:"name"`
	Form            string `json:"form" json:"form"`
	RequestTemplate string `json:"requestTemplate" json:"request_template"`
	ProcInstanceId  string `json:"procInstanceId" json:"proc_instance_id"`
	Reporter        string `json:"reporter" json:"reporter"`
	ReportTime      string `json:"reportTime" json:"report_time"`
	Emergency       string `json:"emergency" json:"emergency"`
	ReportRole      string `json:"reportRole" json:"report_role"`
	AttachFile      string `json:"attachFile" json:"attach_file"`
	Status          string `json:"status" json:"status"`
	Cache           string `json:"cache" json:"cache"`
	Result          string `json:"result" json:"result"`
	DelFlag         int    `json:"delFlag" json:"del_flag"`
}

type AttachFileTable struct {
	Id           string `json:"id" json:"id"`
	Name         string `json:"name" json:"name"`
	S3Url        string `json:"s3Url" json:"s3_url"`
	S3BucketName string `json:"s3BucketName" json:"s3_bucket_name"`
	S3KeyName    string `json:"s3KeyName" json:"s3_key_name"`
	DelFlag      int    `json:"delFlag" json:"del_flag"`
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
