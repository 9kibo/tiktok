package config

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var (
	AppMode string
	Host    string
	Port    string
	JwtKey  string

	DBHost string
	DBPort string
	DBUser string
	DBPwd  string
	DBName string

	CosUrl      string
	CosImageUrl string
	SecretId    string
	SecretKey   string
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}
	Loadserver(file)
	LoadDB(file)
	LoadCos(file)
}

func Loadserver(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	Host = file.Section("server").Key("Host").MustString("43.142.175.143")
	Port = file.Section("server").Key("HttpPort").MustString(":8080")
	JwtKey = file.Section("server").Key("JwtKey").MustString("tiktok")
}
func LoadDB(file *ini.File) {
	DBHost = file.Section("mysql").Key("DBHost").MustString("43.142.175.143")
	DBPort = file.Section("mysql").Key("DBPort").MustString("3306")
	DBUser = file.Section("mysql").Key("DBUser").MustString("root")
	DBPwd = file.Section("mysql").Key("DBPwd").MustString("TXY.zh2425904437")
	DBName = file.Section("mysql").Key("DBName").MustString("tiktok")
}
func LoadCos(file *ini.File) {
	CosUrl = file.Section("cos").Key("CosUrl").MustString("")
	CosImageUrl = file.Section("cos").Key("CosImageUrl").MustString("")
	SecretId = file.Section("cos").Key("SecretId").MustString("")
	SecretKey = file.Section("cos").Key("SecretKey").MustString("")
}
