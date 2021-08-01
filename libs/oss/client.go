package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
)

const (
	regionId        = "cn-hangzhou"
	AccessKeyId     = "LTAI4GBFG46DfFaw1TAg9mvx"
	AccessKeySecret = "RNSIHK8Edai5BhLS4dw8QgrE99DCAy"
	Endpoint        = "http://oss-cn-heyuan.aliyuncs.com"
	EndPointNoHttp  = "oss-cn-heyuan.aliyuncs.com"
	SignName        = "iFortune"
)

type OssClient struct {
	accessKeyId     string
	accessKeySecret string
	endPoint        string
	bucketName      string
	client          *oss.Client
}

func NewOssClient(accessKey, accessSecret, endPoint, bucKeyName string) (*OssClient, error) {
	client, err := oss.New(endPoint, accessKey, accessSecret)
	if err != nil {
		return nil, err
	}
	return &OssClient{
		accessKeyId:     accessKey,
		accessKeySecret: accessSecret,
		endPoint:        endPoint,
		bucketName:      bucKeyName,
		client:          client,
	}, nil
}

func (o *OssClient) GetBucKey() (*oss.Bucket, error) {
	return o.client.Bucket(o.bucketName)
}

func (o *OssClient) PushFile(objectName, localFileName string) error {
	buckey, err := o.GetBucKey()
	if err != nil {
		return err
	}
	if err = buckey.PutObjectFromFile(objectName, localFileName); err != nil {
		return err
	}
	return nil
}

func (o *OssClient) PushFileWithIOReader(objectName string, localFile io.Reader) error {
	buckey, err := o.GetBucKey()
	if err != nil {
		return err
	}
	if err = buckey.PutObject(objectName, localFile); err != nil {
		return err
	}
	return nil
}
