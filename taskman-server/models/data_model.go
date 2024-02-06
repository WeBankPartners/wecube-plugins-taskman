package models

import "time"

type DataModel struct {
	PluginPackageDataModel
	Entities []*DataModelEntity `json:"entities"`
}

type DataModelEntity struct {
	Id               string `json:"id"`               // 唯一标识
	DataModelId      string `json:"dataModelId"`      // 所属数据模型
	DataModelVersion int    `json:"dataModelVersion"` // 版本
	PackageName      string `json:"packageName"`      // 包名
	Name             string `json:"name"`             // 模型名
	DisplayName      string `json:"displayName"`      // 显示名
	Description      string `json:"description"`      // 描述
}

type EntityQueryResult struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    []*EntityDataObj `json:"data"`
}

type EntityDataObj struct {
	Id          string `json:"guid"`
	DisplayName string `json:"key_name"`
	IsNew       bool   `json:"isNew"`
	PackageName string `json:"package_name"`
	Entity      string `json:"entity"`
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

type PluginPackageDataModel struct {
	Id           string    `json:"id"`           // 唯一标识
	Version      int       `json:"version"`      // 版本
	PackageName  string    `json:"packageName"`  // 包名
	IsDynamic    bool      `json:"dynamic"`      // 是否动态
	UpdatePath   string    `json:"updatePath"`   // 请求路径
	UpdateMethod string    `json:"updateMethod"` // 请求方法
	UpdateSource string    `json:"updateSource"` // 来源
	UpdatedTime  time.Time `json:"updatedTime"`  // 更新时间
	UpdateTime   int64     `json:"updateTime"`   // 旧更新时间,毫秒时间戳
}

type ProcEntityAttributeObj struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	DataType          string `json:"dataType"`
	Mandatory         bool   `json:"mandatory"`
	RefPackageName    string `json:"refPackageName"`
	RefEntityName     string `json:"refEntityName"`
	RefAttrName       string `json:"refAttrName"`
	ReferenceId       string `json:"referenceId"`
	Active            bool   `json:"active"`
	EntityId          string `json:"entityId"`
	EntityName        string `json:"entityName"`
	EntityDisplayName string `json:"entityDisplayName"`
	EntityPackage     string `json:"entityPackage"`
	Multiple          string `json:"multiple"`
}

type ProcEntity struct {
	Id          string                    `json:"id"`
	PackageName string                    `json:"packageName"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	DisplayName string                    `json:"displayName"`
	Attributes  []*ProcEntityAttributeObj `json:"attributes"`
}
