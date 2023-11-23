package douyin

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type JsCode2SessionRequest struct {
	AppId         string `json:"appid"`
	Secret        string `json:"secret"`
	Code          string `json:"code"`
	AnonymousCode string `json:"anonymous_code"`
}

type JsCode2SessionResponse struct {
	ErrCode int32  `json:"err_code"`
	ErrMsg  string `json:"err_tips"`
	Data    struct {
		OpenId          string `json:"openid"`
		AnonymousOpenId string `json:"anonymous_openid"`
		SessionKey      string `json:"session_key"`
	} `json:"data"`
}

// https://developer.open-douyin.com/docs/resource/zh-CN/mini-game/develop/guide/game-engine/rd-to-SCgame/server/log-in/code-2-session
// 小程序登录
func JsCode2Session(appId string, secret string, code string, anonymousCode string) (*JsCode2SessionResponse, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	host := "https://minigame.zijieapi.com"
	url := fmt.Sprintf("%s/mgplatform/api/apps/jscode2session", host)
	req := &JsCode2SessionRequest{
		AppId:         appId,
		Secret:        secret,
		Code:          code,
		AnonymousCode: anonymousCode,
	}
	var httpReq *http.Request
	var result *http.Response
	var data []byte
	var err error
	data, err = json.Marshal(req)
	if err != nil {
		return nil, err
	}
	log.Println(string(data))
	httpReq, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
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
