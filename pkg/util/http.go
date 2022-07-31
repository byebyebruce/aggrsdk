package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func HttpGet(u string, data interface{}) error {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = json.Unmarshal(body, data); nil != err {
		return err
	}
	return nil
}

func HttpPostForm(u string, req map[string]interface{}, data interface{}) error {
	param := url.Values{}
	for k, v := range req {
		param.Set(k, fmt.Sprint(v))
	}
	resp, err := http.PostForm(u, param)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = json.Unmarshal(body, data); nil != err {
		return err
	}
	return nil
}

func HTTPPostJSON(u string, req, resp interface{}) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(b))
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 3}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http code %d", httpResp.StatusCode)
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	if err := json.Unmarshal(body, resp); err != nil {
		return err
	}
	return nil
}

func HTTPPostRaw(u string, h map[string]string, r io.Reader, resp interface{}) error {
	httpReq, err := http.NewRequest(http.MethodPost, u, r)
	for k, v := range h {
		httpReq.Header.Set(k, v)
	}

	client := &http.Client{Timeout: time.Second * 3}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, resp); err != nil {
		return err
	}
	return nil
}
