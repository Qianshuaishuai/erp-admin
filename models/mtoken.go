package models

import (
	"fmt"
	"math/rand"
	"time"
	"dreamEbagPaperAdmin/helper"
)

func GetSaveToken() string {
	//token = 随机8位数 + 时间戳
	ntime := time.Now().Unix()
	rnd := rand.New(rand.NewSource(ntime))
	vcode := fmt.Sprintf("%08v", rnd.Int31n(10000000))

	return vcode + helper.Int64ToString(ntime)
}
