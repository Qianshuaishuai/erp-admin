package models

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/astaxie/beego"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

const (
	TYPE_USER_HEAD_ID          = 101
	TYPE_CONNECTION_CARD_ID    = 102
	TYPE_PROJECT_CARD_ID       = 103
	TYPE_PROJECT_BACKGROUND_ID = 104
)

func UploadFile(typeID int, fileName string, file multipart.File) (imageUrl string, err error) {
	fileDir := "./"
	os.RemoveAll(fileDir)
	os.Mkdir(fileDir, 0766)

	f, _ := os.OpenFile(fileDir+fileName, os.O_CREATE|os.O_RDWR, 0766)
	io.Copy(f, file)

	tag := ""

	switch typeID {
	case TYPE_USER_HEAD_ID:
		tag = "headicon"
		break
	case TYPE_CONNECTION_CARD_ID:
		tag = "card"
		break
	case TYPE_PROJECT_CARD_ID:
		tag = "project/card/"
		break
	case TYPE_PROJECT_BACKGROUND_ID:
		tag = "project/background/"
		break
	default:
		tag = "default"
		break
	}

	beego.Debug(tag)

	newFileName := "admin" + "/" + fileName
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
