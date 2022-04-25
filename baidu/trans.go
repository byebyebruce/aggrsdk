package baidu

import (
	"fmt"
	"time"
)

func (this *API_Util) Trans(audioURL string,
	format string, /*wav pcm*/
	pid int /*[80001（中文语音近场识别模型极速版）, 1737（英文模型）] */) (DetailedResult, error) {

	const api = "https://aip.baidubce.com/rpc/2.0/mt/texttrans/v1?access_token=%s"

	this.genCredentials()
	task, err := this.createAASRTask(audioURL, format, pid)
	if err != nil {
		return nil, err
	}

	const maxTry = 120
	for i := 0; i < maxTry; i++ {
		finish, ret, err := this.queryAASRTask(task)
		if err != nil {
			return nil, err
		}
		if finish {
			return ret, nil
		}
		time.Sleep(time.Second * 2)
	}
	return nil, fmt.Errorf("max try")
}
