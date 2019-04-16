package models

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

const (
	TYPE_ADD_PERSON_TAG        = 101
	TYPE_ADD_INDUSTRY_TAG      = 102
	TYPE_PROJECT_CARD_ID       = 103
	TYPE_PROJECT_BACKGROUND_ID = 104
	TYPE_CONNECTION_ICON_ID    = 105
	TYPE_CONNECTION_CARD_ID    = 106
)

func UploadFile(typeID int, fileName string, file multipart.File) (imageUrl string, err error) {
	var snowCurl MSnowflakCurl
	fileDir := "./admin-file/"
	os.RemoveAll(fileDir)
	os.Mkdir(fileDir, 0766)

	f, _ := os.OpenFile(fileDir+fileName, os.O_CREATE|os.O_RDWR, 0766)
	io.Copy(f, file)

	tag := ""
	nFileName := ""

	switch typeID {
	case TYPE_ADD_PERSON_TAG:
		tag = "person/tag"
		nFileName = fileName
		break
	case TYPE_ADD_INDUSTRY_TAG:
		tag = "industry/tag"
		nFileName = fileName
		break
	case TYPE_PROJECT_CARD_ID:
		tag = "project/system/card"
		nFileName = strconv.Itoa(snowCurl.GetIntId(true)) + "-card" + ".png"
		break
	case TYPE_PROJECT_BACKGROUND_ID:
		tag = "project/system/background"
		nFileName = strconv.Itoa(snowCurl.GetIntId(true)) + "-back" + ".png"
		break
	case TYPE_CONNECTION_ICON_ID:
		tag = "connection/system/icon"
		nFileName = strconv.Itoa(snowCurl.GetIntId(true)) + "-icon" + ".png"
		break
	case TYPE_CONNECTION_CARD_ID:
		tag = "connection/system/card"
		nFileName = strconv.Itoa(snowCurl.GetIntId(true)) + "-card" + ".png"
		break
	default:
		tag = "default"
		break
	}

	newFileName := "admin" + "/" + tag + "/" + nFileName
	imageURL, err := UploadFileToQiniu(newFileName, fileDir+fileName)

	os.RemoveAll(fileDir)

	return imageURL, err
}

//七牛云文件上传
func UploadFileToQiniu(fileName, filePath string) (url string, err error) {
	mac := qbox.NewMac(MyConfig.QiniuAccessKey, MyConfig.QiniuSecretKey)

	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", MyConfig.QiniuBucket, fileName),
	}

	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuanan
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	//可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "cache cut exam pic",
		},
	}

	err = formUploader.PutFile(context.Background(), &ret, upToken, fileName, filePath, &putExtra)

	if err != nil {
		return "", err
	}

	imageURL := MyConfig.QiniuBaseURL + ret.Key

	return imageURL, nil
}
