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
	err := a.Text2AudioFile("tmp.mp3", "haha")
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
	tmp := "http://81.70.119.61:8088/fs/a.mp4_audio.wav"
	text, err := a.QueryAASR(tmp, "wav", 1737)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(text)
}
