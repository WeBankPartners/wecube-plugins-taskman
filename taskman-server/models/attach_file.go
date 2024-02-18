package models

type AttachFileTable struct {
	Id           string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Name         string `json:"name" xorm:"name"`
	S3BucketName string `json:"s3BucketName" xorm:"s3_bucket_name"`
	S3KeyName    string `json:"s3KeyName" xorm:"s3_key_name"`
	DelFlag      int    `json:"delFlag" xorm:"del_flag"`
	Request      string `json:"request" xorm:"request"`
	Task         string `json:"task" xorm:"task"`
}

func (AttachFileTable) TableName() string {
	return "attach_file"
}
