package utils

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
	"tiktok/config"
	"time"
)

func UpLoad(f multipart.File, userid string, fileHeader *multipart.FileHeader) (VideoUrl string, ImageUrl string, err error) {
	u, _ := url.Parse(config.CosUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv(config.SecretId),
			SecretKey: os.Getenv(config.SecretKey),
		},
	})
	//获取文件类型
	fileType := path.Ext(fileHeader.Filename)
	//通过 用户id+时间戳+文件名 生成key
	FileName := userid + strconv.FormatInt(time.Now().Unix(), 10) + fileHeader.Filename
	client.Object.Put(context.Background(), FileName, f, nil)
	VideoUrl = config.CosUrl + "/" + FileName
	ImageUrl = config.CosImageUrl + strings.TrimSuffix(FileName, fileType) + "image.jpg"
	return VideoUrl, ImageUrl, err
}
