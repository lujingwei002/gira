package kuaishou

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type JsCode2SessionRequest struct {
	AppId         string `json:"appid"`
	Secret        string `json:"secret"`
	Code          string `json:"code"`
	AnonymousCode string `json:"anonymous_code"`
}

type JsCode2SessionResponse struct {
	Result     int32  `json:"result"`
	OpenId     string `json:"open_id"`
	SessionKey string `json:"session_key"`
}

// https://developer.open-douyin.com/docs/resource/zh-CN/mini-game/develop/guide/game-engine/rd-to-SCgame/server/log-in/code-2-session
// 小程序登录
func JsCode2Session(appId string, secret string, code string, anonymousCode string) (*JsCode2SessionResponse, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	host := "https://open.kuaishou.com"
	params := url.Values{}
	params.Set("app_id", appId)
	params.Set("app_secret", secret)
	params.Set("js_code", code)
	// params.Set("anonymous_code", anonymousCode)

	url := fmt.Sprintf("%s/oauth2/mp/code2session?%s", host, params.Encode())
	log.Println(url)

	var httpReq *http.Request
	var result *http.Response
	var err error
	if err != nil {
		return nil, err
	}
	httpReq, err = http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Content-Type", "application/json")
	client := &http.Client{Transport: tr}
	// Send the HTTP request
	result, err = client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	var body []byte
	body, err = io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	resp := &JsCode2SessionResponse{}
	if err = json.Unmarshal(body, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
