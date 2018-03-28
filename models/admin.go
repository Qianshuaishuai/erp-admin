/**********************************************
** @Des: This file ...
** @Author: haodaquan
** @Date:   2017-09-16 15:42:43
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 11:48:17
***********************************************/
package models

import (
	"dreamEbagPaperAdmin/helper"
	"errors"
)

type Admin struct {
	Id        int
	LoginName string
	Password  string
	Salt      string
	Status    int
}

var defaultAdminList []*Admin

func init() {
	pwd1, salt1 := helper.Password(4, "Admin123")

	defaultAdminList = []*Admin{
		{
			Id:        1000,
			LoginName: "Admin",
			Password:  pwd1,
			Salt:      salt1,
			Status:    1,
		},
	}
}

func AdminGetByName(loginName string) (*Admin, error) {
	for i := range defaultAdminList {
		if defaultAdminList[i].LoginName == loginName {
			return defaultAdminList[i], nil
		}
	}

	return nil,errors.New("没有这个账号")
}

func AdminGetById(id int) (*Admin, error) {
	for i := range defaultAdminList {
		if defaultAdminList[i].Id == id {
			return defaultAdminList[i], nil
		}
	}

	return nil,errors.New("没有这个账号")
}
