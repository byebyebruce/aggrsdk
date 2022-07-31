package baiduai

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const TSN_URL = "http://tsn.baidu.com/text2audio"

type TSNAudioType string

const (
	MP3    TSNAudioType = "3" // 3为mp3格式(默认)
	PCM16K TSNAudioType = "4" // 4为pcm-16k
	PCM8K  TSNAudioType = "5" // 5为pcm-8k
	WAV    TSNAudioType = "6" // 6为wav
)

// Text2AudioFile 文字合成语音
func (this *BaiduAI) Text2AudioFile(filePath, text string, t TSNAudioType) error {
	body, err := this.Text2AudioBuffer(text, t)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, body, 0666)
	if err != nil {
		return err
	}
	return nil
}

// Text2AudioBuffer 文字合成语音buffer
func (this *BaiduAI) Text2AudioBuffer(text string, t TSNAudioType) ([]byte, error) {
	this.genCredentials()
	param := url.Values{}
	param.Set("tex", text)
	param.Set("tok", this.Credentials.Access_token)
	param.Set("cuid", this.Cuid)
	param.Set("ctp", "1")
	param.Set("lan", "zh")
	param.Set("aue", string(t))

	/*
		aue	选填	3为mp3格式(默认)； 4为pcm-16k；5为pcm-8k；6为wav（内容同pcm-16k）; 注意aue=4或者6是语音识别要求的格式，但是音频内容不是语音识别要求的自然人发音，所以识别效果会受影响。
	*/

	response, err := http.PostForm(TSN_URL, param)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	contentType := response.Header.Get("Content-type")
	switch contentType {
	case "audio/wav",
		"audio/mp3",
		"audio/basic;codec=pcm;rate=16000;channel=1",
		"audio/basic;codec=pcm;rate=8000;channel=1":
		return body, nil
	default:
		var errMsg API_Response
		err = json.Unmarshal(body, &errMsg)
		if nil != err {
			return nil, err
		}
		return nil, fmt.Errorf("%+v", errMsg.Err_msg)
	}
}
