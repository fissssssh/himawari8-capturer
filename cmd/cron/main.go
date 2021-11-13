package main

import (
	"context"
	"fmt"
	"himawari8Capturer/core"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func main() {
	endPoint, secId, secKey := os.Args[1], os.Args[2], os.Args[3]

	latestTime, err := core.GetLatestEarthTime()
	if err != nil {
		log.Fatal("get latest earth time failed")
	}
	img, err := core.GetEarthWithShorelines(core.FHD, latestTime, core.Ignore)
	if err != nil {
		panic(err)
	}
	year, month, day := latestTime.Date()
	name := fmt.Sprintf("earth/2200x2200/%v-%v-%v/%v.png", year, int(month), day, latestTime.Unix())

	u, _ := url.Parse(endPoint)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secId,
			SecretKey: secKey,
		},
	})

	_, err = c.Object.Put(context.Background(), name, img, nil)
	if err != nil {
		panic(err)
	}
}
