package models

import (
	"encoding/json"
	"erp-admin/helper"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	loglib "github.com/HYY-yu/LogLib"
)

var (
	curlIdClient *http.Client
)

type CurlReseponIntId struct {
	F_id int `json:"F_id"`
}

func init() {
	curlIdClient = &http.Client{
		Transport: &http.Transport{
			//			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression:    true,
			ResponseHeaderTimeout: time.Second * 0,
		},
	}
}

type MSnowflakCurl struct {
}

//获取发号器发出的ID(int类型,16位)
func (u *MSnowflakCurl) GetIntId(test bool) (id int) {
	id = 0
	uniqueFlag := helper.GetGuid()

	var uri, method string
	var req *http.Request

	if test {
		uri = MyConfig.SnowTestFlakDomain + "/v1/snowflak/intId"
		method = "GET"
		req, _ = http.NewRequest(method, uri, nil)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authentication", MyConfig.SnowTestFlakAuthUser+":"+helper.Md5([]byte(MyConfig.SnowTestFlakAuthUserSecurity)))
	} else {
		uri = MyConfig.SnowProcFlakDomain + "/v1/snowflak/intId"
		method = "GET"
		req, _ = http.NewRequest(method, uri, nil)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authentication", MyConfig.SnowProcFlakAuthUser+":"+helper.Md5([]byte(MyConfig.SnowProcFlakAuthUserSecurity)))
	}

	client := curlIdClient

	//log request
	loglib.GetLogger().LogSnowflakRequest(uniqueFlag, uri, "")

	resp, err := client.Do(req)
	idObj := CurlReseponIntId{}
	if err == nil {
		defer resp.Body.Close()
		bodyByte, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			dec := json.NewDecoder(strings.NewReader(string(bodyByte)))
			dec.UseNumber()
			dec.Decode(&idObj)
		}
		if resp.Status == "200 OK" {
			id = idObj.F_id
		}
		//log response
		loglib.GetLogger().LogSnowflakResponse(uniqueFlag, strconv.Itoa(idObj.F_id), resp.Status, string(bodyByte))
	} else {
		//log err
		loglib.GetLogger().LogErr(err, "snowflak module")
	}
	return
}
