package wechatapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/lujingwei002/gira"
)

type JsCode2SessionResponse struct {
	ErrCode    int32  `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	OpenId     string `json:"openid"`
}

const (
	JsCode2SessionGrantType_AuthorizationCode = "authorization_code "
)

// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html
// 小程序登录
func JsCode2Session(appId string, secret string, jsCode string, grantType string) (*JsCode2SessionResponse, error) {
	host := "https://api.weixin.qq.com"
	params := url.Values{}
	params.Set("appid", appId)
	params.Set("secret", secret)
	params.Set("js_code", jsCode)
	params.Set("grant_type", grantType)
	url := fmt.Sprintf("%s/sns/jscode2session?%s", host, params.Encode())
	var result *http.Response
	var err error
	result, err = http.Get(url)
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

type GetAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

const (
	GetAccessTokenGrantType_AuthorizationCode = "client_credential"
)

// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-access-token/getAccessToken.html
// 小程序登录
func GetAccessToken(appId string, secret string, grantType string) (*GetAccessTokenResponse, error) {
	host := "https://api.weixin.qq.com"
	params := url.Values{}
	params.Set("appid", appId)
	params.Set("secret", secret)
	params.Set("grant_type", grantType)
	url := fmt.Sprintf("%s/cgi-bin/token?%s", host, params.Encode())
	var result *http.Response
	var err error
	result, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	var body []byte
	body, err = io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	resp := &GetAccessTokenResponse{}
	if err = json.Unmarshal(body, resp); err != nil {
		return nil, err
	}
	log.Println("f", resp)
	return resp, nil
}

type GetWxaCodeUnLimitResponse struct {
	ErrCode     int32  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	Buffer      []byte
	ContentType string
}

type GetWxaCodeUnLimitRequest struct {
	Scene      string `json:"scene"`
	Page       string `json:"page"`
	EnvVersion string `json:"env_version"`
}

// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/qr-code/getUnlimitedQRCode.html
// 小程序登录
func GetWxaCodeUnLimit(accessToken string, scene string, page string, envVersion string) (*GetWxaCodeUnLimitResponse, error) {
	host := "https://api.weixin.qq.com"
	params := url.Values{}
	params.Set("access_token", accessToken)
	url := fmt.Sprintf("%s/wxa/getwxacodeunlimit?%s", host, params.Encode())
	req := &GetWxaCodeUnLimitRequest{
		Scene:      scene,
		Page:       page,
		EnvVersion: envVersion,
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
	client := &http.Client{}
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
	resp := &GetWxaCodeUnLimitResponse{
		ContentType: result.Header.Get("Content-Type"),
	}
	log.Println(result.Header)
	log.Println(resp.ContentType)
	if strings.HasPrefix(resp.ContentType, "image/jpeg") {
		resp.Buffer = body
		return resp, nil
	} else if strings.HasPrefix(resp.ContentType, "application/json") {
		// log.Println(string(body))
		if err = json.Unmarshal(body, resp); err != nil {
			return nil, err
		}
		return resp, nil
	} else {
		return nil, gira.ErrTodo
	}
}
