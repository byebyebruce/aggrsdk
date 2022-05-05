package ownthink

import (
	"fmt"

	"github.com/byebyebruce/aggrsdk/pkg/util"
)

const URL = `https://api.ownthink.com/bot?appid=%s&spoken=%s`

type OwnThink struct {
	AppID string
}

func New(appID string) *OwnThink {
	return &OwnThink{
		AppID: appID,
	}
}

type resp struct {
	Message string `json:"message"`
	Data    struct {
		Type int `json:"type"`
		Info struct {
			Text string `json:"text"`
		} `json:"info"`
	} `json:"data"`
}

// Bot 聊天机器人
func (o *OwnThink) Bot(content string) (string, error) {
	url := fmt.Sprintf(URL, o.AppID, content)
	resp := &resp{}
	if err := util.HttpGet(url, resp); err != nil {
		return "", err
	}
	if resp.Message != "success" {
		return "", fmt.Errorf("%s", resp.Message)
	}
	return resp.Data.Info.Text, nil
}
