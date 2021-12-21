package file_server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io/ioutil"
)

type MinioServer struct {
	ServerAddress string        `json:"server_address"`
	AccessKey     string        `json:"access_key"`
	SecretKey     string        `json:"secret_key"`
	SSL           bool          `json:"ssl"`
	Client        *minio.Client `json:"client"`
}

type MinioParam struct {
	Ctx         context.Context `json:"ctx"`
	Bucket      string          `json:"bucket"`
	ObjectName  string          `json:"object_name"`
	FileContent []byte          `json:"file_content"`
}

func (m *MinioServer) Init() error {
	minioClient, err := minio.New(m.ServerAddress, &minio.Options{
		Creds:  credentials.NewStaticV4(m.AccessKey, m.SecretKey, ""),
		Secure: m.SSL,
	})
	if err == nil {
		m.Client = minioClient
	}
	return err
}

func (m *MinioServer) Upload(param MinioParam) error {
	err := m.Client.MakeBucket(param.Ctx, param.Bucket, minio.MakeBucketOptions{Region: ""})
	if err != nil {
		exists, errBucketExists := m.Client.BucketExists(param.Ctx, param.Bucket)
		if errBucketExists != nil {
			return fmt.Errorf("Check bucket if exist fail,%s ", errBucketExists.Error())
		}
		if !exists {
			return fmt.Errorf("Bucket:%s is not exist ", param.Bucket)
		}
	}

	info, putErr := m.Client.PutObject(param.Ctx, param.Bucket, param.ObjectName, bytes.NewReader(param.FileContent), int64(len(param.FileContent)), minio.PutObjectOptions{})
	if putErr != nil {
		return fmt.Errorf("Upload minio file fail,%s ", putErr.Error())
	}
	fmt.Printf("info key:%s \n", info.Key)
	return nil
}

func (m *MinioServer) Download(param MinioParam) (result []byte, err error) {
	obj, getErr := m.Client.GetObject(param.Ctx, param.Bucket, param.ObjectName, minio.GetObjectOptions{})
	if getErr != nil {
		err = fmt.Errorf("Download file fail,%s ", getErr.Error())
		return result, err
	}
	result, err = ioutil.ReadAll(obj)
	if err != nil {
		err = fmt.Errorf("Read file content fail,%s ", err.Error())
	}
	return result, err
}

func (m *MinioServer) Remove(param MinioParam) error {
	err := m.Client.RemoveObject(param.Ctx, param.Bucket, param.ObjectName, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		err = fmt.Errorf("Remove file:%s/%s fail,%s ", param.Bucket, param.ObjectName, err.Error())
	}
	return err
}
