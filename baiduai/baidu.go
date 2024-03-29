package baiduai

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

/*
	客户凭证授权
	获取调用API的token
*/
const Credentials_Url = "https://aip.baidubce.com/oauth/2.0/token"

var Credentials_ResponseErrEnum map[string]Credentials_ResponseErr

func init() {
	Credentials_ResponseErrEnum = map[string]Credentials_ResponseErr{
		"invalid_request":           {Error: "invalid_request", Error_description: "invalid refresh token", Description: "请求缺少某个必需参数，包含一个不支持的参数或参数值，或者格式不正确。"},
		"invalid_client":            {Error: "invalid_client", Error_description: "unknown client id", Description: "client_id或client_secret参数无效"},
		"invalid_grant":             {Error: "invalid_grant", Error_description: "The provided authorization grant is revoked", Description: "提供的Access Grant是无效的、过期的或已撤销的，例如，Authorization Code无效(一个授权码只能使用一次)、Refresh Token无效、redirect_uri与获取Authorization Code时提供的不一致、Devie Code无效(一个设备授权码只能使用一次)等。"},
		"unauthorized_client":       {Error: "unauthorized_client", Error_description: "The client is not authorized to use this authorization grant type", Description: "应用没有被授权，无法使用所指定的grant_type。"},
		"unsupported_grant_type":    {Error: "unsupported_grant_type", Error_description: "The authorization grant type is not supported", Description: "“grant_type”百度OAuth2.0服务不支持该参数。"},
		"invalid_scope":             {Error: "invalid_scope", Error_description: "The requested scope is exceeds the scope granted by the resource owner", Description: "请求的“scope”参数是无效的、未知的、格式不正确的、或所请求的权限范围超过了数据拥有者所授予的权限范围。"},
		"expired_token":             {Error: "expired_token", Error_description: "refresh token has been used", Description: "提供的Refresh Token已过期"},
		"redirect_uri_mismatch":     {Error: "redirect_uri_mismatch", Error_description: "Invalid redirect uri", Description: "“redirect_uri”所在的根域与开发者注册应用时所填写的根域名不匹配。"},
		"unsupported_response_type": {Error: "unsupported_response_type", Error_description: "The response type is not supported", Description: "“response_type”参数值不为百度OAuth2.0服务所支持，或者应用已经主动禁用了对应的授权模式"},
		"slow_down":                 {Error: "slow_down", Error_description: "The device is polling too frequently", Description: "Device Flow中，设备通过Device Code换取Access Token的接口过于频繁，两次尝试的间隔应大于5秒。"},
		"authorization_pending":     {Error: "authorization_pending", Error_description: "User has not yet completed the authorization", Description: "Device Flow中，用户还没有对Device Code完成授权操作。"},
		"authorization_declined":    {Error: "authorization_declined", Error_description: "User has declined the authorization", Description: "Device Flow中，用户拒绝了对Device Code的授权操作。"},
		"invalid_referer":           {Error: "invalid_referer", Error_description: "Invalid Referer", Description: "Implicit Grant模式中，浏览器请求的Referer与根域名绑定不匹配"},
	}
}

type Credentials_Request struct {
	Grant_type    string // 必填参数 固定为“client_credentials”；
	Client_id     string // 必填参数 应用的API Key
	Client_secret string // 必须参数 应用的Secret Key;
	/*
	  非必须参数。
	  以空格分隔的权限列表，采用本方式获取Access Token时只能申请跟用户数据无关的数据访问权限。
	  关于权限的具体信息请参考
	  http://developer.baidu.com/wiki/index.php?title=docs/oauth/list
	*/
	Scope string
}

type Credentials_Response struct {
	Access_token   string `json:"access_token"`   // 要获取的Access Token
	Expires_in     int    `json:"expires_in"`     // Access Token的有效期,以秒为单位
	Refresh_token  string `json:"refresh_token"`  // 用于刷新Access Token 的 Refresh Token,所有应用都会返回该参数;（10年的有效期）
	Session_key    string `json:"session_key"`    // 基于http调用Open API时所需要的Session Key,其有效期与Access Token一致;
	Session_secret string `json:"session_secret"` // 基于http调用Open API时计算参数签名用的签名密钥.
	/*
		Access Token最终的访问范围，
		即用户实际授予的权限列表（用户在授权页面时，有可能会取消掉某些请求的权限），
		关于权限的具体信息参考
		http://developer.baidu.com/wiki/index.php?title=docs/oauth/list
	*/
	Scope string `json:"scope"`
}

type Credentials_ResponseErr struct {
	Error             string `json:"error"`
	Error_description string `json:"error_description"`
	Description       string
}

func GetCredentials(request Credentials_Request) Credentials_Response {

	postValue := url.Values{}
	// postValue.Set("scope", request.Scope)
	postValue.Set("client_id", request.Client_id)
	postValue.Set("grant_type", "client_credentials")
	postValue.Set("client_secret", request.Client_secret)

	postResponse, err := http.PostForm(Credentials_Url, postValue)
	if err != nil {
		panic(err.Error())
	}
	defer postResponse.Body.Close()

	body, err := ioutil.ReadAll(postResponse.Body)
	if err != nil {
		panic(err.Error())
	}

	var result Credentials_ResponseErr
	if err = json.Unmarshal(body, &result); nil == err {
		if description, ok := Credentials_ResponseErrEnum[result.Error]; ok {
			panic(description)
		}
	}

	var response Credentials_Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err.Error())
	}
	return response
}

/*
	获取一个本地的MAC地址作为API的 cuid
*/
func GetCUID() string {

	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err.Error())
	}

	var result string
	switch runtime.GOOS {
	case "windows":
		result = fmt.Sprintf("%s", interfaces[0].HardwareAddr)
	case "linux":
		result = fmt.Sprintf("%s", interfaces[1].HardwareAddr)
	default:
		result = "01:02:03:04:05:06"
	}

	return result
}

type API_Request struct {
	Tex  string `json:"tex"`           // 必填；合成文本，使用UTF-8编码，请注意文本长度必须小于1024
	Lan  string `json:"lan"`           // 必填；语言选择，填写zh
	Tok  string `json:"tok"`           // 必填；开放平台获取到的开发者access_token
	Ctp  int    `json:"ctp"`           // 必填；客户端类型选择,web端填写1
	Cuid string `json:"cuid"`          // 必填；用户唯一标识，用来区分用户，填写机器 MAC 地址或 IMEI 码，长度为60以内
	Spd  int    `json:"spd,omitempty"` // 选填；语速，取值0-9，默认5
	Pit  int    `json:"pit,omitempty"` // 选填；语调，取值0-9，默认5
	Vol  int    `json:"vol,omitempty"` // 选填；音量，取值0-9，默认5
	Per  int    `json:"per,omitempty"` // 选填；发音人选择，取值0-1；默认0女声 1男声
}

type API_Response struct {
	Err_no  int    `json:"err_no"`
	Err_msg string `json:"err_msg"`
	Sn      string `json:"sn"`
	Idx     int    `json:"idx"`
}

type BaiduAI struct {
	api_key, secret_key   string
	Credentials           Credentials_Response
	Cuid                  string
	lastGenCredentialTime int64
}

func NewBaiduAI(api_key, secret_key string) *BaiduAI {
	util := &BaiduAI{api_key: api_key, secret_key: secret_key}
	util.genCredentials()
	return util
}

func (this *BaiduAI) genCredentials() {
	if time.Now().Unix()-this.lastGenCredentialTime < 60*60*24*7 {
		return
	}
	cuid := GetCUID()

	res := GetCredentials(Credentials_Request{
		Client_id: this.api_key, Client_secret: this.secret_key})

	this.Cuid = cuid
	this.Credentials = res
	this.lastGenCredentialTime = time.Now().Unix()
}
