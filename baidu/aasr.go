package baidu

import (
	"fmt"
	"time"

	"github.com/byebyebruce/autopublish/pkg/util"
)

// AASR_TASK_URL https://cloud.baidu.com/doc/SPEECH/s/ck5diijkt

func (this *API_Util) createAASRTask(audioURL string,
	format string, /*wav pcm*/
	pid int /*[80001（中文语音近场识别模型极速版）, 1737（英文模型）] */) (string, error) {

	const (
		AASR_TASK_URL = "https://aip.baidubce.com/rpc/2.0/aasr/v1/create"
	)
	type aasrReq struct {
		SpeechURL string `json:"speech_url"`
		Format    string `json:"format"` // "mp3", "wav", "pcm","m4a","amr"
		Pid       int    `json:"pid"`    // [80001（中文语音近场识别模型极速版）, 1737（英文模型）]
		Rate      int    `json:"rate"`
	}

	type aasrResp struct {
		LogID      int    `json:"log_id"`
		TaskStatus string `json:"task_status"`
		TaskID     string `json:"task_id"`
		ErrorCode  int    `json:"error_code"`
		ErrorMsg   string `json:"error_msg"`
	}

	url := fmt.Sprintf("%s?cuid=%v&access_token=%s", AASR_TASK_URL, this.Cuid, this.Credentials.Access_token)

	req := aasrReq{
		SpeechURL: audioURL,
		Format:    format,
		Pid:       pid,
		Rate:      16000,
	}
	resp := &aasrResp{}
	err := util.HTTPPostJSON(url, req, resp)
	if err != nil {
		return "", err
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("code:%d msg:%v", resp.ErrorCode, resp.ErrorMsg)
	}
	return resp.TaskID, nil
}

type DetailedResult []struct {
	Res       []string `json:"res"`
	EndTime   int      `json:"end_time"`
	BeginTime int      `json:"begin_time"`
}

func (this *API_Util) queryAASRTask(taskId string) (bool, DetailedResult, error) {
	const (
		AASR_TASK_RESULT_URL = "https://aip.baidubce.com/rpc/2.0/aasr/v1/query"
	)

	type aasrResultRequest struct {
		TaskIds []string `json:"task_ids"`
	}

	type aasrResult struct {
		TaskStatus string `json:"task_status"`
		TaskResult struct {
			ErrNo  int    `json:"err_no"`
			ErrMsg string `json:"err_msg"`

			Result         []string       `json:"result"`
			AudioDuration  int            `json:"audio_duration"`
			DetailedResult DetailedResult `json:"detailed_result"`
			CorpusNo       string         `json:"corpus_no"`
		} `json:"task_result"`
		TaskID string `json:"task_id"`
	}
	type aasrResultResp struct {
		LogID     int          `json:"log_id"`
		TasksInfo []aasrResult `json:"tasks_info"`
	}
	url := fmt.Sprintf("%s?cuid=%v&access_token=%s", AASR_TASK_RESULT_URL, this.Cuid, this.Credentials.Access_token)
	req := aasrResultRequest{
		TaskIds: []string{taskId},
	}
	resp := &aasrResultResp{}
	err := util.HTTPPostJSON(url, req, resp)
	if err != nil {
		return true, nil, err
	}
	if len(resp.TasksInfo) == 0 {
		return true, nil, fmt.Errorf("no task")
	}
	taskResult := resp.TasksInfo[0]
	switch taskResult.TaskStatus {
	case "Running":
		return false, nil, nil
	case "Failure":
		return true, nil, fmt.Errorf("code %v msg %v", taskResult.TaskResult.ErrNo, taskResult.TaskResult.ErrMsg)
	case "Success":
		return true, taskResult.TaskResult.DetailedResult, nil
	default:
		return true, nil, fmt.Errorf("unknown")
	}
}

func (this *API_Util) QueryAASR(audioURL string,
	format string, /*wav pcm*/
	pid int /*[80001（中文语音近场识别模型极速版）, 1737（英文模型）] */) (DetailedResult, error) {
	this.genCredentials()
	task, err := this.createAASRTask(audioURL, format, pid)
	if err != nil {
		return nil, err
	}

	const maxTry = 30
	for i := 0; i < maxTry; i++ {
		finish, ret, err := this.queryAASRTask(task)
		if err != nil {
			return nil, err
		}
		if finish {
			return ret, nil
		}
		fmt.Printf("\r %d/%d", i, maxTry)
		time.Sleep(time.Second * 10)
	}
	return nil, fmt.Errorf("max try")
}
