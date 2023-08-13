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

	RedisAddr string
	RedisPwd  string

	KafkaAddr string
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}
	Loadserver(file)
	LoadDB(file)
	LoadCos(file)
	LoadRedis(file)
}

func Loadserver(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("")
	Host = file.Section("server").Key("Host").MustString("")
	Port = file.Section("server").Key("HttpPort").MustString("")
	JwtKey = file.Section("server").Key("JwtKey").MustString("")
}
func LoadDB(file *ini.File) {
	DBHost = file.Section("mysql").Key("DBHosexitt").MustString("")
	DBPort = file.Section("mysql").Key("DBPort").MustString("")
	DBUser = file.Section("mysql").Key("DBUser").MustString("")
	DBPwd = file.Section("mysql").Key("DBPwd").MustString("")
	DBName = file.Section("mysql").Key("DBName").MustString("")
}
func LoadCos(file *ini.File) {
	CosUrl = file.Section("cos").Key("CosUrl").MustString("")
	CosImageUrl = file.Section("cos").Key("CosImageUrl").MustString("")
	SecretId = file.Section("cos").Key("SecretId").MustString("")
	SecretKey = file.Section("cos").Key("SecretKey").MustString("")
}
func LoadRedis(file *ini.File) {
	RedisAddr = file.Section("redis").Key("RedisAddr").MustString("")
	RedisPwd = file.Section("redis").Key("RedisPwd").MustString("")
}
func LoadKafka(file *ini.File) {
	KafkaAddr = file.Section("kafka").Key("KafkaAddr").MustString("")
}
