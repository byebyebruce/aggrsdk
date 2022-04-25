package baidu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const TSN_URL = "http://tsn.baidu.com/text2audio"

func (this *API_Util) Text2AudioFile(filePath, text string) error {
	body, err := this.Text2AudioBytes(text)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, body, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (this *API_Util) Text2AudioBytes(text string) ([]byte, error) {
	this.genCredentials()
	param := url.Values{}
	param.Set("tex", text)
	param.Set("tok", this.Credentials.Access_token)
	param.Set("cuid", this.Cuid)
	param.Set("ctp", "1")
	param.Set("lan", "zh")

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
	if "audio/mp3" == contentType {
		return body, nil
	} else {
		var errMsg API_Response
		err = json.Unmarshal(body, &errMsg)
		if nil != err {
			return nil, err
		}
		return nil, fmt.Errorf("%+v", errMsg.Err_msg)
	}
}
