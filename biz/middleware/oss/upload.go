package oss

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"tiktok/biz/config"
	"time"
)

func UpLoad(f multipart.File, userid string, fileHeader *multipart.FileHeader) (VideoUrl string, ImageUrl string, err error) {
	cosConfig := config.C.Cos
	u, _ := url.Parse(cosConfig.Url)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv(cosConfig.SecretId),
			SecretKey: os.Getenv(cosConfig.SecretKey),
		},
	})
	//获取文件类型
	fileType := path.Ext(fileHeader.Filename)
	//通过 用户id+时间戳+文件名 生成key
	FileName := userid + strconv.FormatInt(time.Now().Unix(), 10) + fileHeader.Filename
	client.Object.Put(context.Background(), FileName, f, nil)
	VideoUrl = cosConfig.Url + "/" + FileName
	ImageUrl = cosConfig.ImageUrl + strings.TrimSuffix(FileName, fileType) + "image.jpg"
	return VideoUrl, ImageUrl, err
}
