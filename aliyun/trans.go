package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alimt"
)

type Translator struct {
	Aliyun
}

func NewTranslator(aliyun Aliyun) Translator {
	return Translator{
		Aliyun: aliyun,
	}
}
func (t Translator) Do(text string, srcLanguage string, targetLanguage string /*en zh jp*/) (string, error) {
	client, err := alimt.NewClientWithAccessKey("cn-qingdao", t.AK, t.AS)
	request := alimt.CreateTranslateGeneralRequest()
	request.Scheme = "https"

	request.SourceLanguage = srcLanguage
	request.TargetLanguage = targetLanguage
	request.SourceText = text
	request.FormatType = "text"

	response, err := client.TranslateGeneral(request)
	if err != nil {
		return "", err
	}
	return response.Data.Translated, nil
}
