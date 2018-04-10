package helper

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	cyrand "crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"math/rand"
)

var uniqueInt64 uint64

func init() {
	uniqueInt64 = 0
}

//类型转化 string  to int
func StrToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

//类型转化 string  to uint64
func StrToUint64(str string) uint64 {
	i, _ := strconv.ParseUint(str, 0, 64)
	return i
}

//类型转化 string  to float64
func StrToFloat64(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

//类型转化 int to string
func IntToString(i int) string {
	return fmt.Sprintf("%d", i)
}

//类型转化 int64 to string
func Int64ToString(i int64) string {
	return fmt.Sprintf("%d", i)
}

//类型转化 uint64 to string
func Uint64ToString(i uint64) string {
	return fmt.Sprintf("%d", i)
}

//类型转化 uint32 to string
func Uint32ToString(i uint32) string {
	return fmt.Sprintf("%d", i)
}

//类型转换inerface to string
func InterfaceToString(data interface{}) string {
	return fmt.Sprintf("%s", data)
}

//get now datatime(Y-m-d H:i:s)
func GetNowDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//get now datatime2(YmdHis)
func GetNowDateTime2() string {
	return time.Now().Format("20060102150405")
}

//get now datatime3(Y-m-d)
func GetNowDateTime3() string {
	return time.Now().Format("2006-01-02")
}

//get now datatime4(H:i:s)
func GetNowDateTime4() string {
	return time.Now().Format("15:04:05")
}

//get yestoday(Y-m-d)
func GetYestoday() string {
	return time.Now().Add(-time.Minute * (time.Duration(24*60) - 1)).Format("2006-01-02")
}

//获取当前的时间字符串
func GetNowDateTimeDefault() string {
	return time.Now().String()
}

//获取当前几分钟前的时间(Y-m-d H:i:s)
func GetDateTimeBeforeMinute(num int) string {
	return time.Now().Add(-time.Minute * time.Duration(num)).Format("2006-01-02 15:04:05")
}

//获取当前几秒钟前的时间(Y-m-d H:i:s)
func GetDateTimeBeforeSecond(num int) string {
	return time.Now().Add(-time.Second * time.Duration(num)).Format("2006-01-02 15:04:05")
}

//获取当前几分钟后的时间(Y-m-d H:i:s)
func GetDateTimeAfterMinute(num int) string {
	return time.Now().Add(time.Minute * time.Duration(num)).Format("2006-01-02 15:04:05")
}

//把一个时间字符串转为unix时间戳
func StrToTimeStamp(timeStr string) int64 {
	//	time = "2015-09-14 16:33:00"
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	return t.Unix()
}

//把一个unix时间戳转为Y-m-d H:i:s格式的日期
func TimeStampToStr(timeStamp int64) string {
	timeObj := time.Unix(timeStamp-int64(8*60*60), int64(0))
	return timeObj.Format("2006-01-02 15:04:05")
}

//把一个unix时间戳转为Y-m-d格式的日期
func TimeStampToStr2(timeStamp int64) string {
	timeObj := time.Unix(timeStamp-int64(8*60*60), int64(0))
	return timeObj.Format("2006-01-02")
}

//切分一个字符串为字符串数组
func Split(str string, flag string) []string {
	return strings.Split(str, flag)
}

//合并字符串数组
func JoinString(list []string, flag string) string {
	result := ""
	if len(list) > 0 {
		for _, v := range list {
			result += v + flag
		}
		result = strings.Trim(result, flag)
	}
	return result
}

//合并int数组
func JoinInt(list []int, flag string) string {
	result := ""
	if len(list) > 0 {
		for _, v := range list {
			result += IntToString(v) + flag
		}
		result = strings.Trim(result, flag)
	}
	return result
}

func JoinInt64(list []int64, flag string) string {
	result := ""
	if len(list) > 0 {
		for _, v := range list {
			result += Int64ToString(v) + flag
		}
		result = strings.Trim(result, flag)
	}
	return result
}

//检查一个字符串是否在字符串数组里面
func StringInArray(value string, list []string) bool {
	result := false
	for _, item := range list {
		if value == item {
			result = true
			break
		}
	}
	return result
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// 检查是否为文件且存在
// 如果由 filename 指定的文件存在则返回 true，否则返回 false
func Exist2(filename string) (exists bool) {
	exists = false
	fileInfo, err := os.Stat(filename)
	if err == nil || os.IsExist(err) {
		if !fileInfo.IsDir() {
			exists = true
		}
	}
	return
}

//sha1
func Sha1(str string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(str)))
}

//获取interface的类型
func GetInterfaceType(i interface{}) string {
	typeObj := reflect.TypeOf(i)
	return typeObj.Kind().String()
}

//检查interface数据是否为string类型
func CheckInterfaceIsString(i interface{}) bool {
	if i != nil {
		if GetInterfaceType(i) == "string" {
			return true
		}
	}
	return false
}

//生成Guid字串(32位)
func GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(cyrand.Reader, b); err != nil {
		return ""
	}
	return Md5([]byte(base64.URLEncoding.EncodeToString(b)))
}

//创建一个20位的唯一数字字符串(flag为两位数的字符数字字符串)
func GetUqunieNumString20(flag string) string {
	nowTime := time.Now().Format("20060102150405")
	atomic.CompareAndSwapUint64(&uniqueInt64, uint64(9999), uint64(0))
	atomic.AddUint64(&uniqueInt64, 1)
	tmp := Uint64ToString(uniqueInt64)
	if len(tmp) == 1 {
		tmp = "000" + tmp
	} else if len(tmp) == 2 {
		tmp = "00" + tmp
	} else if len(tmp) == 3 {
		tmp = "0" + tmp
	}
	id := flag + nowTime + tmp
	return id
}

//copy file
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

//深度复制一个对象
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

//去重一个int数组
func RmDuplicateInt(list *[]int) []int {
	var x []int = []int{}
	for _, i := range *list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

//检查某个汉字有效性
func CheckHanziValid(hanzi string) bool {
	matched, err := regexp.MatchString("^[\u4e00-\u9fa5]{1}$", hanzi)
	if err == nil && matched {
		return true
	}
	return false
}

//uint数组中是否有元素elem
func UIntContainer(values []uint, elem uint) bool {
	found := false
	for _, searchValue := range values {
		if elem == searchValue {
			found = true
			break
		}
	}
	if !found {
		return false
	}
	return true
}

//将[id,id,id]字符串转换成id数组
func TransformStringToInt64Arr(idsString string) ([]int64, error) {
	resourceIdList := make([]int64, 0)
	dec := json.NewDecoder(strings.NewReader(idsString))
	dec.UseNumber()
	errJ := dec.Decode(&resourceIdList)
	return resourceIdList, errJ
}

//将id数组转[id,id,id]字符串
func TransformInt64ArrToString(idsString []int64) string {
	return "[" + JoinInt64(idsString, ",") + "]"
}

func Md5(buf []byte) string {
	hash := md5.New()
	hash.Write(buf)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func Password(len int, pwdO string) (pwd string, salt string) {
	salt = GetRandomString(len)
	defaultPwd := "dream123"
	if pwdO != "" {
		defaultPwd = pwdO
	}
	pwd = Md5([]byte(defaultPwd + salt))
	return pwd, salt
}

//生成随机字符串
func GetRandomString(lens int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lens; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
