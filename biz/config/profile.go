package config

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"tiktok/pkg/constant"
)

// UseProfile 激活active环境
func UseProfile(active string) {
	os.Setenv(constant.ProfileEnv, active)
}

func GetProfilePath(defaultPath, profile string) string {
	logrus.WithField("module", "config").Infof("active profile = %s", profile)
	filename := filepath.Base(defaultPath)
	split := strings.Split(filename, ".")
	return filepath.Join(filepath.Dir(defaultPath), split[0]+"-"+profile+"."+split[1])
}

func GetActiveProfile() string {
	return os.Getenv(constant.ProfileEnv)
}
