package baidu

import (
	"os"
	"testing"
)

var (
	apiKey    = os.Getenv("baidu_key")
	apiSecret = os.Getenv("baidu_secret")
)

func TestText2AudioFile(t *testing.T) {
	a := NewAPI_Util(apiKey, apiSecret)
	err := a.Text2AudioFile("tmp.wav", "请问今天北京的天气怎么样", WAV)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAudio2Text(t *testing.T) {
	a := NewAPI_Util(apiKey, apiSecret)
	text, err := a.Audio2Text("16k.wav")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(text)
}

func TestAASR(t *testing.T) {
	a := NewAPI_Util(apiKey, apiSecret)
	tmp := "http://127.0.0.1:8088/fs/a.mp4_audio.wav"
	text, err := a.QueryAASR(tmp, "wav", 1737)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(text)
}
