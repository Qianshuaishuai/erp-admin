package models

import (
	"net/http"
	"time"
)

var (
	curlCheckAccessTokenClient *http.Client
)

func init() {
	curlCheckAccessTokenClient = &http.Client{
		Transport: &http.Transport{
			//			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression:    true,
			ResponseHeaderTimeout: time.Second * 5,
		},
	}
}

type MAuthCurl struct {
}

//验证accesstoken
//func (u *MAuthCurl) CheckAccessToken(userId string, accesstoken string, userType int, platform int) (respNo int) {
//	respNo = RESP_ERR
//	if len(userId) > 0 && len(accesstoken) > 0 && userType >= ROLE_STUDENT && userType <= ROLE_TEACHER && platform >= PLATFORM_ANDROID && platform <= PLATFORM_WEBCHAT {
//		uniqueLogFlag := helper.GetGuid()
//		//set request data
//		requestData := make(map[string]string)
//		requestData["userType"] = helper.IntToString(userType)
//		requestData["userId"] = userId
//		requestData["platform"] = helper.IntToString(platform)
//		requestData["ACCESSTOKEN"] = accesstoken
//		v := url.Values{}
//		for k1, v1 := range requestData {
//			v.Add(k1, v1)
//		}
//		values := v.Encode()
//		//
//		uri := MyConfig.LoginServerDomain + "/v1/auth/access"
//		method := "POST"
//		req, _ := http.NewRequest(method, uri, nil)
//		req.Body = ioutil.NopCloser(strings.NewReader(values))
//		//	req.Header.Set("Accept", "application/json")
//		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
//		client := curlCheckAccessTokenClient
//
//		//log request
//		var logObj *MLog
//		logObj.LogOtherCurlRequest(uniqueLogFlag, "向login中心'验证accesstoken'发起请求", uri, method, requestData)
//
//		//
//		type InfoTmp struct {
//			F_responseNo int `json:"F_responseNo"`
//		}
//
//		var responInfo InfoTmp
//
//		resp, err := client.Do(req)
//		if err == nil {
//			defer resp.Body.Close()
//			bodyByte, err := ioutil.ReadAll(resp.Body)
//			if err == nil {
//				dec := json.NewDecoder(strings.NewReader(string(bodyByte)))
//				dec.UseNumber()
//				dec.Decode(&responInfo)
//			}
//			if responInfo.F_responseNo > 0 {
//				respNo = responInfo.F_responseNo
//			}
//			//log response
//			logObj.LogOtherCurlResponse(uniqueLogFlag, "向login中心'验证accesstoken'发起请求响应", string(bodyByte), resp.Header, resp.Status)
//		} else {
//			//log err
//			logObj.LogOtherCurlResponseErr(uniqueLogFlag, "向login中心'验证accesstoken'发起请求响应报错:", err)
//		}
//	}
//	return
//}
