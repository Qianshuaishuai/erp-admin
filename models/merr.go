package models

import "github.com/jinzhu/gorm"

func HandleErrByTx(err error, tx *gorm.DB) error {
	if tx != nil{
		tx.Rollback()
	}
	//处理err

	return err
}
