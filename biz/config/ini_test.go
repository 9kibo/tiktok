package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"testing"
)

func TestIniLoad(t *testing.T) {
	// Load 后面的配置源会替换前面的
	load, err := ini.Load("../../config.ini", "../../config-dev.ini")
	if err != nil {
		panic(err)
	}
	section := load.Section("Mysql")
	mysqlHost := section.Key("Host")
	fmt.Println(mysqlHost)

	c := Configuration{}
	err = load.MapTo(&c)
	if err != nil {
		panic(err)
	}
	fmt.Println(c)

}
