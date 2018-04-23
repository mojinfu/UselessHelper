//
//  Created by 摸金校尉 on 17/11/20.
//  Copyright (c) 2018年 摸金校尉. All rights reserved.
//
package UselessHelper
import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	//"hash/crc32"
	"os"
	"path/filepath"
	"github.com/rlds/rlog"
	"strconv"
	"time"

	"strings"
	"regexp"
	"fmt"
)
func Logger(logPath string)bool{
	if !IsExist(logPath){
		fmt.Println("[创建日志文件夹]")
		if !MkAlldir(logPath){
			fmt.Println("[创建日志文件夹失败]")
			return false
		}
	}
	rlog.LogInit(3, GetConfPath(logPath), 1800 * 1024 * 1024, 1)
	rlog.V(1).Info( "即将启动服务:[" + "" + "]")
	return true
}
func UnicodeToChinese(unicodeStr string)string{
	defer RecoverMedicine("UnicodeToChinese")
	if len(unicodeStr)<2 {
		return unicodeStr
	}
	sUnicodev := strings.Split(unicodeStr, `\u`)
	var context string
	for ii, v := range sUnicodev {
		if len(v) < 1 {
			continue
		}
		if ii==0{
			context += v
			continue	
		}
		if len(v)<4{
			context += `\u`+v
			continue
		}
		temp, err := strconv.ParseInt(v[0:4], 16, 32)
		if err != nil {
			context += `\u`+v
			continue
		}
		context += fmt.Sprintf("%c", temp)+v[4:]
	}
	return context
}
func GetObjectValue(myJson string, key string)string{
	defer RecoverMedicine("getJsonValue")
	keyReg := regexp.MustCompile(key+`\s*:\s*"([^"]+)"`)
	QRURLArr :=keyReg.FindAllStringSubmatch(myJson,1)
	if len(QRURLArr)==0{
		return ""
	}
	if len(QRURLArr[0])<2{
		return ""
	}
	return QRURLArr[0][1]
}
func GetJsonValue(myJson string, key string)string{
	defer RecoverMedicine("getJsonValue")
	keyReg := regexp.MustCompile(key+`"\s*:\s*"([^"]+)"`)
	QRURLArr :=keyReg.FindAllStringSubmatch(myJson,1)
	if len(QRURLArr)==0{
		return ""
	}
	if len(QRURLArr[0])<2{
		return ""
	}
	return QRURLArr[0][1]
}
func TrimCannotbeseen (src string)(afterTrim string){
	defer RecoverMedicine("TrimCannotbeseen")
	afterTrim =strings.TrimFunc(src,func(w rune)bool{
		if w<32 {
			return true
		}
		if w=='\n'{
			return true
		}
		if w=='\t'{
			return true
		}
		if w=='\r'{
			return true
		}
		if w==' '{
			return true
		}
		return false
	})
	return 
}
func RecoverMedicine(funcname string){
	if err := recover(); err != nil {
			rlog.V(1).Info(funcname+" panic:",err)    
	}
}

func NowTime_s() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetAllFileData(filepath string) []byte {
	f, err := os.Open(filepath)
	if err != nil {
		//rlog.Info("GetAllFileDataErr_1:" + err.Error())
		return nil
	}
	var n int64

	if fi, err2 := f.Stat(); err2 == nil {
		if size := fi.Size(); size < 1e9 {
			n = size
		}
	} else {
		return nil
	}
	buf := bytes.NewBuffer(make([]byte, 0, n+bytes.MinRead))
	defer buf.Reset()
	_, err = buf.ReadFrom(f)
	f.Close()
	if err != nil {
		//rlog.Info("GetAllFileDataErr_2:" + err.Error())
		return nil
	}
	return buf.Bytes()
}

//存储并替换文件内容
func SaveReplaceFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	f.Write(data)
	f.Close()
	return nil
}

//判断文件是否存在
func IsFile(filepath string) error {
	_, err := os.Stat(filepath)
	return err
}
func IsExist(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}
func DelFile(path string) error {
	return os.Remove(path)
}

func MkAlldir(dir string) bool {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		//rlog.V(1).Info("建立文件夹出现错误[" + err.Error() + "]")
		return false
	}
	return true
}

//删除文件夹
func DelDir(path string) error {
	return os.RemoveAll(path)
}
func GetMd5Str(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

func GetConfPath(ipath string) (path string) {
	//file, _ := exec.LookPath(os.Args[0])
	var err error
	path, err = filepath.Abs(ipath)
	if err!=nil{
		fmt.Println("GetConfPath err")
	}
	// fmt.Println(path)
	// fmt.Println(path + "/../" + ipath)
	// cpath, _ = filepath.Abs(path + "/../" + ipath)
	return
}