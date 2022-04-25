package aliyun

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

type AASR struct {
	AppKey string
	Aliyun
}

func NewAASR(Aliyun Aliyun, AppKey string) AASR {
	return AASR{
		AppKey: AppKey,
		Aliyun: Aliyun,
	}
}

type AASRResult struct {
	TaskID      string    `json:"TaskId"`
	RequestID   string    `json:"RequestId"`
	StatusText  string    `json:"StatusText"`
	BizDuration int       `json:"BizDuration"`
	SolveTime   int64     `json:"SolveTime"`
	StatusCode  int       `json:"StatusCode"`
	Result      Sentences `json:"Result"`
}
type Sentence struct {
	EndTime         int64   `json:"EndTime"`
	SilenceDuration int     `json:"SilenceDuration"`
	BeginTime       int64   `json:"BeginTime"`
	Text            string  `json:"Text"`
	ChannelID       int     `json:"ChannelId"`
	SpeechRate      int     `json:"SpeechRate"`
	EmotionValue    float64 `json:"EmotionValue"`
}
type Sentences struct {
	Sentences []Sentence `json:"Sentences"`
}

func (z *AASR) Do(url string) ([]Sentence, error) {
	// 地域ID，固定值。
	const REGION_ID string = "cn-shanghai"
	const ENDPOINT_NAME string = "cn-shanghai"
	const PRODUCT string = "nls-filetrans"
	const DOMAIN string = "filetrans.cn-shanghai.aliyuncs.com"
	const API_VERSION string = "2018-08-17"
	const POST_REQUEST_ACTION string = "SubmitTask"
	const GET_REQUEST_ACTION string = "GetTaskResult"
	// 请求参数
	const KEY_APP_KEY string = "appkey"
	const KEY_FILE_LINK string = "file_link"
	const KEY_VERSION string = "version"
	const KEY_ENABLE_WORDS string = "enable_words"
	// 响应参数
	const KEY_TASK string = "Task"
	const KEY_TASK_ID string = "TaskId"
	const KEY_STATUS_TEXT string = "StatusText"
	const KEY_RESULT string = "Result"
	// 状态值
	const STATUS_SUCCESS string = "SUCCESS"
	const STATUS_RUNNING string = "RUNNING"
	const STATUS_QUEUEING string = "QUEUEING"

	client, err := sdk.NewClientWithAccessKey(REGION_ID, z.AK, z.AS)
	if err != nil {
		return nil, err
	}

	postRequest := requests.NewCommonRequest()
	postRequest.Domain = DOMAIN
	postRequest.Version = API_VERSION
	postRequest.Product = PRODUCT
	postRequest.ApiName = POST_REQUEST_ACTION
	postRequest.Method = "POST"
	mapTask := make(map[string]interface{})
	mapTask[KEY_APP_KEY] = z.AppKey
	mapTask[KEY_FILE_LINK] = url
	// 新接入请使用4.0版本，已接入（默认2.0）如需维持现状，请注释掉该参数设置。
	mapTask[KEY_VERSION] = "4.0"
	// 设置是否输出词信息，默认为false。开启时需要设置version为4.0。
	mapTask[KEY_ENABLE_WORDS] = "false"
	mapTask["first_channel_only"] = true
	mapTask["enable_timestamp_alignment"] = true
	task, err := json.Marshal(mapTask)
	if err != nil {
		return nil, err
	}
	postRequest.FormParams[KEY_TASK] = string(task)
	postResponse, err := client.ProcessCommonRequest(postRequest)
	if err != nil {
		return nil, err
	}
	postResponseContent := postResponse.GetHttpContentString()
	if postResponse.GetHttpStatus() != 200 {
		return nil, fmt.Errorf("录音文件识别请求失败，Http错误码: %v", postResponse.GetHttpStatus())
	}
	var postMapResult map[string]interface{}
	err = json.Unmarshal([]byte(postResponseContent), &postMapResult)
	if err != nil {
		return nil, err
	}
	var taskId string = ""
	var statusText string = ""
	statusText = postMapResult[KEY_STATUS_TEXT].(string)
	if statusText == STATUS_SUCCESS {
		taskId = postMapResult[KEY_TASK_ID].(string)
	} else {
		return nil, fmt.Errorf("录音文件识别请求失败!")
	}

	getRequest := requests.NewCommonRequest()
	getRequest.Domain = DOMAIN
	getRequest.Version = API_VERSION
	getRequest.Product = PRODUCT
	getRequest.ApiName = GET_REQUEST_ACTION
	getRequest.Method = "GET"
	getRequest.QueryParams[KEY_TASK_ID] = taskId
	statusText = ""

	maxTry := 100000
	for i := 0; i < maxTry; i++ {
		getResponse, err := client.ProcessCommonRequest(getRequest)
		if err != nil {
			return nil, err
		}
		getResponseContent := getResponse.GetHttpContentString()
		if getResponse.GetHttpStatus() != 200 {
			return nil, fmt.Errorf("识别结果查询请求失败，Http错误码：%v", getResponse.GetHttpStatus())
		}
		var getMapResult AASRResult
		err = json.Unmarshal([]byte(getResponseContent), &getMapResult)
		if err != nil {
			return nil, err
		}
		statusText = getMapResult.StatusText
		if statusText == STATUS_SUCCESS {
			return getMapResult.Result.Sentences, nil
		}
		time.Sleep(10 * time.Second)
	}
	return nil, fmt.Errorf("max try %d", maxTry)
}
