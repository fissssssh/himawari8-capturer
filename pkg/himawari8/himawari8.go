package himawari8

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
	"net/http"
	"time"
)

type Quality uint8

const (
	Low Quality = 1 << iota
	HD
	FHD
	QHD
	UHD
)

type ShorelineColor uint8

const (
	Ignore ShorelineColor = iota
	Red
	Green
	Yellow
)

func (c ShorelineColor) String() string {
	switch c {
	case Ignore:
		return ""
	case Red:
		return "ff0000"
	case Green:
		return "00ff00"
	case Yellow:
		return "ffff00"
	default:
		return ""
	}
}

type signal struct {
	X     int
	Y     int
	Image image.Image
	Op    draw.Op
	Error error
}

const concurrency = 10
const tileWidth = 550
const tileEarthUrlTemplate = "https://himawari8.nict.go.jp/img/D531106/%dd/550/%d/%02d/%02d/%02d%d000_%d_%d.png"
const tileShorelineUrlTemplate = "https://himawari8.nict.go.jp/img/D531106/%dd/550/coastline/%s_%d_%d.png"

func GetImage(q Quality, t time.Time, c ShorelineColor) (io.Reader, error) {
	level := int(q)
	year, month, day := t.Date()
	hour, minute := t.Hour(), t.Minute()
	conCh := make(chan struct{}, concurrency)
	imgCh := make(chan signal)
	for x := 0; x < level; x++ {
		for y := 0; y < level; y++ {
			conCh <- struct{}{}
			go func(x int, y int, imgCh chan<- signal, conCh <-chan struct{}) {
				earth, err := getTileEarth(q, year, int(month), day, hour, minute, x, y)
				if err != nil {
					imgCh <- signal{x, y, nil, draw.Src, err}
				} else {
					imgCh <- signal{x, y, earth, draw.Src, nil}
				}
				if c != Ignore {
					shorelines, err := getTileShorelines(q, c, x, y)
					if err != nil {
						imgCh <- signal{x, y, nil, draw.Over, err}
					} else {
						imgCh <- signal{x, y, shorelines, draw.Over, nil}
					}
				}
				<-conCh
			}(x, y, imgCh, conCh)
		}
	}
	r := image.NewRGBA(image.Rect(0, 0, level*tileWidth, level*tileWidth))
	count := level * level
	if c != Ignore {
		count *= 2
	}
	for i := 0; i < count; i++ {
		s := <-imgCh
		if s.Error == nil {
			var t string
			if s.Op == draw.Src {
				t = "earth"
			} else {
				t = "shorelines"
			}
			log.Printf("plot %d-%d tile %s", s.X, s.Y, t)
			draw.Draw(r, s.Image.Bounds().Add(image.Pt(s.X*tileWidth, s.Y*tileWidth)), s.Image, s.Image.Bounds().Min, s.Op)
		}
	}
	var buf bytes.Buffer
	err := png.Encode(&buf, r)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

func getTileEarth(q Quality, year int, month int, day int, hour int, minute int, x int, y int) (image.Image, error) {
	url := fmt.Sprintf(tileEarthUrlTemplate, q, year, month, day, hour, minute/10, x, y)
	log.Printf("get tile %d-%d from %s for earth", x, y, url)
	return requestImage(url)
}

func getTileShorelines(q Quality, c ShorelineColor, x int, y int) (image.Image, error) {
	url := fmt.Sprintf(tileShorelineUrlTemplate, q, c, x, y)
	log.Printf("get tile %d-%d from %s for shorelines", x, y, url)
	return requestImage(url)
}

func requestImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
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
