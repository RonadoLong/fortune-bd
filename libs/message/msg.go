package message

import (
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

const (
	regionId        = "cn-hangzhou"
	AccessKeyId     = "LTAI4GBFG46DfFaw1TAg9mvx"
	AccessKeySecret = "RNSIHK8Edai5BhLS4dw8QgrE99DCAy"
	TemplateCode    = "SMS_190273927"
	SignName        = "iFortune"
)

// SendMsg 发送短信
func SendMsg(phone, code string) error {
	if phone == "" || code == "" {
		return errors.New("phone or code not allow bank")
	}
	client, err := dysmsapi.NewClientWithAccessKey(
		regionId,
		AccessKeyId,
		AccessKeySecret,
	)
	if err != nil {
		return err
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phone
	request.TemplateCode = TemplateCode
	request.TemplateParam = fmt.Sprintf(`{"code":%v}`, code)
	request.SignName = SignName
	response, err := client.SendSms(request)
	if err != nil {
		return err
	}
	if response.Code != "OK" {
		return errors.New(response.Message)
	}
	fmt.Printf("response is %#v\n", response)
	return nil
}
