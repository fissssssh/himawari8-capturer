package core

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type tileImage struct {
	X     int
	Y     int
	Image image.Image
}

// request image concurrent limit
const concurrency = 10
const tileWidth = 550
const tileEarthUrlTemplate = "https://himawari8.nict.go.jp/img/D531106/%dd/550/%d/%02d/%02d/%02d%d000_%d_%d.png"
const tileShorelineUrlTemplate = "https://himawari8.nict.go.jp/img/D531106/%dd/550/coastline/%s_%d_%d.png"
const latestEarthTimeUrlTemplate = "https://himawari8.nict.go.jp/img/FULL_24h/latest.json?_=%d"
const latestDateLayout = "2006-01-02 15:04:05"

var client http.Client = http.Client{}

func GetLatestEarthTime() (time.Time, error) {
	url := fmt.Sprintf(latestEarthTimeUrlTemplate, time.Now().UnixMilli())
	resp, err := http.Get(url)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()
	var latest map[string]string
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return time.Time{}, err
	}
	json.Unmarshal(bodyBytes, &latest)
	if date, exist := latest["date"]; exist {
		return time.Parse(latestDateLayout, date)
	}
	return time.Time{}, errors.New("response not contains date")
}

func GetEarthWithShorelines(q Quality, t time.Time, c Shoreline) (io.Reader, error) {
	earth, err := getEarth(q, t)
	if err != nil {
		return nil, errors.New("get earth image failed")
	}
	r := image.NewRGBA(image.Rect(0, 0, int(q)*tileWidth, int(q)*tileWidth))
	draw.Draw(r, earth.Bounds(), earth, earth.Bounds().Min, draw.Src)
	if c != Ignore {
		shorelines, err := getShorelines(q, c)
		if err != nil {
			return nil, errors.New("get shorelines image failed")
		}
		draw.Draw(r, shorelines.Bounds(), shorelines, shorelines.Bounds().Min, draw.Over)
	}
	var buf bytes.Buffer
	err = png.Encode(&buf, r)
	if err != nil {
		return nil, errors.New("png encode failed")
	}
	return &buf, nil
}

func getEarth(q Quality, t time.Time) (image.Image, error) {
	t = t.UTC()
	year, month, day := t.Date()
	hour, minute := t.Hour(), t.Minute()
	return getTilesAndCombine(q, func(x, y int) (image.Image, error) {
		url := fmt.Sprintf(tileEarthUrlTemplate, q, year, month, day, hour, minute/10, x, y)
		return getImage(url)
	})
}

func getShorelines(q Quality, c Shoreline) (image.Image, error) {
	if c == Ignore {
		return nil, errors.New("shoreline can not be ignore")
	}
	return getTilesAndCombine(q, func(x, y int) (image.Image, error) {
		url := fmt.Sprintf(tileShorelineUrlTemplate, q, c, x, y)
		return getImage(url)
	})
}

func getTilesAndCombine(q Quality, tileGetter func(x, y int) (image.Image, error)) (image.Image, error) {
	level := int(q)
	concurrentControlCh := make(chan struct{}, concurrency)
	tileCh := make(chan tileImage, concurrency)
	defer close(concurrentControlCh)
	defer close(tileCh)
	wg := sync.WaitGroup{}
	// request image from internet
	for x := 0; x < level; x++ {
		for y := 0; y < level; y++ {
			wg.Add(1)
			concurrentControlCh <- struct{}{}
			go func(x int, y int) {
				tile, err := tileGetter(x, y)
				<-concurrentControlCh
				if err == nil {
					tileCh <- tileImage{x, y, tile}
				}
			}(x, y)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// receive tile images from channel for combining
	r := image.NewRGBA(image.Rect(0, 0, level*tileWidth, level*tileWidth))
	go func() {
		for {
			select {
			case t := <-tileCh:
				draw.Draw(r, t.Image.Bounds().Add(image.Pt(t.X*tileWidth, t.Y*tileWidth)), t.Image, t.Image.Bounds().Min, draw.Src)
				wg.Done()
			case <-ctx.Done():
				return
			}
		}
	}()
	wg.Wait()
	return r, nil
}

func getImage(url string) (image.Image, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func SetHttpProxy(proxyUrl string) error {
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return err
	}
	transport := http.Transport{
		Proxy: func(r *http.Request) (*url.URL, error) {
			return proxy, nil
		},
	}
	client.Transport = &transport
	return nil
}
