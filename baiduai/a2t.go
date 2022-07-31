package baiduai

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/byebyebruce/aggrsdk/pkg/util"
)

const (
	VOP_URL = "http://vop.baidu.com/server_api"
)

type VOPResponse struct {
	CorpusNo string   `json:"corpus_no"`
	ErrMsg   string   `json:"err_msg"`
	ErrNo    int64    `json:"err_no"`
	Result   []string `json:"result"`
	Sn       string   `json:"sn"`
}

func (this *BaiduAI) Reader2Text(reader io.Reader) (string, error) {
	this.genCredentials()
	// POST http://vop.baidu.com/server_api?dev_pid=1537&cuid=******&token=1.a6b7dbd428f731035f771b8d********.86400.1292922000-2346678-124328
	url := fmt.Sprintf("%s?cuid=%v&token=%s", VOP_URL, this.Cuid, this.Credentials.Access_token)

	resp := VOPResponse{}
	err := util.HTTPPostRaw(url, map[string]string{
		"Content-Type": "audio/wav;rate=16000",
	}, reader, &resp)
	if err != nil {
		return "", err
	}
	if resp.ErrNo != 0 {
		return "", fmt.Errorf("%v:%s", resp.ErrNo, resp.ErrMsg)
	}
	if len(resp.Result) == 0 {
		return "", fmt.Errorf("no result")
	}
	return resp.Result[0], nil
}

func (this *BaiduAI) Audio2Text(wavFile string) (string, error) {
	r, err := ioutil.ReadFile(wavFile)
	if err != nil {
		return "", err
	}
	return this.Reader2Text(bytes.NewReader(r))
}
